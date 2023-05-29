package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration the server waits for existing connections to finish")
	flag.Parse()

	LoadEnv()

	DatabaseInit()

	templates := template.Must(template.ParseGlob("templates/*.html"))

	router := mux.NewRouter()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	RegisterSignupRoutes(router, templates)
	RegisterLoginRoutes(router, templates)
	RegisterMessengerRoutes(router, templates)

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		templates.ExecuteTemplate(res, "index.html", nil)
	})

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
