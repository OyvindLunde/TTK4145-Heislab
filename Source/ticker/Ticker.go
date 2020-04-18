package ticker

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// This Module handles the functions regarding the ticker
// ------------------------------------------------------------------------------------------------------------------------------------------------------

import (
	"sync"
	"time"
)

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------------------------------------------------------------

var heartbeatThreshold int
var currentOrderThres int

var done chan bool
var ticker *time.Ticker

var orderTicker map[int]int //Keeps track of how long an order has been active
var heartbeat map[int]int   //Keeps track of how long its been since last we heard from an elevator
var _mtx sync.Mutex

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Public functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------
/*Starts ticker and check if the other elevators finishes orders within ticklength * tickTreshold seconds*/
func StartTicker(tickLength time.Duration, heartBeatThreshold int, currentOrderThreshold int) {
	currentOrderThres = currentOrderThreshold
	heartbeatThreshold = heartBeatThreshold
	done = make(chan bool)
	ticker = time.NewTicker(tickLength * time.Second)
	orderTicker = make(map[int]int)
	heartbeat = make(map[int]int)
	_mtx = sync.Mutex{}
	go elevTicker()
}

/*Stops ticker*/
func StoppTicker() {
	ticker.Stop()
	done <- true

}

func ResetOrderTicker(id int) {
	_mtx.Lock()
	defer _mtx.Unlock()
	orderTicker[id] = 0
}

func AddElevToTicker(id int) {
	_mtx.Lock()
	defer _mtx.Unlock()
	heartbeat[id] = 0
	orderTicker[id] = 0
}

func DeleteElevFromTicker(id int) {
	_mtx.Lock()
	defer _mtx.Unlock()
	delete(heartbeat, id)
	delete(orderTicker, id)
}

func ResetHeartBeat(id int) {
	_mtx.Lock()
	defer _mtx.Unlock()
	heartbeat[id] = 0
}

func IsElevAlive(id int) bool {
	return heartbeat[id] < heartbeatThreshold
}

func HasCurrentOrderTimedOut(id int) bool {
	return orderTicker[id] > currentOrderThres
}

func ClearElevTickerInfo() {
	for key, _ := range heartbeat {
		DeleteElevFromTicker(key)
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Private functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*checks if the other elevators finishes orders within ticklength * tickTreshold seconds*/
func elevTicker() {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			_mtx.Lock()

			for key, _ := range heartbeat {
				heartbeat[key]++
				orderTicker[key]++
			}
			_mtx.Unlock()

		}
	}
}
