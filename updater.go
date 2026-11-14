package main

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

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

// Updater holds the update bindings of the Update screen (18.24, ADR 029).
type Updater struct {
	client *update.Client

	// The Wails calls; tests replace them.
	emit func(name string, data any)
	quit func()
	open func(url string)

	mu     sync.Mutex
	ctx    context.Context //nolint:containedctx // ends when the app shuts down
	cancel context.CancelFunc
	url    string // the release page of the last Check
}

// NewUpdater makes the bindings for this build's releases.
func NewUpdater() *Updater {
	ctx, cancel := context.WithCancel(context.Background())
	return &Updater{
		client: update.New(update.Config{
			Owner:   "FourSoulsRepo",
			Repo:    "four_souls",
			Name:    "FourSouls",
			Current: version.App,
			Edition: version.Edition,
		}),
		ctx:    ctx,
		cancel: cancel,
		emit:   func(string, any) {},
		quit:   func() {},
		open:   func(string) {},
	}
}

// startup wires the Wails runtime; downloads end when the app shuts down.
func (u *Updater) startup(ctx context.Context) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.cancel()
	u.ctx, u.cancel = context.WithCancel(ctx)
	u.emit = func(name string, data any) { runtime.EventsEmit(ctx, name, data) }
	u.quit = func() { runtime.Quit(ctx) }
	u.open = func(url string) { runtime.BrowserOpenURL(ctx, url) }
}

// shutdown stops a running check or download.
func (u *Updater) shutdown() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.cancel()
}

// wails is the context and the runtime calls, read under the lock.
func (u *Updater) wails() (ctx context.Context, emit func(string, any), quit func(), open func(string)) {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.ctx, u.emit, u.quit, u.open
}

// Check asks GitHub for the latest release. "No update" is in the
// status, with its reason; an error means GitHub could not be asked.
func (u *Updater) Check() (update.Status, error) {
	ctx, _, _, _ := u.wails()
	st, err := u.client.Check(ctx)
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
	ctx, emit, _, _ := u.wails()
	err := u.client.Download(ctx, func(done, total int64) {
		emit(progressEvent, UpdateProgress{Done: done, Total: total})
	})
	return updateError(err, errUpdateDownload)
}

// ApplyAndQuit starts the helper that installs the download, then quits
// the app so the helper can swap it.
func (u *Updater) ApplyAndQuit() error {
	if err := u.client.Apply(); err != nil {
		return updateError(err, errUpdateInstall)
	}
	_, _, quit, _ := u.wails()
	quit()
	return nil
}

// OpenReleases opens the release page in the browser: the last
// checked release, else the page of all releases.
func (u *Updater) OpenReleases() {
	_, _, _, open := u.wails()
	u.mu.Lock()
	url := u.url
	u.mu.Unlock()
	if url == "" {
		url = u.client.ReleasesURL()
	}
	open(url)
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
	return msgError{msg: key.Msg(), cause: err}
}
