# UC-10: Change settings

**Module:** `internal/settings`, `app.go`, `frontend`
**Status:** In progress
**Actors:** Player
**Goal:** The player's settings are kept across app restarts
**Preconditions:** The app is on the main menu

## Main scenario (happy path)

1. The player opens **Settings** from the main menu.
2. The app shows the saved settings.
   1. Now only the nickname (ST-06).
3. The player types a nickname and presses **Save**.
4. Go trims it and writes the settings file.
   1. File: `settings.json` in the user's config folder, under `FourSouls`.
5. The app shows the stored nickname and "Saved.".
6. After a restart, Settings shows the same nickname.

## Alternative scenarios and errors

* **3a. Bad nickname:** empty after trimming, over 20 characters, or with control characters.
  * The app shows "nickname must be 1–20 characters".
  * The saved settings stay as they were.
* **4a. Missing or broken file:** the app starts with default settings.
* **6a. Back:** **Back** returns to the main menu without saving.

## Postconditions

* The settings file holds the last saved settings.
* Planned in 7.10: theme, window mode, game mat, animation speed, nickname asked on first start.
