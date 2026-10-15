# 11. Who deals damage a player takes

* **Status:** Accepted
* **Date:** 2026-10-13
* **Authors:** @HardDie

---

## Context

1. Play test, record `2026-10-12-175853-g1.fsrec`, step 157.
   1. Boom Fly died: "it deals 1 damage to each player".
   2. The history said "Harddie deals 1 damage to Harddie".
   3. The other player's line was the same, with their name.
2. The engine writes "take n damage" as `DealDamage(n, You)`.
   1. `EachPlayer` runs it once per player.
   2. Each run makes that player the controller.
   3. So each damage item has controller = target = the player.
3. The rules handle that already.
   1. `whenYouDealDamage` ignores damage a player takes from itself.
   2. The item keeps `source` and `source_card`: Boom Fly.
4. The client named the controller as the dealer.
   1. `stackParts` used the controller whenever there was one.
   2. The source card only when there was no controller.
5. The same line came out for every "take damage" effect.
   1. Secret Room "Take 3 damage", curses, monster death triggers.

## Considered options

1. **Client: name the source card for self-damage**
   1. Won. Display only; no engine, protocol or record change.
   2. Old records show the right line too.
2. **Engine: no controller for monster damage**
   1. Lost. Changes `KilledBy`, `By` and priority order.
   2. Card triggers depend on them; risk for a display bug.
   3. Old records would still replay with the old items.
3. **Engine: a separate "dealer" field on damage items**
   1. Lost. Protocol and record change for one line of text.
   2. May come back if a card needs the true dealer.

## Decision

Use option 1.

1. A damage item whose target is its own controller is "taken damage".
2. Its line names the source card: "Boom Fly deals 1 damage to Ann".
3. With no source card: "Ann takes 1 damage".
4. Damage to anyone else still names the controller.
5. Code: `stackParts` in `frontend/src/stack.ts`.

## Consequences

### Positive

1. History strip and stack panel both read right.
2. No change to the engine, the protocol or the records.

### Negative and risks

1. The engine still says the player dealt the damage.
   1. A future "when you deal damage" card must keep ignoring self-damage.
2. The hover's owner line still shows the controller.

### Neutral

1. Monster combat damage already had no controller; unchanged.
