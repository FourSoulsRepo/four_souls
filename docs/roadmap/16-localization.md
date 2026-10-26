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

### [ ] 16.4 Server messages

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

### [ ] 16.5 Russian card names

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

### [ ] 16.6 Card images per language

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

### [ ] 16.7 Translation tooling

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
3. Done when:
   1. `make i18n LOCALE=xx` lists every gap of a half-translated test file.
   2. `make screenshots LOCALE=xx` shows that file's lines on the screens.
   3. An agent fills a `qps`-like test file from `todo.json` alone.

### [ ] 16.8 Setup scripts

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

### [ ] 16.9 Translator guides

1. Goal: anyone can add a language from the wiki (LO-05).
2. Tasks:
   1. `Translate-the-interface.md`: setup, copy `en.json`, translate, check, screenshots.
   2. `Translate-cards.md`: `text.<lang>`, translators (CD-06), images under `images/<lang>/`.
   3. How to send the result: a pull request.
   4. Both pages link `Developer-setup.md` for the tools.
   5. Wiki group for them: a new group "Translators" (owner).
   6. A section on translating with an AI agent: `todo.json`, checks, screenshots.
   7. The new group goes into `Home.md`, `_Sidebar.md`, and `CLAUDE.md` (3.14).
3. Done when:
   1. A walkthrough of each guide on a clean checkout works.
