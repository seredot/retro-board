package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkStatusOK(t *testing.T, code int) {
	if code != http.StatusOK {
		t.Errorf("handler returned status code: %d", code)
	}
}

func checkStatusNotOK(t *testing.T, code int) {
	if code < 300 {
		t.Errorf("handler returned status code: %d", code)
	}
}

func checkResultJSON(t *testing.T, expected interface{}, bytes []byte, v interface{}) {
	err := json.Unmarshal(bytes, v)

	assert.NoError(t, err, "Error parsing JSON result.")
	assert.Equal(t, expected, v)
}
