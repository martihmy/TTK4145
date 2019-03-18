package governor

import (. "../config"
		hw "../hardware_io"
		"fmt"
)


func ElevGovernor(id int, btnPressChan chan ButtonEvent, newOrderChan chan ButtonEvent, lightUpdaterChan chan [NumElevators]Elevator, 
	elevatorChan chan Elevator, servicedFloorChan chan int){

	var(
		elevatorList [NumElevators]Elevator
		servicedOrder	ButtonEvent
	)
	elevatorList[id] = <- elevatorChan

	for {
		select {
		case newLocalOrder := <- btnPressChan:
			newLocalOrder.DesignatedID = id
			if newLocalOrder.Button == Btn_Cab{
				elevatorList[id].Queue[newLocalOrder.Floor][Btn_Cab] = true //byttet Button med Floor
				lightUpdaterChan <- elevatorList
				fmt.Println("Enters governor")
				go func() { newOrderChan <- newLocalOrder }()
			} else {
				if !orderAlreadyInQueue(newLocalOrder, elevatorList, id) {
					newLocalOrder.DesignatedID = costCalculation(newLocalOrder, elevatorList)
					//orderUpdate <- newLocalOrder

				}
			}
		case servicedOrder.Floor = <- servicedFloorChan:
			//Do we need a Done in ButtonEvent struct?
			for button := Btn_Up; button < NumButtons; button++{
				if elevatorList[id].Queue[servicedOrder.Floor][button] {
					servicedOrder.Button = button
				}
				for elev := 0; elev < NumElevators;elev++{
					if button != Btn_Cab || elev == id {
						elevatorList[elev].Queue[servicedOrder.Floor][button] = false
					}
				}
				//Send servicedOrder to update sync
			}
			lightUpdaterChan <- elevatorList //Update lights, assumes that the servicedOrder sent to sync has managed to update the other elevators 
		}
	}
}

func LightUpdater (lightUpdaterChan <- chan [NumElevators]Elevator, id int){
	var orderPlaced [NumElevators]bool
	for{
		elevators := <- lightUpdaterChan
		for floor:= 0; floor < NumFloors; floor++{
			for button := Btn_Up; button < NumButtons; button++{
				for elev := 0; elev < NumElevators; elev++{
					orderPlaced[elev] = false
					if elev != id && button == Btn_Cab{
						continue

					}
					if elevators[elev].Queue[floor][button]{
						hw.SetButtonLamp(button,floor,true)
						orderPlaced[elev] = true

					}

				} 
				if orderPlaced == [NumElevators]bool{}{
					hw.SetButtonLamp(button, floor, false)
				}

			}
		}
	}
}