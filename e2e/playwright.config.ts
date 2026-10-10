// End-to-end UI tests (ADR 008). Vite serves the frontend with fake
// Wails bindings (?fake), so no Wails window is needed; network play
// runs against the real Go server.
// SCREENSHOTS=1 also saves each screen into .cache/screenshots/.
import { defineConfig } from '@playwright/test';
import { serverPort } from './tests/server';

const port = 5199;

export default defineConfig({
  testDir: 'tests',
  fullyParallel: true,
  reporter: 'list',
  use: {
    baseURL: `http://127.0.0.1:${port}`,
    // Same size as the window in main.go.
    viewport: { width: 1024, height: 768 },
  },
  webServer: [
    {
      command: `npx vite --host 127.0.0.1 --port ${port} --strictPort`,
      cwd: '../frontend',
      url: `http://127.0.0.1:${port}`,
      reuseExistingServer: false,
    },
    {
      command: `go run ./cmd/server -addr 127.0.0.1 -port ${serverPort} -records .cache/e2e/records -saves .cache/e2e/saves`,
      cwd: '..',
      port: serverPort,
      reuseExistingServer: false,
    },
  ],
});
