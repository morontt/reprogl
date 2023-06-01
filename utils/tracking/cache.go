package tracking

import (
	"sync"
	"time"

	"xelbot.com/reprogl/container"
)

var cache map[string]int64
var cacheLocker sync.RWMutex
var cleanUpAfter int64

func init() {
	cache = make(map[string]int64)
	cleanUpAfter = time.Now().Add(container.CleanUpInterval).Unix()
}

// returns true if the cache item exists and has not expired
func testItem(id string) bool {
	cacheLocker.RLock()
	defer cacheLocker.RUnlock()

	expiration, found := cache[id]
	if found {
		result := time.Now().Unix() < expiration

		return result
	}

	return false
}

func setItem(id string) {
	cacheLocker.Lock()
	defer cacheLocker.Unlock()

	expiration := time.Now().Add(container.TrackExpiration).Unix()
	cache[id] = expiration

	now := time.Now().Unix()
	if now > cleanUpAfter {
		for key, expirationItem := range cache {
			if now > expirationItem {
				delete(cache, key)
			}
		}
		cleanUpAfter = time.Now().Add(container.CleanUpInterval).Unix()
	}
}
