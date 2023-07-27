package style

func GenerateInfoStyles() string {
	style := "<style>\n"
	style += glyphiconsFont() + "\n"
	style += infoBackground() + "\n"
	style += "    </style>"

	return cdnReplace(style)
}

func infoBackground() string {
	return `      .big-header-container .main-header{background-size:cover;background-position:50% 50%;background-repeat:no-repeat;background-color:#23222d;background-image:url(%cdn%/images/tractor.jpg)}
      @supports (background-image:url(%cdn%/images/tractor.webp)){.big-header-container .main-header{background-image:url(%cdn%/images/tractor.webp)}}
      @media only screen and (max-width:752px){.big-header-container .main-header{background-image:url(%cdn%/images/tractor-752.jpg)}
      @supports (background-image:url(%cdn%/images/tractor-752.webp)){.big-header-container .main-header{background-image:url(%cdn%/images/tractor-752.webp)}}}`
}
