package orderhandler
import(
	. "../config"
	"time"
	//"fmt" //brukes bare til print osv
	hw "../hardware_io"

	)

	var(
		elevatorList [NumElevators]Elevator
		servicedOrder	ButtonEvent
	)
	elevatorList[id] = <- elevatorChan
	updateSynchronizer <- elevatorList[id]


func orderHandler(...) {
	Ack_Timer := time.NewTimer(2*time.Second)
	Ack_Timer.Stop()
	Fulfill_Timer := time.NewTimer(5*time.Second)
	Fulfill_Timer.Stop()



	select {
		case newLocalOrder := <- btnPressChan:
			if newLocalOrder.Button == Btn_Cab {
				elevatorList[id].Queue[newLocalOrder.Floor][Btn_Cab] = true //byttet Button med Floor
				lightUpdaterChan <- elevatorList
				go func() { newOrderChan <- newLocalOrder }()
			} else {
				
			}


		switch newLocalOrder.Status{
		case Empty:
		case Unconfirmed:
		case Confirmed 
		}

		case newNetworkOrder := <- SomeNetworkChan 




	}
}