package main

import (
	"./elevio"
	"./fsm"
	//"./logmanagement"
	"./display"
)

func main() {
	numFloors := 4
	numButtons := 3
	port := 15657

	fsmChannels := fsm.FsmChannels{
		ButtonPress:  make(chan elevio.ButtonEvent),
		FloorReached: make(chan int),
	}

	/*networkChannels := logmanagement.NetworkChannels{
		RcvChannel:   make(chan logmanagement.Log),
		BcastChannel: make(chan logmanagement.Log),
	}*/

	fsm.Initialize(numFloors, port)
	go fsm.RunElevator(fsmChannels, numFloors, numButtons)
	//go logmanagement.Communication(port, networkChannels)

	go display.Display()

	select {}
}
