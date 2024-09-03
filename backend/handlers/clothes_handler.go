package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"com.fukubox/database"
	"com.fukubox/repository"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
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
	CategoryId int    `json:"category_id" validate:"required,gt=0"`
	ImageUrl   string `json:"image_url" validate:"required"`
	TagIds     []int  `json:"tag_ids"`
}

func GetClothes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")

	clothesDto, err := repository.GetClothesByUser(ctx, userId)
	if err != nil {
		log.Printf("Failed to get clothes: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	clothes := []Cloth{}
	for _, cloth := range clothesDto {
		newCloth := Cloth{
			Id:         cloth.Id,
			UserId:     cloth.UserId,
			CategoryId: cloth.CategoryId,
			ImageUrl:   cloth.ImageUrl,
			CreatedAt:  cloth.CreatedAt,
			UpdatedAt:  cloth.UpdatedAt,
		}
		err := json.Unmarshal([]byte(cloth.TagsJson), &newCloth.Tags)
		if err != nil {
			log.Printf("Failed to unmarshall ClothDto.TagsJson: \n\ttagString:%v \n\tErr:%v", cloth.TagsJson, err)
		}

		clothes = append(clothes, newCloth)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(clothes); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetClothesById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")
	clothIdStr := chi.URLParam(r, "id")

	_, err := strconv.Atoi(clothIdStr)
	if err != nil {
		log.Printf("Invalid or empty clothing id: %v", clothIdStr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	clothDto, err := repository.GetClothesByUserAndId(ctx, userId, clothIdStr)
	if err != nil {
		log.Printf("Failed to get cloth by id: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	cloth := Cloth{
		Id:         clothDto.Id,
		UserId:     clothDto.UserId,
		CategoryId: clothDto.CategoryId,
		ImageUrl:   clothDto.ImageUrl,
		CreatedAt:  clothDto.CreatedAt,
		UpdatedAt:  clothDto.UpdatedAt,
	}

	err = json.Unmarshal([]byte(clothDto.TagsJson), &cloth.Tags)
	if err != nil {
		log.Printf("Failed to unmarshall ClothDto.TagsJson: \n\ttagString:%v \n\tErr:%v", clothDto.TagsJson, err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cloth); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func CreateClothes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")

	var req ClothEdit
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request as handlers.ClothEdit: %v", err)
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

	clothId, err := repository.CreateClothWithTags(ctx, userId, repository.ClothEditDto{
		CategoryId: req.CategoryId,
		ImageUrl:   req.ImageUrl,
	}, req.TagIds)

	clothDto, err := repository.GetClothesByUserAndId(ctx, userId, strconv.Itoa(clothId))
	if err != nil {
		log.Printf("Failed to get cloth by id %v: %v", clothId, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cloth := Cloth{
		Id:         clothDto.Id,
		UserId:     clothDto.UserId,
		CategoryId: clothDto.CategoryId,
		ImageUrl:   clothDto.ImageUrl,
		CreatedAt:  clothDto.CreatedAt,
		UpdatedAt:  clothDto.UpdatedAt,
	}

	err = json.Unmarshal([]byte(clothDto.TagsJson), &cloth.Tags)
	if err != nil {
		log.Printf("Failed to unmarshall ClothDto.TagsJson: \n\ttagString:%v \n\tErr:%v", clothDto.TagsJson, err)
	}

	// Query for new item to return

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cloth); err != nil {
		log.Printf("Failed to encode response as JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func UpdateClothes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")

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

	conn := database.AcquireConnection(ctx)
	if conn == nil {
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
	ctx := r.Context()

	userId := r.Header.Get("userId")

	clothIdStr := chi.URLParam(r, "id")
	clothId, err := strconv.Atoi(clothIdStr)
	if err != nil {
		log.Printf("Invalid or empty clothing id: %v", clothIdStr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = repository.DeleteClothWithTags(ctx, userId, clothId)
	if err != nil {
		log.Printf("Failed to delete clothing item %v for user %v: %v", clothId, userId, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
