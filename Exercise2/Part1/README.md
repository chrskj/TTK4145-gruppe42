# Mutex and Channel basics

### What is an atomic operation?
> an atomic operation is one which cannot be (or is not) interrupted by concurrent operations

### What is a semaphore?
> A variable which controls the access to a common resource. It limits the number of consumers to a specific resource.

### What is a mutex?
> A binary variable which locks a common resource. When a process releases it, another can request access to the resource.

### What is the difference between a mutex and a binary semaphore?
> A binary semaphore flags a resource as taken without protecting it. A mutex gives exclusive access to a resource for the user.

### What is a critical section?
> Parts of the program where a shared resource is accessed. This part is protected with only one process being able to access it at once.

### What is the difference between race conditions and data races?
 > A race condition occurs when the timing or ordering of events affect the program's correctness. A data race happens when two concurrent threads write to the same location in memory.

### List some advantages of using message passing over lock-based synchronization primitives.
> 

### List some advantages of using lock-based synchronization primitives over message passing.
> 
