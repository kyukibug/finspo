package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GIVE ME MY CATEGORY GOPHER PLUSHIE"))
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
