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

| ID | What it does | Started by | Claims | Code |
|----|--------------|------------|--------|------|
| `card-pop` | A table card scales and fades in when it first appears. | The view draws a `Card` with a new `objectId`. | Nothing; stays still when the card is claimed `arrive` this beat (19.6). | `components/Card.tsx` |
| `fly` | A card flies from one zone to another. As a `place` effect it stays on the slot until the view is drawn (19.3). | `looted`, `discarded` (`fly`); `card_revealed` of a slot's new top card (`place`). | `arrive` of the card, as `fly` or `place`. | `components/Effects.tsx` (`Fly`) |
| `reveal` | A card shows large in the middle. | `card_revealed`, unless the update before revealed it (19.4). | `show` of the card. | `components/Effects.tsx` |
| `dice` | A die spins in with the changed roll and its label. | `roll_changed`. | Nothing. | `components/Effects.tsx` |
| `dice-throw` | The thrown die falls from above, bounces twice and tumbles to the engine's value; its label fades in after the landing (15.62). | `dice_rolled`: a fresh roll or a reroll. | Nothing. | `components/Effects.tsx` (`ThrownDie`) |
| `banner` | A text banner in the middle. | `turn_started`, `game_won`. | Nothing. | `components/Effects.tsx` |
| `popup` | A short "+N" text rises from a zone. | `gained_cents`, `healed`, `prevented`, `stole`. | Nothing. | `components/Effects.tsx` (`Popup`) |
| `lunge` | A copy of the attacker jumps into the target and back; the card hides meanwhile. | `damaged` with `text: "combat"` and a `source`. | Nothing; it moves a copy. | `components/lunge.ts` |
| `burst` | A starburst "−N" on a damaged card. | The hit of a `lunge`, or of a `hit-shot`. | Nothing. | `components/lunge.ts` (`burst`) |
| `hit-shot` | A glowing shot flies from a damage source to its target. | `damaged` without combat (`hits` effect). | Nothing. | `components/marks.ts` |
| `death-skull` | A big skull pops on a dying card. | `died` (`deaths` effect). | Nothing. | `components/marks.ts` |
| `end-turn-glow` | The End turn button glows while the player may end the turn. | CSS class `.end-turn.idle`. | Nothing; not a game event. | `components/Actions.css` |
| `aim-march` | The aiming arrow's blocks march toward its head. | Aiming a card (15.32). | Nothing; not a game event. | `components/Table.css` |
| `aim-fade-in` | The arrow's first block fades in, in step with the march. | Aiming a card. | Nothing; not a game event. | `components/Table.css` |
| `aim-fade-out` | The arrow's last block fades out, in step with the march. | Aiming a card. | Nothing; not a game event. | `components/Table.css` |
| `aim-why` | The hint under an invalid aim fades out. | Aiming at a card that is not a target; a click on a card or turn icon the player can't use now (15.29). | Nothing; not a game event. | `components/Table.css` |
| `hand-glow` | The green aura behind a playable hand card breathes slowly. | CSS class `.hand-card.usable` (15.27). | Nothing; not a game event. | `components/Table.css` |
| `prompt-glow` | The gold glow behind the prompt banner, and behind Pass when passing answers it, breathes. | A prompt for the viewer (15.28). | Nothing; not a game event. | `components/Actions.css` |
| `throw-glow` | The halo behind the roller's Throw button breathes. | A throw prompt for the viewer (15.61). | Nothing; not a game event. | `components/Throw.css` |
| `throw-shake` | The roller's die shakes in their hand while the others wait. | A throw prompt for another player (15.61). | Nothing; not a game event. | `components/Throw.css` |

CSS `transition`s on hover and focus are not listed: they answer the pointer, not the game, and cannot clash with an effect.

## Adding an animation

1. Pick an id; add the `anim: <id>` comment above the animation in the code.
2. Add its row to the table above.
3. If it brings a card to its place or shows it large, claim the card (`claimEffect` or `claim`).
4. If it replaces an older animation for the same event, remove the old one with its tag and row.
5. Add an e2e test that plays it; the clash fixture then checks it with every other test.
