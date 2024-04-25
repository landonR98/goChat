package controler

import (
	"encoding/json"
	"html/template"
	"landonRyan/goChat/model"
	"landonRyan/goChat/model/db"
	"landonRyan/goChat/util"
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
	user, session := util.GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		InviteId int
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	chatDAO := db.NewChatAccessor(db.GetConection())
	err = chatDAO.AcceptInvite(bodyStruct.InviteId, user.Id)
	if err == db.ErrInviteNotFound {
		res.WriteHeader(http.StatusUnprocessableEntity)
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}

func handleSendInvite(res http.ResponseWriter, req *http.Request) {
	user, session := util.GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		UserId int
		ChatId int
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	chatDao := db.NewChatAccessor(db.GetConection())
	switch chatDao.SendInvite(user.Id, bodyStruct.ChatId, bodyStruct.UserId) {
	case db.ErrDuplicateInvite, nil:
		res.WriteHeader(http.StatusOK)
	case db.ErrUserNotInChat:
		res.WriteHeader(http.StatusUnauthorized)
	default:
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func handleCreateChatroom(res http.ResponseWriter, req *http.Request) {
	user, session := util.GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		Name string
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	chatDAO := db.NewChatAccessor(db.GetConection())
	err = chatDAO.AddChat(bodyStruct.Name, user.Id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}

func handleGetInvites(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	user, session := util.GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusUnauthorized)
		return
	}
	cAcc := db.NewChatAccessor(db.GetConection())
	chatRooms, err := cAcc.GetChatInvitesByUserId(user.Id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	invitesJson, err := json.Marshal(chatRooms)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(invitesJson)
}

func handleGetUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	user, session := util.GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}
	uAcc := db.NewUserAccessor(db.GetConection())
	userList, err := uAcc.GetAllButOne(user.Id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	usersJson, err := json.Marshal(userList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(usersJson)
}

func handleGetChatrooms(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	user, session := util.GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}
	cAcc := db.NewChatAccessor(db.GetConection())
	chatList, err := cAcc.GetRoomsByUserId(user.Id)
	chatsJson, err := json.Marshal(chatList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Write(chatsJson)
}

func handleGetMessages(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	user, session := util.GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		Id          int
		LastMessage int
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	util.CheckErr(err)

	cAcc := db.NewChatAccessor(db.GetConection())

	ch := make(chan bool)
	go cAcc.IsUserInChatConcurrent(user.Id, bodyStruct.Id, ch)

	messageList, err := cAcc.GetChatMessages(bodyStruct.Id)

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
			messageList = make([]model.ChatMessage, 0)
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
	user, session := util.GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	bodyStruct := struct {
		Message string
		ChatId  int
	}{}
	err := json.NewDecoder(req.Body).Decode(&bodyStruct)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	cAcc := db.NewChatAccessor(db.GetConection())
	err = cAcc.NewMessage(bodyStruct.ChatId, user.Id, bodyStruct.Message)
	if err == db.ErrUserNotInChat {
		res.WriteHeader(http.StatusUnauthorized)
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}

func RegisterMessengerRoutes(router *mux.Router, templates *template.Template) {
	router.HandleFunc("/messenger", func(res http.ResponseWriter, req *http.Request) {
		session, err := util.SessionStore.Get(req, "goChat")
		util.CheckErr(err)
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
