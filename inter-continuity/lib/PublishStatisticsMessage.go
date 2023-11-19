package lib

type StatisticFactor = int

const (
	SetTotalCount StatisticFactor = 1 + iota
	IncreaseCompleteCount
	IncreaseFailedCount
)

type PublishStatisticsMessage struct {
	Type  StatisticFactor
	Value int
}
