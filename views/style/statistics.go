package style

func GenerateStatisticsStyles() string {
	style := "<style>\n"
	style += commonStyle() + "\n"
	style += defaultStyleWithoutImage() + "\n"
	style += "    </style>"

	return cdnReplace(style)
}
