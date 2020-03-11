package logmanagement

import (
	"time"

	"../network"
)

const numFloors = 4
const numButtons = 3

//var Id int

/*State enum*/
type State int // Kanskje slette denne?

const (
	IDLE    = 0
	EXECUTE = 1
	LOST    = 2
	RESET   = 3
)

/*OrderStruct*/
type Order struct {
	Floor      int
	ButtonType int
	Active     int // Rename to status
	// Timer   timer
}

/*Elevstruct for keeping info about ther elevs*/
type Elev struct {
	Id           int //endret fra string
	Floor        int
	CurrentOrder Order
	//Lastseen time
	State int
}

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

func InitializeQueue() {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			//queue[i][j] = nil
			OrderQueue[i][j].Floor = i
			OrderQueue[i][j].ButtonType = j
			OrderQueue[i][j].Active = -1
		}
	}
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
		//time.Sleep(100 * time.Millisecond)
		select {
		case a := <-RcvChannel:
			//fmt.Printf("Received: %#v\n", a.Elev.Id)
			if a.Elev.Id != ElevInfo.Id {
				updateElevatorList(a)
				updateQueueFromNetwork(a)
				//fmt.Println(OrderQueue)
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
	//fmt.Println(msg.Orders)
	OrderQueue = msg.Orders
	/*for _, order := range msg.orders {
		OrderQueue[order.Floor][order.ButtonType].Active = order.Active
	}*/

}

/*func createOrderListFromOrderQueue() []Order {
	listedOrders := [][]Order{}
	for i := range OrderQueue {
		for k := range i {
			temp := Order{}
			temp.Floor = i
			temp.ButtonType = k
			temp.Active = OrderQueue[i][k].Active
			listedOrders = append(listedOrders, temp)
		}
	}
	return listedOrders

}*/
func GetMatrixDimensions() (rows, cols int) {
	return numFloors, numButtons
}

func InitializeElevInfo(port int) {
	ElevInfo.Id = port
	ElevInfo.Floor = 0
	ElevInfo.CurrentOrder = Order{Floor: -1, ButtonType: -1, Active: -1}
	ElevInfo.State = 0
}

func UpdateElevInfo(floor *int, order *Order, state *int) {
	for {
		time.Sleep(20 * time.Millisecond)
		ElevInfo.Floor = *floor
		ElevInfo.CurrentOrder = *order
		ElevInfo.State = *state
		//fmt.Println(ElevInfo)

	}

}
