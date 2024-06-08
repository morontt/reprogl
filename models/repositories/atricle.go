package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"xelbot.com/reprogl/models"
)

const RecentPostsCount = 6

type ArticleRepository struct {
	DB *sql.DB
}

func (ar *ArticleRepository) GetBySlug(slug string, isAdmin bool) (*models.Article, error) {
	ds := goqu.Dialect("mysql8").From(goqu.T("posts").As("p")).Select(
		"p.id",
		"p.title",
		"p.url",
		"p.text_post",
		"p.description",
		"p.time_created",
		"p.last_update",
		"p.comments_count",
		"p.views_count",
		"p.disable_comments",
		goqu.I("mf.path").As("image_path"),
		goqu.I("mf.width").As("image_width"),
		"mf.src_set",
		goqu.I("mf.description").As("image_alt"),
		goqu.I("c.name").As("cat_name"),
		goqu.I("c.url").As("cat_url"),
	).InnerJoin(
		goqu.T("category").As("c"),
		goqu.On(goqu.Ex{
			"c.id": goqu.I("p.category_id"),
		}),
	).LeftJoin(
		goqu.T("media_file").As("mf"),
		goqu.On(goqu.Ex{
			"mf.post_id":       goqu.I("p.id"),
			"mf.default_image": goqu.L("1"),
		}),
	).Where(
		goqu.Ex{
			"p.url": slug,
		},
	)

	if !isAdmin {
		ds = ds.Where(goqu.Ex{
			"p.hide": goqu.L("0"),
		})
	}

	query, params, _ := ds.Prepared(true).ToSQL()

	article := &models.Article{}

	err := ar.DB.QueryRow(query, params...).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Text,
		&article.Description,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.CommentsCount,
		&article.Views,
		&article.DisabledComments,
		&article.ImagePath,
		&article.Width,
		&article.SrcSet,
		&article.Alt,
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

func (ar *ArticleRepository) GetIdBySlug(slug string) int {
	var id int

	err := ar.DB.QueryRow(
		`SELECT id FROM posts WHERE url = ?`,
		slug,
	).Scan(&id)

	if err != nil {
		return 0
	}

	return id
}

func (ar *ArticleRepository) GetCollection(page int, isAdmin bool) (*models.ArticlesPaginator, error) {
	countQuery := `
		SELECT
			COUNT(p.id) AS cnt
		FROM posts AS p
		WHERE p.hide = 0`

	ds := goqu.Dialect("mysql8").From(goqu.T("posts").As("p")).Select(
		"p.id",
		"p.title",
		"p.url",
		"p.text_post",
		"p.preview",
		"p.time_created",
		"p.comments_count",
		"p.hide",
		"mf.picture_tag",
		goqu.I("c.name").As("cat_name"),
		goqu.I("c.url").As("cat_url"),
	).InnerJoin(
		goqu.T("category").As("c"),
		goqu.On(goqu.Ex{
			"c.id": goqu.I("p.category_id"),
		}),
	).LeftJoin(
		goqu.T("media_file").As("mf"),
		goqu.On(goqu.Ex{
			"mf.post_id":       goqu.I("p.id"),
			"mf.default_image": goqu.L("1"),
		}),
	).Order(goqu.I("p.time_created").Desc())

	if !isAdmin {
		ds = ds.Where(goqu.Ex{
			"p.hide": goqu.L("0"),
		})
	}

	query, _, _ := ds.ToSQL()
	query += " LIMIT 10 OFFSET ?"

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
			p.preview,
			p.time_created,
			p.comments_count,
			p.hide,
			mf.picture_tag,
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
			p.preview,
			p.time_created,
			p.comments_count,
			p.hide,
			mf.picture_tag,
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
			updated_at
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
		item := models.SitemapItem{}
		err = rows.Scan(
			&item.Slug,
			&item.UpdatedAt)

		if err != nil {
			return nil, err
		}

		articles = append(articles, &item)
	}

	return &articles, nil
}

func (ar *ArticleRepository) GetFeedCollection() (models.FeedItemList, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.url,
			p.text_post,
			mf.src_set,
			p.updated_at,
			p.time_created
		FROM posts AS p
		LEFT JOIN media_file AS mf ON (p.id = mf.post_id AND mf.default_image = 1)
		WHERE
			hide = 0
		ORDER BY time_created DESC
		LIMIT 25
`

	rows, err := ar.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	articles := models.FeedItemList{}

	for rows.Next() {
		item := models.FeedItem{}
		err = rows.Scan(
			&item.ID,
			&item.Title,
			&item.Slug,
			&item.Text,
			&item.SrcSet,
			&item.UpdatedAt,
			&item.CreatedAt)

		if err != nil {
			return nil, err
		}

		articles = append(articles, &item)
	}

	return articles, nil
}

func (ar *ArticleRepository) GetRecentPostsCollection(articleId int) (*models.RecentPostList, error) {
	query := fmt.Sprintf(`
		SELECT
			title,
			url
		FROM posts
		WHERE
			hide = 0
			AND id != ?
		ORDER BY time_created DESC
		LIMIT %d`, RecentPostsCount)

	rows, err := ar.DB.Query(query, articleId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	articles := models.RecentPostList{}

	for rows.Next() {
		item := models.RecentPost{}
		err = rows.Scan(
			&item.Title,
			&item.Slug)

		if err != nil {
			return nil, err
		}

		articles = append(articles, &item)
	}

	return &articles, nil
}

func (ar *ArticleRepository) GetLastRecentPostsID() (int, error) {
	query := fmt.Sprintf(`
		SELECT
			MIN(src.id) AS id
		FROM (
			SELECT id
			FROM posts
			WHERE hide = 0
			ORDER BY time_created DESC
			LIMIT %d) AS src`, RecentPostsCount)

	var id int
	err := ar.DB.QueryRow(query).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ar *ArticleRepository) GetByIdForComment(id int) (*models.ArticleForComment, error) {
	query := `
		SELECT
			p.id,
			p.url,
			p.hide
		FROM posts AS p
		WHERE (p.id = ?)`

	article := &models.ArticleForComment{}

	err := ar.DB.QueryRow(query, id).Scan(
		&article.ID,
		&article.Slug,
		&article.Hidden)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return article, nil
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

	defer rows.Close()

	articles, err := populateArticles(rows)
	if err != nil {
		return nil, err
	}

	return &models.ArticlesPaginator{Items: articles, CurrentPage: page, PageCount: pageCount}, nil
}

func (ar *ArticleRepository) GetMostVisitedArticlesOfMonth() ([]models.ArticleStatItem, error) {
	query := `
		SELECT
			p.title,
			p.url,
			COUNT(t.id) AS cnt
		FROM posts AS p
		INNER JOIN tracking AS t ON t.post_id = p.id
		INNER JOIN tracking_agent AS ta ON t.user_agent_id = ta.id
		WHERE t.time_created > ?
			AND ta.is_bot = 0
		GROUP BY p.id
		ORDER BY cnt DESC
		LIMIT 6`

	from := time.Now().Add(-30 * 24 * time.Hour).Format(time.DateTime)

	rows, err := ar.DB.Query(query, from)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := make([]models.ArticleStatItem, 0, 6)
	for rows.Next() {
		item := models.ArticleStatItem{}
		err = rows.Scan(
			&item.Title,
			&item.Slug,
			&item.Views)

		if err != nil {
			return nil, err
		}

		data = append(data, item)
	}

	return data, nil
}

func (ar *ArticleRepository) GetMostVisitedArticles() ([]models.ArticleStatItem, error) {
	query := `
		SELECT
			p.title,
			p.url,
			p.views_count AS cnt
		FROM posts AS p
		ORDER BY cnt DESC
		LIMIT 6`

	rows, err := ar.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := make([]models.ArticleStatItem, 0, 6)
	for rows.Next() {
		item := models.ArticleStatItem{}
		err = rows.Scan(
			&item.Title,
			&item.Slug,
			&item.Views)

		if err != nil {
			return nil, err
		}

		data = append(data, item)
	}

	return data, nil
}

func populateArticles(rows *sql.Rows) (models.ArticleList, error) {
	var err error

	articles := models.ArticleList{}

	for rows.Next() {
		article := &models.ArticleListItem{}
		err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Text,
			&article.Preview,
			&article.CreatedAt,
			&article.CommentsCount,
			&article.Hidden,
			&article.PictureTag,
			&article.CategoryName,
			&article.CategorySlug)

		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}
