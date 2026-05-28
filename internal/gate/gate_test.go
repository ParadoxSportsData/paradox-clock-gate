package gate

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		tick    int
		maxTick uint16
		wantErr bool
	}{
		{"valid kickoff", 0, 3600, false},
		{"valid mid-game", 1800, 3600, false},
		{"valid at maxTick", 3600, 3600, false},
		{"negative tick", -1, 3600, true},
		{"tick exceeds maxTick by one", 3601, 3600, true},
		{"tick far exceeds maxTick", 999999, 3600, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.tick, tc.maxTick)
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate(%d, %d) error = %v, wantErr %v", tc.tick, tc.maxTick, err, tc.wantErr)
			}
		})
	}
}
