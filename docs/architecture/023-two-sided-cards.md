# 23. Two-sided cards

* **Status:** Accepted
* **Date:** 2026-10-27
* **Authors:** @HardDie

---

## Context

1. Requiem has three two-sided cards (R-CARD-39).
   1. Anima Sola / The Revenant, The Enigma / amginE ehT.
   2. The Harbingers / The Beast!
2. The engine flips them already (R-MECH-33).
   1. Each back is its own `CardDef`: `the_revenant`, `amgine_eht`, `the_beast`.
   2. `Object.Flipped` marks the back up.
3. card_db had one entry per card.
   1. Name "Front/Back"; both sides' effects in one list.
   2. Stats "2/1": front and back in one value.
   3. Only the front image.
4. The view sent the front's ref, flipped or not.
   1. The client never showed the back.
5. Owner, 2026-10-27:
   1. Each side shows its own name, never "Front/Back".
   2. A card that changes side plays a flip animation.
6. The site has a page per card with both images.
   1. The back is the second lightbox link, "Alternate Card Back".

## Considered options

1. **A `back` face inside the card entry**
   1. Won. One entry per printed card; copies and artists stay shared.
2. **Two card entries, linked by `related`**
   1. Lost. Copies would count twice; a back is never dealt.
3. **Split "Front/Back" names in the client**
   1. Lost. Effects and stats need the split too; data is the place.
4. **View sends the front ref and a flag; client finds the back**
   1. Lost. Every client lookup by `card` would need the flag.
5. **View sends the face-up ref and a flag**
   1. Won. Every existing lookup draws the right side.

## Decision

Use options 1 and 5.

1. card_db
   1. A two-sided card has `back`: its other face.
   2. `back.id`: the back's ref, e.g. `the_revenant`.
   3. It is unique among all refs, and equals the engine's back ref.
   4. `back.text.<lang>`, `back.stats`, `back.rewards`, `back.images`.
   5. Type, copies, artists and translators: the front's.
   6. Each face has only its own name and effects.
   7. The loader lists the back as a card of its own.
   8. `Card.Front` on the back names the front; `Back` the other way.
2. Tools (`FourSoulsRepo/tools`)
   1. `cards` reads the back image URL.
   2. `build` splits the name at `/`.
   3. An effect starting "Back name: " starts the back's effects.
   4. A "Front name: " prefix is dropped.
   5. A stat "a/b" gives a to the front and b to the back.
   6. Anything else the site gives once goes to the front.
   7. An override's `back` key replaces the whole back.
   8. `images` writes the back to `<set>/<back id>.webp`.
   9. `ru-names` splits the Russian name the same way.
3. Engine view
   1. `card` is the face up: the back's ref when flipped.
   2. `flipped: true` marks the back up.
4. Client
   1. `CardInfo.back` and `front` link the two faces.
   2. The table draws the face up.
   3. The zoom shows both faces.
   4. A card whose `card` turns to its other face plays a flip.
   5. The flip is `anim: card-flip` (ADR 012).

## Consequences

### Positive

1. Names, text, stats and images are right per side.
2. The history names the side a card flipped to.
3. Every client lookup by `card` draws the side that is up.

### Negative and risks

1. A back ref must match in card_db and the engine.
   1. `TestBackFacesMatchCardData` (`internal/server`) checks both.
2. A site page that names sides differently needs an override.
3. An older client lacks the back's data: it shows the back's ref.

### Neutral

1. A back is never dealt or banned; only its front is.
2. Copies of a card share one back image.
