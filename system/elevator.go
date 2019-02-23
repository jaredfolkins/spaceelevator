package system

import (
	"fmt"
	"sync"
)

// https://www.quora.com/Is-there-any-public-elevator-scheduling-algorithm-standard
type Elevator struct {
	mu               sync.RWMutex
	ID               int
	CurrentFloor     int
	DestinationFloor int
	Direction        string
	PickupQueue      map[string]*Pickup
	Target           *Passenger
	DropoffQueue     map[string]*Dropoff
	Pickup           chan *Pickup
	Dropoff          chan *Dropoff
	Status           chan chan *Status
	Tick             chan interface{}
}

func NewElevator(id int) *Elevator {
	e := &Elevator{
		CurrentFloor: 0,
		ID:           id,
		Direction:    Idle,
		PickupQueue:  make(map[string]*Pickup),
		DropoffQueue: make(map[string]*Dropoff),
		Pickup:       make(chan *Pickup),
		Dropoff:      make(chan *Dropoff),
		Status:       make(chan chan *Status),
		Tick:         make(chan interface{}),
	}
	e.Listener()
	return e
}

func (e *Elevator) Listener() {
	go func() {
		for {
			select {
			case <-e.Tick:
				e.tick()
			case ch := <-e.Status:
				ch <- e.MyStatus()
			case p := <-e.Pickup:
				e.pushPickupQueue(p)
			case d := <-e.Dropoff:
				e.pushDropoffQueue(d)
			}
		}
	}()
}

func (e *Elevator) pushPickupQueue(p *Pickup) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.PickupQueue[p.GetStringUUID()] = p
}

func (e *Elevator) delPickupQueue(uuid string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.PickupQueue[uuid]; ok {
		delete(e.PickupQueue, uuid)
	}
}

func (e *Elevator) pushDropoffQueue(d *Dropoff) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.DropoffQueue[d.GetStringUUID()] = d
}

func (e *Elevator) delDropoffQueue(uuid string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.DropoffQueue[uuid]; ok {
		delete(e.DropoffQueue, uuid)
	}

}

func (e *Elevator) MyStatus() *Status {
	return &Status{
		ElevatorID:       e.ID,
		DestinationFloor: e.DestinationFloor,
		CurrentFloor:     e.CurrentFloor,
		Direction:        e.Direction,
		PickupTotal:      len(e.PickupQueue),
		DropoffTotal:     len(e.DropoffQueue),
	}
}

// dropoff queue takes precident
func (e *Elevator) pickDestinationFloor() {
	if len(e.DropoffQueue) > 0 && e.Target == nil {
		for _, d := range e.DropoffQueue {
			e.Target = d.Passenger
			e.DestinationFloor = d.Passenger.DestinationFloor
		}
	} else if len(e.PickupQueue) > 0 && e.Target == nil {
		for _, p := range e.PickupQueue {
			e.Target = p.Passenger
			e.DestinationFloor = p.Passenger.CurrentFloor
		}
	}
}

func (e *Elevator) orient() {
	if e.DestinationFloor > e.CurrentFloor {
		e.Direction = Up
	} else if e.DestinationFloor < e.CurrentFloor {
		e.Direction = Down
	} else if e.DestinationFloor == e.CurrentFloor {
		e.Direction = Idle
	}
}

// unsafe
func (e *Elevator) move() {
	switch e.Direction {
	case Up:
		if e.CurrentFloor < e.DestinationFloor {
			e.CurrentFloor = e.CurrentFloor + 1
		} else if e.CurrentFloor == e.DestinationFloor {
			e.Direction = Idle
		}
	case Down:
		if e.CurrentFloor > e.DestinationFloor {
			e.CurrentFloor = e.CurrentFloor - 1
		} else if e.CurrentFloor == e.DestinationFloor {
			e.Direction = Idle
		}
	case Idle:
		e.Target = nil
	}
}

// optimization, it will check to see if
// this checks to see if (a) the passenger is in the pickup queue,
// if it is, allow the passenger to board
// then pop the request off the queue
func (e *Elevator) tryPickup() {
	for _, p := range e.PickupQueue {
		if e.CurrentFloor == p.Passenger.CurrentFloor {
			d := NewDropoff(p)
			e.pushDropoffQueue(d)
			e.delPickupQueue(p.GetStringUUID())
		}
	}
}

func (e *Elevator) tryDropoff() {
	for _, d := range e.DropoffQueue {
		if e.CurrentFloor == d.Passenger.DestinationFloor {
			e.delDropoffQueue(d.GetStringUUID())
		}
	}
}

// not thread safe
func (e *Elevator) tick() {
	e.pickDestinationFloor()
	e.orient()
	e.move()
	e.tryDropoff()
	e.tryPickup()
	fmt.Printf("Step() ID: %d, Target: %v, Current: %d, Destination: %d, Direction: %s, PickupQueue: %d, DropoffQueue %d\n", e.ID, e.Target, e.CurrentFloor, e.DestinationFloor, e.Direction, len(e.PickupQueue), len(e.DropoffQueue))
}
