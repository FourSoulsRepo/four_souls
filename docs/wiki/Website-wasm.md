# Website wasm

The rules website (step 12, ADR 026) runs the rules engine in the browser as WebAssembly. There is no server: the page loads one gzipped `.wasm` and asks it.

## Code

* `cmd/website/main.go` — js/wasm only; puts the functions on `globalThis.fourSouls`.
* `internal/website/` — what each function answers, as plain Go, so `go test` covers it.
* `internal/cardinfo/` — builds `CardInfo` from card_db; the app's `Cards()` uses it too.

## Functions

Every function returns a JSON string. An error is `{"error": "…"}`; a Go panic becomes such an error, never an exception in the page.

| Function | Answer |
|----------|--------|
| `version()` | `{"app": "…", "engine": "…"}` |
| `sets()` | the set codes the site carries, e.g. `["b2","g2","fsp2","r"]` |
| `cards()` | every `CardInfo` of those sets; `image` is empty and `images` is `{}`, so cards show as text |
| `check(payload)` | the engine's `Answer` for a [situation payload](Situation-payload) |
| `notice()` | the fan-game notice lines, as the app shows them |

A refusal in `check` carries `rule`, the English `reason` for developers, and `say`: a message key with values (ADR 021) that the page shows in its language.

## Card sets

`internal/website.Sets` lists the set codes by hand: only finished sets, with no stub card left (`go run ./cmd/cardgen -report`). The `cmd/cardgen` test `TestWebsiteSetsAreFinished` fails if a listed set still has a stub. A set joins the list when step 14 finishes it.

## Build

```sh
make site
```

1. Builds `cmd/website` for `js/wasm` with `-tags embed` (card data inside, no images) and `-ldflags "-s -w"`.
2. Gzips it to `.cache/site/site.wasm.gz` and deletes the raw file; the page unpacks it with `DecompressionStream("gzip")`.
3. Prints the raw and gzip sizes and fails above the budget: 3,000,000 bytes of gzip.
4. Runs `cmd/website/smoke.mjs` with Node: it loads the `.gz` with Go's `wasm_exec.js` (`$(go env GOROOT)/lib/wasm/`), calls every function, and checks one refused and one allowed payload from `internal/website/testdata/`.

The output goes to `.cache/site/` because `.cache/` is git-ignored and Vite empties its own `frontend/dist-site/` on every build. CI runs `make site` in the `wasm` job.

`make wasm` still only checks that the website builds; without `-tags embed` its `cards()` answers an error.

Size on 2026-11-06, sets b2, g2, fsp2 and r: 8,874,685 bytes raw, 2,156,660 bytes gzip.
