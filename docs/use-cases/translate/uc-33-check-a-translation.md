# UC-33: Check a translation

**Module:** `frontend`, `Makefile`
**Status:** Implemented
**Actors:** Translator, or an AI agent
**Goal:** The translator sees every line still to do, without reading code
**Preconditions:** A checkout with Node; `make frontend` works

## Main scenario (happy path)

1. The translator runs `make i18n LOCALE=xx`.
2. It prints one plain line per gap.
   1. `missing`: a key not translated yet, with its English.
   2. `placeholders`: `{names}` that differ from English.
   3. `forms`: a plural line lacking a form the language needs.
   4. `shape`: one text where plural forms belong, or the other way.
   5. `stray`: a key `en.json` no longer has.
3. It writes `.cache/translations/xx/todo.json`.
   1. Grouped by area: the key's first word.
   2. Each entry: key, English, note, placeholders, forms needed.
   3. The glossary terms the line uses, with this language's word.
4. The translator edits `frontend/src/locales/xx.json`.
5. They run `make i18n LOCALE=xx` again.
6. When nothing is left, it passes.
7. `make screenshots LOCALE=xx` shows each screen in the language.
   1. Into `.cache/screenshots/xx/`; the translator checks by eye.

## Alternative scenarios and errors

* **1a. `AREA=lobby`:** the list and `todo.json` hold that area only.
* **4a. An agent translates:** it writes `.cache/translations/xx/done.json`.
  * The next `make i18n LOCALE=xx` merges it into `xx.json`, sorted.
  * Keys `en.json` lacks stay out and are listed.
  * A key set to `null` there is deleted: how a stray key goes.
* **1b. `LOCALE=en` or `LOCALE=qps`:** it fails and says why.
* **1c. No `LOCALE`:** every language file is checked.
  * Stray keys, shapes and placeholders; missing lines are fine.
  * Also unused keys in `en.json`; `make lint` runs it.

## Postconditions

* `xx.json` holds every line, sorted by key.
* `xx` is a choice in Settings (UC-32).
