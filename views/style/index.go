package style

func GenerateIndexStyles() string {
	style := "<style>\n"
	style += cdnReplace(glyphiconsFont()) + "\n"
	style += "    </style>"

	return style
}
