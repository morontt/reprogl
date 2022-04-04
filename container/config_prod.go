//go:build prod

package container

const (
	DefaultEsiTTL = 3600 * 24 * 7
)

func GetBuildTag() string {
	return "prod"
}
