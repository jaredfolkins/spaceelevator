package system

import (
	"sync"

	"github.com/google/uuid"
)

type Order interface {
	GetUUID() uuid.UUID
	GetPassenger() *Passenger
}

type Pickup struct {
	mu        sync.Mutex
	UUID      uuid.UUID
	Passenger *Passenger
}

func NewPickup(p *Passenger) *Pickup {
	return &Pickup{
		UUID:      uuid.New(),
		Passenger: p,
	}
}

func (p *Pickup) GetUUID() uuid.UUID {
	return p.UUID
}

func (p *Pickup) GetStringUUID() string {
	return p.UUID.String()
}

func (p *Pickup) GetPassenger() *Passenger {
	return p.Passenger
}

type Dropoff struct {
	mu        sync.Mutex
	UUID      uuid.UUID
	Passenger *Passenger
}

func NewDropoff(p *Pickup) *Dropoff {
	return &Dropoff{
		UUID:      p.UUID,
		Passenger: p.Passenger,
	}
}

func (d *Dropoff) GetUUID() uuid.UUID {
	return d.UUID
}

func (d *Dropoff) GetStringUUID() string {
	return d.UUID.String()
}

func (d *Dropoff) GetPassenger() *Passenger {
	return d.Passenger
}
