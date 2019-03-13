package governor
import . "../config"

func costCalculation (order ButtonEvent, elevatorList [NumElevators]Elevator){
	if order.Btn == BtnCab {
		return id //***********************'''?? number between 0 and 2??
	}
	cheapest := 1000
	bestElevator := 0 // 
	for elev := 0; elev < NumElevators; elev++{
		if elevatorList[elev].Dir ==  Dir_Up && order.Floor > elevatorList[elev].Floor{ //You are going up and the order is above you
			if order.Btn == BtnUp {
				cost := order.Floor - elevatorList[elev].Floor //+ They want to go in the same direction as you are already going
			}else{
				cost := order.Floor - elevatorList[elev].Floor + 2 // They do not want to go in the same direction
			}


		}else if elevatorList[elev].Dir ==  Dir_Down && order.Floor < elevatorList[elev].Floor{ //You are going down and the order is below you
			
			if order.Btn == BtnDown {
				cost := elevatorList[elev].Floor - order.Floor // + They want to go in the same direction as you are already going
			}else{
				cost := elevatorList[elev].Floor - order.Floor + 2// They do not want to go in the same direction
			}

			
		}else if elevatorList[elev].Dir ==  Dir_Stop && order.Floor > elevatorList[elev].Floor{ //You are standing still and the order is above you
			cost := order.Floor - elevatorList[elev].Floor +1
			
		}else if elevatorList[elev].Dir ==  Dir_Stop && order.Floor < elevatorList[elev].Floor{ //You are standing still and the order is below you
			cost := elevatorList[elev].Floor - order.Floor +1

		}else if elevatorList[elev].Dir ==  Dir_Down && order.Floor > elevatorList[elev].Floor{ //You are going down and the order is above you
			cost := order.Floor - elevatorList[elev].Floor +3

		}else if elevatorList[elev].Dir ==  Dir_Up && order.Floor < elevatorList[elev].Floor{ //You are going up and the order is below you
			cost := elevatorList[elev].Floor - order.Floor +3
		}

		if cost < cheapest{
			cheapest = cost 
			bestElevator = 	elev //number (0,1 or 2) of the elevator with the lowest cost
		}
	}
	return bestElevator
}

func orderAlreadyInQueue(order ButtonEvent, elevatorList [NumElevators]Elevator){ // **************** Still not sure about the use of id
	if order.Btn == BtnCab && elevatorList[id].Queue[order.Floor][BtnCab]{
		return true
	} else {
		for elev := 0; elev < NumElevators; elev++{
			if elevatorList[elev].Queue[order.Floor][order.Btn]{
				return true
			}
		}
	} 
	return false

}