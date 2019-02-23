package system

type Status struct {
	ElevatorID       int
	DestinationFloor int
	CurrentFloor     int
	Direction        string
	Score            int
	PickupTotal      int
	DropoffTotal     int
}
