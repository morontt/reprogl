package models

import (
	"strconv"
	"testing"
)

func TestBotFilter(t *testing.T) {
	tests := []struct {
		name   string
		result bool
	}{
		{
			name:   "Mozilla/5.0 (Linux; Android 5.1.1; SM-J120F Build/LMY47X; wv) Mobile Safari/537.36 YandexSearch/7.61 YandexSearchBrowser/7.61",
			result: false,
		},
		{
			name:   "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106",
			result: true,
		},
		{
			name:   "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			result: true,
		},
		{
			name:   "Mozilla/5.0 (compatible; Linux x86_64; Mail.RU_Bot/2.0; +https://help.mail.ru/webmaster/indexing/robots)",
			result: true,
		},
	}

	for idx, item := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			res := isBot(item.name)

			if res != item.result {
				t.Errorf("%s : got %t; want %t", item.name, res, item.result)
			}
		})
	}
}
