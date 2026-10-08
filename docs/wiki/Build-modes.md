# Build modes

## Embedded vs. external files

App files (game mats, sounds and similar) are reached through `internal/assets.FS()`.

| Build | Where files come from |
|-------|----------------------|
| `-tags embed` | packed into the binary from `internal/assets/files/` |
| default (`wails dev`) | `$FOUR_SOULS_ASSETS`, else `assets/` next to the binary, else `internal/assets/files/` |

External mode lets you change files and restart without rebuilding. A missing folder is not fatal: reads simply fail with "not exist", and the app still starts.

Release builds use embed mode:

```sh
wails build -tags embed
```

## Versions

* App version: `internal/version.App`, `dev` by default, stamped at link time:

  ```sh
  wails build -ldflags "-X github.com/FourSoulsRepo/four_souls/internal/version.App=1.2.3"
  ```

* Rules engine version: the constant `rulesengine.Version` in `pkg/rules_engine/version.go`.

The main menu shows both in its corner; the dedicated server prints both at start.

## Platforms

CI builds the app and the server for Linux and Windows (x64 and ARM) and a universal macOS app. Linux on Ubuntu 24.04 needs `libgtk-3-dev`, `libwebkit2gtk-4.1-dev` and the build tag `webkit2_41`.
