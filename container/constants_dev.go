//go:build dev

package container

import "time"

const (
	DefaultEsiTTL = 30
	StatisticsTTL
	RobotsTxtTTL
	StaticFileTTL
	FeedTTL
	AvatarTTL = 300
	devMode   = true

	TrackExpiration time.Duration = time.Minute
	CleanUpInterval time.Duration = 5 * time.Minute
)
