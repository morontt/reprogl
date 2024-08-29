package repositories

import (
	"database/sql"
	"errors"
	"time"

	"xelbot.com/reprogl/models"
)

type EmailSubscriptionRepository struct {
	DB *sql.DB
}

func (es *EmailSubscriptionRepository) Find(id int) (*models.EmailSubscription, error) {
	query := `
		SELECT
			s.id,
			s.email,
			s.block_sending
		FROM subscription_settings AS s
		WHERE (s.id = ?)`

	model := models.EmailSubscription{}
	err := es.DB.QueryRow(query, id).Scan(
		&model.ID,
		&model.Email,
		&model.BlockSending)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return &model, nil
}

func (es *EmailSubscriptionRepository) Unsubscribe(id int) error {
	query := `
		UPDATE
			subscription_settings
		SET
			block_sending = 1,
			last_update = ?
		WHERE
			id = ?`

	_, err := es.DB.Exec(
		query,
		time.Now().Format("2006-01-02 15:04:05.000"),
		id,
	)

	return err
}
