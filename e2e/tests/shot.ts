import type { Page } from '@playwright/test';
import { join } from 'node:path';

const outDir = join(import.meta.dirname, '..', '..', '.cache', 'screenshots');

// shot saves the page as screen-<name>.png when SCREENSHOTS is set.
export async function shot(page: Page, name: string): Promise<void> {
  if (process.env.SCREENSHOTS) {
    await page.screenshot({ path: join(outDir, `screen-${name}.png`) });
  }
}
