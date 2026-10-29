//go:build trace

package tracing

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Recording writes every file when it stops, and the live pages answer
// while it runs.
func TestStartWritesFilesAndServes(t *testing.T) {
	var lc net.ListenConfig
	ln, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	if err = ln.Close(); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "traces")
	stop, err := Start(t.Context(), Options{Dir: dir, Addr: addr})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Get("http://" + addr + "/debug/pprof/") //nolint:noctx // test
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close() //nolint:errcheck // test
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("pprof index: %s", resp.Status)
	}
	if err = stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	want := []string{"trace.out", "cpu.pprof", "stats.csv", "heap.pprof", "allocs.pprof", "goroutine.pprof", "block.pprof", "mutex.pprof"}
	for _, name := range want {
		info, serr := os.Stat(filepath.Join(dir, name))
		if serr != nil {
			t.Errorf("%s: %v", name, serr)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("%s is empty", name)
		}
	}
	csv, err := os.ReadFile(filepath.Join(dir, "stats.csv")) //nolint:gosec // test file
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(csv), "seconds,goroutines,") {
		t.Errorf("stats.csv header: %q", string(csv))
	}
	if resp, err = http.Get("http://" + addr + "/debug/pprof/"); err == nil { //nolint:noctx // test
		_ = resp.Body.Close() //nolint:errcheck // test
		t.Error("pprof still serves after stop")
	}
}
