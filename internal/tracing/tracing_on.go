//go:build trace

package tracing

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	httppprof "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"runtime/metrics"
	"runtime/pprof"
	rtrace "runtime/trace"
	"sync"
	"time"
)

const (
	sampleEvery   = 5 * time.Second // one line of stats.csv
	blockRate     = 10_000          // ns: block events of 10 µs and longer
	mutexFraction = 100             // one mutex event in 100
)

// profiles written when recording stops; the goroutine profile shows what
// is still running after shutdown, which is where leaks show up.
var profiles = []string{"heap", "allocs", "goroutine", "block", "mutex"}

// sampled are the runtime values of stats.csv, read without stopping the
// world.
var sampled = []string{
	"/sched/goroutines:goroutines",
	"/gc/heap/live:bytes",
	"/memory/classes/heap/objects:bytes",
	"/gc/cycles/total:gc-cycles",
}

// start begins what o asks for: files in o.Dir and a listener on o.Addr.
func start(ctx context.Context, o Options) (func() error, error) {
	var stops []func() error
	stop := func() error {
		var errs []error
		for i := len(stops) - 1; i >= 0; i-- {
			errs = append(errs, stops[i]())
		}
		return errors.Join(errs...)
	}
	if o.Dir != "" {
		s, err := startFiles(o.Dir)
		if err != nil {
			_ = stop() //nolint:errcheck // the start failed; its error is returned
			return nil, err
		}
		stops = append(stops, s)
	}
	if o.Addr != "" {
		s, err := startHTTP(ctx, o.Addr)
		if err != nil {
			_ = stop() //nolint:errcheck // the start failed; its error is returned
			return nil, err
		}
		stops = append(stops, s)
	}
	return stop, nil
}

// startFiles records the runtime trace, a CPU profile and stats.csv into
// dir, and turns on the block and mutex profiles.
func startFiles(dir string) (func() error, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("trace folder: %w", err)
	}
	r := &recorder{dir: dir}
	if err := r.start(); err != nil {
		_ = r.stop() //nolint:errcheck // the start failed; its error is returned
		return nil, err
	}
	return r.stop, nil
}

// recorder holds the files of one run.
type recorder struct {
	dir      string
	traceOut *os.File
	cpuOut   *os.File
	statsOut *os.File
	done     chan struct{} // closed to stop sampling
	sampling sync.WaitGroup
}

func (r *recorder) start() error {
	var err error
	if r.traceOut, err = r.create("trace.out"); err != nil {
		return err
	}
	if err = rtrace.Start(r.traceOut); err != nil {
		return fmt.Errorf("start trace: %w", err)
	}
	if r.cpuOut, err = r.create("cpu.pprof"); err != nil {
		return err
	}
	if err = pprof.StartCPUProfile(r.cpuOut); err != nil {
		return fmt.Errorf("start cpu profile: %w", err)
	}
	if r.statsOut, err = r.create("stats.csv"); err != nil {
		return err
	}
	runtime.SetBlockProfileRate(blockRate)
	runtime.SetMutexProfileFraction(mutexFraction)
	r.done = make(chan struct{})
	r.sampling.Add(1)
	go r.sample()
	return nil
}

// stop ends the recording and writes the profiles. It runs even when the
// start failed part way; it then has nothing more to write.
func (r *recorder) stop() error {
	var errs []error
	if r.done != nil {
		close(r.done)
		r.sampling.Wait()
	}
	pprof.StopCPUProfile()
	rtrace.Stop()
	runtime.SetBlockProfileRate(0)
	runtime.SetMutexProfileFraction(0)
	if r.done != nil {
		// A collection first, so the heap profile shows what is still live.
		runtime.GC()
		for _, name := range profiles {
			errs = append(errs, r.writeProfile(name))
		}
	}
	for _, f := range []*os.File{r.traceOut, r.cpuOut, r.statsOut} {
		if f != nil {
			errs = append(errs, f.Close())
		}
	}
	return errors.Join(errs...)
}

func (r *recorder) create(name string) (*os.File, error) {
	f, err := os.OpenFile(filepath.Join(r.dir, name), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600) //nolint:gosec // the owner names the folder with the flag
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", name, err)
	}
	return f, nil
}

// writeProfile writes one named runtime profile to <name>.pprof.
func (r *recorder) writeProfile(name string) error {
	f, err := r.create(name + ".pprof")
	if err != nil {
		return err
	}
	if err = pprof.Lookup(name).WriteTo(f, 0); err != nil {
		_ = f.Close() //nolint:errcheck // the write failed; its error is returned
		return fmt.Errorf("write %s profile: %w", name, err)
	}
	return f.Close()
}

// sample writes one line of stats.csv every sampleEvery, and one now.
func (r *recorder) sample() {
	defer r.sampling.Done()
	samples := make([]metrics.Sample, len(sampled))
	for i, name := range sampled {
		samples[i].Name = name
	}
	fmt.Fprintln(r.statsOut, "seconds,goroutines,heap_live_bytes,heap_objects_bytes,gc_cycles") //nolint:errcheck // a failed write shows when the file closes
	began := time.Now()
	line := func() {
		metrics.Read(samples)
		fmt.Fprintf(r.statsOut, "%.0f,%d,%d,%d,%d\n", time.Since(began).Seconds(), //nolint:errcheck // a failed write shows when the file closes
			value(samples[0]), value(samples[1]), value(samples[2]), value(samples[3]))
	}
	line()
	tick := time.NewTicker(sampleEvery)
	defer tick.Stop()
	for {
		select {
		case <-r.done:
			return
		case <-tick.C:
			line()
		}
	}
}

// value is a sample's number; 0 when this Go version lacks the metric.
func value(s metrics.Sample) uint64 {
	if s.Value.Kind() != metrics.KindUint64 {
		return 0
	}
	return s.Value.Uint64()
}

// startHTTP serves the pprof pages on addr. Its own mux keeps the
// profiles off any other handler.
func startHTTP(ctx context.Context, addr string) (func() error, error) {
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("pprof listen on %s: %w", addr, err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", httppprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", httppprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", httppprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", httppprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", httppprof.Trace)
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second} //nolint:gosec // a debug listener the owner sets up
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("pprof: %v", err)
		}
	}()
	// Close, not Shutdown: a profile still being taken is cut short rather
	// than waited for.
	return srv.Close, nil
}
