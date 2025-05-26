package style

import (
	"strconv"
	"strings"

	"xelbot.com/reprogl/models"
)

func GenerateArticleStyles(article *models.Article, acceptAvif, acceptWebp bool) string {
	style := "<style>\n"

	style += commonStyle() + "\n"
	if article.HasImage() {
		style += styleWithImage(article.FeaturedImage, acceptAvif, acceptWebp)
	} else if article.LjItemID.Valid {
		style += defaultStyleLj() + "\n"
	} else {
		style += defaultStyleWithoutImage() + "\n"
	}
	style += "    </style>"

	return cdnReplace(style)
}

func styleWithImage(image models.FeaturedImage, acceptAvif, acceptWebp bool) string {
	var cssRules string

	srcSet := image.DecodeSrcSet()
	if acceptAvif && image.HasAvif() {
		cssRules = stylesImagesByFormat(srcSet, "avif")
	} else if acceptWebp && image.HasWebp() {
		cssRules = stylesImagesByFormat(srcSet, "webp")
	} else {
		cssRules = stylesImagesByFormat(srcSet, "origin")
		if image.HasWebp() {
			srcSetItem, _ := srcSet["webp"]
			cssRules += "      @supports (background-image:url(%cdn%/uploads/"
			cssRules += srcSetItem.Items[0].Path + ")){\n"
			cssRules += strings.TrimRight(stylesImagesByFormat(srcSet, "webp"), "\n") + "}\n"
		}
	}

	return cssRules
}

func stylesImagesByFormat(srcSet map[string]models.SrcSetItem, format string) string {
	var cssRules string
	var firstImage = true

	if srcSetItem, found := srcSet[format]; found {
		for _, srcImage := range srcSetItem.Items {
			if firstImage {
				cssRules += "      .post-view .post-view-sidebar{background-image:url(%cdn%/uploads/"
				cssRules += srcImage.Path + ")}\n"
				firstImage = false
			} else {
				cssRules += "      @media only screen and (max-width:" + strconv.Itoa(srcImage.Width) + "px){"
				cssRules += ".post-view .post-view-sidebar{background-image:url(%cdn%/uploads/"
				cssRules += srcImage.Path + ")}}\n"
			}
		}
	}

	return cssRules
}
