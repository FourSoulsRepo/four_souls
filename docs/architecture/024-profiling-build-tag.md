# 24. Profiling behind a build tag

* **Status:** Accepted
* **Date:** 2026-10-29
* **Authors:** @HardDie

---

## Context

1. We need to find hot spots, leaks and slow paths in the server.
   1. The server must handle many games at once.
   2. The load tool (cmd/server_load) makes that load.
2. Profiling costs time and memory.
   1. The block and mutex profiles cost the most.
   2. Release binaries must not pay for it.

## Considered options

1. **Build tag `trace` with the Go standard library**
   1. Won. No new dependency; default builds carry none of the code.
2. **Always on, behind a flag**
   1. Lost. The code and its imports ship in every binary.
3. **An external agent or APM service**
   1. Lost. Adds a dependency and sends data out of the server.

## Decision

Use option 1.

1. Package `internal/tracing` holds the API.
   1. `tracing_on.go` has the code under `//go:build trace`.
   2. `tracing_off.go` refuses any request with `ErrNotBuilt`.
2. Each binary takes the same settings.
   1. The server: `-trace_dir` and `-pprof_addr`, also as ENV.
   2. The app: `FOUR_SOULS_TRACE_DIR` and `FOUR_SOULS_PPROF_ADDR`.
3. Recording writes files when the process stops.
   1. `trace.out`, `cpu.pprof`, `stats.csv`.
   2. `heap`, `allocs`, `goroutine`, `block` and `mutex` profiles.
4. `-pprof_addr` serves the live pages on its own mux.
   1. Bind it to localhost unless a proxy guards it.

## Consequences

### Positive

1. Release binaries are unchanged.
2. One tool set covers the app and the server.
3. `stats.csv` shows goroutines and live heap over time, for leaks.

### Negative and risks

1. A default binary refuses trace settings, so a wrong build fails loudly.
2. `trace.out` grows fast under load; record for minutes, not hours.
3. The live listener exposes process details; keep it off public ports.

### Neutral

1. The wiki page [Profiling](Profiling) says how to run it.
