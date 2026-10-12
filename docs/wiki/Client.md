# Client

The desktop client is React + TypeScript in `frontend/` (ADR 002), with three libraries: Zustand for the store, Motion for animations and drag, nothing else (ADR 008). It keeps the KISS rules of [Frontend tests](Frontend-tests): the server decides, the client draws and sends intents.

## Files

| File | What it does |
|---|---|
| `src/bridge/` | The only code that calls Wails: settings, versions, card data, hosting. |
| `src/net.ts` | The one WebSocket to the game server. |
| `src/protocol.ts` | Message and engine view types, mirroring `internal/protocol` and the engine; `wsURL`, `errorText`. |
| `src/lobby.ts` | Pure: how server messages change the session (screen, lobby, error). |
| `src/playback.ts` | Pure: event → effect mapping, speed, and the update queue. |
| `src/table.ts` | Pure: where each player sits. |
| `src/cards.ts` | Pure: card kinds, text icons. |
| `src/store.ts` | The Zustand store; wires the pure parts to the bridge, the socket and playback. |
| `src/screens/` | One component per screen; `App.tsx` picks by `store.screen`. |
| `src/components/` | `Card`, `CardZoom`, `Board`, `PlayerArea`, `Pile`, `Effects`. |
| `src/fakeBridge.ts` | Dev only: fake Wails bindings for tests (`?fake`). |

## Data flow

1. A server message arrives in `store.receive`.
2. `lobby.onMessage` updates the screen and lobby.
3. An `update` also goes to playback, which queues it.
4. Playback plays the update's events as effects, one at a time (`store.playing`, drawn by `Effects`), then sets `store.game`.
5. The table (`GameScreen`) draws `store.game`; it never changes game state itself.

## Effects

`playback.effectsFor` maps engine events to effects: a flying card (looted, discarded), a dice, a popup (cents, healing, damage), a banner (turn, winner), a revealed card. Effects find their places on screen by `data-zone` attributes: `deck-loot`, `discard-monster`, `player-1`, `hand-0`, `obj-42` and so on, falling back to the next zone when one is not drawn.

A dice effect draws `Die`: one face cut out of the player's dice texture (V-19). `App.Dice` lists the looks with their image URLs, and the `dice` setting picks one. A texture is a 3×3 grid in the Tabletop Simulator layout: the top row is unused, then faces 2 4 5 and 1 3 6 (`dice.ts`). Without a texture, or for a modified value outside 1–6, the die is plain with a number.

## Game mats

The table lies on a mat (ST-03, V-24), drawn by `mats.matLook`. Three color mats are CSS: Green felt, Basement and Cathedral. Image mats come from `App.Mats`, one per file in `assets/mats/`; an image mat with a color mat's id takes its place in the menu, and a new id is added after them. The image is drawn under a dark shade, so cards and text stay readable, and over its colors (Basement for a new id), so a missing or unloadable image leaves the colors. The `mat` setting keeps any mat file name, so a saved mat whose image is gone also falls back to colors. An empty `mat` is the default: the first mat image, or Basement without images (`mats.matOf`). Green felt stays a choice but is never the default: with seat boxes it read as a poker table (H-05).

## Board

`Board` is one framed board in the middle of the table (F-01, H-01). Its grid has five columns: the info column (turn badge, stack, history), the monster and shop slots, the bonus souls, the decks, each with its discard (H-16), and the End turn button at the right edge (H-04). Nothing lies over it (F-02); the e2e tests check that every deck and discard is uncovered with 2, 3 and 4 players.

The board has one fixed shape, empty or full (H-06, V-23). `.board-cell` takes the space the players leave and is a size container; the board takes the largest size of its shape that fits there. `--u`, the width of a slot card, sizes everything on it, so no zone grows or shrinks when cards come and go. An empty slot keeps a dashed outline; up to four slots show full size, and more share the row. `--u` is used only outside cards: inside a card, container units would measure the card.

Players sit on the mat with no seat boxes (H-05); the active player's name is in the accent color. Choices open at the bottom edge, over the viewer's own area, so the board and its targets stay in view.

The viewer's own area (`OwnArea` in `PlayerArea.tsx`, 15.17) is a three-column grid: status, character, items, with the hand below. The character has a column of its own, so it stays centered under the board whatever the status and the items take (H-02). The status shows the name, the coins as a big number like Hearthstone's mana (H-18), the souls, and for the active player what is left this turn. The character card carries the live health and attack from `PlayerView.hp` and `atk` as big badges (`Card`'s `stats`, H-15); counters then move to its top corner. The hand box is shorter than its cards: they hang past the bottom edge, where `.game` clips them, so the box itself stays on screen (H-03). `fan` sets the angle and overlap on the hand from the card count, so a big hand still fits; each card turns by its `--spread` from the middle. A hovered card comes to the front; rising and glow wait for 15.27.

Buy and attack are round icons beside the own status (`TurnActions.tsx`, 15.19, F-11): a column right of the coins, so the status never grows taller than the character card and the board keeps its size. An icon is lit when the update allows its intent and the player may act, dim otherwise (F-12). A press marks the icon at once, until the next update (V-15); while the player picks what to buy or attack, that icon stays pressed (F-13). Keys B and A press them; the listener ignores text fields and modified keys. There is no Cancel: a declared purchase passes priority (R-SHOP-02).

End turn is one big button in the board's last column (`EndTurn` in `TurnActions.tsx`, 15.20, H-04); the column is part of the board's fixed shape, so the board never moves when the button appears. On the viewer's turn it is yellow, and dimmed while the engine does not allow ending the turn yet, e.g. while the stack resolves; on another player's turn it is grey and says "Their turn" (H-19). When the update allows nothing but passing and ending the turn (`TurnActions.idle`), it glows green (H-20); with reduced motion the glow stays still. The button looks like the game's cut-out art: the game's hand-drawn letters, a thick ink outline, a lumpy tilted shape and a hard drop shadow that it sinks onto when pressed. Pass and Skip all in the action bar, and the Buy and Attack icons, share the ink outline and shadow. The letters come from `--font-game` in `style.css`: font-souls Title for Latin, then the Russian card font, which the browser uses glyph by glyph for Cyrillic, then the system font. `fonts.ts` registers both faces from `App.Fonts()` URLs; without the private images there are no URLs, no faces, and the system font shows. Keys (`useKey` in `keys.ts`, F-15): E ends the turn, Space passes, as the Pass button does. Keys pressed in a text field, a dialog, or a menu are left alone. Space also never clicks a focused button, so it cannot end the turn by accident.

Opponents (`OpponentArea`, 15.18) show the same facts: a small name outside the board (H-21), coins with the same `Coin`, souls, the character with its `stats` badges, the items, and the hand as fanned backs (F-07, F-08, H-17). At the top the area mirrors the own area: a three-column grid keeps the character centered over the board, and the hand row is shorter than its backs, so they hang past the top edge, upside down as across a real table. In a side column (`.player-left`, `.player-right`, 150 px) the name comes first, clear of the corner menu, then the hand and the cards at `small` size, two to a row. The narrow columns leave the board more width, so with three players, where the top is free, it grows (F-09). The allowance wraps only between its parts.

The coin by the cents is the game's own coin (15.38): `App.Icons` gives `assets/icons/coin.webp`, the ¢ face cut out of the game's `coin.png` by the tools repo, or no URL when the file is missing; the store keeps it in `store.icons`. `Coin` in `PlayerArea.tsx` draws it, and falls back to the drawn ¢ coin when there is none or it fails to load.

## Card backs

Face-down cards show their deck's back (V-26). `App.Backs` gives the loot, treasure and monster back URLs from `assets/backs/<kind>.webp`, or no URL when the file is missing; the store keeps them in `store.backs`. `Card` draws the back image, and falls back to the plain dark back with the deck name when there is none or it fails to load. An opponent's hand is `Pile.HiddenHand`: one loot back per card, fanned, at most eight, with the card count, in the `hand-<id>` zone, so drawn cards fly to it. A back never carries a card ref, so it cannot reveal the card.

Effect durations scale with the animation setting (ST-05) and speed up when updates pile up; past 8 waiting updates they are skipped, so the table always catches up. Each effect has a time limit, so a closed screen never blocks playback.

## Cards

`Card` shows a card image, or a blank card with its text when the image is missing (A-11). Sizes: `hand`, `table`, `small`, `zoom`; everything inside scales with the width through container query units. Hover zooms a card; on touch a long press does (B-06). A card with an `objectId` pops in when it first appears and is a `data-zone` for effects.

## Player actions

`actions.ts` turns the update's allowed intents into what the table offers: buttons (pass), the turn actions (attack, buy, end turn), per-card actions (play, use an ability), a choice, or a discard. The table offers nothing else, so it never sends what the server would not allow; a refused intent still shows the engine's reason and rule ID.

* Cards the player may act on glow yellow; clicking or tapping one opens its action menu.
* A playable loot card can also be dragged onto the board (`data-drop="play"`). The drag is plain pointer events with pointer capture (`HandCard` in `PlayerArea.tsx`).
* Choices show as a dialog; card options show as cards (V-07). Options that point at the table (`view.options`, from the engine) glow blue there and can be clicked in place.
* The player acts only on the latest update (`store.behind`) and once per step (`store.sent`); meanwhile buttons stay but are disabled, so nothing blinks.

## Stack and waiting

`stack.ts` has the words: `stackLine` describes a stack item ("Ann plays A Dime!!", "Attack: Bo rolls 4", "2 damage to Gaper"), and `waitText` says who the game waits for and why ("Your turn", "Bo may respond", "Ann is choosing"). The stack panel sits in the board's info column, above the history, while the stack is not empty, the top item first.

The action bar shows the response timer from `update.deadline` (N-07) and "Skip all" / "Cancel skip" while something is on the stack (N-06); the server does the skipping.

## Info panels

`history.ts` builds the history from the events this viewer received: one entry per resolved stack item (`stack_resolved`), plus turns, loots, rewards, curses and the winner. The server already filters events per viewer, so a hidden loot reads "Bo loots a card". The history lives in `store.history`, grows as playback shows updates, and is lost on leaving the table (replays come in step 10).

The turn badge shows the turn number, the active player and the step. Player areas show what the engine computed: HP, ATK, cents, souls, hand size, and for the active player the loot plays, attacks and purchases left (`PlayerView.loot_plays`, `attacks`, `purchases`). Card counters show as a yellow badge (`CardView.counters`).

## Pause, reconnect, game over

* A lost connection during a game keeps the table (`lobby.lost`); `LostOverlay` offers Reconnect, which says hello again with the seat's token from `welcome`.
* The seat token also goes into the settings (`settings.rejoin`: address, token, game) whenever the player sits at a table. A restarted app sends it on the next connect to that address and lands back in its seat; `seatOf` then reads the seat from the welcome.
* `update.pause` shows `PauseOverlay` to the others: who is away, the votes, and Wait / Go on without them (`vote` messages).
* `view.over` shows `EndOverlay`: the winner and Save record. `App.SaveRecord` (Go, `record.go`) asks where to save and downloads `http://<server>/record?game=…&token=…`; it accepts only such links.

## Save and continue

The in-game **Menu** (`GameMenu`) offers **Save and leave** to seat 1, the host (`save` message), and **Leave game** to all. After `saved`, everyone is back on the server's game list. **Continue a saved game** on the host screen starts the in-app server and shows that list (`store.hostSaved`); **Continue** sends `load`. A loaded table shows its seats as "back" or "not back yet".
