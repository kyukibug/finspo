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

type Category struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryName struct {
	Name string `json:"name"`
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
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

	rows, err := conn.Query(ctx, "SELECT * FROM categories WHERE user_id = $1", userId)
	if err != nil {
		log.Printf("Query failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.Id, &category.UserId, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			log.Printf("Failed to scan row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetCategoriesById(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var category Category
	err = dbpool.QueryRow(ctx, "SELECT id, user_id, name, created_at, updated_at FROM categories WHERE id = $1 AND user_id = $2", categoryId, userId).
		Scan(&category.Id, &category.UserId, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		log.Printf("Failed to query row: %v", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(category); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req CategoryName
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var newCategory Category
	err = conn.QueryRow(ctx, "INSERT INTO categories (user_id, name, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING id, user_id, name, created_at, updated_at",
		userId, req.Name).Scan(&newCategory.Id, &newCategory.UserId, &newCategory.Name, &newCategory.CreatedAt, &newCategory.UpdatedAt)
	if err != nil {
		log.Printf("Failed to insert new category: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newCategory); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var req CategoryName
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var updatedCategory Category
	err = conn.QueryRow(ctx,
		"UPDATE categories SET name = $1, updated_at = now() WHERE id = $2 AND user_id = $3 RETURNING id, user_id, name, created_at, updated_at",
		req.Name, categoryId, userId).Scan(&updatedCategory.Id, &updatedCategory.UserId, &updatedCategory.Name, &updatedCategory.CreatedAt, &updatedCategory.UpdatedAt)
	if err != nil {
		log.Printf("Failed to update category: %v", err)
		http.Error(w, "Category not found or not authorized to update", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedCategory); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()
	ctx := r.Context()

	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		log.Printf("Failed to acquire a database connection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	commandTag, err := conn.Exec(ctx, "DELETE FROM categories WHERE id = $1 AND user_id = $2", categoryId, userId)
	if err != nil {
		log.Printf("Failed to delete category: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Category not found or not authorized to delete", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
