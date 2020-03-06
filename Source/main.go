package main

import (
	"./fsm"
)

func main() {
	// Make elevator and channels (and more?) here. Use as input in RunElevator
	numFloors := 4

	fsm.Initialize(numFloors)
	fsm.RunElevator()
}
