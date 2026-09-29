package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"calculator/internal/httpapi"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	server := &http.Server{
		Addr: addr, Handler: httpapi.NewAppHandler(os.Getenv("STATIC_DIR")),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second,
	}
	log.Printf("Calculator API listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}
