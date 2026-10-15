# Screenshots

`make screenshots` captures every screen into `.cache/screenshots/screen-*.png`. Use it to check UI changes during development.

The folder is git-ignored: the UI changes often, and images would bloat the repository. Screenshots go into the docs once the UI settles.

Screenshots are taken by the Playwright UI tests in `e2e/` (see [Frontend tests](Frontend-tests)): `make screenshots` runs `make e2e` with `SCREENSHOTS=1`, and each test calls `shot(page, name)` where a screen is worth a picture.

1. Vite serves the frontend with `?fake`.
2. That flag loads `frontend/src/fakeBridge.ts`, which installs fake Wails bindings, so no Go process or real window is needed. The fixture is loaded only by the Vite dev server; release builds never contain it.
3. Playwright (headless Chromium, 1024×768 like the app window) opens the page.

Playwright lives in `e2e/`, not in `frontend/`, so normal frontend installs and release builds never download Chromium.

The fixture copies the fan-game notice. `internal/legal` has a test that fails when the two drift apart.

| Screen | File |
|--------|------|
| Splash with the fan-game notice | `screen-splash.png` |
| Main menu with versions | `screen-menu.png` |
| Settings | `screen-settings.png` |
| Host a game, detailed setup | `screen-host.png` |
| Game list on a server | `screen-games.png` |
| Table before the game | `screen-table.png` |
| Card gallery: image and blank cards (dev only) | `screen-cards.png` |
| Zoomed card | `screen-card-zoom.png` |
| Game table, 2 players, 1280×720 | `screen-table-2p.png` |
| Game table, 4 players, 1280×720 | `screen-table-4p.png` |
| Dice roll effect | `screen-effect-dice.png` |
| Your turn: playable cards highlighted | `screen-actions.png` |
| A card's action menu (touch) | `screen-action-menu.png` |
| The stack, seen by a player who may respond | `screen-stack.png` |
| History, turn badge, what is left this turn | `screen-history.png` |
| First start: choose a nickname | `screen-welcome.png` |
| Settings in the light theme | `screen-settings-light.png` |
| Paused: waiting for a player who left | `screen-paused.png` |
| Connection lost, with Reconnect | `screen-connection-lost.png` |
| Game over: the winner and the record | `screen-game-over.png` |
| A saved game's table: who is back | `screen-saved-game.png` |
| History strip with attack and loot badges | `screen-history-strip.png` |
| History strip: hovering an icon shows its card | `screen-history-strip-hover.png` |
| Ability dialog: a card with several abilities | `screen-ability-dialog.png` |
| Own turn: buys and attacks left on their icons | `screen-turn-left.png` |
| History strip: damage a player takes names its card | `screen-history-strip-taken-damage.png` |
| Damage shots on their way to their targets | `screen-damage-shots.png` |
| Damage bursts on every target at once | `screen-damage-bursts.png` |
| Hurt health in red | `screen-damage-red-hp.png` |
| A skull on a dying monster | `screen-death-skull.png` |
| Playable hand cards glow green | `screen-hand-glow.png` |
| A hovered hand card rises large | `screen-hand-hover.png` |
| A monster's zoom right beside it | `screen-monster-hover.png` |
| Keyword tips beside a zoomed card | `screen-zoom-tips.png` |

When a screen is added, call `shot` in its test and add a row here.
