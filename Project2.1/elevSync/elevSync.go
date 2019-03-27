package elevSync

import (
	peers "../network/peers"
	"time"
	. "../config"
	//"math/rand"
	"fmt"
	"strconv"

)

type SyncChannels struct {
	OrderUpdate 				chan ButtonEvent
	UpdateOrderHandler			chan [NumElevators]Elevator
	PeerUpdateChan				chan peers.PeerUpdate
	TransmitEnable 				chan bool
	SyncUpdate 					chan Elevator
	UpdateGovOnlineList			chan [NumElevators]bool
	IncomingUpdateMsg 			chan Msg
	OutgoingUpdateMsg			chan Msg
}


func ElevatorSynchronizer(ch SyncChannels, id int, newOrderChan chan ButtonEvent) {

	var (
		elevatorList [NumElevators]Elevator
		localStatusMatrix [NumFloors][NumButtons-1]OrderInfo
		onlineElevators [NumElevators]bool
		wasTrulyOffline bool
		//change bool
	)


	timeer := make(chan bool)
	go func() {time.Sleep(1*time.Second); timeer<-true}()

	select{
	case init := <- ch.IncomingUpdateMsg:
		elevatorList = init.ElevatorList
		fmt.Println("Should have initilized our own queue from other elevs")
	case <- timeer:
		fmt.Println("ReinitTimer timed out")
	}

	UpdategovTicker := time.NewTicker(250*time.Millisecond)
	BcastTicker := time.NewTicker(200*time.Millisecond)
	for {
		select{


		case newOrder := <- ch.OrderUpdate:
			if newOrder.Finished {
				fmt.Println("Received order with finished")
				elevatorList[id].Queue[newOrder.Floor] = [NumButtons]bool{}
				if newOrder.Button != Btn_Cab {
					localStatusMatrix[newOrder.Floor][newOrder.Button].AckList[id] = 0
					localStatusMatrix[newOrder.Floor][newOrder.Button].DoneList[id] = 1

				}
			} else if newOrder.Button == Btn_Cab && newOrder.DesignatedID == id {
				elevatorList[id].Queue[newOrder.Floor][newOrder.Button] = true   //This update will be broadcasted on the standard msg each time the broadcaster ticker goes off
				fmt.Println("Received order for local cab call")
				ch.UpdateOrderHandler <- elevatorList
			} else {
				fmt.Println("Recieved outside order with DesignatedID:", newOrder.DesignatedID, "which has been acked")
				localStatusMatrix[newOrder.Floor][newOrder.Button].DesignatedID = newOrder.DesignatedID
				localStatusMatrix[newOrder.Floor][newOrder.Button].AckList[id] = 1
			}



		case newPeerUpdate := <- ch.PeerUpdateChan:
			if len(newPeerUpdate.Peers) <= 1  {
				fmt.Println("We went offline")
				onlineElevators[id] = false
				if len(newPeerUpdate.Peers) == 0 {
					wasTrulyOffline = true
				}
				}
			if len(newPeerUpdate.New) > 0 {
				if wasTrulyOffline == true{
					for floor:=0;floor<NumFloors;floor++{
						for btn:=Btn_Up;btn<Btn_Cab;btn++{
							for elev:=0;elev<NumElevators;elev++ {
								localStatusMatrix[floor][btn].AckList[elev] = 0
								localStatusMatrix[floor][btn].DoneList[elev] = 0
							}
						}
					}

				}
				idOfNewPeer,_ := strconv.Atoi(newPeerUpdate.New)
				if idOfNewPeer != id {
					onlineElevators[idOfNewPeer] = true
					onlineElevators[id] = true
					fmt.Println("new peer has been discovered with id:",idOfNewPeer)
				}
				}
			if len(newPeerUpdate.Lost) > 0 {
				idOfLostPeer,_ := strconv.Atoi(newPeerUpdate.Lost[0]) //Only one elevator will loose connection
				onlineElevators[idOfLostPeer] = false
				fmt.Println("lost peer has been discovered with id:",idOfLostPeer)
			}
			//how do we detect that we are offline?

			go func(){ch.UpdateGovOnlineList <- onlineElevators} ()


		case updateOnLocalElev := <- ch.SyncUpdate:
			if updateOnLocalElev.State != Undefined {
				tempQueue := elevatorList[id].Queue
				elevatorList[id] = updateOnLocalElev
				elevatorList[id].Queue = tempQueue
				ch.TransmitEnable <- true
				//Distribute change to other elevators
			} else {
				elevatorList[id] = updateOnLocalElev
				ch.TransmitEnable <- false
			}

		case <- BcastTicker.C:
		//	message := Msg{ElevatorList: elevatorList, StatusMatrix: localStatusMatrix, SenderID: id
			var message Msg
			message.ElevatorList = elevatorList
			message.StatusMatrix = localStatusMatrix
			message.SenderID = id
			ch.OutgoingUpdateMsg <- message

		case <- UpdategovTicker.C:
			ch.UpdateOrderHandler <- elevatorList



		case message := <- ch.IncomingUpdateMsg:
					if message.SenderID ==id{  //ignore messages that I broadcasted myself
							continue
					} else {
					/*	fmt.Println("Received these from senderID:",message.SenderID)
						fmt.Println("AckList:",message.StatusMatrix[1][Btn_Up].AckList[message.SenderID])
						fmt.Println("DoneList:",message.StatusMatrix[1][Btn_Up].DoneList[message.SenderID])*/
							elevatorList[message.SenderID] = message.ElevatorList[message.SenderID]
							for elev := 0; elev < NumElevators; elev++{
								if elev == id || !onlineElevators[elev] { //sjekker andre heiser ack og Done Lister
									continue
								}
								for btn := Btn_Up; btn < Btn_Cab; btn++ { //blar gjennom hele matrisen på deres indexer
									for floor := 0; floor < NumFloors; floor++{
										if floor == 3 && btn == Btn_Down {
											fmt.Println("Local stat:",localStatusMatrix[floor][btn].DesignatedID)
										}
										//	fmt.Println("AckList in local matrix for elevator:",message.SenderID,":",localStatusMatrix[floor][btn].AckList[message.SenderID])
										//	fmt.Println("DoneList in local matrix for elevator:",message.SenderID,":",localStatusMatrix[floor][btn].DoneList[message.SenderID])
										//	fmt.Println(message.StatusMatrix[floor][btn].AckList[id])
										//	fmt.Println(message.StatusMatrix[floor][btn].DoneList[id])

										if allOnSamePage(floor,btn,true,false,false,onlineElevators,message) && localStatusMatrix[floor][btn].DesignatedID==id && localStatusMatrix[floor][btn].DoneList[id] != 1{
											//	fmt.Println("Recieved order for floor",floor,"button:",btn,"That has been acked by all elevators online")
												localStatusMatrix[floor][btn].AckList[id] = 1
												elevatorList[id].Queue[floor][btn] = true
												ch.UpdateOrderHandler <- elevatorList																								//Pass på at denne går til governor + lys
																																				// Jeg kan utøre ordren hvis alle har acket og det er jeg som skal ta den
										}else if allOnSamePage(floor,btn,false,true,false, onlineElevators,message) && localStatusMatrix[floor][btn].AckList[id] != 1 && localStatusMatrix[floor][btn].AckList[message.SenderID] != 1 {
											localStatusMatrix[floor][btn].DoneList[id] = 0
											localStatusMatrix[floor][btn].DoneList[message.SenderID] = 0


										}else if message.StatusMatrix[floor][btn].DoneList[message.SenderID] == 1 && !allOnSamePage(floor,btn,false,true,false, onlineElevators,message) { //finner true i både ack og Done
										/*	fmt.Println(message.StatusMatrix[floor][btn].AckList[id])
											fmt.Println(message.StatusMatrix[floor][btn].DoneList[id])
											fmt.Println(message.StatusMatrix[floor][btn].AckList[message.SenderID])
											fmt.Println(message.StatusMatrix[floor][btn].DoneList[message.SenderID])

											fmt.Println("Found one 1 in a donelist and one in ack")*/

											localStatusMatrix[floor][btn].DoneList[id] = 1 //hvis vi finner vi en bestilling i vår oversikt som andre har meldt som "done", så melder vi den også som "done"
											localStatusMatrix[floor][btn].AckList[id] = 0
											localStatusMatrix[floor][btn].DoneList[message.SenderID] = 1
											localStatusMatrix[floor][btn].AckList[message.SenderID] = 0
									/*		fmt.Println("Change should have happend")
											fmt.Println(localStatusMatrix[floor][btn].AckList[id])
											fmt.Println(localStatusMatrix[floor][btn].DoneList[id])
											fmt.Println(localStatusMatrix[floor][btn].AckList[elev])
											fmt.Println(localStatusMatrix[floor][btn].DoneList[elev])*/

										}else if allOnSamePage(floor,btn,false,false,true, onlineElevators,message) && message.StatusMatrix[floor][btn].AckList[elev] == 1 && localStatusMatrix[floor][btn].DoneList[id] == 0 { ///* && !AllAckOnLocal(floor, btn, onlineElevators, localStatusMatrix*/
											localStatusMatrix[floor][btn].AckList[id] = 1
											localStatusMatrix[floor][btn].AckList[message.SenderID] = 1
											if message.ElevatorList[message.StatusMatrix[floor][btn].DesignatedID].State != Undefined {
												localStatusMatrix[floor][btn].DesignatedID = message.StatusMatrix[floor][btn].DesignatedID
											}
											fmt.Println("New order recieved for floor",floor,"button:",btn,"and has beend acked by",id)

										}else if allOnSamePage(floor,btn,false,true,false, onlineElevators, message) && message.StatusMatrix[floor][btn].DoneList[message.SenderID] == 1 && localStatusMatrix[floor][btn].AckList[elev] == 1 {
											localStatusMatrix[floor][btn].AckList[message.SenderID] = 0
											localStatusMatrix[floor][btn].DoneList[message.SenderID] = 1
											localStatusMatrix[floor][btn].AckList[id] = 0
											localStatusMatrix[floor][btn].DoneList[id] = 1
										}
									}
								}
							}
						}
		}
	}
}
