// Settings screen (7.10): ST-01, ST-03, ST-06, ST-07, V-46.
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { realAssets, shot } from './shot';
import { pseudo } from '../../frontend/src/locale';
import { languageNames, locale, tr, trRe } from './locale';

// One dark theme (15.80, ST-07): nothing to pick, nothing on the page.
test('there is no theme to pick', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  await expect(page.getByRole('combobox', { name: tr('settings.mat') })).toBeVisible();
  await expect(page.getByLabel('Theme')).toHaveCount(0);
  await expect(page.locator('html')).not.toHaveAttribute('data-theme');
});

// Owner review of 15.80: Animations on its own line, then the game mat
// and the dice side by side.
test('animations sit above the game mat and the dice', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  const box = async (name: string) => {
    const b = await page.getByRole('combobox', { name }).boundingBox();
    if (!b) {
      throw new Error(`no ${name} field`);
    }
    return b;
  };
  const [animations, mat, dice] = [await box(tr('settings.animations')), await box(tr('settings.mat')), await box(tr('settings.dice'))];
  expect(animations.y + animations.height).toBeLessThan(mat.y);
  expect(Math.abs(mat.y - dice.y)).toBeLessThan(1);
  expect(mat.x).toBeLessThan(dice.x);
});

test('the first start asks for a nickname, once', async ({ page }) => {
  await realAssets(page);
  await page.goto('/?fake&nosplash&firstrun');
  await expect(page.getByRole('heading', { name: tr('welcome.title') })).toBeVisible();
  await shot(page, 'welcome');
  await page.getByLabel(tr('lobby.nickname')).fill('Ada');
  await page.getByRole('button', { name: tr('lobby.games.continue') }).click();
  await expect(page.getByRole('button', { name: tr('menu.play'), exact: true })).toBeVisible();
  await page.reload();
  await expect(page.getByRole('button', { name: tr('menu.play'), exact: true })).toBeVisible();
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  await expect(page.getByLabel(tr('settings.nickname'))).toHaveValue('Ada');
});

test('the question before ending the turn can be turned off (V-46)', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  const ask = page.getByLabel(tr('settings.end_turn_ask'));
  await expect(ask).toBeChecked();
  await ask.uncheck();
  await page.getByRole('button', { name: tr('settings.save') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('settings.saved'));
  await page.reload();
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  await expect(ask).not.toBeChecked();
});

test('fullscreen and mat apply at once and survive a restart', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  await page.getByLabel(tr('settings.fullscreen')).check();
  await page.getByRole('combobox', { name: tr('settings.mat') }).selectOption('basement');
  await expect(page.getByLabel(tr('settings.mat_preview'))).toHaveClass(/mat-basement/);
  await page.getByRole('button', { name: tr('settings.save') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('settings.saved'));
  expect(await page.evaluate(() => (window as unknown as { fakeFullscreen: boolean }).fakeFullscreen)).toBe(true);

  await page.reload();
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  await expect(page.getByLabel(tr('settings.fullscreen'))).toBeChecked();
  await expect(page.getByRole('combobox', { name: tr('settings.mat') })).toHaveValue('basement');
  await page.getByLabel(tr('settings.fullscreen')).uncheck();
  await page.getByRole('button', { name: tr('settings.save') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('settings.saved'));
  expect(await page.evaluate(() => (window as unknown as { fakeFullscreen: boolean }).fakeFullscreen)).toBe(false);
});

// Mat images (15.10, V-24): the setting puts the image under the table;
// without images, or with a missing one, the table shows colors.
test('a mat image lies under the table', async ({ browser }) => {
  const [mat] = await startGame(browser, ['Mat', 'Ned'], { url: `${appURL()}&mat=basement` });
  if (!mat) {
    throw new Error('no pages');
  }
  const table = mat.getByRole('main', { name: new RegExp(`^${trRe('game.turn').source}`) });
  await expect(table).toHaveClass(/mat-basement mat-image/);
  expect(await table.evaluate((el) => getComputedStyle(el).backgroundImage)).toContain('/cardimg/assets/mats/basement.webp');
  await shot(mat, 'table-mat');
});

test('the mat menu lists the mat images', async ({ page }) => {
  await realAssets(page);
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  const mats = page.getByRole('combobox', { name: tr('settings.mat') });
  await expect(mats.getByRole('option')).toHaveText([tr('mat.basement'), tr('mat.cathedral')]);
  await mats.selectOption('basement');
  await expect(page.getByLabel(tr('settings.mat_preview'))).toHaveClass(/mat-image/);
  await shot(page, 'settings-mat');
});

test('the first mat image is the default', async ({ page }) => {
  await page.goto(appURL());
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  await expect(page.getByRole('combobox', { name: tr('settings.mat') })).toHaveValue('basement');
  await expect(page.getByLabel(tr('settings.mat_preview'))).toHaveClass(/mat-basement mat-image/);
  // The image fails to load here; the basement colors stay under it.
  expect(await page.getByLabel(tr('settings.mat_preview')).evaluate((el) => getComputedStyle(el).backgroundImage)).toContain('rgb(90, 67, 48)');
});

test('without mat images the table shows colors', async ({ browser }) => {
  for (const [query, colors] of [
    ['&nomats', 'mat-basement'],
    ['&nomats&mat=cathedral', 'mat-cathedral'],
    // Green felt is gone (ST-08): a saved one is the default.
    ['&nomats&mat=felt', 'mat-basement'],
    ['&mat=dark_room', 'mat-basement'],
  ]) {
    const page = await (await browser.newContext()).newPage();
    await page.goto(`${appURL()}${query ?? ''}`);
    await page.getByRole('button', { name: tr('menu.settings') }).click();
    await expect(page.getByLabel(tr('settings.mat_preview'))).toHaveClass(new RegExp(`${colors ?? ''}$`));
    await expect(page.getByLabel(tr('settings.mat_preview'))).not.toHaveClass(/mat-image/);
  }
});

// Dice looks (15.64, V-19): the menu groups every d6 by its set, the
// preview shows the picked look's six faces, and the pick is saved; the
// next hello sends it (15.6).
test('the dice menu picks the look of your dice', async ({ page }) => {
  if (!process.env.SCREENSHOTS) {
    // Any small image does: the test checks which texture is asked for.
    const pixel = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=', 'base64');
    await page.route(/\/cardimg\/assets\/dices\//, (route) => route.fulfill({ contentType: 'image/png', body: pixel }));
  }
  await realAssets(page);
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  const dice = page.getByRole('combobox', { name: tr('settings.dice') });
  await expect(dice).toHaveValue('tear');
  expect(await dice.locator('optgroup').evaluateAll((gs) => gs.map((g) => g.getAttribute('label')))).toEqual([tr('settings.dice_set.holy'), tr('settings.dice_set.unholy')]);
  await expect(dice.getByRole('option')).toHaveCount(10);
  const faces = page.getByRole('group', { name: tr('settings.dice_preview') }).locator('.die img');
  await expect(faces).toHaveCount(6);
  await expect(faces.first()).toHaveAttribute('src', /d6_tear\.webp$/);

  await dice.selectOption('poop');
  await expect(faces.first()).toHaveAttribute('src', /unholy_set\/d6_poop\.webp$/);
  await shot(page, 'settings-dice');
  await page.getByRole('button', { name: tr('settings.save') }).click();
  await expect(page.getByRole('status')).toHaveText(tr('settings.saved'));
  await page.reload();
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  await expect(page.getByRole('combobox', { name: tr('settings.dice') })).toHaveValue('poop');
});

test('without dice images there is no dice menu', async ({ page }) => {
  await page.goto('/?fake&nosplash&nodice');
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  await expect(page.getByRole('combobox', { name: tr('settings.mat') })).toBeVisible();
  await expect(page.getByRole('combobox', { name: tr('settings.dice') })).toHaveCount(0);
  await expect(page.getByRole('group', { name: tr('settings.dice_preview') })).toHaveCount(0);
});

// Language (16.2, ADR 021): each in its own name; English and Russian
// for players, the test language in development builds. Russian has no
// interface file yet, so it shows English. It applies on Save.
test('the language applies at once and survives a restart', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  // Found by its choices, as its label changes with the language.
  const lang = page.getByRole('combobox').filter({ has: page.getByRole('option', { name: 'English' }) });
  const save = page.locator('form button[type=submit]');
  // A language file adds its choice (ADR 021, 11): ['English', 'Русский', 'qps (test)'] now.
  await expect(lang.getByRole('option')).toHaveText(languageNames);

  await lang.selectOption('qps');
  await save.click();
  await expect(page.getByRole('heading', { level: 1 })).toHaveText(pseudo('Settings'));
  await expect(page.getByRole('status')).toHaveText(pseudo('Saved.'));
  await expect(page.getByRole('combobox', { name: pseudo('Language') })).toHaveValue('qps');

  await lang.selectOption('ru');
  await save.click();
  // Russian has no interface file yet: English shows.
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Settings');
  await page.reload();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByRole('combobox', { name: 'Language' })).toHaveValue('ru');
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Settings');
});

test('the language follows the OS until one is picked', async ({ browser }) => {
  test.skip(locale !== 'en', 'LOCALE picks the language instead of the OS');
  for (const [os, want] of [
    ['ru-RU', 'ru'],
    ['de-DE', 'en'],
  ]) {
    const page = await (await browser.newContext({ locale: os })).newPage();
    await page.goto('/?fake&nosplash');
    await page.getByRole('button', { name: tr('menu.settings') }).click();
    await expect(page.getByRole('combobox', { name: tr('settings.language') })).toHaveValue(want ?? '');
  }
});
