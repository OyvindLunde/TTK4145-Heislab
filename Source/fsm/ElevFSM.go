package fsm


import "fmt"
import (
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

func Initialize(numFloors int) {
	elevcontroller.InitializeElevator(numFloors)
	elevio.SetFloorIndicator(0) // Evt fix løpende update senere
	logmanagement.InitializeElevInfo()
	orderhandler.InitOrderHandler(15647)
}

func RunElevator() {
	fmt.Println("Hello")
	destination := -1
	dir := 0
	floor := 0
	state := IDLE
	//id := 1
	NoOrder := logmanagement.Order{Floor: -1, ButtonType: -1, Active: -1}
	currentOrder := NoOrder

	ButtonPress := make(chan elevio.ButtonEvent)
	FloorReached := make(chan int)

	go elevio.PollButtons(ButtonPress)
	go elevio.PollFloorSensor(FloorReached)
	go orderhandler.HandleButtonEvents(ButtonPress)
	go logmanagement.UpdateElevInfo(&floor, &currentOrder, &state)

	for {
		switch state {			
		case IDLE:
			currentOrder = orderhandler.GetPendingOrder()
			if currentOrder != NoOrder {
				destination = orderhandler.GetDestination(currentOrder) //tentativt navn
				currentOrder.Active = 1
				ElevList := orderhandler.GetElevList()                  // Må lage denne
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
			case a := <-FloorReached:
				floor = a
				elevcontroller.UpdateFloorIndicator(floor)
				if orderhandler.ShouldElevatorStop(floor, destination, logmanagement.ElevInfo, logmanagement.ElevList) {					
					elevcontroller.ElevStopAtFloor(floor)
					orderhandler.ClearOrdersAtFloor(floor)
					dir = orderhandler.GetDirection(floor, destination)
					elevio.SetMotorDirection(elevio.MotorDirection(dir))
					if dir == 0 {
						destination = -1
						state = IDLE
					}
				}
			default:
				if dir == 0 {
					elevcontroller.OpenCloseDoor(3)
					orderhandler.ClearOrdersAtFloor(floor)
					state = IDLE
				}
			}						

		case RESET:
			//reset elevator

		}

	}

}

/*if dir != 0 && orderhandler.ShouldElevatorStop(floor, destination, logmanagement.ElevInfo, logmanagement.ElevList) {
	elevcontroller.ElevStopAtFloor(floor)
	orderhandler.ClearOrdersAtFloor(floor)
	dir = orderhandler.GetDirection(floor, destination)
	elevio.SetMotorDirection(elevio.MotorDirection(dir))
	if dir == 0 {
		destination = -1
		state = IDLE
	}
}*/


/*func RunElevator() {
	destination := -1
	dir := 0   // declared to make code run, prob might delete later
	floor := 0 // same
	state := IDLE

	ButtonPressed := make(chan elevio.ButtonEvent)
	FloorReached := make(chan int)

	go elevio.PollButtons(ButtonPressed)
	go elevio.PollFloorSensor(FloorReached)

	for {
		switch state {
		case IDLE:
			destination = orderhandler.ShouldElevatorExecuteOrder() //tentativt navn
			if destination != -1 {
				dir = orderhandler.GetMotorDirection(floor, destination)
				state = EXECUTE
			}
			select {
			case a := <-ButtonPressed:
				if logmanagement.GetOrder(a.Floor, int(a.Button)).Active == -1 {
					logmanagement.UpdateOrderQueue(a.Floor, int(a.Button), 0)
					elevio.SetButtonLamp(a.Button, a.Floor, true)
				}

			default:
				newOrder := logmanagement.GetPendingOrder()
				if newOrder.Active == 0 {
					logmanagement.UpdateOrderQueue(newOrder.Floor, newOrder.ButtonType, 1)
					destination = newOrder.Floor
					dir = orderhandler.GetMotorDirection(floor, destination)
					state = EXECUTE
				}
			}

		case EXECUTE:
			elevio.SetMotorDirection(elevio.MotorDirection(dir))
			select {
			case a := <-ButtonPressed:
				if logmanagement.GetOrder(a.Floor, int(a.Button)).Active == -1 {
					logmanagement.UpdateOrderQueue(a.Floor, int(a.Button), 0)
					elevio.SetButtonLamp(a.Button, a.Floor, true)
				}

			case floorReached := <-FloorReached:
				floor = floorReached
				elevio.SetFloorIndicator(floor)
				if orderhandler.ShouldElevatorStop(floor, destination) {
					elevcontroller.ElevStopAtFloor(floor)
					orderhandler.ClearOrdersAtFloor(floor)
					dir = orderhandler.GetMotorDirection(floor, destination)
					elevio.SetMotorDirection(elevio.MotorDirection(dir))
					if dir == 0 {
						destination = -1
						state = IDLE
					}
				}
			default:
				if dir == 0 {
					elevcontroller.OpenCloseDoor(3)
					orderhandler.ClearOrdersAtFloor(floor)
					state = IDLE
				}
			}

		case RESET:
			//reset elevator

		}

	}

} */
