package sm

import (
	. "../config"
	"time"
	"fmt" //brukes bare til print osv
	hw "../hardware_io"
)

type SMChannels struct {
	NewOrder			chan ButtonEvent
	Elevator 			chan Elevator
	FloorArrival		chan int
	OrderComplete		chan ButtonEvent
}





func ElevatorRun(ch SMChannels, initialFloor int) {
	
	doorOpenTimer := time.NewTimer(3*time.Second)
	doorOpenTimer.Stop()

	elevator := Elevator{ //Initialization of elevator object
		State: Idle,
		Dir: Dir_Stop,
		Floor: initialFloor, //  <- hw.PollFloorSensor(ch.FloorChan), //Should perhaps be initialized with some function in hw or at least use channals to get floor signal
		Queue: [NumFloors][NumButtons]bool{},
		//ID: 1 //Maybe assign ID from some online 
	}
	//ch.Elevator <- elevator //ALL OF THESE HAS BEEN COMMENTED OUT BECAUSE THEY WERE BLOCKING THE SM FROM RUNNING
	for{
		select{
		case elevatorcopy := <- ch.Elevator:

		elevator.Queue = elevatorcopy.Queue
			//elevator.Queue[newOrder.Floor][newOrder.Button] = true //Temporarily until we have sorted our Queue system in Governor**
			//Do we have to do some checks for order complete?
			switch elevator.State{
			case Idle:
				fmt.Println("case idle")
				elevator.Dir = chooseDirection(elevator)
				fmt.Println(elevator.Dir)
				hw.SetMotorDirection(elevator.Dir)
				if elevator.Dir == Dir_Stop {
					elevator.State = DoorOpen
					fmt.Println("Should go to door open")
					hw.SetDoorOpenLamp(true)
					doorOpenTimer.Reset(3*time.Second)
					//go func() {ch.OrderComplete <- newOrder.Floor}() -- Send message to governor on OrderComplete channal and ask to turn of all lights for that floor. 
				} else {
					elevator.State = Moving
					fmt.Println(elevator.State)
				}

			case Moving://Keep moving until arrived at floor
				fmt.Println("Moving")

			case DoorOpen:
					//go func() {ch.OrderComplete <- newOrder.Floor}() -- Send message to governor on OrderComplete channal and ask to turn of all lights for that floor.

				
			//ch.Elevator <- elevator

			}

		case elevator.Floor = <- ch.FloorArrival:
			if shouldStop(elevator) {
				hw.SetDoorOpenLamp(true)
				elevator.State = DoorOpen
				DirCopy := elevator.Dir
				hw.SetMotorDirection(Dir_Stop)
				doorOpenTimer.Reset(3*time.Second)
				if DirCopy == Dir_Down {
					ch.OrderComplete <- ButtonEvent{elevator.Floor, BtnDown}
				} else {
					ch.OrderComplete <- ButtonEvent{elevator.Floor, BtnUp}
				}
				//elevator.Queue[elevator.Floor] = [NumButtons]bool{}//Removes floor from queue after finised. See ** above

				//go func() {ch.OrderComplete <- elevator.Floor}() -- Send message to governor on OrderComplete channal and ask to turn of lights
			}
			//ch.Elevator <- elevator

		case <- doorOpenTimer.C:
			hw.SetDoorOpenLamp(false)
			elevator.Dir = chooseDirection(elevator)
			if elevator.Dir == Dir_Stop {
				elevator.State = Idle
			} else {
				elevator.State = Moving
				hw.SetMotorDirection(elevator.Dir)
			}
			//ch.Elevator <- elevator
		}
	}

}
