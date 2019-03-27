package sm

import (
	. "../config"
	"time"
	//"fmt" //brukes bare til print osv
	hw "../hardware_io"
)

type SMChannels struct {
	NewOrder		chan ButtonEvent
	Elevator 		chan Elevator
	FloorArrival	chan int
	ServicedFloor	chan int
}





func ElevatorRun(ch SMChannels, initialFloor int, id int) {
	elevator := Elevator{ //Initialization of elevator object
		State: Idle,
		Dir: Dir_Stop,
		Floor: initialFloor, //  <- hw.PollFloorSensor(ch.FloorChan), //Should perhaps be initialized with some function in hw or at least use channals to get floor signal
		Queue: [NumFloors][NumButtons]bool{},
		ID: id,
	}
	doorOpenTimer := time.NewTimer(3*time.Second)
	doorOpenTimer.Stop()
	ch.Elevator <- elevator //Update gov with initilized struct


	for{
		select{
		case newOrder := <- ch.NewOrder:
			elevator.Queue[newOrder.Floor][newOrder.Button] = true //Temporarily until we have sorted our Queue system in Governor**

			switch elevator.State{
			case Idle:
				elevator.Dir = chooseDirection(elevator)
				hw.SetMotorDirection(elevator.Dir)
				if elevator.Dir == Dir_Stop {
					elevator.State = DoorOpen
					hw.SetDoorOpenLamp(true)
					doorOpenTimer.Reset(3*time.Second)
					go func() {ch.ServicedFloor <- newOrder.Floor}()
					elevator.Queue[elevator.Floor] = [NumButtons]bool{} //-- Send message to governor on OrderComplete channal and ask to turn of all lights for that floor.
				} else {
					elevator.State = Moving
				}

			case Moving://Keep moving until arrived at floor

			case DoorOpen:
				if elevator.Floor == newOrder.Floor {
					doorOpenTimer.Reset(3*time.Second)
					go func() {ch.ServicedFloor <- newOrder.Floor}()
					elevator.Queue[elevator.Floor] = [NumButtons]bool{}

				}
			case Undefined:

			default:
				//fmt.Println("Some error has occurred")

			}

			ch.Elevator <- elevator //to update when change in state

		case elevator.Floor = <- ch.FloorArrival:
			//fmt.Println("Arrived at floor:",elevator.Floor)
			if shouldStop(elevator) {
		//		fmt.Println("Should stop")
				hw.SetDoorOpenLamp(true)
				elevator.State = DoorOpen
				hw.SetMotorDirection(Dir_Stop)
				doorOpenTimer.Reset(3*time.Second)
				elevator.Queue[elevator.Floor] = [NumButtons]bool{}
				go func() {ch.ServicedFloor <- elevator.Floor}() //-- Send message to governor on OrderComplete channal and ask to turn of lights
			//	fmt.Println("Floor:",elevator.Floor,"has been sent to orderHandler")
			}
			ch.Elevator <- elevator

		case <- doorOpenTimer.C:
			hw.SetDoorOpenLamp(false)
			elevator.Dir = chooseDirection(elevator)
			if elevator.Dir == Dir_Stop {
				elevator.State = Idle
			} else {
				elevator.State = Moving
				hw.SetMotorDirection(elevator.Dir)
			}
			ch.Elevator <- elevator //to update when change in state
		}
	}

}
