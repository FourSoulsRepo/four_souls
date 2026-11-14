package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/update"
)

// errKey is the message key of a binding error, as the client reads it.
func errKey(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		t.Fatal("no error")
	}
	var m struct {
		Key string `json:"key"`
	}
	if jerr := json.Unmarshal([]byte(err.Error()), &m); jerr != nil {
		t.Fatalf("%q is no message: %v", err.Error(), jerr)
	}
	return m.Key
}

func TestUpdateErrorKeys(t *testing.T) {
	if updateError(nil, errUpdateCheck) != nil {
		t.Error("nil error gave a message")
	}
	for _, tc := range []struct {
		err  error
		want string
	}{
		{update.ErrNoAsset, "error.update_no_archive"},
		{fmt.Errorf("%w: HTTP 404", update.ErrNoSums), "error.update_no_sums"},
		{fmt.Errorf("%w: got a, want b", update.ErrChecksum), "error.update_checksum"},
		{update.ErrTooLarge, "error.update_too_large"},
		{update.ErrBadArchive, "error.update_bad_archive"},
		{update.ErrNotReady, "error.update_not_ready"},
		{update.ErrNoDownload, "error.update_no_download"},
		{errors.New("connection refused"), "error.update_download"},
	} {
		err := updateError(tc.err, errUpdateDownload)
		if got := errKey(t, err); got != tc.want {
			t.Errorf("%v: key %q, want %q", tc.err, got, tc.want)
		}
		if !errors.Is(err, tc.err) {
			t.Errorf("%v: cause lost", tc.err)
		}
	}
}

// fakeUpdater points the bindings at a fake GitHub that answers with
// status code; it records what the runtime was asked.
func fakeUpdater(t *testing.T, code int) (u *Updater, opened *[]string, quits *int) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if code != http.StatusOK {
			w.WriteHeader(code)
			return
		}
		_, _ = w.Write([]byte(`{"tag_name":"v9.0.0","html_url":"https://example.test/v9","body":"new"}`)) //nolint:errcheck // test server
	}))
	t.Cleanup(ts.Close)
	u = NewUpdater()
	u.client.API = ts.URL
	opened, quits = new([]string), new(int)
	u.open = func(url string) { *opened = append(*opened, url) }
	u.quit = func() { *quits++ }
	u.emit = func(name string, data any) { t.Errorf("event %s %v", name, data) }
	return u, opened, quits
}

func TestOpenReleases(t *testing.T) {
	u, opened, _ := fakeUpdater(t, http.StatusOK)
	u.OpenReleases()
	st, err := u.Check()
	if err != nil {
		t.Fatal(err)
	}
	if st.Latest != "v9.0.0" || st.Notes != "new" {
		t.Errorf("status %+v", st)
	}
	u.OpenReleases()
	want := []string{"https://github.com/FourSoulsRepo/four_souls/releases", "https://example.test/v9"}
	if strings.Join(*opened, " ") != strings.Join(want, " ") {
		t.Errorf("opened %q, want %q", *opened, want)
	}
}

func TestCheckFails(t *testing.T) {
	u, _, _ := fakeUpdater(t, http.StatusInternalServerError)
	_, err := u.Check()
	if got := errKey(t, err); got != "error.update_check" {
		t.Errorf("key %q", got)
	}
}

// Download and install refuse to run out of order; nothing quits.
func TestUpdateOrder(t *testing.T) {
	u, _, quits := fakeUpdater(t, http.StatusOK)
	if got := errKey(t, u.Download()); got != "error.update_not_ready" {
		t.Errorf("download: key %q", got)
	}
	if got := errKey(t, u.ApplyAndQuit()); got != "error.update_no_download" {
		t.Errorf("apply: key %q", got)
	}
	if *quits != 0 {
		t.Error("the app quit")
	}
}
