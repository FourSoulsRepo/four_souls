# UC-32: Change the language

**Module:** `frontend`, `internal/settings`
**Status:** Implemented
**Actors:** Player
**Goal:** The app shows its texts in the player's language
**Preconditions:** The app is on the main menu

## Main scenario (happy path)

1. The player opens **Settings** from the main menu.
2. The Language list shows the language in use.
   1. Each language shows in its own name: English, Русский.
3. The player picks a language and presses **Save**.
4. Go writes it into the settings file (`language`).
5. The screens switch to that language at once.
   1. No restart is needed.
6. After a restart, the app starts in the saved language.
7. Cards show their names in that language (16.5).
   1. Text cards, zoom, history, stack, prompts and choices.
   2. The zoom credits the card's translators under it.
8. Cards show their images in that language (16.6).
   1. Only where the private images have that language's file.

## Alternative scenarios and errors

* **2a. First start:** no language is saved yet.
  * The app takes the OS language, if it is listed.
  * Else English.
* **3a. Russian (16.2):** it has no interface file yet.
  * The texts stay English; card names are Russian (16.5).
* **5a. A line missing in the language:** the English line shows.
* **6a. A saved language this build lacks:** as on the first start.
* **2b. Test language `qps`:** listed in development builds only.
  * Accented, longer English in `[!! … !!]` (ADR 021).
  * Plain English on screen then means a missed string.
* **7a. A card with no name in the language:** its English name shows.
  * No translators are credited then.
* **8a. A card with no image in the language:** its English image shows.
  * Without the private images: the text card, as before.
* **2c. A language with card names only (16.10):** it is listed too.
  * No interface file is needed; the texts stay English.
  * Cards show their names in it; others stay English.
  * It joins the list once the card data has loaded.

## Postconditions

* The settings file holds the chosen language.
* Every translated text shows in that language, or English.
