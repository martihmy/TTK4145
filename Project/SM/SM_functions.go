

func shouldStop(elevator Elev) bool {
	switch elevator.Dir {
	case Dir_Up:
		return elevator.Queue[elevator.Floor][BtnUp] ||
			elevator.Queue[elevator.Floor][BtnCab] ||
			!ordersAbove(elevator)
	case Dir_Down:
		return elevator.Queue[elevator.Floor][BtnDown] ||
			elevator.Queue[elevator.Floor][BtnCab] ||
			!ordersBelow(elevator)
	case Dir_Stop:
	default:
	}
	return false
}

func chooseDirection(elevator Elev) MotorDirection {
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

func ordersAbove(elevator Elev) bool {
	for floor := elevator.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevator Elev) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}