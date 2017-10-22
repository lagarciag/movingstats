package movingstats

import (
	"math"

	"github.com/sirupsen/logrus"
)

/*DmiCal Calculates all DMI components:

 ---------------------------------
 Directional Movement Calcualtion
 ---------------------------------

	To calculate the +DI and -DI you need to find the +DM and -DM (Directional Movement).
	+DM and -DM are calculated using the High, Low and Close for each period.
	You can then calculate the following:

	Current High - Previous High = UpMove
	Current Low - Previous Low = DownMove

	If UpMove > DownMove and UpMove > 0, then +DM = UpMove, else +DM = 0
	If DownMove > Upmove and Downmove > 0, then -DM = DownMove, else -DM = 0
*/
func (ms *MovingStats) dmiCalc() {

	currentHigh := ms.currentWindowHistory.High()
	previousHigh := ms.lastWindowHistory.High()
	currentLow := ms.currentWindowHistory.Low()
	previousLow := ms.lastWindowHistory.Low()

	upMove := currentHigh - previousHigh
	downMove := currentLow - previousLow

	if (upMove > downMove) && (upMove > float64(0)) {
		ms.plusDM = upMove
	} else {
		ms.plusDM = float64(0)
	}

	if (downMove > upMove) && (downMove > float64(0)) {
		ms.minusDM = downMove
	} else {
		ms.minusDM = float64(0)
	}

	ms.plusDMAvr.Add(ms.plusDM / ms.atr.Value())
	ms.minusDMAvr.Add(ms.minusDM / ms.atr.Value())

	//logrus.Debug(ms.plusDMAvr.Value(), ms.minusDMAvr.Value())

	ms.plusDI = ms.plusDMAvr.Value() * float64(100)
	ms.minusDI = ms.minusDMAvr.Value() * float64(100)

	//logrus.Debug(currentHigh, previousHigh, currentLow, previousLow, upMove, downMove, ms.plusDM, ms.minusDM, ms.atr.Value(), ms.plusDI, ms.minusDI)

	ms.adxAvr.Add(math.Abs((ms.plusDI - ms.minusDI) / (ms.plusDI + ms.minusDI)))
	ms.adx = float64(100) * ms.adxAvr.Value()

}
