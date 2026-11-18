package diag

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArgs(t *testing.T) {
	cases := []struct {
		args     []string
		diagnose bool
		probe    string
	}{
		{nil, false, ""},
		{[]string{"--diagnose"}, true, ""},
		{[]string{"-diagnose"}, true, ""},
		{[]string{"--diagnose-probe=no-dmabuf"}, false, "no-dmabuf"},
		// macOS may pass a process serial number; it is not ours.
		{[]string{"-psn_0_12345"}, false, ""},
		{[]string{"diagnose"}, false, ""},
	}
	for _, c := range cases {
		d, p := Args(c.args)
		if d != c.diagnose || p != c.probe {
			t.Errorf("Args(%q) = %v, %q; want %v, %q", c.args, d, p, c.diagnose, c.probe)
		}
	}
}

// ADR 035 item 3.1: five variants on Linux, the default elsewhere.
func TestVariants(t *testing.T) {
	if n := len(Variants("linux")); n != 5 {
		t.Errorf("linux: %d variants, want 5", n)
	}
	if v := Variants("darwin"); len(v) != 1 || v[0].ID != "default" {
		t.Errorf("darwin: %+v, want only default", v)
	}
	if _, ok := VariantByID("darwin", "no-dmabuf"); ok {
		t.Error("no-dmabuf found on darwin")
	}
	if v, ok := VariantByID("linux", "no-dmabuf"); !ok || v.Env[0] != "WEBKIT_DISABLE_DMABUF_RENDERER=1" {
		t.Errorf("no-dmabuf = %+v, %v", v, ok)
	}
}

func TestResultLineRoundTrip(t *testing.T) {
	r := Result{Frames: 300, Seconds: 5, Renderer: "Mesa Intel", Sounds: []string{"mp3"}}
	out := "INF | starting\n" + ResultLine(r) + "\nbye\n"
	got, ok := parseResult([]byte(out))
	if !ok || got.Frames != 300 || got.Renderer != "Mesa Intel" {
		t.Fatalf("parseResult = %+v, %v", got, ok)
	}
	if got.FPS() != 60 {
		t.Errorf("FPS = %v, want 60", got.FPS())
	}
	if _, ok := parseResult([]byte("nothing here\n")); ok {
		t.Error("a result found in output without one")
	}
}

func TestProbeAll(t *testing.T) {
	vs := Variants("linux")[:3]
	run := func(v Variant) ([]byte, error) {
		switch v.ID {
		case "default":
			return []byte(ResultLine(Result{Frames: 100, Seconds: 5})), nil
		case "no-dmabuf":
			return nil, context.DeadlineExceeded
		}
		return nil, errors.New("exit status 1")
	}
	var heard int
	outs := ProbeAll(vs, run, func(Outcome) { heard++ })
	if heard != 3 || len(outs) != 3 {
		t.Fatalf("heard %d, outs %d; want 3, 3", heard, len(outs))
	}
	if outs[0].Result == nil || outs[0].Result.Variant != "default" {
		t.Errorf("default: %+v", outs[0])
	}
	if !strings.Contains(outs[1].Err, "no result in") {
		t.Errorf("timeout: %q", outs[1].Err)
	}
	if outs[2].Err != "exit status 1" {
		t.Errorf("failure: %q", outs[2].Err)
	}
}

func TestBestAndFlags(t *testing.T) {
	outs := []Outcome{
		{Variant: Variant{ID: "default"}, Result: &Result{Frames: 75, Seconds: 5, Renderer: "llvmpipe (LLVM 17)"}},
		{Variant: Variant{ID: "no-dmabuf", Env: []string{"WEBKIT_DISABLE_DMABUF_RENDERER=1"}}, Result: &Result{Frames: 300, Seconds: 5, Sounds: []string{"mp3"}}},
		{Variant: Variant{ID: "gpu-never"}, Err: "exit status 1"},
	}
	best, ok := Best(outs)
	if !ok || best.Variant.ID != "no-dmabuf" {
		t.Fatalf("Best = %+v, %v", best, ok)
	}
	s := System{OS: "linux", GPUs: []GPU{{Node: "/dev/dri/renderD128", Err: "permission denied"}}}
	flags := strings.Join(Flags(s, outs), "\n")
	for _, want := range []string{"does not open", "llvmpipe", "No MP3"} {
		if !strings.Contains(flags, want) {
			t.Errorf("flags miss %q:\n%s", want, flags)
		}
	}
	if strings.Contains(flags, "Every variant is slow") {
		t.Errorf("60 fps called slow:\n%s", flags)
	}
	var b strings.Builder
	Write(&b, s, outs)
	report := b.String()
	for _, want := range []string{"Best: no-dmabuf, 60.0 fps.", "WEBKIT_DISABLE_DMABUF_RENDERER=1 ./FourSouls", "gpu-never", "failed: exit status 1"} {
		if !strings.Contains(report, want) {
			t.Errorf("report misses %q:\n%s", want, report)
		}
	}
}

func TestGPUs(t *testing.T) {
	root := t.TempDir()
	sys, dev := filepath.Join(root, "sys"), filepath.Join(root, "dev")
	device := filepath.Join(sys, "renderD128", "device")
	for _, dir := range []string{device, dev, filepath.Join(root, "drivers", "i915")} {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(filepath.Join(root, "drivers", "i915"), filepath.Join(device, "driver")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(device, "vendor"), []byte("0x8086\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dev, "renderD128"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	got := gpus(sys, dev)
	if len(got) != 1 || got[0].Driver != "i915" || got[0].Vendor != "Intel" || !got[0].Opens {
		t.Errorf("gpus = %+v", got)
	}
}
