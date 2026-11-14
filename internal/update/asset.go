package update

import (
	"fmt"
	"strings"

	"github.com/FourSoulsRepo/four_souls/internal/version"
)

// Target is the archive's system part, e.g. linux-amd64.tar.gz.
// The names come from release.yml and are a contract (ADR 029).
func Target(goos, goarch string) (string, error) {
	switch goos {
	case "darwin":
		return "macos-universal.zip", nil
	case "windows":
		if goarch == "amd64" || goarch == "arm64" {
			return "windows-" + goarch + ".zip", nil
		}
	case "linux":
		if goarch == "amd64" || goarch == "arm64" {
			return "linux-" + goarch + ".tar.gz", nil
		}
	}
	return "", fmt.Errorf("%w: %s/%s", ErrNoAsset, goos, goarch)
}

// AssetName is the exact archive name of one release for one build:
// four-souls-<latest without v>-<edition prefix><target>.
func AssetName(latest, edition, goos, goarch string) (string, error) {
	target, err := Target(goos, goarch)
	if err != nil {
		return "", err
	}
	prefix := ""
	switch edition {
	case version.EditionFull:
	case version.EditionForPoor:
		prefix = "for-poor-"
	default:
		return "", fmt.Errorf("%w: edition %q", ErrNoAsset, edition)
	}
	ver := strings.TrimPrefix(strings.TrimSpace(latest), "v")
	return "four-souls-" + ver + "-" + prefix + target, nil
}

// findAsset returns the asset with exactly this name; never a substring.
func findAsset(assets []ghAsset, name string) (ghAsset, bool) {
	for _, a := range assets {
		if a.Name == name {
			return a, true
		}
	}
	return ghAsset{}, false
}
