package repositories

import (
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	trackmodels "xelbot.com/reprogl/utils/tracking/models"
)

type TrackingRepository struct {
	DB *sql.DB
}

func (tr *TrackingRepository) GetAgentByHash(hash string) (*models.TrackingAgent, error) {
	query := `
		SELECT
			ta.id,
			ta.user_agent,
			ta.hash,
			ta.is_bot,
			ta.created_at
		FROM tracking_agent AS ta
		WHERE (ta.hash = ?)`

	agent := &models.TrackingAgent{}

	err := tr.DB.QueryRow(query, hash).Scan(
		&agent.ID,
		&agent.UserAgent,
		&agent.Hash,
		&agent.IsBot,
		&agent.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return agent, nil
}

func (tr *TrackingRepository) SaveTrackingAgent(activity *trackmodels.Activity) (int, error) {
	query := `INSERT INTO tracking_agent
        (user_agent, hash, is_bot, created_at)
        VALUES(?, ?, ?, ?)`

	result, err := tr.DB.Exec(
		query,
		activity.UserAgent,
		container.MD5(activity.UserAgent),
		activity.IsBot(),
		activity.Time,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (tr *TrackingRepository) SaveTracking(activity *trackmodels.Activity, agentId, articleId int) error {
	data := goqu.Record{
		"ip_addr":     activity.Addr.String(),
		"is_cdn":      activity.IsCDN,
		"status_code": activity.Status,
		"duration":    activity.Duration.Microseconds(),

		"time_created": activity.Time.Format("2006-01-02 15:04:05.000"),
	}

	if agentId > 0 {
		data["user_agent_id"] = agentId
	}
	if articleId > 0 {
		data["post_id"] = articleId
	} else {
		data["request_uri"] = activity.RequestedURI
	}

	if activity.LocationID > 0 {
		data["ip_long"] = activity.LocationID
	}

	if activity.Method != "GET" {
		data["method"] = activity.Method
	}

	ds := goqu.Dialect("mysql8").Insert("tracking").Rows(data)

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}

	_, err = tr.DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
