import type { Page } from '@playwright/test';
import { existsSync } from 'node:fs';
import { join } from 'node:path';

const root = join(import.meta.dirname, '..', '..');
const outDir = join(root, '.cache', 'screenshots');

// shot saves the page as screen-<name>.png when SCREENSHOTS is set.
export async function shot(page: Page, name: string): Promise<void> {
  if (process.env.SCREENSHOTS) {
    await page.screenshot({ path: join(outDir, `screen-${name}.png`) });
  }
}

// realMats serves the mat images from the private images submodule when
// SCREENSHOTS is set and the submodule is there, so screenshots show the
// real mat. Tests never depend on it.
export async function realMats(page: Page): Promise<void> {
  const dir = join(root, 'pkg', 'card_db', 'images', 'assets', 'mats');
  if (!process.env.SCREENSHOTS || !existsSync(dir)) {
    return;
  }
  await page.route('**/cardimg/assets/mats/*.webp', async (route) => {
    const file = join(dir, new URL(route.request().url()).pathname.split('/').pop() ?? '');
    await (existsSync(file) ? route.fulfill({ path: file, contentType: 'image/webp' }) : route.continue());
  });
}
