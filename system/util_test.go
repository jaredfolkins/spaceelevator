package system

import "testing"

func TestNearestCarScoreUsesPickupFloorDistance(t *testing.T) {
	status := &Status{CurrentFloor: 10, Direction: Idle}
	passenger := &Passenger{CurrentFloor: 1, DestinationFloor: 15, Direction: Up}

	got := NearestCarScore(16, status, passenger)
	want := 8 // (15 + 2) - distance from car floor 10 to pickup floor 1.

	if got != want {
		t.Fatalf("NearestCarScore = %d, want %d", got, want)
	}
}

func TestNearestCarScoreMovingAwayIsOne(t *testing.T) {
	status := &Status{CurrentFloor: 10, Direction: Up}
	passenger := &Passenger{CurrentFloor: 4, DestinationFloor: 1, Direction: Down}

	got := NearestCarScore(16, status, passenger)
	if got != 1 {
		t.Fatalf("NearestCarScore = %d, want 1 for a car moving away from the call", got)
	}
}

func TestNearestCarScoreDirectionPreference(t *testing.T) {
	status := &Status{CurrentFloor: 10, Direction: Down}
	sameDirection := &Passenger{CurrentFloor: 4, DestinationFloor: 1, Direction: Down}
	oppositeDirection := &Passenger{CurrentFloor: 4, DestinationFloor: 12, Direction: Up}

	sameScore := NearestCarScore(16, status, sameDirection)
	oppositeScore := NearestCarScore(16, status, oppositeDirection)

	if sameScore <= oppositeScore {
		t.Fatalf("same-direction score = %d, opposite-direction score = %d; want same direction higher", sameScore, oppositeScore)
	}
}
