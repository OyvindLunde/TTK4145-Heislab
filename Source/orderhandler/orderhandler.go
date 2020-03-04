package orderhandler

import (
	"../elevcontroller"
	"../elevio"
	"../logmanagement"
)

// GetDestination returns the floor the elevator should go to
func InitOrderHandler(port int) { //Overflødig per nå
	logmanagement.InitNetwork(port)
	logmanagement.InitializeQueue(logmanagement.OrderQueue)
}

// GetDestination returns the floor the elevator should go to
func GetDestination(order logmanagement.Order) int { //Overflødig per nå
	order = GetPendingOrder()
	if order.Active != -1 {
		return order.Floor
	}
	return -1
}

func GetPendingOrder() logmanagement.Order {
	numFloors, numButtons := logmanagement.GetMatrixDimensions()
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			if logmanagement.OrderQueue[i][j].Active == 0 {
				return logmanagement.OrderQueue[i][j]
			}
		}
	}
	return logmanagement.Order{Floor: -1, ButtonType: -1, Active: -1}
}

//Annet navn enn Get?
// GetMotorDirection returns which direction the elevator should move
func GetMotorDirection(currentfloor int, destination int) int {
	if destination == -1 || destination == currentfloor {
		return 0
	} else if destination > currentfloor {
		return 1
	} else {
		return -1
	}
}

// ShouldElevatorExecuteOrder determines if the elevator should execute a certain order
/*func ShouldElevatorExecuteOrder(order logmanagement.Order, currentfloor int, destination int) bool {
	if destination == -1 || destination == order.Floor {
		return true
	}
	return false
}*/

// ShouldElevatorStop determines if the elevator should stop when it reaches a floor
func ShouldElevatorStop(currentfloor int, destination int) bool {
	dir := GetMotorDirection(currentfloor, destination)
	if dir == 0 {
		return true
	}
	if logmanagement.GetOrder(currentfloor, 2).Active == 0 {
		// Update order queue?
		return true
	}
	if logmanagement.GetOrder(currentfloor, 1).Active == 0 && dir == -1 {
		return true
	}
	if logmanagement.GetOrder(currentfloor, 0).Active == 0 && dir == 1 {
		return true
	}

	return false
}

func ClearOrdersAtFloor(floor int) {
	for i := 0; i < 3; i++ { // Ta inn numButtons ??
		logmanagement.UpdateOrderQueue(floor, i, -1)
		elevio.SetButtonLamp(elevio.ButtonType(i), floor, false)
	}
}

func HandleButtonEvents() {
	ButtonPress := make(chan elevcontroller.Button)
	FloorReached := make(chan int)

	go elevcontroller.ButtonPressed(ButtonPress)
	go elevcontroller.FloorIsReached(FloorReached)

	for {
		select {
		case a := <-ButtonPress:
			order := logmanagement.GetOrder(a.Floor, a.Type)
			if order.Active == -1 {
				logmanagement.UpdateOrderQueue(order.Floor, int(order.ButtonType), 0)
				elevcontroller.UpdateLight(elevcontroller.Button{Floor: order.Floor, Type: int(order.ButtonType)}, true)
			}
		case a := <-FloorReached:

		}
	}
}

func ShouldElevatorExecuteOrder() int {
	order := logmanagement.GetPendingOrder()
	if order.Active == 0 {
		return order.Floor
	}
	return -1
}
