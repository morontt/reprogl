package style

import (
	"strconv"
	"strings"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
)

func GenerateArticleStyles(article *models.Article, acceptAvif, acceptWebp bool) string {
	style := "<style>\n"

	if article.HasImage() {
		style += styleWithImage(article.FeaturedImage, acceptAvif, acceptWebp)
	} else {
		style += strings.Replace(defaultStyleWithoutImage(), "%cdn%", container.GetConfig().CDNBaseURL, -1) + "\n"
	}

	style += "    </style>"

	return style
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
			cssRules += "      @supports (background-image:url(" + container.GetConfig().CDNBaseURL + "/uploads/"
			cssRules += srcSetItem.Items[0].Path + ")){\n"
			cssRules += strings.TrimRight(stylesImagesByFormat(srcSet, "webp"), "\n") + "}\n"
		}
	}

	return cssRules
}

func stylesImagesByFormat(srcSet map[string]models.SrcSetItem, format string) string {
	var cssRules string
	var firstImage = true

	pathPrefix := container.GetConfig().CDNBaseURL + "/uploads/"

	if srcSetItem, found := srcSet[format]; found {
		for _, srcImage := range srcSetItem.Items {
			if firstImage {
				cssRules += "      .post-view .post-view-sidebar{background-image:url("
				cssRules += pathPrefix + srcImage.Path + ")}\n"
				firstImage = false
			} else {
				cssRules += "      @media only screen and (max-width:" + strconv.Itoa(srcImage.Width) + "px){"
				cssRules += ".post-view .post-view-sidebar{background-image:url("
				cssRules += pathPrefix + srcImage.Path + ")}}\n"
			}
		}
	}

	return cssRules
}

func defaultStyleWithoutImage() string {
	return `      .post-view .post-view-sidebar{background-image:url(%cdn%/images/mantis.jpg)}
      @media only screen and (max-width:624px){.post-view .post-view-sidebar{background-image:url(%cdn%/images/mantis/mantis_624w.jpg)}}
      @media only screen and (max-width:448px){.post-view .post-view-sidebar{background-image:url(%cdn%/images/mantis/mantis_448w.jpg)}}
      @media only screen and (max-width:320px){.post-view .post-view-sidebar{background-image:url(%cdn%/images/mantis/mantis_320w.jpg)}}
      @supports (background-image:url(%cdn%/images/mantis/mantis.webp)){
      .post-view .post-view-sidebar{background-image: url(%cdn%/images/mantis/mantis.webp)}
      @media only screen and (max-width:624px){.post-view .post-view-sidebar{background-image:url(%cdn%/images/mantis/mantis_624w.webp)}}
      @media only screen and (max-width:448px){.post-view .post-view-sidebar{background-image:url(%cdn%/images/mantis/mantis_448w.webp)}}
      @media only screen and (max-width:320px){.post-view .post-view-sidebar{background-image:url(%cdn%/images/mantis/mantis_320w.webp)}}}`
}
