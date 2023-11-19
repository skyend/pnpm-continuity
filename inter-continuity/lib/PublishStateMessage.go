package lib

type PublishingState = int

const (
	PublishingStateStart PublishingState = 1 + iota
	PublishingStatePublishing
	PublishingStateFinish
	PublishingStateFailed
)

type PublishStateMessage struct {
	Name    string
	Version string

	State PublishingState

	NpmErrorCode string
}
