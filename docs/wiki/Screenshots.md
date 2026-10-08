# Screenshots

`make screenshots` captures every screen into `docs/wiki/images/screen-*.png`. Use it to check UI changes and to keep these docs current.

It works the same way as in the owner's other Wails projects:

1. Vite serves the frontend with `?screenshot=1`.
2. That flag loads `frontend/src/screenshotBridge.ts`, which installs fake Wails bindings, so no Go process or real window is needed.
3. Playwright (headless Chromium, 1024×768 like the app window) opens the page and saves each screen.

Playwright lives in `scripts/screenshots/`, not in `frontend/`, so normal frontend installs and release builds never download Chromium.

The fixture copies the fan-game notice. `internal/legal` has a test that fails when the two drift apart.

| Screen | File |
|--------|------|
| Splash with the fan-game notice | ![Splash](images/screen-splash.png) |
| Main menu with versions | ![Main menu](images/screen-menu.png) |

When a screen is added, add a step to `scripts/screenshots/capture.mjs` and a row here.
