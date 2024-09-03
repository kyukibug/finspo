package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"com.fukubox/repository"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
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
	ctx := r.Context()

	userId := r.Header.Get("userId")
	tagsDto, err := repository.GetTagsByUser(ctx, userId)
	if err != nil {
		log.Printf("Failed to get clothes for user %v: %v", userId, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tags := []TagItem{}
	for _, tag := range tagsDto {
		newTag := TagItem{
			Id:        tag.Id,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
		tags = append(tags, newTag)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetTagById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")
	tagIdStr := chi.URLParam(r, "id")

	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		log.Printf("Invalid or empty tag id: %v", tagIdStr)
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	tagsDto, err := repository.GetTagByUserAndId(ctx, userId, tagId)
	if err != nil {
		log.Printf("Failed to get clothes for user %v: %v", userId, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tag := TagItem{
		Id:        tagsDto.Id,
		Name:      tagsDto.Name,
		CreatedAt: tagsDto.CreatedAt,
		UpdatedAt: tagsDto.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tag); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func CreateTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")

	var req TagEdit
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request as handlers.TagEdit: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err := validate.Struct(req)
	if err != nil {
		var validationErrors strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors.WriteString(fmt.Sprintf("%s is %s with type %s\n", err.StructField(), err.Tag(), err.Type()))
		}

		http.Error(w, validationErrors.String(), http.StatusBadRequest)
		return
	}

	tagsDto, err := repository.CreateTag(ctx, userId, req.Name)
	if err != nil {
		log.Printf("Failed to get clothes for user %v: %v", userId, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tag := TagItem{
		Id:        tagsDto.Id,
		Name:      tagsDto.Name,
		CreatedAt: tagsDto.CreatedAt,
		UpdatedAt: tagsDto.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tag); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func UpdateTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")
	tagIdStr := chi.URLParam(r, "id")

	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	var req TagEdit
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request as handlers.TagEdit: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(req)
	if err != nil {
		var validationErrors strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors.WriteString(fmt.Sprintf("%s is %s with type %s\n", err.StructField(), err.Tag(), err.Type()))
		}

		http.Error(w, validationErrors.String(), http.StatusBadRequest)
		return
	}

	tagsDto, err := repository.UpdateTag(ctx, userId, repository.TagEditDto{
		Id:   tagId,
		Name: req.Name,
	})
	if err != nil {
		log.Printf("Failed to get clothes for user %v: %v", userId, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tag := TagItem{
		Id:        tagsDto.Id,
		Name:      tagsDto.Name,
		CreatedAt: tagsDto.CreatedAt,
		UpdatedAt: tagsDto.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tag); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")
	tagIdStr := chi.URLParam(r, "id")

	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	err = repository.DeleteTag(ctx, userId, tagId)
	if err != nil {
		log.Printf("Failed to get clothes for user %v: %v", userId, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
