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
					localStatusMatrix[newOrder.Floor][newOrder.Button].StatusList[id].Acked = false
					localStatusMatrix[newOrder.Floor][newOrder.Button].StatusList[id].Done = true
				}
			} else if newOrder.Button == Btn_Cab && newOrder.DesignatedID == id {
				elevatorList[id].Queue[newOrder.Floor][newOrder.Button] = true   //This update will be broadcasted on the standard msg each time the broadcaster ticker goes off
				fmt.Println("Received order for local cab call")
				ch.UpdateOrderHandler <- elevatorList
			} else {
				fmt.Println("Recieved outside order with DesignatedID:", newOrder.DesignatedID, "which has been acked")
				localStatusMatrix[newOrder.Floor][newOrder.Button].DesignatedID = newOrder.DesignatedID
				localStatusMatrix[newOrder.Floor][newOrder.Button].StatusList[id].Acked = true
			}



		case newPeerUpdate := <- ch.PeerUpdateChan:
			if len(newPeerUpdate.New) > 0 {
				idOfNewPeer,_ := strconv.Atoi(newPeerUpdate.New)
				onlineElevators[idOfNewPeer] = true
				fmt.Println("new peer has been discovered with id:",idOfNewPeer)
			}
			if len(newPeerUpdate.Lost) > 0 {
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
			message := Msg{ElevatorList: elevatorList, StatusMatrix: localStatusMatrix, SenderID: id}
			ch.OutgoingUpdateMsg <- message
			fmt.Println("New broadcast message has been sent")


		case message := <- ch.IncomingUpdateMsg:
			fmt.Println("New message received")
			if message.SenderID ==id{  //ignore messages that I broadcasted myself
				continue
			} else {
					elevatorList[message.SenderID] = message.ElevatorList[message.SenderID]
					for elev := 0; elev < NumElevators; elev++{
						if elev == id { //sjekker andre heiser ack og Done Lister
							continue
						}
						for btn := Btn_Up; btn < Btn_Cab; btn++ { //blar gjennom hele matrisen på deres indexer
							for floor := 0; floor < NumFloors; floor++{
								fmt.Println("Has arrived inside received msg and inside all three loops")
								if allOnSamePage(floor,btn,true,false,false,onlineElevators,message) && message.StatusMatrix[floor][btn].DesignatedID==id{
										fmt.Println("Recieved order for floor",floor,"button:",btn,"That has been acked by all elevators online")
										elevatorList[id].Queue[floor][btn] = true
										ch.UpdateOrderHandler <- elevatorList																								//Pass på at denne går til governor + lys
									}																									// Jeg kan utøre ordren hvis alle har acket og det er jeg som skal ta den
																																		// send til statemaskin og oppdater andre

								if allOnSamePage(floor,btn,false,true,false, onlineElevators,message){
									localStatusMatrix[floor][btn].StatusList[id].Done=false
								}


								if !allOnSamePage(floor,btn,false,true,false, onlineElevators,message) && message.StatusMatrix[floor][btn].StatusList[elev].Done{ //finner true i både ack og Done
									localStatusMatrix[floor][btn].StatusList[id].Done = true //hvis vi finner vi en bestilling i vår oversikt som andre har meldt som "done", så melder vi den også som "done"
									localStatusMatrix[floor][btn].StatusList[id].Acked = false
								}
								if allOnSamePage(floor,btn,false,false,true, onlineElevators,message) && message.StatusMatrix[floor][btn].StatusList[elev].Acked {
									localStatusMatrix[floor][btn].StatusList[id].Acked = true
									fmt.Println("New order recieved for floor",floor,"button:",btn,"and has beend acked by",id)
								}
							}
						}
					}
				}
		}
	}
}
