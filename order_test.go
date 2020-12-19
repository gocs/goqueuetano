package goqueuetano_test

import (
	"testing"
	"time"

	"github.com/gocs/goqueuetano"
)

func TestOrderAdd(t *testing.T) {
	cs := goqueuetano.Customers{}

	cs.Add(goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second})
	if len(cs) == 0 {
		t.Errorf("customers must not be empty")
	}

	cs.Add(
		goqueuetano.Customer{Name: "Ball", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Cat", Duration: 5 * time.Second},
	)
	if len(cs) != 3 {
		t.Errorf("customers is expected to have 3 elements")
	}
}

func TestOrderGet(t *testing.T) {
	cs := goqueuetano.Customers{}
	cs.Add(
		goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Ball", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Cat", Duration: 5 * time.Second},
	)
	id := cs[2].ID()
	customer := cs.Get(id)
	if customer.Name != "Cat" {
		t.Errorf("They should be equal: %v", customer)
	}
}

func TestOrderEdit(t *testing.T) {
	cs := goqueuetano.Customers{}
	cs.Add(
		goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Ball", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Cat", Duration: 5 * time.Second},
	)
	clone := goqueuetano.Customers{}
	for i := 0; i < cs.Len(); i++ {
		clone = append(clone, cs[i])
	}

	// intended operation
	c := cs[2]
	c.Name = "Boy"
	c.Duration = 4 * time.Second
	cs.Edit(c)

}

func TestOrderDelete(t *testing.T) {
	cs := goqueuetano.Customers{}
	cs.Add(
		goqueuetano.Customer{Name: "Apple", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Ball", Duration: 5 * time.Second},
		goqueuetano.Customer{Name: "Cat", Duration: 5 * time.Second},
	)

	// intended operation
	id := cs[1].ID()
	cs.Delete(id)

	if cs[1].Name == "Ball" {
		t.Errorf("customer is not deleted")
	}

	// should invoke error
	cs = goqueuetano.Customers{}
	if err := cs.Delete(id); err == nil {
		t.Errorf("cs allows deleting with empty customers")
	} else {
		t.Log(err)
	}

	// should invoke error
	if err := cs.Delete("1"); err == nil {
		t.Errorf("delete should give error")
	} else {
		t.Log(err)
	}
}

func TestNewCustomer(t *testing.T) {
	cs := *goqueuetano.NewCustomers()
	cs.Add(
		goqueuetano.Customer{Name: "Apple", Duration: 1 * time.Second},
		goqueuetano.Customer{Name: "Ball", Duration: 1 * time.Second},
		goqueuetano.Customer{Name: "Cat", Duration: 1 * time.Second},
	)
	time.Sleep(5 * time.Second)
	if len(cs) > 0 {
		t.Errorf("customers are not automatically deleted: %v", cs)
	}
}
