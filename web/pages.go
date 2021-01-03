package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gocs/goqueuetano"
	"github.com/gorilla/csrf"
)

// HomePage is the landing page of the app
func HomePage(app App) http.HandlerFunc {
	type Data struct {
		CustomerNotEmpty bool
		CustomerSize     int
		CSRF             template.HTML
	}
	tmp := template.Must(template.ParseFiles(app.pages["home"]))
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

// AddPage is the template page where you can add a customer
func AddPage(app App) http.HandlerFunc {
	type Data struct {
		Today string
		CSRF  template.HTML
	}
	tmp := template.Must(template.ParseFiles(app.pages["add"]))
	return func(w http.ResponseWriter, r *http.Request) {
		data := Data{
			CSRF: csrf.TemplateField(r),
		}
		if err := tmp.Execute(w, data); err != nil {
			log.Println("err in addPage:", err)
		}
	}
}

// EditPage is the template page where you can add a customer
func EditPage(app *App) http.HandlerFunc {
	type Data struct {
		Customer goqueuetano.Customer
		CSRF     template.HTML
	}
	tmp := template.Must(template.ParseFiles(app.pages["edit"]))
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
		data := Data{
			// the index is decremented in order to input using index of the ordered list
			Customer: customer,
			CSRF:     csrf.TemplateField(r),
		}
		if err := tmp.Execute(w, data); err != nil {
			log.Println("err in editPage:", err)
		}
	}
}
