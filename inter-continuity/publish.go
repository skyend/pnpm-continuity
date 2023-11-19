package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"math"
	"math/rand"
	"os"
	"path"
	"pnpm-inter-continuity/inter-continuity/lib"
	"regexp"
	"strings"
	"sync"
)

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
	num      int
	spinner  spinner.Model

	state              int
	total              int
	failed             int
	completed          int
	publishingPackages []lib.PublishStateMessage
	failedPackages     []lib.PublishStateMessage
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),

		spinner: s,
	}
}
func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return m.spinner.Tick
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.num = rand.Int()
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		}
		break
	case lib.PublishStatisticsMessage:
		if msg.Type == lib.SetTotalCount {
			m.total = msg.Value
		}
		break
	case lib.PublishStateMessage:
		if msg.State == lib.PublishingStateStart {
			//

			m.publishingPackages = append(m.publishingPackages, msg)
		} else if msg.State == lib.PublishingStateFinish {
			//
			m.publishingPackages = lib.Filter(m.publishingPackages, func(message lib.PublishStateMessage) bool {
				return message.Name != msg.Name
			})
			m.completed++
		} else if msg.State == lib.PublishingStateFailed {
			//
			m.publishingPackages = lib.Filter(m.publishingPackages, func(message lib.PublishStateMessage) bool {
				return message.Name != msg.Name
			})
			m.failedPackages = append(m.failedPackages, msg)
			m.failed++
		}
		break
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
func (m model) View() string {
	// The header
	s := fmt.Sprintf("Publishing... %s \n", m.spinner.View())
	s += fmt.Sprintf("Progress: %d / %d, Failed: %d \n", m.completed, m.total, m.failed)

	recentFailedPackages := lib.LastSlice(m.failedPackages, 5)

	computedRowCount := int(math.Max(
		float64(len(m.publishingPackages)),
		float64(len(recentFailedPackages)),
	))
	rows := make([][]string, computedRowCount)
	for i := 0; i < computedRowCount; i++ {
		firstColumn := ""
		secondColumn := ""
		thirdColumn := ""

		if len(m.publishingPackages) > i {
			publishingItem := m.publishingPackages[i]
			firstColumn = publishingItem.Name
		}

		if len(recentFailedPackages) > i {
			publishingItem := recentFailedPackages[i]
			secondColumn = fmt.Sprintf(
				"%s (%s)",
				publishingItem.Name,
				publishingItem.NpmErrorCode,
			)
		}

		rows[i] = []string{firstColumn, secondColumn, thirdColumn}
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		//StyleFunc(func(row, col int) lipgloss.Style {
		//	switch {
		//	case row == 0:
		//		return HeaderStyle
		//	case row%2 == 0:
		//		return EvenRowStyle
		//	default:
		//		return OddRowStyle
		//	}
		//}).
		Headers("Publishing...", "Failed", "Note").
		Rows(rows...)

	// You can also add tables row-by-row
	s += "\n" + t.Render() + "\n"

	// Send the UI for rendering
	return s
}

const WorkerCount = 10

var RemoveTgzExtension = regexp.MustCompile("\\.tgz$")

func main() {
	publishQueue := make(chan lib.PackageTarball)

	p := tea.NewProgram(initialModel())

	publishedTarballDestination := path.Join(lib.PackDestination, "published")

	// published tarball destination
	os.Mkdir(publishedTarballDestination, os.ModePerm)

	go func() {
		dirEntries, err := os.ReadDir(lib.PackDestination)

		if err != nil {
			panic(err)
		}

		tarEntries := lib.Filter(dirEntries, func(entry os.DirEntry) bool {
			return strings.HasSuffix(entry.Name(), ".tgz")
		})
		//fmt.Println(dirEntries)

		//p.Send(lib.PublishStateMessage{
		//	Name:    "",
		//	Version: "",
		//	State:   0,
		//})

		p.Send(lib.PublishStatisticsMessage{
			Type:  lib.SetTotalCount,
			Value: len(tarEntries),
		})

		for _, dirEntry := range tarEntries {
			TarballName := RemoveTgzExtension.ReplaceAll([]byte(dirEntry.Name()), []byte(""))

			publishQueue <- lib.PackageTarball{
				Filename:               dirEntry.Name(),
				TarballName:            string(TarballName),
				WorkingDirRelativePath: path.Join(lib.PackDestination, dirEntry.Name()),
			}
		}

		close(publishQueue)

	}()

	workerWait := sync.WaitGroup{}
	workerWait.Add(WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		go func() {
			for {
				publishTarball, more := <-publishQueue
				if !more {
					workerWait.Done()
					return
				}

				p.Send(lib.PublishStateMessage{
					Name:    publishTarball.Filename,
					Version: "",
					State:   lib.PublishingStateStart,
				})
				p.Send(lib.PublishStateMessage{
					Name:    publishTarball.Filename,
					Version: "",
					State:   lib.PublishingStatePublishing,
				})
				//time.Sleep(time.Second * 2)
				ret := lib.PublishPackedTarball(publishTarball, 0)

				if ret.Error != nil {

					p.Send(lib.PublishStateMessage{
						Name:         publishTarball.Filename,
						Version:      "",
						State:        lib.PublishingStateFailed,
						NpmErrorCode: ret.NpmErrorCode,
					})

					lib.WriteAppend(
						"./publish-fail.log",
						fmt.Sprintf("::%s\n%s", publishTarball.Filename, ret.Message),
					)
				} else {
					os.Rename(
						publishTarball.WorkingDirRelativePath,
						path.Join(publishedTarballDestination, publishTarball.Filename),
					)
					p.Send(lib.PublishStateMessage{
						Name:    publishTarball.Filename,
						Version: "",
						State:   lib.PublishingStateFinish,
					})

				}

			}
		}()
	}

	go func() {

		workerWait.Wait()
		fmt.Println("Publishing process completed")
		p.Kill()
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	return

}
