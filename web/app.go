package web

import (
	"fmt"
	"log"
	"os"

	"github.com/gocs/goqueuetano"
)

// App holds the state of the customers
type App struct {
	customers goqueuetano.Order
	editID    string
	pages     map[string]string
}

// NewApp instantiates the web app with its html templates
func NewApp() *App {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pages := map[string]string{
		"home": "/public/layout.html",
		"add":  "/public/add.html",
		"edit": "/public/edit.html",
	}
	for k, v := range pages {
		pages[k] = fmt.Sprintf("%s%s", dir, v)
	}
	return &App{
		customers: &goqueuetano.Customers{},
		pages:     pages,
	}
}
