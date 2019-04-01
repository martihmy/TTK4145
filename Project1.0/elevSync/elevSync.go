package elevSync

import (
	peers "../network/peers"
	"time"
	. "../config"
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

	UpdategovTicker := time.NewTicker(250*time.Millisecond) //Triggers update-message to governor every 250 millisecond
	BcastTicker := time.NewTicker(200*time.Millisecond)     //Triggers update-message to other elevators every 200 millisecond

	for {

		select{

		case newOrder := <- ch.OrderUpdate:
			if newOrder.Finished {
				fmt.Println("Some order has been fulfilled")
				elevatorList[id].Queue[newOrder.Floor] = [NumButtons]bool{}
				if newOrder.Button != Btn_Cab {
					localStatusMatrix[newOrder.Floor][newOrder.Button].AckList[id] = 0
					localStatusMatrix[newOrder.Floor][newOrder.Button].DoneList[id] = 1

				}
			} else if newOrder.Button == Btn_Cab && newOrder.DesignatedID == id {
				elevatorList[id].Queue[newOrder.Floor][newOrder.Button] = true 
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
								localStatusMatrix[floor][btn].AckList[elev] = 0  //reset Acknowledge- and DoneList if the elevator lose connection
								localStatusMatrix[floor][btn].DoneList[elev] = 0 //...
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
				deadPeer,_ := strconv.Atoi(newPeerUpdate.Lost[0]) //Only one elevator will loose connection
				onlineElevators[deadPeer] = false

				fmt.Println("lost peer has been discovered with id:",deadPeer)
			}

			go func(){ch.UpdateGovOnlineList <- onlineElevators} ()


		case updateOnLocalElev := <- ch.SyncUpdate:
			if updateOnLocalElev.State != Undefined {
				tempQueue := elevatorList[id].Queue
				elevatorList[id] = updateOnLocalElev
				elevatorList[id].Queue = tempQueue
				ch.TransmitEnable <- true
				//Distribute change to other elevators

			} else {
				ch.TransmitEnable <- false
			}

		case <- BcastTicker.C:
		//	Format of broadcasted message:  "message := Msg{ElevatorList: elevatorList, StatusMatrix: localStatusMatrix, SenderID: id} "
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
							elevatorList[message.SenderID] = message.ElevatorList[message.SenderID]
							for elev := 0; elev < NumElevators; elev++{
								if elev == id || !onlineElevators[elev] { //We only care about acknowledgement from other elevators that are online
									continue
								}
								for btn := Btn_Up; btn < Btn_Cab; btn++ {
									for floor := 0; floor < NumFloors; floor++{


										if allOnSamePage(floor,btn,true,false,false,onlineElevators,message) && localStatusMatrix[floor][btn].DesignatedID==id && localStatusMatrix[floor][btn].DoneList[id] != 1{
											//Some order is acknowledged by all fully-functional elevators and I am the designated elevator. Send update to orderHandler
											localStatusMatrix[floor][btn].AckList[id] = 1
											elevatorList[id].Queue[floor][btn] = true
											ch.UpdateOrderHandler <- elevatorList
										
										}else if allOnSamePage(floor,btn,false,true,false, onlineElevators,message) && localStatusMatrix[floor][btn].AckList[id] != 1 && localStatusMatrix[floor][btn].AckList[message.SenderID] != 1 {
											//Reset DoneList to make sure that we are ready to ackowledge a new order on that floor for that specific button type
											localStatusMatrix[floor][btn].DoneList[id] = 0
											localStatusMatrix[floor][btn].DoneList[message.SenderID] = 0


										}else if message.StatusMatrix[floor][btn].DoneList[message.SenderID] == 1 && !allOnSamePage(floor,btn,false,true,false, onlineElevators,message) {
											//We ackowledge that some order has been fulfilled
											localStatusMatrix[floor][btn].DoneList[id] = 1
											localStatusMatrix[floor][btn].AckList[id] = 0
											localStatusMatrix[floor][btn].DoneList[message.SenderID] = 1
											localStatusMatrix[floor][btn].AckList[message.SenderID] = 0
							
										}else if allOnSamePage(floor,btn,false,false,true, onlineElevators,message) && message.StatusMatrix[floor][btn].AckList[elev] == 1 && localStatusMatrix[floor][btn].DoneList[id] == 0 && !AllAckOnLocal(floor, btn, onlineElevators, localStatusMatrix) {
											//Some other elevator has ackowledged an order that I have not. I will Ackowledge it myself 
											localStatusMatrix[floor][btn].AckList[id] = 1
											localStatusMatrix[floor][btn].AckList[message.SenderID] = 1
											localStatusMatrix[floor][btn].DesignatedID = message.StatusMatrix[floor][btn].DesignatedID

										}else if allOnSamePage(floor,btn,false,true,false, onlineElevators, message) && message.StatusMatrix[floor][btn].DoneList[message.SenderID] == 1 && localStatusMatrix[floor][btn].AckList[elev] == 1 {
											//Some order that I have acknowledged has been marked as Done. I will mark it as Done as well
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
