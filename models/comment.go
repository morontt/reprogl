package models

import (
	"database/sql"
	"time"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/utils/hashid"
)

type Commentator struct {
	Name          string
	Website       sql.NullString
	Email         sql.NullString
	CommentatorID sql.NullInt32
	AuthorID      sql.NullInt32
	CommentsCount int
	Gender        int
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
	var id int
	var options hashid.Option

	switch {
	case ctt.CommentatorID.Valid:
		id = int(ctt.CommentatorID.Int32)
		options = hashid.Commentator
	case ctt.AuthorID.Valid:
		id = int(ctt.AuthorID.Int32)
		options = hashid.User
	}

	if ctt.Gender == 1 {
		options |= hashid.Male
	} else {
		options |= hashid.Female
	}

	return container.GetConfig().CDNBaseURL + "/images/avatar/" + hashid.Encode(id, options) + ".png"
}
