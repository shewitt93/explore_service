package entity

import "github.com/stretchr/testify/mock"

type MockCursorService struct {
	mock.Mock
}

func (m *MockCursorService) Encode(c *Cursor) (string, error) {
	args := m.Called(c)
	return args.String(0), args.Error(1)
}

func (m *MockCursorService) Decode(cursor string) (*Cursor, error) {
	args := m.Called(cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Cursor), args.Error(1)
}

func (m *MockCursorService) CreateNextCursor(likers []Liker, limit int) ([]Liker, *Cursor) {
	args := m.Called(likers, limit)
	return args.Get(0).([]Liker), args.Get(1).(*Cursor)
}
