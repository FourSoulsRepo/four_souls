# 4. Card data format

* **Status:** Accepted
* **Date:** 2026-10-09
* **Authors:** @HardDie

---

## Context

1. Clients ship the full card collection; servers send card refs (CD-01).
2. Changed cards keep their ID and get a version (CD-08).
3. Text is English first, other languages later (LO-01, LO-02).
4. Artists and translators are credited by name only (CD-04, CD-06).
5. Images are licensed and never committed here (A-11).
6. The engine needs IDs and effects, not display data (A-06, A-11).
7. Sources: foursouls.com pages and the official TTS table (2.2, 2.3).

## Considered options

1. **One JSON file per set**
   1. Won. Small diffs, easy review, stdlib only, works in wasm.
2. **YAML or a database file**
   1. Lost. Extra dependency or a binary file in git.
3. **Set code inside the ID (`b2_the_d6`)**
   1. Lost. A reprint in another set would look like a new card.
4. **Site slug without the set prefix as the ID (`the_d6`)**
   1. Won. Stable, readable, unique within a set.
5. **Stats and rewards as numbers**
   1. Lost. Cards print values like "X"; this is display data.

## Decision

Use options 1 and 4; stats and rewards stay strings (option 5 lost).

1. Files
   1. `pkg/card_db/data/<set code>.json`; the file name must match the code.
   2. Images under `pkg/card_db/images/<set code>/`, git-ignored.
2. Set file
   1. `set`: `code` (e.g. `b2`) and `name` (e.g. `Base Game V2`).
   2. `cards`: list of cards.
3. Card fields
   1. `id`: lowercase words joined by `_`, e.g. `the_d6`.
   2. `version`: omitted for the first version, else 2, 3, …
   3. `type`: the printed card type, e.g. `Eternal Treasure Card`.
   4. `copies`: copies in the set (from the TTS table).
   5. `text.<lang>`: `name`, `effects`, `footnotes`, `notes`.
   6. Icons in text: `{Icon Name}`, e.g. `{Tap Effect}`.
   7. `stats`, `rewards`: printed name and value strings.
   8. `artists`: role and name; `translators.<lang>`: names.
   9. `image`: path inside `images/`, e.g. `b2/the_d6.webp`.
   10. `related`: card refs, e.g. a character's starting item.
   11. `rebalanced`: true for cards changed by the Rebalance pack.
4. Card ref
   1. `id` for the first version, `id@N` for version N ≥ 2.
   2. The ref is unique across all loaded sets.
5. Loader (`carddb.Load`)
   1. Reads every `data/*.json` from an `fs.FS`.
   2. Rejects bad IDs, duplicate refs, missing English names.
   3. Rejects image paths that leave `images/`.
   4. Reports every problem at once.
6. Build modes
   1. `-tags embed`: `carddb.Embedded()` packs `data/` and `images/`.
   2. Otherwise `assets.CardsFS()` reads `$FOUR_SOULS_CARDS`.
   3. Then `cards/` next to the binary, then `pkg/card_db/` (dev).

## Consequences

### Positive

1. Card files are plain JSON, reviewable in pull requests.
2. A fresh clone loads card data without the private images.
3. The module builds for `js/wasm` and has no dependencies.

### Negative and risks

1. A reprint with the same slug in two sets is a load error.
   1. It must become a version or alt art when that set is added.
2. Embedded release builds grow by the size of all images.

### Neutral

1. Game rules for each card live in the engine (step 4), keyed by ref.
2. The data is produced by the separate `FourSoulsRepo/tools` repo.
   1. Its `scraper build` writes `data/<set>.json` here.
   2. Its `scraper images` writes `images/<set>/` here.
   3. It mirrors these types in `scraper/carddata`; keep both in sync.
