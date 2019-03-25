package config

const (
	NumFloors    = 4
	NumElevators = 3
	NumButtons   = 3
)

type ButtonType int

const (
	Btn_Up   ButtonType 	= 0
	Btn_Down         	 	= 1
	Btn_Cab             	= 2
)

type MotorDirection int

const (
	Dir_Up   MotorDirection = 1
	Dir_Down                = -1
	Dir_Stop                = 0
)


type ButtonEvent struct {
	Floor  int
	Button ButtonType
	DesignatedID int
	Finished bool
}


type ElevState int

const (
	Undefined ElevState  = iota -1
	Idle
	Moving
	DoorOpen
)

type Elevator struct {
	State ElevState
	Dir MotorDirection
	Floor int
	Queue[NumFloors][NumButtons]bool
	ID int
}





type OrderInfo struct{
	DesignatedID int
	DoneList [NumElevators] int
	AckList [NumElevators] int
}
type Msg struct {
	ElevatorList [NumElevators]Elevator
	StatusMatrix	[NumFloors][NumButtons-1]OrderInfo
	SenderID 	int
}
