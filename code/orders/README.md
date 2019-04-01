##Module for handling of orders
This module takes care of handling all incoming and outgoing orders. It distributes orders based on received cost functions from the other elevators. It also saves orders to file to ensure correct behaviour of the elevator after a power outage.

Initialisation includes:

- Getting ID
- Getting the elevator's orders from file
- reinitialize cab orders

For-select loop:

- is the main "while" loop that the function utilises to check for updates/messages from other modules.
- At the top of the loop there is a function, this is there (instead of outside the OrderRoutine function) so that it can use the function-wide variables in OrderRoutine.
- When receiving a message, it executes a behaviour based on the PacketType.

CostCompare:

- Asks for and compares costs from different elevators and assigns the order to the elevator with the lowest cost function.

readFile reads the orders[ID].csv file and saves the data to localOrders.


writeToFile writes the localOrders data to file

 
