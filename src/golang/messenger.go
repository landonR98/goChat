package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type messengerPageData struct {
	Name     string
	LoggedIn bool
}

func handleGetInvites(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	session, err := SessionStore.Get(req, "goChat")
	CheckErr(err)
	if session.IsNew {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	userId := session.Values["userId"]
	sql := "SELECT cp.id, c.name FROM chat_participants cp JOIN users u ON u.id = cp.user_id JOIN chats c ON cp.chat_id = c.id WHERE u.id != ?"
	users, err := DB.Query(sql, userId)
	CheckErr(err)
	defer users.Close()
	var userList []string
	for users.Next() {
		var username string
		if err := users.Scan(&username); err != nil {
			CheckErr(err)
		}
		userList = append(userList, username)
	}
	usersJson, err := json.Marshal(userList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Write(usersJson)
}

func handleGetUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	session, err := SessionStore.Get(req, "goChat")
	CheckErr(err)
	if session.IsNew {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	userId := session.Values["userId"]
	users, err := DB.Query("SELECT username FROM users where id != ?", userId)
	CheckErr(err)
	defer users.Close()
	var userList []string
	for users.Next() {
		var username string
		if err := users.Scan(&username); err != nil {
			CheckErr(err)
		}
		userList = append(userList, username)
	}
	usersJson, err := json.Marshal(userList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Write(usersJson)
}

func handleGetChatrooms(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	session, err := SessionStore.Get(req, "goChat")
	CheckErr(err)
	if session.IsNew {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	userId := session.Values["userId"]
	users, err := DB.Query("SELECT c.name, c.id FROM chats c JOIN chat_participants cp ON cp.chat_id = c.id where cp.user_id != ?", userId)
	CheckErr(err)
	defer users.Close()
	var userList []string
	for users.Next() {
		var username string
		if err := users.Scan(&username); err != nil {
			CheckErr(err)
		}
		userList = append(userList, username)
	}
	usersJson, err := json.Marshal(userList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Write(usersJson)
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
	router.HandleFunc("/users", handleGetUsers).Methods("GET")
	router.HandleFunc("/invites", handleGetInvites).Methods("GET")
	router.HandleFunc("/chatrooms", handleGetChatrooms).Methods("GET")
}
