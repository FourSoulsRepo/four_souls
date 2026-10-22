# UC-10: Change settings

**Module:** `internal/settings`, `app.go`, `frontend`
**Status:** Implemented
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
* **3b. Animations (ST-05):** Normal, Fast, Slow or Off.
  * Saved with the nickname; it applies to the next effect.
* **3c. Fullscreen (ST-01):** the window switches when saved.
  * The next start opens in the saved mode.
* **3d. Theme (ST-07):** no choice; the app has one dark theme (15.80).
  * An older settings file's theme is ignored.
* **3e. Game mat (ST-03):** Basement or Cathedral.
  * A preview shows the mat; the table uses it at once.
  * With the private images, every mat image is a choice (V-24).
  * An image takes the place of the color mat with its name.
  * A missing mat image shows that mat's colors, or Basement.
  * Default: the first mat image; Basement without images.
  * A saved Green felt, a mat before 15.81, gets the default (ST-08).
* **1a. First start (ST-06):** after the splash, a Welcome screen asks for the nickname.
  * Continue saves it and opens the main menu; it is asked only once.
* **3f. Dice (V-19, 15.64):** the look of the player's dice.
  * Ten looks, grouped as Holy and Unholy; default: Tear.
  * A preview shows the picked look's six faces.
  * The next hello sends it; every screen draws this player's rolls in it.
  * Without the private images there is no dice field; dice stay plain.
* **3g. Ask before ending the turn (V-46):** on by default.
  * Off: End turn ends the turn at once, even with plays left (UC-13, 7e).

## Postconditions

* The settings file holds the last saved settings.
* Every setting applies without a restart.
