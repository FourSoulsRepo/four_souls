package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/FourSoulsRepo/record"
)

// ErrRecordURL means a record URL is not a game server's record link.
var ErrRecordURL = errors.New("not a record link")

// ErrRecordSize means a server sent more than a record can be.
var ErrRecordSize = errors.New("record: the file is too large")

// maxRecord is the largest record the app downloads; records are a few
// hundred KB.
const maxRecord = 64 << 20

// SaveRecord downloads a finished game's record from the game server
// and asks where to save it (RP-05). link is http://host:port/record?…;
// the dialog suggests the server's file name, cleaned again. The result
// is the saved file, or "" when the player cancelled.
func (a *App) SaveRecord(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil || u.Scheme != "http" || u.Path != "/record" || u.Host == "" {
		return "", ErrRecordURL
	}
	data, name, err := fetchRecord(a.ctx, u.String())
	if err != nil {
		return "", err
	}
	if name == "" {
		name = "four-souls-" + record.SafeName(u.Query().Get("game")) + record.Ext
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save the match record",
		DefaultFilename: name,
		Filters:         []runtime.FileFilter{{DisplayName: "Four Souls records", Pattern: "*" + record.Ext}},
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
