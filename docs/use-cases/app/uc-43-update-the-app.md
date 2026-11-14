# UC-43: Update the app

**Module:** `internal/update`, `updater.go`, `frontend`
**Status:** Implemented
**Actors:** Player
**Goal:** The app installs the newest release of its own build and restarts
**Preconditions:** The app is on the main menu; a release build (ADR 029)

## Main scenario (happy path)

1. The player presses **Update** in the main menu.
2. The Update screen shows this version.
   1. Nothing is asked of GitHub yet (owner, 2026-11-14).
3. The player presses **Check for updates**.
   1. The status line says "Checking for updates…".
4. The app shows the latest version and its release notes.
   1. Long notes scroll in their own box; the page does not.
   2. The status line says a new version is ready.
5. The player presses **Download**.
   1. A bar fills; under it the percent and MB.
   2. Check and Download are off while it runs.
6. The download is checked against `SHA256SUMS.txt`.
   1. The status line says it is downloaded and checked.
7. The player presses **Install and restart**.
   1. The app starts the installer and closes.
   2. The new version starts by itself.
8. The player presses **Back** at any time to return to the menu.

## Alternative scenarios and errors

* **4a. Already the latest version:** "You have the latest version."
* **4b. A test or `dev` build:** the reason, and **Open the releases page**.
* **4c. No archive for this system:** the reason and the releases button.
* **4d. The app's folder is not writable:** the reason and the releases button.
* **4e. GitHub can't be reached:** the error, and the releases button.
* **5a. The server sends no size:** MB only, no percent.
* **5b. The download breaks off:** the error; Download again, or the releases page.
* **6a. A wrong or missing checksum:** the file is dropped; the error says so.
* **7a. The installer can't start:** the error and the releases button.

## Postconditions

* The new version runs, same OS, arch and edition.
* Or the player knows why not, with the releases page one press away.
