package main

import "sync"

// ErrorResponse used for service responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthCheck
type HealthCheck struct {
	Status string `json:"status"`
}

// BoardSync is the board synchronization data.
type BoardSync struct {
	Mutex sync.Mutex
	Cond  *sync.Cond
}

// Board data.
type Board struct {
	BoardSync `json:"-"`
	Id        string           `json:"id"`
	Items     map[string]*Item `json:"items"`
	Version   uint64           `json:"version"`
}

// Item of a board.
type Item struct {
	Version uint64  `json:"-"`
	Id      string  `json:"id"`
	Text    string  `json:"text"`
	Color   string  `json:"color"`
	Left    float32 `json:"left"`
	Top     float32 `json:"top"`
	Width   float32 `json:"width"`
	Height  float32 `json:"height"`
}
