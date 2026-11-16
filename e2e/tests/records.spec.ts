// The Records screen (10.1, M-01, UC-44): the fake bridge's records
// folder has an unfinished match, a finished one, one made before
// replays and one from a newer app. No real files.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { shot } from './shot';
import { tr, trRe } from './locale';

interface FakeWindow {
  fakeRecords: { path: string }[];
  fakeExported?: string;
  fakePick?: string;
}

// openRecords opens the Records screen from the main menu; before runs
// first, e.g. to empty the fake records folder.
async function openRecords(page: Page, before?: () => void): Promise<void> {
  await page.goto('/?fake&nosplash');
  // The menu shows once the fake bridge has set up its records.
  const records = page.getByRole('button', { name: tr('menu.records'), exact: true });
  await expect(records).toBeVisible();
  if (before) {
    await page.evaluate(before);
  }
  await records.click();
  await expect(page.getByRole('heading', { name: tr('records.heading') })).toBeVisible();
}

// row is the listed record whose text has text, e.g. its players.
function row(page: Page, text: string) {
  return page.getByRole('list', { name: tr('records.list_label') }).getByRole('listitem').filter({ hasText: text });
}

test('the menu opens the list; records that cannot play say why', async ({ page }) => {
  await openRecords(page);
  await expect(page.getByRole('list', { name: tr('records.list_label') }).getByRole('listitem')).toHaveCount(4);

  // Newest first: the unfinished match leads.
  const going = page.getByRole('list', { name: tr('records.list_label') }).getByRole('listitem').first();
  await expect(going).toContainText(tr('records.unfinished'));
  await expect(going).toContainText(tr('records.steps', { n: 12 }));
  const won = row(page, 'Ann, Bo, Cy');
  await expect(won).toContainText(tr('records.won', { players: 'Bo' }));
  await expect(won.getByRole('button', { name: trRe('records.watch_label') })).toBeVisible();

  const old = row(page, tr('error.record_before_replays'));
  await expect(old).toHaveCount(1);
  await expect(old.getByRole('button', { name: trRe('records.watch_label') })).toHaveCount(0);
  await expect(old.getByRole('button', { name: trRe('records.export_label') })).toBeVisible();

  const newer = row(page, 'from-the-future.fsrec');
  await expect(newer).toContainText(tr('error.record_newer'));
  await expect(newer.getByRole('button', { name: trRe('records.watch_label') })).toHaveCount(0);
  await shot(page, 'records');

  await page.getByRole('button', { name: tr('records.back') }).click();
  await expect(page.getByRole('button', { name: tr('menu.records'), exact: true })).toBeVisible();
});

test('Watch opens a record and sums it up; Back returns to the list', async ({ page }) => {
  await openRecords(page);
  await row(page, 'Ann, Bo, Cy').getByRole('button', { name: trRe('records.watch_label') }).click();
  await expect(page.getByText(tr('records.players', { players: 'Ann, Bo, Cy' }))).toBeVisible();
  await expect(page.getByText(tr('records.won', { players: 'Bo' }))).toBeVisible();
  await expect(page.getByText(`${tr('records.turns', { n: 84 })} · ${tr('records.steps', { n: 840 })}`)).toBeVisible();
  await shot(page, 'records-opened');

  await page.getByRole('button', { name: tr('records.close') }).click();
  await expect(page.getByRole('list', { name: tr('records.list_label') })).toBeVisible();
  // Esc goes back too: from the list to the menu.
  await page.keyboard.press('Escape');
  await expect(page.getByRole('button', { name: tr('menu.records'), exact: true })).toBeVisible();
});

test('Delete asks first; Keep keeps the record', async ({ page }) => {
  await openRecords(page);
  const old = row(page, tr('error.record_before_replays'));
  await old.getByRole('button', { name: trRe('records.delete_label') }).click();
  const ask = old.getByRole('group', { name: tr('records.delete_ask') });
  await expect(ask).toBeVisible();
  await shot(page, 'records-delete');
  await ask.getByRole('button', { name: tr('records.delete_no') }).click();
  await expect(ask).toHaveCount(0);
  await expect(old).toHaveCount(1);

  await old.getByRole('button', { name: trRe('records.delete_label') }).click();
  await old.getByRole('button', { name: tr('records.delete_yes'), exact: true }).click();
  await expect(page.getByRole('status')).toHaveText(tr('records.deleted'));
  await expect(row(page, tr('error.record_before_replays'))).toHaveCount(0);
  await expect(page.getByRole('list', { name: tr('records.list_label') }).getByRole('listitem')).toHaveCount(3);
  expect(await page.evaluate(() => (window as unknown as FakeWindow).fakeRecords.length)).toBe(3);
});

test('Export saves a copy where the player picks', async ({ page }) => {
  await openRecords(page);
  await row(page, 'Ann, Bo, Cy').getByRole('button', { name: trRe('records.export_label') }).click();
  const path = '/fake/records/2026-11-15-180000-g2_Ann_Bo_Cy.fsrec';
  await expect(page.getByRole('status')).toHaveText(tr('records.exported', { path: '/fake/export/2026-11-15-180000-g2_Ann_Bo_Cy.fsrec' }));
  expect(await page.evaluate(() => (window as unknown as FakeWindow).fakeExported)).toBe(path);
});

test('Open a file opens a record from anywhere; a cancel stays', async ({ page }) => {
  await openRecords(page);
  const open = page.getByRole('button', { name: tr('records.open_file') });
  await open.click();
  await expect(page.getByRole('list', { name: tr('records.list_label') })).toBeVisible();

  await page.evaluate(() => {
    (window as unknown as FakeWindow).fakePick = '/elsewhere/match.fsrec';
  });
  await open.click();
  await expect(page.getByText(tr('records.players', { players: 'Dee, Eve' }))).toBeVisible();
  await expect(page.getByText(tr('records.won', { players: 'Eve' }))).toBeVisible();
});

test('an empty records folder says so', async ({ page }) => {
  await openRecords(page, () => {
    (window as unknown as FakeWindow).fakeRecords = [];
  });
  await expect(page.getByText(tr('records.none'))).toBeVisible();
  await expect(page.getByRole('list', { name: tr('records.list_label') })).toHaveCount(0);
  await shot(page, 'records-empty');
});
