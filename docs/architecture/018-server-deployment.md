# 18. Server deployment from releases

* **Status:** Accepted
* **Date:** 2026-10-21
* **Authors:** @HardDie

---

## Context

1. Operators run the dedicated server on a VPS (N-19).
2. They want a ready setup, not a page of manual steps.
   1. Some use systemd, some Docker, some want HTTPS.
3. Each setup should be one or a few files to download.
   1. No clone of the repository; no Go or Node on the VPS.
4. The server is plain Go; it needs no cgo, Wails or card images.

## Considered options

1. **Dockerfile builds from source**
   1. Lost. Needs the whole repository and Go in the build.
2. **An image on a registry (GHCR)**
   1. Lost. Owner's choice: no published image for now.
3. **Dockerfile downloads the latest release archive**
   1. Won. One file; the server matches a real release.
4. **Ask the GitHub API for the latest tag**
   1. Lost. `releases/latest/download/<name>` needs no API.
   2. It works only with names without a version.

## Decision

Use option 3.

1. Server archives apart from the app (`release/server.sh`).
   1. `four-souls-server-<os>-<arch>`, no version in the name.
   2. Built with `CGO_ENABLED=0`: static Linux binaries.
   3. App archives keep their versions.
2. Settings by ENV: `FOUR_SOULS_` and the flag name in capitals.
   1. Order: defaults, JSON file, ENV, flags.
   2. An empty variable counts as unset.
3. Three setups under `deploy/`, one wiki page each.
   1. `systemd/four-souls.service`.
   2. `docker/Dockerfile`: Alpine, non-root user, `/data` volume.
   3. `docker/docker-compose.yml` with Caddy for HTTPS (ADR 017).
4. The Dockerfile `ADD`s the `latest/download` link.
   1. Docker fetches it on every build.
   2. A new release reruns the later steps; an old one keeps the cache.
   3. `--build-arg VERSION=1.2.3` downloads that tag instead.
5. `GET /healthz` answers 200 while the lobby goroutine answers.
   1. The Dockerfile and compose healthchecks use it.
   2. Caddy starts once the server is healthy.
6. Base images are pinned: `alpine:3.24.2`, `caddy:2.11.7-alpine`.
   1. Bumped by hand in `deploy/docker/`.

## Consequences

### Positive

1. An operator downloads one to four files and runs one command.
2. Updating is the same build again; nothing else to edit.
3. The Linux server runs on any distribution.

### Negative and risks

1. "Latest" skips pre-releases; a server needs a full release.
2. A pinned version works only from the first release with server archives.
3. Every build downloads the archive to compare it (about 7 MB).
4. Docker does not restart an unhealthy container by itself.

### Neutral

1. `release/server.sh` needs macOS, for the universal macOS binary.
2. Raw file links in the wiki point at `main`.
