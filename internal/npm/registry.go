package npm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type PackageInfo struct {
	Name     string
	Latest   string
	Versions []string
}

type registryResponse struct {
	Name     string            `json:"name"`
	DistTags map[string]string `json:"dist-tags"`
	Versions map[string]any    `json:"versions"`
}

func FetchPackageInfo(name string) (PackageInfo, error) {
	if name == "" {
		return PackageInfo{}, fmt.Errorf("package name is required")
	}

	endpoint := "https://registry.npmjs.org/" + url.PathEscape(name)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return PackageInfo{}, err
	}
	req.Header.Set("Accept", "application/vnd.npm.install-v1+json")

	resp, err := client.Do(req)
	if err != nil {
		return PackageInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return PackageInfo{}, ErrPackageNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return PackageInfo{}, fmt.Errorf("npm registry error: %s", strings.TrimSpace(string(body)))
	}

	var payload registryResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&payload); err != nil {
		return PackageInfo{}, err
	}

	latest := ""
	if payload.DistTags != nil {
		latest = payload.DistTags["latest"]
	}

	versions := make([]string, 0, len(payload.Versions))
	for version := range payload.Versions {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	return PackageInfo{
		Name:     payload.Name,
		Latest:   latest,
		Versions: versions,
	}, nil
}
