package views

import (
	"github.com/CloudyKit/jet"
	"xelbot.com/reprogl/config"
)

var ViewSet *jet.Set

func LoadViewSet() {
	ViewSet = jet.NewHTMLSet("./templates")

	cfg := config.Get()
	ViewSet.SetDevelopmentMode(cfg.DevMode)
}
