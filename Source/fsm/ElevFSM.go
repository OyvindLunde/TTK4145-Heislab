package fsm

import (
	"../elevcontroller"
	"../elevio"
	"../logmanagement"
	"../orderhandler"
)

type State string

const (
	INIT    = "INIT"
	IDLE    = "IDLE"
	EXECUTE = "EXECUTE"
	RESET   = "RESET"
)

func Initialize(numFloors int) {
	elevio.Init("localhost:15657", numFloors)
	elevcontroller.InitializeElevator(numFloors)
	elevio.SetFloorIndicator(0) // Evt fix løpende update senere
	orderhandler.InitOrderHandler(15647)

}

func RunElevator() {
	destination := -1
	dir := 0   // declared to make code run, prob might delete later
	floor := 0 // same
	state := IDLE
	//id := 1
	NoOrder := logmanagement.Order{Floor: -1, ButtonType: -1, Active: -1}
	orderhandler.UpdateElevInfo(floor, NoOrder, 0)

	ButtonPress := make(chan elevcontroller.Button)
	FloorReached := make(chan int)

	go elevcontroller.ButtonPressed(ButtonPress)
	go elevcontroller.FloorIsReached(FloorReached)
	//go elevio.PollButtons(ButtonPressed)
	//go elevio.PollFloorSensor(FloorReached)
	//go orderhandler.HandleButtonEvents(&floor)

	for {
		switch state {
		case IDLE:
			//fmt.Println("In: IDLE")
			currentOrder := orderhandler.GetPendingOrder()
			if currentOrder != NoOrder {
				currentOrder.Active = 1
				orderhandler.UpdateElevInfo(floor, currentOrder, 0)     //Skrive IDLE og ikke 0??
				destination = orderhandler.GetDestination(currentOrder) //tentativt navn
				ElevList := orderhandler.GetElevList()                  // Må lage denne
				if orderhandler.ShouldITakeOrder(currentOrder, logmanagement.ElevInfo, destination, ElevList) {
					orderhandler.UpdateOrderQueue(currentOrder.Floor, int(currentOrder.ButtonType), 1)
					orderhandler.UpdateElevInfo(floor, currentOrder, 1)
					dir = orderhandler.GetDirection(floor, destination)
					state = EXECUTE
				} else {
					currentOrder = NoOrder
					orderhandler.UpdateElevInfo(floor, currentOrder, 0)
				}

			}

		case EXECUTE:
			elevio.SetMotorDirection(elevio.MotorDirection(dir))
			if dir != 0 && orderhandler.ShouldElevatorStop(floor, destination, logmanagement.ElevInfo, logmanagement.ElevList) {
				elevcontroller.ElevStopAtFloor(floor)
				orderhandler.ClearOrdersAtFloor(floor)
				dir = orderhandler.GetDirection(floor, destination)
				elevio.SetMotorDirection(elevio.MotorDirection(dir))
				if dir == 0 {
					destination = -1
					state = IDLE
				}
			}
			if dir == 0 {
				elevcontroller.OpenCloseDoor(3)
				orderhandler.ClearOrdersAtFloor(floor)
				state = IDLE
			}

		case RESET:
			//reset elevator

		}

	}

}

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
