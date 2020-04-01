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
	ButtonPress  chan elevio.ButtonEvent
	FloorReached chan int
	ToggleLights chan elevio.PanelLight
	NewOrder     chan logmanagement.Order
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Init and FSM
// ------------------------------------------------------------------------------------------------------------------------------------------------------

func Initialize(numFloors int, id int, addr int) {
	elevcontroller.InitializeElevator(numFloors, addr)
	elevio.SetFloorIndicator(0)
}

/*Elevator FSM*/
func RunElevator(channels FsmChannels, numFloors int, numButtons int) {
	fmt.Println("AutoHeis assemble")
	//destination := -1
	dir := 0
	floor := 0
	state := IDLE
	NoOrder := logmanagement.Order{Floor: -1, ButtonType: -1, Status: -1, Finished: false}
	currentOrder := NoOrder

	go elevio.PollButtons(channels.ButtonPress) // Kan vi legge denne inn i HandleButtonEvents?
	go elevio.PollFloorSensor(channels.FloorReached)
	go orderhandler.HandleButtonEvents(channels.ButtonPress, channels.ToggleLights, channels.NewOrder)
	go orderhandler.UpdateLightsV2(channels.ToggleLights)
	//go logmanagement.UpdateMyElevInfo(&floor, &currentOrder, &state) // Vurdere å droppe denne? Kjører unødvendig ofte

	for {
		time.Sleep(20 * time.Millisecond)
		switch state {
		case IDLE:
			select {
			case currentOrder = <-channels.NewOrder:
				fmt.Println(currentOrder)
				if orderhandler.IsOrderValid(currentOrder) {
					currentOrder.Status = logmanagement.GetMyElevInfo().Id // Remove this?
					logmanagement.SetMyElevInfo(floor, currentOrder, state)
					orderhandler.UpdateLocalOrders(currentOrder.Floor, int(currentOrder.ButtonType), logmanagement.GetMyElevInfo().Id, false)
					dir = orderhandler.GetDirection(floor, currentOrder.Floor)
					state = EXECUTE
					logmanagement.SetMyElevInfo(floor, currentOrder, state)
				}
			}

		case EXECUTE:
			elevio.SetMotorDirection(elevio.MotorDirection(dir)) // Blir kalt (unødvendig) mange ganger. Men er "sikker"
			select {
			case a := <-channels.FloorReached:
				floor = a
				logmanagement.SetMyElevInfo(floor, currentOrder, state)
				elevio.SetFloorIndicator(floor)
				if orderhandler.ShouldElevatorStop(floor, currentOrder.Floor, logmanagement.GetMyElevInfo(), logmanagement.GetOtherElevInfo()) {
					orderhandler.StopAtFloor(floor, channels.ToggleLights)
					dir = orderhandler.GetDirection(floor, currentOrder.Floor)
					elevio.SetMotorDirection(elevio.MotorDirection(dir))
					if dir == 0 { // Forslag: Legge inn en CheckForCABOrders funksjon, må i så fall inn i default også
						//destination = -1 // Unødvendig?
						state = IDLE
						logmanagement.SetMyElevInfo(floor, NoOrder, state)
					}
				}

			default:
				if dir == 0 {
					orderhandler.StopAtFloor(floor, channels.ToggleLights)
					state = IDLE
					logmanagement.SetMyElevInfo(floor, NoOrder, state)
				}
			}

		case RESET:
			//reset elevator

		}
	}
}
