// Package tracing records Go runtime traces and profiles of the app and
// the dedicated server (ADR 024). The recording code is compiled only
// with the trace build tag. Without it, Start refuses any options, so the
// default binaries carry none of that code.
package tracing

import (
	"context"
	"errors"
	"flag"
	"os"
)

// Environment variables for the app, which takes no flags. The server
// takes the same settings as flags.
const (
	EnvDir  = "FOUR_SOULS_TRACE_DIR"
	EnvAddr = "FOUR_SOULS_PPROF_ADDR"
)

// ErrNotBuilt means the options ask for tracing in a build without the
// trace tag.
var ErrNotBuilt = errors.New("tracing: this build has no tracing; rebuild with -tags trace")

// Options say what to record. The zero value records nothing.
type Options struct {
	Dir  string // folder for trace.out, cpu.pprof and the other files; "" writes none
	Addr string // address for live profiles over HTTP, e.g. 127.0.0.1:6060; "" serves none
}

// Flags adds the tracing flags to fs and returns the options they set.
// Both builds take the flags, so a default build says what it lacks.
func Flags(fs *flag.FlagSet) *Options {
	var o Options
	fs.StringVar(&o.Dir, "trace_dir", "", "folder for runtime traces and profiles, written when the server stops; empty writes none")
	fs.StringVar(&o.Addr, "pprof_addr", "", "address for live profiles over HTTP, e.g. 127.0.0.1:6060; empty serves none")
	return &o
}

// FromEnv reads the options from the environment, for the app.
func FromEnv() Options {
	return Options{Dir: os.Getenv(EnvDir), Addr: os.Getenv(EnvAddr)}
}

// Start begins recording what o asks for. The returned stop writes the
// files and closes the listener; call it once, when the process is done.
// Zero options start nothing.
func Start(ctx context.Context, o Options) (func() error, error) {
	if o == (Options{}) {
		return func() error { return nil }, nil
	}
	return start(ctx, o)
}
