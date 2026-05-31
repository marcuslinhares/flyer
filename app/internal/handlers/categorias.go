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

func ListCategorias(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT id, nome, created_at, updated_at FROM categorias ORDER BY nome ASC`)
	if err != nil {
		log.Printf("Error listing categorias: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categorias := []models.Categoria{}
	for rows.Next() {
		var c models.Categoria
		var createdAt, updatedAt string
		err := rows.Scan(&c.ID, &c.Nome, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Error scanning categoria: %v", err)
			continue
		}
		c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		c.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		categorias = append(categorias, c)
	}

	respondJSON(w, http.StatusOK, categorias)
}

func GetCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	var c models.Categoria
	var createdAt, updatedAt string
	err = database.DB.QueryRow(
		`SELECT id, nome, created_at, updated_at FROM categorias WHERE id = ?`, id).
		Scan(&c.ID, &c.Nome, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Categoria not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error getting categoria: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	c.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	respondJSON(w, http.StatusOK, c)
}

func CreateCategoria(w http.ResponseWriter, r *http.Request) {
	var req models.CategoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Nome == "" {
		http.Error(w, `{"error":"nome is required"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`INSERT INTO categorias (nome) VALUES (?)`, req.Nome)
	if err != nil {
		log.Printf("Error creating categoria: %v", err)
		http.Error(w, `{"error":"Failed to create categoria"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	var c models.Categoria
	var createdAt, updatedAt string
	database.DB.QueryRow(
		`SELECT id, nome, created_at, updated_at FROM categorias WHERE id = ?`, id).
		Scan(&c.ID, &c.Nome, &createdAt, &updatedAt)

	c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	c.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	respondJSON(w, http.StatusCreated, c)
}

func UpdateCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	var req models.CategoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		`UPDATE categorias SET nome=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		req.Nome, id)
	if err != nil {
		log.Printf("Error updating categoria: %v", err)
		http.Error(w, `{"error":"Failed to update categoria"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"Categoria not found"}`, http.StatusNotFound)
		return
	}

	var c models.Categoria
	var createdAt, updatedAt string
	database.DB.QueryRow(
		`SELECT id, nome, created_at, updated_at FROM categorias WHERE id = ?`, id).
		Scan(&c.ID, &c.Nome, &createdAt, &updatedAt)

	c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	c.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	respondJSON(w, http.StatusOK, c)
}

func DeleteCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`DELETE FROM categorias WHERE id = ?`, id)
	if err != nil {
		log.Printf("Error deleting categoria: %v", err)
		http.Error(w, `{"error":"Failed to delete categoria"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"Categoria not found"}`, http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Categoria deleted successfully"})
}
