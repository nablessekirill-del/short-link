package main

import (
	"log"
	"net/http"
	"shortener/internal/handler"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func main() {
	store := storage.NewMemory()
	svc := service.New(store)

	addr := ":8080"
	h := handler.New(svc, "http://localhost:8080")

	mux := http.NewServeMux()
	h.Routes(mux)

	log.Printf("server listening om %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
