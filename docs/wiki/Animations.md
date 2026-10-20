# Animations

Every animation in the client, what starts it, and which card it owns. The decision behind these rules is ADR 012; the code that checks them is `frontend/src/animations.ts`.

## The rules

1. **One card, one animation per job.** A card gets one arrival on the table per update, and is shown large once. Two animations that bring or show the same card *clash*; that was 19.3, 19.4 and 19.6.
2. **Events own the animation; the view follows.** Playback plays an update's events as effects, then draws the update's view (see [Client](Client)). An effect that brings a card to its place *claims* it. When the view then draws that card, it does not animate it again: `Card`'s pop-in only plays for cards no effect brought.
3. **Every animation is listed here.** Each animation in the code has an `anim: <id>` comment on the line above it (`<motion.*>`, `.animate(`, `@keyframes`). Each id has one row in the table below. The Vitest test `animation inventory` (`animations.test.ts`) fails when a site has no tag, or when the code and this table disagree.
4. **A new animation that brings or shows a card claims it.** Map its effect in `claimEffect`, or call `claim(card, job, id)` before it starts. A clash is logged as `animation clash: …` in the console, kept in `clashes`, and the second animation is skipped.
5. **Every e2e test fails on a clash.** Specs import `test` and `expect` from `e2e/tests/fixtures.ts`, not from `@playwright/test`. Its automatic fixture reads `window.fourSouls.clashes` on every page after each test.
6. **Removing or replacing an animation removes its row and tag.** The inventory test makes a leftover row fail.

## Beats and jobs

Claims are counted in *beats*: a beat is one update, from the start of its effects until the next update's effects start. It covers the cards the update's view mounts, so a pop-in sees the effect that brought its card.

| Job | Means | Clashes with |
|-----|-------|--------------|
| `arrive` | The card comes to its place on the table: a hand, a slot, a discard pile. | Another animation's `arrive` of that card in the same beat. The same animation may claim again (two copies of a card). |
| `show` | The card shows large in the middle of the screen. | Any `show` of that card in the same or the previous beat (19.4). |

## Inventory

| ID | What it does | Why: what the player learns | Direction | Started by | Claims | Code |
|----|--------------|-----------------------------|-----------|------------|--------|------|
| `card-pop` | A table card scales and fades in when it first appears. | A card is new on the table. | In place: grows from 60 % where the card sits. | The view draws a `Card` with a new `objectId`. | Nothing; stays still when the card is claimed `arrive` this beat (19.6). | `components/Card.tsx` |
| `fly` | A card flies from one zone to another. As a `place` effect it stays on the slot until the view is drawn (19.3). | Where a card went: who drew loot, what was thrown away, which card filled a slot. | `looted`: loot deck (or loot discard) → the player's hand, else their panel. `discarded`: the card's spot on the table, else the player's hand, else its deck → that kind's discard. `place`: the kind's deck → the monster or shop slot. Right after a held `reveal`: from the middle, large, shrinking to a table card (15.73). | `looted`, `discarded` (`fly`); `card_revealed` of a slot's new top card (`place`). | `arrive` of the card, as `fly` or `place`. | `components/Effects.tsx` (`Fly`) |
| `reveal` | A card grows out of its deck and shows large in the middle (15.73). | Which card just turned up, so everyone can read it. | Its deck → the middle. Then: held in the middle when its next effect moves it; else back onto its deck, shrinking. No deck on screen: fades in and out in the middle. | `card_revealed`, unless the update before revealed it (19.4). | `show` of the card. | `components/Effects.tsx` |
| `dice` | A die spins in with the changed roll and its label. | A roll was changed, and its new value. | Middle of the screen; spins half a turn in from 40 %. Not tied to the roller. | `roll_changed`. | Nothing. | `components/Effects.tsx` |
| `dice-throw` | The thrown die, a 3D cube (`DieCube`, 15.69), falls from above, bounces twice and tumbles to rest with the engine's value in front; its shadow grows on the table, and its label fades in after the landing (15.62). | A die was thrown, and the value it landed on. | Middle of the screen, where the Throw button waited; straight down from above the screen. | `dice_rolled`: a fresh roll or a reroll. | Nothing. | `components/Effects.tsx` (`ThrownDie`) |
| `banner` | A text banner in the middle. | Whose turn starts; who won. | Middle of the screen; fades in place. | `turn_started`, `game_won`. | Nothing. | `components/Effects.tsx` |
| `popup` | A short "+N" text rises from a zone. | A player gained or stole cents; a card or player healed; damage was prevented. | Rises up from the receiver: the player's panel (cents, stolen), else the healed or protected card. | `gained_cents`, `healed`, `prevented`, `stole`. | Nothing. | `components/Effects.tsx` (`Popup`) |
| `lunge` | A copy of the attacker jumps into the target and back; the card hides meanwhile. | Who attacks whom in a fight. | Attacker's card → target card (short of its centre) → back; a small pull-back first. | `damaged` with `text: "combat"` and a `source`. | Nothing; it moves a copy. | `components/lunge.ts` |
| `burst` | A starburst "−N" on a damaged card. | Which card took damage, and how much. | In place, on the damaged card (or the player's character). | The hit of a `lunge`, or of a `hit-shot`. | Nothing. | `components/lunge.ts` (`burst`) |
| `hit-shot` | A glowing shot flies from a damage source to its target. | Where damage outside a fight came from, and what it hit. | Source card on the table, else its stack row → the damaged card. No source on screen: no shot, only the `burst`. | `damaged` without combat (`hits` effect). | Nothing. | `components/marks.ts` |
| `death-skull` | A big skull pops on a dying card. | Which card died. | In place, on the dying card. | `died` (`deaths` effect). | Nothing. | `components/marks.ts` |
| `end-turn-glow` | The End turn button glows while the player may end the turn. | Nothing is left to do; end the turn. | In place, on End turn. | CSS class `.end-turn.idle`. | Nothing; not a game event. | `components/Actions.css` |
| `aim-march` | The aiming arrow's blocks march toward its head. | Which way the play goes: from the card to its target. | Along the arrow: the played card → the arrow's head at the pointer. | Aiming a card (15.32). | Nothing; not a game event. | `components/Table.css` |
| `aim-fade-in` | The arrow's first block fades in, in step with the march. | None of its own; keeps the march seamless. | In place, at the arrow's tail by the card. | Aiming a card. | Nothing; not a game event. | `components/Table.css` |
| `aim-fade-out` | The arrow's last block fades out, in step with the march. | None of its own; keeps the march seamless. | In place, at the arrow's head. | Aiming a card. | Nothing; not a game event. | `components/Table.css` |
| `aim-why` | The hint under an invalid aim fades out. | Why that click does nothing; it leaves by itself. | In place, just above the pointer. | Aiming at a card that is not a target; a click on a card or turn icon the player can't use now (15.29). | Nothing; not a game event. | `components/Table.css` |
| `hand-glow` | The green aura behind a playable hand card breathes slowly. | This hand card can be played now. | In place, behind the hand card. | CSS class `.hand-card.usable` (15.27). | Nothing; not a game event. | `components/Table.css` |
| `prompt-glow` | The gold glow behind the prompt banner, and behind Pass when passing answers it, breathes. | The game waits for you; Pass answers it. | In place, behind the banner and Pass. | A prompt for the viewer (15.28). | Nothing; not a game event. | `components/Actions.css` |
| `throw-glow` | The halo behind the roller's Throw button breathes. | It is your roll: press Throw. | In place, behind Throw in the middle. | A throw prompt for the viewer (15.61). | Nothing; not a game event. | `components/Throw.css` |
| `throw-shake` | The roller's die shakes in their hand while the others wait. | Another player is about to roll; wait. | In place, in the middle where the die will land, under the roller's name. | A throw prompt for another player (15.61). | Nothing; not a game event. | `components/Throw.css` |
| `played` | Another player's loot card shows large in the middle with "Ann plays", then fades (15.31, H-10). | Which loot card another player played. | Middle of the screen; grows from 60 %, shrinks to 50 % as it fades. Not tied to the player's hand. | `loot_played` by another player. | `show` of the card. | `components/Effects.tsx` |

CSS `transition`s on hover and focus are not listed: they answer the pointer, not the game, and cannot clash with an effect.

## Adding an animation

1. Pick an id; add the `anim: <id>` comment above the animation in the code.
2. Add its row to the table above, with its Why and Direction.
3. If it brings a card to its place or shows it large, claim the card (`claimEffect` or `claim`).
4. If it replaces an older animation for the same event, remove the old one with its tag and row.
5. Add an e2e test that plays it; the clash fixture then checks it with every other test.
