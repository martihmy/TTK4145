package governor

import (. "../config"
		//hw "../hardware_io"
)

func ElevGovernor(btnPressChan <-chan ButtonEvent, newOrderChan chan<- ButtonEvent){
	newLocalOrder := <- BtnPressChan
	go func() { newOrderChan <- newLocalOrder }()
}