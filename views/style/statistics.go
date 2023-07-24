package style

func GenerateStatisticsStyles() string {
	style := "<style>\n"
	style += glyphiconsFont() + "\n"
	style += defaultStyleWithoutImage() + "\n"
	style += "    </style>"

	return cdnReplace(style)
}
