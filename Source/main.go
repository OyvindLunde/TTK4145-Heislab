package main

import (
	"./elevio"
	"./fsm"
	//"./logmanagement"
)

func main() {
	// Make elevator and channels (and more?) here. Use as input in RunElevator
	numFloors := 4

	fsmChannels := fsm.FsmChannels{
		ButtonPress:  make(chan elevio.ButtonEvent),
		FloorReached: make(chan int),
	}

	//networkChannels := logmanagement.
	//logmanagement.InitNetwork(15374)
	//go logmanagement.SendLogFromLocal()
	//logmanagement.UpdateLogFromNetwork()
	fsm.Initialize(numFloors)
	fsm.RunElevator(fsmChannels)
}
