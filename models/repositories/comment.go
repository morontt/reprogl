package repositories

import (
	"database/sql"
	"xelbot.com/reprogl/models"
)

type CommentRepository struct {
	DB *sql.DB
}

func (cr *CommentRepository) GetLastUpdate(articleId int) (string, error) {
	var last string
	query := `
		SELECT
			DATE_FORMAT(src.last_update, '%y%j%H%i%s') AS last
		FROM (
			SELECT
				MAX(last_update) AS last_update
			FROM comments
			WHERE (post_id = ?)) AS src`

	err := cr.DB.QueryRow(query, articleId).Scan(&last)

	return last, err
}

func (cr *CommentRepository) GetCollectionByArticleId(articleId int) (*models.CommentList, error) {
	query := `
		SELECT
			c.id,
			COALESCE(t.name, u.username) AS username,
			COALESCE(t.mail, u.mail) AS email,
			t.website,
			c.text,
			c.tree_depth,
			c.time_created,
			c.deleted
		FROM comments AS c
		LEFT JOIN commentators AS t ON c.commentator_id = t.id
		LEFT JOIN users AS u ON c.user_id = u.id
		WHERE (c.post_id = ?)
		ORDER BY c.tree_left_key`

	rows, err := cr.DB.Query(query, articleId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := models.CommentList{}

	for rows.Next() {
		comment := &models.Comment{}
		err = rows.Scan(
			&comment.ID,
			&comment.Name,
			&comment.Email,
			&comment.Website,
			&comment.Text,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.Deleted)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return &comments, nil
}
