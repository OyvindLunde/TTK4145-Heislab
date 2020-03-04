package orderhandler

import (
	"../elevio"
	"../logmanagement"
)

// GetDestination returns the floor the elevator should go to
func GetDestination(order logmanagement.Order) int { //Overflødig per nå
	if order.Active != -1 {
		return order.Floor
	}
	return -1
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
func ShouldElevatorExecuteOrder(order logmanagement.Order, currentfloor int, destination int) bool {
	if destination == -1 || destination == order.Floor {
		return true
	}
	return false
}

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
