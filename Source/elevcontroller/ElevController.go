package elevcontroller

import (
	"strconv"
	"time"

	"../elevio"
)

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

func InitializeElevator(numFloors int, port int) {
	elevio.Init("localhost:"+strconv.Itoa(port), numFloors)
	initializeLights(numFloors)
	elevio.SetMotorDirection(elevio.MD_Down)
	for elevio.GetFloor() != 0 {
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(0)
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
