package elevSync

import (
	//"fmt"
	"time"
	. "../config"
	"math/rand"
	"fmt"

)

type SyncChannels struct {
	UpdateSynchronizer		chan Elevator
	IncomingUpdateMsg		chan Elevator
	OutgoingUpdateMsg 		chan Elevator
	UpdateGov				chan Elevator
	OutgoingOrder 			chan ButtonEvent
	IncomingOrder			chan ButtonEvent
	DoItMySelf				chan ButtonEvent
	SendOrder 				chan ButtonEvent
	OrderAck				chan TimerMsg
	OrderFulfilled 			chan TimerMsg
	ReceiveTimerMsg 		chan TimerMsg
	SendTimerMsg			chan TimerMsg
}


func ElevatorSynchronizer(ch SyncChannels, id int, newOrderChan chan ButtonEvent) {

	var (
		elevatorList [NumElevators]Elevator
		orderID int
	//	someUpdate bool
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
			elevatorList[updateMsg.ID] = updateMsg
			ch.UpdateGov <- updateMsg 	//Update your own elevatorList at the other elevators ID
													//Send the change from the other elevator to SetupF



	/*	case startFulfilmentTimer := <- ch.ReceiveFulfillmentTimer: //From reciever
			orderId := startFulfilmentTimer.OrderID
			Fulfill_Timer.Reset(5*time.Second)
			for {
				select{
					case orderFulfilled := <- Ch.OrderFulfilled: //From reciever
							if orderFulfilled.OrderID == orderId && orderFulfilled.Fulfilled {
								Fulfill_Timer.Stop()
								break
							}
					case <- Fulfill_Timer.C:
						orderToBeFulfilled.DesignatedID = id
						break
					}
				}
			if orderToBeFulfilled.DesignatedID == id {
				ch.DoItMySelf <-orderToBeFulfilled
			}*/

/*
		case orderToBeFulfilled := <- ch.LocalFulfillmentTimer:
			orderId := orderToBeFulfilled.OrderID
			Fulfill_Timer.Reset(5*time.Second)
			for {
				select{
					case orderFulfilled := <- orderFulfilledChan: //From reciever
							if orderFulfilled.OrderID == orderId {
								Fulfill_Timer.Stop()
								break
							}
					case <- Fulfill_Timer.C:
						orderToBeFulfilled.DesignatedID = id
						break
					}
				}
			if orderToBeFulfilled.DesignatedID == id {
				ch.DoItMySelf <-orderToBeFulfilled
			}
*/
		case orderTakeOver := <- ch.DoItMySelf:
			newOrderChan <- orderTakeOver


		case outgoingOrder := <- ch.SendOrder:
			rand.Seed(time.Now().UTC().UnixNano())
			outgoingOrder.OrderID = rand.Intn(1000)			//Broadcast order to other elevators
			fmt.Println("Order with orderID:",outgoingOrder.OrderID,"has been sent")
			if outgoingOrder.DesignatedID != id { 	//Check if sendOrder.DesignatedID is not my id
				orderID = outgoingOrder.OrderID
				ch.OutgoingOrder <- outgoingOrder
				Ack_Timer.Reset(5*time.Second)
				fmt.Println("Ack_timer started for Elevator",outgoingOrder.DesignatedID) 								//Start acknowledgement timer
			}

		case orderAck := <- ch.ReceiveTimerMsg:                     //If the correct ack-message is received, then we stop the ack-timer and start the fulfill-timer
			fmt.Println("Received Timer Msg")
			if orderAck.Ack && orderAck.OrderID == orderID { //
				fmt.Println("Ack Timer stopped for orderID", orderAck.OrderID)
				Ack_Timer.Stop()
				Fulfill_Timer.Reset(5*time.Second)
			}
			if orderAck.OrderID == orderID && orderAck.Fulfilled {
				fmt.Println("Fulfillmenttimer stopped")
				Fulfill_Timer.Stop()
			}


		case <- Ack_Timer.C:									// The ack-message was not returned -> I have to do the order myself 
			fmt.Println("Order Ack Timeout Case")
		//	outGoingOrder.DesignatedID = id


		case <- Fulfill_Timer.C:								// The order was not fulfilled -> I have to do it myself 
			fmt.Println("Order Fulfilled Timeout Case")
		//	outGoingOrder.DesignatedID = id

			//default:
			//fmt.Println("In default, waiting for confirmation on ack/fulfill")

		/*	if outGoingOrder.DesignatedID == id {																// DesignitedID == my ID -> DoIiMySelf
				fmt.Println("I am best elevator, doing it my self")                              
				ch.OutgoingOrder <- outGoingOrder
				ch.DoItMySelf <- outGoingOrder

			}
*/
	/*	case newTimerMsg := <- ch.ReceiveTimerMsg:
			fmt.Println("I received a new timer message with ackstatus:",newTimerMsg.Ack,"and Fulfilled status:", newTimerMsg.Fulfilled,"and ID:",newTimerMsg.OrderID)
			if newTimerMsg.Ack {
				ch.OrderAck <- newTimerMsg
			} else if newTimerMsg.Fulfilled {
					ch.OrderFulfilled <- newTimerMsg
			}
*/


		case newOrder := <- ch.IncomingOrder: 	//New order is recieved from master
			if newOrder.DesignatedID == id {	//Check if we are suppose to take this order (Desginated order id is our id)
				fmt.Println("Elevator",id,"Recieved a new order at floor:",newOrder.Floor,"for button:",newOrder.Button)
				timerMsg := TimerMsg{newOrder.OrderID, true, false}
				go func() {ch.SendTimerMsg <- timerMsg}()
				fmt.Println("timerMsg has been send to ack order")
				newOrderChan <- newOrder
											//When order is finished --> tell other elevators to stop their fulfillmenttimers (might have to be done from some other place)
												//Tell governor to turn of lights
			} else {
				fmt.Println("I have received an order that is not mine, starting fulfillmenttimer")
				//orderId := newOrder.OrderID
				Fulfill_Timer.Reset(20*time.Second)
				/*for {
					select{
						case orderFulfilled := <- ch.OrderFulfilled: //From reciever
								if orderFulfilled.OrderID == orderId && orderFulfilled.Fulfilled {
									Fulfill_Timer.Stop()
									fmt.Println("Fullfilled-message was received")
									break
								}
						case <- Fulfill_Timer.C:
							fmt.Println("Fullfilled-message was not received")
							newOrder.DesignatedID = id
							ch.OutgoingOrder <- newOrder
							//fmt.Println("newOrder has been sent to outgoingorder")
							newOrderChan <- newOrder
							fmt.Println("newOrder has been sent to DoItMySelf")
							
							break
						}
						break
					}*/
					fmt.Println("exited for-loop")
			}




		}
	}
}
