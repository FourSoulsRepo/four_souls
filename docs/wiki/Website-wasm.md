# Website wasm

The rules website (step 12, ADR 026) runs the rules engine in the browser as WebAssembly. There is no server: the page loads one gzipped `.wasm` and asks it.

## Code

* `cmd/website/main.go` — js/wasm only; puts the functions on `globalThis.fourSouls`.
* `internal/website/` — what each function answers, as plain Go, so `go test` covers it.
* `internal/cardinfo/` — builds `CardInfo` from card_db; the app's `Cards()` uses it too.
* `frontend/site.html`, `frontend/src/site/` — the page: a second Vite entry next to the app (12.2).

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

## The page

`frontend/src/site/` is a small React page. It shares `t()`, `say()`, the locale files and the plain card helpers with the app, never the bridge or the store: ESLint's `no-restricted-imports` rejects them in `src/site/**`. Its texts are keys of the area `site` in `frontend/src/locales/en.json`; the unused-key check of `make i18n` reads `src/site/` like the rest of `src/`.

| File | What it does |
|------|--------------|
| `main.tsx`, `Site.tsx` | entry; header, the view the route names, footer with the fan-game notice and versions |
| `engine.ts` | loads `wasm_exec.js` and `site.wasm.gz`, with progress; types of the answers |
| `route.ts` | the routes in the URL fragment |
| `link.ts` | the `<data>` of a check link |
| `CheckView.tsx` | payload box, Check, Copy link, Copy payload, examples, the answer |
| `markdown.ts`, `RulesText.tsx` | the own Markdown reader and its React drawing |
| `rules.ts`, `RulesView.tsx` | the rules files, the topic index, the rules page |
| `site.css` | system fonts; light and dark by the browser's setting |

### Loading the engine

1. The page needs `DecompressionStream`; without it, it says to update the browser.
2. It adds Go's `wasm_exec.js` as a script, then fetches `site.wasm.gz`.
3. A progress bar counts the bytes against `Content-Length`.
4. The page unpacks the stream with `DecompressionStream("gzip")` only when its first bytes are gzip's (`1f 8b`). Some hosts send `.gz` files with `Content-Encoding: gzip` (Vite's dev server and `vite preview` do), and then the browser has unpacked it already; the progress bar then has no total.
5. `WebAssembly.instantiateStreaming` compiles it as it arrives; `go.run` sets `globalThis.fourSouls`.
6. The page reads `version()`, `cards()` and `notice()` once. A failure shows its reason; the rules page works without the engine.

### Routes

All routes live in the URL fragment, so GitHub Pages needs no rewrites and a payload never reaches a server.

| Route | View |
|-------|------|
| `#/<lang>/check/<data>` | the situation of the link, checked at once |
| `#/<lang>/check` | an empty payload box, with example buttons (`internal/website/testdata/`) |
| `#/<lang>/rules/<rule ID>` | the rules page, scrolled to that rule and highlighted |
| `#/<lang>/rules` | the rules page from the top |

`<lang>` is a code of two or three letters; without it, the page takes the browser's language among the locale files, else English, and rewrites the address. A language with no file shows English (`t()` falls back). FAQ (12.6) and stack (12.4) routes come later; anything else opens the check view.

### Check links

`<data>` is the payload JSON without spaces, compressed with `CompressionStream("deflate-raw")`, in base64url with no `=` padding. `decodePayload` reverses it; a cut link says it is broken. Check also replaces the address with the payload's link, so a reload keeps it. Copy link copies the link and shows it in a field; above 2000 characters the page warns that some chats cut long links. Copy payload always works. `link.test.ts` round-trips it; the e2e test opens a link the page copied.

### Answers

Allowed or Refused, with the payload's `title`. A refusal shows `say` through `say()`, as the app shows server messages: cards by their names from `cards()` (text only), players as "Player N". The rule ID links to the rules page. With `expect`, the page says whether the answer matches, and lists `mismatch`. An allowed move lists its events: player, event kind, label, card, amount.

### Rules text

Vite imports `pkg/rules_engine/docs/rules/*.md` as text, except `README.md` and `open-questions.md`. `markdown.ts` reads headings, paragraphs, numbered and bulleted lists nested by indent, tables, bold, italics, code and links; `RulesText.tsx` draws them as React elements, never as HTML. A `**R-XXX-NN**` at the start of a list item is the anchor `rule-R-XXX-NN`; any other rule ID links to it. A link to another rules file opens that file's first rule; one to a file the site lacks stays plain text. The topic index is built from each file's title, in README's order. `markdown.test.ts` checks that every rule ID in the real files has exactly one anchor and that every mentioned ID exists.

### Tests

* Vitest: `src/site/*.test.ts` in `make frontend`.
* Playwright: `e2e/tests/site.spec.ts` against `vite preview --mode site` on port 5197, which serves `frontend/dist-site/`. `make e2e` runs `make site` first, so it tests the current build. Run without a build, the spec skips and says to run `make site`.
* `make screenshots` saves `screen-site-check`, `screen-site-rules` and their `-phone` forms (390 px wide).

## Build

```sh
make site
```

1. Builds `cmd/website` for `js/wasm` with `-tags embed` (card data inside, no images) and `-ldflags "-s -w"`.
2. Gzips it to `.cache/site/site.wasm.gz` and deletes the raw file; the page unpacks it with `DecompressionStream("gzip")`.
3. Prints the raw and gzip sizes and fails above the budget: 3,000,000 bytes of gzip.
4. Runs `cmd/website/smoke.mjs` with Node: it loads the `.gz` with Go's `wasm_exec.js` (`$(go env GOROOT)/lib/wasm/`), calls every function, and checks one refused and one allowed payload from `internal/website/testdata/`.
5. Copies that `wasm_exec.js` into `.cache/site/`; it must match the Go that built the `.wasm`.
6. Runs `vite build --mode site`: `base` is `/four_souls/` (the GitHub Pages path), the public folder is `.cache/site/`, the output `frontend/dist-site/`. Then `site.html` becomes `index.html`.

To work on the page: `make site` once, then `cd frontend && npx vite --mode site` and open `http://localhost:5173/four_souls/site.html`. `npx vite preview --mode site` serves the built `dist-site/` at `/four_souls/`. The app's build (`npm run build`, `wails build`) is unchanged.

The `.wasm` goes to `.cache/site/` because `.cache/` is git-ignored and Vite empties its own `frontend/dist-site/` on every build. CI runs `make site` in the `wasm` job.

`make wasm` still only checks that the website builds; without `-tags embed` its `cards()` answers an error.

Size on 2026-11-06, sets b2, g2, fsp2 and r: 8,874,685 bytes raw, 2,156,660 bytes gzip.
