# Server with Docker Compose and Caddy

The server behind [Caddy](https://caddyserver.com), which gets an HTTPS certificate for your domain by itself and renews it. Traffic is encrypted (`wss://`), and players type only the domain. Read [Dedicated server on a VPS](Dedicated-server-on-a-VPS) first: the ports, options and safety are there.

You need:

* Docker with the Compose plugin (`docker compose version` works).
* A domain, e.g. `games.example.com`, whose DNS `A` record points at the VPS.
* Ports 80 and 443 open; Caddy needs both for the certificate. Port 4774 stays closed.

## 1. Download

The four files live in the repository under `deploy/docker/`; put them into one folder:

```sh
mkdir -p ~/four-souls && cd ~/four-souls
base=https://raw.githubusercontent.com/FourSoulsRepo/four_souls/main/deploy/docker
curl -fsSLO "$base/Dockerfile"
curl -fsSLO "$base/docker-compose.yml"
curl -fsSLO "$base/Caddyfile"
curl -fsSL -o .env "$base/.env.example"
```

## 2. Configure

Edit `.env`: set `DOMAIN` to your domain, and optionally the server's name and other settings. An empty value keeps the default.

```sh
DOMAIN=games.example.com
FOUR_SOULS_NAME=Basement
FOUR_SOULS_RETENTION=
FOUR_SOULS_PRESENCE=
FOUR_SOULS_EMPTY_MINUTES=
FOUR_SOULS_GAME_HOURS=
```

Another `FOUR_SOULS_*` setting from the Options table goes into `.env` and into the `environment` list of `server` in `docker-compose.yml`.

## 3. Start

```sh
docker compose up -d --build
docker compose ps             # server shows (healthy); caddy starts after it
docker compose logs -f        # both outputs; Caddy says when it has the certificate
```

Check it from any computer: `https://games.example.com/healthz` answers `ok`. Players type `games.example.com` in **Play → Online**.

* The healthcheck asks the server's `/healthz` every 30 seconds; `docker compose ps` shows the state, and Caddy starts only once the server is healthy. Docker does not restart an unhealthy container by itself; `restart: unless-stopped` restarts one that exits.
* Records and saved games live in the `data` volume; Caddy's certificates in `caddy_data`. They survive updates.
* To let players also join without HTTPS, uncomment `ports` of `server` in `docker-compose.yml` and open port 4774.

## Updating

```sh
cd ~/four-souls
docker compose up -d --build   # builds the newest release, if any, and replaces the server
docker compose logs --no-log-prefix server | grep '^app'   # the new version
docker image prune -f          # removes the old image
```

The running games go on: the `data` volume keeps them over the replacement, and players see a short "Reconnecting…" (see Updating on [Dedicated server on a VPS](Dedicated-server-on-a-VPS)).

As in [Server in Docker](Server-in-Docker), the build always asks for the latest server and downloads only when a new release is out. To update Caddy or the Alpine base, download the three files again (not `.env`, which holds your settings), then run `docker compose pull caddy` and the same commands.

## Stop

```sh
docker compose down            # keeps the volumes; add -v to delete records and certificates too
```
