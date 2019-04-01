##Communication module
This module handles all communication with the other independent elevators of the system.

## Notes
The heartbeat broadcast operates on port 16569.

The communication broadcast operates on port 16570

The module checks that the message received is not a duplicate of the last one sent and sends the message on where it's supposed to go.

In order to make sure that the message arrives we simply send the message several times.





The heartbeat module (peers) is a premade module made by Anders. This peers module consistently stops the program (effectively crashing it, but leaving out an error message saying what's wrong) after a certain amount of packets received. The same module worked perfectly for others, but not for us for some reason. 
