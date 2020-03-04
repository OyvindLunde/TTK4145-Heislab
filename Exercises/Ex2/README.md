# Mutex and Channel basics

### What is an atomic operation?
> Only one action can happen at once, e.g a bank transfer

### What is a semaphore?
> Integer that locks/unlocks threads depending on whether it is incremented/decremneted. If a semaphore is 
decremented and the result is negative -> thread is locked. Incrementing a semaphore wakes an arbitrary locked thread.

### What is a mutex?
> Mutual exclusive, only one thread can perform a section at once.

### What is the difference between a mutex and a binary semaphore?
> Mutex is exclusive access to a resource.
Binary semaphore is used for synchronization, the binary "giver" simply notifies whoever the "taker" is, that what they
were waiting for happened.

### What is a critical section?
> Section usually within a mutex, something that only one thread can do at once.

### What is the difference between race conditions and data races?
 > Race condition: situation in which the result of an operation depneds on the interleaving of certain individual operations.
 Data race: situation in which at least two threads access a shared variable at the same time. At least one thread tries to
modify the variable.

### List some advantages of using message passing over lock-based synchronization primitives.
> Message passing is more scalable. A shared object can be duplicated. State of mutable/shared are harder to reason about
where multiple threads can run concurrently. 

### List some advantages of using lock-based synchronization primitives over message passing.
> Simpler algorithms. A message passing system that requires resources to be locked will eventually degrade to shared variable.
If algorithms are wait-free, you will see a improved performance and reduced memory footprint (less object allocation in form of new messages).
