package models

import (
	"regexp"
	"strings"
)

var engines = []string{
	"googlebot",
	"aport",
	"msnbot",
	"rambler",
	"yahoo",
	"abachobot",
	"accoona",
	"acoirobot",
	"aspseek",
	"croccrawler",
	"dumbot",
	"fast-webcrawler",
	"geonabot",
	"gigabot",
	"lycos",
	"msrbot",
	"scooter",
	"altavista",
	"webalta",
	"idbot",
	"estyle",
	"mail.ru_bot",
	"scrubby",
	"yandex.com/bots",
	"yadirectbot",
	"ia_archiver",
	"antabot",
	"baiduspider",
	"eltaindexer",
	"gsa-crawler",
	"mihalismbot",
	"ahrefsbot",
	"ezooms",
	"bingbot",
	"panscient",
	"solomono",
	"mj12bot",
	"exabot",
	"sistrix",
	"unisterbot",
	"turnitinbot",
	"wbsearchbot",
	"larbin",
	"npbot",
	"infohelfer",
	"nutch",
	"careerbot",
	"seznam",
	"mlbot",
	"search.goo",
	"semager",
	"pixray",
	"findlinks",
	"beetlebot",
	"adsbot",
	"job-bot",
	"spider",
	"crawler",
	"seokicks",
	"yacybot",
	"apache-httpclient",
	"femtosearchbot",
	"fetch",
	"curl/",
	"wget/",
	"parser",
	"ruby/",
	"libwww-perl/",
	"guzzlehttp",
	"http_request2",
	"gocolly",
	"go-http-client",
	"petalbot",
	"dreamwidth.org",
	"semrushbot",
	"liferea",
	"zoominfobot",
	"tt-rss.org",
	"datanyze",
	"dataforseobot",
	"dotbot",
	"applebot",
	"flipboardrss",
	"crawlson",
	"barkrowler",
	"borneobot",
	"seekportbot",
	"serpstatbot",
	"aspiegelbot",
	"amazonbot",
	"indeedbot",
	"awario.com/bots",
	"neevabot",
	"linkpadbot",
	"keybot",
	"timpibot",
	"mojeekbot",
	"senutobot",
	"cliqzbot",
	"flfbaldrbot",
	"duckduckbot",
	"redirectbot",
	"inetdex.com/bot",
	"webwikibot",
	"paperlibot",
	"clarabot",
	"discordbot",
	"duckduckgo-favicons-bot",
	"bsbot",
	"twitterbot",
	"telegrambot",
	"pandalytics",
	"2ip bot",
	"qwantify",
	"bot@linkfluence.com",
	"bubing",
	"appengine-google",
	"censysinspect",
	"cpp-httplib",
	"leakix.net",
	"gptbot",
	"http.rb/",
	"headless",
	"scrapy",
	"axios/",
	"archive.org_bot",
}

var (
	regexpJava    = regexp.MustCompile(`^java/\d+\.[\d\._]+$`)
	regexpHotJava = regexp.MustCompile(`^hotjava/\d+\.[\d\._]+/jre[\d\._x]+$`)
)

func isBot(agent string) bool {
	agent = strings.ToLower(agent)
	for _, substring := range engines {
		if strings.Contains(agent, substring) {
			return true
		}
	}

	if strings.Contains(agent, "java") {
		switch {
		case regexpJava.MatchString(agent):
			return true
		case regexpHotJava.MatchString(agent):
			return true
		}
	}

	if strings.Contains(agent, "python") &&
		(strings.Contains(agent, "aiohttp") ||
			strings.Contains(agent, "requests") ||
			strings.Contains(agent, "httpx") ||
			strings.Contains(agent, "urllib")) {

		return true
	}

	return false
}
