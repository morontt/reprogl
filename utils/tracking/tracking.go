package tracking

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	trackmodels "xelbot.com/reprogl/utils/tracking/models"
)

func CreateActivity(req *http.Request) *trackmodels.Activity {
	ip := net.ParseIP(container.RealRemoteAddress(req))
	if ip == nil {
		return nil
	}

	return &trackmodels.Activity{
		Time:         time.Now(),
		IsCDN:        container.IsCDN(req),
		Addr:         ip,
		UserAgent:    req.UserAgent(),
		RequestedURI: req.URL.RequestURI(),
	}
}

func SaveActivity(activity *trackmodels.Activity, app *container.Application) {
	var userAgentId int
	if strings.HasPrefix(activity.RequestedURI, "/_fragment/") {
		return
	}

	agentHash := userAgentHash(activity.UserAgent)

	repo := repositories.TrackingRepository{DB: app.DB}
	agent, err := repo.GetAgentByHash(agentHash)
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

	err = repo.SaveTracking(activity, userAgentId)
	if err != nil {
		app.LogError(err)
		return
	}
}

func userAgentHash(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
