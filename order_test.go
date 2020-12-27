package goqueuetano_test

import (
	"testing"
	"time"

	"github.com/gocs/goqueuetano"
)

type Expected struct {
	name     string
	duration time.Duration
}

// integrity_test checks whether the expected ang given values are the same.
//	if it returns false then the datas are mismatched
func integrity_test(expected []Expected, given *goqueuetano.Customers) (bool, error) {
	for i := 0; i < given.Len(); i++ {
		c, e := given.GetByKey(i)
		if e != nil {
			return false, e
		}
		if c.Name != expected[i].name ||
			c.Duration != expected[i].duration {
			return false, nil
		}
	}
	return true, nil
}

func TestOrderGetByKey(t *testing.T) {
	expected := []Expected{
		{name: "Apple", duration: 5 * time.Second},
		{name: "Ball", duration: 4 * time.Second},
		{name: "Cat", duration: 3 * time.Second},
	}
	cs := goqueuetano.Customers{}

	c, err := cs.GetByKey(0)
	if err == nil {
		t.Errorf("customers is expecting to have error")
	}
	if c != (goqueuetano.Customer{}) {
		t.Errorf("customers is expected to be empty")
	}
	cs.Add(goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second})

	c, err = cs.GetByKey(2)
	if err != nil {
		if err.Error() != "key out of range from the list [2]" {
			t.Errorf("unexpected behaviour: %v", err)
		}
	} else {
		if c == (goqueuetano.Customer{}) {
			t.Errorf("customers is expected to have a value: %v", c)
		}
	}

	c, err = cs.GetByKey(0)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	if c == (goqueuetano.Customer{}) {
		t.Errorf("customers is expected to have a value: %v", c)
	}

	ok, err := integrity_test(expected, &cs)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	if !ok {
		t.Errorf("data mismatch: %v, %v", expected, cs.All())
	}
}

func TestOrderAdd(t *testing.T) {
	expected := []Expected{
		{name: "Apple", duration: 5 * time.Second},
		{name: "Ball", duration: 4 * time.Second},
		{name: "Cat", duration: 3 * time.Second},
	}
	cs := goqueuetano.Customers{}

	cs.Add(goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second})
	if cs.Len() == 0 {
		t.Errorf("customers must not be empty")
	}

	cs.Add(
		goqueuetano.Customer{Name: "Ball", Duration: 4 * time.Second},
		goqueuetano.Customer{Name: "Cat", Duration: 3 * time.Second},
	)
	if cs.Len() != 3 {
		t.Errorf("customers is expected to have 3 elements")
	}

	ok, err := integrity_test(expected, &cs)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	if !ok {
		t.Errorf("data mismatch: %v, %v", expected, cs.All())
	}
}

func TestOrderGet(t *testing.T) {
	expected := []Expected{
		{name: "Apple", duration: 5 * time.Second},
		{name: "Ball", duration: 4 * time.Second},
		{name: "Cat", duration: 3 * time.Second},
	}

	cs := goqueuetano.Customers{}
	cs.Add(
		goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Ball", Duration: 4 * time.Second},
		goqueuetano.Customer{Name: "Cat", Duration: 3 * time.Second},
	)
	c, err := cs.GetByKey(0)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	c = cs.Get(c.ID())
	if c == (goqueuetano.Customer{}) {
		t.Errorf("customers is expected to have a value: %v", c)
	}

	ok, err := integrity_test(expected, &cs)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	if !ok {
		t.Errorf("data mismatch: %v, %v", expected, cs.All())
	}
}

func TestOrderEdit(t *testing.T) {
	expected := []Expected{
		{name: "Apple", duration: 5 * time.Second},
		{name: "Boy", duration: 4 * time.Second},
		{name: "Cat", duration: 3 * time.Second},
	}

	cs := goqueuetano.Customers{}
	cs.Add(
		goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Ball", Duration: 4 * time.Second},
		goqueuetano.Customer{Name: "Cat", Duration: 3 * time.Second},
	)

	// intended operation
	c, err := cs.GetByKey(1)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	c.Name = "Boy"
	c.Duration = 4 * time.Second
	err = cs.Edit(c)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}

	// intended mis-operation
	err = cs.Edit(goqueuetano.Customer{
		Name:     "id less",
		Duration: 4 * time.Second,
	})
	if err.Error() != "invalid UUID length: 0" {
		t.Errorf("unexpected behaviour: %v", err)
	}

	ok, err := integrity_test(expected, &cs)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	if !ok {
		t.Errorf("data mismatch: %v, %v", expected, cs.All())
	}
}

func TestOrderDelete(t *testing.T) {
	expected := []Expected{
		{name: "Apple", duration: 5 * time.Second},
		{name: "Cat", duration: 3 * time.Second},
	}

	cs := goqueuetano.Customers{}
	cs.Add(
		goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Ball", Duration: 4 * time.Second},
		goqueuetano.Customer{Name: "Cat", Duration: 3 * time.Second},
	)

	// intended operation
	c, err := cs.GetByKey(1)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	cs.Delete(c.ID())

	for _, v := range cs.All() {
		if v.Name == "Ball" {
			t.Errorf("customer is not deleted, cs: %v, c: %v", cs.All(), c)
		}
	}

	// should invoke error
	cs = goqueuetano.Customers{}
	if err := cs.Delete(c.ID()); err == nil {
		t.Errorf("cs allows deleting with empty customers")
	} else {
		// expected given error messages
		switch err.Error() {
		case "invalid UUID length: 0":
		case "customer list is empty":
		default:
			t.Errorf("unexpected error: %v", err)
		}
	}

	// should invoke error
	if err := cs.Delete("1"); err == nil {
		t.Errorf("delete should give error")
	} else {
		if err.Error() != "invalid UUID length: 1" {
			t.Errorf("unexpected error")
		}
	}

	ok, err := integrity_test(expected, &cs)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	if !ok {
		t.Errorf("data mismatch: %v, %v", expected, cs.All())
	}
}

func TestOrderEditAfterDelete(t *testing.T) {
	expected := []Expected{}

	cs := goqueuetano.Customers{}
	cs.Add(goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second})

	// intended operation
	c, err := cs.GetByKey(0)
	if err != nil {
		if err.Error() != "key out of range from the list [1]" {
			t.Errorf("unexpected behaviour: %v", err)
		}
	}
	cs.Delete(c.ID())
	err = cs.Edit(c)
	if err != nil {
		// expected given error messages
		switch err.Error() {
		case "customer list is empty":
		default:
			t.Errorf("unexpected error: %v", err)
		}
	}

	ok, err := integrity_test(expected, &cs)
	if err != nil {
		t.Errorf("unexpected behaviour: %v", err)
	}
	if !ok {
		t.Errorf("data mismatch: %v, %v", expected, cs.All())
	}
}
