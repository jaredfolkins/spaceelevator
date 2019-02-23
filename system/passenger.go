package system

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Passenger struct {
	UUID             uuid.UUID
	CurrentFloor     int
	DestinationFloor int
	Direction        string
	Active           bool
}

func NewPassenger(MaxFloor int) *Passenger {
	p := &Passenger{}
	p.GenCurrentAndDest(MaxFloor)
	p.CalcDirection()
	p.UUID = uuid.New()
	return p
}

func (p *Passenger) GenCurrentAndDest(MaxFloor int) {
	rand.Seed(time.Now().UTC().UnixNano())
	p.CurrentFloor = rand.Intn(MaxFloor)
	p.DestinationFloor = rand.Intn(MaxFloor)
}

func (p *Passenger) CalcDirection() {
	if p.DestinationFloor > p.CurrentFloor {
		p.Direction = Up
	} else if p.DestinationFloor < p.CurrentFloor {
		p.Direction = Down
	}
	p.Direction = Idle
}
