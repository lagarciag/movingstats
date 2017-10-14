package movingstats

import "math"

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

	if (upMove > downMove) && (upMove > 0) {
		ms.plusDM = upMove
	} else {
		ms.plusDM = 0
	}

	if (downMove > upMove) && (downMove > 0) {
		ms.minusDM = downMove
	} else {
		ms.minusDM = 0
	}

	ms.plusDMAvr.Add(ms.plusDM / ms.atr.Value())
	ms.minusDMAvr.Add(ms.minusDM / ms.atr.Value())

	ms.plusDI = ms.plusDMAvr.Value() * 100
	ms.minusDI = ms.minusDMAvr.Value() * 100

	//fmt.Println(currentHigh, previousHigh, currentLow, previousLow, upMove, downMove,ms.plusDM, ms.minusDM, ms.atr.Value(), ms.plusDI, ms.minusDI)

	ms.adxAvr.Add(math.Abs((ms.plusDI - ms.minusDI) / (ms.plusDI + ms.minusDI)))
	ms.adx = 100 * ms.adxAvr.Value()

}
