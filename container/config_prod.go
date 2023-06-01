//go:build prod

package container

import "time"

const (
	DefaultEsiTTL = 3600 * 24 * 7
	StatisticsTTL = 3600 * 8
	RobotsTxtTTL  = 3600 * 24 * 30
	devMode       = false

	TrackExpiration time.Duration = time.Hour
	CleanUpInterval time.Duration = 24 * time.Hour
)
