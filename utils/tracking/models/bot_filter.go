package models

import (
	"regexp"
	"strings"
)

var engines = []string{
	"bot",
	"spider",
	"crawler",
	"aport",
	"rambler",
	"yahoo",
	"accoona",
	"aspseek",
	"lycos",
	"scooter",
	"altavista",
	"webalta",
	"estyle",
	"scrubby",
	"ia_archiver",
	"eltaindexer",
	"ezooms",
	"panscient",
	"solomono",
	"sistrix",
	"larbin",
	"infohelfer",
	"nutch",
	"seznam",
	"search.goo",
	"semager",
	"pixray",
	"findlinks",
	"seokicks",
	"apache-httpclient",
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
	"dreamwidth.org",
	"liferea",
	"tt-rss.org",
	"datanyze",
	"flipboardrss",
	"crawlson",
	"barkrowler",
	"pandalytics",
	"qwantify",
	"bubing",
	"appengine-google",
	"censysinspect",
	"cpp-httplib",
	"leakix.net",
	"http.rb/",
	"headless",
	"scrap",
	"axios/",
	"googleother",
	"friendica",
	"pleroma",
	"akkoma",
	"flipboardproxy",
	"anyevent-http",
	"faraday",
	"hackney",
	"bw/",
	"windowspowershell/",
	"perplexity",
	"seolyt",
	"okhttp/",
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

	if (strings.Contains(agent, "facebook") && strings.Contains(agent, "externalhit")) ||
		(strings.Contains(agent, "owler") && strings.Contains(agent, "ows.eu")) {
		return true
	}

	return false
}
