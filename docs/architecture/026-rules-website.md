# 26. Rules website

* **Status:** Accepted
* **Date:** 2026-11-06
* **Authors:** @HardDie

---

## Context

1. Step 12: check rules in the browser; no server (A-03).
2. Owner rules for step 12 (2026-11-06):
   1. English and Russian.
   2. The app gets the Russian interface too, if that keeps the code simpler.
   3. Card texts stay English; card names follow the language.
   4. The site carries every finished set.
   5. It deploys on a release tag.
   6. The repo is public; GitHub Pages is free for it.
3. What exists:
   1. `cmd/website` builds for `js/wasm`; it exports only the version.
   2. `engine.RunSituation` answers a payload, with a rule ID.
   3. A refusal carries `Say`: a message key with values (ADR 021).
   4. Rules text: 12 Markdown files in `pkg/rules_engine/docs/rules/`.
   5. `frontend/` has `t()`, the locale loader and the card component.
   6. `cards.go` builds `CardInfo` from card_db for the app.
4. Size, measured 2026-11-06, every set, `-ldflags "-s -w"`:
   1. 7.9 MB raw.
   2. 2.0 MB gzip.
   3. 1.4 MB brotli.
5. Card images, fonts and sounds are private; the site never has them.
6. The app copies situation links later (11.3); the site sets the format.
7. The app has 959 interface keys in `en.json`; none are in Russian yet.
8. It is unknown whether GitHub Pages compresses `.wasm` files.

## Considered options

1. **Site code: a second Vite entry in `frontend/`**
   1. Won. Reuses `t()`, locales, the card component, lint and tests.
2. **Site code: its own small app**
   1. Lost. A second copy of `t()` and of card rendering.
3. **Rules text: own small Markdown renderer in TypeScript**
   1. Won. The rules use a small subset of Markdown.
   2. It builds React elements; no `innerHTML`.
   3. No new dependency (PR-01).
4. **Rules text: a Markdown library (marked, goldmark)**
   1. Lost. A dependency for headings, lists, bold and links.
5. **Card data: from the `.wasm`, built by Go**
   1. Won. One `CardInfo` builder for app and site.
6. **Card data: JSON imported by Vite, mapped in TypeScript**
   1. Lost. A second copy of the builder.
   2. `fakeBridge.ts` does this, but only for development.
7. **Links: payload in the URL fragment, compressed**
   1. Won. Short links; the browser has `CompressionStream`.
   2. Go has `compress/flate` for the app (11.3).
   3. The fragment never reaches a server.
8. **Links: plain JSON in the query string**
   1. Lost. Long URLs; the payload reaches server logs.
9. **Texts: one set of locale files for app and site**
   1. Won. Each key is translated once; one checker, one tool.
   2. The app gets the whole Russian interface in this step.
10. **Texts: separate locale files for the site**
    1. Lost. Engine messages would be translated twice.
    2. A Russian site with an English app splits the work in two.
11. **TinyGo for a smaller `.wasm`**
    1. Lost. It breaks `encoding/json` (ADR 001).
12. **Ship the `.wasm` gzipped; the page unpacks it**
    1. Won. Works on any host; `DecompressionStream` is built in.
13. **Rely on the host to compress the `.wasm`**
    1. Lost. Unknown on GitHub Pages; 7.9 MB per visit if not.

## Decision

Use options 1, 3, 5, 7, 9 and 12.

1. Layout
   1. `cmd/website/main.go` builds the `.wasm`.
   2. `frontend/site.html` and `frontend/src/site/` hold the page.
   3. `vite build --mode site` writes `frontend/dist-site/`.
   4. Site code never imports `frontend/src/bridge/`; lint checks it.
   5. `make site` builds the `.wasm` and the page together.
2. Wasm functions, on `globalThis.fourSouls`
   1. `version()`: the engine and app versions.
   2. `sets()`: the set codes the site carries.
   3. `cards()`: every `CardInfo` of those sets, with no images.
   4. `check(payload)`: the `Answer` as JSON.
   5. A refusal reason goes out as `Say`, a key with values.
   6. `notice()`: the fan-game notice from `internal/legal`.
   7. 12.4 adds `step(state)` and sets its shape.
   8. Errors come back as `{"error": "…"}`; nothing panics into the page.
3. Card data
   1. `cardInfo` moves from `cards.go` to `internal/cardinfo`.
   2. The app and the website both call it.
   3. The website passes no image folder, so cards show as text.
4. Card sets
   1. `internal/website` lists its set codes by hand; `cmd/website` exports them.
   2. A test in `cmd/cardgen` fails if a listed set has a stub card.
   3. A set joins the list when step 14 finishes it.
5. Rules text
   1. Vite imports `pkg/rules_engine/docs/rules/*.md` as text.
   2. The renderer knows headings, numbered lists, bold, italics, links and tables.
   3. A `**R-XXX-NN**` at the start of a list item is an anchor.
   4. Any other `R-XXX-NN` in the text links to that anchor.
   5. `open-questions.md` stays off the site; it is research.
6. Routes, in the URL fragment
   1. `#/<lang>/check/<data>`: a situation and its answer.
   2. `#/<lang>/rules/<rule ID>`: the rules text at that rule.
   3. `#/<lang>/faq/<entry>`: an FAQ entry (12.6).
   4. `#/<lang>/stack/<data>`: step-by-step stack (12.4).
   5. `<lang>` is `en` or `ru`; without it, the browser's language.
   6. Fragment routes need no server rewrites on GitHub Pages.
7. Link data
   1. `<data>`: base64url, no padding, of deflate-raw of the payload JSON.
   2. Above 2000 characters, the page warns that some chats cut links.
   3. Copying the payload itself always works.
8. Russian rules text
   1. Folder: `pkg/rules_engine/docs/rules/ru/`, same file names.
   2. Each Russian rule keeps its English ID.
   3. A missing rule shows the English text, marked (LO-06).
   4. `ru/source.json`: per rule ID, a hash of the English text it translates.
   5. `make i18n` lists missing rules and rules whose English changed.
   6. A Go test: every Russian ID exists in English, none twice.
9. Interface texts
   1. App and site share `frontend/src/locales/en.json` and `ru.json`.
   2. Site keys live in the area `site`.
   3. The unused-key check reads app and site code.
   4. `ru.json` holds the whole interface: app, site, engine messages.
   5. Glossary `glossary/ru.json` comes first; rules and interface use it.
   6. Notes, glossary and checks work as in ADR 021.
   7. Card names come from card data per language; texts stay English.
   8. ADR 021 drops "interface ships in English only" in the same change.
10. Size
    1. The build ships `site.wasm.gz` only.
    2. The page unpacks it with `DecompressionStream("gzip")`.
    3. Budget: 3 MB for `site.wasm.gz`.
    4. `make site` prints the size; CI fails above the budget.
    5. The page shows progress while the `.wasm` loads.
    6. Output: `.cache/site/site.wasm.gz`; Vite empties its own `dist-site/`.
    7. 2026-11-06, four sets: 8.9 MB raw, 2.2 MB gzip (12.1).
    8. The page unpacks only a stream still in gzip (12.2).
    9. A host may send it with `Content-Encoding: gzip`; the browser unpacks it then.
    10. The site ships Go's `wasm_exec.js` next to it; `NOTICE` names it.
11. Hosting
    1. GitHub Pages, deployed by GitHub Actions (`actions/deploy-pages`).
    2. A job in `release.yml` on a `v*` tag; tags exist only on `main`.
    3. URL: `https://foursoulsrepo.github.io/four_souls/`.
    4. Vite `base` matches that path.
12. Legal
    1. The site shows the fan-game notice, like the app.
    2. System fonts only; they cover Cyrillic.
    3. No file from the private submodule reaches the site.

## Consequences

### Positive

1. One engine answers in the app, in tests and on the site.
2. No new dependency; `NOTICE` only adds Go's `wasm_exec.js`.
3. Links work without a server and keep the payload private.
4. Rule IDs link the same way in answers, rules and FAQ.
5. Translators use one file and one tool for app and site.
6. Russian players get a Russian app, not only a Russian site.
7. The download size does not depend on the host.

### Negative and risks

1. A first visit downloads about 2 MB.
2. Every new set grows the `.wasm`; the budget may need raising.
3. An own Markdown renderer must grow with the rules files.
4. The Russian interface is a large translation: 959 keys.
5. Every new interface key now needs a Russian line, or falls back to English.
6. The page unpacks the `.wasm` itself; a little more page code.

### Neutral

1. The app's in-app rules viewer (11.1) can reuse the renderer later.
2. Deploying on every push to `main` is a one-line change later.
