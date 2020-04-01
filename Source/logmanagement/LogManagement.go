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

const timerLength = 20;
var myElevInfo Elev
var otherElevInfo []Elev

var displayUpdates = false // Used to display the system
var orderTimer []ElevTimer
// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Declaration of structs and Enums
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*State enum*/
type State int // Kanskje slette denne? Eksisterer i FSM også
const (
	INIT    = 0
	IDLE    = 1
	EXECUTE = 2
	LOST    = 3
	RESET   = 4
)

type Order struct {
	Floor      int
	ButtonType int
	Status     int
	Finished   bool
	// Timer   timer
	// Confirmed bool
}

/*Elevstruct for keeping info about ther elevs*/
type Elev struct {
	Id           int //endret fra string
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

type ElevTimer struct {
	id int
	timer timer
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Setters and Getters
// ------------------------------------------------------------------------------------------------------------------------------------------------------
func GetOrder(floor int, buttonType int) Order {
	return myElevInfo.Orders[floor][buttonType]
}

func SetOrder(floor int, buttonType int, status OrderStatus, finished bool) {
	myElevInfo.Orders[floor][buttonType].Status = status
	myElevInfo.Orders[floor][buttonType].Finished = finished
}

func GetOrderList() [numFloors][numButtons]Order {
	return myElevInfo.Orders
}

func GetElevList() []Elev {
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


func GetMyElevInfo() Elev {
	return myElevInfo
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
func InitLogManagement(id int, numFloors int, numButtons int) {
	numButtons = numButtons
	numFloors = numFloors
	InitializeMyElevInfo(id)
}

/*Initialises MyElevInfo variable*/
func InitializeMyElevInfo(id int) {
	MyElevInfo.Id = id
	MyElevInfo.Floor = 0
	MyElevInfo.CurrentOrder = Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
	MyElevInfo.State = 0
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			myElevInfo.Orders[i][j].Floor = i
			myElevInfo.Orders[i][j].ButtonType = j
			myElevInfo.Orders[i][j].Status = 2
			myElevInfo.Orders[i][j].Finished = false
		}
	}
	fmt.Println("MyElev initialized")
}

/*Inits network communication*/
func InitCommunication(port int, channels NetworkChannels, toggleLights chan elevio.PanelLight) {
	go network.RecieveMessage(port, channels.RcvChannel)
	go network.BrodcastMessage(port, channels.BcastChannel)
	go SendMyElevInfo(channels.BcastChannel)
	go UpdateFromNetwork(channels.RcvChannel, toggleLights)
	fmt.Printf("Network initialized\n")
}

func StartTicker(){
	for i :=0 ; i < len(orderTimer); i++{
		if i.id == elevId{
			orderTimer[i] := time.NewTimer(timerLength * time.Second)
		}
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Logic functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Sends MyElevInfo on channel in parameter*/
func SendMyElevInfo(BcastChannel chan Elev) {
	for {
		time.Sleep(20 * time.Millisecond)
		BcastChannel <- MyElevInfo
	}
}

/*Updates OtherElevLsit from channel in parameter*/
func UpdateFromNetwork(RcvChannel chan Elev, lightsChannel chan<- elevio.PanelLight) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-RcvChannel:
			if a.Id != MyElevInfo.Id {
				/*fmt.Println("Received:")
				PrintOrderQueue(a.Orders)
				fmt.Println("_____________")*/
				updateOtherElevInfo(a)
				updateOrderList(a, lightsChannel)
			}
		}
	}
}

/*Updates otherelevinfo with info about elev in param*/
func updateOtherElevInfo(msg Elev) {
	for i := 0; i < len(OtherElevInfo); i++ {
		if msg.Id == OtherElevInfo[i].Id {
			otherElevInfo[i].Floor = msg.Floor
			otherElevInfo[i].CurrentOrder = msg.CurrentOrder
			otherElevInfo[i].State = msg.State
			otherElevInfo[i].Orders = msg.Orders
			/*fmt.Println("Other elevs orders:")
			PrintOrderQueue(OtherElevInfo[0].Orders)
			fmt.Println("__________")*/
			return
		}
	}
	otherElevInfo = append(OtherElevInfo, msg)
	displayUpdates = true

}

/*Updates orderlist with data stored in elev-param*/
func updateOrderList(msg Elev, lightsChannel chan<- elevio.PanelLight) {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons-1; j++ {
			if msg.Orders[i][j].Finished == true && MyElevInfo.Orders[i][j].Status != -1 {
				fmt.Println("case 1")
				myElevInfo.Orders[i][j].Status = 2 
				// Replace with finished chan
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: false}
				lightsChannel <-light
			} else if msg.Orders[i][j].Status == 0 && MyElevInfo.Orders[i][j].Status == -1 {
				fmt.Println("case 2")
				myElevInfo.Orders[i][j].Status = 0
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <-light
			} else if msg.Orders[i][j].Status == 1 && MyElevInfo.Orders[i][j].Status == 0 && msg.Orders[i][j].Finished == false {
				fmt.Println("case 3")
				myElevInfo.Orders[i][j].Status = 1
				light := elevio.PanelLight{Floor: i, Button: elevio.ButtonType(j), Value: true}
				lightsChannel <-light
			}
		}
	}
	//DisplayUpdates = true
}

/*Updates MyElevInfo variable from params*/
func UpdateMyElevInfo(floor int, order Order, state int) {
	myElevInfo.Floor = floor
	myElevInfo.CurrentOrder = order
	myElevInfo.State = state
	displayUpdates = true
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
