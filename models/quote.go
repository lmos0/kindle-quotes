package models

import "time"

type Quote struct {
	ID        int64
	BookID    int64
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time

	Book *Book
}
