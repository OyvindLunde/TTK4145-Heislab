package logmanagement

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module contain functions for updating orders and statuses between elevators
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"strconv"
	"time"
	"../network"
)

const numFloors = 4
const numButtons = 3

var MyElevInfo Elev
var OtherElevInfo []Elev

var DisplayUpdates = false // Used to display the system

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
	Status     OrderStatus // Rename to status
	Finished   bool
	// Timer   timer
	// Confirmed bool
}

/*OrderStatus Enum*/
type OrderStatus int
const (
	PENDING  OrderStatus = 0
	ACTIVE   OrderStatus = 1
	INACTIVE OrderStatus = 2
	// ACTIVE = ID?
)

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

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Setters and Getters
// ------------------------------------------------------------------------------------------------------------------------------------------------------
func GetOrder(floor int, buttonType int) Order {
	return MyElevInfo.Orders[floor][buttonType]
}

func SetOrder(floor int, buttonType int, status OrderStatus, finished bool) {
	 MyElevInfo.Orders[floor][buttonType].Status = status
	 MyElevInfo.Orders[floor][buttonType].Finished = finished
}

func GetOrderList() [numFloors][numButtons]Order{
	return MyElevInfo.Orders
}

func GetElevList() []Elev {
	return OtherElevInfo
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

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Init functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Initializes LogManagement*/
func InitLogManagement(id int, numFloors int, numButtons int){
	numButtons = numButtons
	numFloors = numFloors
	InitializeMyElevInfo(id)
}

/*Initialises MyElevInfo variable*/
func InitializeMyElevInfo(id int) {
	MyElevInfo.Id = id
	MyElevInfo.Floor = 0
	MyElevInfo.CurrentOrder = Order{Floor: -1, ButtonType: -1, Status: 2, Finished: false}
	MyElevInfo.State = 0
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			MyElevInfo.Orders[i][j].Floor = i
			MyElevInfo.Orders[i][j].ButtonType = j
			MyElevInfo.Orders[i][j].Status = 2
			MyElevInfo.Orders[i][j].Finished = false
		}
	}
	fmt.Println("MyElev initialized")
}

/*Inits network communication*/
func InitCommunication(port int, channels NetworkChannels) {
	go network.RecieveMessage(port, channels.RcvChannel)
	go network.BrodcastMessage(port, channels.BcastChannel)
	go SendMyElevInfo(channels.BcastChannel)
	go UpdateOtherElevListFromNetwork(channels.RcvChannel)
	fmt.Printf("Network initialized\n")
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Logic functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Sends MyElevInfo on channel in parameter*/
func SendMyElevInfo(BcastChannel chan Elev) {
	for {
		time.Sleep(20 * time.Millisecond)
		//fmt.Println("Sending:")
		//PrintOrderQueue(MyElevInfo.orders)
		BcastChannel <- MyElevInfo
	}
}

/*Updates OtherElevLsit from channel in parameter*/
func UpdateOtherElevListFromNetwork(RcvChannel chan Elev) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-RcvChannel:
			if a.Id != MyElevInfo.Id {
				fmt.Println("Receiving:")
				PrintOrderQueue(a.Orders)
				updateOtherElevInfo(a)
				updateOrderList(a)
			}
		}
	}
}



/*Updates otherelevinfo with info about elev in param*/
func updateOtherElevInfo(msg Elev) {
	for _, i := range OtherElevInfo {
		if msg.Id == i.Id {
			//fmt.Println("Correct ID")
			i.Floor = msg.Floor
			i.CurrentOrder = msg.CurrentOrder
			i.State = msg.State
			i.Orders = msg.Orders
			return
		}
	}
	OtherElevInfo = append(OtherElevInfo, msg)

}


/*Updates orderlist with data stored in elev-param*/
func updateOrderList(msg Elev) {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons-1; j++ {
			if msg.Orders[i][j].Finished == true {
				MyElevInfo.Orders[i][j].Status = 2
			} else if msg.Orders[i][j].Status != 2 {
				MyElevInfo.Orders[i][j].Status = msg.Orders[i][j].Status
			}
		}
	}
	//DisplayUpdates = true
}




/*Updates MyElevInfo variable from params*/
func UpdateMyElevInfo(floor int, order Order, state int) {
	MyElevInfo.Floor = floor
	MyElevInfo.CurrentOrder = order
	MyElevInfo.State = state
	DisplayUpdates = true
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