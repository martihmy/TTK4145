package main

import (. "./config"
//	"time"
	"fmt"
	hw "./hardware_io"
	sm "./sm"
	gov "./orderHandler"
	"flag"
	"strconv"
	sync "./elevSync"
	bcast "./network/bcast"
	peers "./network/peers"

)

func main() {
	var (
		id string
		ID int
		simPort string
	)

	flag.StringVar(&id, "id", "0", "id of this peer")
	flag.StringVar(&simPort, "simPort", "44523", "simulation server port")
	flag.Parse()
	ID, _ = strconv.Atoi(id)
	fmt.Println(ID)

	smChans := sm.SMChannels {
		FloorArrival:		make(chan int),
		Elevator:				make(chan Elevator),
		NewOrder:				make(chan ButtonEvent),
		ServicedFloor:	make(chan int),
	}

	syncChans := sync.SyncChannels {
		OrderUpdate: make(chan ButtonEvent),
		UpdateOrderHandler:	make(chan [NumElevators]Elevator),
		TransmitEnable: make(chan bool),
		PeerUpdateChan: make(chan peers.PeerUpdate),
		SyncUpdate: make(chan Elevator),
		UpdateGovOnlineList: make(chan [NumElevators]bool),
		IncomingUpdateMsg: make(chan Msg),
		OutgoingUpdateMsg: make(chan Msg),

	}


	hw.Init("Localhost:"+simPort, NumFloors)

	var (
		btnPressChan 		= make(chan ButtonEvent)
		lightUpdaterChan 	= make(chan [NumElevators]Elevator)
	)




	go hw.PollFloorSensor(smChans.FloorArrival)
	initFloor := hw.InitElev(smChans.FloorArrival)
	go hw.PollButtons(btnPressChan)
	go gov.OrderHandler(ID, btnPressChan, smChans.NewOrder, lightUpdaterChan, smChans.Elevator, smChans.ServicedFloor, syncChans.UpdateOrderHandler, syncChans.OrderUpdate, syncChans.SyncUpdate, syncChans.UpdateGovOnlineList)
	go gov.LightUpdater(lightUpdaterChan,ID)
	go sm.ElevatorRun(smChans, initFloor,ID)
	go sync.ElevatorSynchronizer(syncChans, ID, smChans.NewOrder)

	//Must handle if an elevator goes down and reinitialized with a zero queue (so it copies its queue from someone else)
	go peers.Transmitter(20344,id,syncChans.TransmitEnable)
	go peers.Receiver(20344, syncChans.PeerUpdateChan)

	go bcast.Transmitter(43034, syncChans.OutgoingUpdateMsg)
	go bcast.Receiver(43034, syncChans.IncomingUpdateMsg)
	select {}
}
