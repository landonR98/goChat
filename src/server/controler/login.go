package controler

import (
	"fmt"
	"html/template"
	"landonRyan/goChat/model/db"
	"landonRyan/goChat/util"
	"net/http"

	"github.com/gorilla/mux"
)

type formData struct {
	Name      string
	Password  string
	FormError bool
}

const loginPage string = "login.html"

func handleLogout(res http.ResponseWriter, req *http.Request) {
	util.DestroyUserSession(res, req)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func handleLogin(res http.ResponseWriter, req *http.Request, templates *template.Template) {
	var tmplData formData
	tmplData.Name = req.PostFormValue("name")
	tmplData.Password = req.PostFormValue("password")
	userDAO := db.NewUserAccessor(db.GetConection())

	err, user := userDAO.GetUserByName(tmplData.Name)
	if err != nil {
		if err == db.ErrUserNotFound {
			tmplData.FormError = true
			templates.ExecuteTemplate(res, loginPage, tmplData)
		} else {
			fmt.Println("unknown error in login database query")
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	} else {
		if user.CheckHash(tmplData.Password) {
			fmt.Println("pass match")
			if err := util.SetUserSession(res, req, user.Name, user.Id); err == util.SessionNotCreated {
				http.Error(res, err.Error(), http.StatusInternalServerError)
			} else {
				http.Redirect(res, req, "/messenger", http.StatusSeeOther)
			}
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
