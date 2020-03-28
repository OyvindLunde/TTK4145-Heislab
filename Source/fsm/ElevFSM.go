package fsm

import (
	"fmt"
	"time"

	"../elevcontroller"
	"../elevio"
	"../logmanagement"
	"../orderhandler"
)

type State int

const (
	INIT    = 0
	IDLE    = 1
	EXECUTE = 2
	LOST    = 3
	RESET   = 4
)

type FsmChannels struct {
	ButtonPress  chan elevio.ButtonEvent
	FloorReached chan int
}

func Initialize(numFloors int, id int, addr int) {
	elevcontroller.InitializeElevator(numFloors, addr)
	elevio.SetFloorIndicator(0)
	//logmanagement.InitializeElevInfo()
	orderhandler.InitOrderHandler(id)
}

func RunElevator(channels FsmChannels, numFloors int, numButtons int) {
	fmt.Println("Hello")
	//destination := -1
	dir := 0
	floor := 0
	state := IDLE
	NoOrder := logmanagement.Order{Floor: -1, ButtonType: -1, Status: 2, Finished: false}
	currentOrder := NoOrder

	go elevio.PollButtons(channels.ButtonPress) // Kan vi legge denne inn i HandleButtonEvents?
	go elevio.PollFloorSensor(channels.FloorReached)
	go orderhandler.HandleButtonEvents(channels.ButtonPress)
	go orderhandler.UpdateLights(numFloors, numButtons)
	//go logmanagement.UpdateElevInfo(&floor, &currentOrder, &state) // Vurdere å droppe denne? Kjører unødvendig ofte

	for {
		time.Sleep(20 * time.Millisecond)
		/*if len(logmanagement.OtherElevInfo) > 0 {
			logmanagement.PrintOrderQueue(logmanagement.OtherElevInfo[0].Orders)
		}*/
		switch state {
		case IDLE:
			currentOrder = orderhandler.GetPendingOrder()
			if currentOrder != NoOrder {
				//destination = orderhandler.GetDestination(currentOrder)
				// currentOrder.Status = 1 // Tror denne linjen er kilden til kommunikasjonsproblemet vårt
				ElevList := orderhandler.GetElevList() // ElevList er public så trenger egt ikke denne?
				if orderhandler.ShouldITakeOrder(currentOrder, logmanagement.MyElevInfo, currentOrder.Floor, ElevList) {
					currentOrder.Status = 1
					orderhandler.UpdateOrderQueue(currentOrder.Floor, int(currentOrder.ButtonType), 1, false)
					dir = orderhandler.GetDirection(floor, currentOrder.Floor)
					state = EXECUTE
					logmanagement.UpdateElevInfo(floor, currentOrder, state)
				}
			}

		case EXECUTE:
			elevio.SetMotorDirection(elevio.MotorDirection(dir)) // Blir kalt (unødvendig) mange ganger. Men er "sikker"
			select {
			case a := <-channels.FloorReached:
				floor = a
				logmanagement.UpdateElevInfo(floor, currentOrder, state)
				elevio.SetFloorIndicator(floor)
				if orderhandler.ShouldElevatorStop(floor, currentOrder.Floor, logmanagement.MyElevInfo, logmanagement.OtherElevInfo) {
					orderhandler.StopAtFloor(floor)
					dir = orderhandler.GetDirection(floor, currentOrder.Floor)
					elevio.SetMotorDirection(elevio.MotorDirection(dir))
					if dir == 0 { // Forslag: Legge inn en CheckForCABOrders funksjon, må i så fall inn i default også
						//destination = -1 // Unødvendig?
						state = IDLE
						logmanagement.UpdateElevInfo(floor, NoOrder, state)
					}
				}

			default:
				if dir == 0 {
					orderhandler.StopAtFloor(floor)
					state = IDLE
					logmanagement.UpdateElevInfo(floor, NoOrder, state)
				}
			}

		case RESET:
			//reset elevator

		}
	}
}
