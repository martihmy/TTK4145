package orderHandler

import (. "../config"
		hw "../hardware_io"
		"fmt"
		//"time"
		//sync "../elevSync"
)


func OrderHandler(id int, btnPressChan chan ButtonEvent, newOrderChan chan ButtonEvent, lightUpdaterChan chan [NumElevators]Elevator,
	elevatorChan chan Elevator, servicedFloorChan chan int, updateFromSync chan [NumElevators]Elevator, orderUpdateChan chan ButtonEvent, syncUpdateChan chan Elevator,
	updateGovOnlineListChan chan [NumElevators]bool){

	var(
		elevatorList [NumElevators]Elevator
		servicedOrder	ButtonEvent
		onlineElevators [NumElevators]bool
	)
	elevatorList[id] = <- elevatorChan
	//onlineElevators[id] = true
	syncUpdateChan <- elevatorList[id] //assures that all elevators in continously updated on direction, floor, orders and states

	for {
		select {
		case updateOnSyncList := <- updateFromSync: //Some update has occurred in sync so the updated elevatorList is sent here
			aNewOrder := false
			for elev:=0; elev < NumElevators; elev++ {
				if elev == id {								//Only check updates on other elevators
					continue
				}
				if elevatorList[elev].Queue != updateOnSyncList[elev].Queue { //if the incomming updated remote elevatorList has a different Queue then it contains a new order
					aNewOrder = true
				}
				elevatorList[elev] = updateOnSyncList[elev] //Replace our information on other elevators with the new updated version
			}
			for floor:=0; floor < NumFloors;floor++{
				for button:=Btn_Up; button < NumButtons; button++{
					if !elevatorList[id].Queue[floor][button] && updateOnSyncList[id].Queue[floor][button] { //Check if the new remote elevatorlist contains an order placed in our queue
						elevatorList[id].Queue[floor][button] = true
						aNewOrder = true
						newOrder := ButtonEvent{floor,button,id,false}
						go func(){newOrderChan <- newOrder} ()
					}
					//Do we need to check the other case? If the updateremote list contains a false in our matrix and our local elevatorList contains a true?
				}
			}
			if aNewOrder {
				lightUpdaterChan <- elevatorList
			}

		case updateOnLocalElev := <- elevatorChan:
			temp := elevatorList[id].Queue 						//Do we need to check for undefined state?
			elevatorList[id] = updateOnLocalElev
			elevatorList[id].Queue = temp
			if onlineElevators[id] {
				syncUpdateChan <- elevatorList[id]
			}


		case updateMyOnlineList := <- updateGovOnlineListChan:
			onlineElevators = updateMyOnlineList
			for i := 0; i<NumElevators;i++ {
				fmt.Println("OrderHandler has registered that elevator:",i,"is",onlineElevators[i])
			}


		case newLocalOrder := <- btnPressChan:
			if !onlineElevators[id] {
				elevatorList[id].Queue[newLocalOrder.Floor][newLocalOrder.Button] = true //byttet Button med Floor
				fmt.Println("New local call for for floor", newLocalOrder.Floor,"at Elevator:",id)
				lightUpdaterChan <- elevatorList
				go func() { newOrderChan <- newLocalOrder }()

			} else {
				if !orderAlreadyInQueue(newLocalOrder, elevatorList, id) {
					newLocalOrder.DesignatedID = costCalculation(newLocalOrder, elevatorList,id,onlineElevators)
					fmt.Println("Best elevator for order is:",newLocalOrder.DesignatedID)
					orderUpdateChan <- newLocalOrder
				}
			}


		case servicedOrder.Floor = <- servicedFloorChan:
			servicedOrder.Finished = true
			fmt.Println("Received serviced order in orderHandler")
			for button := Btn_Up; button < NumButtons; button++{
				if elevatorList[id].Queue[servicedOrder.Floor][button] {
					servicedOrder.Button = button
				}
				for elev := 0; elev < NumElevators;elev++{
					if button != Btn_Cab || elev == id {

						elevatorList[elev].Queue[servicedOrder.Floor][button] = false //Kan vi bare fjerne for oss selv her?
					}
				}
			}
			fmt.Println("If the elevator is online order on floor:",servicedOrder.Floor, "for btn:",servicedOrder.Button,"Should be sent to orderupdate")
			if onlineElevators[id] {
				orderUpdateChan <- servicedOrder
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
