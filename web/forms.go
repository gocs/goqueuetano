package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gocs/goqueuetano"
)

func AddForm(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		// duration's actual value is a datetime-local
		duration := r.FormValue("duration")
		// concat the timezone
		fmtDuration := fmt.Sprintf("%s%s", duration, "+08:00")

		t, err := time.Parse(time.RFC3339Nano, fmtDuration)
		if err != nil {
			log.Println("err:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("cancelled")
			return
		}

		app.customers.Add(goqueuetano.Customer{
			Name:     name,
			Duration: t.Sub(time.Now()),
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func EditForm(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		duration := r.FormValue("duration")

		// concat the timezone
		fmtDuration := fmt.Sprintf("%s%s", duration, "+08:00")
		t, err := time.Parse(time.RFC3339Nano, fmtDuration)
		if err != nil {
			log.Println("err:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("cancelled")
			return
		}

		customer := app.customers.Get(app.editID)
		customer.Name = name
		customer.Duration = t.Sub(time.Now())
		app.customers.Edit(customer)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

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
		app.customers.Delete(c.ID())
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
