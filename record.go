package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/FourSoulsRepo/record"
)

// ErrRecordURL means a record URL is not a game server's record link.
var ErrRecordURL = errors.New("not a record link")

// SaveRecord asks where to save a finished game's record and downloads
// it from the game server (RP-05). link is http://host:port/record?…;
// the result is the saved file, or "" when the player cancelled.
func (a *App) SaveRecord(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil || u.Scheme != "http" || u.Path != "/record" || u.Host == "" {
		return "", ErrRecordURL
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save the match record",
		DefaultFilename: "four-souls-" + u.Query().Get("game") + record.Ext,
		Filters:         []runtime.FileFilter{{DisplayName: "Four Souls records", Pattern: "*" + record.Ext}},
	})
	if err != nil || path == "" {
		return "", err //nolint:wrapcheck // shown to the player as is
	}
	return path, download(a.ctx, u.String(), path)
}

// download saves what link serves into path.
func download(ctx context.Context, link, path string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return fmt.Errorf("record: %w", err)
	}
	resp, err := http.DefaultClient.Do(req) //nolint:gosec // the link was checked to be a game server's record
	if err != nil {
		return fmt.Errorf("record: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // reading only
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512)) //nolint:errcheck // the status says enough
		return fmt.Errorf("record: %s: %s", resp.Status, msg)
	}
	f, err := os.Create(path) //nolint:gosec // the player chose the file
	if err != nil {
		return fmt.Errorf("record: %w", err)
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		_ = f.Close() //nolint:errcheck // the copy error matters
		return fmt.Errorf("record: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("record: %w", err)
	}
	return nil
}
