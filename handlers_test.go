package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := setupRouter()
	rr := callHandler(router, req)

	var result HealthCheck
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("unable to parse response: %s", err)
	}

	checkStatusOK(t, rr.Code)
	assert.Equal(t, result.Status, "ok")
}

func TestCreateBoard(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/board", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := setupRouter()
	rr := callHandler(router, req)

	var result Board
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("unable to parse response: %s", err)
	}

	checkStatusOK(t, rr.Code)
	// Check if an `Id` is returned.
	assert.Greater(t, len(result.Id), 0)
}

func TestGetBoard(t *testing.T) {
	// First, create a board to reach an initial state.
	req, err := http.NewRequest("POST", "/api/board", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := setupRouter()
	rr := callHandler(router, req)

	var created Board
	err = json.Unmarshal(rr.Body.Bytes(), &created)
	if err != nil {
		t.Errorf("unable to parse response: %s", err)
	}

	checkStatusOK(t, rr.Code)
	boardId := created.Id

	// Get the board using returned identifier.
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/board/%s", boardId), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = callHandler(router, req)

	var received Board
	err = json.Unmarshal(rr.Body.Bytes(), &received)
	if err != nil {
		t.Errorf("unable to parse response: %s", err)
	}

	checkStatusOK(t, rr.Code)
	assert.Equal(t, boardId, received.Id)
}

func TestCreateItem(t *testing.T) {
	// First, create a board.
	req, err := http.NewRequest("POST", "/api/board", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := setupRouter()
	rr := callHandler(router, req)

	var created Board
	err = json.Unmarshal(rr.Body.Bytes(), &created)
	if err != nil {
		t.Errorf("unable to parse response: %s", err)
	}

	checkStatusOK(t, rr.Code)
	boardId := created.Id

	// Create an item on the board.
	body := strings.NewReader(`{"text": "foo"}`)
	req, err = http.NewRequest("POST", fmt.Sprintf("/api/board/%s/item", boardId), body)
	if err != nil {
		t.Fatal(err)
	}

	rr = callHandler(router, req)

	var item Item
	err = json.Unmarshal(rr.Body.Bytes(), &item)
	if err != nil {
		t.Errorf("unable to parse response: %s", err)
	}

	checkStatusOK(t, rr.Code)
	assert.Greater(t, len(item.Id), 0)
}

// Helpers
////////////

func callHandler(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func setupRouter() *mux.Router {
	repo := &Repo{}
	setInitialState(repo)

	r := mux.NewRouter()
	mapHandlerFuncs(r, newHandler(repo))

	return r
}

func setInitialState(repo *Repo) {
	repo.Init()
}

func checkStatusOK(t *testing.T, code int) {
	if code != http.StatusOK {
		t.Errorf("handler returned status code: %d", code)
	}
}
