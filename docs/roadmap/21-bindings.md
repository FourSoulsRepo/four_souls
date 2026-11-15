# Step 21. Bindings by domain

Goal: Wails bindings grouped by domain, one package each under `bindings/` (A-17).
Model: `HardDie/ytmemchat_wails`, folder `bindings/`.

Ideas: A-17.

## Rules for this step (owner, 2026-11-15)

1. A refactor: nothing the player sees changes.
   1. The bridge API (`frontend/src/bridge/index.ts` exports) stays the same.
   2. Screens need no change.
   3. Every e2e test passes unchanged.
   4. `fakeBridge.ts` changes only its `window.go` paths (owner, 2026-11-15).
2. One domain move per sub-step; the app builds and runs after each.
3. A move takes its tests along (`*_test.go` next to the code).
4. Gates: `make lint`, `make frontend e2e`.
5. Check without images too: `wails build -m -tags embed`.
6. Branch: the step 18 branch, not its own (owner, 2026-11-15).
7. Design: ADR 030; each sub-step follows it.

## Today (2026-11-15)

`main.go` binds three structs from `package main`: `Bind: []any{app, host, updater}`.

| Struct | Method | File |
|--------|--------|------|
| `App` | `Versions`, `Notice` | `app.go` |
| `App` | `Settings`, `SaveSettings` | `app.go` |
| `App` | `CardSets`, `Challenges`, `DeckRatios` | `app.go` |
| `App` | `Cards` | `cards.go` |
| `App` | `Backs`, `Dice`, `Mats`, `Fonts`, `Icons`, `Sounds` | `backs.go`, `dice.go`, `mats.go`, `fonts.go`, `icons.go`, `sounds.go` |
| `App` | `FindGames` | `lan.go` |
| `App` | `CheckSeat`, `ServerStatus` | `servers.go` |
| `App` | `SaveRecord` | `record.go` |
| `Host` | `Start`, `Stop`, `Info` | `host.go` |
| `Updater` | `Check`, `Download`, `ApplyAndQuit`, `OpenReleases` | `updater.go` |

Shared helpers: `msgError` (`messages.go`), `imageURL`, `soundURLs`, `appContext`, `appPath`.

---

### [x] 21.1 ADR: bindings by domain

1. Goal: one ADR fixes the packages, names and wiring.
2. Tasks:
   1. Written in the main conversation (CLAUDE.md 5.21.2).
   2. Package list and the methods in each; draft below.
   3. Wiring: each package takes a small `api` interface, as ytmemchat does.
   4. `App` in `package main` keeps lifecycle and shared state only.
   5. Where `msgError` and the URL helpers go.
   6. Who gets `startup` and `shutdown` (today `App`, `Updater`, `Host`).
   7. `cardImages` middleware: stays in `main` or moves to `assets`.
   8. Update ADR 003 (layout) in the same change.
3. Draft packages (to settle in the ADR):

   | Package | Struct | Methods |
   |---------|--------|---------|
   | `about` | `About` | `Versions`, `Notice` |
   | `prefs` | `Prefs` | `Settings`, `SaveSettings` |
   | `assets` | `Assets` | `Cards`, `Backs`, `Dice`, `Mats`, `Fonts`, `Icons`, `Sounds` |
   | `setup` | `Setup` | `CardSets`, `Challenges`, `DeckRatios` |
   | `join` | `Join` | `FindGames`, `CheckSeat`, `ServerStatus` |
   | `host` | `Host` | `Start`, `Stop`, `Info` |
   | `records` | `Records` | `SaveRecord` |
   | `updater` | `Updater` | `Check`, `Download`, `ApplyAndQuit`, `OpenReleases` |

4. Known constraint: Wails names TS model namespaces after the Go package.
   1. Taken today: `carddb`, `cardinfo`, `lan`, `legal`, `main`, `protocol`, `settings`, `update`, `version`.
   2. A binding package with one of these names merges with it in `models.ts`.
   3. Hence `prefs` (not `settings`) and `updater` (not `update`) in the draft.
5. **Owner** questions:
   1. Package names: the draft above, or others?
   2. `assets`: one package, or `cards` apart from the media skins?
6. Done when:
   1. The ADR is Accepted and its row is in `INDEX.md`.
7. Result (2026-11-15)
   1. ADR 030, Accepted; the draft packages as they stand.
   2. Owner answers: the draft names; one `assets` package.
   3. `App` stays in `main`, not bound; it gives `Context()` only.
   4. Fixed values go to each `New`; no shared `api` beyond the context.
   5. Player messages: `internal/msgerr`; keys stay in each package.
   6. ADR 003 item 4.4: the `bindings/` folder.

### [x] 21.2 Host and updater

1. Goal: the two structs that already stand alone move first.
2. Tasks:
   1. `bindings/host`, `bindings/updater` with their tests.
   2. `main.go` binds them from there.
   3. Bridge imports the new `wailsjs/go/<package>/<Struct>` paths.
   4. Shared helpers they need go where the ADR says.
3. Done when:
   1. `host.go`, `updater.go` are gone from the repo root.
   2. Gates pass.
4. Result (2026-11-16)
   1. `bindings/host`: `host.New(records, saves)`; `main` passes the `appPath` folders.
      1. Its own key `error.hosting`; `OnShutdown` calls `Stop`.
   2. `bindings/updater`: `updater.New(api)`; `api` gives `Context()`.
      1. No lifecycle methods; runtime calls take that context.
      2. Tests replace `emit`, `quit`, `open`.
   3. `internal/msgerr`: `Error`, `New`, `JSON`, with tests.
      1. `main` keeps `playerError` for settings and LAN.
   4. Departure: `Context()` is on `App.wails`, an unbound `wailsContext` field.
      1. Not `App.Context()`: `App` is bound until 21.5.
      2. Background before startup; ends at shutdown.
      3. `FindGames`, `SaveRecord`, `CheckSeat`, `ServerStatus` use it.
      4. In 21.5 it can fold into `App` (ADR 030 item 4).
   5. Departure: `HostInfo` is now `host.Info` (revive: no stutter).
      1. The bridge keeps its own `HostInfo`; its exports are unchanged.
   6. Generated: `wailsjs/go/host/Host`, `wailsjs/go/updater/Updater`.
      1. `models.ts`: new namespace `host`; no merge.
   7. Rule change (owner, 2026-11-15): `fakeBridge.ts` fakes `window.go.<package>.<Struct>`.
      1. Step rule 1.4 and ADR 030 item 10.2.
      2. Only the paths and the `go` type changed there.
   8. Bridge: two import paths; exports unchanged.
   9. `Makefile` `internal`: chmod covers `wailsjs/go/*/*`.
   10. Docs: `Client.md`, `Repo-layout.md`, UC-03, UC-43, use-case index.
   11. Gates: `make lint`, `make test`, `make frontend`, `wails build -m -tags embed`.
       1. `make e2e`: 562 passed, 2 failed, 18 skipped.
       2. Failures, chromium only: `aim.spec.ts:33`, `lobby.spec.ts:13`; logged as 19.161, 19.162.
       3. Rerun of both: 2 of 2 passed.
       4. A run before the fake bridge change: 446 failed (`window.go.main.Host` gone).

### [ ] 21.3 About and prefs

1. Goal: version, notice and settings bindings leave `App`.
2. Tasks:
   1. `bindings/about`, `bindings/prefs` with their tests.
   2. `App` gives them what they need through `api`.
   3. Bridge imports the new paths.
3. Done when:
   1. `App` has no `Versions`, `Notice`, `Settings`, `SaveSettings`.
   2. Gates pass.

### [ ] 21.4 Assets

1. Goal: card and media bindings leave `App`.
2. Tasks:
   1. Cards, backs, dice, mats, fonts, icons, sounds into `bindings/assets`.
   2. Their tests move along.
   3. Bridge imports the new paths.
3. Done when:
   1. `cards.go`, `backs.go`, `dice.go`, `mats.go`, `fonts.go`, `icons.go`, `sounds.go` are gone from the root.
   2. Cards still show as text without the private images.
   3. Gates pass.

### [ ] 21.5 Setup, join, records

1. Goal: the last `App` bindings move out.
2. Tasks:
   1. `bindings/setup`, `bindings/join`, `bindings/records` with their tests.
   2. Bridge imports the new paths.
3. Done when:
   1. `App` is not bound; it keeps only the Wails context.
   2. `main.go` binds only `bindings/*` structs.
   3. Gates pass.

### [ ] 21.6 Docs

1. Goal: docs match the new layout.
2. Tasks:
   1. `docs/wiki/Repo-layout.md`: the `bindings/` folder.
   2. `docs/wiki/Client.md`: binding names (e.g. `App.Cards`, `Updater`).
   3. `CLAUDE.md` 3.1: `main.go`, `app.go` and `bindings/`.
   4. ADRs that name a moved file or struct: edit in place.
3. Done when:
   1. `grep` finds no stale `App.<Method>` binding names in `docs/`.
