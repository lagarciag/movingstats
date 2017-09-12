package movingstats

import (
	"sync"

	"github.com/VividCortex/ewma"
	"github.com/lagarciag/movingaverage"
	"github.com/lagarciag/ringbuffer"
)

type MovingStats struct {
	mu *sync.Mutex

	sma *movingaverage.MovingAverage

	smaLong *movingaverage.MovingAverage

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

	HistMostRecent float64
	HistOldest     float64
	HistNow        float64
}

func NewMovingStats(size int) *MovingStats {
	window := float64(size)
	ms := &MovingStats{}
	ms.mu = &sync.Mutex{}
	ms.sma = movingaverage.New(size)
	ms.smaLong = movingaverage.New(size * 20)

	ms.sEma = ewma.NewMovingAverage(window)
	ms.sEmaHistory = ringbuffer.NewBuffer(size)

	ms.dEma = ewma.NewMovingAverage(window)
	ms.dEmaHistory = ringbuffer.NewBuffer(size)

	ms.tEma = ewma.NewMovingAverage(window)
	ms.tEmaHistory = ringbuffer.NewBuffer(size)

	ms.emaMacd9 = ewma.NewMovingAverage(window * 9)

	ms.ema12 = ewma.NewMovingAverage(window * 12)
	ms.ema26 = ewma.NewMovingAverage(window * 26)

	return ms
}

func (ms *MovingStats) Add(value float64) {
	ms.mu.Lock()
	ms.sma.Add(value)
	ms.smaLong.Add(value)
	ms.sEma.Add(value)
	ms.HistNow = ms.sEma.Value()
	ms.sEmaHistory.Push(ms.HistNow)

	ms.HistMostRecent = ms.sEmaHistory.MostRecent()
	ms.HistOldest = ms.sEmaHistory.Oldest()

	ms.sEmaSlope = ms.sEmaHistory.MostRecent() - ms.sEmaHistory.Oldest()

	if ms.sEmaSlope > 0 {
		ms.sEmaUp = true
	} else {
		ms.sEmaUp = false
	}

	ms.dEma.Add(value)
	ms.dEmaHistory.Push(ms.sEma.Value())
	ms.dEmaSlope = ms.dEmaHistory.MostRecent() - ms.dEmaHistory.Oldest()
	if ms.dEmaSlope > 0 {
		ms.dEmaUp = true
	} else {
		ms.dEmaUp = false
	}

	ms.tEma.Add(value)
	ms.tEmaHistory.Push(ms.sEma.Value())

	ms.tEmaSlope = ms.tEmaHistory.MostRecent() - ms.tEmaHistory.Oldest()

	if ms.tEmaSlope > 0 {
		ms.tEmaUp = true
	} else {
		ms.tEmaUp = false
	}

	//ms.emaMacd9.Add(value)
	ms.ema12.Add(value)
	ms.ema26.Add(value)

	ms.macd = ms.ema12.Value() - ms.ema26.Value()
	ms.emaMacd9.Add(ms.macd)

	ms.macdDivergence = ms.macd - ms.emaMacd9.Value()

	ms.mu.Unlock()
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

func (ms *MovingStats) StdDevLong() float64 {
	return ms.smaLong.MovingStandardDeviation()
}
