import type { Page } from '@playwright/test';
import { existsSync } from 'node:fs';
import { join, normalize } from 'node:path';

const root = join(import.meta.dirname, '..', '..');
const outDir = join(root, '.cache', 'screenshots');

// shot saves the page as screen-<name>.png when SCREENSHOTS is set.
export async function shot(page: Page, name: string): Promise<void> {
  if (process.env.SCREENSHOTS) {
    await page.screenshot({ path: join(outDir, `screen-${name}.png`) });
  }
}

// realAssets serves the mats and card backs from the private images
// submodule when SCREENSHOTS is set and the submodule is there, so
// screenshots show the real images. Tests never depend on it.
export async function realAssets(page: Page): Promise<void> {
  const dir = join(root, 'pkg', 'card_db', 'images', 'assets');
  if (!process.env.SCREENSHOTS || !existsSync(dir)) {
    return;
  }
  await page.route('**/cardimg/assets/**/*.webp', async (route) => {
    const name = normalize(new URL(route.request().url()).pathname.replace(/^\/cardimg\/assets\//, ''));
    const file = join(dir, name);
    const inside = !name.startsWith('..') && existsSync(file);
    await (inside ? route.fulfill({ path: file, contentType: 'image/webp' }) : route.continue());
  });
}
