package main

import (. "./config"
	//"time"
	"fmt"
	hw "./hardware_io"
	sm "./sm"
	gov "./governor"
	"flag"
	"strconv"
	sync "./elevSync"
	bcast "./network/bcast"

)

func main() {
	var (
		id string
		ID int
	)

	flag.StringVar(&id, "id", "0", "id of this peer")
	flag.Parse()
	ID, _ = strconv.Atoi(id)
	fmt.Println(ID)

	smChans := sm.SMChannels {
		FloorArrival:		make(chan int),
		Elevator:				make(chan Elevator),
		NewOrder:				make(chan ButtonEvent),
		ServicedFloor:	make(chan int),
		OrderFulfilled: make(chan ButtonEvent),
	}

	syncChans := sync.SyncChannels {
	UpdateSynchronizer		make(chan Elevator),
	IncomingUpdateMsg		make(chan Elevator),
	OutgoingUpdateMsg 		make(chan Elevator),
	UpdateGov				make(chan Elevator),
	OutgoingOrder 			make(chan ButtonEvent),
	IncomingOrder			make(chan ButtonEvent),
	DoItMySelf				make(chan ButtonEvent),
	SendOrder 				make(chan ButtonEvent),
	OrderAck				make(chan TimerMsg),
	ReceiveTimerMsg 		make(chan TimerMsg),
	SendTimerMsg			make(chan TimerMsg),
	}


	hw.Init("localhost:15657", NumFloors)

	var (
		btnPressChan 		= make(chan ButtonEvent)
		lightUpdaterChan 	= make(chan [NumElevators]Elevator)
	)




	go hw.PollFloorSensor(smChans.FloorArrival)
	initFloor := hw.InitElev(smChans.FloorArrival)
	go hw.PollButtons(btnPressChan)
	go gov.ElevGovernor(ID, btnPressChan, smChans.NewOrder, lightUpdaterChan, smChans.Elevator, smChans.ServicedFloor, syncChans.SendOrder, syncChans.UpdateSynchronizer, syncChans.UpdateGov)
	go gov.LightUpdater(lightUpdaterChan,ID)
	go sm.ElevatorRun(smChans, initFloor,ID)
	go sync.ElevatorSynchronizer(syncChans, ID, smChans.NewOrder)

	//Must handle if an elevator goes down and reinitialized with a zero queue (so it copies its queue from someone else)


	go bcast.Transmitter(20034, syncChans.OutgoingOrder, syncChans.OutgoingUpdateMsg, syncChans.SendTimerMsg)
	go bcast.Receiver(20034, syncChans.IncomingOrder, syncChans.IncomingUpdateMsg, syncChans.ReceiveTimerMsg)
	select {}
}
