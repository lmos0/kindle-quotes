package dto

import "time"

type AuthorResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Order     *int      `json:"order"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAuthorRequest struct {
	Name string `json:"name"`
}

type UpdateAuthorRequest struct {
	Name string `json:"name"`
}

type AuthorSimple struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
