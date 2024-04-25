package main

import (
	"context"
	"flag"
	"fmt"
	"landonRyan/goChat/model/db"
	"landonRyan/goChat/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration the server waits for existing connections to finish")
	useSqlite := flag.Bool("sqlite", false, "use sqlite database")
	flag.Parse()

	util.LoadEnv()

	if *useSqlite {
		db.DatabaseInitSqlite()
	} else {
		db.DatabaseInitMysql()
	}
	defer db.CloseDB()

	router := NewRouter()

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()
	fmt.Printf("starting server on port %s\n", server.Addr)

	// gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println()
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
