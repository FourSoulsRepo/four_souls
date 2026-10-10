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

When a screen is added, call `shot` in its test and add a row here.
