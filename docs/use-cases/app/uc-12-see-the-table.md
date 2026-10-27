# UC-12: See the table

**Module:** `frontend/src/screens/GameScreen.tsx`, `frontend/src/components`
**Status:** In progress
**Actors:** Player
**Goal:** The player sees the whole game state at a glance (V-01, V-02)
**Preconditions:** The game started (UC-03, UC-04)

## Main scenario (happy path)

1. The server sends an update after every step.
2. The app draws the table from the update's view.
   1. The player's own area is at the bottom: character, items, hand.
   2. The others sit in turn order, clockwise.
   3. Two players: the other one on top.
   4. Three players: left and right.
   5. Four players: left, top, right.
3. The middle shows the board, in one frame (F-01).
   1. Left: the turn, the stack and the history.
   2. Middle: monster slots above shop slots, then bonus souls.
   3. Right edge: monster, treasure and loot decks, each with its discard.
   4. The board keeps one shape; it scales with the window only (V-23).
   5. Nothing lies over the board (F-02).
   6. Left edge: a thin history strip, before the turn and stack (15.22).
4. Each player area shows HP, ATK, cents, souls and hand size.
5. The active player's name is in the accent color.
6. The player's own area reads at a glance (15.17).
   1. Left: name and coins as a big number.
   2. The character sits centered under the board; items to its right.
   3. Health and attack are big numbers on the character card.
   4. The hand fans out at the bottom edge, partly off screen.
   5. Hover a card to read it in full (V-18).
   6. The coins show the game's own coin (15.38).
7. Opponents read at a glance too (15.18).
   1. A small name and coins with the coin icon.
   2. The character with big health and attack; items beside it.
   3. The hand: one back per card, fanned, with the count.
   4. Top: mirrors the own area; the hand hangs past the top edge.
   5. Sides: one narrow column; with three players the board grows.

## Alternative scenarios and errors

* **2a. Used cards:** a character or item that is not charged is dimmed.
* **2b. Damage:** a red counter on a character shows damage taken.
* **2c. Dead player:** the area fades and says "dead".
* **3a. Empty deck or discard:** a dashed outline keeps its place.
* **3b. Covered slot:** a "+N" badge shows the cards under the top one.
* **3c. More than four slots:** the slots share the row; the board stays the same.
* **1a. Bans or a draft first:** "Setting up the match…" until 7.14.
* **1b. Events play first (V-08):** each update's events play one at a time.
  * Looted and discarded cards fly between deck, hand and discard.
  * Dice rolls show the result; cents, healing and damage pop up.
  * A new turn shows a banner; a revealed card shows large.
  * Then the table shows the update; new cards pop in.
* **1c. Many updates at once:** effects play faster; a long queue skips them.
  * The table always ends on the latest update.
* **1d. Animations off:** updates show at once.
* **1e. Dice look (V-19):** a roll shows the face of the player's dice texture.
  * No texture (public clone) or a value past 1–6: a plain die with the number.
  * Every screen draws a roll in the roller's look, not the viewer's (15.6).
  * A seat without a look, as a bot's, gets the default one.
  * A look this app does not have: a plain die.
* **1f. Card backs (V-26):** each deck shows its own back.
  * An opponent's hand is fanned loot backs with its card count.
  * No back images (public clone): plain dark backs with the deck name.
* **1g. No coin image (public clone):** a drawn ¢ coin by the cents.
* **2d. Counters (15.23):** a big green blob with the count, on the card art.
  * A text card shows it in its top corner, clear of the words.
  * Bonus souls show it too: Soul of Wrath's progress (19.29).
* **2e. Changed monster stats (15.23):** a badge per changed value.
  * Red health, a white die for the dice, gold attack.
  * A value as printed shows no badge; the card shows it.
  * A damaged monster shows its lower health instead of a damage mark.
* **2f. Kept die numbers (15.67):** a die on the card art per number.
  * Telepathy for Dummies: the roll it stored this turn.
  * Red Candle and Crystal Ball: the number chosen.
  * Drawn in the owner's dice look; a plain die without one.
  * The die goes when the number is used or the turn ends.
  * Its tooltip reads "Kept dice: 4" (15.83).
* **1h. Attack hit (15.3):** the attacker's card lunges into the target.
  * It winds up, hits, and goes back to its place.
  * A yellow starburst "−N" pops on the card that was hit.
  * The target's health drops at the hit, while the card goes back.
  * A monster's hit on a player lunges the monster at the character.
  * Damage shows no popup when it goes on the stack.
* **1i. Other damage (15.37):** a glowing shot flies from its source.
  * The source: its card on the table, else its stack row.
  * Every damaged card bursts "−N" at once.
* **1j. Hurt health (15.37):** the health number is red below full.
  * The heart turns pale, so the red number reads.
  * A hurt monster shows its health even when it matches the print.
  * A character has no damage mark in its corner; the health shows it.
* **1k. Death (15.37):** a big skull pops on the dying card.
  * The card leaves while the skull fades.
  * A dead player's skull sits on their character.
* **1l. Damage sounds (15.50):** an attack hit sounds as a punch.
  * Other damage sounds as a fireball: a cast, then an impact.
  * All shots of one update share one fireball.
  * No private images (public clone): no sound.
  * With animations off, no sound plays.
* **1m. Soul crystals (15.39):** souls show as gem slots left of the character.
  * One slot per soul the game needs to win; four today.
  * Filled slots are gained souls, from the bottom up.
  * A card worth two souls fills two slots.
  * Opponents get the same column, smaller.
  * Hover the column to read "Souls: N of 4".
  * Pure CSS: the same without the card images.
* **1n. A card put into a slot (19.3):** it flies there from its deck.
  * It stays on the slot while the update's later effects play.
  * A deck attack: the monster is in its slot before the attack roll.
  * Then the table shows the update and the copy goes away.
* **1o. A card revealed twice (19.4):** it shows large only the first time.
  * A deck attack shows its monster large when it is drawn.
  * In the next update it only flies into its slot.
* **1p. A revealed hand (19.5):** the opponent's cards show face up.
  * Guppy's Eye reveals each other player's hand.
  * Same fan and place as the backs, with the count.
  * At the top the cards hang past the edge, turned to their owner.
  * Hover a card to read it in full.
* **1q. A card that flew to its place (19.6):** it does not pop in there again.
  * A deck attack's monster lands in its slot once.
  * A looted or discarded card lands in the hand or pile once.
  * Cards no effect brought still pop in.
* **1r. A deck attack's drawn card (19.9):** it stays in view until placed.
  * It lies face up on the monster deck, for every player.
  * The slot dialog shows it above the slot buttons.
  * It leaves the deck when it flies into the chosen slot.
* **1s. A big hand (19.15):** the fan tightens; the end cards drop only a little.
  * Every card's top stays on screen, even with 12 cards.
  * Each card can be hovered and played.
* **1t. Soul cards behind the crystals (15.56):** gained souls are no cards on the table.
  * Each filled crystal holds its soul card.
  * A card worth two souls holds both of its crystals.
  * Hover a filled crystal to zoom its card; on touch, a long tap.
  * A killed monster zooms as the soul it gave; a bonus soul zooms too.
  * A character's own soul value zooms the character.
  * An effect that picks a soul highlights its crystal; click it to pick.
* **1u. Curses apart from items (15.55):** curses are not in the items row.
  * Own area and top: at the far left of the player area.
  * Not next to the soul crystals.
  * Side columns: their own row above the character, at its left.
  * No curse: the areas keep their size and look as before.
* **1v. A revealed card's path (15.73):** it grows out of its deck to the middle.
  * Moved in the same update, e.g. a slot refill: it flies from the middle.
  * It shrinks to a table card on the way to its place.
  * Not moved: it shrinks back onto its deck.
* **1w. A played loot card's path (15.74):** another player's card rises from their hand.
  * It shows large in the middle, then shrinks into the stack panel.
  * Resolved or cancelled, it flies from the stack to the loot discard.
  * That flight shows for the viewer's own loot cards too.
  * Once it flies, the card is gone from the stack, row and all (19.41).
* **1x. Stolen cents (15.75):** "−N¢" rises on the victim's coins.
  * A coin flies from the victim's coins to the thief's.
  * "+N¢ stolen" rises on the thief's coins when it lands.
* **1y. A changed roll (15.76):** the die shows its old value first.
  * A shot flies to it from the card that changed it.
  * The die turns to the new value: "Ann's roll: 3 → 4".
* **1z. Damage from a card that left (15.77):** the shot comes from its player.
  * From their character card, else their area.
  * Only damage with no known dealer shows the hit alone.
* **3d. The History panel (15.88):** one line per play, under turn lines.
  * Each line names who: "Bo attacks Gaper", "You kill Gaper".
  * Buys: "Bo buys Spoon Bender"; a failed one: "Bo can't buy: the price is 10¢".
  * Declared attacks and purchases, gains, souls, thefts, cents.
  * Deaths, destroyed items, fizzled and cancelled effects.
  * Hovering a line zooms its card.
  * Others' hidden cards stay hidden: "Bo loots a card".
* **3e. The history strip (15.89):** one tile per headline, newest on top.
  * Played cards, abilities and rolls: as on the stack, with a sword or bolt.
  * A buy: the bought card with a coin badge.
  * A kill: the dead card with a skull badge, "Bo kills Gaper".
  * Hover a tile: its card large, with the line under it.
* **3f. What a play did (15.90):** hovering a strip tile shows an arrow and the results.
  * A played card: its target, then what it gave at once, such as cents.
  * Damage: the dealer large; the target with a starburst and the number.
  * A kill: the killer large; the dead card under a skull, then the rewards.
  * A buy: the buyer large; the bought item.
* **3g. New tiles (15.51):** a new strip tile fades in on top.
  * The older tiles slide down; none jumps.
  * With animations off, tiles change at once.
* **3h. Why a stat changed (15.92, V-41):** a changed stat says why.
  * A stat above its print is green, below it red.
  * Characters and monsters alike, for every player.
  * Hover the card: its zoom sums each stat below it.
  * Example: "ATK 3 = 1 printed + 1 Brimstone + 1 XI. Strength (this turn)".
  * Roll changes follow: "Attack rolls: +1 Meat!".
  * Effects of this turn are small icons on the card's right edge.
  * Each icon is cut from its source card's art.
  * Hover an icon: its card zooms, with what it does.
  * Lasting items get no icon; they show only in the sums.
  * A text card's icon shows its initial; a hidden source a "?".
* **3i. Effects that end (15.93):** an effect's icon leaves when it ends.
  * A boost for the next attack roll has a small die on its icon.
  * The roll that uses it makes the icon pop and fade.
  * At the end of the turn, this turn's icons shrink away.
  * With animations off, icons go at once.
* **3j. First attack roll ATK (19.67):** it shows while that roll is to come.
  * Curved Horn, Champion Belt, Polydactyly: "+1 (first attack roll)".
  * It shows from the turn start, before the attack.
  * It stays while the first roll's damage is on the stack.
  * It goes once that damage is dealt, or the roll missed.
* **1aa. An uncovered card (19.72):** it shows up in its slot.
  * A covering monster dies; the card below was there already.
  * It does not show large or fly from the deck.
  * It pops in place, like any new table card.
  * Rooms too: a Red Key room leaves; the room below shows.
* **1ab. A dead monster's flight (19.73):** it flies from its slot.
  * The skull pops on it first.
  * No soul: it flies to the monster discard.
  * A soul: it flies to the player who gains it.
  * Reward rolls may come first; it still flies from its slot.
  * Its slot shows empty while it flies.
  * A bonus soul does not fly.
* **1ac. Loot draw sound (18.7):** a looted card makes a deal sound.
  * One sound per card, as each card flies to the hand.
  * Any cause: start of turn, a monster's reward, a card's effect.
  * A card taken from the loot discard sounds too.
  * No private images (public clone): no sound.
  * With animations off, no sound plays.
* **1ad. Player death sound (18.8):** a dead player makes a death sound.
  * It plays as the skull pops on the character.
  * Only once the death resolves: nothing can undo it then.
  * A death prevented on the stack is silent.
  * A monster's death is silent.
  * No private images (public clone): no sound.
  * With animations off, no sound plays.

## Postconditions

* The table fits 1280×720 and larger windows.
* Planned: the stack panel (7.8), history and turn panels (7.9).
* No screen scrolls as a whole, at any window size (V-10).
  * Long lists and panels scroll on their own.
* A mouse drag selects nothing on the table or in menus.
