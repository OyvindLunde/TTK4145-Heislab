package main

import (
	"./display"
	"./elevio"
	"./fsm"
	"./logmanagement"
)

// Network "doesnt work" (for Øyvind), Windows Firewall?

// Changes made:
// Modified StopAtFloor to also take in "Finished" variable
// Commented out a line in FSM - IDLE that may fix our communication problem
// Made an alternative UpdateElevInfo, that is not a goroutine, but instead only is called when something is changed (pings the system less, which is nice)

func main() {
	numFloors := 4
	numButtons := 3
	id := 3
	port := 20009
	addr := 11112

	fsmChannels := fsm.FsmChannels{
		ButtonPress:  make(chan elevio.ButtonEvent),
		FloorReached: make(chan int),
	}

	networkChannels := logmanagement.NetworkChannels{
		RcvChannel:   make(chan logmanagement.Log),
		BcastChannel: make(chan logmanagement.Log),
	}

	fsm.Initialize(numFloors, id, addr)
	go fsm.RunElevator(fsmChannels, numFloors, numButtons)
	go logmanagement.Communication(port, networkChannels)

	go display.Display()

	select {}
}
