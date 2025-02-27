package entity

import (
	"time"
)

type Decision struct {
	ActorID     string
	RecipientID string
	Liked       bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Liker struct {
	ActorID       string
	UnixTimestamp uint64
}
