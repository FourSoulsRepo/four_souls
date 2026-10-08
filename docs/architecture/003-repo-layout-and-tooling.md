# 3. Repo layout, modules, build modes and linters

* **Status:** Accepted
* **Date:** 2026-10-08
* **Authors:** @HardDie

---

## Context

1. One repo holds the app, the server and the website (A-04).
2. Some code may move to its own repo later (A-01, PR-02).
3. The rules engine must build for the browser (A-03, A-09).
4. Files must be editable without a rebuild during development (B-01).
5. Linters must be strict, with security rules (PR-03).

## Considered options

1. **One Go module for everything**
   1. Lost. Splitting a package out later breaks import paths.
2. **Separate Go modules under `pkg/`, joined by `go.work`**
   1. Won. Each module is ready to move; local work stays simple.
3. **`go.work` only, no `replace` in the app module**
   1. Lost. Builds with `GOWORK=off` would fail.
4. **Embed files always**
   1. Lost. Every asset change needs a rebuild.
5. **Build tag `embed` switches between embedded and disk files**
   1. Won. Releases embed; development reads from disk.
6. **ESLint 9 with `eslint-plugin-react`**
   1. Lost. ESLint 9 is end of life; the plugin lacks ESLint 10 support.
7. **ESLint 10 with `@eslint-react/eslint-plugin`**
   1. Won. Supported and strict type-checked presets exist.

## Decision

Use options 2, 5 and 7.

1. Module path prefix: `github.com/FourSoulsRepo`.
   1. Renaming it later is one search-and-replace.
2. Modules
   1. `.`: the app, module `github.com/FourSoulsRepo/four_souls`.
   2. `pkg/rules_engine`: package `rulesengine`, imported as `engine`.
   3. `pkg/card_db`: package `carddb`.
   4. `pkg/record`: package `record`.
   5. `go.work` joins them; the app also has `replace` lines.
3. Packages in the app module
   1. `internal/protocol`: network messages; never split out.
   2. `internal/version`: app version, stamped at link time.
   3. `internal/legal`: the fan-game notice.
   4. `internal/assets`: embedded or disk files.
4. Entry points
   1. `main.go`: the Wails app.
   2. `cmd/server`: the dedicated server.
   3. `cmd/website`: `js/wasm` only.
5. Versions
   1. App: `-X github.com/FourSoulsRepo/four_souls/internal/version.App=…`.
   2. Engine: constant `rulesengine.Version`.
6. Assets
   1. `-tags embed`: files from `internal/assets/files` in the binary.
   2. Default: `$FOUR_SOULS_ASSETS`, then `assets/` next to the binary.
   3. Then `internal/assets/files` (for `wails dev`).
7. Frontend
   1. Only `frontend/src/bridge` imports Wails bindings.
   2. ESLint enforces it.
8. Go linter
   1. golangci-lint v2 config `.golangci.yml`, one for all modules.
   2. `gosec` and other security linters; strict correctness linters.
   3. `depguard`: `pkg/` never imports the app or Wails.
   4. `depguard`: the engine never imports `os`, `net`, `time`, `math/rand`.
9. Frontend linter
   1. ESLint 10, typescript-eslint strict type-checked.
   2. `@eslint-react` strict type-checked, React Hooks.
   3. `eslint-plugin-security`, `eslint-plugin-no-unsanitized`.
   4. Warnings fail the run.
10. CI (GitHub Actions)
   1. Go build, vet, test, lint for every module.
   2. Engine and website built for `js/wasm`.
   3. Frontend typecheck, lint, build.
   4. App and server for Linux and Windows (x64, ARM), macOS universal.
11. Implements from ADR 001: Wails v2, wasm build in CI.

## Consequences

### Positive

1. `pkg/` modules can move to their own repos without code changes.
2. The engine cannot quietly gain OS, network or clock dependencies.
3. Asset edits show up after a restart, without a rebuild.

### Negative and risks

1. golangci-lint must be built with Go ≥ the local toolchain.
   1. An older binary fails with "export data version" errors.
2. Linux builds need `-tags webkit2_41` on Ubuntu 24.04.
3. CI is verified locally only until the repo has a remote.

### Neutral

1. Commands are in the `Makefile`: `build`, `vet`, `test`, `lint`, `wasm`.
2. Wiki pages: Repo layout, Build modes, Linters.
