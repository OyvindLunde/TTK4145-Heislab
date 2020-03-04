package main

import (
	"encoding/binary"
	. "fmt"
	"log"
	"net"
	"os/exec"
	t "time"
)

var counter uint64
var port = 1111
var buffer = make([]byte, 16)

func spawnBackup() {
	(exec.Command("gnome-terminal", "-x", "sh", "-c", "go run ~/SuperHeis/ex6/Sanntidsprogrammering/Ex6/ex6.go")).Run()

	Println("Backup spawned")
}

func main() {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println("Error")
	}
	imPrimary := false

	log.Println("Backup here")

	//backup
	for !(imPrimary) {
		conn.SetReadDeadline(t.Now().Add(2 * t.Second))
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			imPrimary = true
		} else {
			counter = binary.BigEndian.Uint64(buffer[:n])
		}
	}
	conn.Close()

	spawnBackup()
	Println("Primary here")
	broadcaster, _ := net.DialUDP("udp", nil, addr)
	//primary
	for {
		println("Count:", counter)
		counter++
		binary.BigEndian.PutUint64(buffer, counter)
		_, _ = broadcaster.Write(buffer)
		t.Sleep(100 * t.Millisecond)
	}
}