package gate

import "testing"

func TestValidateNegativeTick(t *testing.T) {
	err := Validate(-1, 3600)
	if err == nil {
		t.Fatal("expected error for negative tick, got nil")
	}
}

func TestValidateZeroTick(t *testing.T) {
	// Tick 0 is valid — kickoff state.
	err := Validate(0, 3600)
	if err != nil {
		t.Fatalf("expected no error for tick 0, got: %v", err)
	}
}

func TestValidateTickAtMaxTick(t *testing.T) {
	err := Validate(3600, 3600)
	if err != nil {
		t.Fatalf("expected no error for tick == maxTick, got: %v", err)
	}
}

func TestValidateTickExceedsMaxTick(t *testing.T) {
	err := Validate(3601, 3600)
	if err == nil {
		t.Fatal("expected error for tick > maxTick, got nil")
	}
}

func TestValidateTickFarExceedsMaxTick(t *testing.T) {
	err := Validate(999999, 3600)
	if err == nil {
		t.Fatal("expected error for tick 999999, got nil")
	}
}
