package orderhandler

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module contain functions for handling local orders
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"

	//"reflect"
	"strconv"
	"strings"
	"time"

	"../elevcontroller"
	"../elevio"
	"../logmanagement"
)

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Advanced Getters
// ------------------------------------------------------------------------------------------------------------------------------------------------------

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
	for i := 0; i < logmanagement.GetNumButtons(); i++ {
		status := int(logmanagement.GetOrder(floor, i).Status)
		if status == 0 || status == logmanagement.GetMyElevInfo().Id {
			UpdateLocalOrders(floor, i, status, true, false)
		}
	} // Nye ordrer i samme etg som kommer inn mens dørene er åpne: Rekker vi å sende at de ordrene er fullført?
	elevcontroller.ElevStopAtFloor(floor)
	for i := 0; i < 3; i++ { // Ta inn numButtons ??ddd
		status := int(logmanagement.GetOrder(floor, i).Status)
		if status == 0 || status == logmanagement.GetMyElevInfo().Id {
			UpdateLocalOrders(floor, i, -1, false, false)
			light := elevio.PanelLight{Floor: floor, Button: elevio.ButtonType(i), Value: false}
			lightsChannel <- light
		}
	}
}

/*Trenger vi egt denne?*/
func HandleButtonEvents(ButtonPress chan elevio.ButtonEvent, lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- logmanagement.Order) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-ButtonPress:
			order := logmanagement.GetOrder(a.Floor, int(a.Button))
			if order.Status == -1 {
				UpdateLocalOrders(order.Floor, int(order.ButtonType), 0, false, false)

				if order.ButtonType == 2 || len(logmanagement.GetOtherElevInfo()) == 0 { // Update lights and newOrder only for CAB orders and for single elev state
					light := elevio.PanelLight{Floor: a.Floor, Button: a.Button, Value: true}
					lightsChannel <- light
					newOrderChannel <- order
				}
				logmanagement.SetDisplayUpdates(true)
			}
		}
	}
}

/* Updates the Local Orders*/
func UpdateLocalOrders(floor int, button int, active int, finished bool, confirm bool) {
	logmanagement.SetOrder(floor, button, active, finished, confirm)
	UpdateCabOrderBackup()
	logmanagement.SetDisplayUpdates(true)
}

/*Updates local lights acording to orders*/
func UpdateLightsV2(lightschannel chan elevio.PanelLight) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-lightschannel:
			elevio.SetButtonLamp(a.Button, a.Floor, a.Value)
		}
	}
}

func CheckForUnconfirmedOrders(lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- logmanagement.Order) {
	orderList := logmanagement.GetOrderList()
	for i := 0; i < logmanagement.GetNumFloors(); i++ {
		for j := 0; j < logmanagement.GetNumButtons()-1; j++ {
			if orderList[i][j].Status == 0 && orderList[i][j].Confirm == false {
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <- light
				order := logmanagement.Order{Floor: i, ButtonType: j, Status: 0, Finished: false}
				newOrderChannel <- order
			}
		}
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Cost Function
// ------------------------------------------------------------------------------------------------------------------------------------------------------

//Returns true if I order is valid
func IsOrderValid(currentOrder logmanagement.Order) bool {
	if logmanagement.GetOrder(currentOrder.Floor, int(currentOrder.ButtonType)).Status == 0 {
		return true
	}
	return false
}

func ShouldITakeOrder(myCurrentOrder logmanagement.Order) bool {
	if myCurrentOrder.ButtonType == 2 { // Check if Cab Order
		return true
	}

	time.Sleep(1000 * time.Millisecond)

	conflictElevs := make([]logmanagement.Elev, 0)

	for _, otherElev := range logmanagement.GetOtherElevInfo() {
		if myCurrentOrder.Floor == otherElev.CurrentOrder.Floor && myCurrentOrder.ButtonType == otherElev.CurrentOrder.ButtonType {
			if otherElev.State != -2 {
				conflictElevs = append(conflictElevs, otherElev)
			}

		}
	}
	if len(conflictElevs) > 0 {
		return solveConflict(myCurrentOrder, logmanagement.GetMyElevInfo(), conflictElevs)
	}

	return true
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Private Functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Returns true if this elevator should take order during conflict*/
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

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Backup of cab orders
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Updates a txt file of cab orders*/
func UpdateCabOrderBackup() {
	filename := "CabOrderBackup" + strconv.Itoa(logmanagement.GetMyElevInfo().Id) + ".txt"

	orders := logmanagement.GetOrderList()
	cabOrders := make([]int, 0)
	for _, row := range orders {
		cabOrders = append(cabOrders, row[2].Status)
	}

	// convert []int to string
	cabOrderString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(cabOrders)), ","), "[]")

	file, err := os.Create(filename)
	checkError(err)

	defer file.Close()

	_, err = file.WriteString(cabOrderString)

}

/*Reads a txt file of backed up cab orders*/
func ReadCabOrderBackup(lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- logmanagement.Order) {
	filename := "CabOrderBackup" + strconv.Itoa(logmanagement.GetMyElevInfo().Id) + ".txt"
	fmt.Println("Checking for existing cab orders")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("No backup found, creating new")
		UpdateCabOrderBackup()
		return
	}
	readCab := string(content)

	cabStringList := strings.Split(readCab, ",")
	cabIntList := []int{}

	// convert []string to []int
	for _, i := range cabStringList {
		j, err := strconv.Atoi(i)
		checkError(err)
		cabIntList = append(cabIntList, j)
	}

	// Add active orders first
	for floor, status := range cabIntList {
		if status != -1 {
			//fmt.Println("Found active order")
			order := logmanagement.Order{Floor: floor, ButtonType: 2, Status: 0, Finished: false}
			logmanagement.SetOrder(floor, 2, 0, false, false)
			light := elevio.PanelLight{Floor: floor, Button: 2, Value: true}
			lightsChannel <- light
			newOrderChannel <- order
		}
	}
	// Add remaining pending orders
	/*for floor, status := range cabIntList {
		if status == 0 {
			fmt.Println("Found pending order")
			order := logmanagement.Order{Floor: floor, ButtonType: 2, Status: 0, Finished: false}
			logmanagement.SetOrder(floor, 2, 0, false, false)
			light := elevio.PanelLight{Floor: floor, Button: 2, Value: true}
			lightsChannel <- light
			newOrderChannel <- order
		}
	}*/
	time.Sleep(1 * time.Second)

}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
