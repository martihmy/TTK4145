package elevSync

import (
	"fmt"
	"time"
	. "./config"
)

type SyncChannels struct {
	UpdateSync	chan Elevator
	UpdateOrder	chan ButtonEvent
}

func ElevatorSynchronizer(ch SyncChannels, id int) {
	var (
		elevatorList [NumElevators]Elevator
	)
	ComfirmationTimer := time.NewTimer(2*time.Seconds)
	ComfirmationTimer.Stop()
	FulfillmentTimer := time.NewTimer(5*time.Seconds)
	FulfillmentTimer.Stop()

	for {
		select{
		case newOrder := <- ch.UpdateOrder:
			 
		}

	}
}