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

## Alternative scenarios and errors

* **2a. First start:** no language is saved yet.
  * The app takes the OS language, if it is listed.
  * Else English.
* **3a. Russian (16.2):** it has no interface file yet.
  * The texts stay English; card names follow in 16.5.
* **5a. A line missing in the language:** the English line shows.
* **6a. A saved language this build lacks:** as on the first start.
* **2b. Test language `qps`:** listed in development builds only.
  * Accented, longer English in `[!! … !!]` (ADR 021).
  * Plain English on screen then means a missed string.

## Postconditions

* The settings file holds the chosen language.
* Every translated text shows in that language, or English.
