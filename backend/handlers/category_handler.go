package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"com.fukubox/database"
	"github.com/go-chi/chi"
)

type Category struct {
	id         int    `json:"id"`
	user_id    int    `json:"user_id"`
	name       string `json:"name"`
	created_at string `json:"created_at"`
	updated_at string `json:"updated_at"`
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	dbpool := database.GetDB()

	user_id := r.Header.Get("userId")
	if user_id == "" {
		http.Error(w, "User ID is not set", http.StatusBadRequest)
		return
	}

	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		http.Error(w, "Failed to acquire a database connection", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT * FROM categories WHERE user_id = $1", user_id)
	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.id, &category.user_id, &category.name, &category.created_at, &category.updated_at); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error after iterating rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
		return
	}
}

func GetCategoriesById(w http.ResponseWriter, r *http.Request) {
	category_id_str := chi.URLParam(r, "category_id")
	category_id, err := strconv.Atoi(category_id_str)
	if err != nil {
		http.Error(w, "Invalid category ID", 400)
		return
	}
	w.Write([]byte(fmt.Sprintf("GOPHERGOPHER:%d", category_id)))
}
