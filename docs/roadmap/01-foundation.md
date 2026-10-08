# Step 1. Foundation

Goal: an empty but complete skeleton.
Everything builds, tests, and runs on all platforms.

Ideas: PR-01, PR-02, A-01, A-02, A-04, A-09, A-10, A-13, A-14, B-01, B-02, L-01, L-02.

---

### [ ] 1.1 Repo hygiene

1. Goal: the repo is ready for public work.
2. Tasks:
   1. Extend `.gitignore`: caches, scraper output, `build/bin`.
   2. Code license: MIT (`LICENSE`, added at project init).
   3. Third-party notices: `NOTICE` (added at project init).
   4. Add `.editorconfig`: tabs for Go, 2 spaces for TS.
3. Done when:
   1. `git status` is clean after a build.

### [ ] 1.2 Go modules and workspace

1. Goal: the module layout from A-01, A-11, A-12, A-13.
2. Tasks:
   1. `pkg/rules_engine`: own `go.mod`, package `rulesengine`.
   2. Imported under the alias `engine`.
   3. `pkg/card_db`: own `go.mod`.
   4. `pkg/record`: own `go.mod`.
   5. `internal/protocol`: package in the main module.
   6. `go.work` ties all modules together for local work.
   7. One module path prefix everywhere: `github.com/FourSoulsRepo`.
   8. Renaming the prefix later must be one search-and-replace.
3. Done when:
   1. `go build ./...` and `go test ./...` pass in every module.
   2. `pkg/` modules import nothing from the main module.

### [ ] 1.3 Entry points

1. Goal: the entry points from A-04 exist and run.
2. Tasks:
   1. `main.go`: the Wails app (exists).
   2. `cmd/server/main.go`: prints notice and versions, then exits.
   3. `cmd/website/main.go`: `js/wasm` build, exposes engine version.
   4. No relay yet (future).
3. Done when:
   1. All three build.
   2. `GOOS=js GOARCH=wasm go build ./cmd/website` passes.

### [ ] 1.4 Versions

1. Goal: app and engine versions in the main menu corner (A-02).
2. Tasks:
   1. App version stamped at link time; default `dev`.
   2. Engine version is a constant in the engine module.
   3. A Wails binding returns both.
   4. Server prints both on start.
3. Done when:
   1. A stamped build shows the stamped version.

### [ ] 1.5 Frontend skeleton and bridge

1. Goal: a clean React + TS app with a thin Wails bridge (A-14).
2. Tasks:
   1. Remove the Wails template demo.
   2. `frontend/src/bridge/`: the only place that imports Wails.
   3. Screens: splash, main menu (empty buttons).
   4. Lint and type check scripts.
3. Done when:
   1. No file outside `bridge/` imports Wails runtime or bindings.
   2. `npm run build` and lint pass.

### [ ] 1.6 Fan-game notice

1. Goal: legal notice on every start (L-01, L-02).
2. Tasks:
   1. One shared notice text in Go.
   2. Names: Edmund McMillen (designer), Maestro Media (publisher).
   3. Copy the exact rights line from the official rulebook or box.
   4. Links to official site and store.
   5. Splash shows it for a few seconds; links open the browser.
   6. Server prints it with links.
3. Done when:
   1. Splash and console show the same text.
   2. Use case written: app start with notice.
4. **Owner:** confirm the final notice text.

### [ ] 1.7 Embedded and external assets

1. Goal: two build types (B-01).
2. Tasks:
   1. One asset access interface in Go.
   2. `-tags embed`: files packed into the binary.
   3. Default: files read from a folder next to the app.
   4. `wails dev` uses external mode.
   5. Missing folder is not fatal; the app still starts.
3. Done when:
   1. Both modes return the same test file.
   2. Tests cover both build tags.

### [ ] 1.8 CI

1. Goal: GitHub Actions checks every push (A-09, B-02).
2. Tasks:
   1. Go vet and tests for every module.
   2. Engine and website built for `js/wasm`.
   3. Frontend type check, lint, build.
   4. Wails build matrix: Windows, Linux, macOS × x64, ARM.
   5. Build artifacts uploaded (not released).
3. Done when:
   1. A push runs all jobs green.

### [ ] 1.9 Foundation docs

1. Goal: the layout is documented.
2. Tasks:
   1. ADR: repo layout, modules, entry points, asset modes.
   2. Accept ADR 001 parts that step 1 implements.
   3. Wiki: `Home`, `_Sidebar`, `_Footer`, "Repo layout".
   4. Update `CLAUDE.md` and `README.md`.
3. Done when:
   1. A new agent can find every folder's purpose in docs.
   2. `NOTICE` lists every dependency added in step 1.
