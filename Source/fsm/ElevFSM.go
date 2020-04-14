package fsm

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module contain functions for executing this elevators orders
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
	Reset 		   chan bool
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Init and FSM
// ------------------------------------------------------------------------------------------------------------------------------------------------------

var state int

func InitFSM(id int, addr int) {
	_id = id
	address = addr
	elevcontroller.InitializeElevator(logmanagement.GetNumFloors(), addr)
	elevio.SetFloorIndicator(0)
}

/*Elevator FSM*/
func RunElevator(channels FsmChannels) {
	fmt.Println("AutoHeis assemble")
	floor := 0
	state = IDLE
	NoOrder := logmanagement.Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
	currentOrder := NoOrder

	go elevio.PollButtons(channels.ButtonPress) // Kan vi legge denne inn i HandleButtonEvents?
	go elevio.PollFloorSensor(channels.FloorReached)
	go orderhandler.HandleButtonEvents(channels.ButtonPress, channels.ToggleLights, channels.NewOrder)
	go orderhandler.UpdateLightsV2(channels.ToggleLights)
	go ResetElev(channels)

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
						orderhandler.UpdateLocalOrders(currentOrder.Floor, int(currentOrder.ButtonType), logmanagement.GetMyElevInfo().Id, false, true)
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

			case a := <-channels.FloorReached: //annet navn enn "a"?
				floor = a
				logmanagement.SetMyElevInfo(floor, currentOrder, state)
				elevio.SetFloorIndicator(floor)
				if orderhandler.ShouldElevatorStop(floor, currentOrder.Floor, logmanagement.GetMyElevInfo(), logmanagement.GetOtherElevInfo()) {
					orderhandler.StopAtFloor(floor, channels.ToggleLights)
					dir := orderhandler.GetDirection(floor, currentOrder.Floor)
					if dir == 0 {
						state = IDLE
						logmanagement.SetMyElevInfo(floor, NoOrder, state)
					} else {
						channels.MotorDirection <- dir
					}
				}
			}

		}
	}
}

func ResetElev(channel FsmChannels){
	for{
		time.Sleep(20 * time.Millisecond)
		select{
		case  <- channel.Reset:
			elevio.SetMotorDirection(elevio.MD_Down)
			for elevio.GetFloor() != 0 { //Fix getFloor problemet
			}
			elevio.SetMotorDirection(elevio.MD_Stop)
			logmanagement.InitLogManagement(_id)
			state = IDLE
			// readcaborderbackup
		}
			
	}
}