# 10. Reprints, copy images and set order

* **Status:** Accepted
* **Date:** 2026-10-13
* **Authors:** @HardDie

---

## Context

1. Step 14 adds Gold Box V2 (`g2`) and Four Souls+ V2 (`fsp2`).
2. Their site slugs clash with Base Game IDs in 15 places.
3. Some clashes are reprints: same name, type and text.
   1. Coins, Bomb!, A Nickel!, Butter Bean!, Lil Battery, Troll Bombs.
   2. Also A Penny!, Soul Heart and Dice Shard.
4. Others are different cards with the same slug.
   1. Pills! variants, Isaac! (epic boss), Bum-Bo (character).
   2. Daddy Long Legs (treasure).
5. ADR 004 made any clash a load error (its risk 1).
6. Copies of one card often have their own art.
   1. 15 copies of 3 Cents!: each with a different coin.
   2. The site lists one copy; the others have numbered pages.
7. **Owner:** a game with Base Game and Gold Box uses every copy.
   1. All 12 Base Game and 2 Gold Box "2 Cents!" go into the deck.
   2. Each copy shows its own image.
8. **Owner:** sets are ordered by release, oldest first.

## Considered options

1. **Reprint entry with the older card's ID**
   1. Won. One card, one engine definition, copies per set.
2. **Set code as a suffix on every clashing ID**
   1. Lost. Reprints would look like new cards (as in ADR 004).
3. **Hand-picked readable suffixes (`isaac_boss`)**
   1. Lost. Needs a rule per card; the tools cannot pick them.
4. **Set code as a suffix only for different cards**
   1. Won. Mechanical and stable; old IDs never change.
5. **One image per card**
   1. Lost. Copies lose their own art.

## Decision

Use options 1 and 4; images become a list.

1. Set order
   1. `set.order`: 1 for the oldest set; unique.
   2. `set.released`: first release as `YYYY-MM`.
   3. Order follows the first release, not the V2 reprint.
   4. `b2` 1 (2018-12), `g2` 2 (2018-12), `fsp2` 3 (2019-06).
   5. Base Game comes before Gold Box, its add-on.
   6. The loader reads sets in this order.
2. Images
   1. `images` replaces `image`.
   2. One path: every copy shows it.
   3. One path per copy: each copy has its own art.
   4. Other counts are a load error.
   5. Copy files: `<id>.webp`, `<id>-2.webp`, `<id>-3.webp`, …
   6. File names come from the site slug, inside the set folder.
3. Reprints
   1. Same type, name, stats and effects as a card of an older set.
   2. Case and a final period do not count ("Gain 5¢" = "Gain 5¢.").
   3. Entry: `{"id", "reprint": true, "copies", "images"}`; nothing else.
   4. The ID must exist in an older set.
   5. The loader adds it to the card's `Prints`.
4. Different cards with a taken ID
   1. ID: slug ID plus `_<set code>`, e.g. `pills_g2`, `isaac_fsp2`.
   2. `related` refs follow the new ID.
5. Deck building
   1. A game's deck takes the copies of every chosen set.
   2. `Card.Prints` lists them: own set first, then reprints.
   3. Each copy will keep its image index for display (14.5).
   4. The engine counts them: `CardSet.Reprints` (14.3); copy images are 14.5.
6. Sources
   1. Copies per card: the official TTS table.
   2. Copy images: the numbered copy pages on foursouls.com.
   3. The tools' `scraper copies` command finds them.

## Consequences

### Positive

1. Every copy keeps its own art.
2. Reprints need no new engine code.
3. Base Game IDs did not change.

### Negative and risks

1. The images submodule grows: one file per copy.
2. A reprint with fixed wording keeps the older set's text.
3. A future set may reprint a card with new text.
   1. It is then a different card and gets the set suffix.

### Neutral

1. ADR 004 still holds for everything else.
2. The tools mirror the format in `scraper/carddata`.
