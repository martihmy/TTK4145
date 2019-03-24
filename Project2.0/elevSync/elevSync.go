package elevSync

import (
	//"fmt"
	"time"
	. "../config"
	//"math/rand"
	//"fmt"

)

type SyncChannels struct {
	OrderUpdate chan ButtonEvent
	UpdateOrderHandler	chan [NumElevators]Elevator
}


func ElevatorSynchronizer(ch SyncChannels, id int, newOrderChan chan ButtonEvent) {

	var (
		elevatorList [NumElevators]Elevator
		onlineElevators [NumElevators]bool
		change bool
	)

	Ack_Timer := time.NewTimer(2*time.Second)
	Ack_Timer.Stop()
	Fulfill_Timer := time.NewTimer(5*time.Second)
	Fulfill_Timer.Stop()



	for {
		select{


		case newOrder := <- ch.OrderUpdate:
			if newOrder.Finished {
				elevatorList[id].Queue[newOrder.Floor][newOrder.Button] = false
				change = true
				if newOrder.Button != Btn_Cab && change {
					//Handle Finished for orders recieved from others
				}
				//Update other elevators on our updated Queue
			} else if newOrder.Button == Btn_Cab && newOrder.DesignatedID == id {
				elevatorList[id].Queue[newOrder.Floor][newOrder.Button] = true   //This update will be broadcasted on the standard msg each time the broadcaster ticker goes off
				ch.UpdateOrderHandler <- elevatorList 
			} 

		case newElevatorUpdate:
		case receivedNewMsg:
		case newPeerUpdate := <- PeerUpdateChan:
			




		}
	}
}
