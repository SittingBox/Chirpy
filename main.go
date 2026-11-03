package main

import (
	"net/http"

	"time"
)

func main() {
	serveMux := http.NewServeMux()

	serveMux.Handle("/", http.FileServer(http.Dir("/")))

	serveMux.Handle("/assets/logo.png", http.FileServer(http.Dir("/")))

	DefaultServeMux := &serveMux

	server := &http.Server{
		Addr:         ":8080",
		Handler:      *DefaultServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
