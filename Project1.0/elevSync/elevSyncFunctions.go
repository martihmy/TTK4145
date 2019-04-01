package elevSync

import . "../config"

func allOnSamePage(floor int, btn ButtonType, AllAckTest bool, AllUnackTest bool, NoneDoneTest bool, onlineElevators[NumElevators]bool, message Msg)bool{
	if AllAckTest{
		for elev:=0;elev<NumElevators;elev++{
			if message.StatusMatrix[floor][btn].AckList[elev] == 0 && onlineElevators[elev]{
				return false

			}
		}

	}else if AllUnackTest{
		for elev:=0;elev<NumElevators;elev++{
			if message.StatusMatrix[floor][btn].AckList[elev] == 1 && onlineElevators[elev]{
				return false

			}
		}

	}else if NoneDoneTest{
		for elev:=0;elev<NumElevators;elev++{
			if message.StatusMatrix[floor][btn].DoneList[elev] == 1 && onlineElevators[elev]{
				return false

			}
		}
	}
	return true
}

func AllAckOnLocal(floor int, btn ButtonType, onlineElevators[NumElevators]bool, localMatrix [NumFloors][NumButtons-1]OrderInfo)bool{
	for elev:=0;elev<NumElevators;elev++{
			if localMatrix[floor][btn].AckList[elev] == 0 && onlineElevators[elev]{
				return false

			}
	}
	return true
}
