package dto

import "time"

type QuoteResponse struct {
	ID        int64      `json:"id"`
	Text      string     `json:"text"`
	Book      BookSimple `json:"book"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CreateQuoteRequest struct {
	BookID int64  `json:"book_id" validate:"required"`
	Text   string `json:"text" validate:"required"`
}

type UpdateQuoteRequest struct {
	Text *string `json:"text" validate:"required,min=3"`
}

type ListQuotesRequest struct {
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	BookID     int64 `json:"book_id,omitempty"`
	CategoryID int   `json:"category_id,omitempty"`
	AuthorID   int64 `json:"author_id,omitempty"`
}

type ListQuotesResponse struct {
	Quotes []QuoteResponse `json:"quotes"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
