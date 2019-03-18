package sm
import(. "../config"
)

func shouldStop(elevator Elevator) bool {
	switch elevator.Dir {
	case Dir_Up:
		return elevator.Queue[elevator.Floor][Btn_Up] ||
			elevator.Queue[elevator.Floor][Btn_Cab] ||
			!ordersAbove(elevator)
	case Dir_Down:
		return elevator.Queue[elevator.Floor][Btn_Down] ||
			elevator.Queue[elevator.Floor][Btn_Cab] ||
			!ordersBelow(elevator)
	case Dir_Stop:
	default:
	}
	return false
}

func chooseDirection(elevator Elevator) MotorDirection {
	switch elevator.Dir {
	case Dir_Stop:
		if ordersAbove(elevator) {
			return Dir_Up
		} else if ordersBelow(elevator) {
			return Dir_Down
		} else {
			return Dir_Stop
		}
	case Dir_Up:
		if ordersAbove(elevator) {
			return Dir_Up
		} else if ordersBelow(elevator) {
			return Dir_Down
		} else {
			return Dir_Stop
		}

	case Dir_Down:
		if ordersBelow(elevator) {
			return Dir_Down
		} else if ordersAbove(elevator) {
			return Dir_Up
		} else {
			return Dir_Stop
		}
	}
	return Dir_Stop
}

func ordersAbove(elevator Elevator) bool {
	for floor := elevator.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevator Elevator) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}