// End-to-end UI tests (ADR 008). Vite serves the frontend with fake
// Wails bindings (?fake), so no Wails window is needed; network play
// runs against the real Go server.
// SCREENSHOTS=1 also saves each screen into .cache/screenshots/.
import { defineConfig } from '@playwright/test';
import { serverPorts } from './tests/ports';

const port = 5199;

export default defineConfig({
  testDir: 'tests',
  fullyParallel: true,
  // Each test opens two to four pages; they close when it ends (19.38).
  // More workers made known flaky tests fail more often (19.38 sweep);
  // E2E_WORKERS in the Makefile overrides it.
  workers: 3,
  reporter: 'list',
  use: {
    baseURL: `http://127.0.0.1:${port}`,
    // Same size as the window in main.go.
    viewport: { width: 1024, height: 768 },
  },
  // Chromium, and WebKit as the nearest check for Linux's WebKitGTK and
  // macOS's WKWebView (7.12).
  projects: [
    { name: 'chromium', use: { browserName: 'chromium' } },
    { name: 'webkit', use: { browserName: 'webkit' } },
  ],
  webServer: [
    {
      command: `npx vite --host 127.0.0.1 --port ${port} --strictPort`,
      cwd: '../frontend',
      // No file watching: see frontend/vite.config.ts.
      env: { FOUR_SOULS_E2E: '1' },
      url: `http://127.0.0.1:${port}`,
      reuseExistingServer: false,
    },
    ...Object.values(serverPorts).map((p) => ({
      // A fresh folder per run: saved games from earlier runs would show.
      // The name shows instead of the address (N-18).
      command: `rm -rf .cache/e2e/${String(p)} && go run ./cmd/server -name "Test server" -addr 127.0.0.1 -port ${String(p)} -records .cache/e2e/${String(p)}/records -saves .cache/e2e/${String(p)}/saves`,
      cwd: '..',
      port: p,
      reuseExistingServer: false,
    })),
  ],
});
