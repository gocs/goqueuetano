package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gocs/goqueuetano"
)

// AddForm is a form handler for adding new customer
func AddForm(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		total := r.FormValue("total")
		t, err := strconv.Atoi(total)
		if err != nil {
			log.Println("err:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		app.customers.Add(goqueuetano.Customer{
			Name:  name,
			Total: t,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// EditForm is a form handler for editing the selected customer
func EditForm(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")

		customer := app.customers.Get(app.editID)
		customer.Name = name

		if err := app.customers.Edit(customer); err != nil {
			log.Println("err:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// DeleteForm is a form handler for deleting the selected customer
func DeleteForm(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		k, err := strconv.Atoi(key)
		if err != nil {
			log.Println("err:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("cancelled")
			return
		}
		// the index is decremented in order to input using index of the ordered list
		c, err := app.customers.GetByKey(k - 1)
		if err != nil {
			log.Println("err upon getbykey:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("cancelled")
			return
		}

		if err := app.customers.Delete(c.ID()); err != nil {
			log.Println("err:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
