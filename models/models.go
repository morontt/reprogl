package models

import (
	"database/sql"
	"errors"
	"time"
)

var RecordNotFound = errors.New("models: no matching record found")

type Article struct {
	ID            int
	Title         string
	Slug          string
	Text          string
	Description   sql.NullString
	CreatedAt     time.Time
	CommentsCount int
	CategoryName  string
	CategorySlug  string
}
