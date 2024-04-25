package model

import "landonRyan/goChat/util"

type User struct {
	Name         string
	Id           int
	passwordHash string
}

func NewUser(name string, id int, passHash string) User {
	return User{Name: name, Id: id, passwordHash: passHash}
}

func (u *User) SetPassword(hash string) {
	u.passwordHash = hash
}

func (u *User) CheckHash(hash string) bool {
	return util.CheckPasswordHash(&hash, &u.passwordHash)
}

type ChatRoom struct {
	Name string
	Id   int
}

type ChatParticepent struct {
	Id     int
	ChatId int
	UserId int
}

type ChatMessage struct {
	Id      int
	Name    string
	Message string
}
