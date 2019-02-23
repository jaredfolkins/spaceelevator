package system

func NearestCarScore(totalFloors int, s *Status, p *Passenger) int {
	n := totalFloors - 1
	if TowardsAndSame(s, p) {
		return (n + 2) - Distance(s.CurrentFloor, p.DestinationFloor)
	} else if TowardsAndDiff(s, p) {
		return (n + 1) - Distance(s.CurrentFloor, p.DestinationFloor)
	}
	d := Distance(s.CurrentFloor, p.DestinationFloor)
	total := n - d
	return total

}

func TowardsAndSame(s *Status, p *Passenger) bool {
	if p.CurrentFloor >= s.CurrentFloor {
		if s.Direction == Up || s.Direction == Idle {
			return true
		}
	}
	if p.CurrentFloor <= s.CurrentFloor {
		if s.Direction == Down || s.Direction == Idle {
			return true
		}
	}
	return false
}

func TowardsAndDiff(s *Status, p *Passenger) bool {
	if p.CurrentFloor > s.CurrentFloor && s.Direction != Up {
		return true
	}
	if p.CurrentFloor < s.CurrentFloor && s.Direction != Down {
		return true
	}
	return false
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
