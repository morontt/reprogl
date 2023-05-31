package models

import (
	"database/sql"
)

type TrackingAgent struct {
	ID        int
	UserAgent string
	Hash      string
	IsHuman   bool
	CreatedAt sql.NullTime
}
