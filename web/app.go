package web

import "github.com/gocs/goqueuetano"

// App holds the state of the customers
type App struct {
	customers goqueuetano.Order
	editID    string
}

func New(cs goqueuetano.Order) *App {
	return &App{customers: cs}
}
