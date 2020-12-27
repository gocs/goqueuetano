package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gocs/goqueuetano/web"
	"github.com/gorilla/csrf"
)

var (
	csrfKey = flag.String("K", "byte array", "csrf key")
)

func main() {
	a := web.New()

	r := chi.NewRouter()
	csrfKey := []byte(*csrfKey)
	r.Use(csrf.Protect(csrfKey, csrf.Secure(false)))
	r.Use(middleware.Logger)

	r.Get("/", web.HomePage(*a))

	r.Get("/add", web.AddPage())
	r.Post("/add", web.AddForm(a))

	r.Get("/edit", web.EditPage(*a))
	r.Post("/edit", web.EditForm(a))

	r.Post("/delete", web.DeleteForm(a))

	// ws
	r.Get("/ws", web.RemainingRealTime(a))

	log.Println("entered...")
	http.ListenAndServe(":3000", r)
	log.Println("exited...")
}
