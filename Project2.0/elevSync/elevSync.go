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
	UpdateOrderHandler	chan [NumElevators]Elevator
	PeerUpdateChan			chan peers.PeerUpdate
	TransmitEnable 			chan bool
	SyncUpdate 					chan Elevator
	UpdateGovOnlineList chan [NumElevators]bool
	IncomingUpdateMsg 	chan Msg
	OutgoingUpdateMsg		chan Msg
	Outgoingtest				chan int
	Incomingtest				chan int
}


func ElevatorSynchronizer(ch SyncChannels, id int, newOrderChan chan ButtonEvent) {

	var (
		elevatorList [NumElevators]Elevator
		localStatusMatrix [NumFloors][NumButtons-1]OrderInfo
		onlineElevators [NumElevators]bool
		//change bool
	)

	Ack_Timer := time.NewTimer(2*time.Second)
	Ack_Timer.Stop()
	Fulfill_Timer := time.NewTimer(5*time.Second)
	Fulfill_Timer.Stop()
	BcastTicker := time.NewTicker(2000*time.Millisecond)


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
			/*if len(newPeerUpdate.Peers) == 1 {
				onlineElevators[id] = false
			} else*/
				if len(newPeerUpdate.New) > 0 {
				idOfNewPeer,_ := strconv.Atoi(newPeerUpdate.New)
				onlineElevators[idOfNewPeer] = true
				fmt.Println("new peer has been discovered with id:",idOfNewPeer)
				}else if len(newPeerUpdate.Lost) > 0 {
					idOfLostPeer,_ := strconv.Atoi(newPeerUpdate.Lost[0]) //Only one elevator will loose connection
					onlineElevators[idOfLostPeer] = false
					fmt.Println("lost peer has been discovered with id:",idOfLostPeer)
				}

			//how do we detect that we are offline?

			ch.UpdateGovOnlineList <- onlineElevators


		case updateOnLocalElev := <- ch.SyncUpdate:
			if updateOnLocalElev.State != Undefined {
				elevatorList[id] = updateOnLocalElev
				ch.TransmitEnable <- true
				//Distribute change to other elevators
			}

		case <- BcastTicker.C:
		//	message := Msg{ElevatorList: elevatorList, StatusMatrix: localStatusMatrix, SenderID: id}
			var message Msg
			message.ElevatorList = elevatorList
			message.StatusMatrix = localStatusMatrix
			message.SenderID = id
			ch.OutgoingUpdateMsg <- message



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
										//	fmt.Println("AckList in local matrix for elevator:",message.SenderID,":",localStatusMatrix[floor][btn].AckList[message.SenderID])
										//	fmt.Println("DoneList in local matrix for elevator:",message.SenderID,":",localStatusMatrix[floor][btn].DoneList[message.SenderID])
										//	fmt.Println(message.StatusMatrix[floor][btn].AckList[id])
										//	fmt.Println(message.StatusMatrix[floor][btn].DoneList[id])

										if allOnSamePage(floor,btn,true,false,false,onlineElevators,message) && message.StatusMatrix[floor][btn].DesignatedID==id && localStatusMatrix[floor][btn].DoneList[id] != 1{
												fmt.Println("Recieved order for floor",floor,"button:",btn,"That has been acked by all elevators online")
												localStatusMatrix[floor][btn].AckList[id] = 1
												elevatorList[id].Queue[floor][btn] = true
												ch.UpdateOrderHandler <- elevatorList																								//Pass på at denne går til governor + lys
																																				// Jeg kan utøre ordren hvis alle har acket og det er jeg som skal ta den
										}else	if allOnSamePage(floor,btn,false,true,false, onlineElevators,message) && localStatusMatrix[floor][btn].AckList[id] != 1 && localStatusMatrix[floor][btn].AckList[message.SenderID] != 1 {
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

										}else	if allOnSamePage(floor,btn,false,false,true, onlineElevators,message) && message.StatusMatrix[floor][btn].AckList[elev] == 1 && localStatusMatrix[floor][btn].DoneList[id] == 0 && !AllAckOnLocal(floor, btn, onlineElevators, localStatusMatrix) {
											localStatusMatrix[floor][btn].AckList[id] = 1
											localStatusMatrix[floor][btn].AckList[message.SenderID] = 1
											localStatusMatrix[floor][btn].DesignatedID = message.StatusMatrix[floor][btn].DesignatedID
											fmt.Println("New order recieved for floor",floor,"button:",btn,"and has beend acked by",id)

										}else if allOnSamePage(floor,btn,false,true,false, onlineElevators, message) && message.StatusMatrix[floor][btn].DoneList[message.SenderID] == 1 && localStatusMatrix[floor][btn].AckList[elev] == 1 {
											localStatusMatrix[floor][btn].AckList[message.SenderID] = 0
											localStatusMatrix[floor][btn].DoneList[message.SenderID] = 1
											localStatusMatrix[floor][btn].AckList[id] = 0
											localStatusMatrix[floor][btn].DoneList[id] = 1
										}
										for i :=0; i<3;i++{
											fmt.Println("Ack: ", i)
											fmt.Println(localStatusMatrix[1][Btn_Up].AckList[i])
											fmt.Println("Done: ", i)
											fmt.Println(localStatusMatrix[1][Btn_Up].DoneList[i])
										}

									}
								}
							}
							fmt.Println("------------------------")
						}
		}
	}
}
