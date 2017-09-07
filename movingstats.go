package movingstats

import (
	"github.com/VividCortex/ewma"
	"github.com/lagarciag/movingaverage"
	"github.com/lagarciag/ringbuffer"
)

type MovingStats struct {
	sma *movingaverage.MovingAverage

	sma20 *movingaverage.MovingAverage

	sEma        ewma.MovingAverage
	sEmaSlope   float64
	sEmaUp      bool
	sEmaHistory *ringbuffer.RingBuffer

	dEma        ewma.MovingAverage
	dEmaSlope   float64
	dEmaUp      bool
	dEmaHistory *ringbuffer.RingBuffer

	tEma        ewma.MovingAverage
	tEmaSlope   float64
	tEmaUp      bool
	tEmaHistory *ringbuffer.RingBuffer

	//MACD
	emaMacd9 ewma.MovingAverage

	ema12 ewma.MovingAverage
	ema26 ewma.MovingAverage

	macd           float64
	macdDivergence float64
}

func NewMovingStats(size int) *MovingStats {
	window := float64(size)
	ms := &MovingStats{}
	ms.sma = movingaverage.New(size)
	ms.sma20 = movingaverage.New(size * 20)

	ms.sEma = ewma.NewMovingAverage(window)
	ms.sEmaHistory = ringbuffer.NewBuffer(size)

	ms.dEma = ewma.NewMovingAverage(window)
	ms.dEmaHistory = ringbuffer.NewBuffer(size)

	ms.tEma = ewma.NewMovingAverage(window)
	ms.tEmaHistory = ringbuffer.NewBuffer(size)

	ms.emaMacd9 = ewma.NewMovingAverage(window * 9)

	ms.ema12 = ewma.NewMovingAverage(window * 12)
	ms.ema26 = ewma.NewMovingAverage(window * 26)

	ms.macd = ms.ema26.Value() - ms.ema12.Value()
	ms.macdDivergence = ms.macd - ms.emaMacd9.Value()

	return ms
}

func (ms *MovingStats) Add(value float64) {
	ms.sma.Add(value)
	ms.sma20.Add(value * 20)

	ms.sEma.Add(value)
	ms.sEmaHistory.Push(value)
	ms.sEmaSlope = ms.sEmaHistory.Head() - ms.sEmaHistory.Tail()
	if ms.sEmaSlope > 0 {
		ms.sEmaUp = true
	} else {
		ms.sEmaUp = false
	}

	ms.dEma.Add(value)
	ms.dEmaHistory.Push(value)
	ms.dEmaSlope = ms.dEmaHistory.Head() - ms.dEmaHistory.Tail()
	if ms.dEmaSlope > 0 {
		ms.dEmaUp = true
	} else {
		ms.dEmaUp = false
	}

	ms.tEma.Add(value)
	ms.tEmaHistory.Push(value)
	ms.tEmaSlope = ms.tEmaHistory.Head() - ms.tEmaHistory.Tail()
	if ms.tEmaSlope > 0 {
		ms.tEmaUp = true
	} else {
		ms.tEmaUp = false
	}

	ms.emaMacd9.Add(value * 9)
	ms.ema12.Add(value * 12)
	ms.ema26.Add(value * 26)

	ms.macd = ms.ema26.Value() - ms.ema12.Value()
	ms.emaMacd9.Add(ms.macd)
}

func (ms *MovingStats) Ema1() float64 {
	return ms.sEma.Value()
}

func (ms *MovingStats) Ema1Slope() float64 {
	return ms.sEmaSlope
}

func (ms *MovingStats) Ema1Up() bool {
	return ms.sEmaUp
}

func (ms *MovingStats) DoubleEma() float64 {
	return ms.dEma.Value()
}

func (ms *MovingStats) DoubleEmaSlope() float64 {
	return ms.dEmaSlope
}

func (ms *MovingStats) DoubleEmaUp() bool {
	return ms.dEmaUp
}

func (ms *MovingStats) TripleEma() float64 {
	return ms.tEma.Value()
}

func (ms *MovingStats) TripleEmaSlope() float64 {
	return ms.tEmaSlope
}

func (ms *MovingStats) TripleEmaUp() bool {
	return ms.tEmaUp
}

func (ms *MovingStats) Macd() float64 {
	return ms.macd
}

func (ms *MovingStats) EmaMacd9() float64 {
	return ms.emaMacd9.Value()
}

func (ms *MovingStats) MacdDiv() float64 {
	return ms.macdDivergence
}

func (ms *MovingStats) SMA1() float64 {
	return ms.sma.SimpleMovingAverage()
}

func (ms *MovingStats) StdDev1() float64 {
	return ms.sma.MovingStandardDeviation()
}

func (ms *MovingStats) StdDev20() float64 {
	return ms.sma20.MovingStandardDeviation()
}
