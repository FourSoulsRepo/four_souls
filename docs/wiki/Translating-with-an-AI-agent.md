# Translating with an AI agent

An AI coding agent (one that can read files and run commands in your copy of the code) can write a first draft of a whole language. It uses the same checks and screenshots as a person ([Translate the interface](Translate-the-interface)), so read that page first. You stay the translator: you fix the glossary, read the result, and send it.

## How it works

1. `make i18n LOCALE=xx` writes `.cache/translations/xx/todo.json`: every line still to do, grouped by area (the key's first word). Each entry has the key, the English, the note, the placeholders, the plural forms needed, and the glossary terms its English uses with your word. Its `how` list tells the agent the rules; `plural_forms` gives the numbers each form is for.
2. The agent writes its translations into `.cache/translations/xx/done.json` only, key → text: `{"lobby.join": "Beitreten"}`, a plural line as an object of forms. It never edits `frontend/src/locales/xx.json` itself.
3. The next `make i18n LOCALE=xx` merges `done.json` into `xx.json`, sorted by key, deletes `done.json`, and lists what is left. A key `en.json` lacks is left out and named; a key set to `null` is deleted (that is how stray keys go).
4. `AREA=lobby` limits the list and `todo.json` to one area. One area per pass keeps each step small enough for the agent to do well; it then runs the check again before the next area.
5. `make screenshots LOCALE=xx` shows the result on every screen, into `.cache/screenshots/xx/`; the English screens (`make screenshots`, into `.cache/screenshots/`) show the agent where a line sits.

## Before you start

1. Set up the tools and fork, as in [Translate the interface](Translate-the-interface#1-get-the-code-and-the-tools).
2. Create `frontend/src/locales/xx.json` with `{}`.
3. Decide the game terms first: write `frontend/src/locales/glossary/xx.json` (term → word) for the terms in `glossary/en.json`. You may ask the agent to propose words, but save the file yourself. Every entry in `todo.json` then carries your words, so the whole translation uses them.
4. Optional: `make screenshots` once, so the English screens are there for the agent.

## The prompt

Paste this into the agent, with your language's name and code in place of `German` and `de`:

```text
Translate the interface of the Four Souls game into German (language code de).
You are in the root of its repository.

Rules:
- Write only .cache/translations/de/done.json. Do not edit any other file.
- done.json is one JSON object: key -> text. A plural line is an object of
  exactly the forms the entry's "forms" lists, e.g. {"one": "...", "other": "..."}.
- Keep every {placeholder} exactly as in the English. Translate no key.
- Names of players and cards stay in their dictionary form: build each sentence
  so that {player}, {card}, {monster} need no case ending.
- Use the glossary words each entry lists; an empty word is not decided yet:
  pick one and use it everywhere.
- The note says where a line shows; keep button and label texts short.

Loop, one area at a time:
1. Run: make i18n LOCALE=de
   It fails while lines are left; that is expected. Read the area names
   under "areas" in .cache/translations/de/todo.json.
2. For the next area: run make i18n LOCALE=de AREA=<area>, read
   .cache/translations/de/todo.json, translate every entry into done.json.
3. Run make i18n LOCALE=de AREA=<area> again. It merges done.json into
   frontend/src/locales/de.json. Fix every line it still lists the same way.
4. Go on with the next area until make i18n LOCALE=de passes.

Then run make screenshots LOCALE=de and look at the PNGs in
.cache/screenshots/de/ next to the English ones in .cache/screenshots/.
Shorten lines that are cut off or overflow their button, through done.json
and make i18n LOCALE=de again. Finally report what you were unsure about.
```

## After the agent

1. Read the file yourself: `frontend/src/locales/xx.json`. An agent can be wrong in tone, game terms, or which word fits a button.
2. Run `make i18n LOCALE=xx` and `make screenshots LOCALE=xx` yourself and look at the screens.
3. Send it as a pull request ([Translate the interface](Translate-the-interface#7-send-it)) and say that an agent wrote the draft and you checked it.

The same loop works later for new lines: `make i18n LOCALE=xx` lists only what is missing, so the agent translates just those.
