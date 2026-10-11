package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("token") != "t0k3n" {
			http.Error(w, "only the game's players get its record", http.StatusForbidden)
			return
		}
		_, _ = w.Write([]byte("RECORD")) //nolint:errcheck // a test server
	}))
	defer srv.Close()
	path := filepath.Join(t.TempDir(), "g1.fsrec")
	if err := download(t.Context(), srv.URL+"/record?game=g1&token=t0k3n", path); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != "RECORD" { //nolint:gosec // the test's temp file
		t.Fatalf("saved %q, %v", data, err)
	}
	err := download(t.Context(), srv.URL+"/record?game=g1&token=wrong", path)
	if err == nil || !strings.Contains(err.Error(), "only the game's players") {
		t.Fatalf("a refused download: %v", err)
	}
}

func TestSaveRecordChecksTheLink(t *testing.T) {
	a := &App{}
	for _, bad := range []string{"https://x/record", "http://x/save", "file:///etc/passwd", "http:///record", "::"} {
		if _, err := a.SaveRecord(bad); !errors.Is(err, ErrRecordURL) {
			t.Errorf("%q: err %v, want ErrRecordURL", bad, err)
		}
	}
}
