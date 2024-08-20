package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"com.fukubox/database"
	"github.com/jackc/pgx/v5"
)

type ClothDto struct {
	Id         int
	UserId     int
	CategoryId int
	ImageUrl   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	TagsJson   string
}

type ClothEditDto struct {
	CategoryId int
	ImageUrl   string
}

func GetClothesByUser(ctx context.Context, userId string) ([]ClothDto, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return nil, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	query := `SELECT id, user_id, category_id, image_url, created_at, updated_at, tags FROM get_clothes_by_user($1)`

	rows, err := conn.Query(ctx, query, userId)
	if err != nil {
		log.Printf("Query failed \n\tQuery: %v \n\tError: %v", query, err)
		return nil, err
	}
	defer rows.Close()

	clothes := []ClothDto{}

	for rows.Next() {
		var cloth ClothDto
		if err := rows.Scan(&cloth.Id, &cloth.UserId, &cloth.CategoryId, &cloth.ImageUrl, &cloth.CreatedAt, &cloth.UpdatedAt, &cloth.TagsJson); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		clothes = append(clothes, cloth)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v", err)
		return nil, err
	}

	return clothes, nil
}

func GetClothesByUserAndId(ctx context.Context, userId string, clothId string) (ClothDto, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return ClothDto{}, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	query := `SELECT id, user_id, category_id, image_url, created_at, updated_at, tags FROM get_clothes_by_user_and_id($1, $2)`

	var cloth ClothDto

	err := conn.QueryRow(ctx, query, userId, clothId).
		Scan(&cloth.Id, &cloth.UserId, &cloth.CategoryId, &cloth.ImageUrl, &cloth.CreatedAt, &cloth.UpdatedAt, &cloth.TagsJson)
	if err != nil {
		log.Printf("Query failed \n\tQuery: %v \n\tError: %v", query, err)
		log.Printf("Failed to query %v with params {user_id: %v, id:%v}: %v", query, userId, clothId, err)
		return ClothDto{}, err
	}

	return cloth, nil
}

func CreateCloth(ctx context.Context, userId string, newCloth ClothEditDto) (int, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return -1, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	query := `INSERT INTO clothing_items (user_id, category_id, image_url, created_at, updated_at)
			  VALUES ($1, $2, $3, now(), now()) RETURNING id, user_id, category_id, image_url, created_at, updated_at
			  RETURNING id`

	var id int
	err := conn.QueryRow(ctx, query, userId, newCloth.CategoryId, newCloth.ImageUrl).Scan(&id)
	if err != nil {
		log.Printf("Failed to insert new clothing item: %v", err)
		return -1, err
	}

	return id, nil
}

func CreateClothTx(tx pgx.Tx, ctx context.Context, userId string, newCloth ClothEditDto) (int, error) {

	query := `INSERT INTO clothing_items (user_id, category_id, image_url, created_at, updated_at)
			  VALUES ($1, $2, $3, now(), now()) 
			  RETURNING id`

	var id int
	err := tx.QueryRow(ctx, query, userId, newCloth.CategoryId, newCloth.ImageUrl).Scan(&id)
	if err != nil {
		log.Printf("Failed to insert new clothing item: %v", err)
		return -1, err
	}

	return id, nil
}

func CreateClothWithTags(ctx context.Context, userId string, newCloth ClothEditDto, tags []int) (int, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return -1, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Begin Transation Failure: %v", err)
		return -1, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	id, err := CreateClothTx(tx, ctx, userId, newCloth)
	if err != nil {
		log.Printf("Failed to create cloth: %v", err)
		return -1, err
	}

	err = BindTagsTx(tx, ctx, id, tags)
	if err != nil {
		log.Printf("Failed to bind cloth tags: %v", err)
		return -1, err
	}

	return id, nil
}

func UpdateCloth(ctx context.Context) (int, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return -1, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	return 1, nil
}

func BindTagsTx(tx pgx.Tx, ctx context.Context, clothId int, tags []int) error {

	query := `INSERT INTO clothing_item_tags (clothing_item_id, tag_id) VALUES ($1, $2)`

	batch := &pgx.Batch{}

	for _, tagId := range tags {
		batch.Queue(query, clothId, tagId)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for _, tagId := range tags {
		_, err := results.Exec()
		if err != nil {
			log.Printf("Failed to insert tag %v: %v", tagId, err)
			return err
		}
	}

	return nil
}
