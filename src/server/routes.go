package main

import (
	"html/template"
	"landonRyan/goChat/controler"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	templates := template.Must(template.ParseGlob("templates/*.html"))

	router := mux.NewRouter()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	controler.RegisterSignupRoutes(router, templates)
	controler.RegisterLoginRoutes(router, templates)
	controler.RegisterMessengerRoutes(router, templates)

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		templates.ExecuteTemplate(res, "index.html", nil)
	})
	return router

}
