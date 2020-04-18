package ticker

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module handles the functions regarding the ticker
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"time"
	"../logmanagement"
	"../elevio"
	"../orderhandler"
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
func StartTicker(tickLength time.Duration, tickTreshold int, lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- logmanagement.Order){
	done = make(chan bool)
	ticker = time.NewTicker(tickLength * time.Second)
	go checkOnOtherElevs(tickTreshold, lightsChannel, newOrderChannel)

}

/*Stops ticker*/
func StoppTicker(){
	ticker.Stop()
	done <- true

}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Private functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*checks if the other elevators finishes orders within ticklength * tickTreshold seconds*/
func checkOnOtherElevs(tickTreshold int, lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- logmanagement.Order) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			logmanagement.IncrementHeartBeat()
			for i := 0; i < len(logmanagement.GetElevTickerInfo()); i++ {
				
				if logmanagement.GetOtherElevInfo()[i].CurrentOrder.Status != -1 && logmanagement.GetOtherElevInfo()[i].CurrentOrder.Status != 0 {
					logmanagement.IncrementElevTickerInfo(i)
					if logmanagement.GetElevTickerInfo()[i] >= tickTreshold && len(logmanagement.GetElevTickerInfo()) != 0{
						fmt.Println("Timer interrupt")
						var floor = logmanagement.GetOtherElevInfo()[i].CurrentOrder.Floor
						var button = logmanagement.GetOtherElevInfo()[i].CurrentOrder.ButtonType
						logmanagement.SetOrder(floor, button, -2, false, true)
						time.Sleep(1 * time.Second)
						logmanagement.RemoveElevFromOtherElevInfo(i)
						logmanagement.RemoveElevFromelevTickerInfo(i)
						logmanagement.RemoveHeartbeat(i)
						logmanagement.SetOrder(floor, button, 0, false, true)
						order := logmanagement.Order{Floor: floor, ButtonType: button, Status: 0, Finished: false}
						newOrderChannel <- order
						orderhandler.CheckForUnconfirmedOrders(lightsChannel, newOrderChannel)
					}
				} else if logmanagement.GetHeartBeat(i) > 1 { //burde være dynamisk
					logmanagement.RemoveElevFromOtherElevInfo(i)
					logmanagement.RemoveElevFromelevTickerInfo(i)
					logmanagement.RemoveHeartbeat(i)
					orderhandler.CheckForUnconfirmedOrders(lightsChannel, newOrderChannel)
				}
			}
		}
	}
}