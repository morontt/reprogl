package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"xelbot.com/reprogl/container"
)

type Commentator struct {
	Name          string
	Website       sql.NullString
	Email         sql.NullString
	CommentatorID sql.NullInt32
	AuthorID      sql.NullInt32
	CommentsCount int
	ForceImage    bool
}

type Comment struct {
	ID        int
	Text      string
	Depth     int
	CreatedAt time.Time
	Deleted   bool
	Commentator
}

type CommentatorForGravatar struct {
	ID        int
	Email     sql.NullString
	FakeEmail sql.NullBool
}

type CommentList []*Comment

type CommentatorList []*Commentator

func (c *Comment) Avatar() (src string) {
	if c.Deleted {
		src = container.GetConfig().CDNBaseURL + "/images/avatar/clown.png"

		return
	}

	return c.Commentator.Avatar()
}

func (ctt *Commentator) Avatar() (src string) {
	var str string
	switch {
	case ctt.CommentatorID.Valid && !ctt.ForceImage:
		return ctt.gravatar()
	case ctt.CommentatorID.Valid && ctt.ForceImage:
		str = fmt.Sprintf("ratava%d", ctt.CommentatorID.Int32)
	case ctt.AuthorID.Valid:
		str = fmt.Sprintf("avatar%d", ctt.AuthorID.Int32)
	}

	hashString := strings.ToUpper(container.MD5(str))

	return container.GetConfig().CDNBaseURL + "/images/avatar/" + hashString[2:8] + ".png"
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

func (ctt *Commentator) gravatarHash() (hash string) {
	if ctt.Email.Valid {
		hash = md5sum(ctt.Email.String)
	} else {
		hash = md5sum(ctt.Name)
		if ctt.Website.Valid {
			hash = md5sum(
				fmt.Sprintf("%s%s", hash, ctt.Website.String),
			)
		}
	}

	return
}

func md5sum(s string) string {
	return container.MD5(strings.ToLower(strings.TrimSpace(s)))
}
