package main

import (
	"./fsm"
)

func main() {
	numFloors := 4
	fsm.Initialize(numFloors)
	fsm.RunElevator()
}
