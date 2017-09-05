package movingstats

import "math"

type MovingAverage struct {
	count  int
	period int

	avgSum      float64
	average     float64
	avgHistBuff *historyBuffer

	avg2Sum     float64
	variance    float64
	varHistBuff *historyBuffer

}

func NewAverage(period int) *MovingAverage {

	avg := &MovingAverage{}
	avg.period = period
	avg.avgHistBuff = NewHistoryBuffer(period)
	avg.varHistBuff = NewHistoryBuffer(period)
	return avg
}

func (avg *MovingAverage) Add(value float64) {
	avg.avg(value)
}


func (avg *MovingAverage) SimpleMovingAverage() float64 {
	return avg.average
}

func (avg *MovingAverage) MovingStandardDeviation() float64 {
	return math.Sqrt(avg.variance)
}

func (avg *MovingAverage) avg(value float64) {
	avg.count++

	lastAvgValue := avg.avgHistBuff.Oldest()
	avg.avgSum = (avg.avgSum - lastAvgValue) + value

	if avg.count < avg.period {
		avg.average = avg.avgSum / float64(avg.count)
	} else {
		avg.average = avg.avgSum / float64(avg.period)
	}

	avg.avgHistBuff.Add(value)

	value2 := value * value
	last2AvgValue := avg.varHistBuff.Oldest()
	avg.avg2Sum = (avg.avg2Sum - last2AvgValue) + value2

	n := float64(avg.period)
	if avg.count < avg.period {
		n = float64(avg.count)
	}

	avg.variance = ((n * avg.avg2Sum) - (avg.avgSum * avg.avgSum)) / (n * (n - 1))

	avg.varHistBuff.Add(value2)
}

type historyBuffer struct {
	buff []float64
	head int
	tail int
	size int
}

func NewHistoryBuffer(size int) *historyBuffer {

	hb := &historyBuffer{}
	hb.size = size
	hb.buff = make([]float64, hb.size+1)
	hb.head = 0
	hb.tail = 1
	return hb
}

func (hb *historyBuffer) Add(value float64) {
	hb.buff[hb.head] = value
	hb.head++
	hb.tail++

	if hb.tail%(hb.size+1) == 0 {
		hb.tail = 0
	}

	if hb.head%(hb.size+1) == 0 {
		hb.head = 0
	}

}

func (hb *historyBuffer) Oldest() float64 {
	return hb.buff[hb.tail]
}

func (hb *historyBuffer) Last() float64 {
	if hb.head == 0 {
		return hb.buff[hb.size-1]
	}
	return hb.buff[hb.head-1]

}
