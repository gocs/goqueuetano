package goqueuetano

import "time"

// Customer utilizes the queue to its fullest
type Customer struct {
	id       string
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
	entry    time.Time
}

// ID customer unique ID; this will make id readonly
func (c *Customer) ID() string {
	return c.id
}

// RemainingTime gives the remaining time from entry
func (c *Customer) RemainingTime() time.Duration {
	elapsed := time.Since(c.entry)
	return c.Duration - elapsed
}
