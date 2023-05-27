//go:build prod

package container

const (
	DefaultEsiTTL = 3600 * 24 * 7
	StatisticsTTL = 3600 * 8
	RobotsTxtTTL  = 3600 * 24 * 30
	devMode       = false
)
