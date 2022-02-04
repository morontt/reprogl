package models

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"
)

type Commentator struct {
	Name    string
	Website sql.NullString
	Email   sql.NullString
}

type Comment struct {
	ID        int
	Text      string
	Depth     int
	CreatedAt time.Time
	Deleted   bool
	Commentator
}

type CommentList []*Comment

func (c *Comment) Avatar() string {
	return fmt.Sprintf("//www.gravatar.com/avatar/%s?s=80&d=monsterid", c.GravatarHash())
}

func (ctt *Commentator) GravatarHash() string {
	var bytes []byte

	if ctt.Email.Valid {
		bytes = md5sum(ctt.Email.String)
	} else {
		bytes = md5sum(ctt.Name)
		if ctt.Website.Valid {
			bytes = md5sum(
				fmt.Sprintf("%x%s", bytes, ctt.Website.String),
			)
		}
	}

	return fmt.Sprintf("%x", bytes)
}

func md5sum(s string) []byte {
	hash := md5.New()
	_, _ = io.WriteString(hash, strings.ToLower(strings.TrimSpace(s)))

	return hash.Sum(nil)
}
