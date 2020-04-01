package network
//-----------------------------------------------------------------------------------------------------------------------------------------------
/*This module conains the functions for sending and recieving messages on the network*/
//--------------------------------------------------------------------------------------------------------------------------------------------------
import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"syscall"
	"time"
)

var port string
var protocol string
var serverIP string

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Nettwork functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*JSON decodes message from network and adds message to param channel.*/
func RecieveMessage(port int, chans ...interface{}) {
	checkArgs(chans...)

	var buf [1024]byte
	conn := dialBroadcastUDP(port)
	for {
		time.Sleep(20 * time.Millisecond)
		//fmt.Println("In: Receive msg")
		n, _, _ := conn.ReadFrom(buf[0:])
		for _, ch := range chans {
			T := reflect.TypeOf(ch).Elem()
			typeName := T.String()
			//fmt.Printf(typeName)
			if strings.HasPrefix(string(buf[0:n])+"{", typeName) {
				v := reflect.New(T)
				//fmt.Println("Receiving:")
				//fmt.Println(v)
				json.Unmarshal(buf[len(typeName):n], v.Interface())
				reflect.Select([]reflect.SelectCase{{
					Dir:  reflect.SelectSend,
					Chan: reflect.ValueOf(ch),
					Send: reflect.Indirect(v),
				}})
			}
		}
	}
}

/* JSON encodes and brodcasts message from param channel to the network on param port*/
func BrodcastMessage(port int, chans ...interface{}) {

	//fmt.Printf("typeName")
	checkArgs(chans...)
	n := 0
	for range chans {
		n++
	}
	//fmt.Println("typeName")
	selectCases := make([]reflect.SelectCase, n)
	typeNames := make([]string, n)
	for i, ch := range chans {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
		typeNames[i] = reflect.TypeOf(ch).Elem().String()
	}

	conn := dialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	for {
		time.Sleep(20 * time.Millisecond)
		chosen, value, _ := reflect.Select(selectCases)
		/*fmt.Println("Sending:")
		fmt.Println(value)
		fmt.Println("_______________")*/
		buf, _ := json.Marshal(value.Interface())
		conn.WriteTo([]byte(typeNames[chosen]+string(buf)), addr)
	}

}




/*Sets up and returns connection*/
func dialBroadcastUDP(port int) net.PacketConn {
	s, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	syscall.Bind(s, &syscall.SockaddrInet4{Port: port})

	f := os.NewFile(uintptr(s), "")
	conn, _ := net.FilePacketConn(f)
	f.Close()

	return conn
}


// ------------------------------------------------------------------------------------------------------------------------------------------------------
// Support Functions
// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*Checks that args to Tx'er/Rx'er are valid on param channels*/
 func checkArgs(chans ...interface{}) {
	n := 0
	for range chans {
		n++
	}
	elemTypes := make([]reflect.Type, n)

	for i, ch := range chans {
		// Must be a channel
		if reflect.ValueOf(ch).Kind() != reflect.Chan {
			panic(fmt.Sprintf(
				"Argument must be a channel, got '%s' instead (arg#%d)",
				reflect.TypeOf(ch).String(), i+1))
		}

		elemType := reflect.TypeOf(ch).Elem()

		// Element type must not be repeated
		for j, e := range elemTypes {
			if e == elemType {
				panic(fmt.Sprintf(
					"All channels must have mutually different element types, arg#%d and arg#%d both have element type '%s'",
					j+1, i+1, e.String()))
			}
		}
		elemTypes[i] = elemType

		// Element type must be encodable with JSON
		switch elemType.Kind() {
		case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.UnsafePointer:
			panic(fmt.Sprintf(
				"Channel element type must be supported by JSON, got '%s' instead (arg#%d)",
				elemType.String(), i+1))
		case reflect.Map:
			if elemType.Key().Kind() != reflect.String {
				panic(fmt.Sprintf(
					"Channel element type must be supported by JSON, got '%s' instead (map keys must be 'string') (arg#%d)",
					elemType.String(), i+1))
			}
		}
	}
}