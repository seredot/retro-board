package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	mock "github.com/stretchr/testify/mock"
)

func TestHandlerHealthCheck(t *testing.T) {
	var repo = &RepoMock{}

	expected := &HealthCheck{Status: "ok"}

	req, _ := http.NewRequest("GET", "/api", nil)
	h := http.HandlerFunc(NewHandler(repo).healthCheck)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusOK(t, rr.Code)
	checkResultJSON(t, expected, rr.Body.Bytes(), &HealthCheck{})
	repo.AssertExpectations(t)
}

func TestHandlerCreateBoard(t *testing.T) {
	var repo = &RepoMock{}

	expected := &Board{
		Id:      "board_id",
		Items:   make(map[string]*Item),
		Version: 0,
	}

	repo.
		On("CreateBoard").
		Return(&Board{
			Id:      "board_id",
			Items:   make(map[string]*Item),
			Version: 0,
		}, nil).Once()

	req, _ := http.NewRequest("POST", "/api/board", nil)
	h := http.HandlerFunc(NewHandler(repo).createBoard)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusOK(t, rr.Code)
	checkResultJSON(t, expected, rr.Body.Bytes(), &Board{})
	repo.AssertExpectations(t)
}

func TestHandlerGetBoard(t *testing.T) {
	var repo = &RepoMock{}

	expected := &Board{
		Id:      "board_id",
		Items:   make(map[string]*Item),
		Version: 0,
	}

	repo.
		On("GetBoard", "board_id").
		Return(&Board{
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
	checkResultJSON(t, expected, rr.Body.Bytes(), &Board{})
	repo.AssertExpectations(t)
}

func TestHandlerGetBoardError(t *testing.T) {
	var repo = &RepoMock{}
	var nilBoard *Board

	expected := &ErrorResponse{
		Error: "board_not_found",
	}

	repo.
		On("GetBoard", "not_existing_board_id").
		Return(nilBoard, errors.New("board_not_found")).Once()

	req, _ := http.NewRequest("GET", "/api/board/board_id", nil)
	req = mux.SetURLVars(req, map[string]string{
		"board-id": "not_existing_board_id",
	})
	h := http.HandlerFunc(NewHandler(repo).getBoard)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusNotOK(t, rr.Code)
	checkResultJSON(t, expected, rr.Body.Bytes(), &ErrorResponse{})
	repo.AssertExpectations(t)
}

func TestHandlerCreateItem(t *testing.T) {
	var repo = &RepoMock{}

	expected := &Item{
		Id:     "item_id",
		Text:   "This is an item",
		Color:  "blue",
		Left:   0,
		Top:    0,
		Width:  0,
		Height: 0,
	}

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
	checkResultJSON(t, expected, rr.Body.Bytes(), &Item{})
	repo.AssertExpectations(t)
}

func TestHandlerCreateItemInputError(t *testing.T) {
	var repo = &RepoMock{}

	expected := &ErrorResponse{
		Error: "Missing input error",
	}

	req, _ := http.NewRequest("POST", "/api/board/board_id/item", nil)
	h := http.HandlerFunc(NewHandler(repo).createItem)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusNotOK(t, rr.Code)
	checkResultJSON(t, expected, rr.Body.Bytes(), &ErrorResponse{})
	repo.AssertExpectations(t)
}

func TestHandlerCreateItemParseError(t *testing.T) {
	var repo = &RepoMock{}

	expected := &ErrorResponse{
		Error: "Parse error",
	}

	req, _ := http.NewRequest("POST", "/api/board/board_id/item", strings.NewReader("invalid: json"))
	h := http.HandlerFunc(NewHandler(repo).createItem)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusNotOK(t, rr.Code)
	checkResultJSON(t, expected, rr.Body.Bytes(), &ErrorResponse{})
	repo.AssertExpectations(t)
}

func TestHandlerCreateItemNotFoundBoardError(t *testing.T) {
	var repo = &RepoMock{}

	expected := &ErrorResponse{
		Error: "board_not_found",
	}

	repo.
		On("CreateItem", mock.Anything, mock.Anything).
		Return(&Item{}, errors.New("board_not_found")).
		Once()

	req, _ := http.NewRequest(
		"POST",
		"/api/board/not_found_board_id/item",
		strings.NewReader("{}"),
	)
	h := http.HandlerFunc(NewHandler(repo).createItem)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusNotOK(t, rr.Code)
	checkResultJSON(t, expected, rr.Body.Bytes(), &ErrorResponse{})
	repo.AssertExpectations(t)
}

func TestHandlerUpdateItem(t *testing.T) {
	var repo = &RepoMock{}

	expected := &Item{
		Id:     "item_id",
		Text:   "This is an updated item",
		Color:  "green",
		Left:   10,
		Top:    10,
		Width:  100,
		Height: 100,
	}

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
	checkResultJSON(t, expected, rr.Body.Bytes(), &Item{})
	repo.AssertExpectations(t)
}

func TestHandlerUpdateItemNotFoundBoardError(t *testing.T) {
	var repo = &RepoMock{}

	expected := &ErrorResponse{
		Error: "board_not_found",
	}

	repo.
		On("UpdateItem", mock.Anything, mock.Anything, mock.Anything).
		Return(&Item{}, errors.New("board_not_found")).
		Once()

	req, _ := http.NewRequest(
		"POST",
		"/api/board/not_found_board_id/item",
		strings.NewReader("{}"),
	)
	h := http.HandlerFunc(NewHandler(repo).updateItem)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusNotOK(t, rr.Code)
	checkResultJSON(t, expected, rr.Body.Bytes(), &ErrorResponse{})
	repo.AssertExpectations(t)
}

func TestHandlerGetBoardUpdates(t *testing.T) {
	var repo = &RepoMock{}

	expected := &Board{
		Id:      "board_id",
		Items:   make(map[string]*Item),
		Version: 0,
	}

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
	checkResultJSON(t, expected, rr.Body.Bytes(), &Board{})
	repo.AssertExpectations(t)
}

func TestHandlerGetBoardUpdatesInputError(t *testing.T) {
	var repo = &RepoMock{}

	expected := &ErrorResponse{
		Error: "invalid_argument_version",
	}

	req, _ := http.NewRequest("GET", "/api/board/board_id/updates/K0", nil)
	req = mux.SetURLVars(req, map[string]string{
		"version": "K0",
	})
	h := http.HandlerFunc(NewHandler(repo).getBoardUpdates)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusNotOK(t, rr.Code)
	checkResultJSON(t, expected, rr.Body.Bytes(), &ErrorResponse{})
	repo.AssertExpectations(t)
}

func TestHandlerGetBoardUpdatesNotFoundBoardError(t *testing.T) {
	var repo = &RepoMock{}
	var nilBoard *Board

	expected := &ErrorResponse{
		Error: "board_not_found",
	}

	repo.
		On("GetBoard", "not_existing_board_id").
		Return(nilBoard, errors.New("board_not_found")).Once()

	req, _ := http.NewRequest("GET", "/api/board/not_existing_board_id/updates/0", nil)
	req = mux.SetURLVars(req, map[string]string{
		"board-id": "not_existing_board_id",
		"version":  "0",
	})
	h := http.HandlerFunc(NewHandler(repo).getBoardUpdates)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusNotOK(t, rr.Code)
	checkResultJSON(t, expected, rr.Body.Bytes(), &ErrorResponse{})
	repo.AssertExpectations(t)
}
