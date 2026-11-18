package diag

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ProbeTimeout is how long the parent waits for one probe (ADR 035,
// item 3.8): start, 1 s to settle, 5 s of frames, and room to spare.
const ProbeTimeout = 30 * time.Second

// Outcome is one probe's result, or why there is none.
type Outcome struct {
	Variant Variant
	Result  *Result
	Err     string
}

// RunFunc starts one probe and returns its stdout.
type RunFunc func(v Variant) ([]byte, error)

// Command starts exe as a probe of v: its own flag, and v's switches
// on top of this environment.
func Command(ctx context.Context, exe string, v Variant) *exec.Cmd {
	cmd := exec.CommandContext(ctx, exe, "--"+ProbeName+"="+v.ID) //nolint:gosec // exe is this app (os.Executable); v from Variants
	cmd.Env = append(os.Environ(), v.Env...)
	cmd.WaitDelay = time.Second
	return cmd
}

// ProbeAll runs each variant in turn; done hears each outcome as it
// comes, so the terminal shows progress.
func ProbeAll(variants []Variant, run RunFunc, done func(Outcome)) []Outcome {
	outs := make([]Outcome, 0, len(variants))
	for _, v := range variants {
		o := probeOne(v, run)
		if done != nil {
			done(o)
		}
		outs = append(outs, o)
	}
	return outs
}

func probeOne(v Variant, run RunFunc) Outcome {
	out, err := run(v)
	if r, ok := parseResult(out); ok {
		r.Variant = v.ID
		return Outcome{Variant: v, Result: &r}
	}
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return Outcome{Variant: v, Err: "no result in " + ProbeTimeout.String()}
	case err != nil:
		return Outcome{Variant: v, Err: failure(err)}
	}
	return Outcome{Variant: v, Err: "the window sent no result"}
}

// failure is a short reason, with the last line the probe wrote to
// stderr when there is one.
func failure(err error) string {
	var exit *exec.ExitError
	if errors.As(err, &exit) && len(exit.Stderr) > 0 {
		lines := strings.Split(strings.TrimSpace(string(exit.Stderr)), "\n")
		return fmt.Sprintf("%v: %s", err, lines[len(lines)-1])
	}
	return err.Error()
}
