package views

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
)

const RegionalIndicatorOffset = 127397

func rawHTML(s string) template.HTML {
	return template.HTML(s)
}

func urlGenerator(routeName string, pairs ...string) string {
	return container.GenerateURL(routeName, pairs...)
}

func absUrlGenerator(routeName string, pairs ...string) string {
	return container.GenerateAbsoluteURL(routeName, pairs...)
}

func tags(tl models.TagList) template.HTML {
	var s string
	if len(tl) > 0 {
		s = "Теги: "
		links := make([]string, len(tl))
		for i, t := range tl {
			links[i] = fmt.Sprintf(
				"<a href=\"%s\">%s</a>",
				urlGenerator("tag-first", "slug", t.Slug),
				t.Name,
			)
		}

		s += strings.Join(links, ", ")
	}

	return template.HTML(s)
}

func nl2br(s string) string {
	return strings.Replace(s, "\n", "<br/>", -1)
}

func authorName() string {
	return container.GetConfig().Author
}

func authorBio() string {
	return container.GetConfig().AuthorBio
}

func authorGithub() string {
	return fmt.Sprintf("https://github.com/%s", container.GetConfig().AuthorGithub)
}

func authorTelegram() string {
	return fmt.Sprintf("https://t.me/%s/", container.GetConfig().AuthorTelegram)
}

func renderESI(routeName string, pairs ...string) template.HTML {
	s := fmt.Sprintf(
		"<esi:include src=\"%s\" onerror=\"continue\"/>",
		urlGenerator(routeName, pairs...),
	)

	return template.HTML(s)
}

func subString(input string, length int) (str string) {
	symbols := []rune(input)

	if len(symbols) >= length {
		str = string(symbols[:length-3]) + "..."
	} else {
		str = input
	}

	return
}

func timeTag(t time.Time) template.HTML {
	var s string = "<time class=\"post-date\" datetime=\"" +
		t.Format(time.RFC3339) + "\">" + t.Format("2 Jan 2006, 15:04:05.000") +
		"</time>"

	return template.HTML(s)
}

func goVersion() string {
	return container.GoVersionNumbers
}

func commentsCountString(cnt int) (str string) {
	modulo := cnt % 10
	if modulo == 1 {
		str = fmt.Sprintf("%d комментарий", cnt)
	}

	if modulo > 1 && modulo < 5 {
		str = fmt.Sprintf("%d комментария", cnt)
	}

	if modulo > 4 || modulo == 0 {
		str = fmt.Sprintf("%d комментариев", cnt)
	}

	modulo100 := cnt % 100
	if modulo100 >= 11 && modulo100 <= 14 {
		str = fmt.Sprintf("%d комментариев", cnt)
	}

	return
}

func timesCountString(cnt int) (str string) {
	modulo := cnt % 10
	if modulo == 1 {
		str = fmt.Sprintf("%d раз", cnt)
	}

	if modulo > 1 && modulo < 5 {
		str = fmt.Sprintf("%d раза", cnt)
	}

	if modulo > 4 || modulo == 0 {
		str = fmt.Sprintf("%d раз", cnt)
	}

	modulo100 := cnt % 100
	if modulo100 >= 12 && modulo100 <= 14 {
		str = fmt.Sprintf("%d раз", cnt)
	}

	return
}

func flagCounterImage(fullSize bool) func() template.HTML {
	var (
		url  string
		w, h int
	)

	w = 162
	h = 82
	if container.IsDevMode() {
		url = "/images/flagcounter.png"
		if !fullSize {
			url = "/images/flagcounter_mini.png"
			w = 160
			h = 20
		}
	} else {
		url = "//s05.flagcounter.com/count2/D9g3/bg_23222D/txt_FFFFFF/border_FFFFFF/columns_2/maxflags_4/viewers_3/labels_0/pageviews_1/flags_0/percent_1/"
		if !fullSize {
			url = "//s05.flagcounter.com/mini/D9g3/bg_23222D/txt_FFFFFF/border_23222D/flags_0/"
			w = 160
			h = 20
		}
	}

	return func() template.HTML {
		return template.HTML(
			fmt.Sprintf("<img src=\"%s\" alt=\"Free counters!\" width=\"%d\" height=\"%d\">", url, w, h),
		)
	}
}

func emojiFlag(countryCode string) string {
	if countryCode == "-" {
		// https://apps.timwhitlock.info/unicode/inspect?s=%F0%9F%8F%B4%E2%80%8D%E2%98%A0%EF%B8%8F
		return string([]rune{'\U0001F3F4', '\u200D', '\u2620', '\uFE0F'})
	}

	if len(countryCode) != 2 {
		return countryCode
	}

	countryCode = strings.ToUpper(countryCode)

	resultBytes := make([]rune, 0, 2)
	for _, b := range []byte(countryCode) {
		resultBytes = append(resultBytes, rune(RegionalIndicatorOffset+int(b)))
	}

	return string(resultBytes)
}
