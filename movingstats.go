package movingstats

import (
	"sync"

	"github.com/VividCortex/ewma"
	"github.com/lagarciag/movingaverage"
	"github.com/lagarciag/ringbuffer"
)

type MovingStats struct {
	mu *sync.Mutex

	windowSize int

	currentWindowHistory *ringbuffer.RingBuffer
	lastWindowHistory    *ringbuffer.RingBuffer

	// Simple Moving Average
	sma *movingaverage.MovingAverage

	// True Range Average
	atr ewma.MovingAverage

	// Directional Movement Index
	plusDMAvr  ewma.MovingAverage
	minusDMAvr ewma.MovingAverage
	adxAvr     ewma.MovingAverage

	smaLong *movingaverage.MovingAverage

	/*
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
	*/

	sEma *emaContainer

	dEma *emaContainer

	tEma *emaContainer

	//MACD
	emaMacd9 ewma.MovingAverage

	ema12 ewma.MovingAverage
	ema26 ewma.MovingAverage

	macd           float64
	macdDivergence float64

	//Directional Movement
	cHigh float64
	cLow  float64
	pHigh float64
	pLow  float64

	plusDM  float64
	minusDM float64
	plusDI  float64
	minusDI float64
	adx     float64

	HistMostRecent float64
	HistOldest     float64
	HistNow        float64
	count          int
}

func NewMovingStats(size int) *MovingStats {
	window := float64(size)
	ms := &MovingStats{}
	ms.mu = &sync.Mutex{}
	ms.windowSize = size
	ms.currentWindowHistory = ringbuffer.NewBuffer(size, true)
	ms.lastWindowHistory = ringbuffer.NewBuffer(size, true)

	ms.sma = movingaverage.New(size)
	ms.atr = ewma.NewMovingAverage(window * 14)
	ms.plusDMAvr = ewma.NewMovingAverage(window * 14)
	ms.minusDMAvr = ewma.NewMovingAverage(window * 14)
	ms.adxAvr = ewma.NewMovingAverage(window * 14)

	ms.smaLong = movingaverage.New(size * 20)

	/*
		ms.sEma = ewma.NewMovingAverage(window)
		ms.sEmaHistory = ringbuffer.NewBuffer(size, false)

		ms.dEma = ewma.NewMovingAverage(window)
		ms.dEmaHistory = ringbuffer.NewBuffer(size, false)

		ms.tEma = ewma.NewMovingAverage(window)
		ms.tEmaHistory = ringbuffer.NewBuffer(size, false)
	*/
	ms.sEma = newEmaContainer(size, 1)
	ms.dEma = newEmaContainer(size, 1)
	ms.tEma = newEmaContainer(size, 1)

	ms.emaMacd9 = ewma.NewMovingAverage(window * 9)

	ms.ema12 = ewma.NewMovingAverage(window * 12)
	ms.ema26 = ewma.NewMovingAverage(window * 26)

	return ms
}

func (ms *MovingStats) Add(value float64) {
	ms.mu.Lock()

	ms.sma.Add(value)
	ms.smaLong.Add(value)

	// ------------------------------------------------
	// Calculate Multiple Exponential Moving Averages
	// ------------------------------------------------
	ms.emaCalc(value)

	// --------------------------
	// Calculate MACD indicator
	// --------------------------
	ms.macdCalc(value)

	// ------------------------------------------
	// Calculate True Range & Average True Range
	// ------------------------------------------
	ms.atrCalc(value)

	// -----------------------------------------
	// Calculate Directional Movement Indicator
	// -----------------------------------------
	ms.dmiCalc()

	ms.count++
	ms.mu.Unlock()
}

func (ms *MovingStats) Ema1() float64 {
	return ms.sEma.Ema.Value()
}

func (ms *MovingStats) Ema1Slope() float64 {
	return ms.sEma.EmaSlope
}

func (ms *MovingStats) Ema1Up() bool {
	return ms.sEma.EmaUp
}

func (ms *MovingStats) DoubleEma() float64 {
	return ms.dEma.Ema.Value()
}

func (ms *MovingStats) DoubleEmaSlope() float64 {
	return ms.dEma.EmaSlope
}

func (ms *MovingStats) DoubleEmaUp() bool {
	return ms.dEma.EmaUp
}

func (ms *MovingStats) TripleEma() float64 {
	return ms.tEma.Ema.Value()
}

func (ms *MovingStats) TripleEmaSlope() float64 {
	return ms.tEma.EmaSlope
}

func (ms *MovingStats) TripleEmaUp() bool {
	return ms.tEma.EmaUp
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

func (ms *MovingStats) Atr() float64 {
	return ms.atr.Value()
}

func (ms *MovingStats) Adx() float64 {
	return ms.adx
}

func (ms *MovingStats) PlusDI() float64 {
	return ms.plusDI
}

func (ms *MovingStats) MinusDI() float64 {
	return ms.minusDI
}

func (ms *MovingStats) PlusDM() float64 {
	return ms.plusDM
}

func (ms *MovingStats) MinusDM() float64 {
	return ms.minusDM
}

func (ms *MovingStats) CHigh() float64 {
	return ms.cHigh
}

func (ms *MovingStats) CLow() float64 {
	return ms.cLow
}

func (ms *MovingStats) PHigh() float64 {
	return ms.pHigh
}

func (ms *MovingStats) PLow() float64 {
	return ms.pLow
}
