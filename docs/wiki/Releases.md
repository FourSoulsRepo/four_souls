# Releases

A version tag builds every platform and makes a **draft** GitHub release (`.github/workflows/release.yml`). Nothing is public until the owner publishes the draft.

## Make a release

```sh
git tag v0.1.0
git push origin v0.1.0
```

A tag with a suffix, like `v0.1.0-test`, is a trial: its draft is marked as a pre-release, never the latest one. Delete trial drafts and tags when done.

The version shown in the app's main menu is the tag without the `v`.

## What it builds

| Archive | Contents |
|---|---|
| `four-souls-<version>-windows-amd64.zip`, `-windows-arm64.zip` | `four_souls.exe`, `four_souls_server.exe` |
| `four-souls-<version>-macos-universal.zip` | `four_souls.app`, `four_souls_server` (Intel and Apple silicon) |
| `four-souls-<version>-linux-amd64.tar.gz`, `-linux-arm64.tar.gz` | `four_souls`, `four_souls_server` |

Every archive also has `README.txt` (from `release/README.txt`), `LICENSE` and `NOTICE`.

The app is built with `-tags "embed cardimages"`; the server only with `embed` (it needs card text, not images). Linux builds use WebKitGTK 4.1 (`webkit2_41`), so players need `libwebkit2gtk-4.1` (Ubuntu 22.04 and later).

## Card images

A release needs the private card images. The workflow checks them out with `.github/actions/card-images`, using the `CARD_DB_DEPLOY_KEY` repository secret (a read-only deploy key of `FourSoulsRepo/card_db`). Without the secret the release fails on purpose; CI builds still work without it.

## Shared build steps

`.github/actions/build` builds the app and the server for one target; CI (`ci.yml`, job `app`) and the release workflow both use it.
