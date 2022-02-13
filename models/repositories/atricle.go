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

func (ar *ArticleRepository) GetCollection(page int) (*models.ArticlesPaginator, error) {
	countQuery := `
		SELECT
			COUNT(p.id) AS cnt
		FROM posts AS p
		WHERE p.hide = 0`

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

	params := make([]interface{}, 0)

	return ar.newPaginator(countQuery, query, page, params...)
}

func (ar *ArticleRepository) GetCollectionByCategory(category *models.Category, page int) (*models.ArticlesPaginator, error) {
	countQuery := `
		SELECT
			COUNT(p.id) AS cnt
		FROM posts AS p
		INNER JOIN category AS c ON c.id = p.category_id
		WHERE p.hide = 0
			AND c.tree_left_key >= ?
			AND c.tree_right_key <= ?`

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
			AND c.tree_left_key >= ?
			AND c.tree_right_key <= ?
		ORDER BY time_created DESC
		LIMIT 10 OFFSET ?`

	params := make([]interface{}, 0)
	params = append(params, category.LeftKey, category.RightKey)

	return ar.newPaginator(countQuery, query, page, params...)
}

func (ar *ArticleRepository) GetCollectionByTag(tag *models.Tag, page int) (*models.ArticlesPaginator, error) {
	countQuery := `
		SELECT
			COUNT(p.id) AS cnt
		FROM posts AS p
		INNER JOIN relation_topictag AS at ON p.id = at.post_id
		WHERE p.hide = 0
			AND at.tag_id = ?`

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

	params := make([]interface{}, 0)
	params = append(params, tag.ID)

	return ar.newPaginator(countQuery, query, page, params...)
}

func (ar *ArticleRepository) GetSitemapCollection() (*models.SitemapItemList, error) {
	query := `
		SELECT
			url,
			last_update
		FROM posts
		WHERE
			hide = 0
		ORDER BY time_created DESC
`

	rows, err := ar.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	articles := models.SitemapItemList{}

	for rows.Next() {
		item := &models.SitemapItem{}
		err = rows.Scan(
			&item.Slug,
			&item.UpdatedAt)

		if err != nil {
			return nil, err
		}

		articles = append(articles, item)
	}

	return &articles, nil
}

func (ar *ArticleRepository) newPaginator(countQuery, query string, page int, params ...interface{}) (*models.ArticlesPaginator, error) {
	var articleCount int

	err := ar.DB.QueryRow(countQuery, params...).Scan(&articleCount)
	if err != nil {
		return nil, err
	}

	pageCount := articleCount / 10
	if articleCount%10 != 0 {
		pageCount += 1
	}

	if page > pageCount {
		return nil, models.RecordNotFound
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
