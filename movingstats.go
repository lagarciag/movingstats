package movingstats

import "github.com/VividCortex/ewma"

type MovingStats struct {
	sma   *MovingAverage
	sma20 *MovingAverage
	sEma  ewma.MovingAverage
	dEma  ewma.MovingAverage
	tEma  ewma.MovingAverage

	//MACD
	emaMacd9 ewma.MovingAverage
	ema12    ewma.MovingAverage
	ema26    ewma.MovingAverage
	macd     float64
}

func NewMovingStats(size int) *MovingStats {
	window := float64(size)
	mstats := &MovingStats{}
	mstats.sma = NewAverage(size)
	mstats.sma20 = NewAverage(size * 20)
	mstats.sEma = ewma.NewMovingAverage(window)
	mstats.dEma = ewma.NewMovingAverage(window)
	mstats.tEma = ewma.NewMovingAverage(window)

	mstats.emaMacd9 = ewma.NewMovingAverage(window * 9)
	mstats.ema12 = ewma.NewMovingAverage(window * 12)
	mstats.ema26 = ewma.NewMovingAverage(window * 26)

	return mstats
}

func (ms *MovingStats) Add(value float64) {
	ms.sma.Add(value)
	ms.sma20.Add(value)
	ms.sEma.Add(value)
	ms.dEma.Add(value)
	ms.tEma.Add(value)

	ms.emaMacd9.Add(value)
	ms.ema12.Add(value)
	ms.ema26.Add(value)

}

func (ms *MovingStats) SimpleMovingAverage() float64 {
	return ms.sma.SimpleMovingAverage()
}

func (ms *MovingStats) MovingStandardDeviation() float64 {
	return ms.sma.MovingStandardDeviation()
}

func (ms *MovingStats) MovingStandardDeviation20() float64 {
	return ms.sma20.MovingStandardDeviation()
}
