package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type formData struct {
	Name      string
	Password  string
	FormError bool
}

const loginPage string = "login.html"

func handleLogout(res http.ResponseWriter, req *http.Request) {
	_, session := GetUserSession(req)
	session.Options.MaxAge = -1
	session.Save(req, res)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func handleLogin(res http.ResponseWriter, req *http.Request, templates *template.Template) {
	var tmplData formData
	var passHash string
	var userId int
	tmplData.Name = req.PostFormValue("name")
	tmplData.Password = req.PostFormValue("password")

	err := DB.QueryRow("SELECT password_hash, id FROM users WHERE username = ?", tmplData.Name).Scan(&passHash, &userId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no rows")
			tmplData.FormError = true
			templates.ExecuteTemplate(res, loginPage, tmplData)
		} else {
			fmt.Println("unknown error in login database query")
		}
	} else {
		if CheckPasswordHash(&tmplData.Password, &passHash) {
			fmt.Println("pass match")
			session, err := SessionStore.Get(req, "goChat")
			if err != nil {
				fmt.Printf("Error creating login session. Err: %s\n", err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			session.Options = &sessions.Options{
				Path:   "/",
				MaxAge: 86400,
			}

			session.Values["name"] = tmplData.Name
			session.Values["userId"] = userId

			err = session.Save(req, res)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(res, req, "/messenger", http.StatusSeeOther)
		} else {
			tmplData.FormError = true
			templates.ExecuteTemplate(res, loginPage, tmplData)
		}
	}

}

func RegisterLoginRoutes(router *mux.Router, templates *template.Template) {
	router.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) {
		var tmplData formData
		templates.ExecuteTemplate(res, loginPage, tmplData)
	}).Methods("GET")

	router.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) { handleLogin(res, req, templates) }).Methods("POST")
	router.HandleFunc("/logout", handleLogout)
}
