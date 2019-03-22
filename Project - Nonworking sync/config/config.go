package config

const (
	NumFloors    = 4
	NumElevators = 1
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

type OrderStatus int

const (
	Unconfirmed OrderStatus = iota -1
	Empty
	Confirmed
	Finished
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
	DesignatedID int
	OrderID int
}

//Only a tryout:
type ButtonPress struct {
	Floor int
	Button ButtonType
	DesignatedID int
	OrderID int
	Status OrderStatus
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
	ID int //May have to be moved to some other struct
}


type OrderMsg struct {
	Elevs [NumElevators]Elevator
	ID int
	Order ButtonEvent
}

type TimerMsg struct {
	OrderID 					int
	Ack							bool
	Fulfilled					bool
}
