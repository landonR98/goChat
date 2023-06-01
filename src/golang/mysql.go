package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DatabaseInit() {
	cfg := mysql.Config{
		User:   os.Getenv("MYSQL_USERNAME"),
		Passwd: os.Getenv("MYSQL_PASSWD"),
		Net:    "tcp",
		Addr:   os.Getenv("MYSQL_ADDR"),
		DBName: "gochat",
	}
	// Get a database handle.
	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Database Connected!")
}
