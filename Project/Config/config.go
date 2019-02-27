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

type State int

const (
	Undefined State  = iota -1
	Idle
	Moving
	DoorOpen
)

type Elevator struct {
	ElevState State
	Dir MotorDirection
	Floor int
	Queue[NumFloors][NumButtons]bool
}

