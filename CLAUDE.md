# Four Souls — agent notes

1. Project
   1. Tabletop card game "Four Souls" as a desktop app.
   2. Network play, 2–4 players.
   3. Wails v2: Go backend, React + TypeScript frontend.
2. Stage
   1. Planning. No game code yet.
   2. Ideas live in `docs/ideas.md`.
   3. Roadmap: `docs/roadmap.md`; detailed steps in `docs/roadmap/`.
   4. Decisions: `docs/architecture/` (see `INDEX.md`).
3. Layout
   1. `main.go`, `app.go`: Wails entry point and bindings.
   2. `frontend/`: React + TS (Vite).
   3. `docs/architecture/`: ADRs.
   4. `docs/use-cases/`: user-visible scenarios.
   5. `docs/wiki/`: package notes.
4. Commands
   1. `wails dev`: run with hot reload.
   2. `wails build`: release binary to `build/bin/`.
   3. `go build ./...`: quick backend check.
5. Workflow
   1. Work follows `docs/roadmap/`, one sub-step at a time.
   2. Commit after every finished sub-step.
   3. Stop after every finished global step for owner review.
   4. Never start the next step without approval.
   5. One branch per global step: `step-NN-short-name`.
   6. After approval: merge into `main` with `--no-ff`, never squash.
6. Licenses
   1. Code is MIT (`LICENSE`).
   2. `NOTICE` lists every third-party item and its license.
   3. Update `NOTICE` in the same commit as any dependency change.
   4. Card images are never committed (`pkg/card_db/images/`).
