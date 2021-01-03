package goqueuetano

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Order defines how customers uses the queue
type Order interface {
	Add(customer ...Customer)
	All() []Customer
	Len() int
	Get(id string) (customer Customer)
	GetByKey(k int) (customer Customer, err error)
	Edit(customer Customer) error
	Delete(id string) error
}

// Customers is the collection of the users
type Customers struct {
	list []Customer
	mu   sync.Mutex
}

// Add is a new entry to the queue; ID is non editable
func (c *Customers) Add(customer ...Customer) {
	for _, v := range customer {
		cust := Customer{
			id:    uuid.New().String(),
			Name:  v.Name,
			Total: v.Total,
			entry: time.Now(),
		}

		c.mu.Lock()
		c.list = append(c.list, cust)
		c.mu.Unlock()
	}
}

// All of the available customers is returned out from the interface
func (c *Customers) All() []Customer {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.list
}

// Len returns the number of available customers
func (c *Customers) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.list)
}

// GetByKey identifies the customer by its order
func (c *Customers) GetByKey(k int) (customer Customer, err error) {
	if c.Len() == 0 {
		return Customer{}, errors.New("list not made")
	}
	if k > c.Len()-1 {
		return Customer{}, errors.Errorf("key out of range from the list [%v]", k)
	}
	return c.list[k], nil
}

// Get gives the customer by its id
func (c *Customers) Get(id string) (customer Customer) {
	for _, cust := range c.list {
		if cust.id == id {
			customer = cust
			break
		}
	}
	return
}

// Edit fulfill the customers changes
func (c *Customers) Edit(customer Customer) error {

	index, err := findKeyByIndex(c, customer.ID())
	if err != nil {
		return err
	}
	currCust := c.list[index]
	// prevent id and entry to be modified
	c.list[index] = Customer{
		id:      currCust.ID(),
		Name:    customer.Name,
		Total:   currCust.Total,
		Current: currCust.Current,
		entry:   currCust.entry,
	}
	return nil
}

// Delete will acknowledges the departure of the customer
func (c *Customers) Delete(id string) error {
	// index should never error
	index, err := findKeyByIndex(c, id)
	if err != nil {
		return err
	}

	// concat lists before and after the specific item
	frontC := c.list[:index]
	backC := c.list[index+1:]
	c.list = append(frontC, backC...)
	return nil
}

// findKeyByIndex
func findKeyByIndex(c *Customers, index string) (key int, err error) {
	key = -1
	err = errors.New("item not found")

	if _, err = uuid.Parse(index); err != nil {
		return
	}

	// if there are no more items, recreate
	if c.Len() == 0 {
		err = errors.New("customer list is empty")
		return
	}

	c.mu.Lock()
	clist := c.list
	c.mu.Unlock()

	for k, v := range clist {
		if v.ID() == index {
			key = k
			err = nil
		}
	}

	return
}
