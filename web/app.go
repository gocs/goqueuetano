package web

import "github.com/gocs/goqueuetano"

// App holds the state of the customers
type App struct {
	customers goqueuetano.Order
	editID    string
}

func New() *App {
	return &App{customers: &goqueuetano.Customers{}}
}
