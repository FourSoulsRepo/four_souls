# Step 12. Rules website

Goal: check rules in the browser, in English and Russian; no server.

Ideas: A-03, LR-08, LR-09, LR-10.

Order: 12.7, 12.1, 12.2, 12.3, 12.5, 12.4, 12.6, 12.10, 12.8, 12.11, 12.9.

Moved before step 18 (owner, 2026-11-06).
12.5 and 12.6 came from step 11 (11.8, 11.7).

Rules of this step (owner, 2026-11-06):
1. The site gets a Russian interface; the app gets it too (12.11).
   1. Code quality beats step 16 rule 1 (English-only interface).
   2. App and site share one set of locale files (ADR 026).
2. Card texts stay English on the Russian site, for now.
   1. Card names are Russian (16.5).
3. An agent drafts the Russian rules; the owner checks them.
   1. Compare the draft with the Russian rules on foursouls.ru.
4. The site carries every finished set.
5. The site deploys on a release tag.
   1. Every push to `main` may come later.

---

### [x] 12.1 Engine in WebAssembly

1. Goal: `cmd/website` exposes the engine to JS.
2. Tasks:
   1. Function: payload in, answer with reasons out.
   2. Function: engine version.
   3. Measure and report the `.wasm` size.
   4. Wrap `engine.RunSituation`; every finished card set.
   5. Reasons and events go out as message keys and values (16.4).
3. Done when:
   1. A JS test evaluates a sample payload.
4. Result: `cmd/website` exports `fourSouls.*`; logic in `internal/website`.
   1. `version`, `sets`, `cards`, `check`, `notice`; errors as `{"error"}`.
   2. `engine.Answer` gains `say`: the refusal's key with values.
   3. A situation names a set by code (`b2`) or name.
   4. `internal/cardinfo` builds `CardInfo` for app and site.
   5. Sets b2, g2, fsp2, r; a `cmd/cardgen` test rejects stubs.
   6. `make site`: `.cache/site/site.wasm.gz`, budget 3,000,000 bytes.
   7. Size: 8,874,685 bytes raw, 2,156,660 bytes gzip.
   8. Node smoke test `cmd/website/smoke.mjs` runs in `make site` and CI.
   9. Wiki: `Website-wasm`.

### [x] 12.2 Site UI

1. Goal: a simple page (A-03).
2. Tasks:
   1. Paste payload or open a link.
   2. Show the answer, reasons, rule links.
   3. Rules text with anchors by rule ID.
   4. FAQ: added in 12.6.
   5. Copy a link to the shown situation; format from 12.7.
   6. Cards show as text from card_db; the site never has images.
   7. Use case: check a situation on the website.
3. Done when:
   1. A link copied from the site opens the same situation.
4. Result: `frontend/site.html` and `frontend/src/site/`.
   1. `make site` builds the page into `frontend/dist-site/`, base `/four_souls/`.
   2. It ships `site.wasm.gz` and Go's `wasm_exec.js`; `NOTICE` names it.
   3. Unpacks only a gzip stream: `vite preview` sends `Content-Encoding`.
   4. Routes `#/<lang>/check/<data>`, `#/<lang>/rules/<ID>`.
   5. Check view: payload, examples, answer, rule link, expect, events.
   6. Copy link, Copy payload; a warning above 2000 characters.
   7. Rules view: own Markdown reader, one anchor per rule ID.
   8. Texts in area `site` of `en.json`; lint bans bridge and store imports.
   9. Vitest: links, routes, renderer; e2e `site.spec.ts` on `vite preview`.
   10. `make e2e` and `make screenshots` run `make site` first.
   11. Use case UC-38; wiki `Website-wasm` covers the page.

### [ ] 12.3 GitHub Pages

1. Goal: the site is online.
2. Tasks:
   1. Workflow builds and deploys on release.
   2. App config gets the site base URL (11.3).
   3. Releases come only from `main` (CLAUDE.md 5.23).
3. Done when:
   1. The public URL loads and answers payloads.

### [ ] 12.4 Step-by-step stack

1. Goal: users learn how the stack resolves (LR-10).
2. Tasks:
   1. Wasm function: resolve one step, return the new state.
   2. Each step shows what happened and why, with rule links.
   3. User may add a response between steps.
   4. Step back to the previous state.
   5. Needs stacked events in payloads (12.5).
3. Done when:
   1. A three-event stack resolves step by step and back.
   2. Use case: solve a stack step by step.

### [ ] 12.5 Any event on the stack

1. Goal: payloads start with a ready stack (LR-09).
2. Tasks:
   1. `setup.stack` takes loot cards, abilities, attack rolls, triggers.
   2. Each event: owner, targets, choices made.
   3. Bottom first, ready to resolve.
   4. The site (12.2) takes these payloads; the sandbox (11.2) later too.
   5. Update the wiki page `Situation-payload`.
   6. Moved from 11.8.
3. Done when:
   1. A payload with three stacked events resolves as expected.
   2. Old dice-only payloads still load.

### [ ] 12.6 FAQ from frequent questions

1. Goal: answer what players ask most (LR-08).
2. Tasks:
   1. Group similar questions from the step 3 sources by topic.
   2. Rank them by how often they are asked.
   3. Write each answer in our own words: how it works and why.
   4. Link every answer to its rule IDs (R-…).
   5. Add a situation payload where the case can be tried.
   6. Show the FAQ on the site, searchable; the rules viewer (11.1) reuses it.
   7. Include the owner's picks from `open-questions.md` (e.g. Q-09).
   8. Read the BoardGameGeek rules forum too (R-02):
      https://boardgamegeek.com/boardgame/255664/the-binding-of-isaac-four-souls/forums/66
   9. Its answered questions also become engine tests (R-03).
   10. Moved from 11.7.
3. Done when:
   1. The top questions have entries; every link opens its rule.
   2. Payload examples pass as tests.

### [x] 12.7 Website ADR

1. Goal: decide how the site is built, before any code.
2. Tasks:
   1. Site code: a second Vite entry in `frontend/`, or its own small app.
   2. Reuse `t()` and the locale loader (ADR 021) either way.
   3. Rules Markdown: rendered at build time or in the browser.
   4. Link format: payload in the URL fragment; compression; size limit.
   5. Card sets in the `.wasm`; a size budget.
   6. Where Russian rules live; how their IDs stay in sync.
   7. Hosting: GitHub Pages; when it deploys.
   8. Any new library goes into `NOTICE`.
   9. Written in the main conversation (CLAUDE.md 5.21).
   10. Russian engine messages: shared with the app, or site-only.
3. Done when:
   1. ADR accepted by the **Owner**.
4. Result: ADR 026, accepted 2026-11-06.
   1. One set of locale files for app and site; the app gets Russian (12.11).
   2. The `.wasm` ships gzipped; the page unpacks it.

### [ ] 12.8 Russian rules text

1. Goal: our rules text in Russian, with the same rule IDs.
2. Tasks:
   1. Translate every file in `pkg/rules_engine/docs/rules/`.
   2. Folder and file layout from 12.7.
   3. Card names as in the Russian card data (16.5).
   4. A Go test: every Russian rule ID exists in English, none twice.
   5. A missing Russian rule falls back to English (LO-06).
   6. `make i18n` reports missing and outdated Russian rules.
   7. Wiki: how to translate the rules, in the translators group.
   8. Compare terms with the foursouls.ru rules; write our own words.
   9. Any download of foursouls.ru lives in `FourSoulsRepo/tools`.
3. Done when:
   1. Every English rule ID has a Russian text.
   2. **Owner** checks the Russian text.

### [ ] 12.9 Russian on the website

1. Goal: the site works in English and Russian.
2. Tasks:
   1. Language switch; default from the browser; English fallback.
   2. Site texts through `t()`, in the shared locale files (area `site`).
   3. Russian site keys; engine messages are done in 12.11.
   4. FAQ answers (12.6) in Russian.
   5. Rule links open the rule in the chosen language.
   6. The link keeps the language.
   7. `make screenshots` covers the site in both languages.
   8. Use case: read the rules in Russian.
3. Done when:
   1. One payload answers in both languages, with working rule links.

### [ ] 12.10 Russian glossary

1. Goal: one Russian word per game term, before any translation.
2. Tasks:
   1. `frontend/src/locales/glossary/ru.json` for every term in `glossary/en.json`.
   2. Terms as in the Russian card data (16.5) and foursouls.ru (S-RU).
   3. Where they disagree, list both; the **Owner** picks.
3. Done when:
   1. Every English glossary term has a Russian word.

### [ ] 12.11 Russian app interface

1. Goal: the app in Russian, from the same `ru.json` as the site.
2. Tasks:
   1. Fill `frontend/src/locales/ru.json` with `make i18n LOCALE=ru`, one area per pass.
   2. Every key: app texts, history, engine messages (ADR 021).
   3. Words from the glossary (12.10).
   4. `make screenshots LOCALE=ru`: every screen, no cut text.
   5. ADR 021: drop "interface ships in English only".
   6. Update the use case "change the language".
3. Done when:
   1. `make i18n LOCALE=ru` reports no missing key.
   2. A full game plays in Russian.
   3. **Owner** checks the screenshots.
