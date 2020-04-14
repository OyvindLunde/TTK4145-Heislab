# Elevator FSM

The finite state machine has 3 states and 3 events:

| States        | Events            |
| ------------- |:-----------------:|
| INIT          | NewOrder          |
| IDLE          | MotorDirection    |
| EXECUTE       | FloorReached      |


**INIT:** Initalizes the elevator and its modules, and runs the elevator down to ground floor.

**IDLE:** The elevator is standing still at a floor and is awaiting new orders.

**EXECUTE:** The elevator is currently executing an order.

**NewOrder:** The elevator has received a new order, and will check its validity and availability before executing the order.

**MotorDirection:** The elevator has gotten a new MotorDirection, and will start to run in the given direction.

**FloorReached:** The elevator has reached a floor, and will check whether it should stop or not.
