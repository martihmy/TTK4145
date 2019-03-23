package governor

import (. "../config"
		hw "../hardware_io"
		"fmt"
		//"time"
		//sync "../elevSync"
)


func ElevGovernor(id int, btnPressChan chan ButtonEvent, newOrderChan chan ButtonEvent, lightUpdaterChan chan [NumElevators]Elevator,
	elevatorChan chan Elevator, servicedFloorChan chan int, sendOrderChan chan ButtonEvent, syncUpdateChan chan Elevator, updateGovChan chan Elevator){

	var(
		elevatorList [NumElevators]Elevator
		servicedOrder	ButtonEvent
	)
	elevatorList[id] = <- elevatorChan
	syncUpdateChan <- elevatorList[id] //assures that all elevators in continously updated on direction, floor, orders and states

	for {
		select {
		case updateOnOtherElev := <- updateGovChan:
			elevatorList[updateOnOtherElev.ID]=updateOnOtherElev
			for elev:= 0;elev < NumElevators;elev++ {
				if elev == id {
					continue
				}
					for floor:=0; floor < NumFloors; floor++{
						for btn := Btn_Up; btn < NumButtons; btn++ {
							fmt.Println("State of Elevator:",elev,"is", elevatorList[elev].State, "Dir:",elevatorList[elev].Dir,"Floor:",elevatorList[elev].Floor)
						}
					}
				}


		case updateOnThisElev := <- elevatorChan:
			elevatorList[id]=updateOnThisElev
			fmt.Println("State of current Elevator:", elevatorList[id].State, "Dir:",elevatorList[id].Dir,"Floor:",elevatorList[id].Floor )

/*		for elev:= 0;elev < NumElevators;elev++ {
			for floor:=0; floor < NumFloors; floor++{
				for btn := Btn_Up; btn < NumButtons; btn++ {
					fmt.Println("Queue on recieved order",elevatorList[id].Queue[floor][btn])
				}
			}
		}*/
		case newLocalOrder := <- btnPressChan:
			newLocalOrder.DesignatedID = id
			newLocalOrder.OrderID = 0
			if newLocalOrder.Button == Btn_Cab{
				elevatorList[id].Queue[newLocalOrder.Floor][Btn_Cab] = true //byttet Button med Floor
				fmt.Println("New local Cabcall for for floor", newLocalOrder.Floor,)
				lightUpdaterChan <- elevatorList
				for floor:=0; floor < NumFloors; floor++{
					for btn := Btn_Up; btn < NumButtons; btn++ {
						fmt.Println("Queue on recieved order",elevatorList[id].Queue[floor][btn])
					}
				}
				go func() { newOrderChan <- newLocalOrder }()
			} else {
				if !orderAlreadyInQueue(newLocalOrder, elevatorList, id) {
					newLocalOrder.DesignatedID = costCalculation(newLocalOrder, elevatorList)
					fmt.Println("Best elevator for order is:",newLocalOrder.DesignatedID)
					sendOrderChan <- newLocalOrder

				/*	if newLocalOrder.DesignatedID == id {
						sendOrderChan <- newLocalOrder
						time.Sleep(300*time.Millisecond)
						newOrderChan <- newLocalOrder
					} else {
						sendOrderChan <- newLocalOrder
					}
					//orderUpdate <- newLocalOrder
*/

				}
			}


				//Might want to make this a go function
			//Need to set lights at some point
			//orderUpdate <- newLocalOrder
		case servicedOrder.Floor = <- servicedFloorChan:
			//Do we need a Done in ButtonEvent struct?
			for button := Btn_Up; button < NumButtons; button++{
				if elevatorList[id].Queue[servicedOrder.Floor][button] {
					servicedOrder.Button = button
				}
				for elev := 0; elev < NumElevators;elev++{
					if button != Btn_Cab || elev == id {
						//fmt.Println("Order for elevator:",id,"On floor", servicedOrder.Floor, "for Button", button, "has been serviced")
						elevatorList[elev].Queue[servicedOrder.Floor][button] = false
					}
				}
					}
			for floor:=0; floor < NumFloors; floor++{
				for btn := Btn_Up; btn < NumButtons; btn++ {
					fmt.Println("Queue after completed order for elevator:",id,elevatorList[id].Queue[floor][btn])
				}
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
