package models

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Commentator struct {
	Name          string
	Website       sql.NullString
	Email         sql.NullString
	CommentatorID sql.NullInt32
	AuthorID      sql.NullInt32
	CommentsCount int
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

type CommentatorList []*Commentator

func (c *Comment) Avatar() (src string) {
	if c.Deleted {
		src = cdnBaseURL + "/images/avatar/clown.png"

		return
	}

	return c.Commentator.Avatar()
}

func (ctt *Commentator) Avatar() (src string) {
	if ctt.AuthorID.Valid {
		hash := md5.New()
		hash.Write([]byte(fmt.Sprintf("avatar%d", ctt.AuthorID.Int32)))
		hashString := fmt.Sprintf("%X", hash.Sum(nil))

		src = cdnBaseURL + "/images/avatar/" + hashString[2:8] + ".png"
	} else {
		src = ctt.gravatar()
	}

	return
}

func (ctt *Commentator) gravatar() string {
	defaults := make(map[int32]string)
	defaults[0] = "wavatar"
	defaults[1] = "monsterid"

	var idx int32
	if ctt.CommentatorID.Valid {
		idx = ctt.CommentatorID.Int32 % 2
	}

	return fmt.Sprintf("//www.gravatar.com/avatar/%s?s=80&d=%s", ctt.gravatarHash(), defaults[idx])
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
	hash.Write([]byte(strings.ToLower(strings.TrimSpace(s))))

	return hash.Sum(nil)
}
