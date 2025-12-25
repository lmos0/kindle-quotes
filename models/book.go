package models

import "time"

type Book struct {
	ID            int64
	Title         string
	ISBN          *string
	PublishedYear int
	Publisher     *string
	Pages         int
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Authors    []Author
	Categories []Category
}
