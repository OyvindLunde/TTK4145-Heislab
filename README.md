# Superheis

This Project has been developed by:
```
Gulleik L Olsen         gulleik@hotmail.com      Software Engineer @ Sanntid Supersquad
Jens E Walmsness        jensemil97@gmail.com     Software Engineer @ Sanntid Supersquad
Øyvind R Lunde          oylunrl@gmail.com        Software Engineer @ Sanntid Supersquad
```

# Elevator project
![Elevator Recording](Media/ElevatorRecording.gif)

## Project description
**Create software for controlling `n` elevators working in parallel across `m` floors.**  
We were free to implement our solution however we wanted, but there were some system requirements which had to be fulfilled:

  - **No orders are lost**  
     Once the light on a hall call button is turned on, an elevator should arrive at that floor.
  - **Multiple elevators should be more efficient than one**  
     Similarly for a cab call, but only the elevator at that specific workspace should take the order.
  - **An Inidividual elevator should behave sensibly and efficiently**  
     This means handling network packet loss, losing network connection entirely, software that crashes, and losing power - both to the elevator motor and the machine that controls the elevator.
  - **The lights should function as expected**  
     If the elevator is disconnected from the network, it should still serve all the currently active orders (ie. whatever lights are showing)

## Usage
**Prerequisites for running the program**
- Go
- Linux/Ubuntu
- Elevator simulator  - https://github.com/TTK4145/Simulator-v2
- Package for Display module - https://github.com/golang/exp/tree/master/shiny/driver

**Adding and running elevators**  
For each elevator do the following steps.

1. Open a terminal and run SimElevatorServer, specify a 5 digit port  
   `./SimElevatorServer --port xxxxx`
2. Open another terminal and run main.go  
   `go run main.go`
3. Enter an ID for the elevator and set the adress equal to the port specified in SimElevatorServer  
   `Enter Id: x`  
   `Enter Address: xxxxx`  
   
**NOTE:** Each elevator has to have a unique ID and port/address.  

Hotkeys for controlling the elevator can be found at https://github.com/TTK4145/Simulator-v2

## Our system

### Main
Main() sets the id and addr parameters, initializes the system and creates all channels and starts all goroutines.

### Elevator FSM

The finite state machine executes the given order, and switches between the different states based on the generated events. It has 3 states with 4 events:


| States        | Events            |
|:-------------:|:-----------------:|
| IDLE          | NewOrder          |
| EXECUTE       | MotorDirection    |
| RESET         | FloorReached      |
|               | OrderTimeout      |



  - **IDLE:** The elevator is standing still at a floor and is awaiting new orders.

  - **EXECUTE:** The elevator is currently executing an order.

  - **RESET:** The elevator has timed out while executing an order, and needs to be reset.

  - **NewOrder:** The elevator has received a new order, and will check its validity and availability before executing the order.

  - **MotorDirection:** The elevator has gotten a new MotorDirection, and will start to run in the given direction.

  - **FloorReached:** The elevator has reached a floor, and will check whether it should stop or not.

  - **OrderTimeout:** If the elevator had taken to long time to finish and order, typically because of motor power loss, an OrderTimeout event is generated to tell the system to reset itself.

### OrderHandler
The OrderHandler module controls all the logic for the elevator, such as deciding if the elevator should take an order, which direction it should go in, and whether or not it should stop when reaching a floor. It also handles incoming orders from button presses. 

### Communication
Communication transmits to and receives data from the network, and handles the received information.

### Network
The most basic functions for transmitting and receiving data.  
Mainly delivered code, with some changes.

### LogManagement
Logmanagement stores all information about both the local elevator and other elevators on the network, such as their state, orders and current floor.

### Ticker
A basic timer module.

### ElevController
Module containing basic functions for controlling the elevator.

### ElevIO
The most low-level functionality for interacting with the hardware.  
Mainly delivered code.

### Display
Due to the coronavirus we programmed and implemented a display that showed all connected elevators, as well as their state, orders and other relevant information. This made it much easier to see what was going on when running the system and debug it when errors occured.
See the README file in the Source/Display folder for a more thourough explanation of this module.
