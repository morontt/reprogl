package repositories

import (
	"database/sql"
	"errors"
	"xelbot.com/reprogl/models"
)

type ArticleRepository struct {
	DB *sql.DB
}

func (ar *ArticleRepository) GetBySlug(slug string) (*models.Article, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.url,
			p.text_post,
			p.description,
			p.time_created,
			p.comments_count,
			c.name AS cat_name,
			c.url AS cat_url
		FROM posts AS p
		INNER JOIN category AS c ON c.id = p.category_id
		WHERE (p.url = ?)
			AND (p.hide = 0)`

	article := &models.Article{}

	err := ar.DB.QueryRow(query, slug).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Text,
		&article.Description,
		&article.CreatedAt,
		&article.CommentsCount,
		&article.CategoryName,
		&article.CategorySlug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return article, nil
}
