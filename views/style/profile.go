package style

func GenerateProfileStyles() string {
	style := "<style>\n"
	style += commonStyle() + "\n"
	style += symbolaFont() + "\n"
	style += defaultStyleWithoutImage() + "\n"
	style += "    </style>"

	return cdnReplace(style)
}

func symbolaFont() string {
	return `      @font-face {font-family: 'Symbola'; src: url('%cdn%/assets/fonts/symbola.eot');
      src: url('%cdn%/assets/fonts/symbola.eot?#iefix') format('embedded-opentype'), url('%cdn%/assets/fonts/symbola.woff2') format('woff2'),
      url('%cdn%/assets/fonts/symbola.woff') format('woff'), url('%cdn%/assets/fonts/symbola.ttf') format('truetype'),
      url('%cdn%/assets/fonts/symbola.svg#Symbola') format('svg'); font-weight: normal; font-style: normal}`
}
