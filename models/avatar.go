package models

import (
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/utils/hashid"
)

func avatarLink(id int, options hashid.Option) string {
	return container.GetConfig().CDNBaseURL + "/images/avatar/" + hashid.Encode(id, options) + ".png"
}
