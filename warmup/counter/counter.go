// Package counter demonstrates the difference between value and pointer receivers.
package counter

// Counter is a simple integer counter. The zero value is ready to use.
type Counter struct {
	count int
}

func (c Counter) IncValue() {
	c.count++ //nolint:staticcheck // deliberate: value receiver demo
}

// Inc adds 1 to the counter.
func (c *Counter) Inc() {
	c.count++
}

// Value returns the current count.
func (c *Counter) Value() int {
	return c.count
}
