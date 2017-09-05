package movingstats_test

import (
	"testing"

	"math/rand"

	"os"
	"time"

	"fmt"

	"github.com/lagarciag/golang-moving-average"
	"github.com/lagarciag/movingstats"
)

func TestMain(m *testing.M) {
	seed := time.Now().UTC().UnixNano()
	fmt.Println("SEED:", seed)
	rand.Seed(seed)
	os.Exit(m.Run())
}

func TestSimpleMovingAverage(t *testing.T) {

	period := rand.Intn(10) + rand.Intn(1000000)
	//period = 2

	t.Log("period:", period)
	size := period + rand.Intn(1000000)
	//size = 5
	t.Log("size:", size)
	movingStats := movingstats.NewAverage(period)
	movingAverage := movingaverage.New(uint(period))

	floatList := make([]float64, size)

	for i := range floatList {

		floatList[i] = rand.Float64() * float64(rand.Intn(100000))
	}

	//floatList = []float64{1,1,1,2,2,2}

	for _, value := range floatList {
		movingStats.Add(value)
		movingAverage.Add(value)
	}

	avg1 := movingStats.SimpleMovingAverage()
	avg2 := movingAverage.Avg()

	if uint(avg1) != uint(avg2) {

		t.Error("Mistmatch: ", avg1, avg2)
	} else {
		t.Log("Match: ", avg1, avg2)
	}

}

func TestSimpleMovingAverageFromStats(t *testing.T) {

	period := rand.Intn(10) + rand.Intn(1000000)
	//period = 2

	t.Log("period:", period)
	size := period + rand.Intn(1000000)
	//size = 5
	t.Log("size:", size)
	movingStats := movingstats.NewMovingStats(period)
	movingAverage := movingaverage.New(uint(period))

	floatList := make([]float64, size)

	for i := range floatList {

		floatList[i] = rand.Float64() * float64(rand.Intn(100000))
	}

	//floatList = []float64{1,1,1,2,2,2}

	for _, value := range floatList {
		movingStats.Add(value)
		movingAverage.Add(value)

	}

	avg1 := movingStats.SimpleMovingAverage()
	avg2 := movingAverage.Avg()

	if uint(avg1) != uint(avg2) {

		t.Error("Mistmatch: ", avg1, avg2)
	} else {
		t.Log("Match: ", avg1, avg2)
	}

}

func TestStandardDeviation(t *testing.T) {

	for n := 0; n < 10; n++ {

		period := rand.Intn(10) + rand.Intn(100)
		//period = 10
		size := period + rand.Intn(100000)
		//size = 10
		movingStats := movingstats.NewMovingStats(period)
		movingAverage := movingaverage.New(uint(period))

		floatList := make([]float64, size)

		for i := range floatList {

			floatList[i] = rand.Float64()*10 + float64(rand.Intn(100000))
		}

		/*
			floatList = []float64{2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
				2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
				2, 1, 2, 1, 2, 2, 2, 1, 2, 1}
		*/
		for _, value := range floatList {
			movingStats.Add(value)
			movingAverage.Add(value)
		}

		std2 := movingAverage.StdDev()
		std3 := movingStats.MovingStandardDeviation()
		avg := movingStats.SimpleMovingAverage()

		error := 100 - (std2 / std3 * 100)

		if error > 5 {
			t.Log("AVG:", avg)
			t.Log("GOLD:", std2)
			t.Log("std3: ", std3, std3-std2, 100-(std2/std3*100))
			t.Error("Error is too high")

		}
	}
}
