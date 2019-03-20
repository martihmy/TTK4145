package elevSync

import (
	//"fmt"
	"time"
	. "../config"
	"math/rand"

)

type SyncChannels struct {
	UpdateSynchronizer		chan Elevator
	RecieveFulfillmentTimer chan ButtonEvent
	SendFulfillmentTimer 	chan ButtonEvent
	LocalFulfillmentTimer	chan ButtonEvent
	OutgoingOrder 			chan ButtonEvent
	IncommingOrder			chan ButtonEvent
	OrderAcknowledged		chan ButtonEvent
	DoItMySelf				chan ButtonEvent
	IncomingUpdateMsg		chan Elevator
	OutgoingUpdateMsg 		chan Elevator
	SendOrder 				chan ButtonEvent
	OrderFulfilled 			chan ButtonEvent
}


func ElevatorSynchronizer(ch SyncChannels, id int, newOrderChan chan ButtonEvent) {

	var (
		elevatorList [NumElevators]Elevator
		//someUpdate bool
	)


	Ack_Timer := time.NewTimer(2*time.Second)
	Ack_Timer.Stop()
	Fulfill_Timer := time.NewTimer(5*time.Second)
	Fulfill_Timer.Stop()



	for {
		select{


		case newElevUpdate := <- ch.UpdateSynchronizer:		//Some update in Elevator has occurred (state, floor, dir, queue)
			ch.OutgoingUpdateMsg <- newElevUpdate				//Broadcast the change to other elevators, make sure to check that you do not send an empty queue if you have just been reinitialized

		case updateMsg := <- ch.IncomingUpdateMsg: 	//change in some other elevator struct has occurred
			elevatorList[updateMsg.ID] = updateMsg 	//Update your own elevatorList at the other elevators ID 
													//Send the change from the other elevator to SetupF



		case orderToBeFulfilled := <- ch.RecieveRemoteFulfillmentTimer: //From reciever
			orderId := orderToBeFulfilled.OrderID
			Fulfill_Timer.Reset(5*time.Second)
			for {
				select{
					case orderFulfilled := <- orderFulfilledChan: //From reciever
							if orderFulfilled.OrderId == orderId {
								Fulfill_Timer.Stop()
								break
							}
					case <- Fulfill_Timer.C:
						orderToBeFulfilled.DesignatedID = id
						break
					}
				}
			if orderToBeFulfilled.DesignatedID == id {
				ch.DoItMySelf	:= <-orderToBeFulfilled 
			}

		case orderToBeFulfilled := <- ch.LocalFulfillmentTimer:
			orderId := orderToBeFulfilled.OrderID
			Fulfill_Timer.Reset(5*time.Second)
			for {
				select{
					case orderFulfilled := <- orderFulfilledChan: //From reciever
							if orderFulfilled.OrderId == orderId {
								Fulfill_Timer.Stop()
								break
							}
					case <- Fulfill_Timer.C:
						orderToBeFulfilled.DesignatedID = id
						break
					}
				}
			if orderToBeFulfilled.DesignatedID == id {
				ch.DoItMySelf	:= <-orderToBeFulfilled 
			}

		case orderTakeOver := <- ch.DoItMySelf:
			ch.RemoteFulfillTimer <- orderTakeOver
			newOrderChan <- orderTakeOver

		
		case outGoingOrder := <- ch.SendOrder:			//Broadcast order to other elevators
			if outGoingOrder.DesignatedID != id { 	//Check if sendOrder.DesignatedID is not my id
				outGoingOrder.OrderID = rand.Intn(1000)
				orderId := outGoingOrder.OrderId
				ch.OutgoingOrder <- outGoingOrder 
				Ack_Timer.Reset(2*time.Second) 		//Start acknowledgement timer
				for {
					select{
						case orderAck := <- ch.OrderAcknowledged:
							if orderAck.OrderId == orderId {
								Ack_Timer.Stop()
								break
							}
						case <- Ack_Timer.C:
							outGoingOrder.DesignatedID = id
							break
					}
				}
			} 
			if outGoingOrder.DesignatedID == id {
				ch.OutGoingOrder <- outGoingOrder

			}

					//Recieve acknowledgement through a start fulfillmenttimer signal
					//Stop ack_timer and start fulfillment timer
			//Stop timer when fullfilment is comfirmed


		case newOrder := <- ch.IncomingOrder: 	//New order is recieved from master
			if newOrder.DesignatedID == id {	//Check if we are suppose to take this order (Desginated order id is our id)
				ch.OrderAcknowledged <- newOrder
				ch.SendRemoteFulfillTimer <- newOrder //asks other elevators to start their fulfill timers
				newOrderChan <- newOrder
												
												//When order is finished --> tell other elevators to stop their fulfillmenttimers (might have to be done from some other place)
												//Tell governor to turn of lights
			} else {
				ch.LocalFulfillmentTimer <- newOrder 
			}



		}
	}
}
