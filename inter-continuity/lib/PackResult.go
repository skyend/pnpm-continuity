package lib

import "fmt"

type PackResult struct {
	Err                 error
	Success             bool
	InputPackageName    string
	InputPackageVersion string

	Filename     string
	ShaSum       string
	PackageSize  string
	UnpackedSize string
	Integrity    string
	TotalFiles   int
}

func PrintFailedPacks(packResults []PackResult) string {
	printing := ""
	printing += fmt.Sprintf("# FailedPackages\n")
	for _, packResult := range packResults {
		if !packResult.Success {
			printing += fmt.Sprintf("[%s@%s]\n", packResult.InputPackageName, packResult.InputPackageVersion)
			printing += fmt.Sprintf("-> %s\n", packResult.Err)
		}
	}
	printing += "\n"

	return printing
}

type PackResultAggregate struct {
	Total   int
	Failed  int
	Success int
}

func (p PackResultAggregate) Print() string {
	return fmt.Sprintf("# Mirroring Result \n") +
		fmt.Sprintf("total: %d \n", p.Total) +
		fmt.Sprintf("success: %d \n", p.Success) +
		fmt.Sprintf("failed: %d \n", p.Failed) +
		fmt.Sprintf("\n")
}

func AggregatePackResults(packResults []PackResult) PackResultAggregate {
	agg := PackResultAggregate{
		Total:   0,
		Failed:  0,
		Success: 0,
	}

	for _, packResult := range packResults {
		agg.Total++
		if packResult.Success {
			agg.Success++
		} else {
			agg.Failed++
		}
	}

	return agg
}
