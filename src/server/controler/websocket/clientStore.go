package websocket

import (
	"database/sql"
	"errors"
	"log"
)

type clientStore struct {
	db *sql.DB
}

var store clientStore

type Client struct {
	UserId     int
	ChatRoomId int
}

var (
	ErrClientNotFound = errors.New("Client not found.")
)

// ClientStoreInit initializes the in memory database for managing connected users.
func ClientStoreInit() {
	var err error
	store.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	store.db.Exec(`CREATE TABLE client(
		userId INT PRIMARY KEY
		chatRoomId INT
		)`)
}

// addClient adds a new client or replaces if it already exists.
func (c *clientStore) addClient(client Client) error {
	sql := `INSERT OR REPLACE INTO client
	(userId, chatRoomId)
	VALUES(?,?)`
	_, err := c.db.Exec(sql, client.UserId, client.ChatRoomId)
	return err
}

// getClientsByUserId retrieves client with matching userId.
func (c *clientStore) getClientsByUserId(userId int) (client Client, err error) {
	query := `SELECT userId, chatRoomId FROM client WHERE userId = ?`
	row := c.db.QueryRow(query, userId)
	err = row.Scan(&client.UserId, &client.ChatRoomId)
	if err == sql.ErrNoRows {
		err = ErrClientNotFound
	}
	return client, err
}

// getClientsByChatRoomId retrieves clients in with matching chatRoomId.
func (c *clientStore) getClientsByChatRoomId(chatRoomId int) (clients []Client, err error) {
	query := `SELECT userId, chatRoomId FROM client WHERE chatRoomId = ?`
	rows, err := c.db.Query(query, chatRoomId)
	if err != nil {
		return nil, err
	}

	clients = make([]Client, 0)
	for rows.Next() {
		client := Client{}
		err = rows.Scan(&client.UserId, &client.ChatRoomId)
		if err != nil {
			rows.Close()
			return nil, err
		}

		clients = append(clients, client)
	}
	err = rows.Err()
	return clients, err
}

// updateClient updates the client in the in memory db.
func (c *clientStore) updateClient(client Client) error {
	sql := `UPDATE client
	SET chatRoomId = ?
	WHERE userId = ?`
	_, err := c.db.Exec(sql, client.ChatRoomId, client.UserId)
	return err
}

// removeClientById removes client with matching userId.
func (c *clientStore) removeClientById(userId int) error {
	sql := `DELETE FROM client
	WHERE userId = ? LIMIT 1`
	_, err := c.db.Exec(sql, userId)
	return err
}
