# Developer setup

One script installs every tool `make` needs, on macOS and Linux (16.8). `make doctor` then checks them. Windows is not covered yet.

## Run it

From the repository root:

| System | Command |
|---|---|
| macOS (Homebrew) | `scripts/setup.sh --profile developer` |
| Debian, Ubuntu (apt) | `scripts/setup.sh --profile developer` |
| Fedora (dnf) | `scripts/setup.sh --profile developer` |

Use `--profile translator` for the smaller set. Then:

```sh
make doctor                      # the developer profile
make doctor PROFILE=translator   # only what a translator needs
```

The script first checks every tool and prints a table: found version, wanted version, status. Then it prints its plan, every command it will run, `sudo` ones included, and asks once. It changes nothing before that answer.

| Flag | What it does |
|---|---|
| `--profile translator` or `developer` | which tools; `developer` is the default |
| `--check` | only the table, no installs; `make doctor` runs this |
| `--yes` | no question; for scripts and containers |

A second run skips what is already there at the right version, so it is safe to run again after a `git pull`.

## What each profile installs

| Tool | translator | developer | macOS | Linux |
|---|---|---|---|---|
| git, make | yes | yes | Xcode Command Line Tools | apt / dnf |
| Go (version from `go.mod`, or newer) | yes | yes | `brew install go` | archive from go.dev into `~/.local/go` |
| Node.js (`NODE_VERSION`, or newer) | yes | yes | `brew install node` | archive from nodejs.org into `~/.local/node` |
| npm deps of `frontend/` and `e2e/` | yes | yes | `npm ci` | `npm ci` |
| Playwright Chromium and WebKit | yes | yes | `npx playwright install` | same, plus their system libraries; Fedora: Chromium only |
| C compiler, GTK 3, WebKitGTK 4.1 dev | | yes | Xcode Command Line Tools | apt / dnf |
| Wails CLI (`WAILS_VERSION`) | | yes | `go install` into `$(go env GOPATH)/bin` | same |
| golangci-lint (`GOLANGCI_LINT_VERSION`) | | yes | `go install` into `.cache/bin` | same |
| shellcheck (optional) | | yes | `brew install shellcheck` | apt / dnf |

Translators need Go too: the UI tests behind `make screenshots` start the real game server with `go run`.

Versions live in one file, `scripts/versions.env`; the script, `make doctor` and CI read it. Go's version comes from `go.mod`, Playwright's from `e2e/package.json`. Go and Node pass when they are newer; the Wails CLI, golangci-lint and Playwright must match exactly. golangci-lint must also be built with a Go at least as new as the repo's ([Linters](Linters)); `make doctor` checks that too.

On Linux, Go and Node come from their official archives, checked against the published SHA-256 sums, into `~/.local` without `sudo`: distribution packages are often too old. Add the folders the script names at the end to `PATH` in your shell profile, for example:

```sh
export PATH="$HOME/.local/go/bin:$HOME/.local/node/bin:$HOME/go/bin:$PATH"
```

## What each make target needs

| Target | Needs |
|---|---|
| `help`, `doctor` | make, bash |
| `build`, `vet`, `test` | Go; Linux: C compiler, GTK 3 and WebKitGTK 4.1 dev packages, `TAGS=webkit2_41`; macOS: Xcode Command Line Tools |
| `wasm`, `fuzz`, `fuzz-*`, `sims` | Go |
| `lint` | golangci-lint, Go, plus what `i18n` needs; shellcheck if found |
| `i18n` | Node, `frontend/` npm deps |
| `frontend` | Node, `frontend/` npm deps |
| `e2e`, `e2e-games` | Node, npm deps, Playwright Chromium and WebKit, Go (the test server) |
| `screenshots` | the same, Chromium only |
| `internal` | Wails CLI, Go, Node, what `build` needs, the private images submodule |
| `server-archives` | Go, `zip`, macOS (`lipo` for the universal binary) |

`make e2e` and `make screenshots` also run `npm install` and `playwright install` in `e2e/` themselves; the script does it once ahead.

On Linux, the Go targets that link the app take `TAGS=webkit2_41`: `make build vet test lint TAGS=webkit2_41`, and `make internal TAGS=webkit2_41`. Without it Wails looks for WebKitGTK 4.0, which current distributions no longer ship.

## By hand

When the script cannot run (another distribution, no `sudo`), install these yourself, then run `make doctor`:

1. git and make.
2. Go, at least the version on the `go` line of `go.mod`: [go.dev/dl](https://go.dev/dl/).
3. Node.js, at least `NODE_VERSION` from `scripts/versions.env`: [nodejs.org](https://nodejs.org/).
4. `cd frontend && npm ci`, then `cd e2e && npm ci`.
5. `cd e2e && npx playwright install --with-deps chromium webkit` (on apt systems; elsewhere install the libraries Playwright names when a browser fails to start).
6. Developers only:
   1. A C compiler and pkg-config; on Linux the GTK 3 and WebKitGTK 4.1 development packages (Arch: `base-devel gtk3 webkit2gtk-4.1`).
   2. `go install github.com/wailsapp/wails/v2/cmd/wails@<WAILS_VERSION>`.
   3. `GOBIN=$PWD/.cache/bin GOWORK=off go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@<GOLANGCI_LINT_VERSION>`.
   4. shellcheck, if you change shell scripts.

The private card images submodule is separate: without access the app builds and shows cards as text ([Build modes](Build-modes)).

## Troubleshooting

* **"Unsupported system"**: the script knows Homebrew, apt and dnf, on x86_64 and arm64. Install by hand as above; `make doctor` still works anywhere.
* **"not on PATH"**: the tool is in `~/.local/go/bin`, `~/.local/node/bin` or `~/go/bin`, but your shell does not look there. Add the folder to `PATH` and open a new terminal.
* **golangci-lint "wrong (built with goX < goY)"**, or lint errors like "export data version 4 is greater than maximum supported version": the binary is older than your Go. Run the script again; it rebuilds `.cache/bin/golangci-lint`. `make lint` uses that one when it is there.
* **Go "wrong" although a newer Go is installed**: an older `go` comes first on `PATH`. `which -a go` shows them all.
* **npm deps "stale"**: `package-lock.json` changed since the last install, for example after a `git pull`. Run the script again, or `npm ci` in that folder.
* **"browser libs" missing**, or a browser fails to start on Linux with a missing `.so` file: the browsers' system libraries are missing. Run the script again; on apt systems it runs `npx playwright install-deps chromium webkit`, which asks for `sudo` itself. `make doctor` asks Playwright (apt) or rpm (dnf) whether they are there.
* **WebKit on Fedora**: Playwright's WebKit build links Ubuntu's libraries (ICU 74, GStreamer and more), so it does not start there. The script installs Chromium only, and `make doctor` shows WebKit as skipped. `make screenshots` and `make i18n` work; `make e2e`, which runs WebKit too, needs Debian or Ubuntu (a container works). `cd e2e && npx playwright test --project chromium` runs the Chromium half.
* **`make build` on Linux fails on `webkit2gtk-4.0`**: add `TAGS=webkit2_41`.
* **macOS asks for the Xcode Command Line Tools**: run `xcode-select --install`, then the script again.
* **Homebrew is missing**: install it from [brew.sh](https://brew.sh), then run the script again. The script never installs Homebrew itself.
