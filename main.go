package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var repo = Repo{}

func main() {
	repo.Init()

	r := mux.NewRouter()
	r.HandleFunc("/api", handleHealthCheck).Methods("GET")
	r.HandleFunc("/api/board", handleCreateBoard).Methods("POST")
	r.HandleFunc("/api/board/{board-id}", handleGetBoard).Methods("GET")
	r.HandleFunc("/api/board/{board-id}/item", handleCreateItem).Methods("POST")
	r.HandleFunc("/api/board/{board-id}/item/{item-id}", handleUpdateItem).Methods("PUT")
	r.HandleFunc("/api/board/{board-id}/updates/{version}", handleGetBoardUpdates).Methods("GET")
	log.Println("Running server")
	log.Fatal(http.ListenAndServe(":8080", r))
}
