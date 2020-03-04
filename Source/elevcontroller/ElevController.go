package elevcontroller

import (
	"time"

	"../elevio"
)

var ButtonPress chan elevio.ButtonEvent
var FloorReached chan int

type Button struct {
	Floor int
	Type  int
}

func initializeLights(numFloors int) {
	for i := 0; i < numFloors; i++ {
		if i != 0 {
			elevio.SetButtonLamp(elevio.BT_HallDown, i, false)
		}
		if i != numFloors {
			elevio.SetButtonLamp(elevio.BT_HallUp, i, false)
		}
		elevio.SetButtonLamp(elevio.BT_Cab, i, false)
	}
	elevio.SetDoorOpenLamp(false)

}

func InitializeElevator(numFloors int) {
	elevio.Init("localhost:15657", numFloors)
	initializeLights(numFloors)
	elevio.SetMotorDirection(elevio.MD_Down)
	for elevio.GetFloor() != 0 { //Fix getFloor problemet
	}
	elevio.SetMotorDirection(elevio.MD_Stop)

	/*ButtonPress = make(chan elevio.ButtonEvent)
	FloorReached = make(chan int)

	go elevio.PollButtons(ButtonPress)
	go elevio.PollFloorSensor(FloorReached)*/
}

func ButtonPressed(receiver chan<- Button) {
	for {
		ButtonPress := make(chan elevio.ButtonEvent)
		elevio.PollButtons(ButtonPress)
		select {
		case a := <-ButtonPress:
			receiver <- Button{Floor: a.Floor, Type: int(a.Button)}
		}
	}
}

func FloorIsReached(receiver chan<- int) {
	for {
		FloorReached := make(chan int)
		elevio.PollFloorSensor(FloorReached)
		select {
		case a := <-FloorReached:
			receiver <- a
		}
	}
}

func OpenCloseDoor(seconds time.Duration) {
	elevio.SetDoorOpenLamp(true)
	time.Sleep(seconds * 1000 * time.Millisecond)
	elevio.SetDoorOpenLamp(false)
}

func ElevStopAtFloor(floor int) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	OpenCloseDoor(3)
}

func UpdateLight(button Button, value bool) {
	elevio.SetButtonLamp(elevio.ButtonType(button.Type), button.Floor, value)
}
