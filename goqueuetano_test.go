package goqueuetano_test

import (
	"testing"
	"time"

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

// TestRemainingTime tests the remaining time since sleeping
func TestRemainingTime(t *testing.T) {
	expected := 5 * time.Second
	sleepTime := time.Second
	margin := 50 * time.Millisecond

	cs := goqueuetano.Customers{}
	cs.Add(goqueuetano.Customer{Duration: 6 * time.Second})

	time.Sleep(sleepTime)

	c, err := cs.GetByKey(0)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}

	actual := c.RemainingTime()
	mexpected := expected - margin
	if expected < actual {
		t.Errorf("expected: %s < actual: %s", expected, actual)
	}
	if mexpected > actual {
		t.Errorf("margin+expected: %s > actual: %s", mexpected, actual)
	}

}
