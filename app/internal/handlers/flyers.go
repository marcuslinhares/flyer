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

func ListFlyers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT id, titulo, descricao, preco, data_valida, imagem,
		        COALESCE(estabelecimento_id, 0), COALESCE(categoria_id, 0),
		        created_at, updated_at FROM flyers ORDER BY created_at DESC`)
	if err != nil {
		log.Printf("Error listing flyers: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	flyers := []models.Flyer{}
	for rows.Next() {
		var f models.Flyer
		var estID, catID int64
		var createdAt, updatedAt string
		err := rows.Scan(&f.ID, &f.Titulo, &f.Descricao, &f.Preco, &f.DataValida, &f.Imagem,
			&estID, &catID, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Error scanning flyer: %v", err)
			continue
		}
		if estID > 0 {
			f.EstabelecimentoID = &estID
		}
		if catID > 0 {
			f.CategoriaID = &catID
		}
		f.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		f.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		flyers = append(flyers, f)
	}

	respondJSON(w, http.StatusOK, flyers)
}

func GetFlyer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	var f models.Flyer
	var estID, catID int64
	var createdAt, updatedAt string
	err = database.DB.QueryRow(
		`SELECT id, titulo, descricao, preco, data_valida, imagem,
		        COALESCE(estabelecimento_id, 0), COALESCE(categoria_id, 0),
		        created_at, updated_at FROM flyers WHERE id = ?`, id).
		Scan(&f.ID, &f.Titulo, &f.Descricao, &f.Preco, &f.DataValida, &f.Imagem,
			&estID, &catID, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Flyer not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error getting flyer: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	if estID > 0 {
		f.EstabelecimentoID = &estID
	}
	if catID > 0 {
		f.CategoriaID = &catID
	}
	f.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	f.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	respondJSON(w, http.StatusOK, f)
}

func CreateFlyer(w http.ResponseWriter, r *http.Request) {
	var req models.FlyerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Titulo == "" {
		http.Error(w, `{"error":"titulo is required"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		`INSERT INTO flyers (titulo, descricao, preco, data_valida, imagem, estabelecimento_id, categoria_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.Titulo, req.Descricao, req.Preco, req.DataValida, req.Imagem,
		req.EstabelecimentoID, req.CategoriaID)
	if err != nil {
		log.Printf("Error creating flyer: %v", err)
		http.Error(w, `{"error":"Failed to create flyer"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	var f models.Flyer
	var estID, catID int64
	var createdAt, updatedAt string
	database.DB.QueryRow(
		`SELECT id, titulo, descricao, preco, data_valida, imagem,
		        COALESCE(estabelecimento_id, 0), COALESCE(categoria_id, 0),
		        created_at, updated_at FROM flyers WHERE id = ?`, id).
		Scan(&f.ID, &f.Titulo, &f.Descricao, &f.Preco, &f.DataValida, &f.Imagem,
			&estID, &catID, &createdAt, &updatedAt)

	if estID > 0 {
		f.EstabelecimentoID = &estID
	}
	if catID > 0 {
		f.CategoriaID = &catID
	}
	f.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	f.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	respondJSON(w, http.StatusCreated, f)
}

func UpdateFlyer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	var req models.FlyerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		`UPDATE flyers SET titulo=?, descricao=?, preco=?, data_valida=?, imagem=?,
		 estabelecimento_id=?, categoria_id=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		req.Titulo, req.Descricao, req.Preco, req.DataValida, req.Imagem,
		req.EstabelecimentoID, req.CategoriaID, id)
	if err != nil {
		log.Printf("Error updating flyer: %v", err)
		http.Error(w, `{"error":"Failed to update flyer"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"Flyer not found"}`, http.StatusNotFound)
		return
	}

	// Return updated flyer
	var f models.Flyer
	var estID, catID int64
	var createdAt, updatedAt string
	database.DB.QueryRow(
		`SELECT id, titulo, descricao, preco, data_valida, imagem,
		        COALESCE(estabelecimento_id, 0), COALESCE(categoria_id, 0),
		        created_at, updated_at FROM flyers WHERE id = ?`, id).
		Scan(&f.ID, &f.Titulo, &f.Descricao, &f.Preco, &f.DataValida, &f.Imagem,
			&estID, &catID, &createdAt, &updatedAt)

	if estID > 0 {
		f.EstabelecimentoID = &estID
	}
	if catID > 0 {
		f.CategoriaID = &catID
	}
	f.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	f.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	respondJSON(w, http.StatusOK, f)
}

func DeleteFlyer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`DELETE FROM flyers WHERE id = ?`, id)
	if err != nil {
		log.Printf("Error deleting flyer: %v", err)
		http.Error(w, `{"error":"Failed to delete flyer"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"Flyer not found"}`, http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Flyer deleted successfully"})
}
