//go:build prod

package container

const (
	DefaultEsiTTL = 3600 * 24 * 7
	devMode       = false
)
