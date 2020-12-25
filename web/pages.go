package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gocs/goqueuetano"
	"github.com/gorilla/csrf"
)

func HomePage(app App) http.HandlerFunc {
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

func AddPage() http.HandlerFunc {
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

func EditPage(app App) http.HandlerFunc {
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
			log.Println("err upon getbykey:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("cancelled")
			return
		}

		customer, err := app.customers.GetByKey(k - 1)
		if err != nil {
			log.Println("err upon getbykey:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("cancelled")
			return
		}
		app.editID = customer.ID()
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
