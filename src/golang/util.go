package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name string
	Id   int
}

func (user *User) InChatAsync(chatId int, ch chan bool) {
	var exists bool
	sql := "SELECT count(1) FROM chats c JOIN chat_participants cp ON cp.chat_id = c.id where cp.user_id = ? AND cp.accepted_invite = 1 AND c.id = ? LIMIT 1"
	CheckErr(DB.QueryRow(sql, user.Id, chatId).Scan(&exists))
	ch <- exists
	close(ch)
}

func (user *User) InChat(chatId int) bool {
	var userInChat bool
	sql := "SELECT count(1) FROM chats c JOIN chat_participants cp ON cp.chat_id = c.id where cp.user_id = ? AND cp.accepted_invite = 1 AND c.id = ? LIMIT 1"
	CheckErr(DB.QueryRow(sql, user.Id, chatId).Scan(&userInChat))
	return userInChat
}

var SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func LoadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file. Err: %s", err)
	}
}

func HashPassword(password *string) *string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("PASSWD_SALT")+*password), 14)
	if err != nil {
		fmt.Printf("hash password Error: %s", err)
	}
	hashString := string(bytes)
	return &hashString
}

func CheckPasswordHash(password, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(os.Getenv("PASSWD_SALT")+*password))
	return err == nil
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func GetUserSession(req *http.Request) (*User, *sessions.Session) {
	var user User
	session, err := SessionStore.Get(req, "goChat")
	CheckErr(err)
	user.Id = session.Values["userId"].(int)
	user.Name = session.Values["name"].(string)
	return &user, session
}

func SendJson(res http.ResponseWriter, req *http.Request, makeQuery func(*User) []any) {
	res.Header().Set("Content-Type", "application/json")
	user, session := GetUserSession(req)
	if session.IsNew {
		http.Error(res, "you are not logged in", http.StatusBadRequest)
		return
	}

	queryList := makeQuery(user)

	invitesJson, err := json.Marshal(queryList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Write(invitesJson)
}
