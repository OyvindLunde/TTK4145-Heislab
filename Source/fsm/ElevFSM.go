package fsm

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Elevator Finite State Machine
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"time"

	"../elevcontroller"
	"../elevio"
	"../logmanagement"
	"../orderhandler"
)

var address int
var _id int

//var state int

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Structs and enums
// ------------------------------------------------------------------------------------------------------------------------------------------------------

type State int

const (
	INIT    = 0
	IDLE    = 1
	EXECUTE = 2
	LOST    = 3
	RESET   = 4
)

type FsmChannels struct {
	ButtonPress    chan elevio.ButtonEvent
	FloorReached   chan int
	MotorDirection chan int
	ToggleLights   chan elevio.PanelLight
	NewOrder       chan logmanagement.Order
	Reset          chan bool
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Init and FSM
// ------------------------------------------------------------------------------------------------------------------------------------------------------

func InitFSM(id int, addr int) {
	_id = id
	address = addr
	elevcontroller.InitializeElevator(logmanagement.GetNumFloors(), addr)
	fmt.Println("FSM Initialized")
}

/*Elevator FSM*/
func RunElevator(channels FsmChannels) {
	floor := 0
	state := IDLE
	NoOrder := logmanagement.Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
	currentOrder := NoOrder

	go elevio.PollButtons(channels.ButtonPress)
	go elevio.PollFloorSensor(channels.FloorReached)
	go orderhandler.HandleButtonEvents(channels.ButtonPress, channels.ToggleLights, channels.NewOrder)
	go orderhandler.UpdateLights(channels.ToggleLights)

	for {
		time.Sleep(20 * time.Millisecond)
		switch state {
		case IDLE:
			select {
			case currentOrder = <-channels.NewOrder:
				if orderhandler.IsOrderValid(currentOrder) {
					currentOrder.Status = logmanagement.GetMyElevInfo().Id
					logmanagement.SetMyElevInfo(floor, currentOrder, state)
					if orderhandler.ShouldITakeOrder(currentOrder) {
						orderhandler.UpdateOrder(currentOrder.Floor, int(currentOrder.ButtonType), logmanagement.GetMyElevInfo().Id, false, true)
						channels.MotorDirection <- orderhandler.GetDirection(floor, currentOrder.Floor)
						state = EXECUTE
						logmanagement.SetMyElevInfo(floor, currentOrder, state)
						break
					}
				}
				logmanagement.SetMyElevInfo(floor, NoOrder, state)
			}

		case EXECUTE:
			select {
			case dir := <-channels.MotorDirection:
				elevio.SetMotorDirection(elevio.MotorDirection(dir))
				if dir == 0 {
					orderhandler.StopAtFloor(floor, channels.ToggleLights)
					state = IDLE
					logmanagement.SetMyElevInfo(floor, NoOrder, state)
				}

			case floor = <-channels.FloorReached:
				logmanagement.SetMyElevInfo(floor, currentOrder, state)
				elevio.SetFloorIndicator(floor)
				if orderhandler.ShouldElevatorStop(floor, currentOrder.Floor) {
					orderhandler.StopAtFloor(floor, channels.ToggleLights)
					dir := orderhandler.GetDirection(floor, currentOrder.Floor)
					if dir == 0 {
						state = IDLE
						logmanagement.SetMyElevInfo(floor, NoOrder, state)
						for len(channels.MotorDirection) > 0 { // Emptying MotorDirection channel
							<-channels.MotorDirection
						}
					} else {
						channels.MotorDirection <- dir
					}
				}
			case <-channels.Reset:
				state = RESET
			}

		case RESET:
			logmanagement.InitLogManagement(_id, channels.ToggleLights, channels.NewOrder)
			elevcontroller.InitializeElevator(logmanagement.GetNumFloors(), address)
			orderhandler.ReadCabOrderBackup(channels.ToggleLights, channels.NewOrder)
			floor = 0
			state = IDLE
			logmanagement.SetMyElevInfo(floor, NoOrder, state)
		}
	}
}
