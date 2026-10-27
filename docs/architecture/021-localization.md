# 21. Localization

* **Status:** Accepted
* **Date:** 2026-10-26
* **Authors:** @HardDie

---

## Context

1. Every text is English and written in code (LO-01).
   1. About 450 strings in the frontend.
   2. About 130 in the engine: refusals, prompts, event texts.
   3. About 150 card questions with their options in card scripts.
2. Each client must show text in its own language (LO-04).
   1. Two players at one table may use two languages.
3. Translating must be easy for anyone, not only developers (LO-05).
   1. One file per language.
   2. Each line says where it shows.
   3. A translator can see a line in place on a screenshot.
4. English is always the fallback (LO-06).
5. Card data already has `text.<lang>` and `translators.<lang>` (ADR 004).
6. Card images will differ per language (LO-07).
7. Some languages have several plural forms.
   1. English: "1 card", "2 cards".
   2. Russian: "1 карта", "2 карты", "5 карт".
8. History is already built in the client from events.
9. Keep it simple (PR-01); few dependencies.

## Considered options

1. **i18next with react-i18next**
   1. Lost. About 40 kB and many features we do not need.
   2. The translator still edits plain JSON; the library adds nothing there.
2. **FormatJS (react-intl), ICU messages**
   1. Lost. ICU syntax is hard for translators to read.
3. **Lingui**
   1. Lost. Needs a compile step and macros.
4. **gettext `.po` files**
   1. Lost. Needs tools in both Go and TS.
   2. Editors like Poedit help, but plain JSON edits anywhere, GitHub too.
5. **Own small `t()`**
   1. Won. About 60 lines; plurals come from the browser.
6. **English text as the key**
   1. Lost. No room for a note; one English word may need two translations.
7. **Named keys with a note**
   1. Won. The note tells the translator where the line shows.
8. **Notes in a separate file**
   1. Lost. Two files drift apart.
9. **One file per area per language**
   1. Lost. The owner wants one file per language.
10. **Server sends ready text in each client's language**
    1. Lost. The server would need every language file.
    2. Records and spectators would freeze one language.
11. **Server sends a key and values**
    1. Won. The client renders; the server stays language-free.
12. **Plurals by `Intl.PluralRules`**
    1. Won. Built into every browser and webview; knows each language.
13. **Avoid numbers in sentences instead of plurals**
    1. Lost. Too many lines need a count.
14. **Images per language in one folder: `images/<lang>/<set>/`**
    1. Won. One language stays in one place (owner).
15. **Images with a suffix: `b2/the_d6.ru.webp`**
    1. Lost. Languages mix in one folder.

## Decision

Use options 5, 7, 11, 12 and 14.

1. Files
   1. One file per language: `frontend/src/locales/<lang>.json`.
   2. `<lang>` is an ISO 639-1 code, as card_db `Lang`: `en`, `ru`.
   3. The frontend loads every file in that folder; no code lists them.
   4. Adding a language is adding a file.
2. `en.json`
   1. It is the source; every key must be in it.
   2. Each key holds `text` and `note`.
   3. `note`: where the line shows and what the values are.
   4. A line whose meaning changes gets a new key.
   5. Its old translations then show as stray keys, the new key as missing.
3. Other languages
   1. Each key holds the text only.
   2. A missing key falls back to English.
4. Example

```json
{
  "lobby.join": {
    "text": "Join",
    "note": "Button on a game row in the lobby list."
  },
  "refuse.attack_cost": {
    "text": "Attacking costs {cents}¢",
    "note": "Hint on the Attack button when the player can't pay."
  },
  "hand.count": {
    "text": { "one": "{n} card", "other": "{n} cards" },
    "note": "Under a player's hand. n: cards in it."
  }
}
```

```json
{
  "lobby.join": "Войти",
  "hand.count": { "one": "{n} карта", "few": "{n} карты", "many": "{n} карт" }
}
```

5. Keys
   1. Lowercase words joined by `.` and `_`.
   2. The first word is the area: `menu`, `lobby`, `settings`, `table`, `history`.
   3. Server areas: `refuse`, `prompt`, `event`, `ask`.
   4. Card questions: `card.<card id>.<name>`, e.g. `card.cheese_grater.put_back`.
   5. Questions shared by many cards stay in `ask`, e.g. `ask.which_deck`.
   6. The file is sorted by key, so one screen's lines sit together.
   7. Server and app errors use the area `error`.
   8. A key may hold digits: card ids do, `card.chest_2.…`.
6. Values
   1. Placeholders: `{name}`; the same names in every language.
   2. Values are numbers, plain strings, or refs.
   3. A card ref: `{"card": "<ref>"}`; shown as the name in the viewer's language.
   4. A player: `{"player": N}`; shown as the name in the seat color.
   5. Text is never HTML; React escapes it.
7. Plurals
   1. A plural line holds one text per form, by CLDR name.
   2. Forms: `zero`, `one`, `two`, `few`, `many`, `other`.
   3. `Intl.PluralRules(lang)` picks the form from the value `n`.
   4. Each language fills the forms its rules list.
   5. A missing form uses `other`, then English.
8. `t(key, values)`
   1. Lookup order: the chosen language, then English, then the key itself.
   2. A key missing from English logs a warning in development.
   3. Frontend keys are typed from `en.json`; a typo fails `tsc`.
   4. Frontend code passes a key as a string literal, so checks can find it.
9. Server messages
   1. The engine sends `{key, args}` wherever it sent text.
   2. Refusals, prompt questions, choice options, event texts.
   3. Engine keys are declared in Go through one registry function.
   4. A Go test writes the registry to `frontend/src/locales/server-keys.json`.
   5. Its `-update` flag rewrites the file; CI fails when it is stale.
   6. The dedicated server console stays in English.
   7. The protocol version goes up with this change.
   8. The registry function is `engine.NewKey`; `engine.Msg` is the message.
   9. `Key` is no string type: a bare string never compiles as a key.
   10. The writing test is `TestServerKeys` in the app module (`keys_test.go`).
   11. It links every package that declares keys: engine, cards, server, app.
   12. A Go test also fails on a declared key with no line in `en.json`.
   13. A refusal with no words for players reads `refuse.not_usable`.
   14. A server error carries its message as `text`.
   15. A Go error the app shows is a message as JSON; `messageOf` reads it.
   16. An ability's own words stay card data: a `{text}` value until card texts.
   17. Event `text` stays a code for the client and card triggers, never shown.
   18. Tests name an option by its plain form: a ref, a number, or key and values.
10. Checks
    1. Every key used in code is in `en.json`.
    2. Every key in `en.json` is used: frontend literals or `server-keys.json`.
    3. No other language has a key that `en.json` lacks.
    4. Placeholders match English in every language.
    5. `make i18n LOCALE=xx` also fails on untranslated keys and missing plural forms.
    6. Required plural forms: those `Intl.PluralRules` picks for 0–1000.
    7. So Russian needs `one`, `few`, `many`; `other` is for fractions only.
    8. Placeholders are compared per key, not per plural form.
    9. No form may use a name English lacks; all forms together use every one.
    10. One plural form may leave out `{n}`: "одна карта".
    11. The checks live in `locale.ts`; `make i18n` runs them through Vitest.
11. Language setting
    1. Stored in the app settings file.
    2. Default: the OS language from `navigator.language`, if supported.
    3. Else English.
    4. A language is listed if it has an interface file or card names.
    5. Each is listed by its own name, via `Intl.DisplayNames`.
    6. Card languages come from the card names Go sends (16.10).
12. Test language `qps`
    1. Built at run time from `en.json`; no file.
    2. `Join` becomes `[!! Ĵôîñ ~~ !!]`: accented, about 30% longer.
    3. Placeholders stay as they are.
    4. Plural rules as in English.
    5. Shown only in development builds.
    6. Plain English on screen in `qps` means a missed string.
13. Cards
    1. Names and texts stay in card_db: `text.<lang>`.
    2. A missing language falls back to English.
    3. Images: `images/<lang>/<set>/<file>`; `en` is English.
    4. Card data paths stay language-free: `b2/the_d6.webp`.
    5. The loader adds `<lang>/`, then falls back to `en/`.
    6. Fonts, backs, dice and mats stay outside language folders.
    7. Go sends each card's names and translators in every language.
    8. The client shows the chosen language's name; English if missing.
    9. The zoom credits that language's translators (CD-06).
14. Screenshots for translators
    1. `make screenshots LOCALE=xx`: one PNG per screen, in that language.
    2. The translator checks the lines on them by eye.
    3. Per-key screenshots are paused (owner, 2026-10-27); idea LO-05.
15. Writing lines
    1. Names stay in their dictionary form in every language.
    2. E.g. "{player} → {monster}: attack", not "{player} attacked {monster}".
    3. So languages with cases need no forms of names.
16. Tools serve people and agents alike (owner, 2026-10-26)
    1. `make i18n LOCALE=xx` prints plain lines for people.
    2. It also writes `.cache/translations/xx/todo.json`.
    3. Each entry: key, English, note, placeholders, plural forms needed.
    4. An agent fills the language file from it, then runs the check again.
    5. The screenshots of English give an agent the same context.
    6. Glossary: `frontend/src/locales/glossary/<lang>.json`.
    7. It fixes game terms per language: loot, soul, treasure, cents.
    8. Each entry lists the glossary terms its English line uses.
    9. No check of the words; case endings make it unreliable (owner, 2026-10-27).
    10. `todo.json` is grouped by area: the key's first word.
    11. `make i18n LOCALE=xx AREA=lobby` writes one area only.
    12. An agent translates one area per pass, then runs the check again.
    13. The agent writes only `.cache/translations/xx/done.json`.
    14. `make i18n LOCALE=xx` merges it into `xx.json`, sorted by key.
    15. A file in `locales/` is a choice in Settings, so no test file sits there.
    16. Unit tests check test tables; a trial language stays on a scratch branch.
    17. `glossary/en.json`: each English term with a note.
    18. `glossary/<lang>.json`: each term's word in that language.
    19. The glossary folder and `server-keys.json` are no languages.
    20. The runner is `translations.test.ts`; `make lint` runs `make i18n`.

## Consequences

### Positive

1. A translator edits one JSON file and needs no code.
2. Notes and screenshots give each line its context.
3. Each client shows its own language at one table.
4. No new dependency; `NOTICE` stays as it is.
5. Records stay language-free; replays show the viewer's language.
6. An agent can translate a new language with the same checks.

### Negative and risks

1. Every UI string and engine message moves to a key: a large change.
2. Card scripts lose their inline English questions; text lives in `en.json`.
3. Sentences in names' dictionary form may read stiffly.
4. Moving images changes the private repo and ADRs 004 and 010.
5. A translator finds a line on a screen by eye; no key points to it.

### Neutral

1. LO-03 (a card in other languages in game) is not part of this.
2. The rules text and the rules website are not translated here.
3. Interface ships in English only; Russian gets card names first.
