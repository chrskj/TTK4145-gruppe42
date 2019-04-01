##Elevator state machine
In our state machine we have the following states:
-Idle
-Running
-DoorOpen
-Emergency stop

We also have several events, that triggers specific behaviour
-Recieved order from Order-module
-Recieved request for cost function
-Button pushed
-Floor reached
-engineWatchdog timeout
-doorTimer timeout