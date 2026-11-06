// The rules website (12.2, ADR 026): the page make site builds, served
// by vite preview (playwright.config.ts). A check link shows the
// engine's answer with its rule; the rule link opens the rules text.
// Its lines are English: the site picks its language from the link,
// and a built page has no qps (LOCALE runs keep English here).
import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import type { Page } from '@playwright/test';
import { compact, encodePayload } from '../../frontend/src/site/link';
import { expect, test } from './fixtures';
import { sitePort } from './ports';
import { shot } from './shot';
import { trIn } from './locale';

const root = join(import.meta.dirname, '..', '..');
const built = (() => {
  try {
    readFileSync(join(root, 'frontend', 'dist-site', 'index.html'));
    return true;
  } catch {
    return false;
  }
})();
const site = `http://127.0.0.1:${String(sitePort)}/four_souls/`;

const en = (key: string, values?: Record<string, string | number>) => trIn('en', key, values);

// checkLink is the site's link to a payload of internal/website/testdata.
async function checkLink(name: string): Promise<string> {
  const payload = readFileSync(join(root, 'internal', 'website', 'testdata', name), 'utf8');
  return `${site}#/en/check/${await encodePayload(compact(payload))}`;
}

async function expectRefusal(page: Page): Promise<void> {
  const answer = page.locator('.answer');
  await expect(answer.getByRole('heading', { name: en('site.check.refused') })).toBeVisible();
  await expect(answer).toContainText(en('refuse.loot_your_turn'));
  await expect(answer.getByRole('link', { name: 'R-CARD-08' })).toBeVisible();
  await expect(answer).toContainText(en('site.check.matches'));
}

test.skip(!built, 'the rules website is not built: run make site (make e2e does)');

test('a check link shows the refusal; its rule link opens the rule', async ({ page }) => {
  await page.goto(await checkLink('loot-needs-a-loot-play.json'));
  await expectRefusal(page);
  await expect(page.locator('.site-footer')).toContainText('fan-made');
  await shot(page, 'site-check');

  await page.locator('.answer').getByRole('link', { name: 'R-CARD-08' }).click();
  await expect(page).toHaveURL(/#\/en\/rules\/R-CARD-08$/);
  const rule = page.locator('#rule-R-CARD-08');
  await expect(rule).toHaveAttribute('aria-current', 'true');
  await expect(rule).toBeInViewport();
  await expect(rule).toContainText('Playing a loot card needs priority');
  await shot(page, 'site-rules');
});

test('a copied link opens the same situation, on a phone too', async ({ browser }) => {
  const context = await browser.newContext({ viewport: { width: 390, height: 844 } });
  const page = await context.newPage();
  await page.goto(`${site}#/en/check`);
  await page.getByRole('button', { name: 'Can a non-active player play loot' }).click();
  await expectRefusal(page);
  await page.getByRole('button', { name: en('site.check.copy_link') }).click();
  const link = await page.getByRole('textbox', { name: en('site.check.link') }).inputValue();
  expect(link).toMatch(/#\/en\/check\/[A-Za-z0-9_-]+$/);
  await expect(page.getByRole('status')).toBeVisible();
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  await shot(page, 'site-check-phone');

  const other = await context.newPage();
  await other.goto(link);
  await expectRefusal(other);

  await other.getByRole('button', { name: 'Does the active player gain 1 cent' }).click();
  await expect(other.locator('.answer').getByRole('heading', { name: en('site.check.allowed') })).toBeVisible();
  await expect(other.locator('.answer')).toContainText(en('site.check.matches'));

  await other.getByRole('link', { name: en('site.nav.rules'), exact: true }).click();
  await expect(other.getByRole('heading', { name: en('site.rules.title'), level: 1 })).toBeVisible();
  expect(await other.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  await shot(other, 'site-rules-phone');
});
