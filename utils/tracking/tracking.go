package tracking

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/xelbot/yetacache"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	trackmodels "xelbot.com/reprogl/utils/tracking/models"
)

var (
	regexpArticle = regexp.MustCompile(`^\/article\/(?P<slug>[^/?#]+)`)
	slugIndex     = regexpArticle.SubexpIndex("slug")
	cache         *yetacache.Cache
)

func init() {
	cache = yetacache.New(container.TrackExpiration, container.CleanUpInterval)
}

func CreateActivity(req *http.Request) *trackmodels.Activity {
	ip := net.ParseIP(container.RealRemoteAddress(req))
	if ip == nil {
		return nil
	}

	activity := &trackmodels.Activity{
		Time:         time.Now(),
		IsCDN:        container.IsCDN(req),
		Addr:         ip,
		UserAgent:    req.UserAgent(),
		RequestedURI: req.URL.RequestURI(),
	}

	setupBrowserPassiveFingerprint(req, activity)

	return activity
}

func SaveActivity(activity *trackmodels.Activity, app *container.Application) {
	var (
		userAgentId, articleId int
	)

	if strings.HasPrefix(activity.RequestedURI, "/_fragment/") && activity.Status == http.StatusOK {
		return
	}

	repo := repositories.TrackingRepository{DB: app.DB}
	agent, err := repo.GetAgentByHash(container.MD5(activity.UserAgent))
	if err != nil {
		if errors.Is(err, models.RecordNotFound) {
			userAgentId, err = repo.SaveTrackingAgent(activity)
			if err != nil {
				app.LogError(err)
				return
			}
		} else {
			app.LogError(err)
			return
		}
	} else {
		userAgentId = agent.ID
	}

	geolocationRepo := repositories.GeolocationRepository{DB: app.DB}
	location, err := geolocationRepo.FindByIP(activity.Addr)
	if err == nil && location != nil {
		activity.LocationID = location.ID
	}

	matches := regexpArticle.FindStringSubmatch(activity.RequestedURI)
	if matches != nil {
		articleRepo := repositories.ArticleRepository{DB: app.DB}
		articleId = articleRepo.GetIdBySlug(matches[slugIndex])
	}

	if !cache.TestItem(activity.FingerPrint) {
		err = repo.SaveTracking(activity, userAgentId, articleId)
		if err != nil {
			app.LogError(err)
			return
		}

		cache.SetItem(activity.FingerPrint)
	}
}

func setupBrowserPassiveFingerprint(req *http.Request, a *trackmodels.Activity) {
	a.FingerPrint = container.MD5(
		fmt.Sprintf(
			"%s:%s:%s:%d",
			a.UserAgent,
			a.Addr.String(),
			req.URL.RequestURI(),
			a.Status,
		),
	)
}
