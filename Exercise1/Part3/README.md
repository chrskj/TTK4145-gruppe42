# Reasons for concurrency and parallelism


To complete this exercise you will have to use git. Create one or several commits that adds answers to the following questions and push it to your groups repository to complete the task.

When answering the questions, remember to use all the resources at your disposal. Asking the internet isn't a form of "cheating", it's a way of learning.

 ### What is concurrency? What is parallelism? What's the difference?
 > Concurrency means that progress is made on several tasks before one is completed. The tasks doesn't neccesarily have to be worked on at the same time. Parallelism means that a task is split into smaller subtasks which are worked on at the same time.
 
 ### Why have machines become increasingly multicore in the past decade?
 > Increasingly faster CPU speeds required higher power consumption for reasons such as heat losses and leakage current. It proved more economical to increase the number of cores instead of increasing speeds.
 
 ### What kinds of problems motivates the need for concurrent execution?
 (Or phrased differently: What problems do concurrency help in solving?)
 > Since multicore processors is becoming more prevalent, concurrent programming will let a programmer make use of more than one of the cores for a program.
 
 ### Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
 (Come back to this after you have worked on part 4 of this exercise)
 > Concurrent programming will make it harder to know in what order the program will run in. Therefore methods such as message passing and mutex are neccesary, but makes the programmer's life harder.
 
 ### What are the differences between processes, threads, green threads, and coroutines?
 > Process: Instance of a computer program that is being executed. Can consist of several threads
 > Thread: OS-managed execution of a small sequence of instruction.
 > Green Thread: Threads that are managed by a runtime library or virutal machine and exist in the user space.
 > Coroutine: Similar to threads but are cooperatively multitasked, which means they provide concurrency but not parallelism.
 
 ### Which one of these do `pthread_create()` (C/POSIX), `threading.Thread()` (Python), `go` (Go) create?
 > pthread_create() creates a thread, threading.Thread() creates a thread and go creates a coroutine.
 
 ### How does pythons Global Interpreter Lock (GIL) influence the way a python Thread behaves?
 > GIL is a mutex that protects access to python objects, preventing multiple threads from executing python bytecodes at once. This can become a bottleneck for multicore systems. 
 
 ### With this in mind: What is the workaround for the GIL (Hint: it's another module)?
 > The multiprocessing package offers both local and remote concurrency, effectively side-stepping the GIL.
 
 ### What does `func GOMAXPROCS(n int) int` change? 
 > The GOMAXPROCS variable limits the number of operating threads that can execute code simultaneously.
