package controler

import (
	"fmt"
	"html/template"
	"landonRyan/goChat/model/db"
	"landonRyan/goChat/util"
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
		uAcc := db.NewUserAccessor(db.GetConection())

		tmplData.Name = req.PostFormValue("name")
		tmplData.Password = req.PostFormValue("password")
		tmplData.RepeatPass = req.PostFormValue("repeat password")
		tmplData.NameTaken = uAcc.UsernameExists(tmplData.Name)
		tmplData.PasswordNoMatch = (tmplData.Password != tmplData.RepeatPass)

		if tmplData.PasswordNoMatch || tmplData.NameTaken {
			templates.ExecuteTemplate(res, "signup.html", tmplData)
			fmt.Printf("%s, %t\n", tmplData.Name, tmplData.NameTaken)
		} else {
			hashedPass := util.HashPassword(&tmplData.Password)
			fmt.Printf("%s\n", *hashedPass)

			if uAcc.AddUser(tmplData.Name, *hashedPass) {
				templates.ExecuteTemplate(res, "signupSuccess.html", tmplData.Name)
			} else {
				http.Error(res, "Failed to create user", http.StatusInternalServerError)
			}
		}
	}).Methods("POST")
}
