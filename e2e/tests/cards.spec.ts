// Card component (7.4) in the dev-only card gallery.
import { expect, test } from './fixtures';
import { realAssets, shot } from './shot';
import { cardName, locale, tr, trIn } from './locale';

const gallery = '/?fake&nosplash&screen=cards';

test('cards show their image, or a readable blank card', async ({ page }) => {
  await realAssets(page);
  await page.goto(gallery);
  const hand = page.getByLabel('Hand size');
  // Isaac has an image.
  await expect(hand.getByRole('img', { name: cardName('isaac') }).locator('img')).toBeVisible();
  // The D6's image fails to load: it falls back to its text.
  await expect(hand.getByRole('img', { name: cardName('the_d6') })).toContainText('Eternal Treasure Card');
  // Big Spider has no image: name, stats and reward.
  const spider = hand.getByRole('img', { name: cardName('big_spider') });
  await expect(spider).toContainText('HP 3');
  await expect(spider).toContainText('1x Loot');
  await expect(hand.getByRole('img', { name: tr('deck.face_down.treasure') })).toBeVisible();
  await shot(page, 'cards');
});

test('a face-down card shows its deck back, or a plain back without it', async ({ page }) => {
  // Only the treasure back is there; the loot back fails to load. The
  // loot route goes after the fixture's real assets, so it wins (shot.ts).
  const svg = '<svg xmlns="http://www.w3.org/2000/svg" width="63" height="88"><rect width="63" height="88" fill="#a70"/></svg>';
  await page.route('**/cardimg/assets/backs/treasure.webp', (route) => route.fulfill({ body: svg, contentType: 'image/svg+xml' }));
  await page.route('**/cardimg/assets/backs/loot.webp', (route) => route.abort());
  await page.goto(gallery);
  const back = page.getByLabel('Hand size').getByRole('img', { name: tr('deck.face_down.treasure') });
  await expect(back.locator('img')).toHaveAttribute('src', '/cardimg/assets/backs/treasure.webp');
  await expect(back.locator('img')).toBeVisible();
  await expect(back).not.toContainText(tr('deck.back.treasure'));
  await expect(page.getByLabel('Hand size').getByRole('img', { name: tr('deck.face_down.loot') })).toHaveText(tr('deck.back.loot'));
});

test('hovering a card zooms it with its full text', async ({ page }) => {
  await page.goto(gallery);
  const zoom = page.getByTestId('card-zoom');
  await page.getByLabel('Table size').getByRole('img', { name: cardName('big_spider') }).hover();
  await expect(zoom).toContainText('When this dies, the active player may attack');
  await shot(page, 'card-zoom');
  await page.getByRole('heading', { name: 'Cards' }).hover();
  await expect(zoom).toBeHidden();
});

test('a long press zooms a card on touch; a tap closes it', async ({ page }) => {
  await page.goto(gallery);
  const zoom = page.getByTestId('card-zoom');
  const d6 = page.getByLabel('Table size').getByRole('img', { name: cardName('the_d6') });
  await d6.dispatchEvent('pointerdown', { pointerType: 'touch', clientX: 100, clientY: 100 });
  await expect(zoom).toContainText("This can't be destroyed");
  await zoom.dispatchEvent('pointerdown', { pointerType: 'touch' });
  await expect(zoom).toBeHidden();
  // A short tap does not zoom.
  await d6.dispatchEvent('pointerdown', { pointerType: 'touch', clientX: 100, clientY: 100 });
  await d6.dispatchEvent('pointerup', { pointerType: 'touch' });
  await page.waitForTimeout(600);
  await expect(zoom).toBeHidden();
});

// Russian names come from card_db (16.5, LO-02); the interface stays
// English. Big Spider has no image, so its text card shows the name;
// the zoom credits the translators (CD-06).
test('with Russian picked, cards show their Russian names', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: tr('menu.settings') }).click();
  await page.getByRole('combobox').filter({ has: page.getByRole('option', { name: 'English' }) }).selectOption('ru');
  await page.locator('form button[type=submit]').click();
  await page.goto(gallery);
  const table = page.getByLabel('Table size');
  const spider = table.getByRole('img', { name: 'Большой Паук' });
  await expect(spider).toContainText('Большой Паук');
  await expect(table.getByRole('img', { name: 'Исаак' })).toBeVisible();
  await spider.hover();
  const zoom = page.getByTestId('card-zoom');
  await expect(zoom).toContainText(trIn('ru', 'zoom.translators', { names: 'Alchba, Mitboy, d0m1n1c132, Blacknenko, N-Chance' }));
  await shot(page, 'card-zoom-ru');
  // A card with no Russian name keeps its English one.
  await expect(table.getByRole('img', { name: 'Soul of Greed' })).toBeVisible();
});

// Card images per language (16.6): the fake Isaac has a Russian test
// image; with Russian picked it shows that one, with English the other.
// Each language is picked first: pages start in the tests' LOCALE.
test('with Russian picked, a card shows its Russian image', async ({ page }) => {
  const isaacArt = (lang: string) => page.getByLabel('Hand size').getByRole('img', { name: lang === 'ru' ? 'Исаак' : 'Isaac' }).locator('img');
  let shown = locale;
  const pick = async (lang: string) => {
    await page.goto('/?fake&nosplash');
    await page.getByRole('button', { name: trIn(shown, 'menu.settings') }).click();
    shown = lang;
    await page.getByRole('combobox').filter({ has: page.getByRole('option', { name: 'English' }) }).selectOption(lang);
    await page.locator('form button[type=submit]').click();
    await page.goto(gallery);
  };
  await pick('en');
  await expect(isaacArt('en')).toHaveAttribute('src', /Isaac%20art/);
  await pick('ru');
  await expect(isaacArt('ru')).toHaveAttribute('src', /Isaac%20ru/);
  await expect(isaacArt('ru')).toBeVisible();
});
