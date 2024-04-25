package db

import (
	"database/sql"
	"errors"
	"fmt"
	"landonRyan/goChat/model"
)

var (
	// ErrUserNotFound error used when requested user don't exist.
	ErrUserNotFound = errors.New("User not found.")
)

// userAccessor contains functions for accessing and manipulating database tables related to users.
type userAccessor struct {
	db *sql.DB
}

// NewUserAccessor returns a new userAccessor
func NewUserAccessor(db *sql.DB) userAccessor {
	return userAccessor{db: db}
}

// GetUserByName gets user with provided userName.
func (u *userAccessor) GetUserByName(userName string) (error, model.User) {
	var passHash string
	var userId int
	err := u.db.QueryRow("SELECT password_hash, id FROM users WHERE username = ?", userName).Scan(&passHash, &userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound, model.User{}
		} else {
			return err, model.User{}
		}
	} else {
		return nil, model.NewUser(userName, userId, passHash)
	}
}

// UsernameExists checks database for provided username.
// Returns true if it exists.
func (u *userAccessor) UsernameExists(username string) (nameTaken bool) {
	err := u.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM users WHERE username = ?", username).Scan(&nameTaken)
	if err != nil {
		fmt.Println(err)
	}
	return nameTaken
}

// AddUser adds a new user to the database.
// returns true if successful.
func (u *userAccessor) AddUser(name, hash string) bool {
	_, err := u.db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", name, hash)
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}
}

// GetAllButOne gets all users from the database except user with id of userId.
func (u *userAccessor) GetAllButOne(userId int) ([]model.User, error) {
	users, err := u.db.Query("SELECT username, id FROM users where id != ?", userId)
	if err != nil {
		return nil, err
	}
	defer users.Close()
	var userList []model.User
	for users.Next() {
		var user model.User
		if err := users.Scan(&user.Name, &user.Id); err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}
	return userList, nil
}
