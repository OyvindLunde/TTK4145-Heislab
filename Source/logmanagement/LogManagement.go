package logmanagement

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module contain functions for updating orders and statuses between elevators
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"time"

	"../elevio"
	"../ticker"
)

const numFloors = 4
const numButtons = 3

var myElevInfo Elev
var otherElevInfo []Elev

var displayUpdates = false // Used to know when to update the display

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Declaration of structs and Enums
// ------------------------------------------------------------------------------------------------------------------------------------------------------

type Order struct {
	Floor      int
	ButtonType int
	Status     int
	Finished   bool
	Confirm    bool
}

type OrderStatus int

const (
	OrderTimeout OrderStatus = -2
	Inactive                 = -1
	Pending                  = 0
	Active                   = 1
)

type Elev struct {
	Id           int
	Floor        int
	CurrentOrder Order
	State        int
	Orders       [numFloors][numButtons]Order
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Init functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

func InitLogManagement(id int, toggleLights chan elevio.PanelLight, newOrderChannel chan Order) {
	initializeMyElevInfo(id)
	go checkForTimeout(toggleLights, newOrderChannel)
}

func initializeMyElevInfo(id int) {
	myElevInfo.Id = id
	myElevInfo.Floor = 0
	myElevInfo.CurrentOrder = Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
	myElevInfo.State = 1
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			myElevInfo.Orders[i][j].Floor = i
			myElevInfo.Orders[i][j].ButtonType = j
			myElevInfo.Orders[i][j].Status = -1
			myElevInfo.Orders[i][j].Finished = false
			myElevInfo.Orders[i][j].Confirm = false
		}
	}
	fmt.Println("MyElev initialized")
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Setters and Getters
// ------------------------------------------------------------------------------------------------------------------------------------------------------

func GetOrder(floor int, buttonType int) Order {
	return myElevInfo.Orders[floor][buttonType]
}

func SetOrder(floor int, buttonType int, status int, finished bool, confirm bool) {
	myElevInfo.Orders[floor][buttonType].Status = status
	myElevInfo.Orders[floor][buttonType].Finished = finished
	myElevInfo.Orders[floor][buttonType].Confirm = confirm
}

func GetOrderList() [numFloors][numButtons]Order {
	return myElevInfo.Orders
}

func GetOtherElevInfo() []Elev {
	return otherElevInfo
}

func GetNumFloors() int {
	return numFloors
}

func GetNumButtons() int {
	return numButtons
}

func GetMyElevInfo() Elev {
	return myElevInfo
}

func SetMyElevInfo(floor int, order Order, state int) {
	myElevInfo.Floor = floor
	myElevInfo.CurrentOrder = order
	myElevInfo.State = state
	displayUpdates = true
}

func GetDisplayUpdates() bool {
	return displayUpdates
}

func SetDisplayUpdates(value bool) {
	displayUpdates = value
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Public functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

func UpdateMyElevInfo(floor int, order Order, state int) {
	myElevInfo.Floor = floor
	myElevInfo.CurrentOrder = order
	myElevInfo.State = state
	displayUpdates = true
}

func RemoveElevFromOtherElevInfo(i int) {
	if len(otherElevInfo) != 0 {
		copy(otherElevInfo[i:], otherElevInfo[i+1:])         // Shift a[i+1:] left one index.
		otherElevInfo = otherElevInfo[:len(otherElevInfo)-1] // Truncate slice.
	}
}

func UpdateOtherElevInfo(msg Elev) {
	shouldIUpdateDisplay := checkForRemoteUpdates(msg)
	for i := 0; i < len(otherElevInfo); i++ {
		if msg.Id == otherElevInfo[i].Id {
			ticker.ResetHeartBeat(msg.Id)
			if otherElevInfo[i].CurrentOrder != msg.CurrentOrder {
				ticker.ResetOrderTicker(msg.Id)
			}
			if !ShouldIReset(msg) {
				otherElevInfo[i].Floor = msg.Floor
				otherElevInfo[i].CurrentOrder = msg.CurrentOrder
				otherElevInfo[i].State = msg.State
				otherElevInfo[i].Orders = msg.Orders

				if shouldIUpdateDisplay {
					SetDisplayUpdates(true)
				}
			}

			return
		}
	}
	ticker.AddElevToTicker(msg.Id)
	otherElevInfo = append(otherElevInfo, msg)
	displayUpdates = true
}

func UpdateOrderList(msg Elev, lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- Order) {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons-1; j++ {
			if msg.Orders[i][j].Finished == true && myElevInfo.Orders[i][j].Status != -1 { // Order finished by other elev
				myElevInfo.Orders[i][j].Status = -1
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: false}
				lightsChannel <- light
			} else if msg.Orders[i][j].Status == 0 && myElevInfo.Orders[i][j].Status == -1 && msg.Orders[i][j].Finished == false { // New order received
				myElevInfo.Orders[i][j].Status = 0
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <- light
				newOrderChannel <- msg.Orders[i][j]
			} else if msg.Orders[i][j].Status == msg.Id && (myElevInfo.Orders[i][j].Status == 0 || myElevInfo.Orders[i][j].Status == -1) && msg.Orders[i][j].Finished == false { // Order taken by other elev
				myElevInfo.Orders[i][j].Status = msg.Id
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <- light
			} else if msg.Orders[i][j].Status == 0 && myElevInfo.Orders[i][j].Status == 0 && myElevInfo.Orders[i][j].Confirm == false { // Order confirmed by other elev
				myElevInfo.Orders[i][j].Confirm = true
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <- light
				newOrderChannel <- msg.Orders[i][j]
			}
		}
	}
}

func ShouldIReset(msg Elev) bool {
	var orderFloor = myElevInfo.CurrentOrder.Floor
	var orderButton = myElevInfo.CurrentOrder.ButtonType
	if orderFloor != -1 && orderButton != -1 {
		if msg.Orders[orderFloor][orderButton].Status == -2 && myElevInfo.State != -2 {
			myElevInfo.State = -2
			otherElevInfo = otherElevInfo[:0]
			ticker.ClearElevTickerInfo()
			fmt.Println("resetting")
			return true
		}
	}
	return false
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Private functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

// Used for Display: Checks for changes from other elevators
func checkForRemoteUpdates(msg Elev) bool {
	for i := 0; i < len(otherElevInfo); i++ {
		if msg.Id == otherElevInfo[i].Id {
			if msg.Floor != otherElevInfo[i].Floor {
				return true
			}
			if msg.State != otherElevInfo[i].State {
				return true
			}
			if msg.CurrentOrder != otherElevInfo[i].CurrentOrder {
				return true
			}
			for j := 0; j < numFloors; j++ {
				for k := 0; k < numButtons; k++ {
					if msg.Orders[j][k] != otherElevInfo[i].Orders[j][k] {
						return true
					}
				}
			}
		}
	}
	return false
}

func checkForUnconfirmedOrders(lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- Order) {
	orderList := myElevInfo.Orders
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons-1; j++ {
			if orderList[i][j].Status == 0 && orderList[i][j].Confirm == false {
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <- light
				order := Order{Floor: i, ButtonType: j, Status: 0, Finished: false}
				newOrderChannel <- order
			}
		}
	}
}

func checkForTimeout(lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- Order) {
	for {
		time.Sleep(20 * time.Millisecond)
		for i := 0; i < len(otherElevInfo); i++ {
			key := otherElevInfo[i].Id
			value := otherElevInfo[i]
			if value.CurrentOrder.Status == -1 {
				ticker.ResetOrderTicker(key)
			}
			if ticker.HasCurrentOrderTimedOut(key) || !ticker.IsElevAlive(key) {
				fmt.Println("Timer interrupt")

				var floor = value.CurrentOrder.Floor
				var button = value.CurrentOrder.ButtonType
				if floor != -1 && button != -1 {
					SetOrder(floor, button, -2, false, true)
					time.Sleep(1000 * time.Millisecond) // We have to give the other elevs time to realise that they've timed out
					if button != 2 {
						SetOrder(floor, button, 0, false, true)
						order := Order{Floor: floor, ButtonType: button, Status: 0, Finished: false}
						newOrderChannel <- order
					} else {
						SetOrder(floor, button, -1, false, true)
					}
				}

				RemoveElevFromOtherElevInfo(i)
				ticker.DeleteElevFromTicker(key)
				checkForUnconfirmedOrders(lightsChannel, newOrderChannel)
				SetDisplayUpdates(true)

			}

		}
	}
}
