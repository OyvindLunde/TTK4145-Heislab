# Exercise 9 - Scheduling

## Task 1

#### 1. Why do we assign priorities to tasks?
Hierarchy of importance. If more tasks asks for the same resources, the priorities determines which task will be executed.
#### 2. What features must a scheduler have for it to be usable for real-time systems?
A scheduler must be predictable in order to be usable in real-time systems.

## Task 2
Without priority inheritence

| Task\Time | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10| 11| 12| 13| 14|
|-----------|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---| 
| a | | | | |E| | | | | | |Q|V|E| | 
| b | | |E|V| |V|E|E|E| | | | | | | 
| c |E|Q| | | | | | | |Q|Q| | | |E|


With inheritence

| Task\Time | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10| 11| 12| 13| 14|
|-----------|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---| 
| a | | | | |E| | |Q| |V|E| | | | | 
| b | | |E|V| | | | |V| | |E|E|E| | 
| c |E|Q| | | |Q|Q| | | | | | | |E|

## Task 3

#### 1. What is priority inversion? What is unbounded priority inversion?
- Priority inversion: when a high priority task is indirectly blocked by a lower priority, effectively inverting their priorities. 
- Unbounded Priority inversion: A high priority task may end up waiting indefinitely to acquire the locked resource.

#### 2. Does priority inheritance avoid deadlocks?
- A deadlock is when a circular set of chains are blocking eachothers resources resulting in no task being executed. Priority inheritance does not avoid this

## Task 4

#### 1.
* The application is assumed to consist of a fixed set of tasks.
* All tasks are periodic, with known periods. 
* The tasks are completely independent of each other. 
* All system overheads, context-switching times and so on are ignored (i.e, assumed to have zero cost). 
* Deadlines equal to their periods.
* Constant worse-case execution times. 
* All tasks run on a single processor. 

Most of the assumptions seems reasonable. Completely independent tasks, tasks have constant worse-case execution times might be difficult to realize.

#### 2. 
Utilization test: 15/50 + 10/30 + 5/20 = 0.8833 >  0.7798  = 3*(2^(1/3) - 1)
The left hand side is bigger than the right hand side so the test fails and thus the task set might not be schedulable.

#### 3.
Task c:
* w0 = 5
* Rc = 5 < 20

Task b:
* w0 = 10
* w1 = 10 + ceil(10/20 * 5) = 15
* w2 = 10 + ceil(15/20 * 5) = 15
* Rb = 15 < 30

Task a:
* w0 = 15
* w1 = 15 + ceil(15/30 * 10) + ceil(15/20 * 5) = 30
* w2 = 15 + ceil(30/30 * 10) + ceil(30/20 * 5) = 35
* w3 = 45
* w4 = 50
* w5 = w4 = 50
* Ra = 50 <= 50

All tasks pass the response-time analysis, showing that the task set is schedulable. The reason why the utilization test says otherwise is because it is only a sufficient test, whereas the response-time analysis is both sufficient and necessary.
