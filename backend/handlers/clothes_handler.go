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
		return
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

	clothId, err := strconv.Atoi(clothIdStr)
	if err != nil {
		log.Printf("Invalid or empty clothing id: %v", clothIdStr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	clothDto, err := repository.GetClothesByUserAndId(ctx, userId, clothId)
	if err != nil {
		log.Printf("Failed to get cloth by id: %v", err)
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

	clothDto, err := repository.GetClothesByUserAndId(ctx, userId, clothId)
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

	repository.UpdateCloth(ctx, userId, clothId, repository.ClothEditDto{
		CategoryId: req.CategoryId,
		ImageUrl:   req.ImageUrl,
	})

	clothDto, err := repository.GetClothesByUserAndId(ctx, userId, clothId)
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cloth); err != nil {
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

func AddTagToCloth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")
	clothIdStr := chi.URLParam(r, "id")
	tagIdStr := chi.URLParam(r, "tagId")

	clothId, err := strconv.Atoi(clothIdStr)
	if err != nil {
		log.Printf("Invalid or empty clothing id: %v", clothIdStr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		log.Printf("Invalid or empty tag id: %v", tagIdStr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Ensure clothing item belongs to user
	result, err := repository.GetClothesByUserAndId(ctx, userId, clothId)
	if err != nil {
		log.Printf("Failed to get cloth by id: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if result.Id == 0 {
		// failed to find item
		log.Printf("Cloth item %v not found for user %v", clothId, userId)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = repository.AddTagToCloth(ctx, clothId, tagId)
	if err != nil {
		log.Printf("Failed to get cloth by id: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteTagFromCloth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.Header.Get("userId")
	clothIdStr := chi.URLParam(r, "id")
	tagIdStr := chi.URLParam(r, "tagId")

	clothId, err := strconv.Atoi(clothIdStr)
	if err != nil {
		log.Printf("Invalid or empty clothing id: %v", clothIdStr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		log.Printf("Invalid or empty tag id: %v", tagIdStr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Ensure clothing item belongs to user
	result, err := repository.GetClothesByUserAndId(ctx, userId, clothId)
	if err != nil {
		log.Printf("Failed to get cloth by id: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if result.Id == 0 {
		// failed to find item
		log.Printf("Cloth item %v not found for user %v", clothId, userId)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = repository.DeleteTagFromCloth(ctx, clothId, tagId)
	if err != nil {
		log.Printf("Failed to get cloth by id: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
