# 29. In-app update

* **Status:** Accepted
* **Date:** 2026-11-14
* **Authors:** @HardDie

---

## Context

1. Owner, 2026-11-14: an Update button in the main menu (B-08).
   1. Today users download every new version by hand.
   2. Reference: the owner's `ytmemchat_wails`, package `pkg/update`.
2. Releases come from `release.yml` on a `vX.Y.Z` tag.
   1. App archives: `four-souls-<ver>-<edition><target>.<ext>`.
   2. `<ver>` is the tag without `v`; it is also `version.App`.
   3. Editions: full (`embed cardimages`) and `for-poor-` (`embed`).
   4. Targets: `linux-amd64`, `linux-arm64`, `windows-amd64`, `windows-arm64`, `macos-universal`.
   5. macOS and Windows: `.zip`; Linux: `.tar.gz`.
   6. Each archive holds one folder: the app, `README.txt`, `LICENSE`, `NOTICE`.
   7. The app is `FourSouls`, `FourSouls.exe` or `FourSouls.app`.
3. A tag with a suffix (`v0.1.0-test`) is a prerelease.
4. The release has no checksum file yet.
5. The reference picks the archive by substring.
   1. `windows-amd64.zip` also matches `for-poor-windows-amd64.zip`.
   2. A full build could become for-poor, or the other way.
6. The reference code is GPL-3; this repo is MIT.
   1. The owner wrote both and allows it under MIT here.

## Considered options

1. **Port the reference flow: check, download, apply**
   1. Won. Tried in the owner's app; small; no new dependency.
2. **Wails or third-party self-update library**
   1. Lost. Wails v2 has none; libraries do not know our editions.
3. **Only a link to the releases page**
   1. Lost. Users still unpack and replace by hand.
4. **Pick the archive by its exact name**
   1. Won. Same OS, arch and edition, by construction.
5. **Pick the archive by substring, as in the reference**
   1. Lost. Swaps editions (Context 5).
6. **Edition from build tags**
   1. Won. Known at compile time; cannot be wrong.
7. **Edition from the images present at runtime**
   1. Lost. A guess; breaks if image loading changes.
8. **Check on button press**
   1. Won. Owner, 2026-11-14; no network call at start.
9. **Check once at start**
   1. Lost. Owner chose the button.

## Decision

Use options 1, 4, 6 and 8.

1. Package `internal/update`
   1. Ported from the reference; no Wails import.
   2. `Config`: owner `FourSoulsRepo`, repo `four_souls`, name `FourSouls`.
   3. Also: current version, edition, `GOOS`, `GOARCH`.
2. Source: GitHub API `repos/FourSoulsRepo/four_souls/releases/latest`.
   1. It skips drafts and prereleases.
   2. No token; anonymous rate limit is enough for a button.
3. Archive name: built, then compared exactly.
   1. `four-souls-<latest>-<edition><target>.<ext>`.
   2. `<latest>`: the tag without `v`.
   3. `<edition>`: empty for full, `for-poor-` for text cards.
4. Target from `GOOS` and `GOARCH`
   1. darwin, any arch: `macos-universal.zip`.
   2. windows amd64, arm64: `windows-<arch>.zip`.
   3. linux amd64, arm64: `linux-<arch>.tar.gz`.
   4. Anything else: no update.
5. Edition in `internal/version`
   1. `Edition` constant: `full` or `for-poor`.
   2. Two files with build tags set it.
   3. `embed && cardimages`: `full`; any other build: `for-poor`.
6. Versions
   1. Only `X.Y.Z` (optional `v`) counts.
   2. A suffix after `-` or `+` is ignored when comparing.
   3. `dev` or any other text: never updates.
   4. Update only when latest is greater than current.
7. Download
   1. Archive and `SHA256SUMS.txt` from the release.
   2. Size limit: 500 MB per archive; full builds are near 120 MB and grow.
   3. A wrong or missing hash drops the file.
   4. Progress: bytes done and total, a few times a second.
   5. Total from `Content-Length`; zero when unknown.
8. Unpack
   1. Into a temp folder; paths that leave it are refused.
   2. Payload: `FourSouls.app`, `FourSouls.exe` or `FourSouls`.
   3. Only the payload is installed.
9. Apply
   1. A helper script waits for the app to quit.
   2. It moves the old install aside and the new one in.
   3. Then it starts the new app and removes the old one.
   4. Unix: `/bin/sh`; Windows: PowerShell, retries while files are locked.
   5. macOS replaces the whole `.app`.
   6. The app quits right after starting the helper.
10. No update: the screen says why and links the releases page.
   1. Not an `X.Y.Z` version, or latest is not newer.
   2. No archive for this target and edition.
   3. The install folder is not writable.
11. Release workflow
   1. Job `release` writes `SHA256SUMS.txt` over `dist/*`.
   2. Uploaded with the archives.
12. Bindings and UI
   1. Methods: `Check`, `Download`, `ApplyAndQuit`, `OpenReleases`.
   2. Progress goes to the frontend as a Wails event.
   3. Only `frontend/src/bridge/` uses them.
   4. Main menu button **Update** opens the Update screen.

## Consequences

### Positive

1. Users update with three clicks.
2. The update keeps OS, arch and edition.
3. A broken or swapped download never installs.

### Negative and risks

1. `README.txt`, `LICENSE`, `NOTICE` next to the app stay old.
2. Unsigned builds: the OS may warn again after an update (B-03).
3. A read-only install (e.g. `/Applications` without rights) cannot update.
4. A changed archive name in `release.yml` breaks updates of old builds.
   1. Archive names are now a contract; change them only with a new ADR.

### Neutral

1. Server archives are not updated by the app.
2. Older builds without the button stay manual.
