package style

import (
	"strings"

	"xelbot.com/reprogl/container"
)

func GenerateStatisticsStyles() string {
	style := "<style>\n"
	style += strings.Replace(defaultStyleWithoutImage(), "%cdn%", container.GetConfig().CDNBaseURL, -1) + "\n"
	style += "    </style>"

	return style
}
