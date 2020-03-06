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
	//logmanagement.InitNetwork(20021)
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
	//go logmanagement.SendLogFromLocal()
	//go logmanagement.UpdateLogFromNetwork()

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
				elevio.SetFloorIndicator(floor)
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
