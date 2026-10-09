# UC-08: Get a match record

**Module:** `internal/server`, `pkg/record`
**Status:** In progress
**Actors:** Players; server owner
**Goal:** Every match is kept, and its players can download it after the end
**Preconditions:** A game ran on the server

## Main scenario (happy path)

1. When a game starts, the server opens its record (ADR 007).
2. Every applied step goes into it, with all hands and a checksum.
3. When the game ends, the record is closed as finished.
4. A player of that game asks for it with their token (`/record`).
5. The server sends the `.fsrec` file.

## Alternative scenarios and errors

* **3a. The server stops before the end:** the record is closed as unfinished (RP-08).
* **3b. The server crashes:** the record ends at its last step and still reads.
* **4a. The game is not over:** refused; a record would reveal hands (RP-05).
* **4b. Someone who did not play asks:** refused.
* **5a. The record is older than the retention:** it was deleted (RP-07).

## Postconditions

* The record replays the match exactly (A-08).
