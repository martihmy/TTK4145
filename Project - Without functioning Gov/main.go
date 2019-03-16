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
		OrderComplete: 	make(chan ButtonEvent),
	}
	govChans := gov.GovernorChannels {
		ButtonPress:	make(chan ButtonEvent),
	}
	hw.Init("localhost:15657", NumFloors)


	go hw.PollFloorSensor(smChans.FloorArrival)
	initFloor := hw.InitElev(smChans.FloorArrival)
	go hw.PollButtons(govChans.ButtonPress)
	go gov.ElevGovernor(initFloor, smChans.Elevator, smChans.OrderComplete, govChans)
	go sm.ElevatorRun(smChans, initFloor)

	select {}
}
