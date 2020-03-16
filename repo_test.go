package main

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepoCreateBoard(t *testing.T) {
	r := NewMemoryRepo()
	b := r.CreateBoard()

	assert.Equal(t, len(b.Id), 36)
	assert.Zero(t, b.Version)
	assert.NotNil(t, b.Items)
}

func TestRepoGetBoard(t *testing.T) {
	r := NewMemoryRepo()
	b := r.CreateBoard()
	result, err := r.GetBoard(b.Id)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, b.Id, result.Id)

	result, err = r.GetBoard("not_existing_board_id")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRepoUpdateBoard(t *testing.T) {
	r := NewMemoryRepo()
	b := r.CreateBoard()
	i := Item{}
	r.UpdateBoard(b, &i)

	assert.EqualValues(t, 1, b.Version)
	assert.Equal(t, b.Version, i.Version)

	r.UpdateBoard(b, nil)
	assert.EqualValues(t, 2, b.Version)
}

func TestRepoGetBoardUpdates(t *testing.T) {
	r := NewMemoryRepo()
	b := r.CreateBoard()
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		// Waiting for updates...
		r.GetBoardUpdates(b, 0)
		// Updates received.
		wg.Done()
	}()

	<-time.After(time.Millisecond * 10)
	r.UpdateBoard(b, nil)

	wg.Wait()
}

func TestRepoCreateItem(t *testing.T) {
	item := &Item{
		Id:      "foo",
		Version: 0,
		Color:   "red",
		Text:    "bar",
		Left:    0,
		Top:     1,
		Width:   2,
		Height:  3,
	}

	r := NewMemoryRepo()
	b := r.CreateBoard()
	created, err := r.CreateItem(b.Id, item)

	assert.NoError(t, err)
	assert.Equal(t, len(created.Id), 36)
	assert.EqualValues(t, 1, created.Version)
	assert.EqualValues(t, item.Color, created.Color)
	assert.EqualValues(t, item.Text, created.Text)
	assert.EqualValues(t, item.Left, created.Left)
	assert.EqualValues(t, item.Top, created.Top)
	assert.EqualValues(t, item.Width, created.Width)
	assert.EqualValues(t, item.Height, created.Height)

	created, err = r.CreateItem("not_existing_board_id", item)

	assert.Nil(t, created)
	assert.Error(t, err)
}

func TestRepoGetItem(t *testing.T) {
	item := &Item{
		Id:      "foo",
		Version: 0,
		Color:   "red",
		Text:    "bar",
		Left:    0,
		Top:     1,
		Width:   2,
		Height:  3,
	}

	r := NewMemoryRepo()
	b := r.CreateBoard()
	created, _ := r.CreateItem(b.Id, item)
	result, err := r.GetItem(b, created.Id)

	assert.EqualValues(t, result, created)
	assert.NoError(t, err)

	result, err = r.GetItem(b, "not_existing_item_id")

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestRepoUpdateItem(t *testing.T) {
	createInput := &Item{
		Id:      "",
		Version: 0,
		Color:   "red",
		Text:    "foo",
		Left:    0,
		Top:     1,
		Width:   2,
		Height:  3,
	}

	updateInput := &Item{
		Id:      "",
		Version: 1,
		Color:   "green",
		Text:    "bar",
		Left:    4,
		Top:     5,
		Width:   6,
		Height:  7,
	}

	r := NewMemoryRepo()
	b := r.CreateBoard()
	created, _ := r.CreateItem(b.Id, createInput)
	updated, err := r.UpdateItem(b.Id, created.Id, updateInput)

	assert.NoError(t, err)
	assert.EqualValues(t, updateInput.Color, updated.Color)
	assert.EqualValues(t, updateInput.Text, updated.Text)
	assert.EqualValues(t, updateInput.Left, updated.Left)
	assert.EqualValues(t, updateInput.Top, updated.Top)
	assert.EqualValues(t, updateInput.Width, updated.Width)
	assert.EqualValues(t, updateInput.Height, updated.Height)

	result, err := r.GetItem(b, updated.Id)
	assert.NoError(t, err)
	assert.EqualValues(t, result, updated)

	notFound, err := r.UpdateItem("not_existing_board_id", updated.Id, updateInput)
	assert.Error(t, err)
	assert.Nil(t, notFound)

	notFound, err = r.UpdateItem(b.Id, "not_existing_item_id", updateInput)
	assert.Error(t, err)
	assert.Nil(t, notFound)

	notFound, err = r.UpdateItem(b.Id, updated.Id, nil)
	assert.Error(t, err)
	assert.Nil(t, notFound)
}
