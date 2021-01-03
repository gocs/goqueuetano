package goqueuetano

import "time"

// Customer utilizes the queue to its fullest
type Customer struct {
	id      string
	Name    string `json:"name"`
	Total   int    `json:"total"`
	Current int    `json:"current"`
	entry   time.Time
}

// ID customer unique ID; this will make id readonly
func (c *Customer) ID() string {
	return c.id
}

// Done updates customer status
// Done updates customer status
func (c *Customer) Done() bool {
	if c.Current >= c.Total {
		return true
	}
	c.Current++
	return false
}
