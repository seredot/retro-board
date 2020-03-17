package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestHandlerHealthCheck(t *testing.T) {
	var repo = &RepoMock{}

	req, _ := http.NewRequest("GET", "/api", nil)
	h := http.HandlerFunc(NewHandler(repo).healthCheck)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusOK(t, rr.Code)
	assert.Equal(t, "{\"status\":\"ok\"}\n", rr.Body.String())
	repo.AssertExpectations(t)
}

func TestHandlerCreateBoard(t *testing.T) {
	var repo = &RepoMock{}

	repo.On("CreateBoard").Return(&Board{
		Id:      "board_id",
		Items:   make(map[string]*Item),
		Version: 0,
	}, nil).Once()

	req, _ := http.NewRequest("POST", "/api/board", nil)
	h := http.HandlerFunc(NewHandler(repo).createBoard)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusOK(t, rr.Code)
	assert.Equal(t, "{\"id\":\"board_id\",\"items\":{},\"version\":0}\n", rr.Body.String())
	repo.AssertExpectations(t)
}

func TestHandlerGetBoard(t *testing.T) {
	var repo = &RepoMock{}

	repo.On("GetBoard", "board_id").Return(&Board{
		Id:      "board_id",
		Items:   make(map[string]*Item),
		Version: 0,
	}, nil).Once()

	req, _ := http.NewRequest("GET", "/api/board/board_id", nil)
	req = mux.SetURLVars(req, map[string]string{
		"board-id": "board_id",
	})
	h := http.HandlerFunc(NewHandler(repo).getBoard)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusOK(t, rr.Code)
	assert.Equal(t, "{\"id\":\"board_id\",\"items\":{},\"version\":0}\n", rr.Body.String())
	repo.AssertExpectations(t)
}

func TestHandlerCreateItem(t *testing.T) {
	var repo = &RepoMock{}

	repo.On("CreateItem", "board_id", mock.Anything).Return(&Item{
		Id:     "item_id",
		Text:   "This is an item",
		Color:  "blue",
		Left:   0,
		Top:    0,
		Width:  0,
		Height: 0,
	}, nil).Once()

	input := `{
    "text": "This is an item",
    "color": "blue"
	}`

	req, _ := http.NewRequest("POST", "/api/board/board_id/item", strings.NewReader(input))
	req = mux.SetURLVars(req, map[string]string{
		"board-id": "board_id",
	})
	h := http.HandlerFunc(NewHandler(repo).createItem)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusOK(t, rr.Code)
	assert.Equal(t, "{\"id\":\"item_id\",\"text\":\"This is an item\",\"color\":\"blue\",\"left\":0,\"top\":0,\"width\":0,\"height\":0}\n", rr.Body.String())
	repo.AssertExpectations(t)
}

func TestHandlerUpdateItem(t *testing.T) {
	var repo = &RepoMock{}

	repo.On("UpdateItem", "board_id", "item_id", mock.Anything).Return(&Item{
		Id:     "item_id",
		Text:   "This is an updated item",
		Color:  "green",
		Left:   10,
		Top:    10,
		Width:  100,
		Height: 100,
	}, nil).Once()

	input := `{
    "id": "item_id",
    "text": "This is an updated item",
    "color": "green",
    "left": 10,
    "top": 10,
    "width": 100,
    "height": 100
  }`

	req, _ := http.NewRequest("PUT", "/api/board/board_id/item/item_id", strings.NewReader(input))
	req = mux.SetURLVars(req, map[string]string{
		"board-id": "board_id",
		"item-id":  "item_id",
	})
	h := http.HandlerFunc(NewHandler(repo).updateItem)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusOK(t, rr.Code)
	assert.Equal(t, "{\"id\":\"item_id\",\"text\":\"This is an updated item\",\"color\":\"green\",\"left\":10,\"top\":10,\"width\":100,\"height\":100}\n", rr.Body.String())
	repo.AssertExpectations(t)
}

func TestHandlerGetBoardUpdates(t *testing.T) {
	var repo = &RepoMock{}

	repo.On("GetBoard", "board_id").Return(&Board{
		Id:      "board_id",
		Items:   make(map[string]*Item),
		Version: 0,
	}, nil).Once()

	repo.On("GetBoardUpdates", mock.Anything, mock.Anything).Return().Once()

	req, _ := http.NewRequest("GET", "/api/board/board_id/updates/0", nil)
	req = mux.SetURLVars(req, map[string]string{
		"board-id": "board_id",
		"version":  "0",
	})
	h := http.HandlerFunc(NewHandler(repo).getBoardUpdates)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusOK(t, rr.Code)
	assert.Equal(t, "{\"id\":\"board_id\",\"items\":{},\"version\":0}\n", rr.Body.String())
	repo.AssertExpectations(t)
}