package orderHandler
import . "../config"

func costCalculation (order ButtonEvent, elevatorList [NumElevators]Elevator, id int, onlineElevators[NumElevators]bool) int{
	cheapest := 1000
	bestElevator := 0 //
	if order.Button == Btn_Cab {
		return id
	}

	for elev := 0; elev < NumElevators; elev++{
		cost := 1000
		if !onlineElevators[elev] {
			continue
		}
		if elevatorList[elev].Dir == Dir_Stop && order.Floor == elevatorList[elev].Floor {
			return elev



		}else if elevatorList[elev].Dir ==  Dir_Up && order.Floor > elevatorList[elev].Floor{
			if order.Button == Btn_Up {
				cost = order.Floor - elevatorList[elev].Floor
			}else{
				cost = order.Floor - elevatorList[elev].Floor + 2
			}


		}else if elevatorList[elev].Dir ==  Dir_Down && order.Floor < elevatorList[elev].Floor{

			if order.Button == Btn_Down {
				cost = elevatorList[elev].Floor - order.Floor
			}else{
				cost = elevatorList[elev].Floor - order.Floor + 2
			}

		}else if elevatorList[elev].Dir ==  Dir_Stop && order.Floor > elevatorList[elev].Floor{
			cost = order.Floor - elevatorList[elev].Floor +1

		}else if elevatorList[elev].Dir ==  Dir_Stop && order.Floor < elevatorList[elev].Floor{
			cost = elevatorList[elev].Floor - order.Floor +1

		}else if elevatorList[elev].Dir ==  Dir_Down && order.Floor > elevatorList[elev].Floor{
			cost = order.Floor - elevatorList[elev].Floor +3

		}else if elevatorList[elev].Dir ==  Dir_Up && order.Floor < elevatorList[elev].Floor{
			cost = elevatorList[elev].Floor - order.Floor +3
		}

		if cost < cheapest{
			cheapest = cost
			bestElevator = 	elev 
		}
	}
	return bestElevator
}

func orderAlreadyInQueue(order ButtonEvent, elevatorList [NumElevators]Elevator, id int) bool{
	if order.Button == Btn_Cab && elevatorList[id].Queue[order.Floor][Btn_Cab]{
		return true
	} else {
		for elev := 0; elev < NumElevators; elev++{
			if elevatorList[elev].Queue[order.Floor][order.Button] && order.Button != Btn_Cab{
				return true
			}
		}
	}
	return false

}
