package tracking

import (
	"strconv"
	"testing"
)

func TestCallousVersions(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Apache-HttpClient/4.2 (java 1.5)",
			want: "Apache-HttpClient/4.2 (java 1.5)",
		},
		{
			name: "UCWEB/2.0 (Java; U; MIDP-2.0; ru; NokiaC5-00) U2/1.0.0 UCBrowser/9.4.1.377 U2/1.0.0 Mobile UNTRUSTED/1.0 3gpp-gba",
			want: "UCWEB/2.0 (Java; U; MIDP-2.0; ru; NokiaC5-00) U2/1.0.0 UCBrowser/9.4.x.x U2/1.0.0 Mobile UNTRUSTED/1.0 3gpp-gba",
		},
		{
			name: "Python/3.9 aiohttp/3.8.1",
			want: "Python/3.9 aiohttp/3.8.1",
		},
		{
			name: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36",
			want: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.x.x Safari/537.36",
		},
		{
			name: "Scrapy/1.7.3 (+https://scrapy.org)",
			want: "Scrapy/1.7.3 (+https://scrapy.org)",
		},
		{
			name: "",
			want: "",
		},
		{
			name: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4138.884 Safari/537.36",
			want: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.x.x Safari/537.36",
		},
	}

	for idx, item := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			res := callousVersions(item.name)

			if res != item.want {
				t.Errorf("got %s; want %s", res, item.want)
			}
		})
	}
}
