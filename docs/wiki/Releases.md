# Releases

A version tag builds every platform and makes a **draft** GitHub release (`.github/workflows/release.yml`). Nothing is public until the owner publishes the draft.

## Make a release

1. Move the `Unreleased` items in `CHANGELOG.md` under a new `## X.Y.Z - YYYY-MM-DD` section and commit.
2. Tag that commit and push the tag:

```sh
git tag v0.1.0
git push origin v0.1.0
```

A tag with a suffix, like `v0.1.0-test`, is a trial: its draft is marked as a pre-release, never the latest one. Delete trial drafts and tags when done.

The version shown in the app's main menu is the tag without the `v`.

## Versions

* **App:** semantic versions, from the tag (`v1.2.3` → `1.2.3`). Minor for new features, patch for fixes; 1.0.0 is for the owner to call.
* **Rules engine:** `rulesengine.Version` in `pkg/rules_engine/version.go`, bumped with every change to rules behavior. Saved games only load on the same engine version (N-11).
* The release notes start with both: **App 1.2.3 · Rules engine 0.4.1**, then the changelog section (`cmd/relnotes`). A tag that is not a version, or a release without its changelog section, fails before anything is built.
* The first release of an engine version also tags it `rules_engine/vX.Y.Z`.

## What it builds

| Archive | Contents |
|---|---|
| `four-souls-<version>-windows-amd64.zip`, `-windows-arm64.zip` | `FourSouls.exe` |
| `four-souls-<version>-macos-universal.zip` | `FourSouls.app` (Intel and Apple silicon) |
| `four-souls-<version>-linux-amd64.tar.gz`, `-linux-arm64.tar.gz` | `FourSouls` |
| `four-souls-server-windows-amd64.zip`, `-windows-arm64.zip` | `four_souls_server.exe` |
| `four-souls-server-macos-universal.zip` | `four_souls_server` (Intel and Apple silicon) |
| `four-souls-server-linux-amd64.tar.gz`, `-linux-arm64.tar.gz` | `four_souls_server`, statically linked |

Each platform also gets a `for-poor` archive, e.g. `four-souls-<version>-for-poor-linux-amd64.tar.gz`: the same game without card images, just for fun. Cards, backs, dice and mats then show as text and plain colors.

Every archive also has `README.txt` (from `release/README.txt`, or `release/SERVER.txt` for the server), `LICENSE` and `NOTICE`.

Server archives carry no version in their names (8.17), so `https://github.com/FourSoulsRepo/four_souls/releases/latest/download/four-souls-server-linux-amd64.tar.gz` always gets the newest server; the Dockerfile relies on it. "Latest" skips drafts and pre-releases. App archives keep their versions.

The app is built with `-tags "embed cardimages"` after `go generate ./pkg/card_db/images` (`for-poor`: `embed` only); the server only with `embed` (it needs card text, not images) and `CGO_ENABLED=0`, so its Linux binaries run on any distribution, Alpine included. Linux builds use WebKitGTK 4.1 (`webkit2_41`), so players need `libwebkit2gtk-4.1` (Ubuntu 22.04 and later).

## App identity

Every release has the same identity (B-03), so the OS treats it as the same app:

* Bundle ID `io.github.foursoulsrepo.foursouls` on macOS (`build/darwin/Info.plist`) and in the Windows manifest (`build/windows/wails.exe.manifest`).
* Product name "Four Souls", company "FourSoulsRepo" (`wails.json`, `info`); Windows shows them in the file's details.
* The build stamps the release version into `wails.json` `info.productVersion`; Wails puts it into `Info.plist` (`CFBundleVersion`, `CFBundleShortVersionString`) and the Windows file and product versions (`build/windows/info.json`). These take numbers only: a trial tag drops its suffix (`v1.2.3-test` shows 1.2.3), and CI builds show 0.0.0.
* No code signing; see `README.txt` for the first start.

`identity_test.go` fails if any of these drift.

## Card images

A release needs the private card images. The workflow checks them out with `.github/actions/card-images`, using the `CARD_DB_DEPLOY_KEY` repository secret (a read-only deploy key of `FourSoulsRepo/card_db`). Without the secret the release fails on purpose; CI builds still work without it.

## Shared build steps

`.github/actions/build` builds the app for one target; CI (`ci.yml`, job `app`) and the release workflow both use it.

`release/server.sh <version> <dir>` builds and packs every server archive on one macOS runner (it needs `lipo` for the universal binary); CI (job `server`), the release workflow and `make server-archives` use it. CI also checks that the Linux servers are statically linked.

## Server setups

`deploy/` holds the operator setups (ADR 018): the systemd unit, the Dockerfile, and compose with Caddy. The Dockerfile downloads the server of the latest release, so a release needs no change there. A release must be published as a full release, not a pre-release, or `releases/latest` skips it.

The base images are pinned and bumped by hand: `ARG ALPINE=` in `deploy/docker/Dockerfile` and `image: caddy:…` in `deploy/docker/docker-compose.yml`. Pick a tag at least two weeks old, then check both files once: `docker build` and `docker compose up -d --build` with `DOMAIN=localhost` in `.env` (`curl -k https://localhost/healthz` answers `ok`). Operators get the new pins by downloading the files again, as their wiki pages say.
