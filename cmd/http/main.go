package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gocs/goqueuetano"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/gorilla/websocket"
)

var (
	csrfKey = flag.String("K", "byte array", "csrf key")
)

func main() {
	a := App{
		customers: goqueuetano.NewCustomers(),
	}

	r := chi.NewRouter()
	csrfKey := []byte(*csrfKey)
	r.Use(csrf.Protect(csrfKey, csrf.Secure(false)))
	r.Use(middleware.Logger)

	r.Get("/", homePage(a))

	r.Get("/add", addPage(a))
	r.Post("/add", addForm(a))

	r.Get("/edit", editPage(a))
	r.Post("/edit", editForm(a))

	r.Post("/delete", deleteForm(a))

	// ws
	r.Get("/ws", RemainingRealTime(a))

	http.ListenAndServe(":3000", r)
}

// App holds the state of the customers
type App struct {
	customers goqueuetano.Order
}

func homePage(app App) http.HandlerFunc {
	type Data struct {
		CustomerNotEmpty bool
		CustomerSize     int
		CSRF             template.HTML
	}
	tmp := template.Must(template.ParseFiles("./public/layout.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		data := Data{
			CustomerNotEmpty: app.customers.Len() > 0,
			CustomerSize:     app.customers.Len(),
			CSRF:             csrf.TemplateField(r),
		}

		if err := tmp.Execute(w, data); err != nil {
			log.Println("err in homePage:", err)
		}
	}
}

// RemainingRealTime gives real-time remaining time to client
func RemainingRealTime(app App) http.HandlerFunc {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	type Customer struct {
		Index     string        `json:"i"`
		Name      string        `json:"name"`
		Duration  time.Duration `json:"duration"`
		Countdown time.Duration `json:"countdown"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("err:", err)
			return
		}

		for {
			// prevent instant refresh
			ticker := time.NewTicker(50 * time.Millisecond)
			for range ticker.C {
				var cs []Customer
				for _, c := range app.customers.All() {
					cs = append(cs, Customer{
						Index:     c.ID,
						Name:      c.Name,
						Duration:  c.Duration,
						Countdown: c.Countdown(),
					})
				}

				msg, err := json.Marshal(cs)
				if err != nil {
					log.Println("err:", err)
					return
				}

				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func addPage(app App) http.HandlerFunc {
	type Data struct {
		Today string
		CSRF  template.HTML
	}
	tmp := template.Must(template.ParseFiles("./public/add.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		layout := "2006-01-02T15:04:05"
		data := Data{
			Today: time.Now().Format(layout),
			CSRF:  csrf.TemplateField(r),
		}
		if err := tmp.Execute(w, data); err != nil {
			log.Println("err in addPage:", err)
		}
	}
}

func addForm(app App) http.HandlerFunc {
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
		log.Println("redirect to `/`")
	}
}

func editPage(app App) http.HandlerFunc {
	type Data struct {
		Customer goqueuetano.Customer
		DeadLine string
		Today    string
		CSRF     template.HTML
	}
	tmp := template.Must(template.ParseFiles("./public/edit.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		k, err := strconv.Atoi(key)
		if err != nil {
			log.Println("err:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("cancelled")
			return
		}

		customer := app.customers.GetByKey(k - 1)
		layout := "2006-01-02T15:04:05"
		data := Data{
			// the index is decremented in order to input using index of the ordered list
			Customer: customer,
			DeadLine: time.Now().Add(customer.Duration).Format(layout),
			Today:    time.Now().Format(layout),
			CSRF:     csrf.TemplateField(r),
		}
		if err := tmp.Execute(w, data); err != nil {
			log.Println("err in editPage:", err)
		}
	}
}

func editForm(app App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := r.FormValue("id")
		name := r.FormValue("name")
		duration := r.FormValue("duration")

		id, err := uuid.Parse(ID)
		if err != nil {
			log.Println("err:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("cancelled")
			return
		}

		// concat the timezone
		fmtDuration := fmt.Sprintf("%s%s", duration, "+08:00")
		t, err := time.Parse(time.RFC3339Nano, fmtDuration)
		if err != nil {
			log.Println("err:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("cancelled")
			return
		}

		app.customers.Edit(goqueuetano.Customer{
			ID:       id.String(),
			Name:     name,
			Duration: t.Sub(time.Now()),
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("redirect to `/`")
	}
}

func deleteForm(app App) http.HandlerFunc {
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
		c := app.customers.GetByKey(k - 1)
		app.customers.Delete(c.ID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("redirect to `/`")
	}
}
