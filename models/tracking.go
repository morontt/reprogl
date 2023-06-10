package models

import (
	"database/sql"
)

type TrackingAgent struct {
	ID        int
	UserAgent string
	Hash      string
	IsBot     bool
	CreatedAt sql.NullTime
}
