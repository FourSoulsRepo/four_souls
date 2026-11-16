# UC-44: Manage records

**Module:** `bindings/records`, `internal/replay`, `pkg/record`, `frontend`
**Status:** Implemented
**Actors:** Player
**Goal:** The player finds a saved match, opens it, keeps a copy or deletes it
**Preconditions:** The app is on the main menu

## Main scenario (happy path)

1. The player presses **Records** in the main menu.
2. The Records screen lists the records folder, newest first.
   1. It holds hosted games' records and saved records.
   2. Each row: start time and players.
   3. Under it: the winners or "Unfinished", and the step count.
3. The player presses **Watch** on a row.
   1. The app opens the record without the rules engine (ADR 031).
   2. It plays the record on the table (UC-45).
4. The player presses **Back to records** on the replay bar.
5. The player presses **Export** on a row.
   1. A save dialog asks where.
   2. The app copies the file; the status line names the copy.
6. The player presses **Delete** on a row.
   1. The row asks "Delete this record for good?".
   2. **Delete** deletes the file; the row goes.
7. The player presses **Back** or Esc to return to the menu.

## Alternative scenarios and errors

* **2a. No records yet:** the screen says so.
* **2b. A format 1 record:** "made before replays"; no Watch button.
* **2c. A record from a newer app:** "update the app"; no Watch button.
* **2d. A broken file:** its file name and "damaged"; no Watch button.
* **3a. Open a file:** a file dialog picks a record anywhere; then as 3.
* **3b. The file can't be watched:** the status line says why.
* **5a. The player cancels the dialog:** nothing is copied.
* **6a. The player presses Keep:** the record stays.

## Postconditions

* The list shows the folder as it is now.
* A watched record is open in the app for the player (UC-45).
