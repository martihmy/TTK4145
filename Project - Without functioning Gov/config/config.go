package config

const (
	NumFloors    = 4
	NumElevators = 3
	NumButtons   = 3
)

type ButtonType int

const (
	BtnUp   ButtonType 	= 0
	BtnDown           	= 1
	BtnCab             	= 2
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
	//ID int //May have to be moved to some other struct 
}

