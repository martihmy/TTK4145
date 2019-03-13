package main

import (. "./config"
	//"time"
	//"fmt"
	hw "./hardware_io"
	sm "./sm"
)

func main() {
	smChans := sm.SMChannels {
		FloorArrival:	make(chan int),
		Elevator:		make(chan Elevator),
		NewOrder:		make(chan ButtonEvent),
	}
	hw.Init("localhost:15657", NumFloors)


	go hw.PollFloorSensor(smChans.FloorArrival)
	initFloor := hw.InitElev(smChans.FloorArrival)
	go hw.PollButtons(smChans.NewOrder)
	go sm.ElevatorRun(smChans, initFloor)

	select {}
}
