package repositories

import (
	"database/sql"
	"errors"
	"time"

	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/utils/database"
)

type EmailSubscriptionRepository struct {
	DB *database.DB
}

func (es *EmailSubscriptionRepository) Find(id int) (*models.EmailSubscription, error) {
	query := `
		SELECT
			s.id,
			s.email,
			s.subs_type,
			s.block_sending
		FROM subscription_settings AS s
		WHERE s.id = ?`

	model := models.EmailSubscription{}
	err := es.DB.QueryRow(query, id).Scan(
		&model.ID,
		&model.Email,
		&model.Type,
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

func (es *EmailSubscriptionRepository) FindOrCreate(email string, subscrType int) (*models.EmailSubscription, error) {
	query := `
		SELECT
			s.id,
			s.email,
			s.subs_type,
			s.block_sending
		FROM subscription_settings AS s
		WHERE
			s.email = ?
			AND s.subs_type = ?`

	model := models.EmailSubscription{}
	err := es.DB.QueryRow(query, email, subscrType).Scan(
		&model.ID,
		&model.Email,
		&model.Type,
		&model.BlockSending)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return es.Create(email, subscrType)
		} else {
			return nil, err
		}
	}

	return &model, nil
}

func (es *EmailSubscriptionRepository) Create(email string, subscrType int) (*models.EmailSubscription, error) {
	query := `INSERT INTO subscription_settings (email, subs_type, block_sending)
				VALUES (?, ?, ?)`

	stmtResult, err := es.DB.Exec(
		query,
		email,
		subscrType,
		0,
	)
	if err != nil {
		return nil, err
	}

	id, err := stmtResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.EmailSubscription{
		ID:    int(id),
		Email: email,
		Type:  subscrType,
	}, nil
}

func (es *EmailSubscriptionRepository) Unsubscribe(id int) error {
	return es.changeBlockSending(id, 1)
}

func (es *EmailSubscriptionRepository) Subscribe(id int) error {
	return es.changeBlockSending(id, 0)
}

func (es *EmailSubscriptionRepository) changeBlockSending(id, value int) error {
	query := `
		UPDATE
			subscription_settings
		SET
			block_sending = ?,
			last_update = ?
		WHERE
			id = ?`

	_, err := es.DB.Exec(
		query,
		value,
		time.Now().Format("2006-01-02 15:04:05.000"),
		id,
	)

	return err
}
