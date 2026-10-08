#!/usr/bin/env node
// Captures each screen into docs/wiki/images/screen-*.png.
// Uses ?screenshot=1 so Wails bindings are faked (no Go/Wails window needed).

import { spawn } from 'node:child_process';
import { createServer } from 'node:net';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { chromium } from 'playwright';

const here = dirname(fileURLToPath(import.meta.url));
const repo = join(here, '..', '..');
const frontend = join(repo, 'frontend');
const outDir = join(repo, 'docs', 'wiki', 'images');

// Same size as the window in main.go.
const width = 1024;
const height = 768;

function unusedPort() {
  return new Promise((resolve, reject) => {
    const s = createServer();
    s.listen(0, '127.0.0.1', () => {
      const { port } = s.address();
      s.close((err) => (err ? reject(err) : resolve(port)));
    });
    s.on('error', reject);
  });
}

function startVite(port) {
  const child = spawn('npx', ['vite', '--host', '127.0.0.1', '--port', String(port), '--strictPort'], {
    cwd: frontend,
    stdio: ['ignore', 'pipe', 'pipe'],
    env: { ...process.env, BROWSER: 'none' },
    // Own process group, so stop() also kills the vite process npx starts.
    detached: true,
  });
  let output = '';
  child.stdout.on('data', (b) => (output += b.toString()));
  child.stderr.on('data', (b) => (output += b.toString()));

  const ready = (async () => {
    const start = Date.now();
    while (Date.now() - start < 30000) {
      if (child.exitCode !== null) {
        throw new Error(`vite exited ${child.exitCode}:\n${output}`);
      }
      try {
        const res = await fetch(`http://127.0.0.1:${port}/`);
        if (res.ok) {
          return;
        }
      } catch {
        // Not ready yet.
      }
      await new Promise((r) => setTimeout(r, 100));
    }
    throw new Error(`vite did not start within 30s:\n${output}`);
  })();

  const stop = () => {
    try {
      process.kill(-child.pid, 'SIGKILL');
    } catch {
      // Already gone.
    }
  };

  return { ready, stop };
}

async function shot(page, name) {
  const dest = join(outDir, `screen-${name}.png`);
  await page.screenshot({ path: dest, type: 'png' });
  console.log('wrote', dest);
}

async function main() {
  const port = await unusedPort();
  const vite = startVite(port);
  const browser = await chromium.launch();
  try {
    await vite.ready;
    const page = await browser.newPage({ viewport: { width, height }, deviceScaleFactor: 1 });
    await page.goto(`http://127.0.0.1:${port}/?screenshot=1`, { waitUntil: 'networkidle' });

    // Splash with the fan-game notice.
    await page.getByText('Four Souls is an unofficial').waitFor();
    await shot(page, 'splash');

    // Main menu appears after the splash timer.
    await page.getByRole('button', { name: 'Host a game' }).waitFor({ timeout: 10000 });
    await page.getByText('rules engine').waitFor();
    await shot(page, 'menu');
  } finally {
    await browser.close().catch(() => {});
    vite.stop();
  }
}

main().then(
  () => process.exit(0),
  (err) => {
    console.error(err);
    process.exit(1);
  },
);
