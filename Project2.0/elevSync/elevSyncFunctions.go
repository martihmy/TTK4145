package elevSync

import (
	//"fmt"
	//"time"
	. "../config"
	//"math/rand"
	//"fmt"

)
func allOnSamePage(floor int, btn ButtonType, AllAckTest bool, AllUnackTest bool, NoneDoneTest bool, onlineElevators[NumElevators]bool, message Msg)bool{ // er alle unfinished, er alle unacked, er alle acked??
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
