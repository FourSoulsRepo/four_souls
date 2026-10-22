// A server restart in the middle of a game (N-20, 8.24): the app
// reconnects by itself and the game goes on. These tests start and stop
// their own server, apart from the shared test servers.
import { spawn, execFileSync, type ChildProcess } from 'node:child_process';
import { connect } from 'node:net';
import { rmSync } from 'node:fs';
import { join } from 'node:path';
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { restartPorts } from './ports';
import { appURL } from './server';
import { shot } from './shot';

const root = join(import.meta.dirname, '..', '..');

// restartServer is this project's own server: built once, then started
// and stopped by the tests on the same folders.
class RestartServer {
  readonly port: number;
  private readonly dir: string;
  private readonly bin: string;
  private proc: ChildProcess | null = null;

  constructor(project: string) {
    this.port = restartPorts[project] ?? 4795;
    this.dir = join(root, '.cache', 'e2e', `restart-${project}`);
    this.bin = join(this.dir, 'server');
  }

  build(): void {
    rmSync(this.dir, { recursive: true, force: true });
    execFileSync('go', ['build', '-o', this.bin, './cmd/server'], { cwd: root, stdio: 'inherit' });
  }

  async start(): Promise<void> {
    const args = ['-addr', '127.0.0.1', '-port', String(this.port), '-records', join(this.dir, 'records'), '-saves', join(this.dir, 'saves')];
    this.proc = spawn(this.bin, args, { cwd: root, stdio: 'ignore' });
    for (let i = 0; i < 100; i++) {
      if (await listening(this.port)) {
        return;
      }
      await new Promise((r) => setTimeout(r, 100));
    }
    throw new Error('the restart server did not start');
  }

  // stop sends SIGTERM, as systemd and Docker do, and waits for the exit.
  async stop(): Promise<void> {
    const p = this.proc;
    this.proc = null;
    if (!p || p.exitCode !== null) {
      return;
    }
    const exited = new Promise((r) => p.once('exit', r));
    p.kill('SIGTERM');
    await exited;
  }
}

function listening(port: number): Promise<boolean> {
  return new Promise((resolve) => {
    const s = connect(port, '127.0.0.1');
    s.once('connect', () => {
      s.destroy();
      resolve(true);
    });
    s.once('error', () => {
      resolve(false);
    });
  });
}

// step is the update step the page shows.
function step(page: Page): Promise<number> {
  return page.evaluate(() => {
    interface S {
      game: { step: number } | null;
    }
    return (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.step ?? -1;
  });
}

let server: RestartServer;
// pages are this test's; they close with it, or they would keep
// reconnecting to the stopped server for the rest of the run.
let pages: Page[] = [];

test.describe.configure({ mode: 'serial' });

test.beforeAll(() => {
  server = new RestartServer(test.info().project.name);
  server.build();
});

test.afterEach(async () => {
  for (const page of pages) {
    await page.context().close();
  }
  pages = [];
  await server.stop();
});

test('a server restart mid-game: the game goes on, no dialog', async ({ browser }) => {
  await server.start();
  pages = await startGame(browser, ['Amy', 'Ben'], { port: server.port, url: `${appURL(server.port)}&anim=off` });
  const [amy, ben] = pages;
  if (!amy || !ben) {
    throw new Error('no pages');
  }
  const before = await step(amy);

  await server.stop();
  const badge = amy.getByRole('status', { name: 'Reconnecting' });
  await expect(badge).toContainText('The server is restarting');
  await shot(amy, 'reconnecting');

  await server.start();
  for (const page of [amy, ben]) {
    await expect(page.getByRole('status', { name: 'Reconnecting' })).toHaveCount(0, { timeout: 15000 });
    await expect(page.getByRole('dialog', { name: 'Paused' })).toHaveCount(0);
  }
  await expect(amy.locator('.player-bottom')).toHaveAccessibleName('Amy (you)');
  await expect(ben.locator('.player-bottom')).toHaveAccessibleName('Ben (you)');
  expect(await step(amy)).toBe(before);
  await expect(amy.getByRole('dialog', { name: 'Connection lost' })).toHaveCount(0);
  await expect(ben.getByRole('dialog', { name: 'Connection lost' })).toHaveCount(0);
});

test('a server that stays down: the badge becomes the dialog; Reconnect takes the seat back', async ({ browser }) => {
  await server.start();
  pages = await startGame(browser, ['Cy', 'Dot'], { port: server.port, url: `${appURL(server.port)}&anim=off` });
  const [cy, dot] = pages;
  if (!cy || !dot) {
    throw new Error('no pages');
  }
  for (const page of [cy, dot]) {
    await page.evaluate(() => {
      (window as unknown as { fourSouls: { lostDialogAfter: (ms: number) => void } }).fourSouls.lostDialogAfter(1500);
    });
  }
  await server.stop();
  const lost = cy.getByRole('dialog', { name: 'Connection lost' });
  await expect(lost).toBeVisible();
  await shot(cy, 'connection-lost');

  await server.start();
  await lost.getByRole('button', { name: 'Reconnect' }).click();
  await expect(lost).toHaveCount(0);
  await expect(cy.locator('.player-bottom')).toHaveAccessibleName('Cy (you)');
  await expect(dot.locator('.player-bottom')).toHaveAccessibleName('Dot (you)', { timeout: 15000 });
});

test('server notices: a newer engine, a lost game', async ({ browser }) => {
  await server.start();
  pages = await startGame(browser, ['Eve', 'Fin'], { port: server.port, url: `${appURL(server.port)}&anim=off` });
  const [eve] = pages;
  if (!eve) {
    throw new Error('no pages');
  }
  const receive = (kind: string) =>
    eve.evaluate((k) => {
      interface S {
        receive: (env: unknown) => void;
      }
      (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().receive({ type: 'notice', data: { kind: k, game: 'g1' } });
    }, kind);

  await receive('engine_updated');
  const engine = eve.getByRole('dialog', { name: 'The server was updated' });
  await expect(engine).toContainText('behave differently');
  await shot(eve, 'notice-engine');
  await engine.getByRole('button', { name: 'OK' }).click();
  await expect(engine).toHaveCount(0);
  await expect(eve.locator('.player-bottom')).toBeVisible();

  await receive('game_lost');
  const sorry = eve.getByRole('dialog', { name: 'Sorry, the game ended' });
  await expect(sorry).toContainText('start a new one');
  await shot(eve, 'notice-lost');
  await sorry.getByRole('button', { name: 'OK' }).click();
  await expect(eve.getByRole('heading', { name: 'Games', exact: true })).toBeVisible();
});
