package models

import (
	"strconv"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/utils/hashid"
)

func AvatarLink(id int, options hashid.Option, size ...int) string {
	s := 80
	if len(size) > 0 {
		s = size[0]
	}

	var postfix string
	if s != 80 {
		postfix = ".w" + strconv.Itoa(s)
	}

	return container.GetConfig().CDNBaseURL + "/images/avatar/" + hashid.Encode(id, options) + postfix + ".png"
}
