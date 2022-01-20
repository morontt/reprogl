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
			c.url,
			c.tree_left_key,
			c.tree_right_key,
			c.tree_depth
		FROM category AS c
		WHERE (c.url = ?)`

	category := &models.Category{}

	err := cr.DB.QueryRow(query, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.LeftKey,
		&category.RightKey,
		&category.Depth)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return category, nil
}

func (cr *CategoryRepository) GetCategoryTree() (*models.CategoryList, error) {
	query := `
		SELECT c.id,
			c.name,
			c.url,
			c.tree_depth
		FROM category AS c,
			(SELECT c0.id
				FROM category AS c0,
					category AS c1
				INNER JOIN posts AS p ON c1.id = p.category_id
				WHERE c0.tree_left_key <= c1.tree_left_key
					AND c0.tree_right_key >= c1.tree_right_key
					AND p.hide = 0
				GROUP BY c0.id) AS cnt
		WHERE c.id = cnt.id
		ORDER BY c.tree_left_key`

	rows, err := cr.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := models.CategoryList{}

	for rows.Next() {
		category := &models.Category{}
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Depth)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return &categories, nil
}
