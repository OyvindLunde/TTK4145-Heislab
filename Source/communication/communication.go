package communication

import (
	"time"

	"../elevio"
	"../logmanagement"
	"../network"
)

/*Broadcast and recieve channel*/
type NetworkChannels struct {
	RcvChannel   chan logmanagement.Elev
	BcastChannel chan logmanagement.Elev
}

func Transmit(port int, channels NetworkChannels) {
	go SendMyElevInfo(channels.BcastChannel)
	go network.BrodcastMessage(port, channels.BcastChannel)
}

func Receive(port int, channels NetworkChannels, toggleLights chan elevio.PanelLight, newOrderChannel chan logmanagement.Order, resetChannel chan bool) {
	go UpdateFromNetwork(channels.RcvChannel, toggleLights, newOrderChannel, resetChannel)
	go network.RecieveMessage(port, channels.RcvChannel)
}

/*Sends MyElevInfo on channel in parameter*/
func SendMyElevInfo(BcastChannel chan logmanagement.Elev) { // Move to Coms
	for {
		time.Sleep(2 * time.Millisecond)
		BcastChannel <- logmanagement.GetMyElevInfo()
	}
}

/*Updates orderList and otherElevInfo based on the received message*/
func UpdateFromNetwork(RcvChannel chan logmanagement.Elev, lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- logmanagement.Order, resetChannel chan bool) {
	for {
		time.Sleep(2 * time.Millisecond)
		select {
		case a := <-RcvChannel:
			if a.Id != logmanagement.GetMyElevInfo().Id {
				if logmanagement.ShouldIReset(a) {
					resetChannel <- true
				}
				logmanagement.UpdateOtherElevInfo(a)
				logmanagement.UpdateOrderList(a, lightsChannel, newOrderChannel)
			}

		}
	}
}
