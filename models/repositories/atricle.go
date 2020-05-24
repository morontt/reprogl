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
			mf.path AS image_path,
			mf.description AS image_description,
			c.name AS cat_name,
			c.url AS cat_url
		FROM posts AS p
		INNER JOIN category AS c ON c.id = p.category_id
		LEFT JOIN media_file mf on p.id = mf.post_id
		WHERE (p.url = ?)
			AND (mf.id IS NULL OR mf.default_image = 1)
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
		&article.ImagePath,
		&article.ImageDescription,
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

func (ar *ArticleRepository) GetCollection(page int) (models.ArticleList, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.url,
			p.text_post,
			p.time_created,
			mf.path AS image_path,
			mf.description AS image_description,
			c.name AS cat_name,
			c.url AS cat_url
		FROM posts AS p
		INNER JOIN category AS c ON c.id = p.category_id
		LEFT JOIN media_file mf ON (p.id = mf.post_id AND mf.default_image = 1)
		WHERE p.hide = 0
		ORDER BY time_created DESC
		LIMIT 10 OFFSET ?`

	offset := 10 * (page - 1)
	rows, err := ar.DB.Query(query, offset)
	if err != nil {
		return nil, err
	}

	articles, err := populateArticles(rows)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (ar *ArticleRepository) GetCollectionByCategory(category *models.Category, page int) (models.ArticleList, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.url,
			p.text_post,
			p.time_created,
			mf.path AS image_path,
			mf.description AS image_description,
			c.name AS cat_name,
			c.url AS cat_url
		FROM posts AS p
		INNER JOIN category AS c ON c.id = p.category_id
		LEFT JOIN media_file mf ON (p.id = mf.post_id AND mf.default_image = 1)
		WHERE p.hide = 0
			AND c.id = ?
		ORDER BY time_created DESC
		LIMIT 10 OFFSET ?`

	offset := 10 * (page - 1)
	rows, err := ar.DB.Query(query, category.ID, offset)
	if err != nil {
		return nil, err
	}

	articles, err := populateArticles(rows)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (ar *ArticleRepository) GetCollectionByTag(tag *models.Tag, page int) (models.ArticleList, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.url,
			p.text_post,
			p.time_created,
			mf.path AS image_path,
			mf.description AS image_description,
			c.name AS cat_name,
			c.url AS cat_url
		FROM posts AS p
		INNER JOIN category AS c ON c.id = p.category_id
		LEFT JOIN media_file mf ON (p.id = mf.post_id AND mf.default_image = 1)
		INNER JOIN relation_topictag AS at ON p.id = at.post_id
		WHERE p.hide = 0
			AND at.tag_id = ?
		ORDER BY time_created DESC
		LIMIT 10 OFFSET ?`

	offset := 10 * (page - 1)
	rows, err := ar.DB.Query(query, tag.ID, offset)
	if err != nil {
		return nil, err
	}

	articles, err := populateArticles(rows)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (ar *ArticleRepository) newPaginator(countQuery, query string, page int, params ...interface{}) (*models.ArticlesPaginator, error) {
	var pageCount int

	err := ar.DB.QueryRow(countQuery, params...).Scan(&pageCount)
	if err != nil {
		return nil, err
	}

	offset := 10 * (page - 1)
	params = append(params, offset)
	rows, err := ar.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}

	articles, err := populateArticles(rows)
	if err != nil {
		return nil, err
	}

	return &models.ArticlesPaginator{Items: articles, CurrentPage: page, PageCount: pageCount}, nil
}

func populateArticles(rows *sql.Rows) (models.ArticleList, error) {
	var err error
	defer rows.Close()

	articles := models.ArticleList{}

	for rows.Next() {
		article := &models.ArticleListItem{}
		err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Text,
			&article.CreatedAt,
			&article.ImagePath,
			&article.ImageDescription,
			&article.CategoryName,
			&article.CategorySlug)

		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}
