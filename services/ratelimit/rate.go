package ratelimit

import (
	"fmt"
	"hash/adler32"
	"time"

	"github.com/xelbot/yetacache"
	"xelbot.com/reprogl/container"
	trackmodels "xelbot.com/reprogl/utils/tracking/models"
)

const (
	blockingTime   = 15 * time.Minute
	itemExpiration = 20 * time.Minute
	controlWindow  = 30 * time.Second
	hitLimit       = 7
)

var (
	limiterCache *yetacache.Cache[string, rateLimitItem]
)

func init() {
	limiterCache = yetacache.New[string, rateLimitItem](itemExpiration, container.CleanUpInterval)
}

type rateLimitItem struct {
	blockTo int64
	hits    []int64
}

func ActivityHash(act *trackmodels.Activity) string {
	if act == nil {
		return ""
	}

	if act.IsInternalRequest() {
		return ""
	}

	return fmt.Sprintf("%s_%x", act.Addr.String(), adler32.Checksum([]byte(act.UserAgent)))
}

func IsBlocked(hash string) bool {
	if len(hash) == 0 {
		return false
	}

	return true
}
