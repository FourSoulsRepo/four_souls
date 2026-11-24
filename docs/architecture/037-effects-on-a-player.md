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
8. Owner review of 18.29–18.32 (2026-11-24): show only what the table hides.
   1. A card on the table already shows its effect: no icon, no list line.
   2. E.g. Curse of Tiny Hands lies by the player; no "max hand size 2" icon.
   3. A room's number change shows in the badge and its zoom sum line.
   4. Hidden: a gone card's effect, a shield, a used item's this-turn effect.

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
      4. Hidden ones show as item 3.18; cards on the table show nothing more.
   10. The DC to attack a player sends only the effect that counts (18.30).
       1. The newest wins (R-MECH-06), so no sum.
   11. A flag stat's amount is 1, even when the engine stores more.
       1. E.g. the death-penalty item's receiver (Lust for Blood).
2. Level 1, badges: numbers that change a fight.
   1. HP and ATK, as now.
   2. New: damage taken, red, "+N", only while above 0.
   3. No other stat gets a badge without a new decision.
3. Level 2, icons: one per source card, hidden effects only (owner, 2026-11-24).
   1. A lasting effect has no icon: its card lies on the table.
   2. The owner's rule of 2026-10-22 (V-41.4) holds again.
   3. Each icon's hover lists its effects, one line each.
   4. One source's lines of one stat add up: "+3", not three "+1".
   5. An icon whose lines are all one number stat shows that sum.
   6. Above 4 icons: the first 4, then a "+K" chip.
   7. The chip opens the character's zoom (level 3).
   8. A this-turn effect keeps its icon, its card's place aside (owner, 2026-11-24).
   9. The chip opens the zoom on hover, as an icon does.
   10. Where a lasting effect's card lies no longer matters (owner, 2026-11-24).
       1. 18.31 read the own area from `in_play`; item 3.1 replaces that.
   11. Monsters keep icons for this turn's effects only (18.31).
   12. The number shows only for a change of a number (18.31).
       1. None for flags, caps, the DC to be attacked, shields.
       2. None for any number of attacks or a free shop item.
   13. A lasting "on your turn" effect says so (18.31).
       1. The server sends it on other turns too.
   14. Three icons per column; a second column to their left (18.31).
   15. All of a player's shields share one icon (owner, 2026-11-24).
       1. Each shield blocks one damage instance, oldest first (R-MECH-46).
       2. So the number is how many shields, not a sum of damage.
       3. Its hover lists each, in the order they will block.
   16. The "+K" chip is a pill, "⋯ K", with no card art (owner, 2026-11-24).
       1. So it never looks like an icon's "+N".
   17. Hidden effects that are not stats get icons too (owner, 2026-11-24).
       1. Skips their next turn: Famine, Haunted!, Overflow!.
       2. Takes an extra turn: XIX. The Sun.
       3. The next loot is copied: Blank Card.
       4. Lingering on the player: Black Heart.
   18. Game-wide hidden effects: an icon on the user's character (owner, 2026-11-24).
       1. The user: the controller of the card or ability that made it.
       2. Its line says it is for every player.
       3. A Lucky Penny!, Eternal D6, Compost, Lemon Mishap, Plan C.
   19. Built in 18.34.
       1. The shields' icon shows its count, even 1.
       2. So item 3.12.1 no longer covers shields.
       3. Its lines name each shield's card, in blocking order.
       4. Its art and zoom: the first shield with a known card.
       5. Item 3.13's "on your turn" went with the lasting icons.
4. Level 3, the zoomed character: "Effects on you".
   1. The same effects as the icons: hidden ones only (owner, 2026-11-24).
   2. "This turn" first, then "while in play".
   3. Each line: what it does, then the source card's name.
   4. No repeats of a sum line: HP, ATK, rolls, damage taken (owner, 2026-11-24).
   5. The list stands beside the card, above its keyword notes (18.32).
      1. Under the sums it left a 720 px window.
   6. Above 8 lines: the first 8, then a count of the rest (18.32).
   7. A hidden source is "a hidden card" (18.32).
   8. Only "This turn" is left (18.34).
      1. Item 4.2's "while in play" part is gone: lasting effects have no line.
      2. With nothing hidden, the zoom shows no list.
5. Words are keys in `en.json`, one per stat (ADR 021).
6. Every player sees every player's effects: they are public.
7. "Must attack" shows on the attack button (owner, 2026-11-24).
   1. The button turns red while an attack is owed.
   2. Its tooltip names the target: a monster, or the monster deck.
   3. The red border on a monster stays.
   4. The tooltip repeats the banner's line, from one source (`mustAttackLine`).
      1. Targets also: a player (Lust for Blood), any target (Envy).
      2. A random minion (Pilfering Krampus): the view's `random_attack`.
   5. Red, lit or dim; only on the owing player's own turn.
8. Red Candle shows its chosen number as a die on its card (owner, 2026-11-24).
   1. Till end of turn; not on any character.
9. A free shop item already shows as price 0 (Credit Card); no icon.

## Consequences

### Positive

1. The damage-taken total is on the character card.
2. A room's number change shows in the badge sum on the players it hits.
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
