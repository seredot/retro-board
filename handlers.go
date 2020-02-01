package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// writeError returns an error for the response.
func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

// handleHealthCheck returns a health check message
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(HealthCheck{Status: "ok"})
}

// handleCreateBoard creates a new board and returns the newly created board.
func handleCreateBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b := repo.CreateBoard()
	json.NewEncoder(w).Encode(b)
}

// handleGetBoard returns a board with the specified id.
func handleGetBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["board-id"]

	b, err := repo.GetBoard(id)
	if err != nil {
		writeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(b)
}

// handleCreateItem creates a new item for the specified board using the item info in body.
func handleCreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	boardId := mux.Vars(r)["board-id"]
	var item = Item{}
	json.NewDecoder(r.Body).Decode(&item)

	retItem, err := repo.CreateItem(boardId, &item)
	if err != nil {
		writeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(retItem)
}

// handleUpdateItem updates an item in specified board id and item id using item info in body.
func handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	boardId := mux.Vars(r)["board-id"]
	itemId := mux.Vars(r)["item-id"]

	item := Item{}
	json.NewDecoder(r.Body).Decode(&item)

	retItem, err := repo.UpdateItem(boardId, itemId, &item)
	if err != nil {
		writeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(retItem)
}

// handleGetBoardUpdates long polls for the changes in specified board id.
// Returns a board object with changed items.
// If no changes happen after an amount of time, returns the board object
// with an empty items array.
func handleGetBoardUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["board-id"]
	sVersion := mux.Vars(r)["version"]

	// Parse version number.
	version, err := strconv.ParseUint(sVersion, 10, 64)
	if err != nil {
		writeError(w, errors.New("invalid_argument_version"))
		return
	}

	// Does the board exist?
	b, err := repo.GetBoard(id)
	if err != nil {
		writeError(w, err)
		return
	}

	// Poll for changes on the board.
	repo.GetBoardUpdates(b, version)

	json.NewEncoder(w).Encode(b)
}
