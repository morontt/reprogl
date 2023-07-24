package style

func GenerateInfoStyles() string {
	style := "<style>\n"
	style += cdnReplace(glyphiconsFont()) + "\n"
	style += "    </style>"

	return style
}
