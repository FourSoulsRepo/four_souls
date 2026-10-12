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
* **3d. Theme (ST-02):** Dark or Light, for menus and panels.
  * The game table keeps its dark overlays on the mat.
* **3e. Game mat (ST-03):** Green felt, Basement or Cathedral.
  * A preview shows the mat; the table uses it at once.
  * With the private images, every mat image is a choice (V-24).
  * An image takes the place of the color mat with its name.
  * A missing mat image shows that mat's colors, or Basement.
  * Default: the first mat image; Basement without images.
* **1a. First start (ST-06):** after the splash, a Welcome screen asks for the nickname.
  * Continue saves it and opens the main menu; it is asked only once.

## Postconditions

* The settings file holds the last saved settings.
* Every setting applies without a restart.
