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
	BlockingTime   = 15 * time.Minute
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

func HandleRequest(hash string, statusCode int) (blocked bool) {
	if len(hash) == 0 {
		return false
	}

	if statusCode < 400 || statusCode >= 500 {
		return false
	}

	hit := time.Now().UnixMilli()
	hitFrom := hit - int64(controlWindow/time.Millisecond)

	hits := []int64{hit}
	cachedData := getOrCreate(hash)
	for _, v := range cachedData.hits {
		if v > hitFrom {
			hits = append(hits, v)
		}
	}

	cachedData.hits = hits
	if len(hits) >= hitLimit {
		blocked = true
		cachedData.blockTo = time.Now().Add(BlockingTime).UnixMilli()
	}

	limiterCache.Set(hash, cachedData, yetacache.DefaultTTL)

	return
}

func IsBlocked(hash string) bool {
	if len(hash) == 0 {
		return false
	}

	if item, ok := limiterCache.Get(hash); ok {
		return item.blockTo > time.Now().UnixMilli()
	}

	return false
}

func getOrCreate(hash string) rateLimitItem {
	if item, ok := limiterCache.Get(hash); ok {
		return item
	}

	return rateLimitItem{
		hits: make([]int64, 0),
	}
}
