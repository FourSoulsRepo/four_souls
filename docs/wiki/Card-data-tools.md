# Card data tools

The card data in `pkg/card_db/data/` and the card images in `pkg/card_db/images/` (a submodule, see below) are produced by the separate repository [FourSoulsRepo/tools](https://github.com/FourSoulsRepo/tools). The download tooling lives there on purpose: it fetches copyrighted pages and art, and keeping it apart keeps that work out of the game's code and history.

What the tools do, in short:

1. Download card pages from foursouls.com through a real Chrome (the site is behind a Cloudflare check).
2. Read copies per card from the official Tabletop Simulator table.
3. `scraper build -game <this repo>` merges both with `pkg/card_db/overrides/<set>.json` and writes `pkg/card_db/data/<set>.json` (format: ADR 004).
4. `scraper images -game <this repo>` writes 600 px WebP images to `pkg/card_db/images/<set>/`.

After regenerating data, validate it here:

```sh
go test ./internal/assets/ ./pkg/card_db/...
```

Manual fixes go into `pkg/card_db/overrides/<set>.json`, which the tools read but never write. Card data is committed here; card images never are.

## The images submodule

`pkg/card_db/images` is a git submodule pointing to the **private** repository `FourSoulsRepo/card_db`. It holds only images (one folder per set) and a README.

* With access: `git clone --recurse-submodules …` or `git submodule update --init`.
* Without access (public clones, forks): the folder stays empty, everything still builds, and cards show as text.
* Release builds pack images with `-tags "embed cardimages"`; plain `-tags embed` packs only card data.
* CI: `.github/actions/card-images` checks the images out at the pinned submodule commit, using the secret `CARD_DB_DEPLOY_KEY` (a read-only deploy key of `FourSoulsRepo/card_db`). Every app build runs twice, text-only and with images; without the secret (e.g. pull requests from forks) the image builds skip themselves.
* Its history is a single commit: after regenerating images, commit inside the submodule with `git commit --amend`, force-push it, then commit the new submodule pointer here.
