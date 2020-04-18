package ticker

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module handles the functions regarding the ticker
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"time"
)

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------------------------------------------------------------

var Done chan bool
var Ticker *time.Ticker

var elevTickerInfo []int
var heartbeat []int

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Public functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

func GetElevTickerInfo() []int { // Move/Delete
	return elevTickerInfo
}

func IncrementElevTickerInfo(elev int) {
	elevTickerInfo[elev] += 1
}

func ResetElevTickerInfo(elev int) {
	elevTickerInfo[elev] = 0
}

func ClearElevTickerInfo() {
	elevTickerInfo = elevTickerInfo[:0]
}

func AppendToElevTickerInfo() {
	elevTickerInfo = append(elevTickerInfo, 0)
}

func RemoveElevFromelevTickerInfo(i int) {
	copy(elevTickerInfo[i:], elevTickerInfo[i+1:])          // Shift a[i+1:] left one index.
	elevTickerInfo = elevTickerInfo[:len(elevTickerInfo)-1] // Truncate slice.
}

func GetHeartBeat(index int) int {
	return heartbeat[index]
}

func IncrementHeartBeat() {
	for i := 0; i < len(heartbeat); i++ {
		heartbeat[i]++
	}
}

func ResetHeartBeat(i int) {
	heartbeat[i] = 0
}

func AppendToHeartBeat() {
	heartbeat = append(heartbeat, 0)
}

func RemoveHeartbeat(i int) {
	copy(heartbeat[i:], heartbeat[i+1:])     // Shift a[i+1:] left one index.
	heartbeat = heartbeat[:len(heartbeat)-1] // Truncate slice.
}
