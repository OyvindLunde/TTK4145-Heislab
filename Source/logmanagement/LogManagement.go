package logmanagement

import (
	"fmt"
	"time"
)

const numFloors = 4
const numButtons = 3

var id string

/*State enum*/
type State int

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
/*type NetworkChannels struct {
	RcvChannel   chan Log
	bcastChannel chan Log
}*/

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
func SendLogFromLocal(bcastChannel chan Log) {
	var message Log

	for {
		message.Orders = OrderQueue
		message.Elev = ElevInfo
		//message.orders = createOrderListFromOrderQueue()
		//message.orders = OrderQueue
		//message.Elev = ElevList[0]
		//message.Elev = ElevInfo
		//fmt.Printf("ElevInfo (msg): %#v\n", message.Elev)
		bcastChannel <- message
		//fmt.Printf("Sending: %#v\n", ElevInfo)
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
			fmt.Printf("Received: %#v\n", a.Elev)
			//updateElevatorQueue(a)
			//updateQueueFromNetwork(a)
		}
	}
}

/**
 * @brief initiates channels and creates coroutines for brodcasting and recieving
 * @param port; port to listen and read on
 */
/*func InitNetwork(port int) {
	RcvChannel = make(chan Log)
	bcastChannel = make(chan Log)
	go network.BrodcastMessage(port, bcastChannel)
	go network.RecieveMessage(port, RcvChannel)
	fmt.Printf("Network initialized\n")
}*/

/**
 * @brief Set id of elev.
 */
/*func setElevID() {

	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
}*/

func updateElevatorQueue(msg Log) {
	var elev = msg.Elev
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

func InitializeElevInfo() {
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
