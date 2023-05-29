package main

import (
	"fmt"
	"log"
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
