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
	regexpArticle  = regexp.MustCompile(`^\/article\/(?P<slug>[^/?#]+)`)
	slugIndex      = regexpArticle.SubexpIndex("slug")
	trackingCache  *yetacache.Cache[string, int8]
	regexpVersions = regexp.MustCompile(`(.*\/\d+\.\d+\.)\d+\.\d+$`)
)

func init() {
	trackingCache = yetacache.New[string, int8](container.TrackExpiration, container.CleanUpInterval)
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
		UserAgent:    callousVersions(req.UserAgent()),
		RequestedURI: req.URL.RequestURI(),
		Method:       req.Method,
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

	location, err := findLocationByIP(activity.Addr, app)
	if err == nil {
		activity.LocationID = location.ID
	}

	matches := regexpArticle.FindStringSubmatch(activity.RequestedURI)
	if matches != nil {
		articleRepo := repositories.ArticleRepository{DB: app.DB}
		articleId = articleRepo.GetIdBySlug(matches[slugIndex])
	}

	if !trackingCache.Has(activity.FingerPrint) ||
		activity.Status != http.StatusOK ||
		activity.Method != "GET" {
		err = repo.SaveTracking(activity, userAgentId, articleId)
		if err != nil {
			app.LogError(err)

			ipKey := ipAddrKey(activity.Addr)
			cache := app.GetIntCache()
			cache.Delete(ipKey)

			return
		}

		trackingCache.Set(activity.FingerPrint, 1, yetacache.DefaultTTL)
	}
}

func findLocationByIP(ip net.IP, app *container.Application) (*models.Geolocation, error) {
	ipKey := ipAddrKey(ip)

	cache := app.GetIntCache()
	if locationID, found := cache.Get(ipKey); found {
		app.InfoLog.Printf("[CACHE] IP %s HIT\n", ip.String())

		return &models.Geolocation{IpAddr: ip.String(), ID: locationID}, nil
	} else {
		app.InfoLog.Printf("[CACHE] IP %s MISS\n", ip.String())
	}

	geolocationRepo := repositories.GeolocationRepository{DB: app.DB}
	location, err := geolocationRepo.FindByIP(ip)
	if err != nil {
		return nil, err
	}

	cache.Set(ipKey, location.ID, 168*time.Hour)

	return location, nil
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

func ipAddrKey(ip net.IP) string {
	return "IP_" + ip.String()
}

func callousVersions(agentName string) string {
	parts := strings.Split(agentName, " ")
	for idx, part := range parts {
		parts[idx] = regexpVersions.ReplaceAllString(part, "${1}x.x")
	}

	return strings.Join(parts, " ")
}
