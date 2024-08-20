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

type Cloth struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	CategoryId int       `json:"category_id"`
	ImageUrl   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Tags       []Tag     `json:"tags"`
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ClothEdit struct {
	CategoryId int    `json:"category_id"`
	ImageUrl   string `json:"image_url"`
	TagIds     []int  `json:"tag_ids"`
}

func GetClothes(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT id, user_id, category_id, image_url, created_at, updated_at
	FROM clothing_items WHERE user_id = $1`, userId)
	if err != nil {
		log.Printf("Query failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	clothes := []Cloth{}

	for rows.Next() {
		var cloth Cloth
		if err := rows.Scan(&cloth.Id, &cloth.UserId, &cloth.CategoryId, &cloth.ImageUrl, &cloth.CreatedAt, &cloth.UpdatedAt); err != nil {
			log.Printf("Failed to scan row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tagConn, err := dbpool.Acquire(ctx)
		if err != nil {
			log.Printf("Failed to acquire a database connection for tags: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tagRows, err := tagConn.Query(ctx, `SELECT t.id, t.name FROM tags t
			JOIN clothing_item_tags cit ON t.id = cit.tag_id
			WHERE cit.clothing_item_id = $1`, cloth.Id)
		if err != nil {
			tagConn.Release()
			log.Printf("Failed to query tags: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var tags []Tag
		for tagRows.Next() {
			var tag Tag
			if err := tagRows.Scan(&tag.Id, &tag.Name); err != nil {
				tagRows.Close()
				tagConn.Release()
				log.Printf("Failed to scan tag: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			tags = append(tags, tag)
		}

		tagRows.Close()
		tagConn.Release()

		cloth.Tags = tags
		clothes = append(clothes, cloth)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(clothes); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetClothesById(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	clothIdStr := chi.URLParam(r, "id")
	clothId, err := strconv.Atoi(clothIdStr)
	if err != nil {
		http.Error(w, "Invalid clothing item ID", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var cloth Cloth
	err = dbpool.QueryRow(ctx, "SELECT id, user_id, category_id, image_url, created_at, updated_at FROM clothing_items WHERE id = $1 AND user_id = $2", clothId, userId).
		Scan(&cloth.Id, &cloth.UserId, &cloth.CategoryId, &cloth.ImageUrl, &cloth.CreatedAt, &cloth.UpdatedAt)
	if err != nil {
		log.Printf("Failed to query row: %v", err)
		http.Error(w, "Clothing Item not found", http.StatusNotFound)
		return
	}

	tagRows, err := conn.Query(ctx, `SELECT t.id, t.name FROM tags t
		JOIN clothing_item_tags cit ON t.id = cit.tag_id
		WHERE cit.clothing_item_id = $1`, cloth.Id)
	if err != nil {
		log.Printf("Failed to query tags: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tagRows.Close()

	var tags []Tag
	for tagRows.Next() {
		var tag Tag
		if err := tagRows.Scan(&tag.Id, &tag.Name); err != nil {
			log.Printf("Failed to scan tag: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tags = append(tags, tag)
	}
	cloth.Tags = tags

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cloth); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func CreateClothes(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req ClothEdit
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.CategoryId == 0 || req.ImageUrl == "" {
		http.Error(w, "Category ID and Image URL are required", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var newCloth Cloth
	err = conn.QueryRow(ctx, `INSERT INTO clothing_items (user_id, category_id, image_url, created_at, updated_at)
	VALUES ($1, $2, $3, now(), now()) RETURNING id, user_id, category_id, image_url, created_at, updated_at`,
		userId, req.CategoryId, req.ImageUrl).Scan(&newCloth.Id, &newCloth.UserId, &newCloth.CategoryId, &newCloth.ImageUrl, &newCloth.CreatedAt, &newCloth.UpdatedAt)
	if err != nil {
		log.Printf("Failed to insert new clothing item: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, tagId := range req.TagIds {
		_, err := conn.Exec(ctx, "INSERT INTO clothing_item_tags (clothing_item_id, tag_id) VALUES ($1, $2)", newCloth.Id, tagId)
		if err != nil {
			log.Printf("Failed to insert tag association: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	tagRows, err := conn.Query(ctx, `SELECT t.id, t.name FROM tags t
		JOIN clothing_item_tags cit ON t.id = cit.tag_id
		WHERE cit.clothing_item_id = $1`, newCloth.Id)
	if err != nil {
		log.Printf("Failed to query tags: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tagRows.Close()

	var tags []Tag
	for tagRows.Next() {
		var tag Tag
		if err := tagRows.Scan(&tag.Id, &tag.Name); err != nil {
			log.Printf("Failed to scan tag: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tags = append(tags, tag)
	}
	newCloth.Tags = tags

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newCloth); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func UpdateClothes(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	clothIdStr := chi.URLParam(r, "id")
	clothId, err := strconv.Atoi(clothIdStr)
	if err != nil {
		http.Error(w, "Invalid clothing item ID", http.StatusBadRequest)
		return
	}

	var req ClothEdit
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.CategoryId == 0 || req.ImageUrl == "" {
		http.Error(w, "Category ID and Image URL are required", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var updatedCloth Cloth
	err = conn.QueryRow(ctx,
		`UPDATE clothing_items SET category_id = $1, image_url = $2,
		updated_at = now() WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, category_id, image_url, created_at, updated_at`,
		req.CategoryId, req.ImageUrl, clothId, userId).Scan(
		&updatedCloth.Id, &updatedCloth.UserId, &updatedCloth.CategoryId, &updatedCloth.ImageUrl,
		&updatedCloth.CreatedAt, &updatedCloth.UpdatedAt)
	if err != nil {
		log.Printf("Failed to update clothing item: %v", err)
		http.Error(w, "Clothing item not found or not authorized to update", http.StatusNotFound)
		return
	}

	_, err = conn.Exec(ctx, "DELETE FROM clothing_item_tags WHERE clothing_item_id = $1", clothId)
	if err != nil {
		log.Printf("Failed to delete existing tags: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, tagId := range req.TagIds {
		_, err := conn.Exec(ctx, "INSERT INTO clothing_item_tags (clothing_item_id, tag_id) VALUES ($1, $2)", clothId, tagId)
		if err != nil {
			log.Printf("Failed to insert new tag associations: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	tagRows, err := conn.Query(ctx, `SELECT t.id, t.name FROM tags t
		JOIN clothing_item_tags cit ON t.id = cit.tag_id
		WHERE cit.clothing_item_id = $1`, clothId)
	if err != nil {
		log.Printf("Failed to query tags: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tagRows.Close()

	var tags []Tag
	for tagRows.Next() {
		var tag Tag
		if err := tagRows.Scan(&tag.Id, &tag.Name); err != nil {
			log.Printf("Failed to scan tag: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tags = append(tags, tag)
	}
	updatedCloth.Tags = tags

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedCloth); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func DeleteClothes(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	clothIdStr := chi.URLParam(r, "id")
	clothId, err := strconv.Atoi(clothIdStr)
	if err != nil {
		http.Error(w, "Invalid clothing item ID", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "DELETE FROM clothing_item_tags WHERE clothing_item_id = $1", clothId)
	if err != nil {
		log.Printf("Failed to delete associated tags: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	commandTag, err := conn.Exec(ctx, "DELETE FROM clothing_items WHERE id = $1 AND user_id = $2", clothId, userId)
	if err != nil {
		log.Printf("Failed to delete clothing item: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Clothing item not found or not authorized to delete", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
