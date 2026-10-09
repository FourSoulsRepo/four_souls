# Step 2. Card data

Goal: all Base Game card data and images, collected by scripts.
The scraping tools work for later sets too.

Ideas: CD-01, CD-04, CD-06, CD-07, CD-08, A-11.

---

### [x] 2.1 Browser session tool

1. Goal: fetch pages behind Cloudflare (CD-07).
2. Tasks:
   1. Go tool in the separate repo `FourSoulsRepo/tools`.
   2. Opens a visible Chrome with a persistent profile folder.
   3. The user passes the Cloudflare check once.
   4. The tool reuses that session to fetch pages.
   5. Polite rate limit; retries with backoff.
   6. Raw pages cached in a git-ignored folder.
3. Done when:
   1. foursouls.com pages download after one manual check.
   2. A second run reads from the cache.
4. Docs: the tools repo README; wiki page "Card data tools" here.

### [x] 2.2 Official site card scraper

1. Goal: card data from foursouls.com.
2. Tasks:
   1. Find the card list and per-card pages.
   2. Extract: name, set, type, text, stats, artist, image URL.
   3. Stats: HP, dice, attack, souls, rewards, as present.
   4. Extract per-card FAQ or rulings if the site has them.
   5. Output raw JSON per set.
3. Done when:
   1. Base Game cards are all extracted.
   2. A report lists cards with missing fields.
4. Result:
   1. Base Game means **Base Game V2** (set code `b2`), the current edition.
   2. 287 unique cards (the site's 340 counts copies); 0 with missing fields.
   3. The site has no per-card FAQ or rulings; step 3 uses other sources.
   4. 7 cards carry the *Rebalanced* flag; relevant for card versions (CD-08).
   5. Open: copies per card (e.g. how many A Penny) are not on the page.

### [x] 2.3 Official TTS table import

1. Goal: names and images from the official TTS table.
2. Tasks:
   1. Download the save (Steam Workshop 2501791757).
   2. Decode BSON; list cards and image sheets.
   3. Download sheets; cut them into single cards.
   4. Match TTS cards to site cards by name.
3. Done when:
   1. A match report shows matched and unmatched cards.
4. Result:
   1. Bag "Base Game": 282 names, 340 cards with copies (= site count).
   2. All names match the site; 0 unmatched on either side.
   3. Copies per card come from here (answers the open point in 2.2).
   4. 4 names cover several site cards: Chest, Dark Chest, Gold Chest, Pills!.
   5. 339 faces cut at 659×921 px PNG (294 MB in the cache).

### [x] 2.4 `card_db` schema and loader

1. Goal: the card display data module (A-11).
2. Tasks:
   1. Card fields: ID, version, set, type, name, text, stats.
   2. Also: artist, translators, image file name.
   3. Text is keyed by language; only `en` for now (LO-01).
   4. Card reference format: `id` or `id@2` (CD-08).
   5. JSON files, one per set.
   6. Loader works with both asset modes (1.7).
   7. Validation: unique refs, required fields.
3. Done when:
   1. Loader tests pass on a sample set.
   2. ADR: card data format.

### [x] 2.5 Build Base Game data

1. Goal: final Base Game JSON in `card_db`.
2. Tasks:
   1. Merge site data and TTS data.
   2. Manual fixes in a separate overrides file.
   3. Scripts never overwrite overrides.
3. Done when:
   1. Every Base Game card validates.
   2. **Owner** reviews a sample of cards.
4. Result:
   1. The tools repo writes `pkg/card_db/data/b2.json` (`scraper build -game`).
   2. 287 cards, 340 with copies; loads with `carddb.Load`.
   3. Overrides: `pkg/card_db/overrides/b2.json` (empty so far).
   4. Copies of shared names (Pills!, chests) split evenly; override if wrong.
   5. Images will come from the site (962×1312), not TTS (659×921).

### [x] 2.6 Image pipeline

1. Goal: small, readable WebP images (A-11).
2. Tasks:
   1. Convert to WebP at one fixed size.
   2. Pick the size by checking text readability when zoomed.
   3. One folder per set.
   4. For now images live in `pkg/card_db/images/`.
   5. That folder is git-ignored; never committed here.
   6. Private image repo comes later (A-11).
3. Done when:
   1. Base Game images total is reported (MB).
   2. Card text is readable on zoom.
4. Result:
   1. Source: site images, 962×1312 PNG, downloaded through Chrome.
   2. Output: 600×818 WebP; Catmull-Rom resize, then `cwebp -q 85 -m 6 -pass 10`.
   3. Base Game: 287 images, 14.4 MB total, 52 KB average.
   4. Dense text (e.g. Pandora's Box roll table) stays readable.

### [x] 2.7 Images in builds

1. Goal: builds find images; work without them.
2. Tasks:
   1. Local: external mode reads `pkg/card_db/images/`.
   2. Without images: build passes; cards fall back to text.
   3. CI builds without images for now.
3. Done when:
   1. A fresh clone builds and runs with text-only cards.
4. Result:
   1. Fresh clone: all tests pass, 287 cards load, 0 images; embed build OK.
   2. Local: 287 of 287 images found in both build modes.
   3. `carddb.HasImage` tells the client when to draw a blank card.
   4. The app does not link card data yet; the client uses it in step 7.
5. Later (with the private repo):
   1. CI clones the image repo with `--depth 1` via a secret.
   2. Script publishes the image repo as one amended commit.
   3. **Owner:** create the repo and add the secret.
