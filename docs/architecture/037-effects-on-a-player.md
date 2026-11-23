# 37. Effects on a player

* **Status:** Accepted
* **Date:** 2026-11-19
* **Authors:** @HardDie

---

## Context

1. Owner report, 2026-11-19 (19.232, moved to 18.29):
   1. Trick or Treat: Ghostly Trickery's 1–2 adds "+1 damage taken".
   2. Each attack start rolls again; the +1s add up.
   3. Nothing on the table shows the total.
2. The engine keeps effects on players as stat parts.
   1. Statics of cards in play, and boosts till end of turn (`Boost`).
   2. `Parts(stat, player, subject)` lists them, each with its source.
   3. A static's `Applies` already picks the players it is for.
3. The view sends only some of them (V-41, 15.54).
   1. `playerSources`: HP, ATK, rolls, attack rolls, attacked DC.
   2. About 30 other player stats are never sent.
   3. E.g. damage taken, attack tax, max hand size, locks, no souls.
4. The client shows sent sources two ways (V-41, SB-03).
   1. Badges: HP and ATK numbers on the character card.
   2. Icons: one per source card, only for this turn's effects.
   3. Owner, 2026-10-22: lasting sources only in the badge's hover sum.
5. Effects on a player come from three places:
   1. A. A card in the player's own area (Daddy Haunt, Steamy Sale).
   2. B. A card elsewhere: rooms, monsters, another player's item, a rule card.
   3. C. A card already gone: till end of turn (Get Out of Jail, Algiz).
6. Today A shows its card; B and C show nothing for most stats.
   1. B rooms with a condition are worst: "fewest souls" players.
   2. Only the engine knows whom such a room applies to now.
7. Shields are no stat parts (owner question, 2026-11-24).
   1. `Game.Shields`: "prevent the next damage" till end of turn.
   2. Soul Heart, Holy Chest, The Poop: their card is gone once played.

## Considered options

1. **One list of effects per player; three levels on screen**
   1. Won. One engine path for every stat.
   2. Reuses the badges and icons players already know.
2. **A marker on each source card (e.g. the rule card)**
   1. Lost. A shared card can't say which player it hits.
   2. A card already gone (C) has nothing to mark.
3. **A badge on the character for every stat**
   1. Lost. A character could carry ten badges.
   2. Most effects are not numbers: locks, no souls.
4. **Only a text list in the zoom**
   1. Lost. Nothing on the table says that an effect is on.

## Decision

Use option 1.

1. Engine: every player effect goes into `PlayerView.Sources`.
   1. Each stat a player can have gets a key in `statKeys`.
   2. `StatUsed` stays out: it is engine bookkeeping.
   3. Monster stats stay as they are (`MonsterStats.Sources`).
   4. A flag effect (a lock, no souls) has amount 1.
   5. A conditional static is sent only to the players it applies to now.
   6. Every boost names its source; a hidden one sends no card.
   7. Effects on one object (Algiz's eternal items) are not in this ADR.
   8. Shields on a player are effects too (owner, 2026-11-24).
      1. Sent in `Sources`: stat `shield`, or `shield_cap` for a cap.
      2. Amount is the shield's amount or cap; lasts `turn`.
      3. Shields on monsters or items stay out.
      4. A shield of the whole next damage has no amount (18.30).
   9. Game-wide effects with no player stay out (18.30).
      1. Rewards doubled: A Lucky Penny!, Blessing of the Sack.
      2. A roll of 2 or 3 becomes 1: Eternal D6.
      3. Two dice, the left player picks: Cursed Lil Haunt.
      4. The owner decides later where they show.
   10. The DC to attack a player sends only the effect that counts (18.30).
       1. The newest wins (R-MECH-06), so no sum.
   11. A flag stat's amount is 1, even when the engine stores more.
       1. E.g. the death-penalty item's receiver (Lust for Blood).
2. Level 1, badges: numbers that change a fight.
   1. HP and ATK, as now.
   2. New: damage taken, red, "+N", only while above 0.
   3. No other stat gets a badge without a new decision.
3. Level 2, icons: one per source card, in groups B and C.
   1. Lasting group A effects have no icon: the card lies next to the character.
   2. This changes the owner's rule of 2026-10-22 (V-41.4) for B.
   3. Each icon's hover lists its effects, one line each.
   4. One source's lines of one stat add up: "+3", not three "+1".
   5. An icon whose lines are all one number stat shows that sum.
   6. Above 4 icons: the first 4, then a "+K" chip.
   7. The chip opens the character's zoom (level 3).
   8. A this-turn effect keeps its icon, its card's place aside (owner, 2026-11-24).
   9. The chip opens the zoom on hover, as an icon does.
4. Level 3, the zoomed character: "Effects on you".
   1. Every effect, groups A, B and C.
   2. "This turn" first, then "while in play".
   3. Each line: what it does, then the source card's name.
   4. No repeats of a sum line: HP, ATK, rolls, damage taken (owner, 2026-11-24).
5. Words are keys in `en.json`, one per stat (ADR 021).
6. Every player sees every player's effects: they are public.

## Consequences

### Positive

1. The damage-taken total is on the character card.
2. A room's condition shows on the players it hits.
3. A till-end-of-turn effect stays seen after its card is gone.
4. A new stat needs one key and one text; no new screen code.

### Negative and risks

1. The view grows by a few sources per player.
   1. Sources are sent only while an effect is on.
2. A busy character may hide icons behind "+K".
3. Every locale needs the new effect texts (`make i18n`).

### Neutral

1. Monster and item stats keep their own badges.
2. Old records replay; a missing stat key shows no line.
