# Reasons for concurrency and parallelism

 ### What is concurrency? What is parallelism? What's the difference?
 > Concurrency means multiple operiations are starting, running, and completing in overlapping time periods, in no specific order.
Parallelism is when multiple tasks or different parts of the same task literally run at the same time.
 
 ### Why have machines become increasingly multicore in the past decade?
 > Because the clock rate of processors have an upper bound due to heat. It is cheaper, as well as more efficient, to increase the amount
 of cores instead of pushing this upper bound.

 
 ### What kinds of problems motivates the need for concurrent execution?
 > The majority of the current systems are multicore, i.e needing concurrency.In addition concurrency helps for real time problems,
 e.g multiple elevators communicating, different timing problems etc. Furthermore concurrency overall makes the system/program go faster.

 
 ### Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
 > The programmer's life will probably be harder, as the tradeoff
 for the benfits listed above is increased complexity, as shared resources are harder to manage.

 
 ### What are the differences between processes, threads, green threads, and coroutines?
 > - Processes: independent sequences of execution, runs in separate memory spaces.
 > - Threads: independent sequences of execution run in a shared memory space.
 > - Green threads: threads scheduled by a runtime library or virtual machine instead of the native underlying OS.
 > - Coroutines: general structure where flow control is passed between different threads/routines without returning.
      I.e the coroutines are collaborative; only one coroutine is working
      at a given time and they work in tandem with shared global and other information.

 
 ### Which one of these do `pthread_create()` (C/POSIX), `threading.Thread()` (Python), `go` (Go) create?
 > - pthread_create(): creates thread
 > - threading.Thread(): creates thread
 > - go: creates green thread
 
 ### How does pythons Global Interpreter Lock (GIL) influence the way a python Thread behaves?
 > The GIL is a mutex; it allows only one thread to hold the control of the python interpreter,
meaning only one thread can be in a state of execution at any given time. It basically bottlenecks multithreading.

 
 ### With this in mind: What is the workaround for the GIL (Hint: it's another module)?
 > Create multiple processes instead of multiple threads (multiprocessing module).

 
 ### What does `func GOMAXPROCS(n int) int` change? 
 > Limits the amount of (operating system) threads that can execute user-level Go code simultaneously.

