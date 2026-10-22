// Settings screen (7.10): ST-01, ST-03, ST-06, ST-07, V-46.
import { expect, test } from './fixtures';
import { startGame } from './game';
import { appURL } from './server';
import { realAssets, shot } from './shot';

// One dark theme (15.80, ST-07): nothing to pick, nothing on the page.
test('there is no theme to pick', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByRole('combobox', { name: 'Game mat' })).toBeVisible();
  await expect(page.getByLabel('Theme')).toHaveCount(0);
  await expect(page.locator('html')).not.toHaveAttribute('data-theme');
});

// Owner review of 15.80: Animations on its own line, then the game mat
// and the dice side by side.
test('animations sit above the game mat and the dice', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  const box = async (name: string) => {
    const b = await page.getByRole('combobox', { name }).boundingBox();
    if (!b) {
      throw new Error(`no ${name} field`);
    }
    return b;
  };
  const [animations, mat, dice] = [await box('Animations'), await box('Game mat'), await box('Dice')];
  expect(animations.y + animations.height).toBeLessThan(mat.y);
  expect(Math.abs(mat.y - dice.y)).toBeLessThan(1);
  expect(mat.x).toBeLessThan(dice.x);
});

test('the first start asks for a nickname, once', async ({ page }) => {
  await realAssets(page);
  await page.goto('/?fake&nosplash&firstrun');
  await expect(page.getByRole('heading', { name: 'Welcome!' })).toBeVisible();
  await shot(page, 'welcome');
  await page.getByLabel('Nickname').fill('Ada');
  await page.getByRole('button', { name: 'Continue' }).click();
  await expect(page.getByRole('button', { name: 'Play', exact: true })).toBeVisible();
  await page.reload();
  await expect(page.getByRole('button', { name: 'Play', exact: true })).toBeVisible();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByLabel('Nickname')).toHaveValue('Ada');
});

test('the question before ending the turn can be turned off (V-46)', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  const ask = page.getByLabel('Ask before ending the turn while plays are left');
  await expect(ask).toBeChecked();
  await ask.uncheck();
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  await page.reload();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(ask).not.toBeChecked();
});

test('fullscreen and mat apply at once and survive a restart', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel('Fullscreen').check();
  await page.getByRole('combobox', { name: 'Game mat' }).selectOption('basement');
  await expect(page.getByLabel('Game mat preview')).toHaveClass(/mat-basement/);
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  expect(await page.evaluate(() => (window as unknown as { fakeFullscreen: boolean }).fakeFullscreen)).toBe(true);

  await page.reload();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByLabel('Fullscreen')).toBeChecked();
  await expect(page.getByRole('combobox', { name: 'Game mat' })).toHaveValue('basement');
  await page.getByLabel('Fullscreen').uncheck();
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  expect(await page.evaluate(() => (window as unknown as { fakeFullscreen: boolean }).fakeFullscreen)).toBe(false);
});

// Mat images (15.10, V-24): the setting puts the image under the table;
// without images, or with a missing one, the table shows colors.
test('a mat image lies under the table', async ({ browser }) => {
  const [mat] = await startGame(browser, ['Mat', 'Ned'], { url: `${appURL()}&mat=basement` });
  if (!mat) {
    throw new Error('no pages');
  }
  const table = mat.getByRole('main', { name: /^Turn / });
  await expect(table).toHaveClass(/mat-basement mat-image/);
  expect(await table.evaluate((el) => getComputedStyle(el).backgroundImage)).toContain('/cardimg/assets/mats/basement.webp');
  await shot(mat, 'table-mat');
});

test('the mat menu lists the mat images', async ({ page }) => {
  await realAssets(page);
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  const mats = page.getByRole('combobox', { name: 'Game mat' });
  await expect(mats.getByRole('option')).toHaveText(['Basement', 'Cathedral']);
  await mats.selectOption('basement');
  await expect(page.getByLabel('Game mat preview')).toHaveClass(/mat-image/);
  await shot(page, 'settings-mat');
});

test('the first mat image is the default', async ({ page }) => {
  await page.goto(appURL());
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByRole('combobox', { name: 'Game mat' })).toHaveValue('basement');
  await expect(page.getByLabel('Game mat preview')).toHaveClass(/mat-basement mat-image/);
  // The image fails to load here; the basement colors stay under it.
  expect(await page.getByLabel('Game mat preview').evaluate((el) => getComputedStyle(el).backgroundImage)).toContain('rgb(90, 67, 48)');
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
    await page.getByRole('button', { name: 'Settings' }).click();
    await expect(page.getByLabel('Game mat preview')).toHaveClass(new RegExp(`${colors ?? ''}$`));
    await expect(page.getByLabel('Game mat preview')).not.toHaveClass(/mat-image/);
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
  await page.getByRole('button', { name: 'Settings' }).click();
  const dice = page.getByRole('combobox', { name: 'Dice' });
  await expect(dice).toHaveValue('tear');
  expect(await dice.locator('optgroup').evaluateAll((gs) => gs.map((g) => g.getAttribute('label')))).toEqual(['Holy', 'Unholy']);
  await expect(dice.getByRole('option')).toHaveCount(10);
  const faces = page.getByRole('group', { name: 'Dice preview' }).locator('.die img');
  await expect(faces).toHaveCount(6);
  await expect(faces.first()).toHaveAttribute('src', /d6_tear\.webp$/);

  await dice.selectOption('poop');
  await expect(faces.first()).toHaveAttribute('src', /unholy_set\/d6_poop\.webp$/);
  await shot(page, 'settings-dice');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  await page.reload();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByRole('combobox', { name: 'Dice' })).toHaveValue('poop');
});

test('without dice images there is no dice menu', async ({ page }) => {
  await page.goto('/?fake&nosplash&nodice');
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByRole('combobox', { name: 'Game mat' })).toBeVisible();
  await expect(page.getByRole('combobox', { name: 'Dice' })).toHaveCount(0);
  await expect(page.getByRole('group', { name: 'Dice preview' })).toHaveCount(0);
});
