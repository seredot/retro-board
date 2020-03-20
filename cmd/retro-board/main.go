package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	repo := NewMemoryRepo()
	handler := NewHandler(repo)
	router := mux.NewRouter()
	mapHandlerFuncs(router, handler)

	log.Println("Running server")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func mapHandlerFuncs(r *mux.Router, handler Handler) {
	r.HandleFunc("/api", handler.healthCheck).Methods("GET")
	r.HandleFunc("/api/board", handler.createBoard).Methods("POST")
	r.HandleFunc("/api/board/{board-id}", handler.getBoard).Methods("GET")
	r.HandleFunc("/api/board/{board-id}/item", handler.createItem).Methods("POST")
	r.HandleFunc("/api/board/{board-id}/item/{item-id}", handler.updateItem).Methods("PUT")
	r.HandleFunc("/api/board/{board-id}/updates/{version}", handler.getBoardUpdates).Methods("GET")
}
