package views

import (
	"github.com/CloudyKit/jet"
	"reprogl/config"
)

var ViewSet *jet.Set

func LoadViewSet() {
	ViewSet = jet.NewHTMLSet("./templates")

	cfg := config.Get()
	ViewSet.SetDevelopmentMode(cfg.DevMode)
}
