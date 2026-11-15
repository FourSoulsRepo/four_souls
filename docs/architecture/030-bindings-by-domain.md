# 30. Bindings by domain

* **Status:** Accepted
* **Date:** 2026-11-15
* **Authors:** @HardDie

---

## Context

1. Owner, 2026-11-15: group the Wails bindings by domain (A-17).
   1. Model: `HardDie/ytmemchat_wails`, folder `bindings/`.
2. Today `main.go` binds `App`, `Host` and `Updater` from `package main`.
   1. `App` alone has 17 methods across eight areas.
   2. Its files sit in the repo root next to `main.go`.
3. Wails names the generated files after the Go package.
   1. Methods: `frontend/wailsjs/go/<package>/<Struct>.js`.
   2. Types: `export namespace <package>` in `models.ts`.
   3. Two Go packages with one name share one namespace there.
4. Taken namespaces: `carddb`, `cardinfo`, `lan`, `legal`, `main`.
   1. Also `protocol`, `settings`, `update`, `version`.
5. Wails binds every exported method of a bound struct.
   1. A lifecycle hook on a bound struct would reach the frontend.
6. Most bindings need only fixed values: the settings file, the card files.
   1. Only records, join and update need the Wails context.
   2. Update also needs a context that ends at shutdown.

## Considered options

1. **One package per domain under `bindings/`**
   1. Won. Each area is small, found by name, tested alone.
2. **One package per screen**
   1. Lost. Cards, settings and assets serve many screens.
3. **Keep `package main`, split files only**
   1. Lost. Generated paths and types stay under `main`.
4. **Bindings reach a big `App` core through interfaces**
   1. Lost. Here they share no mutable state; constructors suffice.

## Decision

Use option 1.

1. Packages:

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

2. Names avoid the taken namespaces.
   1. `prefs`, not `settings`; `updater`, not `update`.
   2. A new binding package checks `models.ts` first.
3. Each package has `New(…)`; it takes what it needs.
   1. Fixed values as arguments: settings file, card files, folders.
   2. The Wails context through an `api` interface: `Context() context.Context`.
4. `App` stays in `package main`, not bound.
   1. It keeps the Wails context from `OnStartup`.
   2. Its `Context()` ends at shutdown; background before startup.
   3. `main.go` builds every binding and lists them in `Bind`.
5. No exported lifecycle method on a bound struct.
   1. `OnShutdown` calls `Host.Stop`; the updater ends with `Context()`.
6. Player messages (ADR 021) go to `internal/msgerr`.
   1. It holds the message error type and the JSON list helper.
   2. Each binding package keeps its own keys and its error mapping.
7. `assets` keeps the card files and the URL prefix.
   1. It also serves the files: `cardImages` becomes its middleware.
   2. `internal/assets` keeps the embedded or on-disk file system.
8. `appPath`, `slotSettings` and the instance slot stay in `main`.
9. Each package has its tests; `wails.json` and key tests stay in `main`.
10. Only the bridge and the fake bridge change in the frontend.
    1. Its exports and types stay the same for screens.
    2. `fakeBridge.ts` fakes `window.go.<package>.<Struct>` to match.

## Consequences

### Positive

1. A binding is found by its domain, not in a 17-method `App`.
2. Generated files group the same way: `wailsjs/go/assets/Assets`.
3. Each package is tested without the others.

### Negative and risks

1. Two packages named `assets` (binding and `internal`); one gets an alias.
2. A wrong package name merges two `models.ts` namespaces.
3. Each binding is one more line in `main.go`.

### Neutral

1. The JS calls are the same; only their paths move.
2. `server-keys.json` still lists every key: `main` links every package.
