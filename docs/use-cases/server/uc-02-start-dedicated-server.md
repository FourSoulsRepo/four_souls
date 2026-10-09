# UC-02: Start the dedicated server

**Module:** `cmd/server`, `internal/legal`
**Status:** In progress
**Actors:** Server owner
**Goal:** The server starts and shows the fan-game notice
**Preconditions:** The server binary is installed

## Main scenario (happy path)

1. The owner runs the server binary.
2. The console shows the fan-game notice with links.
3. The console shows app and rules engine versions.
4. The server listens on its port (4774 by default) for players.
5. Players host and join games in its lobby (UC-03, UC-04).
6. Ctrl+C or SIGTERM stops it; every connection closes cleanly.

## Alternative scenarios and errors

* **1a. The port is taken:** the server prints the error and exits.
* **1b. A bad flag or config file:** the server prints the problem and exits.

## Postconditions

* The notice and versions are printed.
* Saving running games on stop comes with step 6.9.
