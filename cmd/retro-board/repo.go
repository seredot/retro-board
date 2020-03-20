package main

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

// Repo interface.
type Repo interface {
	CreateBoard() *Board
	GetBoard(id string) (*Board, error)
	UpdateBoard(b *Board, it *Item)
	GetBoardUpdates(b *Board, version uint64)
	CreateItem(boardId string, item *Item) (*Item, error)
	GetItem(b *Board, itemId string) (*Item, error)
	UpdateItem(boardId string, itemId string, item *Item) (*Item, error)
}

// memoryRepo is an in-memory data store.
type memoryRepo struct {
	boards map[string]*Board
}

// NewMemoryRepo initializes the repo.
func NewMemoryRepo() Repo {
	r := memoryRepo{}
	r.boards = make(map[string]*Board)

	return &r
}

// CreateBoard creates a new board.
func (r *memoryRepo) CreateBoard() *Board {
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
func (r *memoryRepo) GetBoard(id string) (*Board, error) {
	b := r.boards[id]
	if b == nil {
		return nil, errors.New("board_not_found")
	}
	return b, nil
}

// UpdateBoard updates the board version and broadcasts the update to listeners.
func (r *memoryRepo) UpdateBoard(b *Board, it *Item) {
	// Increment the board version.
	v := atomic.AddUint64(&b.Version, 1)

	// Update item version.
	if it != nil {
		atomic.StoreUint64(&it.Version, v)
	}

	// Broadcast listeners.
	b.Cond.Broadcast()
}

func (r *memoryRepo) GetBoardUpdates(b *Board, version uint64) {
	b.Cond.L.Lock()
	b.Cond.Wait()
	b.Cond.L.Unlock()
}

// CreateItem creates a new item.
func (r *memoryRepo) CreateItem(boardId string, item *Item) (*Item, error) {
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
func (r *memoryRepo) GetItem(b *Board, itemId string) (*Item, error) {
	item, ok := b.Items[itemId]
	if !ok {
		return nil, errors.New("item_not_found")
	}
	return item, nil
}

// UpdateItem updates an existing item.
func (r *memoryRepo) UpdateItem(boardId string, itemId string, item *Item) (*Item, error) {
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

	if item == nil {
		return nil, errors.New("input_error")
	}

	// Copy data from received item.
	*oItem = *item
	// No highjacking.
	oItem.Id = itemId

	// Notify listeners.
	r.UpdateBoard(b, oItem)

	return oItem, nil
}
