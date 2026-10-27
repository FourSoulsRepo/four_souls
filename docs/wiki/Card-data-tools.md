# Card data tools

The card data in `pkg/card_db/data/` and the card images in `pkg/card_db/images/` (a submodule, see below) are produced by the separate repository [FourSoulsRepo/tools](https://github.com/FourSoulsRepo/tools). The download tooling lives there on purpose: it fetches copyrighted pages and art, and keeping it apart keeps that work out of the game's code and history.

What the tools do, in short:

1. Download card pages from foursouls.com through a real Chrome (the site is behind a Cloudflare check).
2. Read copies per card from the official Tabletop Simulator table.
3. `scraper copies` finds the numbered page of every copy on foursouls.com, so each copy gets its own image.
4. `scraper build -game <this repo>` merges all of it with `pkg/card_db/overrides/<set>.json` and writes `pkg/card_db/data/<set>.json` (format: ADR 004 and ADR 010).
5. `scraper images -game <this repo>` writes 600 px WebP images to `pkg/card_db/images/en/<set>/`, one per copy when copies differ.
6. `scraper fonts -game <this repo>` writes the game's fonts to `pkg/card_db/images/assets/fonts/`.
7. `scraper ru-names -migrations <dir> -game <this repo>` writes Russian card names and their translators (16.5); see below.

Sets are built oldest first (Base Game, Gold Box, Four Souls+, …), because `build` looks at the older sets' files: a card with the same name, type, stats and text becomes a reprint entry that adds copies and images to the older card, and a different card whose ID is taken gets the set code as a suffix (`isaac_fsp2`). See ADR 010.

In the official TTS table, Gold Box is the bag named "Kickstarter" (`scraper tts -origin g2 -bag Kickstarter`).
Requiem is the bag "Requiem" (`scraper tts -origin r -bag Requiem`). A two-sided card is named "Front/Back" on the site and only by its front in TTS (`The Enigma`); the tools match both. When the site lists one card twice (Requiem's `r-eden` and `r-eden_2`), the build makes one card with both images.

After regenerating data, validate it here:

```sh
go test ./internal/assets/ ./pkg/card_db/...
```

Manual fixes go into `pkg/card_db/overrides/<set>.json`, which the tools read but never write. Card data is committed here; card images never are.

## Russian card names

Russian names come from the old foursouls.ru site's database migrations (goose `*.sql` files; a local copy sits in this repo's git-ignored `.cache/migrations/`). The tools read them; this repo gets only the result in `pkg/card_db/data/*.json`:

```sh
cd ../tools && go run ./scraper ru-names -migrations ../four_souls/.cache/migrations -game ../four_souls
```

For every card it writes `text.ru.name` and `translators.ru` (Alchba, Mitboy, d0m1n1c132, Blacknenko, N-Chance, the foursouls.ru translators). The site's `text_id` is `<set>-<card id>` with this repo's set suffix dropped (`fsp2-pills_2_fsp2` is `fsp2-pills_2`), then a fixed list of the site's misspelled ids (owner, 16.5). A two-sided card gets "Front/Back" like its English name, and every *Pills!* is *Пилюля!*. Reprint entries show the older card's name and are skipped. The command prints per set how many cards got a name and every card it could not match; rerunning it gives the same files.

`scraper build` keeps every language but English from the current data file, by card ID and version, so a new scrape keeps the Russian names. A translator's own fixes can go straight into `text.ru` the same way ([Translate cards](Translate-cards)).

## The images submodule

`pkg/card_db/images` is a git submodule pointing to the **private** repository `FourSoulsRepo/card_db`. It holds only images and a README: card images in one folder per language and in it one per set (`en/b2/`), and the game's own assets (fonts, backs, dice, mats) in `assets/`, outside the language folders. A translated card image goes to `<lang>/<set>/` under the English file name; the game shows it with that language chosen, else the English one (16.6).

* With access: `git clone --recurse-submodules …` or `git submodule update --init`.
* Without access (public clones, forks): the folder stays empty, everything still builds, and cards show as text.
* Release builds pack images with `-tags "embed cardimages"`; plain `-tags embed` packs only card data.
* The submodule embeds its own files (a Go package inside it, 8.29). Run `go generate ./pkg/card_db/images` before a `cardimages` build or test; `make internal` and CI do it.
* CI: `.github/actions/card-images` checks the images out at the pinned submodule commit, using the secret `CARD_DB_DEPLOY_KEY` (a read-only deploy key of `FourSoulsRepo/card_db`). Every app build runs twice, text-only and with images; without the secret (e.g. pull requests from forks) the image builds skip themselves.
* Its history is a single commit: after regenerating images, commit inside the submodule with `git commit --amend`, force-push it, then commit the new submodule pointer here.
