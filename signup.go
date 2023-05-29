package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type formErrorData struct {
	NameTaken       bool
	PasswordNoMatch bool
	Name            string
	Password        string
	RepeatPass      string
}

func RegisterSignupRoutes(router *mux.Router, templates *template.Template) {
	router.HandleFunc("/signup", func(res http.ResponseWriter, req *http.Request) {
		var tmplData formErrorData
		templates.ExecuteTemplate(res, "signup.html", tmplData)
	}).Methods("GET")

	router.HandleFunc("/signup", func(res http.ResponseWriter, req *http.Request) {
		var tmplData formErrorData
		tmplData.Name = req.PostFormValue("name")
		tmplData.Password = req.PostFormValue("password")
		tmplData.RepeatPass = req.PostFormValue("repeat password")

		err := DB.QueryRow("SELECT IF(COUNT(*),'true','false') FROM users WHERE username = ?", tmplData.Name).Scan(&tmplData.NameTaken)
		if err != nil {
			fmt.Println(err)
		}

		tmplData.PasswordNoMatch = (tmplData.Password != tmplData.RepeatPass)

		if tmplData.PasswordNoMatch || tmplData.NameTaken {
			templates.ExecuteTemplate(res, "signup.html", tmplData)
			fmt.Printf("%s, %t\n", tmplData.Name, tmplData.NameTaken)
		} else {
			hashedPass := HashPassword(&tmplData.Password)
			fmt.Printf("%s\n", *hashedPass)
			_, err := DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", tmplData.Name, hashedPass)
			if err != nil {
				fmt.Printf("User insert error: %s\n", err)
			} else {
				templates.ExecuteTemplate(res, "signupSuccess.html", tmplData.Name)
			}
		}
	}).Methods("POST")
}
