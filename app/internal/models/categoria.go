package models

import "time"

type Categoria struct {
	ID        int64     `json:"id"`
	Nome      string    `json:"nome"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoriaRequest struct {
	Nome string `json:"nome"`
}
