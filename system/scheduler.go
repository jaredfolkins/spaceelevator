package system

import (
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	mu        sync.Mutex
	Floors    []*Floor
	Elevators []*Elevator
	Orders    []*Order
	Add       chan interface{}
	Pickup    chan *Pickup
	Dropoff   chan *Dropoff
	Cancel    chan *Order
	Blastoff  chan int
	Paint     chan chan *Paint
	Status    chan chan []*Status
}

func NewScheduler(floors, elevators int) *Scheduler {

	s := &Scheduler{
		Add:      make(chan interface{}),
		Pickup:   make(chan *Pickup),
		Dropoff:  make(chan *Dropoff),
		Cancel:   make(chan *Order),
		Blastoff: make(chan int),
		Paint:    make(chan chan *Paint),
		Status:   make(chan chan []*Status),
	}

	for i := 1; i <= floors; i++ {
		s.AddFloor(i)
	}

	for i := 1; i <= elevators; i++ {
		s.AddElevator(i)
	}

	return s
}

func (s *Scheduler) blastOff(ei int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.Elevators {
		for i := 0; i < ei; i++ {
			p := s.new()
			e.Pickup <- p
		}
	}
}

func (s *Scheduler) AddFloor(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Floors = append(s.Floors, NewFloor(i-1)) // off by one, frick, watch out
}

func (s *Scheduler) AddElevator(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Elevators = append(s.Elevators, NewElevator(i-1)) // off by one, frick, watch out
}

// TODO, optimize later with fan out, or something
func (s *Scheduler) NearestElevator(p *Pickup) *Elevator {
	s.mu.Lock()
	defer s.mu.Unlock()
	target := &Status{}
	tf := len(s.Floors)
	for _, e := range s.Elevators {
		ch := make(chan *Status)
		e.Status <- ch
		challenger := <-ch
		close(ch)
		challenger.Score = NearestCarScore(tf, challenger, p.Passenger)
		target = s.compare(target, challenger)
	}
	return s.Elevators[target.ElevatorID]
}

func (s *Scheduler) allStatus() []*Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	var st []*Status
	for _, e := range s.Elevators {
		ch := make(chan *Status)
		e.Status <- ch
		status := <-ch
		st = append(st, status)
		close(ch)
	}
	return st
}

// <= is critical, compare MUST always result in at least
// one real elevator ID to send an Order
func (s *Scheduler) compare(incumbant *Status, challenger *Status) *Status {
	if incumbant.Score <= challenger.Score {
		return challenger
	}
	return incumbant
}

func (s *Scheduler) new() *Pickup {
	max := len(s.Floors)
	pas := NewPassenger(max)
	pik := NewPickup(pas)
	return pik
}

func (s *Scheduler) pickup(p *Pickup) {
	e := s.NearestElevator(p)
	e.Pickup <- p
}

func (s *Scheduler) Render() *Paint {
	s.mu.Lock()
	defer s.mu.Unlock()

	cpy := make([]*Elevator, len(s.Elevators))
	copy(cpy, s.Elevators)

	pt := NewPaint()
	pt.CalcAliens(cpy)
	return pt
}
func (s *Scheduler) cancel(o *Order) {}

//func (s *Scheduler) completed(o *Order) {}

func (s *Scheduler) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.Elevators {
		e.Tick <- nil
	}
}

func (s *Scheduler) Run() {
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("------------------===TICK===----------------------")
				s.Tick()
			case <-s.Add:
				p := s.new()
				s.pickup(p)
			case i := <-s.Blastoff:
				s.blastOff(i)
			case order := <-s.Pickup:
				s.pickup(order)
			case order := <-s.Cancel:
				s.cancel(order)
			case drop := <-s.Dropoff:
				log.Println(drop)
			case ch := <-s.Status:
				st := s.allStatus()
				ch <- st
			case ch := <-s.Paint:
				ch <- s.Render()
			}
		}
	}()
}
