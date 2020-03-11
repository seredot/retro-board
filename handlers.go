package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type handler struct {
	repo *Repo
}

type Handler interface {
	healthCheck(w http.ResponseWriter, r *http.Request)
	createBoard(w http.ResponseWriter, r *http.Request)
	getBoard(w http.ResponseWriter, r *http.Request)
	createItem(w http.ResponseWriter, r *http.Request)
	updateItem(w http.ResponseWriter, r *http.Request)
	getBoardUpdates(w http.ResponseWriter, r *http.Request)
}

func newHandler(r *Repo) Handler {
	return &handler{repo: r}
}

// writeError returns an error for the response.
func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

// healthCheck returns a health check message
func (h *handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(HealthCheck{Status: "ok"})
}

// createBoard creates a new board and returns the newly created board.
func (h *handler) createBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b := h.repo.CreateBoard()
	json.NewEncoder(w).Encode(b)
}

// getBoard returns a board with the specified id.
func (h *handler) getBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["board-id"]

	b, err := h.repo.GetBoard(id)
	if err != nil {
		writeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(b)
}

// createItem creates a new item for the specified board using the item info in body.
func (h *handler) createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	boardId := mux.Vars(r)["board-id"]
	var item = Item{}
	json.NewDecoder(r.Body).Decode(&item)

	retItem, err := h.repo.CreateItem(boardId, &item)
	if err != nil {
		writeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(retItem)
}

// updateItem updates an item in specified board id and item id using item info in body.
func (h *handler) updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	boardId := mux.Vars(r)["board-id"]
	itemId := mux.Vars(r)["item-id"]

	item := Item{}
	json.NewDecoder(r.Body).Decode(&item)

	retItem, err := h.repo.UpdateItem(boardId, itemId, &item)
	if err != nil {
		writeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(retItem)
}

// getBoardUpdates long polls for the changes in specified board id.
// Returns a board object with changed items.
// If no changes happen after an amount of time, returns the board object
// with an empty items array.
func (h *handler) getBoardUpdates(w http.ResponseWriter, r *http.Request) {
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
	b, err := h.repo.GetBoard(id)
	if err != nil {
		writeError(w, err)
		return
	}

	// Poll for changes on the board.
	h.repo.GetBoardUpdates(b, version)

	json.NewEncoder(w).Encode(b)
}
