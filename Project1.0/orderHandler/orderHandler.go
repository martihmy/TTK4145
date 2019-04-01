package orderHandler

import (. "../config"
		hw "../hardware_io"

)


func OrderHandler(id int, btnPressChan chan ButtonEvent, newOrderChan chan ButtonEvent, lightUpdaterChan chan [NumElevators]Elevator,
	elevatorChan chan Elevator, servicedFloorChan chan int, updateFromSync chan [NumElevators]Elevator, orderUpdateChan chan ButtonEvent, syncUpdateChan chan Elevator,
	updateOrderOnlineListChan chan [NumElevators]bool){

	var(
		elevatorList [NumElevators]Elevator
		servicedOrder	ButtonEvent
		onlineElevators [NumElevators]bool
	)
	elevatorList[id] = <- elevatorChan
	syncUpdateChan <- elevatorList[id]

	for {
		select {
		case updateOnSyncList := <- updateFromSync:
			aNewOrder := false
			for elev:=0; elev < NumElevators; elev++ {
				if elev == id {
					continue
				}
				if elevatorList[elev].Queue != updateOnSyncList[elev].Queue {
					aNewOrder = true
				}
				elevatorList[elev] = updateOnSyncList[elev]
			}
			for floor:=0; floor < NumFloors;floor++{
				for button:=Btn_Up; button < NumButtons; button++{
					if !elevatorList[id].Queue[floor][button] && updateOnSyncList[id].Queue[floor][button] {
						elevatorList[id].Queue[floor][button] = true
						aNewOrder = true
						newOrder := ButtonEvent{floor,button,id,false}
						go func(){newOrderChan <- newOrder} ()
					}
				}
			}
			if aNewOrder {
				lightUpdaterChan <- elevatorList
			}

		case updateOnLocalElev := <- elevatorChan:
			temp := elevatorList[id].Queue
			elevatorList[id] = updateOnLocalElev
			elevatorList[id].Queue = temp
			if onlineElevators[id] {
				syncUpdateChan <- elevatorList[id]
			}


		case updateMyOnlineList := <- updateOrderOnlineListChan:
			onlineElevators = updateMyOnlineList
			for i := 0; i<NumElevators;i++ {
			}


		case newLocalOrder := <- btnPressChan:
			if !onlineElevators[id] {
				elevatorList[id].Queue[newLocalOrder.Floor][newLocalOrder.Button] = true
				lightUpdaterChan <- elevatorList
				go func() { newOrderChan <- newLocalOrder }()

			} else {
				if !orderAlreadyInQueue(newLocalOrder, elevatorList, id) {
					newLocalOrder.DesignatedID = costCalculation(newLocalOrder, elevatorList,id,onlineElevators)
					orderUpdateChan <- newLocalOrder
				}
			}


		case servicedOrder.Floor = <- servicedFloorChan:
			servicedOrder.Finished = true
			for button := Btn_Up; button < NumButtons; button++{
				if elevatorList[id].Queue[servicedOrder.Floor][button] {
					servicedOrder.Button = button
				}
				for elev := 0; elev < NumElevators;elev++{
					if button != Btn_Cab || elev == id {

						elevatorList[elev].Queue[servicedOrder.Floor][button] = false
					}
				}
			}

			if onlineElevators[id] {
				orderUpdateChan <- servicedOrder
			} else {
				orderUpdateChan <- servicedOrder
			}
			lightUpdaterChan <- elevatorList
		}
	}
}

func LightUpdater (lightUpdaterChan <-chan [NumElevators]Elevator, id int){
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
