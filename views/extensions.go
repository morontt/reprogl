package views

import (
	"fmt"
	"html/template"
	"math/rand"
	"strings"
	"time"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/utils"
	"xelbot.com/reprogl/views/style"
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

func authorBio() string {
	emoij := []rune{
		rune(0x1F41C), // ant
		rune(0x1FAB0), // fly
		rune(0x1F41D), // bee
		rune(0x1F980), // crab
		rune(0x1F997), // cricket
		rune(0x1F577), // spider
		rune(0x1F982), // scorpion
		rune(0x1F990), // shrimp
	}

	return container.GetConfig().Author.Bio + " " + string(emoij[rand.Intn(len(emoij))])
}

func authorDataPart(item string) (str string) {
	author := container.GetConfig().Author

	switch item {
	case "name":
		str = author.FullName
	case "github":
		str = fmt.Sprintf("https://github.com/%s", author.GithubUser)
	case "telegram":
		str = fmt.Sprintf("https://t.me/%s/", author.TelegramChannel)
	case "mastodon":
		str = author.MastodonLink
	case "gitverse":
		str = fmt.Sprintf("https://gitverse.ru/%s", author.GitVerseUser)
	default:
		str = "N/A"
	}

	return
}

func authorLocation() template.HTML {
	location := container.GetConfig().AuthorLocationRu
	s := make([]string, 0, 3)

	if location.City != "" {
		s = append(s, "<span class=\"locality\">"+location.City+"</span>")
	}
	if location.Region != "" {
		s = append(s, "<span class=\"region\">"+location.Region+"</span>")
	}
	if location.Country != "" {
		s = append(s, "<span class=\"country-name\">"+location.Country+"</span>")
	}

	return template.HTML(strings.Join(s, ", "))
}

func authorJob() template.HTML {
	jobs := container.GetConfig().Jobs
	job := jobs.Last()
	s := fmt.Sprintf(
		"<span class=\"title\">%s</span> в <a class=\"org\" href=\"%s\">%s</a>",
		job.Title,
		job.Link,
		job.Company,
	)

	return template.HTML(s)
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
	var s = "<time class=\"post-date\" datetime=\"" +
		t.Format(time.RFC3339) + "\">" +
		t.Format("2 ") +
		utils.RuMonthName(t.Month(), true) +
		t.Format(" 2006, 15:04:05.000") +
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
	cdn := container.GetConfig().CDNBaseURL
	if container.IsDevMode() {
		url = cdn + "/images/flagcounter.png"
		if !fullSize {
			url = cdn + "/images/flagcounter_mini.png"
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

func articleStyles(article *models.Article, acceptAvif, acceptWebp bool) template.HTML {
	return template.HTML(style.GenerateArticleStyles(article, acceptAvif, acceptWebp))
}

func statisticsStyles() template.HTML {
	return template.HTML(style.GenerateStatisticsStyles())
}

func indexStyles() template.HTML {
	return template.HTML(style.GenerateIndexStyles())
}

func infoStyles() template.HTML {
	return template.HTML(style.GenerateInfoStyles())
}

func profileStyles() template.HTML {
	return template.HTML(style.GenerateProfileStyles())
}

func isDev() bool {
	return container.IsDevMode()
}
