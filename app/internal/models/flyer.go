package models

import "time"

type Flyer struct {
	ID              int64     `json:"id"`
	Titulo          string    `json:"titulo"`
	Descricao       string    `json:"descricao"`
	Preco           float64   `json:"preco"`
	DataValida      string    `json:"data_valida"`
	Imagem          string    `json:"imagem"`
	EstabelecimentoID *int64  `json:"estabelecimento_id,omitempty"`
	CategoriaID     *int64    `json:"categoria_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type FlyerRequest struct {
	Titulo          string  `json:"titulo"`
	Descricao       string  `json:"descricao"`
	Preco           float64 `json:"preco"`
	DataValida      string  `json:"data_valida"`
	Imagem          string  `json:"imagem"`
	EstabelecimentoID *int64 `json:"estabelecimento_id,omitempty"`
	CategoriaID     *int64  `json:"categoria_id,omitempty"`
}
