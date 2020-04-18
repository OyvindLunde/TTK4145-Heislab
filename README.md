# Superheis

This Project has been developed by:
```
Gulleik L Olsen         gulleik@hotmail.com      Software Engineer @ Sanntid Supersquad
Jens E Walmsness        jensemil97@gmail.com     Software Engineer @ Sanntid Supersquad
Øyvind R Lunde          oylunrl@gmail.com        Software Engineer @ Sanntid Supersquad
```

# Elevator project
**Software for controlling n elevators with m floors** Evt kort intro

- Video av kjørende system med display 

## System requirements
- Add requirements
- Linux/ubuntu
- Simulator v2
- Husk display package

## Project description

## Content (evt annet navn)
- Brief explanation of our system, its modules and how they are connected


### Elevator FSM

The finite state machine has 3 states and 3 events:

| States        | Events            |
|:-------------:|:-----------------:|
| IDLE          | NewOrder          |
| EXECUTE       | MotorDirection    |
| RESET         | FloorReached      |
|               | OrderTimeout      |



**IDLE:** The elevator is standing still at a floor and is awaiting new orders.

**EXECUTE:** The elevator is currently executing an order.

**RESET:** The elevator has timed out whle executing an order, and needs to be reset.

**NewOrder:** The elevator has received a new order, and will check its validity and availability before executing the order.

**MotorDirection:** The elevator has gotten a new MotorDirection, and will start to run in the given direction.

**FloorReached:** The elevator has reached a floor, and will check whether it should stop or not.

**OrderTimeout**


## Display

