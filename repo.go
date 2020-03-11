package main

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

// Repo is an in-memory data store.
type Repo struct {
	boards map[string]*Board
}

// Init initializes the repo.
func (r *Repo) Init() {
	r.boards = make(map[string]*Board)
}

// CreateBoard creates a new board.
func (r *Repo) CreateBoard() *Board {
	id := uuid.New().String()
	b := &Board{
		Id:      id,
		Items:   make(map[string]*Item),
		Version: 0,
	}
	b.Cond = sync.NewCond(&b.Mutex)
	r.boards[id] = b
	return b
}

// GetBoard gets a board.
func (r *Repo) GetBoard(id string) (*Board, error) {
	b := r.boards[id]
	if b == nil {
		return nil, errors.New("board_not_found")
	}
	return b, nil
}

// UpdateBoard updates the board version and broadcasts the update to listeners.
func (r *Repo) UpdateBoard(b *Board, it *Item) {
	// Update board version.
	v := atomic.AddUint64(&b.Version, 1)

	// Update item version.
	if it != nil {
		atomic.StoreUint64(&it.Version, v)
	}

	// Broadcast listeners.
	b.Cond.Broadcast()
}

func (r *Repo) GetBoardUpdates(b *Board, version uint64) {
	b.Cond.L.Lock()
	b.Cond.Wait()
	b.Cond.L.Unlock()
}

// CreateItem creates a new item.
func (r *Repo) CreateItem(boardId string, item *Item) (*Item, error) {
	b, err := r.GetBoard(boardId)
	if err != nil {
		return nil, err
	}
	retItem := *item
	retItem.Id = uuid.New().String()
	b.Items[retItem.Id] = &retItem

	// Notify listeners.
	r.UpdateBoard(b, &retItem)

	return &retItem, nil
}

// GetItem gets an item.
func (r *Repo) GetItem(b *Board, itemId string) (*Item, error) {
	item, ok := b.Items[itemId]
	if !ok {
		return nil, errors.New("item_not_found")
	}
	return item, nil
}

// UpdateItem updates an existing item.
func (r *Repo) UpdateItem(boardId string, itemId string, item *Item) (*Item, error) {
	// Find the board
	b, err := r.GetBoard(boardId)
	if err != nil {
		return nil, err
	}

	// Get the existing item.
	oItem, err := r.GetItem(b, itemId)
	if err != nil {
		return nil, err
	}

	// Copy data from received item.
	*oItem = *item
	// No highjacking.
	oItem.Id = itemId

	// Notify listeners.
	r.UpdateBoard(b, oItem)

	return oItem, nil
}
