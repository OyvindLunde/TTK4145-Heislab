package logmanagement

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module contain functions for updating orders and statuses between elevators
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"strconv"
	"time"

	"../elevio"
	"../network"
)

const numFloors = 4
const numButtons = 3

var myElevInfo Elev
var otherElevInfo []Elev
var elevTickerInfo []int

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

/*Elevstruct for keeping info about ther elevs*/
type Elev struct {
	Id           int
	Floor        int
	CurrentOrder Order
	//Lastseen time
	State  int
	Orders [numFloors][numButtons]Order
}

/*Broadcast and recieve channel*/
type NetworkChannels struct {
	RcvChannel   chan Elev
	BcastChannel chan Elev
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

func SetLocalOrderStatus(floor int, button int, status int) {
	myElevInfo.Orders[floor][button].Status = status
}

func GetOrderList() [numFloors][numButtons]Order {
	return myElevInfo.Orders
}

func GetOtherElevInfo() []Elev {
	return otherElevInfo
}

func GetElevInfo(elev Elev) (id, floor int, currentOrder Order, state int) {
	return elev.Id, elev.Floor, elev.CurrentOrder, elev.State
}

func GetNumFloors() int {
	return numFloors
}

func GetNumButtons() int {
	return numButtons
}

func GetElevTickerInfo() []int {
	return elevTickerInfo
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
// Init functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Initializes LogManagement*/

func InitLogManagement(id int) {
	initializeMyElevInfo(id)
}

/*Inits network communication*/
func InitCommunication(port int, channels NetworkChannels, toggleLights chan elevio.PanelLight, newOrderChannel chan Order) {
	go network.RecieveMessage(port, channels.RcvChannel)
	go network.BrodcastMessage(port, channels.BcastChannel)
	go SendMyElevInfo(channels.BcastChannel)
	go UpdateFromNetwork(channels.RcvChannel, toggleLights, newOrderChannel)
	fmt.Printf("Network initialized\n")
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Additional public functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

func IncrementElevTickerInfo(elev int) {
	elevTickerInfo[elev] += 1
}

/*Sends MyElevInfo on channel in parameter*/
func SendMyElevInfo(BcastChannel chan Elev) {
	for {

		time.Sleep(2 * time.Millisecond)
		//fmt.Println(myElevInfo.Id)
		BcastChannel <- myElevInfo
	}
}

func checkForUpdates(msg Elev) bool {
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

/*Updates OtherElevLsit from channel in parameter*/
func UpdateFromNetwork(RcvChannel chan Elev, lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- Order) {
	for {
		time.Sleep(2 * time.Millisecond)
		select {
		case a := <-RcvChannel:
			if a.Id != myElevInfo.Id {
				updateOtherElevInfo(a)
				updateOrderList(a, lightsChannel, newOrderChannel)
			}
		}
	}
}

/*Updates MyElevInfo variable from params*/
func UpdateMyElevInfo(floor int, order Order, state int) {
	myElevInfo.Floor = floor
	myElevInfo.CurrentOrder = order
	myElevInfo.State = state
	displayUpdates = true
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Private functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Updates otherelevinfo with info about elev in param*/
func updateOtherElevInfo(msg Elev) {
	bool1 := checkForUpdates(msg)
	for i := 0; i < len(otherElevInfo); i++ {
		if msg.Id == otherElevInfo[i].Id {
			if otherElevInfo[i].CurrentOrder != msg.CurrentOrder {
				elevTickerInfo[i] = 0
			}
			if otherElevInfo[i].State != -2 {
				otherElevInfo[i].Floor = msg.Floor
				otherElevInfo[i].CurrentOrder = msg.CurrentOrder
				otherElevInfo[i].State = msg.State
				otherElevInfo[i].Orders = msg.Orders
				if bool1 {
					SetDisplayUpdates(true)
				}
			}

			return
		}
	}
	elevTickerInfo = append(elevTickerInfo, 0)
	otherElevInfo = append(otherElevInfo, msg)
	displayUpdates = true
}

/*Updates orderlist with data stored in elev-param*/
func updateOrderList(msg Elev, lightsChannel chan<- elevio.PanelLight, newOrderChannel chan<- Order) {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons-1; j++ {
			if msg.Orders[i][j].Finished == true && myElevInfo.Orders[i][j].Status != -1 { // Order finished by other elev
				//fmt.Println("Case 1: Order finished by other elevator")
				myElevInfo.Orders[i][j].Status = -1
				// Replace with finished chan
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: false}
				lightsChannel <- light
			} else if msg.Orders[i][j].Status == 0 && myElevInfo.Orders[i][j].Status == -1 && msg.Orders[i][j].Finished == false { // New order received
				//fmt.Println("Case 2: New Order received from other elevator")
				myElevInfo.Orders[i][j].Status = 0
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <- light
				newOrderChannel <- msg.Orders[i][j]
			} else if msg.Orders[i][j].Status == msg.Id && (myElevInfo.Orders[i][j].Status == 0 || myElevInfo.Orders[i][j].Status == -1) && msg.Orders[i][j].Finished == false { // Other elev taken order
				//fmt.Println("Case 3: Order taken by other elevator")
				myElevInfo.Orders[i][j].Status = msg.Id
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <- light
			} else if myElevInfo.Orders[i][j].Status == -2 {
				newOrderChannel <- Order{Floor: i, ButtonType: j, Status: 0, Finished: false}
				myElevInfo.Orders[i][j].Status = 0
			} else if msg.Orders[i][j].Status == 0 && myElevInfo.Orders[i][j].Status == 0 && myElevInfo.Orders[i][j].Confirm == false { // Order confirmed by other elev
				myElevInfo.Orders[i][j].Confirm = true
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <- light
				newOrderChannel <- msg.Orders[i][j]
			}
		}
	}
}

/*Initialises MyElevInfo variable*/
func initializeMyElevInfo(id int) {
	myElevInfo.Id = id
	myElevInfo.Floor = 0
	myElevInfo.CurrentOrder = Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
	myElevInfo.State = 0
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
// Dev functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

func PrintOrderQueue(queue [numFloors][numButtons]Order) {
	fmt.Println("Orders:")
	for i := numFloors - 1; i >= 0; i-- {
		string := strconv.Itoa(int(queue[i][0].Status)) + strconv.Itoa(int(queue[i][1].Status)) + strconv.Itoa(int(queue[i][2].Status))
		fmt.Println(string)
	}
	fmt.Println("____________")
}
