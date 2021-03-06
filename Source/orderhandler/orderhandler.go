package orderhandler

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module contain functions for handling orders
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
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

func ShouldElevatorStop(currentfloor int, destination int) bool {
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

/*Stops elevator and updates Orders and Lights*/
func StopAtFloor(floor int, lightsChannel chan<- elevio.PanelLight) {
	for i := 0; i < logmanagement.GetNumButtons(); i++ {
		status := int(logmanagement.GetOrder(floor, i).Status)
		if status == 0 || status == logmanagement.GetMyElevInfo().Id {
			UpdateOrder(floor, i, status, true, false)
		}
	}

	elevcontroller.ElevStopAtFloor(floor)

	for i := 0; i < logmanagement.GetNumButtons(); i++ {
		status := int(logmanagement.GetOrder(floor, i).Status)
		if status == 0 || status == logmanagement.GetMyElevInfo().Id || logmanagement.GetOrder(floor, i).Finished == true {
			UpdateOrder(floor, i, -1, false, false)
			light := elevio.PanelLight{Floor: floor, Button: elevio.ButtonType(i), Value: false}
			lightsChannel <- light
		}
	}
}

func HandleButtonEvents(ButtonPress chan elevio.ButtonEvent, lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- logmanagement.Order) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case btn := <-ButtonPress:
			order := logmanagement.GetOrder(btn.Floor, int(btn.Button))
			if order.Status == -1 {
				UpdateOrder(order.Floor, int(order.ButtonType), 0, false, false)
				if order.ButtonType == 2 || len(logmanagement.GetOtherElevInfo()) == 0 { // Update lights and newOrder without confirmation only for CAB orders and for single elev state
					light := elevio.PanelLight{Floor: btn.Floor, Button: btn.Button, Value: true}
					lightsChannel <- light
					newOrderChannel <- order
				}
				logmanagement.SetDisplayUpdates(true)
			}
		}
	}
}

func UpdateOrder(floor int, button int, active int, finished bool, confirm bool) {
	logmanagement.SetOrder(floor, button, active, finished, confirm)
	UpdateCabOrderBackup()
	logmanagement.SetDisplayUpdates(true)
}

func UpdateLights(lightschannel chan elevio.PanelLight) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-lightschannel:
			elevio.SetButtonLamp(a.Button, a.Floor, a.Value)
		}
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Cost Function
// ------------------------------------------------------------------------------------------------------------------------------------------------------

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

	time.Sleep(1000 * time.Millisecond) // Sleep to give the elevators sufficient time to synchronize CurrentOrder

	conflictElevs := detectConflictElevs(myCurrentOrder)
	if len(conflictElevs) > 0 {
		return solveConflict(myCurrentOrder, logmanagement.GetMyElevInfo(), conflictElevs)
	}

	return true
}

func detectConflictElevs(myCurrentOrder logmanagement.Order) []logmanagement.Elev {
	conflictElevs := make([]logmanagement.Elev, 0)

	for _, otherElev := range logmanagement.GetOtherElevInfo() {
		if myCurrentOrder.Floor == otherElev.CurrentOrder.Floor && myCurrentOrder.ButtonType == otherElev.CurrentOrder.ButtonType {
			if otherElev.State != -2 {
				conflictElevs = append(conflictElevs, otherElev)
			}
		}
	}

	return conflictElevs
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Private Functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Returns true if this elevator has the smallest distance to the order*/
func solveConflict(order logmanagement.Order, elev logmanagement.Elev, conflictElevs []logmanagement.Elev) bool {
	myId := logmanagement.GetMyElevInfo().Id
	floor := logmanagement.GetMyElevInfo().Floor
	myDist := math.Abs(float64(floor - order.Floor))
	for _, conflictElev := range conflictElevs {
		theirDist := math.Abs(float64(conflictElev.Floor - order.Floor))
		if myDist > theirDist {
			return false
		} else if myDist == theirDist && myId > conflictElev.Id {
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

	cabOrderString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(cabOrders)), ","), "[]") // convert []int to string

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

	for floor, status := range cabIntList {
		if status != -1 {
			order := logmanagement.Order{Floor: floor, ButtonType: 2, Status: 0, Finished: false}
			logmanagement.SetOrder(floor, 2, 0, false, false)
			light := elevio.PanelLight{Floor: floor, Button: 2, Value: true}
			lightsChannel <- light
			newOrderChannel <- order
		}
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
