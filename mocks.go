package main

import mock "github.com/stretchr/testify/mock"

type RepoMock struct {
	mock.Mock
}

// CreateBoard provides a mock function with given fields:
func (_m *RepoMock) CreateBoard() *Board {
	ret := _m.Called()

	return ret.Get(0).(*Board)
}

// CreateItem provides a mock function with given fields: boardId, item
func (_m *RepoMock) CreateItem(boardId string, item *Item) (*Item, error) {
	ret := _m.Called(boardId, item)

	return ret.Get(0).(*Item), ret.Error(1)
}

// GetBoard provides a mock function with given fields: id
func (_m *RepoMock) GetBoard(id string) (*Board, error) {
	ret := _m.Called(id)

	return ret.Get(0).(*Board), ret.Error(1)
}

// GetBoardUpdates provides a mock function with given fields: b, version
func (_m *RepoMock) GetBoardUpdates(b *Board, version uint64) {
	_m.Called(b, version)
}

// GetItem provides a mock function with given fields: b, itemId
func (_m *RepoMock) GetItem(b *Board, itemId string) (*Item, error) {
	ret := _m.Called(b, itemId)

	return ret.Get(0).(*Item), ret.Error(1)
}

// UpdateBoard provides a mock function with given fields: b, it
func (_m *RepoMock) UpdateBoard(b *Board, it *Item) {
	_m.Called(b, it)
}

// UpdateItem provides a mock function with given fields: boardId, itemId, item
func (_m *RepoMock) UpdateItem(boardId string, itemId string, item *Item) (*Item, error) {
	ret := _m.Called(boardId, itemId, item)

	return ret.Get(0).(*Item), ret.Error(1)
}
