package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// ErrServerURL means a base is not a game server's http(s) address.
var ErrServerURL = errors.New("not a server address")

// askWait is how long one server address may take to answer.
const askWait = 3 * time.Second

// maxAnswer caps a small JSON answer from a server.
const maxAnswer = 4 << 10

// CheckSeat asks a game server whether token still holds a seat there
// (15.99), so the menu offers Rejoin. bases are the server's http(s)
// addresses to try in order, as the frontend builds them from the
// address it joined (wss first, then ws, ADR 017).
func (a *App) CheckSeat(bases []string, token string) (protocol.SeatInfo, error) {
	var info protocol.SeatInfo
	err := askServer(a.appContext(), bases, "/seat", url.Values{"token": {token}}, &info)
	return info, err
}

// appContext is the app's context, or a background one before startup.
func (a *App) appContext() context.Context {
	if a.ctx == nil {
		return context.Background()
	}
	return a.ctx
}

// askServer gets path from the first of bases that answers, as JSON
// into out.
func askServer(ctx context.Context, bases []string, path string, query url.Values, out any) error {
	err := ErrServerURL
	for _, b := range bases {
		u, perr := url.Parse(b)
		if perr != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || (u.Path != "" && u.Path != "/") || u.RawQuery != "" {
			return ErrServerURL
		}
		u.Path, u.RawQuery = path, query.Encode()
		if err = getJSON(ctx, u.String(), out); err == nil {
			return nil
		}
	}
	return err
}

func getJSON(ctx context.Context, link string, out any) error {
	ctx, cancel := context.WithTimeout(ctx, askWait)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}
	resp, err := http.DefaultClient.Do(req) //nolint:gosec // the base was checked to be an http(s) server address
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // reading only
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server: %s", resp.Status)
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxAnswer)).Decode(out); err != nil {
		return fmt.Errorf("server: %w", err)
	}
	return nil
}
