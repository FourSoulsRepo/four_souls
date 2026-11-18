# UC-48: Run the diagnostics

**Module:** `internal/diag`, `bindings/probe`, `frontend`
**Status:** Implemented
**Actors:** Player
**Goal:** A report says how fast the game draws and what may slow it
**Preconditions:** The game is installed; a terminal is open in its folder

## Main scenario (happy path)

1. The player starts `FourSouls --diagnose`.
   1. The game itself does not start.
   2. The terminal says test windows will open one after another.
2. The app reads the system (ADR 035 item 2).
   1. Linux: WebKitGTK version, session, GPUs, groups.
   2. Webview switches already set.
3. A test window opens for each variant.
   1. Linux: five variants; other systems: `default` only.
   2. The window is 1280×800 and shows the main menu over the sky.
   3. No splash, no nickname screen; nothing is saved.
   4. After 1 s it counts frames for 5 s, then closes.
   5. The terminal prints each variant's frames per second.
4. The report prints to the terminal.
   1. The system, one row per variant, the best variant.
   2. A command to start the game with the best switch.
   3. The problems found.
5. The report is saved as `diagnose-<date>-<time>.txt` in the app folder.
   1. The last line names the file.

## Alternative scenarios and errors

* **3a. A window sends nothing in 30 s:** its row says "failed"; the next one runs.
* **3b. A window crashes:** its row says "failed" with the last error line.
* **4a. No variant sent a result:** "Best: none".
* **4b. The best is `default`:** no command is printed.
* **4c. A probe's GPU name:** WebKit hides it; the row says "hidden by WebKit" (19.228).
* **5a. The file can't be written:** the error; exit code 1.

## Postconditions

* The report is on screen and in the file; no setting changed.
