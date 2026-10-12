# Four Souls — agent notes

1. Project
   1. Tabletop card game "Four Souls" as a desktop app.
   2. Network play, 2–4 players.
   3. Wails v2: Go backend, React + TypeScript frontend.
2. Stage
   1. Steps 1–7 done; engine, server, and the desktop client in `frontend/`.
   2. Ideas live in `docs/ideas.md`.
   3. Roadmap: `docs/roadmap.md`; detailed steps in `docs/roadmap/`.
   4. Decisions: `docs/architecture/` (see `INDEX.md`).
3. Layout
   1. `main.go`, `app.go`: Wails entry point and bindings.
   2. `frontend/`: React + TS (Vite).
   3. `docs/architecture/`: ADRs.
   4. `docs/use-cases/`: user-visible scenarios.
   5. `docs/wiki/`: package notes.
   6. `cmd/server`, `cmd/website`: server and wasm website entry points.
   7. `pkg/rules_engine`, `pkg/card_db`, `pkg/record`: own Go modules.
   8. `internal/`: app-only packages (protocol, version, legal, assets).
   9. `frontend/src/bridge/`: the only code that imports Wails.
   10. Full map: `docs/wiki/Repo-layout.md`; decision: ADR 003.
   11. Card data tools: separate repo `FourSoulsRepo/tools`.
   12. Never add download or scraping code to this repo.
   13. `e2e/`: Playwright UI tests; also makes screenshots.
   14. Wiki groups: players, server admins, developers.
   15. A new wiki page goes into one group in `Home.md` and `_Sidebar.md`.
4. Commands
   1. `wails dev`: run with hot reload.
   2. `wails build`: release binary to `build/bin/`.
   3. `go build ./...`: quick backend check.
   4. `make build vet test lint wasm`: all Go modules.
   5. `cd frontend && npm run typecheck && npm run lint`.
   6. Release-style build: `wails build -tags "embed cardimages"`.
   7. Linter setup and versions: `docs/wiki/Linters.md`.
   8. `make screenshots`: UI PNGs in `.cache/screenshots/` (not tracked).
   9. Use it to check UI changes; a real-window screenshot shows nothing.
   10. `go run ./cmd/cardgen -set b2`: engine card stubs from card_db.
   11. `go run ./cmd/cardgen -report`: card status report; rerun after card work.
   12. `make` alone prints help: every target and variable.
   13. `make fuzz`, `fuzz-sets`, `sims`: long local checks; see `docs/wiki/Fuzzing.md`.
   14. `make frontend`: typecheck, lint, Vitest.
   15. `make e2e`: Playwright UI tests in `e2e/`.
   16. Frontend testing rules: `docs/wiki/Frontend-tests.md`.
   17. `make internal`: play-test build with card images; see `docs/wiki/Internal-builds.md`.
5. Workflow
   1. Work follows `docs/roadmap/`, one sub-step at a time.
   2. Commit after every finished sub-step.
   3. Stop after every finished global step for owner review.
   4. Never start the next step without approval.
   5. One branch per global step: `step-NN-short-name`.
   6. After approval: merge into `main` with `--no-ff`, never squash.
   7. Before every commit: tests and lint must pass.
   8. Gate on the tool's own exit code: `make lint && git commit …`.
   9. Never pipe lint through `grep` or `tail` in a gate; it hides failures.
   10. Frontend changes also gate on `make frontend e2e`.
   11. Ideas from the owner are not tasks; never implement them unasked.
   12. Ideas go to `docs/ideas.md` and get planned in `docs/roadmap/`.
   13. Open questions wait for the owner at that roadmap step.
6. Licenses
   1. Code is MIT (`LICENSE`).
   2. `NOTICE` lists every third-party item and its license.
   3. Update `NOTICE` in the same commit as any dependency change.
   4. Card images live only in the private submodule `pkg/card_db/images`.
   5. Submodule repo: private `FourSoulsRepo/card_db` (images only).
   6. Public clones build without it; cards then show as text.
   7. Every change must still build and run without the private images.
   8. That covers cards, card backs, dice, and mats.
   9. Without them: text cards, plain backs, plain dice, color table.
   10. Check: `wails build -tags embed` (no `cardimages`) and `make e2e`.
