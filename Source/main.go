package main

import (
	"./elevio"
	"./fsm"
	"./logmanagement"
)

func main() {
	numFloors := 4
	port := 15657

	fsmChannels := fsm.FsmChannels{
		ButtonPress:  make(chan elevio.ButtonEvent),
		FloorReached: make(chan int),
	}

	networkChannels := logmanagement.NetworkChannels{
		RcvChannel:   make(chan logmanagement.Log),
		BcastChannel: make(chan logmanagement.Log),
	}

	fsm.Initialize(numFloors, port)
	go fsm.RunElevator(fsmChannels)
	go logmanagement.Communication(port, networkChannels)

	select {}
}
