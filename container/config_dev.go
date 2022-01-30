//go:build dev
// +build dev

package container

const (
	DefaultEsiTTL = 30
)

func GetBuildTag() string {
	return "dev"
}
