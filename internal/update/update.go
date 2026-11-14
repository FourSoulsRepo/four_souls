// Package update finds this build's archive in the latest GitHub release,
// checks its SHA-256 and installs it with a helper script (ADR 029).
//
// It is ported from the owner's ytmemchat_wails (pkg/update) and has no
// Wails import; the app's bindings drive it.
package update

import (
	"errors"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/version"
)

// maxAssetBytes is the download limit for one archive.
const maxAssetBytes = 500 << 20

// sumsName is the checksum file of every release (18.22).
const sumsName = "SHA256SUMS.txt"

var (
	// ErrNoAsset is returned when this OS, arch and edition have no archive.
	ErrNoAsset = errors.New("update: no archive for this system")
	// ErrNoSums is returned when the release has no SHA256SUMS.txt.
	ErrNoSums = errors.New("update: the release has no " + sumsName)
	// ErrChecksum is returned when the archive's SHA-256 is missing or wrong.
	ErrChecksum = errors.New("update: checksum mismatch")
	// ErrTooLarge is returned when the archive is over the size limit.
	ErrTooLarge = errors.New("update: archive too large")
	// ErrBadArchive is returned when the archive has a path that leaves its
	// folder, or no app in it.
	ErrBadArchive = errors.New("update: bad archive")
	// ErrNotReady is returned when Download runs before a Check that found
	// an update.
	ErrNotReady = errors.New("update: no update to download; check first")
	// ErrNoDownload is returned when Apply runs before a verified download.
	ErrNoDownload = errors.New("update: download the update first")
)

// Reason says why a checked release cannot be installed.
type Reason string

// Reasons, in the order Check tests them.
const (
	// ReasonNone: the update can be downloaded and installed.
	ReasonNone Reason = ""
	// ReasonNotVersion: this build or the release is not X.Y.Z (dev builds).
	ReasonNotVersion Reason = "not_version"
	// ReasonNotNewer: the release is not newer than this build.
	ReasonNotNewer Reason = "not_newer"
	// ReasonNoArchive: the release has no archive for this target and edition.
	ReasonNoArchive Reason = "no_archive"
	// ReasonNotWritable: the install folder cannot be written.
	ReasonNotWritable Reason = "not_writable"
)

// Config names the app and the build being updated.
type Config struct {
	// Owner is the GitHub user or organization: FourSoulsRepo.
	Owner string
	// Repo is the GitHub repository: four_souls.
	Repo string
	// Name is the app's file stem: Name, Name.exe or Name.app.
	Name string
	// Current is this build's version (version.App).
	Current string
	// Edition is version.EditionFull or version.EditionForPoor.
	Edition string
	// GOOS and GOARCH pick the archive; empty uses runtime's.
	GOOS   string
	GOARCH string
}

// Status is the result of one Check.
type Status struct {
	// Current is this build's version.
	Current string `json:"current"`
	// Latest is the release tag, e.g. v0.5.0.
	Latest string `json:"latest"`
	// Notes is the release text, cut to 2000 letters.
	Notes string `json:"notes"`
	// URL is the release's page.
	URL string `json:"url"`
	// Releases is the page of all releases, for a manual update.
	Releases string `json:"releases"`
	// Asset is this build's archive name; empty without one.
	Asset string `json:"asset"`
	// Newer is true when Latest is greater than Current.
	Newer bool `json:"newer"`
	// CanInstall is true when Reason is ReasonNone.
	CanInstall bool `json:"canInstall"`
	// Reason says why it cannot install.
	Reason Reason `json:"reason"`
}

// Progress gets bytes done and the total; total is 0 when unknown.
// Download calls it a few times a second and once at the end.
type Progress func(done, total int64)

// Client checks, downloads and applies one update.
// Its methods are safe to call from several goroutines; they run one at a time.
type Client struct {
	// HTTP makes the requests.
	HTTP *http.Client
	// API is the GitHub API root; tests point it at a fake.
	API string
	// Executable is this install; empty uses os.Executable.
	Executable string
	// TempDir holds a new folder per download; empty uses the OS's.
	TempDir string
	// Start runs the helper script; tests replace it.
	Start func(name string, args ...string) error
	// ProgressEvery is the least time between two Progress calls.
	ProgressEvery time.Duration

	cfg Config

	mu       sync.Mutex
	last     Status
	assetURL string
	sumsURL  string
	payload  string
	workDir  string
}

// New makes a client for the public GitHub API.
func New(cfg Config) *Client {
	cfg.Owner = strings.TrimSpace(cfg.Owner)
	cfg.Repo = strings.TrimSpace(cfg.Repo)
	cfg.Name = strings.TrimSpace(cfg.Name)
	if cfg.GOOS == "" {
		cfg.GOOS = runtime.GOOS
	}
	if cfg.GOARCH == "" {
		cfg.GOARCH = runtime.GOARCH
	}
	rt := http.DefaultTransport
	if tr, ok := rt.(*http.Transport); ok {
		t := tr.Clone()
		t.ResponseHeaderTimeout = 30 * time.Second
		rt = t
	}
	return &Client{
		// No overall timeout: a large archive on a slow line takes minutes.
		HTTP:          &http.Client{Transport: rt},
		API:           "https://api.github.com",
		Start:         startDetached,
		ProgressEvery: 200 * time.Millisecond,
		cfg:           cfg,
	}
}

// ReleasesURL is the page of all releases, for a manual update.
func (c *Client) ReleasesURL() string {
	return "https://github.com/" + c.cfg.Owner + "/" + c.cfg.Repo + "/releases"
}

func (c *Client) identity() error {
	cfg := c.cfg
	if cfg.Owner == "" || cfg.Repo == "" || strings.Contains(cfg.Owner, "/") || strings.Contains(cfg.Repo, "/") {
		return errors.New("update: set the GitHub owner and repo")
	}
	if !validName(cfg.Name) {
		return errors.New("update: set an app name without path characters")
	}
	if cfg.Edition != version.EditionFull && cfg.Edition != version.EditionForPoor {
		return errors.New("update: unknown edition " + cfg.Edition)
	}
	return nil
}

func validName(name string) bool {
	if name == "" || name == "." || name == ".." {
		return false
	}
	return !strings.ContainsAny(name, `/\:*?"<>|`)
}

func (c *Client) userAgent() string {
	return c.cfg.Name + "-updater"
}
