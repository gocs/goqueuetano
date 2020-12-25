package goqueuetano

import (
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
type Customers []Customer

// NewCustomers will ensure customers will leave after the duration times out
func NewCustomers() *Customers {
	c := &Customers{}
	go func() {
		for {
			// clean all the customers' table when they are done
			for _, cust := range c.All() {
				if cust.RemainingTime() < 0 {
					c.Delete(cust.ID())
				}
			}
		}
	}()
	return c
}

// Add is a new entry to the queue; ID is non editable
func (c *Customers) Add(customer ...Customer) {
	for _, v := range customer {
		*c = append(*c, Customer{
			id:       uuid.New().String(),
			Name:     v.Name,
			Duration: v.Duration,
			entry:    time.Now(),
		})
	}
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
func (c *Customers) GetByKey(k int) (customer Customer, err error) {
	if c.Len() == 0 {
		return Customer{}, errors.New("list not made")
	}
	return (*c)[k], nil
}

// Get gives the customer by its id
func (c *Customers) Get(id string) (customer Customer) {
	for _, cust := range *c {
		if cust.id == id {
			customer = cust
			break
		}
	}
	return
}

// Edit fulfill the customers changes
func (c *Customers) Edit(customer Customer) error {
	// filter invalid values
	if customer.id == "" {
		return errors.New("ID must not be Empty")
	}
	_, err := uuid.Parse(customer.id)
	if err != nil {
		return err
	}

	// if there are no more items, recreate
	if len(*c) < 1 {
		return errors.New("customer list is empty")
	}

	index, err := getIndex(*c, customer.ID())
	if err != nil {
		return errors.Wrap(err, "edit can't get index")
	}
	// prevent id and entry to be modified
	(*c)[index] = Customer{
		id:       (*c)[index].ID(),
		Name:     customer.Name,
		Duration: customer.Duration,
		entry:    (*c)[index].entry,
	}
	return nil
}

// Delete will acknowledges the departure of the customer
func (c *Customers) Delete(id string) error {
	// filter invalid values
	if id == "" {
		return errors.New("ID must not be Empty")
	}
	_, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// if there are no more items, recreate
	if len(*c) < 1 {
		return errors.New("customer list is empty")
	}

	// index should never error
	index, err := getIndex(*c, id)
	if err != nil {
		return errors.Wrap(err, "delete can't get index")
	}

	// concat lists before and after the specific item
	frontC := (*c)[:index]
	backC := (*c)[index+1:]
	*c = append(frontC, backC...)
	return nil
}

// getIndex
func getIndex(c Customers, id string) (int, error) {
	for k, v := range c {
		if v.ID() == id {
			return k, nil
		}
	}
	return -1, errors.New("item not found")
}
