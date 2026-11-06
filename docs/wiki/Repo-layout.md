# Repo layout

```
main.go, app.go        Wails desktop app and its bindings
cmd/server/            dedicated server (no UI)
cmd/server_load/       load test: plays whole games against a running server
cmd/website/           rules website, js/wasm only
cmd/cardgen/           writes engine card stubs from card_db data
internal/assets/       app files, embedded or read from disk
internal/cardinfo/     card display info from card_db, for app and website
internal/instance/     one slot per running app copy
internal/lan/          finds hosted games on the LAN (ADR 009)
internal/legal/        fan-game notice used by app and server
internal/protocol/     network messages (never split out)
internal/server/       game rooms and the WebSocket transport (ADR 006)
internal/settings/     the player's app settings file
internal/version/      app version stamped at link time
internal/website/      what the website's .wasm answers, testable in plain Go
pkg/rules_engine/      Go module: the rules engine (package rulesengine)
pkg/rules_engine/cards/  card sets: one package per set, one file per card
pkg/card_db/           Go module: card display data; images/ is a private submodule
pkg/record/            Go module: match record format
frontend/              React + TypeScript (Vite)
e2e/                   Playwright UI tests and screenshots
frontend/src/bridge/   the only code that talks to Wails
frontend/src/site/     the rules website's page: a second Vite entry (site.html)
docs/                  ideas, roadmap, ADRs, use cases, this wiki
release/               READMEs for the release archives; server.sh packs the servers
deploy/                server setups: systemd unit, Dockerfile, compose with Caddy
.github/               CI, release workflow, shared build and card images actions
scripts/               setup.sh installs the tools; versions.env holds their versions
```

## Go modules

Every folder under `pkg/` is its own Go module with its own `go.mod`, so it can later move to a separate repository unchanged. `go.work` at the root joins all modules for local work, and the app's `go.mod` also has `replace` lines so builds work with `GOWORK=off`.

The module path prefix is `github.com/FourSoulsRepo`. The rules engine is imported under the alias `engine`:

```go
import engine "github.com/FourSoulsRepo/rules_engine"
```

Rules enforced by the linter (see [Linters](Linters)):

* `pkg/` modules never import the app module or Wails.
* The rules engine never imports `os`, `net`, `time` or `math/rand`, so it stays deterministic and builds for the browser.

## Frontend bridge

`frontend/src/bridge/index.ts` wraps the generated Wails bindings (`frontend/wailsjs/`). Screens import from the bridge only; ESLint rejects any other import of `wailsjs`. A move to Wails v3, or a browser build for the web client, replaces the bridge and nothing else.

The rules website (`frontend/site.html`, `frontend/src/site/`) shares `t()`, the locale files and the plain helpers with the app, but has no Wails: ESLint rejects any import of the bridge or the store from `src/site/` ([Website wasm](Website-wasm)).

## Make targets

| Target | What it does |
|--------|--------------|
| `make build` | `go build` in every module |
| `make vet` | `go vet` in every module, plus embed mode |
| `make test` | `go test` in every module, plus embed mode |
| `make lint` | golangci-lint in every module, plus embed mode |
| `make wasm` | builds the engine and the website for `js/wasm` |
| `make site` | the rules website: `site.wasm.gz` into `.cache/site/` (size budget, Node smoke test), then the page into `frontend/dist-site/` ([Website wasm](Website-wasm)) |
| `make frontend` | typecheck, lint and Vitest for the frontend |
| `make e2e` | `make site`, then Playwright UI tests (the app and the rules website) |
| `make internal` | play-test app with card images; `PLATFORM=windows/amd64` |
| `make screenshots` | `make e2e`, saving UI screenshots into `.cache/screenshots/` |
| `make doctor` | checks the tools each target needs ([Developer setup](Developer-setup)) |
| `make fuzz-engine`, `fuzz-focused`, `fuzz-sets` | one fuzz target for `FUZZTIME` (default 10m) |
| `make fuzz` | `fuzz-engine`, then `fuzz-sets` |
| `make sims` | `SIMS` (default 3000) whole random games |

`make` alone prints this list with the variables.

On Ubuntu 24.04 pass `TAGS=webkit2_41`.
