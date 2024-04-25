package db

import (
	"database/sql"
	"errors"
	"fmt"
	"landonRyan/goChat/model"
)

// chatAccessor contains functions for accessing and manipulating database tables related to chat rooms
type chatAccessor struct {
	db *sql.DB
}

// NewChatAccessor returns a new chatDAO object
// db *sql.DB sql database connection
func NewChatAccessor(db *sql.DB) chatAccessor {
	return chatAccessor{db: db}
}

var (
	// ErrInviteNotFound error used when requested invite don't exist
	ErrInviteNotFound = errors.New("Invite not found")
	// ErrDuplicateInvite error used when user has already been invited to a chat room.
	ErrDuplicateInvite = errors.New("Invite already exists")
	// ErrUserNotInChat error used when user is trying to access a chat room they are not in.
	ErrUserNotInChat = errors.New("User not in chat")
)

// IsUserInChatConcurrent sends true into ch if userId is in chatId. Or false if not.
func (u *chatAccessor) IsUserInChatConcurrent(userId, chatId int, ch chan bool) {
	var exists bool
	sql := "SELECT count(1) FROM chats c JOIN chat_participants cp ON cp.chat_id = c.id where cp.user_id = ? AND cp.accepted_invite = 1 AND c.id = ? LIMIT 1"
	err := u.db.QueryRow(sql, userId, chatId).Scan(&exists)
	if err != nil {
		exists = false
		fmt.Println(err)
	}
	ch <- exists
	close(ch)
}

// IsUserInChat returns true if userId is in chatId.
func (u *chatAccessor) IsUserInChat(userId, chatId int) bool {
	var userInChat bool
	sql := "SELECT count(1) FROM chats c JOIN chat_participants cp ON cp.chat_id = c.id where cp.user_id = ? AND cp.accepted_invite = 1 AND c.id = ? LIMIT 1"
	err := u.db.QueryRow(sql, userId, chatId).Scan(&userInChat)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return userInChat
}

// AcceptInvite marks an invite as accepted by userId
// returns ErrInvalidInvite if invite is not found
// also can returns sql errors
// returns nil if successful
func (c *chatAccessor) AcceptInvite(inviteId, userId int) error {
	var inviteExists bool
	sql := "SELECT count(1) FROM chat_participants where id = ? AND user_id = ? LIMIT 1"
	err := c.db.QueryRow(sql, inviteId, userId).Scan(&inviteExists)
	if err != nil {
		return err
	}

	if inviteExists {
		sql = "UPDATE chat_participants SET accepted_invite = 1 WHERE id = ?"
		_, err := c.db.Exec(sql, inviteId)
		return err
	} else {
		return ErrInviteNotFound
	}
}

// SendInvite sends invite from user with id of fromId to user with toUserId for chat room with chatId.
// Can return with ErrUserNotInChat if fromId is not in chatId.
// Can also return sql errors
func (c *chatAccessor) SendInvite(fromId, chatId, toUserId int) error {
	userInChatCn := make(chan bool)
	c.IsUserInChatConcurrent(fromId, chatId, userInChatCn)

	var alreadyInvited bool
	sql := "SELECT count(1) FROM chats c JOIN chat_participants cp ON cp.chat_id = c.id where cp.user_id = ? AND c.id = ? LIMIT 1"
	err := c.db.QueryRow(sql, toUserId, chatId).Scan(&alreadyInvited)
	if err != nil {
		return err
	}

	if !<-userInChatCn {
		return ErrUserNotInChat
	}

	if !alreadyInvited {
		sql = "INSERT INTO chat_participants (user_id, chat_id) VALUES(?,?)"
		_, err := c.db.Exec(sql, toUserId, chatId)
		return err
	} else {
		return ErrDuplicateInvite
	}
}

// AddChat adds a new chat room to the database.
// Can return sql errors.
func (c *chatAccessor) AddChat(chatName string, creatorId int) error {
	sql := "INSERT INTO chats (name, creator) VALUES(?,?)"
	query, err := c.db.Exec(sql, chatName, creatorId)
	if err != nil {
		return err
	}
	insertId, err := query.LastInsertId()
	if err != nil {
		return err
	}
	sql = "INSERT INTO chat_participants (user_id, chat_id, accepted_invite) VALUES(?,?,1)"
	_, err = c.db.Exec(sql, creatorId, insertId)
	return err
}

// GetChatInvitesByUserId returns all chat rooms that userId has been invited to.
// Can return sql errors.
func (c *chatAccessor) GetChatInvitesByUserId(userId int) ([]model.ChatRoom, error) {
	sql := "SELECT cp.id, c.name FROM chat_participants cp JOIN chats c ON cp.chat_id = c.id WHERE cp.user_id = ? AND cp.accepted_invite = 0"
	inviteQuery, err := c.db.Query(sql, userId)
	if err != nil {
		return nil, err
	}
	defer inviteQuery.Close()
	var inviteList []model.ChatRoom
	for inviteQuery.Next() {
		var invite model.ChatRoom
		if err := inviteQuery.Scan(&invite.Id, &invite.Name); err != nil {
			return nil, err
		}
		inviteList = append(inviteList, invite)
	}
	return inviteList, nil
}

// GetRoomsByUserId returns all chat rooms where userId is a participant.
// Can return sql errors
func (c *chatAccessor) GetRoomsByUserId(userId int) ([]model.ChatRoom, error) {
	sql := "SELECT c.name, c.id FROM chats c JOIN chat_participants cp ON cp.chat_id = c.id where cp.user_id = ? AND cp.accepted_invite = 1;"
	chatQuery, err := c.db.Query(sql, userId)
	if err != nil {
		return nil, err
	}
	defer chatQuery.Close()
	var chatList []model.ChatRoom
	for chatQuery.Next() {
		var chat model.ChatRoom
		if err := chatQuery.Scan(&chat.Name, &chat.Id); err != nil {
			return nil, err
		}
		chatList = append(chatList, chat)
	}
	return chatList, nil
}

// GetChatMessages returns all messages from chat room with id of chatId.
// Can return sql errors.
func (c *chatAccessor) GetChatMessages(chatId int) ([]model.ChatMessage, error) {
	sql := "SELECT u.username, m.message, m.id FROM message m JOIN users u on m.user_id = u.id WHERE m.chat_id = ? ORDER BY m.id"
	messageQuery, err := c.db.Query(sql, chatId)
	if err != nil {
		return nil, err
	}
	defer messageQuery.Close()
	var messageList []model.ChatMessage
	for messageQuery.Next() {
		var message model.ChatMessage
		if err := messageQuery.Scan(&message.Name, &message.Message, &message.Id); err != nil {
			return nil, err
		}
		messageList = append(messageList, message)
	}
	return messageList, nil
}

// NewMessage adds message from userId to chatId to the database.
// Can return ErrUserNotInChat if userId is not in chatId.
// Can also return an sql error
func (c *chatAccessor) NewMessage(chatId, userId int, message string) error {
	if !c.IsUserInChat(userId, chatId) {
		return ErrUserNotInChat
	}
	sql := "INSERT INTO message (chat_id, user_id, message) VALUES(?,?,?)"
	result, err := c.db.Exec(sql, chatId, userId, message)
	fmt.Println(result.RowsAffected())
	return err
}
