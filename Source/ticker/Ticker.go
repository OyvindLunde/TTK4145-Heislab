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
// VariablesS
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
			//fmt .Println("tick")
			for i := 0; i < len(logmanagement.GetOtherElevInfo()); i++ {
				if logmanagement.GetOtherElevInfo()[i].CurrentOrder.Status != -1 && logmanagement.GetOtherElevInfo()[i].CurrentOrder.Status != 0 {
					logmanagement.GetOtherElevInfo()[i].CurrentOrder.TimeTicks += 1
					if logmanagement.GetOtherElevInfo()[i].CurrentOrder.TimeTicks >= tickTreshold {
						fmt.Println("Timer interupt")
						for j := 0; j < len(logmanagement.GetOtherElevInfo()); j++ {
							logmanagement.GetOtherElevInfo()[j].Orders[logmanagement.GetOtherElevInfo()[i].CurrentOrder.Floor][logmanagement.GetOtherElevInfo()[i].CurrentOrder.ButtonType].Status = 0
						}
						logmanagement.GetOtherElevInfo()[i].CurrentOrder = logmanagement.Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
					}
				}
			}
		}
	}
}
