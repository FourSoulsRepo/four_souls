// Package probe is the binding of a `--diagnose` probe window (ADR 035):
// it tells the page which variant runs, and takes the frames it measured.
package probe

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/FourSoulsRepo/four_souls/internal/diag"
)

// api is what the binding needs from the app (ADR 030, item 3).
type api interface {
	// Context is the Wails context; Report quits the app with it.
	Context() context.Context
}

// Probe is bound in every start (ADR 035, item 5); only a probe start
// has a variant.
type Probe struct {
	api     api
	variant string

	// Where the result goes, and how the app quits; tests replace them.
	out  io.Writer
	quit func(ctx context.Context)
}

// New makes the binding; variant is "" in a normal start.
func New(a api, variant string) *Probe {
	return &Probe{api: a, variant: variant, out: os.Stdout, quit: runtime.Quit}
}

// Variant is the probe's variant; "" in a normal start, so the page
// starts as usual.
func (p *Probe) Variant() string {
	return p.variant
}

// Report prints the page's result for the parent and quits; in a
// normal start it does nothing.
func (p *Probe) Report(r diag.Result) {
	if p.variant == "" {
		return
	}
	r.Variant = p.variant
	_, _ = fmt.Fprintln(p.out, diag.ResultLine(r)) //nolint:errcheck // the parent reports a missing result
	p.quit(p.api.Context())
}
