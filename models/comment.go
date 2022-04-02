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
	Name          string
	Website       sql.NullString
	Email         sql.NullString
	CommentatorID sql.NullInt32
	AuthorID      sql.NullInt32
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

func (c *Comment) Avatar() (src string) {
	if c.Deleted {
		src = cdnBaseURL + "/images/avatar/clown.png"

		return
	}

	if c.AuthorID.Valid {
		hash := md5.New()
		_, _ = io.WriteString(hash, fmt.Sprintf("avatar%d", c.AuthorID.Int32))
		hashString := fmt.Sprintf("%X", hash.Sum(nil))

		src = cdnBaseURL + "/images/avatar/" + hashString[2:8] + ".png"
	} else {
		src = c.gravatar()
	}

	return
}

func (c *Comment) gravatar() string {
	defaults := make(map[int32]string)
	defaults[0] = "wavatar"
	defaults[1] = "monsterid"

	var idx int32
	if c.CommentatorID.Valid {
		idx = c.CommentatorID.Int32 % 2
	}

	return fmt.Sprintf("//www.gravatar.com/avatar/%s?s=80&d=%s", c.gravatarHash(), defaults[idx])
}

func (ctt *Commentator) gravatarHash() string {
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
