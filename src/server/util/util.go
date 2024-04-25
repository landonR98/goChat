package util

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

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

type userSession = struct {
	Id   int
	Name string
}

var SessionNotCreated = errors.New("Session not created")

func SetUserSession(res http.ResponseWriter, req *http.Request, name string, id int) error {
	session, err := SessionStore.Get(req, "goChat")
	if err != nil {
		fmt.Printf("Error creating login session. Err: %s\n", err)
		return SessionNotCreated
	}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
	}

	session.Values["name"] = name
	session.Values["userId"] = id

	err = session.Save(req, res)
	if err != nil {
		return SessionNotCreated
	}
	return nil
}

func DestroyUserSession(res http.ResponseWriter, req *http.Request) {
	_, session := GetUserSession(req)
	session.Options.MaxAge = -1
	session.Options.SameSite = http.SameSiteStrictMode
	session.Save(req, res)
}

func GetUserSession(req *http.Request) (userSession, *sessions.Session) {
	session, err := SessionStore.Get(req, "goChat")
	CheckErr(err)
	id := session.Values["userId"].(int)
	name := session.Values["name"].(string)
	user := userSession{Id: id, Name: name}
	return user, session
}
