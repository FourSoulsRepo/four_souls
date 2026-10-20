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
* **1c. A server name (N-18):** `-name` or `"name"` in the config file.
  * Players see it instead of the address (UC-04).
  * A name over 40 characters or with control characters stops the start (1b).
* **1d. ENV settings (N-19):** `FOUR_SOULS_` and a flag's name in capitals.
  * Example: `FOUR_SOULS_PORT=5000` acts like `-port 5000`.
  * ENV overrides the config file; flags override ENV.
  * An empty variable counts as unset.
  * A bad value names the variable and stops the start (1b).
* **1e. In Docker (8.19):** the owner builds `deploy/docker/Dockerfile`.
  * The image gets the latest release's server; no source code needed.
  * Records and saves live in the `/data` volume.
  * `GET /healthz` answers 200 while the lobby runs; Docker marks it healthy.
  * `docker stop` stops it as SIGTERM does (6).
* **1f. Compose with Caddy (8.20):** `docker compose up -d --build`.
  * `.env` sets `DOMAIN` and the server's settings.
  * Caddy starts once the server is healthy and gets an HTTPS certificate.
  * Players type only the domain and connect with `wss://` (UC-27).
  * `docker compose ps` shows whether the server is healthy.

## Postconditions

* The notice and versions are printed.
* On stop, running games' records are closed as unfinished (UC-08).
