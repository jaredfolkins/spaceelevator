package system

import "testing"

func TestCalcDirection(t *testing.T) {
	tests := []struct {
		name string
		p    Passenger
		want string
	}{
		{
			name: "up",
			p:    Passenger{CurrentFloor: 1, DestinationFloor: 9},
			want: Up,
		},
		{
			name: "down",
			p:    Passenger{CurrentFloor: 9, DestinationFloor: 1},
			want: Down,
		},
		{
			name: "idle",
			p:    Passenger{CurrentFloor: 4, DestinationFloor: 4},
			want: Idle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.CalcDirection()
			if tt.p.Direction != tt.want {
				t.Fatalf("Direction = %q, want %q", tt.p.Direction, tt.want)
			}
		})
	}
}

func TestNewPassengerGeneratesRealTrips(t *testing.T) {
	for i := 0; i < 200; i++ {
		p := NewPassenger(Floors)
		if p.CurrentFloor < 0 || p.CurrentFloor >= Floors {
			t.Fatalf("CurrentFloor = %d, want 0 through %d", p.CurrentFloor, Floors-1)
		}
		if p.DestinationFloor < 0 || p.DestinationFloor >= Floors {
			t.Fatalf("DestinationFloor = %d, want 0 through %d", p.DestinationFloor, Floors-1)
		}
		if p.CurrentFloor == p.DestinationFloor {
			t.Fatalf("generated same-floor trip on floor %d", p.CurrentFloor)
		}
		if p.Direction == Idle {
			t.Fatalf("generated non-moving passenger with direction %q: %#v", p.Direction, p)
		}
	}
}

func TestNewPassengerSingleFloorIsIdle(t *testing.T) {
	p := NewPassenger(1)
	if p.CurrentFloor != 0 || p.DestinationFloor != 0 || p.Direction != Idle {
		t.Fatalf("single-floor passenger = %#v, want floor 0 idle", p)
	}
}
