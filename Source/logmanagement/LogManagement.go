package logmanagement

import (
	"fmt"
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
	LOST	= 3
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
	INACTIVE    = -1
	PENDING 	= 0
	ACTIVE		= 1
	// ACTIVE = ID?
)

/*Elevstruct for keeping info about ther elevs*/
type Elev struct {
	Id           int //endret fra string
	Floor        int
	CurrentOrder Order
	//Lastseen time
	State int
	// Orders?
}

// LogList? Må kunne sende en reset heis cab orders

/*Log to be sendt over the network*/
type Log struct {
	Orders [numFloors][numButtons]Order
	Elev   Elev
	//version time
}

/*Declaration of local log*/
var log1 Log
var ElevInfo Elev
var ElevList []Elev

/*Broadcast and recieve channel*/
type NetworkChannels struct {
	RcvChannel   chan Log
	BcastChannel chan Log
}

//var RcvChannel chan Log
//var bcastChannel chan Log

var OrderQueue = [numFloors][numButtons]Order{}

var Updates = false // Rename?

func InitializeQueue() {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			OrderQueue[i][j].Floor = i
			OrderQueue[i][j].ButtonType = j
			OrderQueue[i][j].Status = -1
			OrderQueue[i][j].Finished = false
		}
	}
	fmt.Println("OrderQueue initialized")
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
func SendLogFromLocal(BcastChannel chan Log) {
	var message Log

	for {
		time.Sleep(20 * time.Millisecond)
		message.Orders = OrderQueue
		message.Elev = ElevInfo
		BcastChannel <- message
		time.Sleep(1000 * time.Millisecond)
	}
}

/**
 * @brief reads message from RcvChannel and does hit width it
 */
func UpdateLogFromNetwork(RcvChannel chan Log) {
	for {
		time.Sleep(20 * time.Millisecond)
		select {
		case a := <-RcvChannel:
			//fmt.Printf("Received: %#v\n", a.Elev.Id)
			if a.Elev.Id != ElevInfo.Id {
				updateElevatorList(a)
				updateQueueFromNetwork(a)
				//fmt.Println("Order4down: ")
				//fmt.Println(OrderQueue[3][1].Active)
				//fmt.Printf("Received: %#v\n", a.Elev.CurrentOrder)
			}
			//fmt.Printf("Received: %#v\n", a.Elev)

		}
	}
}

func Communication(port int, channels NetworkChannels) {
	//RcvChannel := make(chan Log)
	//bcastChannel := make(chan Log)

	go network.RecieveMessage(port, channels.RcvChannel)
	go network.BrodcastMessage(port, channels.BcastChannel)
	go SendLogFromLocal(channels.BcastChannel)
	go UpdateLogFromNetwork(channels.RcvChannel)
	//fmt.Printf("Network initialized\n")
}

func updateElevatorList(msg Log) {
	//fmt.Println("In: ElevatorList")
	//fmt.Println(ElevList)
	var elev = msg.Elev
	if len(ElevList) == 0 {
		ElevList = append(ElevList, elev)
	}
	for _, i := range ElevList {
		if elev.Id == i.Id {
			i.Floor = elev.Floor
			i.State = elev.State
			return
		} else {
			ElevList = append(ElevList, elev)
		}
	}

}

func updateQueueFromNetwork(msg Log) {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons-1; j++ {
			if msg.Orders[i][j].Finished == true {
				OrderQueue[i][j].Status = -1
			} else if msg.Orders[i][j].Status != -1 {
				OrderQueue[i][j].Status = msg.Orders[i][j].Status
			}
		}
	}
	//fmt.Println(msg.Orders)
	//OrderQueue = msg.Orders

}

func GetMatrixDimensions() (rows, cols int) {
	return numFloors, numButtons
}

func InitializeElevInfo(port int) {
	ElevInfo.Id = port
	ElevInfo.Floor = 0
	ElevInfo.CurrentOrder = Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
	ElevInfo.State = 0
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
	ElevInfo.Floor = floor
	ElevInfo.CurrentOrder = order
	ElevInfo.State = state
	Updates = true
}
