package main

import (. "./config"
	//"time"
	//"fmt"
	hw "./hardware_io"
	sm "./sm"
	gov "./governor"
)

func main() {
	smChans := sm.SMChannels {
		FloorArrival:	make(chan int),
		Elevator:		make(chan Elevator),
		NewOrder:		make(chan ButtonEvent),
		ServicedFloor:	make(chan int),
	}
	hw.Init("localhost:15657", NumFloors)

	var (
		btnPressChan 		= make(chan ButtonEvent)
		lightUpdaterChan 	= make(chan [NumElevators]Elevator)
	)


	go hw.PollFloorSensor(smChans.FloorArrival)
	initFloor := hw.InitElev(smChans.FloorArrival)
	go hw.PollButtons(btnPressChan)
	go gov.ElevGovernor(0, btnPressChan, smChans.NewOrder, lightUpdaterChan, smChans.Elevator, smChans.ServicedFloor)
	go gov.LightUpdater(lightUpdaterChan,0)
	go sm.ElevatorRun(smChans, initFloor)

	select {}
}
