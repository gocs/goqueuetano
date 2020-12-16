package goqueuetano

import (
	"time"

	"github.com/google/uuid"
)

// Customer utilizes the queue to its fullest
type Customer struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
}

// Order defines how customers uses the queue
type Order interface {
	Add(customer Customer)
	All() []Customer
	Len() int
	GetByKey(k int) (customer Customer)
	Edit(customer Customer)
	Delete(id string)
}

// Customers is the collection of the users
type Customers []Customer

// Add is a new entry to the queue
func (c *Customers) Add(customer Customer) {
	customer.ID = uuid.New().String()
	*c = append(*c, customer)
}

// All of itself is returned out from the interface
func (c *Customers) All() []Customer {
	return *c
}

// Len is a interface container len; interface can't be len, `Customers` can
func (c *Customers) Len() int {
	return len(*c)
}

// GetByKey identifies the customer by its order
func (c *Customers) GetByKey(k int) (customer Customer) {
	return (*c)[k]
}

// Edit fulfill the customers changes
func (c Customers) Edit(customer Customer) {
	index := getIndex(c, customer.ID)
	c[index] = customer
}

// Delete will acknowledges the departure of the customer
func (c *Customers) Delete(id string) {
	index := getIndex(*c, id)
	// concat lists before and after the specific item
	frontC := (*c)[:index]
	backC := (*c)[index+1:]
	*c = append(frontC, backC...)
}

// getIndex
func getIndex(c Customers, id string) int {
	for k, v := range c {
		if v.ID == id {
			return k
		}
	}
	return -1
}
