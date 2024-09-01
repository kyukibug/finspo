package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"com.fukubox/database"
)

// the same as the TagItem struct in the handlers package
type TagItemDto struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// the same as the TagEdit struct in the handlers package
type TagEditDto struct {
	Name string `json:"name"`
}

func GetTagsByUser(ctx context.Context, userId string) ([]TagItemDto, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return nil, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	query := `SELECT id, name, created_at, updated_at FROM get_tags_by_user($1)`

	rows, err := conn.Query(ctx, query, userId)
	if err != nil {
		log.Printf("Query failed \n\tQuery: %v \n\tError: %v", query, err)
		return nil, err
	}
	defer rows.Close()
	return nil, nil
}
