package repositories

import (
	"database/sql"
	"errors"
	"xelbot.com/reprogl/models"
)

type CategoryRepository struct {
	DB *sql.DB
}

func (cr *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.url
		FROM category AS c
		WHERE (c.url = ?)`

	category := &models.Category{}

	err := cr.DB.QueryRow(query, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return category, nil
}
