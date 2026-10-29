# Profiling

Profiling finds where the server and the app spend time and memory: hot spots, leaks and contention. It is compiled in only with the `trace` build tag (ADR 024). Release binaries have none of it.

## Build with the tag

```sh
go build -tags trace ./cmd/server
wails build -m -tags trace
```

A binary built without the tag refuses the settings below with an error, and records nothing.

## What is recorded

With `-trace_dir DIR` (ENV `FOUR_SOULS_TRACE_DIR`) the files are written when the process stops, on SIGTERM or Ctrl-C:

| File | What it shows |
|------|---------------|
| `trace.out` | Runtime trace: goroutines, scheduling, GC, syscalls, over time. Open with `go tool trace trace.out`. |
| `cpu.pprof` | Where the CPU time goes. |
| `heap.pprof` | Memory still in use at the stop, after a collection. |
| `allocs.pprof` | Every allocation since start: what churns memory. |
| `goroutine.pprof` | Goroutines still alive at the stop. Many left over after shutdown is a leak. |
| `block.pprof` | Time goroutines wait on channels and locks. |
| `mutex.pprof` | Contended mutexes. |
| `stats.csv` | Every 5 seconds: goroutines, live heap, heap objects, GC cycles. A rising column over a quiet server is a leak. |

With `-pprof_addr 127.0.0.1:6060` (ENV `FOUR_SOULS_PPROF_ADDR`) the same profiles are served live while the process runs.

## Profile a load test

1. Start the tagged server with recording on:

   ```sh
   ./server -port 4774 -records ./records -trace_dir ./traces -pprof_addr 127.0.0.1:6060
   ```

2. Run the load tool against it ([Load test](Load-test)):

   ```sh
   go run ./cmd/server_load -url ws://127.0.0.1:4774 -max-games 16
   ```

3. Take a live CPU profile during the load, 30 seconds:

   ```sh
   go tool pprof -http=:8080 "http://127.0.0.1:6060/debug/pprof/profile?seconds=30"
   ```

4. Stop the server with Ctrl-C, so the files are written.

5. Read the files:

   ```sh
   go tool pprof -top -cum traces/cpu.pprof
   go tool pprof -top -sample_index=alloc_space traces/allocs.pprof
   go tool pprof -top traces/block.pprof
   go tool trace traces/trace.out
   ```

## Look for

* **Hot spots**: the top of `cpu.pprof` by flat and by cumulative time. Start with `-cum`, then drill with `-list FUNC`.
* **Churn**: `allocs.pprof`; a big allocator inside a per-step path costs GC time.
* **Leaks**: `stats.csv` climbs while the games end; `goroutine.pprof` keeps a game's goroutines after shutdown.
* **Contention**: `block.pprof` and `mutex.pprof`; a hot lock shows as one function with most of the wait.

## Things to know

* The block and mutex profiles sample at a rate; the trace costs CPU. Record for minutes, not hours.
* `trace.out` grows fast under load: about 85 MB for 3.5 minutes with 8 games at once on a laptop.
* The live listener serves process internals. Bind it to `127.0.0.1` and reach it through a tunnel, or behind a proxy that guards it. Never on a public port.
* A profile is only as good as its load: run the same load before and after a change.
