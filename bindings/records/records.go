// Package records is the binding of match records (ADR 030): it saves
// a finished game's record from the game server (RP-05), and lists,
// opens, exports and deletes saved records (M-01, ADR 031).
package records

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/FourSoulsRepo/four_souls/internal/replay"
	"github.com/FourSoulsRepo/record"
)

// ErrRecordURL means a record URL is not a game server's record link.
var ErrRecordURL = errors.New("not a record link")

// ErrRecordSize means a server sent more than a record can be.
var ErrRecordSize = errors.New("record: the file is too large")

// maxRecord is the largest record the app downloads; records are a few
// hundred KB.
const maxRecord = 64 << 20

// api is what the binding needs from the app (ADR 030, item 3).
type api interface {
	// Context is the Wails context for the save dialog; it ends when
	// the app shuts down, which stops a running download.
	Context() context.Context
}

// Records saves, lists and opens match records.
type Records struct {
	api api
	dir string // the records folder: hosted games' records, and saves by default

	mu   sync.Mutex
	open *replay.Replay // the record Open opened, for the player (10.2)
}

// New makes the binding; dir is the records folder.
func New(a api, dir string) *Records {
	return &Records{api: a, dir: dir}
}

// SaveRecord downloads a finished game's record from the game server
// and asks where to save it (RP-05). link is http://host:port/record?…,
// or https://… for a server behind an HTTPS proxy (8.18). The dialog
// suggests the server's file name, cleaned again, in the records
// folder, so the Records screen lists it (M-01). The result is the
// saved file, or "" when the player cancelled.
func (r *Records) SaveRecord(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Path != "/record" || u.Host == "" {
		return "", ErrRecordURL
	}
	data, name, err := fetchRecord(r.api.Context(), u.String())
	if err != nil {
		return "", err
	}
	if name == "" {
		name = "four-souls-" + record.SafeName(u.Query().Get("game")) + record.Ext
	}
	_ = os.MkdirAll(r.dir, 0o750) //nolint:errcheck // the dialog starts elsewhere without it
	path, err := runtime.SaveFileDialog(r.api.Context(), runtime.SaveDialogOptions{
		Title:            "Save the match record",
		DefaultDirectory: r.dir,
		DefaultFilename:  name,
		Filters:          []runtime.FileFilter{{DisplayName: "Four Souls records", Pattern: "*" + record.Ext}},
	})
	if err != nil || path == "" {
		return "", err //nolint:wrapcheck // shown to the player as is
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("record: %w", err)
	}
	return path, nil
}

// fetchRecord downloads what link serves and the file name it suggests,
// cleaned by record.SafeFileName; "" when it suggests none.
func fetchRecord(ctx context.Context, link string) ([]byte, string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return nil, "", fmt.Errorf("record: %w", err)
	}
	resp, err := http.DefaultClient.Do(req) //nolint:gosec // the link was checked to be a game server's record
	if err != nil {
		return nil, "", fmt.Errorf("record: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // reading only
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512)) //nolint:errcheck // the status says enough
		return nil, "", fmt.Errorf("record: %s: %s", resp.Status, msg)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxRecord+1))
	if err != nil {
		return nil, "", fmt.Errorf("record: %w", err)
	}
	if len(data) > maxRecord {
		return nil, "", ErrRecordSize
	}
	name := ""
	if _, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition")); err == nil {
		name = record.SafeFileName(params["filename"])
	}
	return data, name, nil
}
