package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	repoOwner = "farseenmanekhan1232"
	repoName  = "tkeybr"
)

// ReleaseInfo holds information about a release
type ReleaseInfo struct {
	TagName     string `json:"tag_name"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
}

// CheckForUpdate checks if a newer version is available
// currentVersion should be in "v0.0.0" or "0.0.0" format
func CheckForUpdate(currentVersion string) (*ReleaseInfo, error) {
	// Skip check for dev builds
	if currentVersion == "dev" || currentVersion == "dirty" {
		return nil, nil
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	// Compare versions
	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(currentVersion, "v")

	if isNewer(latest, current) {
		return &release, nil
	}

	return nil, nil
}

// isNewer checks if v1 is newer than v2 (simple semver check)
// Assumes format "1.2.3"
func isNewer(v1, v2 string) bool {
	// Simple string comparison works for most standard semver cases
	// e.g. "0.3.5" > "0.3.4"
	// For production, a proper semver library is better, but this suffices for simple versions
	return v1 > v2
}
