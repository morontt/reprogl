package style

import (
	"strings"

	"xelbot.com/reprogl/container"
)

func cdnReplace(str string) string {
	return strings.Replace(str, "%cdn%", container.GetConfig().CDNBaseURL, -1)
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

func glyphiconsFont() string {
	return `      @font-face{font-family:"Glyphicons Halflings";
      src:url(%cdn%/assets/fonts/glyphicons-halflings-regular.eot);src:url(%cdn%/assets/fonts/glyphicons-halflings-regular.eot?#iefix) format("embedded-opentype"),
      url(%cdn%/assets/fonts/glyphicons-halflings-regular.woff2) format("woff2"),url(%cdn%/assets/fonts/glyphicons-halflings-regular.woff) format("woff"),
      url(%cdn%/assets/fonts/glyphicons-halflings-regular.ttf) format("truetype"),url(%cdn%/assets/fonts/glyphicons-halflings-regular.svg#glyphicons_halflingsregular) format("svg");
      font-display: swap}`
}

func commonStyle() string {
	st := glyphiconsFont()

	return st
}
