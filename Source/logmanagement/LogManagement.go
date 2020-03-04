package logmanagement

import (
	"fmt"
	"time"

	network "../network"
)

const numFloors = 4
const numButtons = 3

var id string

/*State enum*/
type State int

const (
	Idle = 0
	Exec = 1
	Lost = 2
)

/*OrderStruct*/
type Order struct {
	floor      int //Remove
	ButtonType int //Remove
	active     int
	// Timer?
}

/*Elevstruct for keeping info about ther elevs*/
type Elev struct {
	id    string
	floor int
	//Lastseen time
	state int
}

/*Log to be sendt over the network*/
type log struct {
	orders  []Order
	Elev Elev
	//version time
}

/*Declaration of local log*/
var log1 log

var elevList []Elev

/*Broadcast and recieve channel*/
var RcvChannel chan log
var bcastChannel chan log

func NewOrder(floor int, buttonType int, active int) Order { // Overflødig per nå
	order := Order{Floor: floor, ButtonType: buttonType, Active: active}
	return order
}

var OrderQueue = &[numFloors][numButtons]Order{}

func InitializeQueue(queue *[numFloors][numButtons]Order) {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			//queue[i][j] = nil
			queue[i][j].Floor = i
			queue[i][j].ButtonType = j
			queue[i][j].Active = -1
		}
	}
}

/*func CheckForOrders(queue *[numFloors][numButtons]Order, receiver chan<- Order) { // Velger den første i lista, ikke den eldste ordren
	for { // Legg i orderHandler?
		time.Sleep(20 * time.Millisecond)
		for i := 0; i < numFloors; i++ {
			for j := 0; j < numButtons; j++ {
				if queue[i][j].Active == 0 {
					fmt.Printf("%+v\n", queue[i][j])
					receiver <- queue[i][j]
				}
			}
		}
	}
}*/

// UpdateOrderQueue updates the order queue
func UpdateOrderQueue(floor int, button int, active int) {
	OrderQueue[floor][button].Active = active
}

// GetActiveOrder returns the first found active order
func GetPendingOrder() Order {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			if OrderQueue[i][j].Active == 0 {
				return OrderQueue[i][j]
			}
		}
	}
	return Order{Floor: -1, ButtonType: -1, Active: -1}
}

func GetOrder(floor int, buttonType int) Order {
	return OrderQueue[floor][buttonType]
}

/**
 * @brief puts message on bcastChannel
 * @param Message; message to be transmitted
 */
func UpdateLogFromLocal() {
	var message log
	message.orders = createOrderListFromOrderQueue()
	message.Elev = elevList[0]
	for {
		bcastChannel <- message
		time.Sleep(1 * time.Second)
	}
}

/**
 * @brief reads message from RcvChannel and does hit width it
 */
func UpdateLogFromNetwork() {
	for {
		a := <-RcvChannel
		updateElevatorQueue(a)
		updateOrderQueue(a)
		fmt.Printf("Received: %#v\n", a)
	}
}

/**
 * @brief initiates channels and creates coroutines for brodcasting and recieving
 * @param port; port to listen and read on
 */
func InitNetwork(port int) {
	RcvChannel = make(chan log)
	bcastChannel = make(chan log)
	go network.BrodcastMessage(port, bcastChannel)
	go network.RecieveMessage(port, RcvChannel)
	fmt.Printf("log initialized\n")
}

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

func updateElevatorQueue(msg log) {
	var elev = msg.Elev
	for _, i := range elevList {
		if elev.id == i.id {
			i.floor = elev.floor
			i.state = elev.state
			return
		}else{
			elevList = append(elevList, elev)
		}
	}
}

func updateOrderQueue(msg log) {
	for  _, order := range msg.orders{
		OrderQueue[order.floor][order.ButtonType].active = order.active
	}
	
}

func createOrderListFromOrderQueue() []order {
	listedOrders []order
	for  i := range OrderQueue {
        for k := range i{
			var temp order
			temp.floor = i
			temp.buttonType = k
			temp.active = OrderQueue[i][k].active
			listedOrders = append(listedOrders, temp)
		}
	}
	return listedOrders
	
}