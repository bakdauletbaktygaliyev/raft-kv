package counter

import "testing"

func TestInc(t *testing.T) {
	var c Counter
	c.Inc()
	c.Inc()
	c.Inc()
	if got := c.Value(); got != 3 {
		t.Errorf("Value() = %d, want 3", got)
	}
}

func TestIncValue(t *testing.T) {
	c := Counter{count: 41}
	c.IncValue() // value receiver: increments a copy
	if got := c.Value(); got != 41 {
		t.Errorf("after IncValue(), Value() = %d, want 41 (unchanged)", got)
	}
}
