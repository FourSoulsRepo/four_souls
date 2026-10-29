# Server from a local build

Run a server you built yourself, from your own checkout, in the Docker setup of [Server in Docker](Server-in-Docker) or [Server with Docker Compose and Caddy](Server-with-Docker-Compose-and-Caddy). Use it to try a fix on your server before it is in a release.

The idea: keep the release image as the base and put your binary on top of it in a second, tiny image `four-souls-server:local`. Everything else (user, `/data` volume, healthcheck, settings) stays as in the release image. Going back to the release is one command.

## 1. Build the binary on your computer

The server is plain Go; with `CGO_ENABLED=0` it is a static Linux binary that runs on Alpine. Cross-compiling works from macOS, Windows or Linux; the private card images are not needed (the server needs card text only).

Check the server's CPU first: `uname -m` on the server says `x86_64` (use `amd64`) or `aarch64` (use `arm64`).

```sh
cd four_souls_wails
mkdir -p .cache/server-local
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags embed \
  -ldflags "-X github.com/FourSoulsRepo/four_souls/internal/version.App=$(git describe --tags --dirty --always)" \
  -o .cache/server-local/four_souls_server ./cmd/server
```

* `-tags embed` puts the card text into the binary, as the release build does (`release/server.sh`).
* The version comes from `git describe`, for example `v0.2.0-4-g1d363d3` or `…-dirty` with uncommitted changes. The server prints it at start (the `app` log line), so you can tell your build from a release.
* `make server-archives VERSION=…` builds the full release archives instead; it needs macOS. The direct `go build` above is enough here.

## 2. Copy it to the server

Into the folder that already holds the release `Dockerfile` (`~/four-souls` in the other pages):

```sh
scp .cache/server-local/four_souls_server you@server:four-souls/
```

## 3. Build the local image on the server

The release image `four-souls-server` must exist on the server; it does once you have followed [Server in Docker](Server-in-Docker) or the Compose page. Next to the release `Dockerfile`, create `Dockerfile.local`:

```dockerfile
# A locally built server on top of the release image.
FROM four-souls-server
COPY four_souls_server /usr/local/bin/four_souls_server
```

Build it. `COPY` keeps the file's mode, so make it executable first; this also works with Docker's legacy builder (no buildx), which does not know `COPY --chmod`:

```sh
cd ~/four-souls
chmod 755 four_souls_server
docker build -f Dockerfile.local -t four-souls-server:local .
```

`FROM four-souls-server` takes the image from the server itself. If Docker says it tries to pull `docker.io/library/four-souls-server`, the release image is missing: run `docker build -t four-souls-server .` first.

## 4a. Run it with plain Docker

Players connect as they do to any server: **Play → Online**, then the server's address.

![The Online screen, with a server that answers](images/screen-online-status.png)

Replace the container as in Updating on [Server in Docker](Server-in-Docker), with the `:local` image. The same volume keeps records and running games.

```sh
docker stop four-souls && docker rm four-souls
docker run -d --name four-souls --restart unless-stopped \
  -p 4774:4774 -v four-souls-data:/data \
  -e FOUR_SOULS_NAME=Basement \
  four-souls-server:local
docker logs four-souls | grep '^app'   # your git describe version
```

## 4b. Run it with Docker Compose

Next to `docker-compose.yml`, create `docker-compose.override.yml`. Compose reads it on its own and merges it into the main file:

```yaml
# The server from Dockerfile.local; delete this file to go back to the release.
services:
  server:
    build:
      context: .
      dockerfile: Dockerfile.local
    image: four-souls-server:local
```

```sh
docker compose up -d --build
docker compose logs --no-log-prefix server | grep '^app'   # your version
```

While the override file exists, `docker compose up -d --build` builds only the local image; it does not look for a new release.

## Updating your build

Repeat steps 1 and 2, then the build and run of step 3 and 4a, or `docker compose up -d --build` for 4b. `docker image prune -f` removes the old images afterwards.

## Back to the release

* Plain Docker: run the Updating commands of [Server in Docker](Server-in-Docker); they rebuild `four-souls-server` and start the container from it.
* Compose: `rm docker-compose.override.yml`, then `docker compose up -d --build`.

A saved game loads only on a server with the same rules engine version (N-11). If your build changed the engine version, games saved by it do not load in the release, and the other way round. Finish running games first when in doubt.

To refresh the base image under your build (new release, new Alpine), run `docker build -t four-souls-server .` (plain Docker) and then build `:local` again. With Compose, use the same plain `docker build` command, since the override file stops Compose from building the release image.
