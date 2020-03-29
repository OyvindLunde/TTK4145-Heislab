package logmanagement

import (
	"fmt"
	"strconv"
	"time"
	"../network"
)

const numFloors = 4
const numButtons = 3

//var Id int

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

// LogList? Må kunne sende en reset heis cab orders

/*Log to be sendt over the network*/

/*Declaration of local log*/

var MyElevInfo Elev
var OtherElevInfo []Elev

/*Broadcast and recieve channel*/
type NetworkChannels struct {
	RcvChannel   chan Elev
	BcastChannel chan Elev
}

//var RcvChannel chan Log
//var bcastChannel chan Log

var OrderQueue = [numFloors][numButtons]Order{}

var DisplayUpdates = false // Used to display the system

func InitializeQueue() {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			OrderQueue[i][j].Floor = i
			OrderQueue[i][j].ButtonType = j
			OrderQueue[i][j].Status = 2
			OrderQueue[i][j].Finished = false
		}
	}
	fmt.Println("OrderQueue initialized")
	//fmt.Println(OrderStatus(ACTIVE))
}

func GetOrder(floor int, buttonType int) Order {
	return OrderQueue[floor][buttonType]
}

func GetElevInfo(elev Elev) (id, floor int, currentOrder Order, state int) {
	return elev.Id, elev.Floor, elev.CurrentOrder, elev.State
}

/**
 * @brief puts message on bcastChannel
 * @param Message; message to be transmitted
 */
func SendMyElevInfo(BcastChannel chan Elev) {
	for {
		time.Sleep(20 * time.Millisecond)
		//fmt.Println("Sending:")
		//PrintOrderQueue(OrderQueue)
		BcastChannel <- MyElevInfo
	}
}

/**
 * @brief reads message from RcvChannel and does hit width it
 */
func UpdateLogFromNetwork(RcvChannel chan Elev) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-RcvChannel:
			//fmt.Printf("Received: %#v\n", a.Elev.Id)
			if a.Id != MyElevInfo.Id {
				//fmt.Println("Receiving:")
				//PrintOrderQueue(a.Orders)
				updateElevatorList(a)
				updateQueueFromNetwork(a)
			}
		}
	}
}

func PrintOrderQueue(queue [numFloors][numButtons]Order) {
	fmt.Println("Orders:")
	for i := numFloors - 1; i >= 0; i-- {
		string := strconv.Itoa(int(queue[i][0].Status)) + strconv.Itoa(int(queue[i][1].Status)) + strconv.Itoa(int(queue[i][2].Status))
		fmt.Println(string)
	}
	fmt.Println("____________")
}

func Communication(port int, channels NetworkChannels) {
	//RcvChannel := make(chan Log)
	//bcastChannel := make(chan Log)

	go network.RecieveMessage(port, channels.RcvChannel)
	go network.BrodcastMessage(port, channels.BcastChannel)
	go SendMyElevInfo(channels.BcastChannel)
	go UpdateLogFromNetwork(channels.RcvChannel)
	//fmt.Printf("Network initialized\n")
}

func updateElevatorList(msg Elev) {
	//fmt.Println("In: ElevatorList")
	//PrintOrderQueue(msg.Orders)
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

func updateQueueFromNetwork(msg Elev) {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons-1; j++ {
			if msg.Orders[i][j].Finished == true {
				OrderQueue[i][j].Status = 2
			} else if msg.Orders[i][j].Status != 2 {
				OrderQueue[i][j].Status = msg.Orders[i][j].Status
			}
		}
	}
	//DisplayUpdates = true
}

func GetMatrixDimensions() (rows, cols int) {
	return numFloors, numButtons
}

func InitializeElevInfo(id int) {
	MyElevInfo.Id = id
	MyElevInfo.Floor = 0
	MyElevInfo.CurrentOrder = Order{Floor: -1, ButtonType: -1, Status: 2, Finished: false}
	MyElevInfo.State = 0
	MyElevInfo.Orders = OrderQueue
}

/*func UpdateElevInfo(floor *int, order *Order, state *int) {
	for {
		time.Sleep(5 * time.Millisecond)
		ElevInfo.Floor = *floor
		ElevInfo.CurrentOrder = *order
		ElevInfo.State = *state
		//fmt.Println(ElevInfo)

	}

}*/

func UpdateElevInfo(floor int, order Order, state int) {
	MyElevInfo.Floor = floor
	MyElevInfo.CurrentOrder = order
	MyElevInfo.State = state
	MyElevInfo.Orders = OrderQueue
	DisplayUpdates = true
}
