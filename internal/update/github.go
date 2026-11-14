package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

type ghRelease struct {
	TagName string    `json:"tag_name"`
	HTMLURL string    `json:"html_url"`
	Body    string    `json:"body"`
	Assets  []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

// Check reads the latest release and says whether this build can update.
// An error means GitHub could not be asked; "no update" is a Status.
func (c *Client) Check(ctx context.Context) (Status, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.identity(); err != nil {
		return Status{}, err
	}
	c.last, c.assetURL, c.sumsURL, c.payload = Status{}, "", "", ""
	rel, err := c.fetchLatest(ctx)
	if err != nil {
		return Status{}, err
	}
	st := Status{
		Current:  c.cfg.Current,
		Latest:   strings.TrimSpace(rel.TagName),
		Notes:    truncate(rel.Body, 2000),
		URL:      rel.HTMLURL,
		Releases: c.ReleasesURL(),
		Newer:    Newer(rel.TagName, c.cfg.Current),
	}
	if st.URL == "" {
		st.URL = st.Releases
	}
	var asset ghAsset
	if name, err := AssetName(st.Latest, c.cfg.Edition, c.cfg.GOOS, c.cfg.GOARCH); err == nil {
		if a, ok := findAsset(rel.Assets, name); ok {
			asset = a
			st.Asset = a.Name
		}
	}
	switch {
	case !IsVersion(st.Latest) || !IsVersion(st.Current):
		st.Reason = ReasonNotVersion
	case !st.Newer:
		st.Reason = ReasonNotNewer
	case st.Asset == "":
		st.Reason = ReasonNoArchive
	case !c.installWritable():
		st.Reason = ReasonNotWritable
	}
	st.CanInstall = st.Reason == ReasonNone
	if st.CanInstall {
		c.assetURL = asset.URL
		if sums, ok := findAsset(rel.Assets, sumsName); ok {
			c.sumsURL = sums.URL
		}
	}
	c.last = st
	return st, nil
}

func (c *Client) fetchLatest(ctx context.Context) (ghRelease, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	u := strings.TrimRight(c.API, "/") + "/repos/" + c.cfg.Owner + "/" + c.cfg.Repo + "/releases/latest"
	res, err := c.get(ctx, u)
	if err != nil {
		return ghRelease{}, fmt.Errorf("update: github: %w", err)
	}
	defer closeBody(res)
	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return ghRelease{}, fmt.Errorf("update: github: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return ghRelease{}, fmt.Errorf("update: github HTTP %d", res.StatusCode)
	}
	var rel ghRelease
	if err := json.Unmarshal(body, &rel); err != nil {
		return ghRelease{}, fmt.Errorf("update: github json: %w", err)
	}
	if strings.TrimSpace(rel.TagName) == "" {
		return ghRelease{}, errors.New("update: the release has no tag")
	}
	return rel, nil
}

func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", c.userAgent())
	hc := c.HTTP
	if hc == nil {
		hc = http.DefaultClient
	}
	return hc.Do(req)
}

// closeBody closes a response body after it was read; there is nothing to
// do when that fails.
func closeBody(res *http.Response) {
	_ = res.Body.Close() //nolint:errcheck // read already; nothing to do
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n]) + "…"
}
