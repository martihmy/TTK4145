package SM

import {
	"../Config"
	"time"
	"fmt"
	hw "../Hardware"

}

type SMChannels struct {
	NewOrder		chan ButtonEvent
	Elevator 		chan Elevator
	FloorArrival	chan int 
}




func elevatorRun(ch SMChannels) {
	elevator := Elevator{ //Initialization of elevator object
		State: Idle
		MotorDirection: Dir_Stop
		Floor: getFloor() //Should perhaps be initialized with some function in hw or at least use channals to get floor signal
		Queue: [NumFloors][NumButtons]bool{}
		ID: //Maybe assign ID from some online 
	}
	
	doorOpenTimer := time.NewTimer(3*time.Second)
	doorOpenTimer.Stop()

	ch.Elevator <- elevator

	for{
		select{
		case newOrder := <- ch.NewOrder:
			//Do we have to do some checks for order complete?
			switch elevator.State{
			case Idle:
				elevator.MotorDirection = chooseDirection(elevator)
				hw.SetMotorDirection(elevator.MotorDirection)
				if elevator.MotorDirection == Dir_Stop {
					elevator.State = DoorOpen
					hw.SetDoorOpenLamp(1)
					doorOpenTimer.Reset(3*time.Second)
					//go func() {ch.OrderComplete <- newOrder.Floor}() -- Send message to governor on OrderComplete channal and ask to turn of all lights for that floor. 
				} else {
					elevator.State = Moving
				}
			case Moving: //Keep moving until arrived at floor
			case DoorOpen:
				if elevator.Floor == newOrder.Floor {
					doorOpenTimer.Reset(3*time.Second)
					//go func() {ch.OrderComplete <- newOrder.Floor}() -- Send message to governor on OrderComplete channal and ask to turn of all lights for that floor.

				}
				}
			ch.Elevator <- elevator

			}
		case elevator.Floor = <- ch.FloorArrival:
			if shouldStop(elevator) {
				hw.SetDoorOpenLamp(1)
				elevator.State = DoorOpen
				hw.SetMotorDirection(Dir_stop)
				doorOpenTimer.Reset(3*time.Second)
				//go func() {ch.OrderComplete <- elevator.Floor}() -- Send message to governor on OrderComplete channal and ask to turn of lights
			}
			ch.Elevator <- elevator
			}
		case <- doorOpenTimer.C
			hw.SetDoorOpenLamp(0)
			elevator.Dir = chooseDirection(elevator)
			if elevator.Dir == Dir_Stop {
				elevator.State = Idle
			} else {
				elevator.State = Moving
				hw.SetMotorDirection(elevator.Dir)
			}
			ch.Elevator <- elevator
		}
	}

}
