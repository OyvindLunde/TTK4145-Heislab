package main

import (
	"fmt"
	"strconv"

	"./display"
	"./elevio"
	"./fsm"
	"./logmanagement"
	"./orderhandler"
)

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------------------------------------------------------------

// numFloors is declared in Logmanagement
// numButtons is declard in Logmangagement
const port = 20009       // address for network, do not change
const tickRate = 2       // seconds per tick
const tickThreshold = 10 // number of ticks needed to generate an interrupt
const heartbeatThreshold = 4

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Main
// ------------------------------------------------------------------------------------------------------------------------------------------------------

func main() {
	id, addr := setParameters() //Function to take in parameters from user

	fsmChannels := fsm.FsmChannels{
		ButtonPress:    make(chan elevio.ButtonEvent),
		FloorReached:   make(chan int),
		MotorDirection: make(chan int, 2),
		ToggleLights:   make(chan elevio.PanelLight, 100*logmanagement.GetNumFloors()*logmanagement.GetNumButtons()),
		NewOrder:       make(chan logmanagement.Order, 100*logmanagement.GetNumFloors()*logmanagement.GetNumButtons()),
		Reset:          make(chan bool),
	}

	networkChannels := logmanagement.NetworkChannels{
		RcvChannel:   make(chan logmanagement.Elev),
		BcastChannel: make(chan logmanagement.Elev),
	}

	fsm.InitFSM(id, addr)
	logmanagement.InitLogManagement(id)
	orderhandler.ReadCabOrderBackup(fsmChannels.ToggleLights, fsmChannels.NewOrder)
	logmanagement.StartTicker(tickRate, tickThreshold, heartbeatThreshold, fsmChannels.ToggleLights, fsmChannels.NewOrder)

	go fsm.RunElevator(fsmChannels)
	go logmanagement.Communication(port, networkChannels, fsmChannels.ToggleLights, fsmChannels.NewOrder, fsmChannels.Reset)

	go display.Display()

	select {} // Select to stop main form exiting scope
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Input Function
// ------------------------------------------------------------------------------------------------------------------------------------------------------

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
