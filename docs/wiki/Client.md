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

Effect durations scale with the animation setting (ST-05) and speed up when updates pile up; past 8 waiting updates they are skipped, so the table always catches up. Each effect has a time limit, so a closed screen never blocks playback.

## Cards

`Card` shows a card image, or a blank card with its text when the image is missing (A-11). Sizes: `hand`, `table`, `small`, `zoom`; everything inside scales with the width through container query units. Hover zooms a card; on touch a long press does (B-06). A card with an `objectId` pops in when it first appears and is a `data-zone` for effects.

## Player actions

`actions.ts` turns the update's allowed intents into what the table offers: buttons (pass, end turn, attack, buy), per-card actions (play, use an ability), a choice, or a discard. The table offers nothing else, so it never sends what the server would not allow; a refused intent still shows the engine's reason and rule ID.

* Cards the player may act on glow yellow; clicking or tapping one opens its action menu.
* A playable loot card can also be dragged onto the board (`data-drop="play"`). The drag is plain pointer events with pointer capture (`HandCard` in `PlayerArea.tsx`).
* Choices show as a dialog; card options show as cards (V-07). Options that point at the table (`view.options`, from the engine) glow blue there and can be clicked in place.
* The player acts only on the latest update (`store.behind`) and once per step (`store.sent`); meanwhile buttons stay but are disabled, so nothing blinks.

## Stack and waiting

`stack.ts` has the words: `stackLine` describes a stack item ("Ann plays A Dime!!", "Attack: Bo rolls 4", "2 damage to Gaper"), and `waitText` says who the game waits for and why ("Your turn", "Bo may respond", "Ann is choosing"). The stack panel sits in the board's top-right corner while the stack is not empty, the top item first.

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
