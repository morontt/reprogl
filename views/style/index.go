package style

func GenerateIndexStyles() string {
	style := "<style>\n"
	style += glyphiconsFont() + "\n"
	style += indexBackground() + "\n"
	style += "    </style>"

	return cdnReplace(style)
}

func indexBackground() string {
	return `      .big-header-container .main-header{background-size:cover;background-position:50% 50%;background-repeat:no-repeat;background-color:#23222d;background-image:url(%cdn%/images/1500x500.jpg)}
      @supports (background-image:url(%cdn%/images/1500x500.webp)){.big-header-container .main-header{background-image:url(%cdn%/images/1500x500.webp)}}
      @media only screen and (max-width:752px){.big-header-container .main-header{background-image:url(%cdn%/images/kravchik.jpg)}
      @supports (background-image:url(%cdn%/images/kravchik.webp)){.big-header-container .main-header{background-image:url(%cdn%/images/kravchik.webp)}}}`
}
