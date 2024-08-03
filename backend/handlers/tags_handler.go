package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"com.fukubox/database"
	"github.com/go-chi/chi"
)

type TagItem struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TagEdit struct {
	Name string `json:"name"`
}

func GetTags(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT id, name, created_at, updated_at FROM tags")
	if err != nil {
		log.Printf("Query failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tags := []TagItem{}

	for rows.Next() {
		var tag TagItem
		if err := rows.Scan(&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			log.Printf("Failed to scan row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetTagById(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	tagIdStr := chi.URLParam(r, "id")
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var tag TagItem
	err = dbpool.QueryRow(ctx, "SELECT id, name, created_at, updated_at FROM tags WHERE id = $1", tagId).
		Scan(&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		log.Printf("Failed to query row: %v", err)
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tag); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func CreateTag(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	var req TagEdit
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Tag name is required", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var newTag TagItem
	err = conn.QueryRow(
		ctx,
		`INSERT INTO tags (name, created_at, updated_at)
     	VALUES ($1, now(), now())
     	RETURNING id, name, created_at, updated_at`,
		req.Name).Scan(&newTag.Id, &newTag.Name, &newTag.CreatedAt, &newTag.UpdatedAt)
	if err != nil {
		log.Printf("Failed to insert new tag: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newTag); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func UpdateTag(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	tagIdStr := chi.URLParam(r, "id")
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	var req TagEdit
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Tag name is required", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var updatedTag TagItem
	err = conn.QueryRow(ctx,
		"UPDATE tags SET name = $1, updated_at = now() WHERE id = $2 RETURNING id, name, created_at, updated_at",
		req.Name, tagId).Scan(&updatedTag.Id, &updatedTag.Name, &updatedTag.CreatedAt, &updatedTag.UpdatedAt)
	if err != nil {
		log.Printf("Failed to update tag: %v", err)
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTag); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	tagIdStr := chi.URLParam(r, "id")
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	commandTag, err := conn.Exec(ctx, "DELETE FROM tags WHERE id = $1", tagId)
	if err != nil {
		log.Printf("Failed to delete tag: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
