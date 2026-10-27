# Build modes

## Embedded vs. external files

App files (sounds and similar) are reached through `internal/assets.FS()`.

| Build | Where files come from |
|-------|----------------------|
| `-tags embed` | packed into the binary from `internal/assets/files/` |
| default (`wails dev`) | `$FOUR_SOULS_ASSETS`, else `assets/` next to the binary, else `internal/assets/files/` |

Card data and images (`pkg/card_db`, see ADR 004) follow the same idea through `internal/assets.CardsFS()`: card data is embedded with `-tags embed`, and images too with `-tags "embed cardimages"` (needs the private images submodule, prepared first with `go generate ./pkg/card_db/images`); otherwise `$FOUR_SOULS_CARDS`, else `cards/` next to the binary, else `pkg/card_db/` (for `wails dev`).

In every mode, embedded too, `$FOUR_SOULS_IMAGES` can name one more folder of card images, such as a translator's own images repo (16.12, [Translate cards](Translate-cards#card-images)). It is laid out like `images/`: `<lang>/<set>/<file>`, e.g. `de/b2/the_d6.webp`. `carddb.WithImages` lays it over the card files: a file there wins over the built-in one, and a card it lacks falls back as before, to the built-in image in that language, then to English. Only language folders (two or three small letters) are read from it; card data, `assets/` and anything else there are ignored. The folder is opened with `os.Root`, so no read leaves it, not even through a symbolic link. A folder that is not there is logged and skipped. The app reads it at start; restart after adding images.

The app serves card images to the frontend under `/cardimg/<lang>/<set>/<file>.webp` (an asset server middleware in `cards.go`). Card data names images without the language (`b2/the_d6.webp`); `carddb.FindImage` adds `<lang>/` and falls back to `en/` (16.6, ADR 021). `App.Cards()` gives each card's text, its English image URL in `image`, and in `images` one URL per language folder under `images/` (English when the card has none of its own there), or no URL when the file is missing.

Dice textures, game mats and card backs are private images too, under `/cardimg/assets/`. `App.Mats()` lists every `.webp` in `images/assets/mats/`, so a new mat needs only a new file. Without the images the list is empty and the table uses the color mats. `App.Backs()` gives `images/assets/backs/<kind>.webp` for loot, treasure and monster; without them face-down cards are plain. `App.Icons()` gives `images/assets/icons/coin.webp` and `sword.webp`; without them the coin and the sword are drawn. It also gives `images/assets/banner.webp`, the foursouls.ru logo; without it the home screen shows its text title (15.41). Its `sky` lists the six layers in `images/assets/background/`, the foursouls.ru starry sky; when any is missing the list is empty and the menu screens keep their plain color (15.42). `App.Fonts()` gives the game's TrueType fonts in `images/assets/fonts/`, served as `font/ttf`; without them the system font shows (15.20). `App.Sounds()` gives the Hearthstone damage sounds in `images/assets/sounds/`, served as `audio/mpeg`; without them the game is silent (15.50).

External mode lets you change files and restart without rebuilding. A missing folder is not fatal: reads simply fail with "not exist", and the app still starts.

Release builds use embed mode:

```sh
wails build -m -tags embed
```

`-m` skips `go mod tidy`. Tidy reads files of every build tag, so without the private images submodule it fails on the import in `pkg/card_db/embed_images.go` (19.59). `wails dev -m` likewise.

## Versions

* App version: `internal/version.App`, `dev` by default, stamped at link time:

  ```sh
  wails build -m -ldflags "-X github.com/FourSoulsRepo/four_souls/internal/version.App=1.2.3"
  ```

* Rules engine version: the constant `rulesengine.Version` in `pkg/rules_engine/version.go`.

The main menu shows both in its corner; the dedicated server prints both at start.

## Platforms

CI builds the app and the server for Linux and Windows (x64 and ARM) and a universal macOS app. Linux on Ubuntu 24.04 needs `libgtk-3-dev`, `libwebkit2gtk-4.1-dev` and the build tag `webkit2_41`.
