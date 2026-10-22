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
   16. `deploy/`: server setups: systemd, Docker, compose with Caddy (ADR 018).
4. Commands
   1. `wails dev -m`: run with hot reload.
   2. `wails build -m`: release binary to `build/bin/`.
   3. `go build ./...`: quick backend check.
   4. `make build vet test lint wasm`: all Go modules.
   5. `cd frontend && npm run typecheck && npm run lint`.
   6. Release-style build: `wails build -m -tags "embed cardimages"`.
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
   18. `make e2e-games`: whole random games, every set, 2–4 players (slow).
   19. Before a `cardimages` build or test: `go generate ./pkg/card_db/images`.
   20. `wails build` and `wails dev` always take `-m` (19.59).
       1. It skips `go mod tidy`.
       2. Tidy fails without the private images submodule.
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
   14. A step file may add its own rules in its header; they win.
   15. Step 15: stop after every sub-step, with screenshots, for review.
   16. Unlucky deals in e2e tests: fixed by host bans (14.10).
       1. A random deal can break a test that expects a plain board.
       2. Such a test bans the cards: `startGame(…, { bans })`.
       3. A deal-dependent failure is a bug; find why the card breaks it.
       4. A test that fails again on rerun is a real failure.
       5. Ban only the fewest cards the test needs.
       6. Each ban needs a reason: what the card changes that the test checks.
       7. Write the reason in a comment next to the bans.
       8. Never ban a card only because the test failed on it.
       9. Such a failure may be an engine or client bug; fix that instead.
   17. Bug reports go to `docs/roadmap/19-bugs.md`, one item each.
       1. Log the item before any code change.
       2. Fix it under the roadmap rules and that file's rules.
   18. Animations: one card, one animation per job (ADR 012).
       1. Rules and inventory: `docs/wiki/Animations.md`.
       2. Every animation has an `anim: <id>` tag and a row there.
       3. One that brings or shows a card claims it (`animations.ts`).
       4. A replaced animation loses its code, tag and row.
       5. e2e specs import `test` from `e2e/tests/fixtures.ts`.
       6. Its fixture fails any test that logs an animation clash.
   19. Keep a failed e2e run's evidence before the next run (19.12).
       1. A run start deletes `.cache/e2e/<port>/` and `e2e/test-results/`.
       2. Copy the match records and test results out first.
       3. A random test logs its seed when it starts, not only on a pass.
       4. `games.spec.ts` keeps its failed records in `.cache/e2e-games/`.
       5. Never rerun to "check again" before the evidence is saved.
   20. Window layers: one ordered list (ADR 015).
       1. Rules and inventory: `docs/wiki/Layers.md`.
       2. Numbers live only in `frontend/src/layers.css`, top first.
       3. Code uses `z-index: var(--z-<id>)`, never a bare number.
       4. Inside one element only: a number tagged `/* layer: local */`.
       5. A new window gets a row there, with why it covers what.
       6. Pairs that overlap get a test in `e2e/tests/layers.spec.ts`.
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
   10. Check: `wails build -m -tags embed` (no `cardimages`) and `make e2e`.
