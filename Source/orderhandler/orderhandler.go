package orderhandler

import (
	//"fmt"
	"math"
	"time"

	"../elevcontroller"
	"../elevio"
	"../logmanagement"
)

// GetDestination returns the floor the elevator should go to
func InitOrderHandler(port int) { //Overflødig per nå
	//logmanagement.InitNetwork(port)
	logmanagement.InitializeQueue()
	logmanagement.InitializeElevInfo(1) // Finn en løsning for å sette ID
}

// GetDestination returns the floor the elevator should go to
func GetDestination(order logmanagement.Order) int {
	if order.Status == 0 { // != -1 ???
		return order.Floor
	}
	return -1
}

func GetPendingOrder() logmanagement.Order {
	numFloors, numButtons := logmanagement.GetMatrixDimensions()
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			if logmanagement.OrderQueue[i][j].Status == 0 {
				return logmanagement.OrderQueue[i][j]
			}
		}
	}
	return logmanagement.Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
}

//Annet navn enn Get?
// GetMotorDirection returns which direction the elevator should move
func GetDirection(currentfloor int, destination int) int {
	if destination == -1 || destination == currentfloor {
		return 0
	} else if destination > currentfloor {
		return 1
	} else {
		return -1
	}
}

// ShouldElevatorStop determines if the elevator should stop when it reaches a floor
func ShouldElevatorStop(currentfloor int, destination int, elev logmanagement.Elev, elevlist []logmanagement.Elev) bool {
	dir := GetDirection(currentfloor, destination)
	if dir == 0 {
		return true
	}
	if logmanagement.GetOrder(currentfloor, 2).Status == 0 {
		// Update order queue?
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

func StopAtFloor(floor int) {
	for i := 0; i < 3; i++ {
		if logmanagement.OrderQueue[floor][i].Status != -1 {
			logmanagement.OrderQueue[floor][i].Finished = true
		}
	}
	elevcontroller.ElevStopAtFloor(floor)
	for i := 0; i < 3; i++ { // Ta inn numButtons ??ddd
		UpdateOrderQueue(floor, i, -1)
		//elevio.SetButtonLamp(elevio.ButtonType(i), floor, false)
	}
	logmanagement.OrderQueue[floor][0].Finished = false
	logmanagement.OrderQueue[floor][1].Finished = false
	logmanagement.OrderQueue[floor][2].Finished = false
}

func HandleButtonEvents(ButtonPress chan elevio.ButtonEvent) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-ButtonPress:
			//fmt.Println(a)
			order := logmanagement.GetOrder(a.Floor, int(a.Button))
			if order.Status == -1 {
				UpdateOrderQueue(order.Floor, int(order.ButtonType), 0)
				elevio.SetButtonLamp(a.Button, a.Floor, true)
				logmanagement.Updates = true
			}

		}
	}
}

func ShouldITakeOrder(order logmanagement.Order, elev logmanagement.Elev, destination int, elevlist []logmanagement.Elev) bool {
	//fmt.Printf("Elev: %v\n", elev)
	//fmt.Printf("Elevlist: %v\n", elevlist)
	//fmt.Printf("Destination: %v\n", destination)
	if destination == -1 || order.Status == -1 {
		return false
	}
	/*fmt.Println("In: Should i Take Order")
	fmt.Println(elev)
	fmt.Println(elevlist)*/
	conflictElevs := []logmanagement.Elev{}
	for _, elev := range elevlist {
		if elev.State == 0 {
			conflictElevs = append(conflictElevs, elev)
		}
		/*_, _, currentOrder, _ := logmanagement.GetElevInfo(elev)
		if order == currentOrder {
			conflictElevs = append(conflictElevs, elev)
		}*/
	}
	if len(conflictElevs) != 0 {
		return solveConflict(order, elev, destination, conflictElevs)
	}
	return true
}

func solveConflict(order logmanagement.Order, elev logmanagement.Elev, destination int, conflictElevs []logmanagement.Elev) bool {
	id, floor, _, _ := logmanagement.GetElevInfo(elev)
	myDist := math.Abs(float64(floor - destination))
	for _, conflictElev := range conflictElevs {
		conflictID, conflictfloor, _, _ := logmanagement.GetElevInfo(conflictElev)
		if myDist > math.Abs(float64(conflictfloor-destination)) {
			return false
		} else if myDist == math.Abs(float64(conflictfloor-destination)) && id > conflictID {
			return false
		}
	}
	return true
}

func GetElevList() []logmanagement.Elev {
	return logmanagement.ElevList
}

// UpdateOrderQueue updates the order queue
func UpdateOrderQueue(floor int, button int, active int) { //Må kanskje endre til active OrderStatus
	logmanagement.OrderQueue[floor][button].Status = logmanagement.OrderStatus(active)
	logmanagement.Updates = true
}

func UpdateLights(numFloors int, numButtons int) {
	for {
		time.Sleep(20 * time.Millisecond)
		for i := 0; i < numFloors; i++ {
			for j := 0; j < numButtons; j++ {
				if logmanagement.OrderQueue[i][j].Status == -1 {
					elevio.SetButtonLamp(elevio.ButtonType(j), i, false)
				} else {
					elevio.SetButtonLamp(elevio.ButtonType(j), i, true)
				}
			}
		}
	}

}
