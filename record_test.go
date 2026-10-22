package main

import (
	"errors"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestFetchRecord(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("token") != "t0k3n" {
			http.Error(w, "only the game's players get its record", http.StatusForbidden)
			return
		}
		if name := r.URL.Query().Get("name"); name != "" {
			w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))
		}
		_, _ = w.Write([]byte("RECORD")) //nolint:errcheck // a test server
	}))
	defer srv.Close()
	for name, want := range map[string]string{
		"": "",
		"2026-10-10-120000-g1_Ann_Иван.fsrec": "2026-10-10-120000-g1_Ann_Иван.fsrec",
		"../../.ssh/authorized_keys":          "authorized_keys.fsrec",
		"/":                                   "",
	} {
		data, got, err := fetchRecord(t.Context(), srv.URL+"/record?game=g1&token=t0k3n&name="+url.QueryEscape(name))
		if err != nil || string(data) != "RECORD" || got != want {
			t.Errorf("name %q: got %q %q, %v; want %q", name, data, got, err, want)
		}
	}
	_, _, err := fetchRecord(t.Context(), srv.URL+"/record?game=g1&token=wrong")
	if err == nil || !strings.Contains(err.Error(), "only the game's players") {
		t.Fatalf("a refused download: %v", err)
	}
}

func TestSaveRecordChecksTheLink(t *testing.T) {
	a := &App{}
	for _, bad := range []string{"ftp://x/record", "https://x/save", "http://x/save", "file:///etc/passwd", "http:///record", "::"} {
		if _, err := a.SaveRecord(bad); !errors.Is(err, ErrRecordURL) {
			t.Errorf("%q: err %v, want ErrRecordURL", bad, err)
		}
	}
}
