package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerHealthCheck(t *testing.T) {
	var repo = &RepoMock{}

	req, _ := http.NewRequest("GET", "/health-check", nil)
	h := http.HandlerFunc(NewHandler(repo).healthCheck)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	checkStatusOK(t, rr.Code)

	expected := `{"status":"ok"}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestHandlerCreateBoard(t *testing.T) {

}

func TestHandlerGetBoard(t *testing.T) {

}

func TestHandlerCreateItem(t *testing.T) {

}

func TestHandlerUpdateItem(t *testing.T) {

}

func TestHandlerGetBoardUpd(t *testing.T) {

}
