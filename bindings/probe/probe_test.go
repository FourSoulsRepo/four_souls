package probe

import (
	"context"
	"strings"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/diag"
)

type fakeAPI struct{}

func (fakeAPI) Context() context.Context { return context.Background() }

func newTest(variant string) (*Probe, *strings.Builder, *int) {
	var out strings.Builder
	quits := 0
	p := New(fakeAPI{}, variant)
	p.out = &out
	p.quit = func(context.Context) { quits++ }
	return p, &out, &quits
}

// ADR 035 item 5: a normal start has no variant and Report does nothing.
func TestNormalStart(t *testing.T) {
	p, out, quits := newTest("")
	if p.Variant() != "" {
		t.Errorf("Variant = %q, want empty", p.Variant())
	}
	p.Report(diag.Result{Frames: 10})
	if out.Len() != 0 || *quits != 0 {
		t.Errorf("normal start printed %q, quit %d times", out.String(), *quits)
	}
}

// ADR 035 item 5.2: a probe prints its result, with its variant, and quits.
func TestProbeReports(t *testing.T) {
	p, out, quits := newTest("no-dmabuf")
	p.Report(diag.Result{Frames: 10, Seconds: 5})
	if !strings.Contains(out.String(), `"variant":"no-dmabuf"`) || *quits != 1 {
		t.Errorf("printed %q, quit %d times", out.String(), *quits)
	}
}
