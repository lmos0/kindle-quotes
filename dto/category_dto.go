package dto

import "time"

type CategoryResponse struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

type CategorySimple struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
