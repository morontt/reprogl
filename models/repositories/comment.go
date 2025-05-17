package repositories

import (
	"database/sql"
	"errors"

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
				COALESCE(MAX(last_update), '1981-07-18 00:00:00') AS last_update
			FROM comments
			WHERE (post_id = ?)) AS src`

	err := cr.DB.QueryRow(query, articleId).Scan(&last)

	return last, err
}

func (cr *CommentRepository) GetCollectionByArticleId(articleId int) (models.CommentList, error) {
	query := `
		SELECT
			c.id,
			COALESCE(t.name, u.display_name, u.username) AS username,
			COALESCE(t.mail, u.mail) AS email,
			t.website,
			COALESCE(t.gender, u.gender) AS gender,
			c.commentator_id,
			c.user_id,
			c.text,
			c.tree_depth,
			COALESCE(c.force_created_at, c.time_created) AS time_created,
			COALESCE(t.rotten_link, 0) AS rotten_link,
			COALESCE(u.avatar_variant, t.avatar_variant) AS avatar_variant,
			c.deleted
		FROM comments AS c
		LEFT JOIN commentators AS t ON c.commentator_id = t.id
		LEFT JOIN users AS u ON c.user_id = u.id
		INNER JOIN (
			SELECT
				c1.id
			FROM
				comments AS c1,
				comments AS c2
			WHERE
				c1.post_id = ?
				AND c2.post_id = ?
				AND c1.tree_left_key <= c2.tree_left_key
				AND c1.tree_right_key >= c2.tree_right_key
				AND c2.deleted = 0
			GROUP BY c1.id
		) AS cc ON c.id = cc.id
		ORDER BY c.tree_left_key`

	rows, err := cr.DB.Query(query, articleId, articleId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := models.CommentList{}

	for rows.Next() {
		comment := models.Comment{}
		err = rows.Scan(
			&comment.ID,
			&comment.Name,
			&comment.Email,
			&comment.Website,
			&comment.Gender,
			&comment.CommentatorID,
			&comment.AuthorID,
			&comment.Text,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.RottenLink,
			&comment.AvatarVariant,
			&comment.Deleted)

		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

// GetCollectionWithExtraDataByArticleId TODO use query builder with GetCollectionByArticleId
func (cr *CommentRepository) GetCollectionWithExtraDataByArticleId(articleId int) (models.CommentList, error) {
	query := `
		SELECT
			c.id,
			COALESCE(t.name, u.display_name, u.username) AS username,
			COALESCE(t.mail, u.mail) AS email,
			t.website,
			COALESCE(t.gender, u.gender) AS gender,
			c.commentator_id,
			c.user_id,
			c.text,
			c.tree_depth,
			COALESCE(c.force_created_at, c.time_created) AS time_created,
			c.ip_addr,
			COALESCE(gco.country_code, '-') AS country_code,
			ta.user_agent,
			COALESCE(t.rotten_link, 0) AS rotten_link,
			COALESCE(u.avatar_variant, t.avatar_variant) AS avatar_variant,
			c.deleted
		FROM comments AS c
		LEFT JOIN commentators AS t ON c.commentator_id = t.id
		LEFT JOIN users AS u ON c.user_id = u.id
		LEFT JOIN geo_location AS gl ON c.ip_long = gl.ip_long
		LEFT JOIN geo_location_city AS gci ON gl.city_id = gci.id
		LEFT JOIN geo_location_country AS gco ON gci.country_id = gco.id
		LEFT JOIN tracking_agent ta on c.user_agent_id = ta.id
		INNER JOIN (
			SELECT
				c1.id
			FROM
				comments AS c1,
				comments AS c2
			WHERE
				c1.post_id = ?
				AND c2.post_id = ?
				AND c1.tree_left_key <= c2.tree_left_key
				AND c1.tree_right_key >= c2.tree_right_key
				AND c2.deleted = 0
			GROUP BY c1.id
		) AS cc ON c.id = cc.id
		ORDER BY c.tree_left_key`

	rows, err := cr.DB.Query(query, articleId, articleId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := models.CommentList{}

	for rows.Next() {
		comment := models.Comment{}
		err = rows.Scan(
			&comment.ID,
			&comment.Name,
			&comment.Email,
			&comment.Website,
			&comment.Gender,
			&comment.CommentatorID,
			&comment.AuthorID,
			&comment.Text,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.IP,
			&comment.CountryCode,
			&comment.UserAgent,
			&comment.RottenLink,
			&comment.AvatarVariant,
			&comment.Deleted)

		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (cr *CommentRepository) GetMostActiveCommentators() (*models.CommentatorList, error) {
	query := `
		SELECT
			src.cnt,
			COALESCE(t.name, u.display_name, u.username) AS username,
			COALESCE(t.mail, u.mail) AS email,
			t.website,
			COALESCE(t.gender, u.gender) AS gender,
			COALESCE(t.rotten_link, 0) AS rotten_link,
			COALESCE(u.avatar_variant, t.avatar_variant) AS avatar_variant,
			src.commentator_id,
			src.user_id
		FROM (
			SELECT commentator_id,
				c.user_id,
				COUNT(c.id)         AS cnt,
				MAX(c.time_created) AS last_time
			FROM comments AS c
			WHERE
				c.deleted = 0
				AND (c.user_id IS NULL OR c.user_id <> 1)
			GROUP BY commentator_id, user_id) AS src
		LEFT JOIN commentators AS t ON src.commentator_id = t.id
		LEFT JOIN users AS u ON src.user_id = u.id
		ORDER BY src.cnt DESC, src.last_time DESC
		LIMIT 8`

	rows, err := cr.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	commentators := models.CommentatorList{}

	for rows.Next() {
		commentator := &models.Commentator{}
		err = rows.Scan(
			&commentator.CommentsCount,
			&commentator.Name,
			&commentator.Email,
			&commentator.Website,
			&commentator.Gender,
			&commentator.RottenLink,
			&commentator.AvatarVariant,
			&commentator.CommentatorID,
			&commentator.AuthorID)

		if err != nil {
			return nil, err
		}

		commentators = append(commentators, commentator)
	}

	return &commentators, nil
}

func (cr *CommentRepository) FindForGravatar(id int) (*models.CommentatorForGravatar, error) {
	query := `
		SELECT
			c.id,
			c.mail,
			c.fake_email
		FROM commentators AS c
		WHERE (c.id = ?)`

	commentator := models.CommentatorForGravatar{}

	err := cr.DB.QueryRow(query, id).Scan(
		&commentator.ID,
		&commentator.Email,
		&commentator.FakeEmail)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return &commentator, nil
}
