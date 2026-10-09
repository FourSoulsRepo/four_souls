# Repo layout

```
main.go, app.go        Wails desktop app and its bindings
cmd/server/            dedicated server (no UI)
cmd/website/           rules website, js/wasm only
internal/assets/       app files, embedded or read from disk
internal/legal/        fan-game notice used by app and server
internal/protocol/     network messages (never split out)
internal/version/      app version stamped at link time
pkg/rules_engine/      Go module: the rules engine (package rulesengine)
pkg/card_db/           Go module: card display data; images/ is a private submodule
pkg/record/            Go module: match record format
frontend/              React + TypeScript (Vite)
scripts/screenshots/   Playwright UI screenshots
frontend/src/bridge/   the only code that talks to Wails
docs/                  ideas, roadmap, ADRs, use cases, this wiki
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

## Make targets

| Target | What it does |
|--------|--------------|
| `make build` | `go build` in every module |
| `make vet` | `go vet` in every module, plus embed mode |
| `make test` | `go test` in every module, plus embed mode |
| `make lint` | golangci-lint in every module, plus embed mode |
| `make wasm` | builds the engine and the website for `js/wasm` |

On Ubuntu 24.04 pass `TAGS=webkit2_41`.
