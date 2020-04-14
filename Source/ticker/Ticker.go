package ticker

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module handles the functions regarding the ticker
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"time"

	"../logmanagement"
)

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------------------------------------------------------------

var done chan bool
var ticker *time.Ticker

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Public functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------
/*Starts ticker and check if the other elevators finishes orders within ticklength * tickTreshold seconds*/
func StartTicker(tickLength time.Duration, tickTreshold int) {
	done = make(chan bool)
	ticker = time.NewTicker(tickLength * time.Second)
	go checkOnOtherElevs(tickTreshold)

}

/*Stops ticker*/
func StoppTicker() {
	ticker.Stop()
	done <- true

}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Private functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*checks if the other elevators finishes orders within ticklength * tickTreshold seconds*/
func checkOnOtherElevs(tickTreshold int) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			for i := 0; i < len(logmanagement.GetElevTickerInfo()); i++ {
				if logmanagement.GetOtherElevInfo()[i].CurrentOrder.Status != -1 && logmanagement.GetOtherElevInfo()[i].CurrentOrder.Status != 0 {
					logmanagement.IncrementElevTickerInfo(i)
					if logmanagement.GetElevTickerInfo()[i] >= tickTreshold {
						fmt.Println("Timer interupt")
						logmanagement.GetOtherElevInfo()[i].State = -2
						var floor = logmanagement.GetOtherElevInfo()[i].CurrentOrder.Floor
						var button = logmanagement.GetOtherElevInfo()[i].CurrentOrder.ButtonType
						logmanagement.SetOrder(floor, button, 0, false, false)
						logmanagement.SetLocalOrderStatus(floor, button, -2)
						//logmanagement.GetOtherElevInfo()[i].CurrentOrder.Status = -2
					}
				}
			}
		}
	}
}
