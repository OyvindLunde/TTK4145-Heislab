package orderhandler

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module contain functions for handeling local orders
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	//"fmt"
	"math"
	"time"

	"../elevcontroller"
	"../elevio"
	"../logmanagement"
)

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Advanced Getters
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Returns a pending order if one exsist, Othervise returns a false order*/
func GetPendingOrder() logmanagement.Order {
	numFloors := logmanagement.GetNumFloors()
	numButtons := logmanagement.GetNumButtons()
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			if logmanagement.GetOrder(i, j).Status == 0 {
				return logmanagement.GetOrder(i, j)
			}
		}
	}
	return logmanagement.Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
}

/* Returns which direction the elevator should move*/
func GetDirection(currentfloor int, destination int) int {
	if destination == -1 || destination == currentfloor {
		return 0
	} else if destination > currentfloor {
		return 1
	} else {
		return -1
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// OrderHandling
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/* Returns true if the elevator should stop when it reaches a floor*/
func ShouldElevatorStop(currentfloor int, destination int, elev logmanagement.Elev, elevlist []logmanagement.Elev) bool {
	dir := GetDirection(currentfloor, destination)
	if dir == 0 {
		return true
	}
	if logmanagement.GetOrder(currentfloor, 2).Status == 0 {
		return true
	}
	if logmanagement.GetOrder(currentfloor, 1).Status == 0 && dir == -1 {
		return true
	}
	if logmanagement.GetOrder(currentfloor, 0).Status == 0 && dir == 1 {
		return true
	}

	return false

}

/*Stops elevator and updates LocalOrders acording to floor in param*/
func StopAtFloor(floor int, lightsChannel chan<- elevio.PanelLight) {
	for i := 0; i < 3; i++ {
		status := int(logmanagement.GetOrder(floor, i).Status)
		if status != -1 {
			UpdateLocalOrders(floor, i, status, true)
		}
	} // Nye ordrer i samme etg som kommer inn mens dørene er åpne: Rekker vi å sende at de ordrene er fullført?
	elevcontroller.ElevStopAtFloor(floor)
	for i := 0; i < 3; i++ { // Ta inn numButtons ??ddd
		if logmanagement.GetOrder(floor, i).Status != -1 {
			UpdateLocalOrders(floor, i, -1, false)
			light := elevio.PanelLight{Floor: floor, Button: elevio.ButtonType(i), Value: false}
			lightsChannel <- light
		}
	}
}

/*Trenger vi egt denne?*/
func HandleButtonEvents(ButtonPress chan elevio.ButtonEvent, lightsChannel chan<- elevio.PanelLight) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-ButtonPress:
			//fmt.Println(a)
			order := logmanagement.GetOrder(a.Floor, int(a.Button))
			if order.Status == -1 {
				UpdateLocalOrders(order.Floor, int(order.ButtonType), 0, false)
				light := elevio.PanelLight{Floor: a.Floor, Button: a.Button, Value: true}
				lightsChannel <- light
				//elevio.SetButtonLamp(a.Button, a.Floor, true)
				logmanagement.SetDisplayUpdates(true)
			}
		}
	}
}

/* Updates the Local Orders*/
func UpdateLocalOrders(floor int, button int, active int, finished bool) {
	logmanagement.SetOrder(floor, button, active, finished)
	logmanagement.SetDisplayUpdates(true)
}

/*Updates local lights acording to orders*/
func UpdateLightsV2(lightschannel chan elevio.PanelLight) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-lightschannel:
			//fmt.Println(a)
			elevio.SetButtonLamp(a.Button, a.Floor, a.Value)
		}
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Cost Function
// ------------------------------------------------------------------------------------------------------------------------------------------------------

//Returns true if I shuold take order
func ShouldITakeOrder(order logmanagement.Order, elev logmanagement.Elev, elevlist []logmanagement.Elev) bool {
	if order.Floor == -1 {
		return false
	}

	if int(order.ButtonType) == 2 {
		return true
	}

	conflictElevs := []logmanagement.Elev{}
	for _, elev := range elevlist {
		if elev.State == 0 {
			conflictElevs = append(conflictElevs, elev)
		}
	}
	if len(conflictElevs) != 0 {
		return solveConflict(order, elev, conflictElevs)
	}
	return true
}



// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Private Functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Retruns true if this elevator should take order during conflict*/
func solveConflict(order logmanagement.Order, elev logmanagement.Elev, conflictElevs []logmanagement.Elev) bool {
	id, floor, _, _ := logmanagement.GetElevInfo(elev)
	myDist := math.Abs(float64(floor - order.Floor))
	for _, conflictElev := range conflictElevs {
		conflictID, conflictfloor, _, _ := logmanagement.GetElevInfo(conflictElev)
		if myDist > math.Abs(float64(conflictfloor-order.Floor)) {
			return false
		} else if myDist == math.Abs(float64(conflictfloor-order.Floor)) && id > conflictID {
			return false
		}
	}
	return true
}
