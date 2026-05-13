package system

import (
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	passengerRandMu sync.Mutex
	passengerRand   = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
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
	if MaxFloor <= 0 {
		return
	}

	p.CurrentFloor = randomFloor(MaxFloor)
	if MaxFloor == 1 {
		p.DestinationFloor = p.CurrentFloor
		return
	}

	p.DestinationFloor = randomFloor(MaxFloor - 1)
	if p.DestinationFloor >= p.CurrentFloor {
		p.DestinationFloor++
	}
}

func (p *Passenger) CalcDirection() {
	switch {
	case p.DestinationFloor > p.CurrentFloor:
		p.Direction = Up
	case p.DestinationFloor < p.CurrentFloor:
		p.Direction = Down
	default:
		p.Direction = Idle
	}
}

func randomFloor(max int) int {
	passengerRandMu.Lock()
	defer passengerRandMu.Unlock()
	return passengerRand.Intn(max)
}
