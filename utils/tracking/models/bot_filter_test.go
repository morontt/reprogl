package models

import (
	"strconv"
	"testing"
)

func TestBotFilter(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Mozilla/5.0 (Linux; Android 5.1.1; SM-J120F Build/LMY47X; wv) Mobile Safari/537.36 YandexSearch/7.61 YandexSearchBrowser/7.61",
			want: false,
		},
		{
			name: "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106",
			want: true,
		},
		{
			name: "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			want: true,
		},
		{
			name: "Apache-HttpClient/4.2 (java 1.5)",
			want: true,
		},
		{
			name: "UCWEB/2.0 (Java; U; MIDP-2.0; ru; NokiaC5-00) U2/1.0.0 UCBrowser/9.4.1.377 U2/1.0.0 Mobile UNTRUSTED/1.0 3gpp-gba",
			want: false,
		},
		{
			name: "ICE Browser/5.05 (Java 1.4.0; Windows 2000 5.0 x86)",
			want: false,
		},
		{
			name: "Java/1.8.0_144",
			want: true,
		},
		{
			name: "Java/11.0.10",
			want: true,
		},
		{
			name: "HotJava/1.0.1/JRE1.1.x",
			want: true,
		},
		{
			name: "Mozilla/5.0 (compatible; oBot/2.3.1; http://www.xforce-security.com/crawler/)",
			want: true,
		},
		{
			name: "Feedly/1.0 (+http://www.feedly.com/fetcher.html; 2 subscribers; like FeedFetcher-Google)",
			want: true,
		},
		{
			name: "GuzzleHttp/6.2.1 curl/7.29.0 PHP/7.0.27",
			want: true,
		},
		{
			name: "feedfinder/1.371 Python-urllib/1.17 +http://www.aaronsw.com/2002/feedfinder/",
			want: true,
		},
		{
			name: "Python/3.9 aiohttp/3.8.1",
			want: true,
		},
		{
			name: "python-requests/2.0.1 CPython/2.7.3 Linux/3.2.0-41-virtual",
			want: true,
		},
		{
			name: "Mechanize/2.7.5 Ruby/2.3.3p222 (http://github.com/sparklemotion/mechanize/)",
			want: true,
		},
		{
			name: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36 AppEngine-Google; (+http://code.google.com/appengine; appid: s~feedly-nikon3)",
			want: true,
		},
		{
			name: "Mozilla/5.0 (Linux; Android 5.0) AppleWebKit/537.36 (KHTML, like Gecko) Mobile Safari/537.36 (compatible; Bytespider; spider-feedback@bytedance.com)",
			want: true,
		},
		{
			name: "Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)",
			want: true,
		},
		{
			name: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)",
			want: true,
		},
		{
			name: "http.rb/5.1.1 (Mastodon/4.1.4; +https://mastodon.b12e.be/) Bot",
			want: true,
		},
		{
			name: "Scrapy/1.7.3 (+https://scrapy.org)",
			want: true,
		},
		{
			name: "Mozilla/5.0 (compatible; archive.org_bot +http://archive.org/details/archive.org_bot) Zeno/0569f25 warc/v0.8.33",
			want: true,
		},
	}

	for idx, item := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			res := isBot(item.name)

			if res != item.want {
				t.Errorf("%s : got %t; want %t", item.name, res, item.want)
			}
		})
	}
}
