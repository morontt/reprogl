package repositories

import (
	"database/sql"
	"errors"
	"strings"

	"xelbot.com/reprogl/models"
)

type TagRepository struct {
	DB *sql.DB
}

func (tr *TagRepository) GetBySlug(slug string) (*models.Tag, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.url
		FROM tags AS c
		WHERE (c.url = ?)`

	tag := &models.Tag{}

	err := tr.DB.QueryRow(query, slug).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Slug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return tag, nil
}

func (tr *TagRepository) GetCollectionByArticle(article *models.Article) (models.TagList, error) {
	query := `
		SELECT
			t.id,
			t.name,
			t.url
		FROM tags AS t
		INNER JOIN relation_topictag AS tt ON tt.tag_id = t.id
		WHERE tt.post_id = ?`

	rows, err := tr.DB.Query(query, article.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tags := models.TagList{}

	for rows.Next() {
		tag := &models.Tag{}
		err = rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Slug)

		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (tr *TagRepository) PopulateTagsToArticles(articles models.ArticleList) error {
	idx := make([]interface{}, 0, len(articles))
	for _, a := range articles {
		idx = append(idx, a.ID)
	}

	query := `
		SELECT
			t.id,
			t.name,
			t.url,
			tt.post_id
		FROM tags AS t
		INNER JOIN relation_topictag AS tt ON tt.tag_id = t.id
		WHERE tt.post_id IN (?` + strings.Repeat(", ?", len(idx)-1) + ")"

	rows, err := tr.DB.Query(query, idx...)
	if err != nil {
		return err
	}

	defer rows.Close()

	var articleId int
	var tags models.TagList
	var ok bool

	tagsMap := make(map[int]models.TagList)

	for rows.Next() {
		tag := &models.Tag{}
		err = rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Slug,
			&articleId)

		if err != nil {
			return err
		}

		if _, ok = tagsMap[articleId]; ok {
			tagsMap[articleId] = append(tagsMap[articleId], tag)
		} else {
			tagsMap[articleId] = models.TagList{tag}
		}
	}

	for _, a := range articles {
		if tags, ok = tagsMap[a.ID]; ok {
			a.Tags = tags
		} else {
			a.Tags = make(models.TagList, 0)
		}
	}

	return nil
}
