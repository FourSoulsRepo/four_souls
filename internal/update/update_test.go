package update

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/version"
)

// fakeGitHub serves one latest release with the given files.
type fakeGitHub struct {
	tag   string
	files map[string][]byte // asset name → body
	chunk bool              // send without Content-Length
}

func (f *fakeGitHub) serve(t *testing.T) *httptest.Server {
	t.Helper()
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/FourSoulsRepo/four_souls/releases/latest" {
			rel := ghRelease{TagName: f.tag, HTMLURL: "https://example.test/rel", Body: "notes"}
			for name := range f.files {
				rel.Assets = append(rel.Assets, ghAsset{Name: name, URL: ts.URL + "/dl/" + name})
			}
			if err := json.NewEncoder(w).Encode(rel); err != nil {
				t.Error(err)
			}
			return
		}
		body, ok := f.files[strings.TrimPrefix(r.URL.Path, "/dl/")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		if f.chunk {
			// A flush before the end makes the response chunked.
			half := len(body) / 2
			_, _ = w.Write(body[:half]) //nolint:errcheck // test server
			if fl, ok := w.(http.Flusher); ok {
				fl.Flush()
			}
			_, _ = w.Write(body[half:]) //nolint:errcheck // test server
			return
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		_, _ = w.Write(body) //nolint:errcheck // test server
	}))
	t.Cleanup(ts.Close)
	return ts
}

// sums writes SHA256SUMS.txt over the given files, as release.yml does.
func sums(files map[string][]byte) []byte {
	var b strings.Builder
	for name, body := range files {
		h := sha256.Sum256(body)
		b.WriteString(hex.EncodeToString(h[:]) + "  " + name + "\n")
	}
	return []byte(b.String())
}

var allTargets = []string{"linux-amd64.tar.gz", "linux-arm64.tar.gz", "windows-amd64.zip", "windows-arm64.zip", "macos-universal.zip"}

// fullRelease has every archive of both editions, each with its own body.
func fullRelease(ver string) map[string][]byte {
	files := map[string][]byte{}
	for _, target := range allTargets {
		for _, prefix := range []string{"", "for-poor-"} {
			name := "four-souls-" + ver + "-" + prefix + target
			files[name] = []byte(name)
		}
	}
	files["four-souls-server-"+ver+"-linux-amd64.tar.gz"] = []byte("server")
	files[sumsName] = sums(files)
	return files
}

func newTestClient(t *testing.T, ts *httptest.Server, current, edition, goos, goarch string) *Client {
	t.Helper()
	c := New(Config{
		Owner: "FourSoulsRepo", Repo: "four_souls", Name: "FourSouls",
		Current: current, Edition: edition, GOOS: goos, GOARCH: goarch,
	})
	c.API = ts.URL
	c.HTTP = ts.Client()
	c.TempDir = t.TempDir()
	c.Executable = filepath.Join(t.TempDir(), "FourSouls")
	if err := os.WriteFile(c.Executable, []byte("old"), 0o700); err != nil { //nolint:gosec // test file
		t.Fatal(err)
	}
	c.Start = func(string, ...string) error { return nil }
	return c
}

func TestNewer(t *testing.T) {
	for _, tc := range []struct {
		latest, current string
		want            bool
	}{
		{"v0.5.0", "0.4.2", true},
		{"v0.4.10", "v0.4.9", true},
		{"v1.0.0", "0.99.99", true},
		{"v0.4.2", "0.4.2", false},
		{"v0.4.1", "0.4.2", false},
		{"v0.5.0", "dev", false},
		{"v0.5.0", "internal-abc123", false},
		{"nightly", "0.4.2", false},
		{"v0.5.0", "0.4.2-test", true},
		{"v0.4.2+build", "0.4.2", false},
		{"v0.5", "0.4.2", false},
		{"v0.5.x", "0.4.2", false},
	} {
		if got := Newer(tc.latest, tc.current); got != tc.want {
			t.Errorf("Newer(%q, %q) = %v, want %v", tc.latest, tc.current, got, tc.want)
		}
	}
}

func TestAssetName(t *testing.T) {
	for _, tc := range []struct {
		edition, goos, goarch, want string
	}{
		{version.EditionFull, "darwin", "arm64", "four-souls-1.2.0-macos-universal.zip"},
		{version.EditionFull, "darwin", "amd64", "four-souls-1.2.0-macos-universal.zip"},
		{version.EditionFull, "windows", "amd64", "four-souls-1.2.0-windows-amd64.zip"},
		{version.EditionFull, "windows", "arm64", "four-souls-1.2.0-windows-arm64.zip"},
		{version.EditionFull, "linux", "amd64", "four-souls-1.2.0-linux-amd64.tar.gz"},
		{version.EditionFull, "linux", "arm64", "four-souls-1.2.0-linux-arm64.tar.gz"},
		{version.EditionForPoor, "darwin", "arm64", "four-souls-1.2.0-for-poor-macos-universal.zip"},
		{version.EditionForPoor, "windows", "arm64", "four-souls-1.2.0-for-poor-windows-arm64.zip"},
		{version.EditionForPoor, "linux", "amd64", "four-souls-1.2.0-for-poor-linux-amd64.tar.gz"},
	} {
		got, err := AssetName("v1.2.0", tc.edition, tc.goos, tc.goarch)
		if err != nil || got != tc.want {
			t.Errorf("%s %s/%s = %q, %v; want %q", tc.edition, tc.goos, tc.goarch, got, err, tc.want)
		}
	}
	for _, sys := range [][2]string{{"linux", "386"}, {"windows", "386"}, {"freebsd", "amd64"}} {
		if _, err := AssetName("v1.2.0", version.EditionFull, sys[0], sys[1]); !errors.Is(err, ErrNoAsset) {
			t.Errorf("%s/%s: err = %v, want ErrNoAsset", sys[0], sys[1], err)
		}
	}
}

func TestCheckPicksOwnArchive(t *testing.T) {
	ts := (&fakeGitHub{tag: "v1.2.0", files: fullRelease("1.2.0")}).serve(t)
	for _, edition := range []string{version.EditionFull, version.EditionForPoor} {
		for _, sys := range [][2]string{
			{"darwin", "amd64"},
			{"darwin", "arm64"},
			{"windows", "amd64"},
			{"windows", "arm64"},
			{"linux", "amd64"},
			{"linux", "arm64"},
		} {
			c := newTestClient(t, ts, "1.1.0", edition, sys[0], sys[1])
			st, err := c.Check(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			want, _ := AssetName("v1.2.0", edition, sys[0], sys[1]) //nolint:errcheck // checked in TestAssetName
			if st.Asset != want || !st.CanInstall || st.Reason != ReasonNone {
				t.Errorf("%s %s/%s: %+v, want asset %s", edition, sys[0], sys[1], st, want)
			}
			poor := strings.Contains(st.Asset, "for-poor-")
			if poor != (edition == version.EditionForPoor) {
				t.Errorf("%s build picked %s", edition, st.Asset)
			}
		}
	}
}

// Each edition sees only the other edition's archive: no update.
func TestCheckNeverSwapsEdition(t *testing.T) {
	for _, tc := range []struct{ edition, only string }{
		{version.EditionFull, "four-souls-1.2.0-for-poor-windows-amd64.zip"},
		{version.EditionForPoor, "four-souls-1.2.0-windows-amd64.zip"},
	} {
		files := map[string][]byte{tc.only: []byte("x")}
		files[sumsName] = sums(files)
		ts := (&fakeGitHub{tag: "v1.2.0", files: files}).serve(t)
		c := newTestClient(t, ts, "1.1.0", tc.edition, "windows", "amd64")
		st, err := c.Check(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if st.Asset != "" || st.CanInstall || st.Reason != ReasonNoArchive {
			t.Errorf("%s: %+v, want no archive", tc.edition, st)
		}
		if err := c.Download(context.Background(), nil); !errors.Is(err, ErrNotReady) {
			t.Errorf("%s: Download err = %v, want ErrNotReady", tc.edition, err)
		}
	}
}

func TestCheckReasons(t *testing.T) {
	ts := (&fakeGitHub{tag: "v1.2.0", files: fullRelease("1.2.0")}).serve(t)
	for _, tc := range []struct {
		current, goarch string
		want            Reason
	}{
		{"dev", "amd64", ReasonNotVersion},
		{"internal-abc123", "amd64", ReasonNotVersion},
		{"1.2.0", "amd64", ReasonNotNewer},
		{"v1.3.0", "amd64", ReasonNotNewer},
		{"1.1.0", "386", ReasonNoArchive},
		{"1.1.0", "amd64", ReasonNone},
	} {
		c := newTestClient(t, ts, tc.current, version.EditionFull, "linux", tc.goarch)
		st, err := c.Check(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if st.Reason != tc.want || st.CanInstall != (tc.want == ReasonNone) {
			t.Errorf("current %s, %s: reason %q canInstall %v, want %q", tc.current, tc.goarch, st.Reason, st.CanInstall, tc.want)
		}
		if st.Releases != "https://github.com/FourSoulsRepo/four_souls/releases" || st.Latest != "v1.2.0" {
			t.Errorf("status %+v", st)
		}
	}
}

func TestCheckNotWritable(t *testing.T) {
	if runtime.GOOS == "windows" || os.Getuid() == 0 {
		t.Skip("needs a folder this user cannot write")
	}
	ts := (&fakeGitHub{tag: "v1.2.0", files: fullRelease("1.2.0")}).serve(t)
	c := newTestClient(t, ts, "1.1.0", version.EditionFull, "linux", "amd64")
	dir := filepath.Dir(c.Executable)
	if err := os.Chmod(dir, 0o500); err != nil { //nolint:gosec // test folder
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) }) //nolint:errcheck,gosec // test cleanup
	st, err := c.Check(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if st.Reason != ReasonNotWritable || st.CanInstall {
		t.Errorf("%+v, want not_writable", st)
	}
}

func TestCheckGitHubDown(t *testing.T) {
	ts := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(ts.Close)
	c := newTestClient(t, ts, "1.1.0", version.EditionFull, "linux", "amd64")
	if _, err := c.Check(context.Background()); err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("err = %v, want HTTP 404", err)
	}
}

// linuxRelease is a release whose linux-amd64 archive really unpacks.
func linuxRelease(t *testing.T) map[string][]byte {
	t.Helper()
	name := "four-souls-1.2.0-linux-amd64.tar.gz"
	files := map[string][]byte{name: makeTarGz(t, releaseEntries("linux"))}
	files[sumsName] = sums(files)
	return files
}

func TestDownloadAndApply(t *testing.T) {
	fake := &fakeGitHub{tag: "v1.2.0", files: linuxRelease(t)}
	ts := fake.serve(t)
	c := newTestClient(t, ts, "1.1.0", version.EditionFull, "linux", "amd64")
	if err := c.Apply(); !errors.Is(err, ErrNoDownload) {
		t.Fatalf("Apply before download: %v", err)
	}
	if err := c.Download(context.Background(), nil); !errors.Is(err, ErrNotReady) {
		t.Fatalf("Download before check: %v", err)
	}
	if _, err := c.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := c.Download(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if filepath.Base(c.payload) != "FourSouls" {
		t.Fatalf("payload %s", c.payload)
	}
	var started []string
	c.Start = func(name string, args ...string) error {
		started = append([]string{name}, args...)
		return nil
	}
	if err := c.Apply(); err != nil {
		t.Fatal(err)
	}
	if len(started) != 2 || started[0] != "/bin/sh" {
		t.Fatalf("started %q", started)
	}
	script, err := os.ReadFile(started[1])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(script), "NEW="+shQuote(c.payload)) {
		t.Errorf("helper does not name the payload:\n%s", script)
	}
}

func TestDownloadErrors(t *testing.T) {
	name := "four-souls-1.2.0-linux-amd64.tar.gz"
	good := linuxRelease(t)
	for _, tc := range []struct {
		label string
		sums  []byte // nil: no SHA256SUMS.txt
		want  error
	}{
		{"no sums", nil, ErrNoSums},
		{"wrong hash", []byte(strings.Repeat("b", 64) + "  " + name + "\n"), ErrChecksum},
		{"not listed", []byte(strings.Repeat("b", 64) + "  four-souls-1.2.0-for-poor-linux-amd64.tar.gz\n"), ErrChecksum},
		{"bad line", []byte("xyz  " + name + "\n"), ErrChecksum},
	} {
		files := map[string][]byte{name: good[name]}
		if tc.sums != nil {
			files[sumsName] = tc.sums
		}
		ts := (&fakeGitHub{tag: "v1.2.0", files: files}).serve(t)
		c := newTestClient(t, ts, "1.1.0", version.EditionFull, "linux", "amd64")
		if _, err := c.Check(context.Background()); err != nil {
			t.Fatal(err)
		}
		err := c.Download(context.Background(), nil)
		if !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.label, err, tc.want)
		}
		if _, err := os.Stat(filepath.Join(c.workDir, name)); err == nil {
			t.Errorf("%s: the bad archive was kept", tc.label)
		}
		if err := c.Apply(); !errors.Is(err, ErrNoDownload) {
			t.Errorf("%s: Apply err = %v, want ErrNoDownload", tc.label, err)
		}
	}
}

func TestDownloadProgress(t *testing.T) {
	for _, chunk := range []bool{false, true} {
		files := linuxRelease(t)
		size := int64(len(files["four-souls-1.2.0-linux-amd64.tar.gz"]))
		ts := (&fakeGitHub{tag: "v1.2.0", files: files, chunk: chunk}).serve(t)
		c := newTestClient(t, ts, "1.1.0", version.EditionFull, "linux", "amd64")
		c.ProgressEvery = time.Hour // only the first and the last call
		if _, err := c.Check(context.Background()); err != nil {
			t.Fatal(err)
		}
		var calls [][2]int64
		err := c.Download(context.Background(), func(done, total int64) {
			calls = append(calls, [2]int64{done, total})
		})
		if err != nil {
			t.Fatal(err)
		}
		wantTotal := size
		if chunk {
			wantTotal = 0
		}
		last := calls[len(calls)-1]
		if len(calls) != 2 || last[0] != size || last[1] != wantTotal {
			t.Errorf("chunked %v: calls %v, want 2 ending at %d of %d", chunk, calls, size, wantTotal)
		}
	}
}

func TestProgressThrottle(t *testing.T) {
	var n int
	pw := &progressWriter{report: func(int64, int64) { n++ }, every: time.Hour}
	pw.call()
	for range 1000 {
		if _, err := pw.Write(make([]byte, 10)); err != nil {
			t.Fatal(err)
		}
	}
	pw.call()
	if n != 2 || pw.done != 10000 {
		t.Errorf("calls %d, done %d", n, pw.done)
	}
}

func TestChecksumExactName(t *testing.T) {
	a, b := strings.Repeat("a", 64), strings.Repeat("b", 64)
	list := a + "  four-souls-1.2.0-for-poor-windows-amd64.zip\n" + b + "  four-souls-1.2.0-windows-amd64.zip\n"
	if got, err := checksumFor(list, "four-souls-1.2.0-windows-amd64.zip"); err != nil || got != b {
		t.Errorf("full: %q %v", got, err)
	}
	if got, err := checksumFor(list, "four-souls-1.2.0-for-poor-windows-amd64.zip"); err != nil || got != a {
		t.Errorf("for-poor: %q %v", got, err)
	}
	if _, err := checksumFor(list, "windows-amd64.zip"); !errors.Is(err, ErrChecksum) {
		t.Errorf("suffix matched: %v", err)
	}
}

func TestQuote(t *testing.T) {
	if got := shQuote(`it's "here" $HOME`); got != `'it'\''s "here" $HOME'` {
		t.Errorf("shQuote = %s", got)
	}
	if got := psQuote("C:\\O'Neil\u2019s"); got != "'C:\\O''Neil\u2019\u2019s'" {
		t.Errorf("psQuote = %s", got)
	}
}

func TestWindowsHelper(t *testing.T) {
	dest := `C:\Users\Олег\Four Souls\FourSouls.exe`
	got := string(windowsHelper(dest, `C:\Temp\new\FourSouls.exe`, dest+".old", 42))
	if !strings.HasPrefix(got, "\uFEFF") {
		t.Error("no UTF-8 BOM")
	}
	for _, want := range []string{psQuote(dest), "Wait-Process -Id 42", "Start-Process"} {
		if !strings.Contains(got, want) {
			t.Errorf("helper has no %q:\n%s", want, got)
		}
	}
}

// The sh helper really swaps an install whose path has spaces and quotes.
func TestUnixHelperRuns(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sh helper")
	}
	dir := filepath.Join(t.TempDir(), `it's a "dir" $HOME`)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "FourSouls")
	payload := filepath.Join(t.TempDir(), "FourSouls")
	marker := filepath.Join(dir, "started")
	if err := os.WriteFile(dest, []byte("#!/bin/sh\n"), 0o700); err != nil { //nolint:gosec // test file
		t.Fatal(err)
	}
	newApp := "#!/bin/sh\necho new > " + shQuote(marker) + "\n"
	if err := os.WriteFile(payload, []byte(newApp), 0o700); err != nil { //nolint:gosec // test file
		t.Fatal(err)
	}
	// A process that already quit stands for the app.
	done := exec.CommandContext(context.Background(), "true")
	if err := done.Run(); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(t.TempDir(), "apply.sh")
	body := unixHelper(dest, payload, dest+".FourSouls-old", done.Process.Pid, false)
	if err := os.WriteFile(script, body, 0o700); err != nil { //nolint:gosec // test file
		t.Fatal(err)
	}
	if out, err := exec.CommandContext(context.Background(), "/bin/sh", script).CombinedOutput(); err != nil { //nolint:gosec // test script
		t.Fatalf("%v: %s", err, out)
	}
	got, err := os.ReadFile(dest) //nolint:gosec // test file
	if err != nil || string(got) != newApp {
		t.Fatalf("install = %q, %v", got, err)
	}
	if _, err := os.Stat(dest + ".FourSouls-old"); err == nil {
		t.Error("old install left behind")
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the new app did not start")
		}
		time.Sleep(20 * time.Millisecond)
	}
}
