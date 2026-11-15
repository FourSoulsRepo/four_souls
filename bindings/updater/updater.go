// Package updater is the binding of the Update screen (18.24, ADR 029,
// ADR 030).
package updater

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	"github.com/FourSoulsRepo/four_souls/internal/update"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// progressEvent is the Wails event Download sends while it works; the
// Update screen draws its bar from it (18.24, ADR 029).
const progressEvent = "update:progress"

// Update errors the player reads in their own language (ADR 021). The
// cause stays in the log.
var (
	errUpdateCheck      = engine.NewKey("error.update_check")
	errUpdateDownload   = engine.NewKey("error.update_download")
	errUpdateNoArchive  = engine.NewKey("error.update_no_archive")
	errUpdateNoSums     = engine.NewKey("error.update_no_sums")
	errUpdateChecksum   = engine.NewKey("error.update_checksum")
	errUpdateTooLarge   = engine.NewKey("error.update_too_large")
	errUpdateBadArchive = engine.NewKey("error.update_bad_archive")
	errUpdateNotReady   = engine.NewKey("error.update_not_ready")
	errUpdateNoDownload = engine.NewKey("error.update_no_download")
	errUpdateInstall    = engine.NewKey("error.update_install")
)

// UpdateProgress is the payload of progressEvent: bytes done and the
// total, 0 when the server did not say.
type UpdateProgress struct {
	Done  int64 `json:"done"`
	Total int64 `json:"total"`
}

// api is what the binding needs from the app (ADR 030, item 3).
type api interface {
	// Context is the Wails context; it ends when the app shuts down,
	// which stops a running check or download.
	Context() context.Context
}

// Updater holds the update bindings of the Update screen (18.24, ADR 029).
type Updater struct {
	api    api
	client *update.Client

	// The Wails runtime calls; tests replace them.
	emit func(ctx context.Context, name string, data any)
	quit func(ctx context.Context)
	open func(ctx context.Context, url string)

	mu  sync.Mutex
	url string // the release page of the last Check
}

// New makes the bindings for this build's releases.
func New(a api) *Updater {
	return &Updater{
		api: a,
		client: update.New(update.Config{
			Owner:   "FourSoulsRepo",
			Repo:    "four_souls",
			Name:    "FourSouls",
			Current: version.App,
			Edition: version.Edition,
		}),
		emit: func(ctx context.Context, name string, data any) { runtime.EventsEmit(ctx, name, data) },
		quit: runtime.Quit,
		open: runtime.BrowserOpenURL,
	}
}

// Check asks GitHub for the latest release. "No update" is in the
// status, with its reason; an error means GitHub could not be asked.
func (u *Updater) Check() (update.Status, error) {
	st, err := u.client.Check(u.api.Context())
	if err != nil {
		return st, updateError(err, errUpdateCheck)
	}
	u.mu.Lock()
	u.url = st.URL
	u.mu.Unlock()
	return st, nil
}

// Download fetches and checks the archive of the last Check, sending
// progressEvent as it goes.
func (u *Updater) Download() error {
	ctx := u.api.Context()
	err := u.client.Download(ctx, func(done, total int64) {
		u.emit(ctx, progressEvent, UpdateProgress{Done: done, Total: total})
	})
	return updateError(err, errUpdateDownload)
}

// ApplyAndQuit starts the helper that installs the download, then quits
// the app so the helper can swap it.
func (u *Updater) ApplyAndQuit() error {
	if err := u.client.Apply(); err != nil {
		return updateError(err, errUpdateInstall)
	}
	u.quit(u.api.Context())
	return nil
}

// OpenReleases opens the release page in the browser: the last
// checked release, else the page of all releases.
func (u *Updater) OpenReleases() {
	u.mu.Lock()
	url := u.url
	u.mu.Unlock()
	if url == "" {
		url = u.client.ReleasesURL()
	}
	u.open(u.api.Context(), url)
}

// updateError gives an update error the player's words; other is the
// line for an error the package does not name (network, disk).
func updateError(err error, other engine.Key) error {
	if err == nil {
		return nil
	}
	log.Print(err)
	key := other
	for _, k := range []struct {
		err error
		key engine.Key
	}{
		{update.ErrNoAsset, errUpdateNoArchive},
		{update.ErrNoSums, errUpdateNoSums},
		{update.ErrChecksum, errUpdateChecksum},
		{update.ErrTooLarge, errUpdateTooLarge},
		{update.ErrBadArchive, errUpdateBadArchive},
		{update.ErrNotReady, errUpdateNotReady},
		{update.ErrNoDownload, errUpdateNoDownload},
	} {
		if errors.Is(err, k.err) {
			key = k.key
			break
		}
	}
	return msgerr.New(key.Msg(), err)
}
