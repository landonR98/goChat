package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func CloseDB() {
	db.Close()
}

func DatabaseInitMysql() {
	cfg := mysql.Config{
		User:   os.Getenv("MYSQL_USERNAME"),
		Passwd: os.Getenv("MYSQL_PASSWD"),
		Net:    "tcp",
		Addr:   os.Getenv("MYSQL_ADDR"),
		DBName: "gochat",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Database Connected!")
}

func DatabaseInitSqlite() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createSqliteTables()

	fmt.Println("Sqlite3 database connected!")
}

func createSqliteTables() {
	userSql := `CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL,
  password_hash TEXT NOT NULL
  );`

	chatSql := `CREATE TABLE IF NOT EXISTS chats (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  is_private BOOLEAN NOT NULL DEFAULT '0',
  creator INTEGER NOT NULL,
  FOREIGN KEY (creator) REFERENCES users (id)
  )`

	chatParticipantsSql := `CREATE TABLE IF NOT EXISTS chat_participants (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  chat_id INTEGER NOT NULL,
  accepted_invite BOOLEAN NOT NULL DEFAULT '0',
  FOREIGN KEY (chat_id) REFERENCES chats (id),
  FOREIGN KEY (user_id) REFERENCES users (id)
	)`

	messageSql := `CREATE TABLE IF NOT EXISTS message (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  chat_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  message TEXT NOT NULL,
  FOREIGN KEY (chat_id) REFERENCES chats (id),
  FOREIGN KEY (user_id) REFERENCES users (id)
	)`

	_, err := db.Exec(userSql)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(chatSql)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(chatParticipantsSql)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(messageSql)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConection() *sql.DB {
	return db
}
