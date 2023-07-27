package repositories

import (
	"database/sql"
	"errors"

	"xelbot.com/reprogl/models"
)

type SystemParametersRepository struct {
	DB *sql.DB
}

func (sp *SystemParametersRepository) FindByKey(key string) (string, error) {
	query := `
		SELECT
			sp.value
		FROM sys_parameters AS sp
		WHERE sp.optionkey = ?`

	var value string

	err := sp.DB.QueryRow(query, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return value, models.RecordNotFound
		} else {
			return value, err
		}
	}

	return value, nil
}
