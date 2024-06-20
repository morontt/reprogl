package container

import (
	_ "embed"

	"github.com/BurntSushi/toml"
)

type AppData struct {
	HeaderText string `toml:"header_text"`
	Author     Author `toml:"author"`

	AuthorLocationRu Location `toml:"author_location"`
	AuthorLocationEn Location `toml:"author_location_en"`
}

type Author struct {
	FullName string `toml:"fn"`
	Bio      string `toml:"bio"`
	Email    string `toml:"email"`

	GithubUser      string `toml:"github"`
	TelegramChannel string `toml:"telegram"`
	MastodonLink    string `toml:"mastodon"`
}

type Location struct {
	City    string `toml:"city"`
	Region  string `toml:"region"`
	Country string `toml:"country"`
}

var (
	//go:embed app_data.toml
	dataSource string
)

func loadAppData() AppData {
	var data AppData
	_, err := toml.Decode(dataSource, &data)
	if err != nil {
		panic(err)
	}

	return data
}
