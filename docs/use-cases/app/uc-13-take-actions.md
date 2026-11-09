# UC-13: Take actions

**Module:** `frontend/src/actions.ts`, `frontend/src/components`
**Status:** Implemented
**Actors:** Player
**Goal:** The player plays a whole turn with mouse or touch
**Preconditions:** The game runs (UC-12); the player must act

## Main scenario (happy path)

1. The table highlights what the player may do now.
   1. Yellow: cards to play or use; green in the hand (15.27).
   2. The stack panel: Pass.
   3. Icons by own name: Buy and Attack; lit when usable, dim when not.
   4. A big End turn button at the board's right edge.
2. The player plays a loot card on the board: a drag, or a click, then a click there.
3. Others may respond; the player passes or acts again.
4. The player presses Attack.
5. The monsters that can be attacked glow blue; the player clicks one.
6. Rolls and responses play out.
7. The player presses End turn; the next player's turn starts.

## Alternative scenarios and errors

* **2a. Click or tap:** a table card opens its action menu: Use, Cancel.
  * A hand card has no menu; a click picks it up (2j).
* **2b. Abilities:** a character or item in play shows "Use" or "Use ability N".
* **5a. A choice dialog:** every other choice shows as a dialog (5i).
  * Card options show as cards (V-07); text options as buttons.
  * A question shows its text as the title.
* **5b. Pick a player:** a "Choose <name>" button in that player's area.
* **3a. Discard to hand size:** pick exactly N hand cards, then Discard.
* **1a. Not the player's move:** the bar says "Waiting for <name>".
* **1b. Updates still playing:** buttons stay but are disabled.
* **2c. Refused:** the bar shows the engine's reason and rule, e.g. "(R-TURN-06)".
* **4a. Keys:** B buys, A attacks; each icon's tooltip names its key.
* **4b. Pressed:** a pressed icon shows at once, before the server answers.
* **4c. Choosing:** while the player picks what to buy or attack, that icon stays pressed.
* **7a. Their turn:** the big button is grey and says "Their turn".
* **7b. Not yet:** on own turn, End turn stays dim until the engine allows it.
* **7c. Nothing left:** End turn glows green when only passing or ending is left.
* **7d. Keys:** E ends the turn; Space passes.
* **2d. Several abilities (15.24):** the card opens a dialog with the card large and one line per ability.
  * The player clicks a line to use that ability.
  * Lines not usable now stay greyed; Cancel or a click outside closes it.
* **4d. What is left (15.24):** on own turn the Buy and Attack icons carry a count badge.
  * Loot plays left show by the player's name as "Loot N".
* **2e. Aim (15.25):** a card that needs a target is aimed by two clicks.
  * The first click starts aiming; the card stays lifted in the hand.
  * The valid targets glow red: monsters, players, items, stack items.
  * With a mouse, a red arrow points from the card to the pointer.
  * The second click, on a target, plays the card at it; no dialog asks.
  * A click on empty space, a right click, or Esc cancels; nothing is sent.
  * Such a card opens no menu; it may be dragged too (2n).
  * The arrow is a chain of red dashes with a big head (15.32, 20.23).
  * Over a valid target a red ring marks the head.
  * A click on the aiming card again does nothing (15.33).
  * A click on a card, pile, or stack item that is no target keeps aiming.
  * A short note above the click says what the card may target.
  * E.g. "Pick a monster or a player."; it fades after 2.5 s.
  * A player is aimed at by the character card, as a monster is.
  * The character glows red; the player's area stays plain.
  * A card that may target itself, as The D20, picks itself on that click.
  * The arrow bows up the screen; its dashes march toward the head (20.23).
  * An arrow straight up or down stays straight.
  * It shows at once where the mouse is; no move is needed.
  * The system pointer hides while it shows; the head is the pointer.
  * The valid target under the head glows brighter.
  * No card zooms while the arrow moves over it.
  * Resting still on a card for 0.8 s zooms it, beside the card.
  * Any move closes that zoom.
  * Another player's arrow glides out from its source, with no ring.
  * It stays 0.25 s after they stop aiming.
  * Animation speed Off or reduced motion: no march, no glide.
  * Touch: no arrow, as before.
* **2f. Place, then aim (15.34):** a card with no target is dragged onto the board; no arrow.
  * A card that asks its mode first is dragged too; the mode dialog asks.
  * Then the target is aimed from the card, as in 2e; no dialog asks.
  * The card stays lifted in the hand until the target is picked.
  * Esc, a right click, or empty space cancels the whole play.
  * The card stays in hand; nothing goes on the stack.
  * A card with no valid target cannot be played; the engine decides.
  * A target that is not on the table keeps the dialog.
* **5c. Order of "instead" effects (19.17):** two cards replace one gain, e.g. Bum-Bo! and Counterfeit Penny.
  * The title names the gain: "Gaining 10¢: which card applies first?".
  * Each option is a card, with what picking it first gives below.
  * E.g. "→ Bum-Bo! levels up by 11".
  * A card whose effect then does nothing says so.
  * With three or more cards: "Then you pick the next card".
* **5d. Order of triggers (19.19):** two of a player's abilities trigger at once.
  * The title: "Choose which goes on the stack first; it resolves last".
  * Each option is the ability's card, with its name and ability text.
  * No preview: other effects may come before a trigger resolves.
* **6a. Throw (15.61):** every roll waits for its roller's throw.
  * The roller sees a big glowing Throw button where the die lands.
  * The button holds the player's die; the roll's name shows below.
  * The player clicks it, or presses Space.
  * The die lands with the engine's value; then the others may respond.
  * A reroll asks the roll's own player to throw again.
* **6b. The die lands (15.62):** the die falls from above and bounces twice.
  * It tumbles as a 3D cube and stops on the engine's value (15.69).
  * The roll's name shows below it once it lands.
  * The animation speed setting makes it faster or slower; Off skips it.
  * Without a dice texture it is a plain die with numbers.
  * A changed roll ("the next roll gets +1") only spins in, as before.
  * At rest it is tilted, so two more faces show (15.69).
  * Its shadow on the table grows as it falls.
  * It clatters as it first lands; silent without the private sound.
  * At rest the value is on top, seen from above (15.82).
  * The top fills most of the die; the front and right show as strips.
* **5e. Deck attack reveals an unattackable monster (19.23):** e.g. Stoney.
  * It stays in the chosen slot; the attack is cancelled (R-ATK-10).
  * No attack roll; the player goes on with the turn.
* **5f. Choice window (15.29):** a choice opens in the middle over the dimmed table.
  * Its options are big cards; nothing else takes clicks.
  * Discards and orders open the same way.
  * The cards lie in one row up to six; smaller from five (20.18).
  * A card option is the card itself: a click picks it; no name button.
* **5g. View board (15.29):** the player presses "View board".
  * The cards shrink away in 0.2 s; a small bar at the bottom names the choice.
  * Options on the table glow and can be clicked there.
  * "Back to the choice" opens the window again; the cards grow back in 0.2 s.
* **2g. Why not (15.29):** a card or icon that can't be used says why (UC-22).
* **7e. Plays left (V-46, 15.70):** End turn asks "End turn?" while it does not glow.
  * The question lists what is left: Attack, Buy, Play loot, Use abilities.
  * End turn, E, or Enter ends the turn.
  * Keep playing, Esc, or a click outside keeps the turn; nothing is sent.
  * A glowing End turn ends the turn at once.
  * Settings can turn the question off (UC-10, 3g).
* **2h. Choose one- (15.58):** a card with modes asks like 2d, e.g. Dagaz or Chaos Card.
  * The card shows large; one line per mode, in the card's words.
  * A mode with no valid target now stays, greyed: "No valid target now".
  * Cancel or a click outside cancels the play; the card stays.
  * "View board" hides the dialog, as in 5g.
  * A picked mode that needs a target then aims, as in 2f.
* **5i. Attack, buy and slot on the table (18.11):** no dialog asks.
  * What to attack, what to buy, a deck attack's slot.
  * The options glow on the table; the rest dims.
  * A deck option: the monster deck or the treasure deck glows.
  * A deck attack's drawn card stays face up on the monster deck.
  * One click picks; no cancel.
  * An option not on the table, e.g. an empty slot, keeps the dialog.
  * What to attack: an arrow runs from the own character (20.23).
  * It looks and acts as in 2e; the options keep their blue glow.
  * The other players see it, as any aim.
* **2i. Drag tilt (20.6):** a dragged card tilts by its motion, as if held.
  * The edge it moves toward leans away; faster moves lean more.
  * Held still, it lies flat again in about 0.2 s.
  * Animation speed Off or reduced motion: it stays flat.
  * A soft loop plays while it is held; a whoosh grows with its speed.
  * Both fade out on the drop; silent without the private sounds.
* **7f. Nothing left, said aloud (20.10):** End turn turns green and glows stronger.
  * 1.5 s later a voice says "Job's done!", once per turn.
  * No voice if a play became possible meanwhile.
  * Silent without the private sounds.
* **5j. Price on the coins (20.11):** while choosing what to buy, hover an option.
  * The own coins show "−N" beside the number: its price.
  * Red when the cents are short; the number dims meanwhile.
  * It goes when the pointer leaves the option.

* **5k. Cards dealt from a deck (20.18):** e.g. XVIII The Moon, Sleight of Hand.
  * Each card flies from that deck to its place and turns face up.
  * About 1 s each, 0.09 s apart; a card sound as each leaves.
  * Other card choices grow at their places in 0.2 s instead.
  * Animation speed Off or reduced motion: the cards are there at once.
* **5l. Hover and pick (20.18):** a hovered card lifts a little and glows brighter.
  * Its zoom waits as usual; a right click opens it at once (UC-11 2d).
  * A pick sounds; silent without the private sounds.
  * Cards from a deck then fly back to it, the kept one first, on top.
  * An order goes back bottom card first, top card last.
* **5m. Others watch (20.18):** other players see small backs over that deck.
  * One per card; they say "Ann looks at 5 cards"; never which.
  * They sink back into the deck when the choice ends.

* **2j. Pick up (20.20):** a click on a playable hand card picks it up.
  * It follows the pointer, button up; it grows a little; a grab sounds.
  * A drag past 8 px picks it up too.
  * A held card passes over other cards without zooming them.
* **2k. Over the board (20.20):** the held card's glow turns gold.
  * It turns green again off the board; no zone lights up.
  * A release or a click on the board plays it.
  * A click on a button there, e.g. End turn, puts it back.
* **2l. Put back (20.20):** a release or click off the board, a right click, or Esc.
  * The card flies home into the fan in 0.5 s; nothing is sent.
* **2m. Played (20.20):** the card stays where it was let go until the server answers.
  * Then it shrinks into its stack row in 0.4 s; the hand closes the gap.
  * Refused: it flies home; the reason shows (2c) and a refusal sounds.
  * A card that asks its mode goes home while the dialog asks (2h).
* **2n. Aim by drag (20.20):** a card that needs a target may be dragged.
  * On the board's edge it stops; the arrow runs from it (2e).
  * Back below the board, the arrow goes and the drag goes on.
  * Let go on a target, it plays there; elsewhere the aim stays.
  * A click on such a card still aims from the fan, as in 2e.
* **2o. Touch (20.20):** a tap lifts a playable hand card in the fan.
  * A tap on the board plays it; a tap anywhere else puts it back.
  * A drag still plays it; a long press zooms it.
* **2p. Sounds (20.20):** a cancelled aim sounds; so does every why note (2g).
  * Silent without the private sounds.
* **2q. Speed Off or reduced motion (20.20):** no growth, flight, shrink or slide.
  * A played card still waits where it was let go for the answer.
* **2r. No keys (20.20):** no key picks up or plays a hand card.

## Postconditions

* Only intents the server allows are offered; the server still checks each one.
