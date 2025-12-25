package strategy

import "testing"

func TestNewStrategy_Unknown(t *testing.T) {
	_, err := NewStrategy("unknown")
	if err == nil {
		t.Fatalf("expected error for unknown strategy")
	}
}
