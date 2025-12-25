package dto

import "time"

type BookResponse struct {
	ID            int64            `json:"id"`
	Title         string           `json:"title"`
	ISBN          *string          `json:"isbn,omitempty"`
	PublishedYear int              `json:"published_year"`
	Publisher     *string          `json:"publisher,omitempty"`
	Pages         *int             `json:"pages,omitempty"`
	Authors       []AuthorResponse `json:"authors"`
	Categories    []CategorySimple `json:"categories"`
	CreatedAt     time.Time        `json:"created_at"`
}

type CreateBookRequest struct {
	Title         string  `json:"title"`
	ISBN          *string `json:"isbn,omitempty"`
	PublishedYear int     `json:"published_year" validate:"required, min=1000, max=2100"`
	Publisher     *string `json:"publisher,omitempty"`
	Pages         *int    `json:"pages,omitempty"`
	AuthorIDs     []int   `json:"author_ids" validate:"required,min=1,"`
	CategoriesIDs []int   `json:"categories_ids,omitempty" validate:"required,min=1,max=2"`
}

type UpdateBookRequest struct {
	Title         string  `json:"title"`
	ISBN          *string `json:"isbn,omitempty"`
	PublishedYear int     `json:"published_year" validate:"required, min=1000, max=2100"`
	Publisher     *string `json:"publisher,omitempty"`
	Pages         *int    `json:"pages,omitempty"`
	AuthorIDs     []int   `json:"author_ids" validate:"required,min=1,max=2"`
	CategoriesIDs []int   `json:"categories_ids" validate:"required,min=1,max=2"`
}

type BookSimple struct {
	ID            int64   `json:"id"`
	Title         string  `json:"title"`
	PublishedYear int     `json:"published_year"`
	Isbn          *string `json:"isbn,omitempty"`
}
