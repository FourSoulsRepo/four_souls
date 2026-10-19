# Layers

Which window draws over which. The decision behind these rules is ADR 015; the numbers live in `frontend/src/layers.css`, and the Vitest test `layer inventory` (`layers.test.ts`) checks them against this page.

## The rules

1. **One list, top first.** Every layer is a CSS variable `--z-<id>` in `layers.css`. Only that file holds z-index numbers. The table below lists the same ids in the same order; the first row is drawn on top.
2. **Code names a layer, never a number.** A `z-index` elsewhere is `var(--z-<id>)`. The one exception is a number that only orders the parts of one element (a glow behind its card, a badge on its button); it carries `/* layer: local */` on its line.
3. **Every layer has a reason.** Its row says what it must cover, and why. A new window gets a row before it gets a number.
4. **The test fails** when a bare z-index turns up, when the numbers are out of order, when a layer is unused, or when this table and `layers.css` disagree.
5. **The pairs a player meets get an e2e test.** `e2e/tests/layers.spec.ts` checks which one is drawn on top where two overlap.

A layer only counts on the window's own stacking: these elements sit in `.game` or the app root, neither of which starts a stacking context. An element inside a positioned parent with its own z-index stacks only among its siblings; that is what `layer: local` marks.

## Inventory

| ID | What | Above everything below because | Code |
|----|------|--------------------------------|------|
| `aim-why` | The note why an aim is no target. | It answers the pointer, right under it; it covers its own arrow. | `components/Table.css` |
| `aim` | The aiming arrow. | It is the pointer while aiming; the target under it may be any window's card. | `components/Table.css` |
| `trail` | The attack lunge and its starburst. | It plays over the effects of the same hit. | `components/Effects.css` |
| `effects` | Event effects: flying cards, reveal, dice, banner, popups. | A game event must be seen, even over a zoom or an overlay. It never catches the pointer. | `components/Effects.css` |
| `zoom` | The large copy of a hovered or held card. | The player reads a card from any window: discard, choice, look, order, a dialog (19.26). A mouse zoom lets the pointer through. | `components/Card.css` |
| `overlay` | Connection lost, paused, game over, the game menu. | Modal: the table and its windows wait under it. | `components/Overlays.css` |
| `drag` | A hand card being dragged to play. | It follows the pointer over the table and its windows. | `components/Table.css` |
| `throw` | The roll waiting for its throw, with its button. | The thrown die lands here; only the button catches clicks. | `components/Throw.css` |
| `menu` | A card's action menu and the ability dialog, with their backdrop. | The player opened it on purpose, over any prompt. | `components/Actions.css` |
| `look` | The "you look at" window in the middle. | It may come with a choice; the cards to look at come first. | `components/Actions.css` |
| `choice` | Choice, discard and order windows at the bottom. | A prompt the game waits on; it covers the player's own area. | `components/Actions.css` |
| `hand` | A hovered or aiming hand card. | It comes to the front of the fan, but never over a prompt's buttons. | `components/Table.css` |
| `menu-button` | The Menu button. | Over the table, under every window. | `components/Overlays.css` |
| `action-bar` | The status, timer and skip buttons, bottom right. | Over the table, under every window. | `components/Actions.css` |

## Adding a layer

1. Pick an id; decide which rows it must cover and which must cover it.
2. Add `--z-<id>` to `layers.css` between them; keep the numbers falling.
3. Add its row to the table above, in the same place.
4. Use `z-index: var(--z-<id>)` in the code.
5. If it can overlap another window, add the pair to `layers.spec.ts`.
