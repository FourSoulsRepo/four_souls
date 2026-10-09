# Linters

Both linters are strict and include security rules. CI fails on any finding.

## Go: golangci-lint

Config: `.golangci.yml` (golangci-lint config v2), shared by every module.

```sh
make lint
```

Enabled groups:

* Security: `gosec`, `bodyclose`, `noctx`, `errorlint`, `errchkjson`, `depguard` and others.
* Correctness: `errcheck` (with type assertions and blank checks), all of `govet`, `staticcheck`, `nilerr`, `exhaustive`, `forcetypeassert`, `contextcheck`, `containedctx` and others.
* Style: `revive`, `misspell`, `godot`, `whitespace`; formatters `gofumpt` and `goimports`.

`depguard` encodes architecture rules: `pkg/` modules never import the app or Wails, and the rules engine never imports `os`, `net`, `time` or `math/rand`.

golangci-lint must be built with a Go version at least as new as your local toolchain. An older binary fails with errors like "export data version 4 is greater than maximum supported version". Install a matching one into the project cache and point make at it:

```sh
GOBIN=$PWD/.cache/bin GOWORK=off go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0
make lint GOLANGCI=$PWD/.cache/bin/golangci-lint
```

## Frontend: ESLint

Config: `frontend/eslint.config.js` (ESLint 10 flat config).

```sh
cd frontend
npm run lint
```

Presets: typescript-eslint strict and stylistic type-checked, `@eslint-react` strict type-checked, React Hooks, `eslint-plugin-security`, `eslint-plugin-no-unsanitized`. Extra rules:

* `eval`, `new Function`, string timers and `javascript:` URLs are errors.
* `dangerouslySetInnerHTML` is an error.
* Only `src/bridge/` may import the generated Wails bindings.

Warnings fail the run (`--max-warnings 0`).
