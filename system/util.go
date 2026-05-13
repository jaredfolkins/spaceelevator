package system

func NearestCarScore(totalFloors int, s *Status, p *Passenger) int {
	n := totalFloors - 1
	d := Distance(s.CurrentFloor, p.CurrentFloor)
	if TowardsAndSame(s, p) {
		return (n + 2) - d
	} else if TowardsAndDiff(s, p) {
		return (n + 1) - d
	}
	return 1
}

func TowardsAndSame(s *Status, p *Passenger) bool {
	if !towardsCall(s, p) {
		return false
	}
	if s.Direction == Idle {
		return true
	}
	return s.Direction == p.Direction
}

func TowardsAndDiff(s *Status, p *Passenger) bool {
	if !towardsCall(s, p) || s.Direction == Idle {
		return false
	}
	return s.Direction != p.Direction
}

func towardsCall(s *Status, p *Passenger) bool {
	switch s.Direction {
	case Up:
		return p.CurrentFloor >= s.CurrentFloor
	case Down:
		return p.CurrentFloor <= s.CurrentFloor
	default:
		return true
	}
}

func Distance(cur, dest int) int {
	var d int
	if cur > dest {
		d = cur - dest
	} else if dest > cur {
		d = dest - cur
	}
	if d > 0 {
		return d
	}
	return 0
}
