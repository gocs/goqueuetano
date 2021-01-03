package goqueuetano_test

import (
	"testing"

	"github.com/gocs/goqueuetano"
	"github.com/google/uuid"
)

// TestCustomerID tests for the existence and validity of uuid in ID
func TestCustomerID(t *testing.T) {
	cs := goqueuetano.Customers{}
	cs.Add(goqueuetano.Customer{})

	c, err := cs.GetByKey(0)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	if c.ID() == "" {
		t.Errorf("ID should not be empty: %s", c.ID())
	}

	_, err = uuid.Parse(c.ID())
	if err != nil {
		t.Errorf("ID is not valid uuid: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	c := goqueuetano.Customer{Name: "test", Total: 3}

	if c.Done() {
		t.Errorf("Update should not be done: %v", c)
	}
	if c.Done() {
		t.Errorf("Update should not be done: %v", c)
	}
	if c.Done() {
		t.Errorf("Update should not be done: %v", c)
	}
	if !c.Done() {
		t.Errorf("Update should be done: %v", c)
	}
}
