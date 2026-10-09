# Card data tools

The card data in `pkg/card_db/data/` and the card images in `pkg/card_db/images/` are produced by the separate repository [FourSoulsRepo/tools](https://github.com/FourSoulsRepo/tools). The download tooling lives there on purpose: it fetches copyrighted pages and art, and keeping it apart keeps that work out of the game's code and history.

What the tools do, in short:

1. Download card pages from foursouls.com through a real Chrome (the site is behind a Cloudflare check).
2. Read copies per card from the official Tabletop Simulator table.
3. `scraper build -game <this repo>` merges both with `pkg/card_db/overrides/<set>.json` and writes `pkg/card_db/data/<set>.json` (format: ADR 004).
4. `scraper images -game <this repo>` writes 600 px WebP images to `pkg/card_db/images/<set>/` (git-ignored).

After regenerating data, validate it here:

```sh
go test ./internal/assets/ ./pkg/card_db/...
```

Manual fixes go into `pkg/card_db/overrides/<set>.json`, which the tools read but never write. Card data is committed here; card images never are.
