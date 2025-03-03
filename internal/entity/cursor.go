package entity

type CursorService interface {
	Encode(c *Cursor) (string, error)
	Decode(cursor string) (*Cursor, error)
	CreateNextCursor(likers []Liker, limit int) ([]Liker, *Cursor)
}

type DefaultCursorService struct{}

func NewCursorService() CursorService {
	return &DefaultCursorService{}
}

func (s *DefaultCursorService) Encode(c *Cursor) (string, error) {
	return EncodeCursor(c)
}

func (s *DefaultCursorService) Decode(cursor string) (*Cursor, error) {
	return DecodeCursor(cursor)
}

func (s *DefaultCursorService) CreateNextCursor(likers []Liker, limit int) ([]Liker, *Cursor) {
	return CreateNextCursor(likers, limit)
}
