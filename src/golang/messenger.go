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

type chatroomSidebarResponse struct {
	Name string
	Id   int
}

type messageStruct struct {
	Name    string
	Message string
	Id      int
}

func handleAcceptInvite(res http.ResponseWriter, req *http.Request) {
	user, session := GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		InviteId int
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	CheckErr(err)

	var inviteExists bool
	sql := "SELECT count(1) FROM chat_participants where id = ? AND user_id = ? LIMIT 1"
	CheckErr(DB.QueryRow(sql, bodyStruct.InviteId, user.Id).Scan(&inviteExists))

	if inviteExists {
		sql = "UPDATE chat_participants SET accepted_invite = 1 WHERE id = ?"
		_, err = DB.Exec(sql, bodyStruct.InviteId)
		CheckErr(err)
	}
	res.WriteHeader(http.StatusOK)
}

func handleSendInvite(res http.ResponseWriter, req *http.Request) {
	user, session := GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		UserId int
		ChatId int
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	CheckErr(err)

	userCh := make(chan bool)
	go user.InChatAsync(bodyStruct.ChatId, userCh)

	var alreadyInvited bool
	sql := "SELECT count(1) FROM chats c JOIN chat_participants cp ON cp.chat_id = c.id where cp.user_id = ? AND c.id = ? LIMIT 1"
	CheckErr(DB.QueryRow(sql, bodyStruct.UserId, bodyStruct.ChatId).Scan(&alreadyInvited))

	userInChat := <-userCh

	if userInChat && !alreadyInvited {
		sql = "INSERT INTO chat_participants (user_id, chat_id) VALUES(?,?)"
		_, err = DB.Exec(sql, bodyStruct.UserId, bodyStruct.ChatId)
		CheckErr(err)
	}
	res.WriteHeader(http.StatusOK)
}

func handleCreateChatroom(res http.ResponseWriter, req *http.Request) {
	user, session := GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		Name string
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	CheckErr(err)

	sql := "INSERT INTO chats (name, creator) VALUES(?,?)"
	query, err := DB.Exec(sql, bodyStruct.Name, user.Id)
	CheckErr(err)
	insertId, err := query.LastInsertId()
	CheckErr(err)
	sql = "INSERT INTO chat_participants (user_id, chat_id, accepted_invite) VALUES(?,?,1)"
	_, err = DB.Exec(sql, user.Id, insertId)
	CheckErr(err)

	res.WriteHeader(http.StatusOK)
}

func handleGetInvites(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	user, session := GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}
	sql := "SELECT cp.id, c.name FROM chat_participants cp JOIN chats c ON cp.chat_id = c.id WHERE cp.user_id = ? AND cp.accepted_invite = 0"
	inviteQuery, err := DB.Query(sql, user.Id)
	CheckErr(err)
	defer inviteQuery.Close()
	var inviteList []chatroomSidebarResponse
	for inviteQuery.Next() {
		var invite chatroomSidebarResponse
		if err := inviteQuery.Scan(&invite.Id, &invite.Name); err != nil {
			CheckErr(err)
		}
		inviteList = append(inviteList, invite)
	}
	invitesJson, err := json.Marshal(inviteList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Write(invitesJson)
}

func handleGetUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	user, session := GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}
	users, err := DB.Query("SELECT username, id FROM users where id != ?", user.Id)
	CheckErr(err)
	defer users.Close()
	var userList []chatroomSidebarResponse
	for users.Next() {
		var user chatroomSidebarResponse
		if err := users.Scan(&user.Name, &user.Id); err != nil {
			CheckErr(err)
		}
		userList = append(userList, user)
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
	user, session := GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}
	sql := "SELECT c.name, c.id FROM chats c JOIN chat_participants cp ON cp.chat_id = c.id where cp.user_id = ? AND cp.accepted_invite = 1;"
	chatQuery, err := DB.Query(sql, user.Id)
	CheckErr(err)
	defer chatQuery.Close()
	var chatList []chatroomSidebarResponse
	for chatQuery.Next() {
		var chat chatroomSidebarResponse
		if err := chatQuery.Scan(&chat.Name, &chat.Id); err != nil {
			CheckErr(err)
		}
		chatList = append(chatList, chat)
	}
	chatsJson, err := json.Marshal(chatList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Write(chatsJson)
}

func handleGetMessages(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	user, session := GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		Id          int
		LastMessage int
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	CheckErr(err)

	ch := make(chan bool)
	go user.InChatAsync(bodyStruct.Id, ch)

	sql := "SELECT u.username, m.message, m.id FROM message m JOIN users u on m.user_id = u.id WHERE m.chat_id = ? ORDER BY m.id"
	messageQuery, err := DB.Query(sql, bodyStruct.Id)
	CheckErr(err)
	defer messageQuery.Close()
	var messageList []messageStruct
	for messageQuery.Next() {
		var message messageStruct
		if err := messageQuery.Scan(&message.Name, &message.Message, &message.Id); err != nil {
			CheckErr(err)
		}
		messageList = append(messageList, message)
	}

	userInChat := <-ch
	if !userInChat {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if bodyStruct.LastMessage != -1 {
		var lastMessageIndex int
		for i, message := range messageList {
			if message.Id == bodyStruct.LastMessage {
				lastMessageIndex = i
			}
		}
		if len(messageList) != lastMessageIndex {
			messageList = messageList[lastMessageIndex+1:]
		} else {
			messageList = make([]messageStruct, 0)
		}
	}

	messageJson, err := json.Marshal(messageList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Write(messageJson)
}

func handleSendMessage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	user, session := GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		Message string
		ChatId  int
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	CheckErr(err)

	if !user.InChat(bodyStruct.ChatId) {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	sql := "INSERT INTO message (chat_id, user_id, message) VALUES(?,?,?)"
	result, err := DB.Exec(sql, bodyStruct.ChatId, user.Id, bodyStruct.Message)
	CheckErr(err)
	fmt.Println(result.RowsAffected())

	res.WriteHeader(http.StatusOK)
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

		templates.ExecuteTemplate(res, "messenger.html", tmplData)
	})
	router.HandleFunc("/users", handleGetUsers).Methods("GET")
	router.HandleFunc("/invites", handleGetInvites).Methods("GET")
	router.HandleFunc("/chatrooms", handleGetChatrooms).Methods("GET")
	router.HandleFunc("/messages", handleGetMessages).Methods("POST")
	router.HandleFunc("/send", handleSendMessage).Methods("POST")
	router.HandleFunc("/createChatroom", handleCreateChatroom).Methods("POST")
	router.HandleFunc("/sendInvite", handleSendInvite).Methods("POST")
	router.HandleFunc("/acceptInvite", handleAcceptInvite).Methods("POST")
}
