package system

type Paint struct {
	Graph [][]string
}

func NewPaint() *Paint {
	return &Paint{
		Graph: td(),
	}
}

func td() [][]string {
	x := Elevators
	y := Floors
	td := make([][]string, y)
	for i := 0; i < y; i++ {
		td[i] = make([]string, x)
		for j := 0; j < x; j++ {
			td[i][j] = EmptyCell
		}
	}
	return td
}

// Calculates the aliens/elevators positions
func (pt *Paint) CalcAliens(elevators []*Elevator) {
	padding := 0
	for _, e := range elevators {
		switch e.Direction {
		case Up:
			pt.Graph[e.CurrentFloor][e.ID+padding] = UpCell
		case Down:
			pt.Graph[e.CurrentFloor][e.ID+padding] = DownCell
		case Idle:
			pt.Graph[e.CurrentFloor][e.ID+padding] = IdleCell
		}
	}
}
