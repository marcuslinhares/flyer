package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/marcuslinhares/flyer/app/internal/database"
	"github.com/marcuslinhares/flyer/app/internal/models"
)

func ListEstabelecimentos(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT id, nome, endereco, telefone, logo, created_at, updated_at
		 FROM estabelecimentos ORDER BY nome ASC`)
	if err != nil {
		log.Printf("Error listing estabelecimentos: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	estabs := []models.Estabelecimento{}
	for rows.Next() {
		var e models.Estabelecimento
		var createdAt, updatedAt string
		err := rows.Scan(&e.ID, &e.Nome, &e.Endereco, &e.Telefone, &e.Logo, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Error scanning estabelecimento: %v", err)
			continue
		}
		e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		estabs = append(estabs, e)
	}

	respondJSON(w, http.StatusOK, estabs)
}

func GetEstabelecimento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	var e models.Estabelecimento
	var createdAt, updatedAt string
	err = database.DB.QueryRow(
		`SELECT id, nome, endereco, telefone, logo, created_at, updated_at
		 FROM estabelecimentos WHERE id = ?`, id).
		Scan(&e.ID, &e.Nome, &e.Endereco, &e.Telefone, &e.Logo, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Estabelecimento not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error getting estabelecimento: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	respondJSON(w, http.StatusOK, e)
}

func CreateEstabelecimento(w http.ResponseWriter, r *http.Request) {
	var req models.EstabelecimentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Nome == "" {
		http.Error(w, `{"error":"nome is required"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		`INSERT INTO estabelecimentos (nome, endereco, telefone, logo) VALUES (?, ?, ?, ?)`,
		req.Nome, req.Endereco, req.Telefone, req.Logo)
	if err != nil {
		log.Printf("Error creating estabelecimento: %v", err)
		http.Error(w, `{"error":"Failed to create estabelecimento"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	var e models.Estabelecimento
	var createdAt, updatedAt string
	database.DB.QueryRow(
		`SELECT id, nome, endereco, telefone, logo, created_at, updated_at
		 FROM estabelecimentos WHERE id = ?`, id).
		Scan(&e.ID, &e.Nome, &e.Endereco, &e.Telefone, &e.Logo, &createdAt, &updatedAt)

	e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	respondJSON(w, http.StatusCreated, e)
}

func UpdateEstabelecimento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	var req models.EstabelecimentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		`UPDATE estabelecimentos SET nome=?, endereco=?, telefone=?, logo=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		req.Nome, req.Endereco, req.Telefone, req.Logo, id)
	if err != nil {
		log.Printf("Error updating estabelecimento: %v", err)
		http.Error(w, `{"error":"Failed to update estabelecimento"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"Estabelecimento not found"}`, http.StatusNotFound)
		return
	}

	var e models.Estabelecimento
	var createdAt, updatedAt string
	database.DB.QueryRow(
		`SELECT id, nome, endereco, telefone, logo, created_at, updated_at
		 FROM estabelecimentos WHERE id = ?`, id).
		Scan(&e.ID, &e.Nome, &e.Endereco, &e.Telefone, &e.Logo, &createdAt, &updatedAt)

	e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	respondJSON(w, http.StatusOK, e)
}

func DeleteEstabelecimento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`DELETE FROM estabelecimentos WHERE id = ?`, id)
	if err != nil {
		log.Printf("Error deleting estabelecimento: %v", err)
		http.Error(w, `{"error":"Failed to delete estabelecimento"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"Estabelecimento not found"}`, http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Estabelecimento deleted successfully"})
}
