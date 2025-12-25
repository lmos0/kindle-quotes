package models

import "time"

type Author struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
