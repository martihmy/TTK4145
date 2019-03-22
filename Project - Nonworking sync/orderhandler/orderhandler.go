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

for {
	select {
	case newOrder := <- btnPressChan:
		newOrder.Status = Uncomfirmed
		

				switch newOrder.Status{
				case Finished:
				case Empty:
				case Unconfirmed:
				case Confirmed
		}

		case newNetworkOrder := <- SomeNetworkChan




	}
}

}
