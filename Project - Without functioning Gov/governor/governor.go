package governor

import (. "../config"
	"fmt"
		//hw "../hardware_io"
)
 type GovernorChannels struct {
 	ButtonPress chan ButtonEvent
 }

func ElevGovernor(/*id int,*/initialFloor int, elevatorChan chan Elevator, orderCompleteChan chan ButtonEvent, ch GovernorChannels){
	/*var (
		elevatorList	[NumElevators]Elevator
		//onlineList		[NumElevators]bool

	)*/
	elevator := Elevator{ //Initialization of elevator object
		State: Idle,
		Dir: Dir_Stop,
		Floor: initialFloor, //  <- hw.PollFloorSensor(ch.FloorChan), //Should perhaps be initialized with some function in hw or at least use channals to get floor signal
		Queue: [NumFloors][NumButtons]bool{},
		//ID: 1 //Maybe assign ID from some online 
	}

	//do some update to the sync
	switch {
	case newLocalOrder <- ch.ButtonPress:
		fmt.Println("New order registered")
		if newLocalOrder.Button == BtnCab {

			elevator.Queue[newLocalOrder.Floor][newLocalOrder.Button] = true
			for floor := 0; floor < 4; floor++ {
				for btn := BtnUp; btn < NumButtons; btn++ {
					fmt.Println(elevator.Queue[floor][btn])
				}
			}

			go func() {elevatorChan <- elevator}()

			//fmt.Println(elevator.Queue[newLocalOrder.Floor][newLocalOrder.Button])
			fmt.Println("Turn on lights")

			orderComplete := <- orderCompleteChan

			elevator.Queue[orderComplete.Floor][orderComplete.Button] = false
			elevator.Queue[orderComplete.Floor][BtnCab] = false
			fmt.Println("Turn of lights")
			for floor := 0; floor < 4; floor++ {
				for btn := BtnUp; btn < NumButtons; btn++ {
					fmt.Println(elevator.Queue[floor][btn])
				}
			}
			//orderComplete <- orderCompleteChan


		}
	}
}