package governor
import . "../config"

func costCalculation (order ButtonEvent, elevatorList [NumElevators]Elevator) int{
	//Can consider to implement a check for buttoncab. should not be nessesary
	cheapest := 1000
	bestElevator := 0 // 
	for elev := 0; elev < NumElevators; elev++{
		cost := 1000
		if elevatorList[elev].Dir ==  Dir_Up && order.Floor > elevatorList[elev].Floor{ //You are going up and the order is above you
			if order.Button == Btn_Up {
				cost = order.Floor - elevatorList[elev].Floor //+ They want to go in the same direction as you are already going
			}else{
				cost = order.Floor - elevatorList[elev].Floor + 2 // They do not want to go in the same direction
			}


		}else if elevatorList[elev].Dir ==  Dir_Down && order.Floor < elevatorList[elev].Floor{ //You are going down and the order is below you
			
			if order.Button == Btn_Down {
				cost = elevatorList[elev].Floor - order.Floor // + They want to go in the same direction as you are already going
			}else{
				cost = elevatorList[elev].Floor - order.Floor + 2// They do not want to go in the same direction
			}

		}else if elevatorList[elev].Dir == Dir_Stop && order.Floor == elevatorList[elev].Floor { //Standing still and order is at the same floor
			return elev

		}else if elevatorList[elev].Dir ==  Dir_Stop && order.Floor > elevatorList[elev].Floor{ //You are standing still and the order is above you
			cost = order.Floor - elevatorList[elev].Floor +1
			
		}else if elevatorList[elev].Dir ==  Dir_Stop && order.Floor < elevatorList[elev].Floor{ //You are standing still and the order is below you
			cost = elevatorList[elev].Floor - order.Floor +1

		}else if elevatorList[elev].Dir ==  Dir_Down && order.Floor > elevatorList[elev].Floor{ //You are going down and the order is above you
			cost = order.Floor - elevatorList[elev].Floor +3

		}else if elevatorList[elev].Dir ==  Dir_Up && order.Floor < elevatorList[elev].Floor{ //You are going up and the order is below you
			cost = elevatorList[elev].Floor - order.Floor +3
		}

		if cost < cheapest{
			cheapest = cost 
			bestElevator = 	elev //number (0,1 or 2) of the elevator with the lowest cost
		}
	}
	return bestElevator
}

func orderAlreadyInQueue(order ButtonEvent, elevatorList [NumElevators]Elevator, id int) bool{ // **************** Still not sure about the use of id
	if order.Button == Btn_Cab && elevatorList[id].Queue[order.Floor][Btn_Cab]{
		return true
	} else {
		for elev := 0; elev < NumElevators; elev++{
			if elevatorList[elev].Queue[order.Floor][order.Button]{
				return true
			}
		}
	} 
	return false

}