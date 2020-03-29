package main

import (
	"fmt"
	"strconv"

	"./display"
	"./elevio"
	"./fsm"
	"./logmanagement"
)

// Changes made:
// Modified StopAtFloor to also take in "Finished" variable
// Commented out a line in FSM - IDLE that may fix our communication problem
// Made an alternative UpdateElevInfo, that is not a goroutine, but instead only is called when something is changed (pings the system less, which is nice)

func main() {
	numFloors := 4
	numButtons := 3
	//id := 1       // elevator id, change for each elevator
	port := 20009 // address for network, do not change
	//addr := 11111 // address for tcp connection to simulator, change for each elevator

	id, addr := setParameters()

	fsmChannels := fsm.FsmChannels{
		ButtonPress:  make(chan elevio.ButtonEvent),
		FloorReached: make(chan int),
		ToggleLights: make(chan elevio.PanelLight),
	}

	networkChannels := logmanagement.NetworkChannels{
		RcvChannel:   make(chan logmanagement.Elev),
		BcastChannel: make(chan logmanagement.Elev),
	}

	fsm.Initialize(numFloors, id, addr)
	logmanagement.InitLogManagement(id, numFloors, numButtons)
	go fsm.RunElevator(fsmChannels, numFloors, numButtons)
	go logmanagement.InitCommunication(port, networkChannels)

	go display.Display()

	select {}
}

func setParameters() (int, int) {
	var input string
	fmt.Print("Enter Id: ")
	fmt.Scanf("%s", &input)
	id, _ := strconv.Atoi(input)
	fmt.Print("Enter Address: ")
	fmt.Scanf("%s", &input)
	addr, _ := strconv.Atoi(input)
	//fmt.Println(id)
	//fmt.Println(addr)

	return id, addr
}
