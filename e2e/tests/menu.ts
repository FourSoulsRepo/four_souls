// The main menu's paths (M-03): Play opens LAN and Online in place.
import type { Page } from '@playwright/test';
import { tr } from './locale';

// play opens the Play group, unless it shows already.
export async function play(page: Page): Promise<void> {
  const button = page.getByRole('button', { name: 'Play', exact: true });
  if ((await button.getAttribute('aria-expanded')) !== 'true') {
    await button.click();
  }
}

// openLan opens the LAN screen from the main menu.
export async function openLan(page: Page): Promise<void> {
  await play(page);
  await page.getByRole('button', { name: 'LAN', exact: true }).click();
}

// openOnline opens the Online screen from the main menu.
export async function openOnline(page: Page): Promise<void> {
  await play(page);
  await page.getByRole('button', { name: 'Online', exact: true }).click();
}

// openHost opens the host setup from the LAN screen.
export async function openHost(page: Page): Promise<void> {
  await openLan(page);
  await page.getByRole('button', { name: 'Host a game' }).click();
}

// backToMenu presses Back until the main menu shows. Settings' Back
// is in the tests' language (16.2).
export async function backToMenu(page: Page): Promise<void> {
  const menu = page.getByRole('button', { name: 'Play', exact: true });
  const back = page.getByRole('button', { name: 'Back' }).or(page.getByRole('button', { name: tr('settings.back') })).first();
  for (let i = 0; i < 4 && !(await menu.isVisible()); i++) {
    await back.click();
    await menu.or(back).first().waitFor();
  }
  await menu.waitFor();
}
