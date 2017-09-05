package movingstats

type MovingStats struct {
	movingAverage *MovingAverage
}

func NewMovingStats(size int) *MovingStats {
	mstats := &MovingStats{}
	mstats.movingAverage = NewAverage(size)

	return mstats
}

func (ms *MovingStats) Add(value float64) {
	ms.movingAverage.Add(value)
}

func (ms *MovingStats) SimpleMovingAverage() float64 {
	return ms.movingAverage.SimpleMovingAverage()
}

func (ms *MovingStats) MovingStandardDeviation() float64 {
	return ms.movingAverage.MovingStandardDeviation()
}
