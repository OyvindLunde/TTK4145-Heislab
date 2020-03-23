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
	RESET   = 3
)

type FsmChannels struct {
	ButtonPress  chan elevio.ButtonEvent
	FloorReached chan int
}

func Initialize(numFloors int, port int) {
	elevcontroller.InitializeElevator(numFloors, port)
	elevio.SetFloorIndicator(0)
	//logmanagement.InitializeElevInfo()
	orderhandler.InitOrderHandler(port)
}

func RunElevator(channels FsmChannels, numFloors int, numButtons int) {
	fmt.Println("Hello")
	destination := -1
	dir := 0
	floor := 0
	state := IDLE
	NoOrder := logmanagement.Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
	currentOrder := NoOrder

	go elevio.PollButtons(channels.ButtonPress)
	go elevio.PollFloorSensor(channels.FloorReached)
	go orderhandler.HandleButtonEvents(channels.ButtonPress)
	go orderhandler.UpdateLights(numFloors, numButtons)
	go logmanagement.UpdateElevInfo(&floor, &currentOrder, &state)

	for {
		time.Sleep(20 * time.Millisecond)
		switch state {
		case IDLE:
			/*fmt.Println("Orders FSM: ")
			fmt.Println(logmanagement.OrderQueue)
			fmt.Println("...")*/
			//fmt.Println(logmanagement.OrderQueue)
			currentOrder = orderhandler.GetPendingOrder()
			//fmt.Println(currentOrder)
			if currentOrder != NoOrder {
				//fmt.Println("I got an order")
				destination = orderhandler.GetDestination(currentOrder)
				currentOrder.Status = 1
				ElevList := orderhandler.GetElevList()
				if orderhandler.ShouldITakeOrder(currentOrder, logmanagement.ElevInfo, destination, ElevList) {
					orderhandler.UpdateOrderQueue(currentOrder.Floor, int(currentOrder.ButtonType), 1)
					dir = orderhandler.GetDirection(floor, destination)
					state = EXECUTE
				} else {
					currentOrder = NoOrder
				}
			}

		case EXECUTE:
			elevio.SetMotorDirection(elevio.MotorDirection(dir))
			select {
			case a := <-channels.FloorReached:
				floor = a
				logmanagement.ElevInfo.Floor = a // Added temporarily to Display the correct floor
				logmanagement.Updates = true
				elevio.SetFloorIndicator(floor)
				if orderhandler.ShouldElevatorStop(floor, destination, logmanagement.ElevInfo, logmanagement.ElevList) {
					//elevcontroller.ElevStopAtFloor(floor)
					orderhandler.StopAtFloor(floor)
					logmanagement.Updates = true
					dir = orderhandler.GetDirection(floor, destination)
					elevio.SetMotorDirection(elevio.MotorDirection(dir))
					if dir == 0 {
						destination = -1
						state = IDLE
					}
				}
			default:
				if dir == 0 {
					//elevcontroller.OpenCloseDoor(3)
					orderhandler.StopAtFloor(floor)
					logmanagement.Updates = true
					state = IDLE
				}
			}

		case RESET:
			//reset elevator

		}
	}
}
