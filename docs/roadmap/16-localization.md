# Step 16. Localization

Goal: menus, messages, and cards in the player's language.

Ideas: LO-01 – LO-08, B-07.

Rules of this step (owner, 2026-10-26):
1. The interface ships in English only; no interface translation yet.
2. Russian gets card names only; card texts come later.
3. English is always the fallback (LO-06).
4. Translating must be easy for anyone (LO-05).
5. Design: ADR 021.
6. Every sub-step after 16.2 is checked with the test language.
7. Per-key screenshots are paused; one screenshot per screen (owner, 2026-10-27).
8. Translation files stay JSON; YAML was weighed and lost (owner, 2026-10-27).

Owner notes (2026-10-27):

1. Owner reviewed step 16 and marks it solved.
2. A language with own images only is not listed in Settings; kept as is.
3. The images and tools repos are pushed.

---

### [x] 16.1 Translation ADR

1. Goal: pick how texts are translated.
2. Tasks:
   1. Compare i18n libraries for React; keep it small (PR-01).
   2. File format for translations.
   3. Fallback to English for missing keys.
   4. Update `NOTICE` with the chosen library.
3. Done when:
   1. ADR accepted by the **Owner**.
4. Result: ADR 021.
   1. Own `t()`, no library; `NOTICE` stays as it is.

### [x] 16.2 Language setting and test language

1. Goal: the player picks the language; developers check with a test one.
2. Tasks:
   1. `t()` and the locale loader from ADR 021, with Vitest tests.
   2. `en.json` with the first few keys: the settings screen.
   3. Language setting in the settings screen; default from the OS.
   4. The list shows English and Russian (owner).
   5. Russian has no interface file yet; it falls back to English.
   6. Test language `qps` from ADR 021; hidden from players.
   7. `make screenshots LOCALE=qps` runs the scenes in it.
   8. Use case: change the language.
3. Done when:
   1. The setting is saved and survives a restart.
   2. The settings screen shows in the test language.

### [x] 16.3 App texts

1. Goal: every UI string comes from `en.json`.
2. Tasks:
   1. Move menu, screen, settings, and table texts to keys.
   2. History lines too; it is built in the client (`history.ts`).
   3. Every key gets a note in `en.json` (ADR 021).
   4. Plural lines use plural forms (ADR 021).
   5. Lint: no hard-coded UI text left in `.tsx`.
   6. Check: no unused key; no used key missing from `en.json`.
3. Done when:
   1. The app looks the same in English.
   2. Screenshots in `qps` show no plain English text.

### [x] 16.4 Server messages

1. Goal: each client shows messages in its own language (LO-04).
2. Tasks:
   1. Refusal reasons, prompts, choices, event texts: key plus values.
   2. Card questions and options too; keys per card (ADR 021).
   3. Card names in messages go as card refs, not names.
   4. Protocol change in `protocol.go` and `protocol.ts`.
   5. Client renders the text; English as fallback.
   6. Dedicated server console stays in English.
   7. A Go test lists every message key; the unused-key check reads it.
3. Done when:
   1. Two clients, English and `qps`, play one game.
4. Result: `engine.Msg` and `engine.NewKey`; protocol 6, engine 0.4.0.
   1. `languages.spec.ts` plays with one page in each language.
   2. Ability words stay card data until card texts are translated.

### [x] 16.5 Russian card names

1. Goal: cards show Russian names (LO-02, LO-08).
2. Tasks:
   1. Importer in `FourSoulsRepo/tools`; this repo gets its output only.
   2. Source: the old site's migrations (`title_ru` by `text_id`).
   3. `text_id` is `<set>-<card id>`; about 20 cards do not match.
   4. The unmatched list and a proposed mapping go to the **Owner**.
   5. Write `text.ru.name` into `pkg/card_db/data/*.json`.
   6. Client shows names in the chosen language; English if missing.
   7. Translators shown with the card when known (CD-06).
3. Done when:
   1. Base Game cards show Russian names with Russian selected.
4. Owner answers (2026-10-27):
   1. Mapping: drop our set suffix, then a fixed list of the old site's typos.
   2. Suffix examples: `fsp2-pills_2_fsp2` → `fsp2-pills_2`, `r-pills_r` → `r-pills`.
   3. The fixed list, ours → old site:

| Ours | Old site |
|---|---|
| `b2-bum_bo` | `b2-bum-bo` |
| `b2-chest_2` | `b2-chest-2` |
| `b2-dark_chest_2` | `b2-dark_chest-2` |
| `b2-gold_chest_2` | `b2-gold_chest-2` |
| `b2-xl_floor` | `b2-xl_fLoor` |
| `fsp2-curse_of_impulse` | `fsp2-curse_of_Impulse` |
| `fsp2-whore_of_babylon` | `fsp2-whore_of_babylony` |
| `r-curse_of_the_soulless` | `r-curse_of_the_souless` |

   4. Every "Pills!" card is "Пилюля!"; the site's "Пилюли!" is a typo.
   5. Other names stay as on the site.
      1. "I Am Error": "Я - ошибка" (`r`) and "Я ошибка!" (`g2`) differ on purpose.
      2. "The Chest" and "Chest" are both "Сундук".
      3. "Flash Исаак" is correct.
   6. `translators.ru` for every imported name, as on foursouls.ru:
      1. Alchba, Mitboy, d0m1n1c132, Blacknenko, N-Chance.
5. Result: `scraper ru-names` in the tools repo; all 660 cards matched.
   1. `scraper build` keeps `text.ru` and `translators.ru` on a re-scrape.
   2. `App.Cards` sends `names` and `translators`; the store picks the language.
   3. The zoom credits the translators under the card.

### [x] 16.6 Card images per language

1. Goal: each language has its own card images (LO-07).
2. Tasks:
   1. Images move to `images/<lang>/<set>/`; English is `en`.
   2. Card data keeps paths without the language: `b2/the_d6.webp`.
   3. The loader adds `<lang>/`, then falls back to `en/`.
   4. Fonts, backs, dice, mats stay where they are.
   5. Commit in the private `FourSoulsRepo/card_db` repo; bump the submodule.
   6. Update ADR 004 and ADR 010 in place.
   7. Builds without the private images still work.
3. Done when:
   1. Every card shows its English image, as before.
   2. A test image under `ru/` shows with Russian selected.
4. Result: images in `images/en/<set>/`; `carddb.FindImage` adds the language.
   1. `App.Cards` sends `images` per language folder; `inLanguage` picks one.
   2. Tests: `TestFindImage`, `TestCardsImagesPerLanguage`, `cards.spec.ts` (fake Isaac).
   3. Tools: `scraper images` writes `images/en/<set>/`.

### [x] 16.7 Translation tooling

1. Goal: a translator checks a translation without reading code (LO-05).
2. Tasks:
   1. `make i18n LOCALE=xx`: missing keys, wrong placeholders, missing plural forms.
   2. It fails while anything is untranslated.
   3. `make i18n` without `LOCALE`: unused keys, stray keys; part of `make lint`.
   4. `todo.json` for agents and the glossary check (ADR 021, item 16).
   5. Screenshots: `make screenshots LOCALE=xx`, one per screen (ADR 021, item 14).
   6. The translator checks them by eye; no per-key pages for now.
   7. `todo.json` is split by area; `AREA=lobby` limits it to one (owner, 2026-10-27).
   8. Glossary terms go into each `todo.json` entry; no word check (owner, 2026-10-27).
   9. Checks live in `locale.ts` and run through Vitest; no new dependency.
   10. An agent writes `done.json`; `make i18n` merges it into `xx.json`, sorted.
   11. Required plural forms: only those whole numbers 0–1000 pick.
   12. Placeholders are checked per key; one plural form may leave out `{n}`.
   13. An English line whose meaning changes gets a new key (ADR 021, item 2).
   14. Trials: unit tests use test tables; a real file stays on a scratch branch.
3. Done when:
   1. `make i18n LOCALE=xx` lists every gap of a half-translated test file.
   2. `make screenshots LOCALE=xx` shows that file's lines on the screens.
   3. An agent fills a `qps`-like test file from `todo.json` alone.
4. Result: `make i18n`; checks in `locale.ts`, runner `translations.test.ts`.
   1. Glossary: `locales/glossary/en.json`, term → note; other files term → word.
   2. `done.json`: key → text; `null` deletes a stray key.
   3. `make lint` runs `make i18n`; the help lists it and `AREA`.
   4. e2e: `cardName(ref)`; `make screenshots LOCALE=ru` passes.
   5. e2e reads `locales/<LOCALE>.json` when there is one.
   6. Trial on a scratch `pl.json`: half file gaps, screenshots, agent fill.
   7. Use case UC-33; 19.79 found by the trial.

### [x] 16.8 Setup scripts

1. Goal: one script installs every tool `make` needs (B-07).
2. Tasks:
   1. `scripts/setup.sh` for macOS and Linux.
   2. Profiles: `translator` (Node, Playwright, make) and `developer` (all).
   3. `make doctor`: lists each tool, its version, and what is missing.
   4. Versions as in `docs/wiki/Linters.md`; keep one source of truth.
   5. No Windows for now (owner, 2026-10-27).
   6. Wiki: `Developer-setup.md`, per OS, and what each target needs.
3. Done when:
   1. A clean macOS and a clean Linux run the scripts, then `make doctor` passes.
4. Result: `scripts/setup.sh`; `make doctor` runs its `--check`.
   1. Versions: `scripts/versions.env`; Go from `go.mod`, Playwright from `e2e/package.json`.
   2. CI reads `versions.env` too; Linters.md points to it.
   3. Linux: Go and Node from their archives into `~/.local`, no sudo.
   4. golangci-lint goes into `.cache/bin`; `make lint` uses it when there.
   5. `make lint` runs shellcheck when it is installed.
   6. Fedora: Playwright's WebKit does not run there; skipped.
   7. Tested in clean Ubuntu 24.04 and Fedora 44 containers, both profiles.
   8. macOS: only `--check` on the owner's Mac; nothing installed there.
   9. The Ubuntu screenshot run found 19.82 (dice test without images).

### [x] 16.9 Translator guides

1. Goal: anyone can add a language from the wiki (LO-05).
2. Tasks:
   1. `Translate-the-interface.md`: setup, start from `{}`, translate, check, screenshots.
   2. `Translate-cards.md`: `text.<lang>`, translators (CD-06), images under `images/<lang>/`.
   3. How to send the result: a pull request.
   4. Both pages link `Developer-setup.md` for the tools.
   5. Wiki group for them: a new group "Translators" (owner).
   6. A section on translating with an AI agent: `todo.json`, checks, screenshots.
   7. The new group goes into `Home.md`, `_Sidebar.md`, and `CLAUDE.md` (3.14).
3. Done when:
   1. A walkthrough of each guide on a clean checkout works.
4. Result: three wiki pages in the new group Translators.
   1. `Translate-the-interface`, `Translate-cards`, `Translating-with-an-AI-agent`.
   2. A new language starts as `{}`; a copy of `en.json` fails the check.
   3. Card-only languages: Settings lists a language with an interface file.
   4. Walkthrough: clean Ubuntu 24.04 container, public clone, `de` trial.
   5. Screenshots on two cores: `E2E_WORKERS=2`; the guide says so.
   6. Card images: sent to the maintainers through the pull request.
   7. The group is in `Home.md`, `_Sidebar.md`, `CLAUDE.md` 3.14.
5. Owner answers after review (2026-10-27):
   1. A card-only language is listed in Settings: WIP translations are welcome.
   2. Card images: the guide explains how to make an own images repo.
   3. A half-done interface file ships; English fills the rest.
   4. A new language starts as `{}`, not a copy of `en.json`.
   5. Files stay JSON, not YAML.
      1. Placeholders like `{n} cards` would need quotes in YAML.
      2. Every reader (checks, merge, e2e, Go) would need a YAML library.

### [x] 16.10 Card-only languages in Settings

1. Goal: a language with card names only is a choice in Settings (owner).
2. Tasks:
   1. The list comes from card data, not the fixed `cardLanguages = ['ru']`.
   2. A language is listed if it has an interface file or card names.
   3. Go sends the card languages; the client merges them with locale files.
   4. Update ADR 021 item 11.4 in place if the wording drifts.
   5. Translator guides: drop the empty `{}` workaround for card-only languages.
3. Done when:
   1. A card-only test language shows in Settings without a locale file.
4. Result: card languages from the `names` of `App.Cards`; no new binding.
   1. `cardLanguagesOf` in `locale.ts`; `cardLanguages` is gone.
   2. `loadCards` calls `setCardLanguages`; the store keeps `languages`.
   3. A saved card-only language shows once the cards load.
   4. Fake bridge `?cardlang=xx` names Isaac in `xx`; card files untouched.
   5. Tests: `TestCardsNamesCardOnlyLanguage`, Vitest, `settings.spec.ts`.
   6. ADR 021 item 11.6, UC-32 2c, `Client.md`, `Translate-cards.md`.

### [x] 16.11 Clear JSON errors for translators

1. Goal: a syntax mistake in a language file says where and what to fix.
2. Tasks:
   1. `make i18n` names the file, line, and column of a JSON error.
   2. Plain hints for common mistakes: a trailing comma, a missing quote.
   3. Card data JSON errors in the loader give the line too.
   4. Guides: replace the `node` one-liner with the new message.
3. Done when:
   1. A trailing comma in `xx.json` and in a set file gives the line and a hint.
4. Result: `jsonError.ts` and `pkg/card_db/jsonerr.go`, the same hints.
   1. Message: `file:line:col: JSON error: hint`, then the line and a caret.
   2. The place comes from an own scan; `JSON.parse` messages differ per Node.
   3. Go: `SyntaxError.Offset` and `UnmarshalTypeError` (a value of the wrong kind).
   4. `make i18n` reads every file first; a broken one skips the checks and the merge.
   5. The loader names every broken set file at once, then stops.
   6. The app: Vite's own JSON error is clear already; documented in `Client.md`.
   7. `overrides/*.json`: no code in this repo reads them; nothing to do.
   8. Tests: `jsonError.test.ts`, `TestJSONError*`; ADR 004 item 5.9, UC-33 5a.

### [x] 16.12 Own card images repo for translators

1. Goal: a translator adds card images without access to the private repo (owner).
2. Tasks:
   1. Research how a build or `wails dev` can read images from another folder.
   2. E.g. `$FOUR_SOULS_CARDS` or `cards/` next to the binary (ADR 004, item 6).
   3. If the app can't do it yet, add the simplest way, with tests.
   4. `Translate-cards.md`: make a repo with `images/<lang>/<set>/`, point the app at it.
   5. Say how to share it with the maintainers.
3. Done when:
   1. A walkthrough on a clean checkout shows a translator's own image.
4. Result: `$FOUR_SOULS_IMAGES` names an own images folder, any build mode.
   1. Research: `$FOUR_SOULS_CARDS` swaps data and images together.
   2. It is skipped in embedded builds; so a new variable was needed.
   3. Layout as the private repo: `<lang>/<set>/<file>`.
   4. Order: own image, built-in image in that language, then English.
   5. `carddb.WithImages` layers it; only language folders are read.
   6. `os.Root` keeps reads inside the folder, symbolic links too.
   7. Tests: `TestWithImages*`, `TestCardsFSOwnImages`, `TestCardsOwnImages`.
   8. Images alone do not list a language in Settings; names do (16.10).
   9. Walkthrough: clean Ubuntu 24.04, public clone, `wails dev` under Xvfb.
   10. The own `de` image showed in the gallery; ADR 004 6.7, ADR 021 13.10, UC-34.
