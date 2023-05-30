package repositories

import (
	"database/sql"
	"errors"
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
			ta.bot_filter,
			ta.created_at
		FROM tracking_agent AS ta
		WHERE (ta.hash = ?)`

	agent := &models.TrackingAgent{}

	err := tr.DB.QueryRow(query, hash).Scan(
		&agent.ID,
		&agent.UserAgent,
		&agent.Hash,
		&agent.IsHuman,
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
        (user_agent, hash, bot_filter, created_at)
        VALUES(?, ?, ?, ?)`

	result, err := tr.DB.Exec(
		query,
		activity.UserAgent,
		activity.AgentHash(),
		1,
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

func (tr *TrackingRepository) SaveTracking(activity *trackmodels.Activity, agentId int) error {
	query := `INSERT INTO tracking
        (user_agent_id, ip_addr, time_created, timestamp_created, is_cdn, request_uri, status_code)
        VALUES(?, ?, ?, ?, ?, ?, ?)`

	_, err := tr.DB.Exec(
		query,
		agentId,
		activity.Addr.String(),
		activity.Time,
		int(activity.Time.Unix()),
		activity.IsCDN,
		activity.RequestedURI,
		activity.Status,
	)
	if err != nil {
		return err
	}

	return nil
}
