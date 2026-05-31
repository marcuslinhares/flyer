package models

import "time"

type Estabelecimento struct {
	ID        int64     `json:"id"`
	Nome      string    `json:"nome"`
	Endereco  string    `json:"endereco"`
	Telefone  string    `json:"telefone"`
	Logo      string    `json:"logo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EstabelecimentoRequest struct {
	Nome     string `json:"nome"`
	Endereco string `json:"endereco"`
	Telefone string `json:"telefone"`
	Logo     string `json:"logo"`
}
