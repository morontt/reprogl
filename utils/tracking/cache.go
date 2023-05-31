package tracking

import (
	"sync"
	"time"
)

var cache map[string]int64
var cacheLocker sync.RWMutex

func init() {
	cache = make(map[string]int64)
}

// returns true if the cache item exists and has not expired
func testItem(id string) bool {
	cacheLocker.RLock()
	expiration, found := cache[id]
	if found {
		result := time.Now().Unix() < expiration
		cacheLocker.RUnlock()

		return result
	}

	cacheLocker.RUnlock()

	return false
}

func setItem(id string) {
	expiration := time.Now().Add(time.Hour).Unix()
	cacheLocker.Lock()
	cache[id] = expiration
	cacheLocker.Unlock()
}
