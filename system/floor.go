package system

type Floor struct {
	ID int
}

func NewFloor(i int) *Floor {
	return &Floor{ID: i}
}
