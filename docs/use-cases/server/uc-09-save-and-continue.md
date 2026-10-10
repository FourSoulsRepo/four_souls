# UC-09: Save and continue a game

**Module:** `internal/server`, `pkg/record`
**Status:** In progress
**Actors:** Host; the other players
**Goal:** Friends finish a game on another day (N-11)
**Preconditions:** A game is running (UC-03)

## Main scenario (happy path)

1. The host saves the game, at any moment, even with the stack open.
2. The server writes the save; every player gets "saved" and is back in the lobby.
3. Another day, the lobby lists the save with its players' nicknames.
4. One of its players loads it; a table opens with the saved seats.
5. Each player takes their seat back: by token, or by nickname.
6. When every seat is back and ready, the game goes on exactly as saved.
7. The match record goes on in the same file.

## Alternative scenarios and errors

* **1a. A player who is not the host saves:** refused with `not_host`.
* **4a. The save is from another rules engine version:** refused with `save_version`.
* **4b. Someone who did not play loads or joins it:** refused with `not_in_save`.
* **5a. A player leaves the table before the start:** the seat waits for them.

## Postconditions

* The game is the same as when saved; the save file is removed once it goes on.
* Players never get the save; it holds every hand (RP-05).
