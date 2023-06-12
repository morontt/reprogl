package models

import "strings"

var engines []string = []string{
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
}

func isBot(agent string) bool {
	for _, substring := range engines {
		if strings.Contains(strings.ToLower(agent), substring) {
			return true
		}
	}

	return false
}
