# Four Souls — agent notes

1. Project
   1. Tabletop card game "Four Souls" as a desktop app.
   2. Network play, 2–4 players.
   3. Wails v2: Go backend, React + TypeScript frontend.
2. Stage
   1. Planning. No game code yet.
   2. Ideas live in `docs/ideas.md`.
   3. Roadmap comes after the ideas list is agreed.
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
