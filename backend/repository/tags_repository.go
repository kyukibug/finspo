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
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// the same as the TagEdit struct in the handlers package
type TagEditDto struct {
	Id   int
	Name string
}

func GetTagsByUser(ctx context.Context, userId string) ([]TagItemDto, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return nil, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	query := `SELECT id, name, created_at, updated_at FROM tags WHERE user_id=$1`

	rows, err := conn.Query(ctx, query, userId)
	if err != nil {
		log.Printf("Query failed \n\tQuery: %v \n\tError: %v", query, err)
		return nil, err
	}
	defer rows.Close()
	return nil, nil
}

func GetTagByUserAndId(ctx context.Context, userId string, tagId int) (TagItemDto, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return TagItemDto{}, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	query := `SELECT id, name, created_at, updated_at FROM tags WHERE user_id=$1 AND id=$2`

	var tag TagItemDto

	err := conn.QueryRow(ctx, query, userId, tagId).Scan(
		&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)

	if err != nil {
		log.Printf("Query failed \n\tQuery: %v \n\tError: %v", query, err)
		return TagItemDto{}, err
	}

	return tag, nil
}

func CreateTag(ctx context.Context, userId string, tagName string) (TagItemDto, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return TagItemDto{}, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	query := `INSERT INTO tags (user_id, name) VALUES($1, $2)
			  RETURNING id, name, created_at, updated_at`

	var tag TagItemDto

	err := conn.QueryRow(ctx, query, userId, tagName).Scan(
		&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)

	if err != nil {
		log.Printf("Query failed \n\tQuery: %v \n\tError: %v", query, err)
		return TagItemDto{}, err
	}

	return tag, nil
}

func UpdateTag(ctx context.Context, userId string, updatedTag TagEditDto) (TagItemDto, error) {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return TagItemDto{}, errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	query := `UPDATE tags SET name = $1, updated_at = now() WHERE id = $2 AND user_id = $3
			  RETURNING id, name, created_at, updated_at`

	var tag TagItemDto

	err := conn.QueryRow(ctx, query, updatedTag.Name, updatedTag.Id, userId).Scan(
		&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)

	if err != nil {
		log.Printf("Query failed \n\tQuery: %v \n\tError: %v", query, err)
		return TagItemDto{}, err
	}

	return tag, nil
}

func DeleteTag(ctx context.Context, userId string, tagId int) error {
	conn := database.AcquireConnection(ctx)
	if conn == nil {
		return errors.New("failed to acquire database connection")
	}
	defer conn.Release()

	query := `DELETE FROM tags WHERE id = $1 AND user_id = $2`

	commandTag, err := conn.Exec(ctx, query, tagId, userId)

	if err != nil {
		log.Printf("Query failed \n\tQuery: %v \n\tError: %v", query, err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		log.Printf("Tag %v not found for user %v", tagId, userId)
		return errors.New("Tag not found")
	}

	return nil
}
