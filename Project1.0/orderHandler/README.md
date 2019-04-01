The orderHandler is responsible for handling all orders for a local elevator based on information on for the local elevator and other elevators supplied by the sync module.
For a local button presses it will calculate the best elevator and redistribute that order if it is an outside call and relay it to the synchronizer in order to update other elevators. All orders that is to be executed by the local elevator is sent to the state machine in order to move the actual elevator. Additionally, the orderHandler
is responsible for setting/turning of lights on the local elevator.
