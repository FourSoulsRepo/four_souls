import type { Page } from '@playwright/test';
import { existsSync } from 'node:fs';
import { extname, join, normalize } from 'node:path';

const root = join(import.meta.dirname, '..', '..');
const outDir = join(root, '.cache', 'screenshots');

// shot saves the page as screen-<name>.png when SCREENSHOTS is set.
export async function shot(page: Page, name: string): Promise<void> {
  if (process.env.SCREENSHOTS) {
    await page.screenshot({ path: join(outDir, `screen-${name}.png`) });
  }
}

// realAsset tells whether realAssets serves this private file.
export function realAsset(name: string): boolean {
  return Boolean(process.env.SCREENSHOTS) && existsSync(join(root, 'pkg', 'card_db', 'images', name));
}

// realAssets serves cards, backs, mats, dice, fonts and sounds from the
// private images submodule when SCREENSHOTS is set and the submodule is
// there, so screenshots show the real images. Tests never depend on it.
export async function realAssets(page: Page): Promise<void> {
  const dir = join(root, 'pkg', 'card_db', 'images');
  if (!process.env.SCREENSHOTS || !existsSync(dir)) {
    return;
  }
  await page.route(/\/cardimg\/.+\.(webp|ttf|mp3)$/, async (route) => {
    const name = normalize(new URL(route.request().url()).pathname.replace(/^\/cardimg\//, ''));
    const file = join(dir, name);
    const inside = !name.startsWith('..') && existsSync(file);
    const types: Record<string, string> = { '.ttf': 'font/ttf', '.mp3': 'audio/mpeg' };
    const contentType = types[extname(name)] ?? 'image/webp';
    await (inside ? route.fulfill({ path: file, contentType }) : route.continue());
  });
}
