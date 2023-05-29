package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type messengerPageData struct {
	Name     string
	LoggedIn bool
}

func RegisterMessengerRoutes(router *mux.Router, templates *template.Template) {

	router.HandleFunc("/messenger", func(res http.ResponseWriter, req *http.Request) {
		session, err := SessionStore.Get(req, "goChat")
		CheckErr(err)
		if session.IsNew {
			http.Redirect(res, req, "/login", http.StatusSeeOther)
			return
		}

		var tmplData messengerPageData

		tmplData.Name = session.Values["name"].(string)
		tmplData.LoggedIn = true

		fmt.Println(tmplData)

		templates.ExecuteTemplate(res, "messenger.html", tmplData)
	})
}
