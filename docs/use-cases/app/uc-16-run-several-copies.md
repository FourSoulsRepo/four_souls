# UC-16: Run several copies on one computer

**Module:** `internal/instance`, `main.go`, `app.go`
**Status:** Implemented
**Actors:** Players sharing one computer; a tester
**Goal:** Each copy plays as its own player
**Preconditions:** One copy of the app is running

## Main scenario (happy path)

1. The player starts the app again.
   1. Windows, Linux: start it as usual.
   2. macOS: `open -n FourSouls.app`.
2. The new copy takes the lowest free slot.
   1. A slot is a locked file `instance-N.lock` in the config folder.
   2. The first copy holds slot 1.
3. The new copy uses its own settings file.
   1. Slot 1: `settings.json`, as before.
   2. Slot N: `settings-N.json`.
4. A new slot file copies the first copy's choices.
   1. Nickname and the last game's seat are left empty.
5. The window title shows the slot: "Four Souls (2)".
6. The Welcome screen asks the new copy for a nickname.
7. Each copy joins a game with its own seat token.

## Alternative scenarios and errors

* **2a. All 16 slots busy:** the copy shares the first copy's settings.
* **2b. A copy closes or crashes:** the OS drops its lock.
  * The next copy started takes that slot and its settings.
* **7a. Two copies host:** the second fails; the port is taken.

## Postconditions

* Each running copy holds one slot.
* No copy takes another copy's seat on restart.
