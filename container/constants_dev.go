//go:build dev

package container

import "time"

const (
	DefaultEsiTTL = 30
	StatisticsTTL
	RobotsTxtTTL
	FeedTTL
	devMode = true

	TrackExpiration time.Duration = time.Minute
	CleanUpInterval time.Duration = 5 * time.Minute
)
