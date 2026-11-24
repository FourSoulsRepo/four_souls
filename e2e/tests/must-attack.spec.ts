// "Must attack" on the attack button (18.35, ADR 037 item 7, UC-12): a
// player who owes an attack sees the attack icon red, its tooltip naming
// the target; the monster keeps its dashed red frame. Saved games give
// the cards, so nothing depends on the deal.
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { repeatNames, startGame } from './game';
import { cardName, tr } from './locale';
import { editSave, saveAndContinue, saveEditContinue } from './saves';
import { appURL } from './server';
import { shot } from './shot';
import { activePage, canEndTurn, click, Intent, Purpose, seen, settle, turnAction, untilPick } from './turn';

test.setTimeout(150_000);

// The engine's monster deck (state.go).
const MONSTER_DECK = 2;

// Owed is what the page's view says the turn owes, and the first slot's
// monster.
interface Owed {
  mustAttack: number;
  mustAttacks: number;
  mustAttackDeck: number;
  attackOpen: boolean;
  monster: number;
}

async function owed(page: Page): Promise<Owed> {
  return page.evaluate(() => {
    interface S {
      game: {
        view: {
          turn: { must_attack?: number; must_attacks?: number; must_attack_deck?: number };
          attack?: unknown;
          monsters: { top?: { id: number } }[];
        };
      } | null;
    }
    const v = (window as unknown as { fourSouls: { store: { getState: () => S } } }).fourSouls.store.getState().game?.view;
    return {
      mustAttack: v?.turn.must_attack ?? 0,
      mustAttacks: v?.turn.must_attacks ?? 0,
      mustAttackDeck: v?.turn.must_attack_deck ?? 0,
      attackOpen: v?.attack !== undefined,
      monster: v?.monsters[0]?.top?.id ?? 0,
    };
  });
}

// free: the page may act in its turn's action step: the stack empty, no
// attack open, and Attack allowed. Then the duty is the one to show.
async function free(page: Page): Promise<boolean> {
  const s = await seen(page);
  return s.ready && s.stack.length === 0 && s.allowed.includes(Intent.attack) && !(await owed(page)).attackOpen;
}

// tip is the attack button's tooltip with the owed attack's line.
const tip = (line: string) => `${tr('turn.attack')} (A) · ${line}`;

test('Monster Manual: the attack button is red and names the monster (18.35)', async ({ browser }) => {
  const names = repeatNames(['Abe', 'Bea']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  // Monster Manual in play; Fatty, a plain monster, in the first slot.
  await saveAndContinue(pages, names.join(', '), 'bomb', 'monster_manual', false, {}, undefined, 'fatty');
  const active = await activePage(pages, click);
  const other = pages.find((p) => p !== active);
  if (!other) {
    throw new Error('no other page');
  }
  const fatty = (await owed(active)).monster;
  await active.mouse.move(1, 1);
  await shot(active, 'must-attack-none');
  // No duty yet: the icon is not red.
  await expect(turnAction(active, 'Attack')).not.toHaveClass(/\bmust\b/);

  // Monster Manual: "Choose a monster. The active player must attack
  // that monster this turn if able." Its one ability aims (15.25).
  const manual = active.getByRole('group', { name: tr('table.your_items') }).getByRole('button', { name: cardName('monster_manual'), exact: true });
  await click(active, manual);
  await click(active, active.getByRole('group', { name: tr('board.monsters') }).getByLabel(cardName('fatty'), { exact: true }));
  // Wait for the exact state (5.16.10): Fatty owed, the ability resolved,
  // the turn back in its action step.
  await settle(pages, click, async () => (await owed(active)).mustAttack === fatty && (await free(active)));
  const button = turnAction(active, 'Attack');
  await expect(button).toHaveClass(/\bmust\b/);
  // The card names the monster and no count (Monster Manual gives no
  // extra attack: its situation "gives-no-attack").
  await expect(button).toHaveAttribute('title', tip(tr('banner.must_attack', { target: cardName('fatty') })));
  // The monster keeps its dashed red frame, for everyone (19.180).
  for (const page of pages) {
    await expect(page.locator(`.card[data-zone="obj-${String(fatty)}"]`)).toHaveClass(/\bmust-attack\b/);
  }
  // The other player's icon stays as it was: the duty is not theirs.
  await expect(turnAction(other, 'Attack')).not.toHaveClass(/\bmust\b/);
  await active.mouse.move(1, 1);
  await shot(active, 'must-attack-manual');
  await button.hover();
  await shot(active, 'must-attack-manual-hover');
});

test('Ambush!: the attack button is red and names the monster deck (18.35)', async ({ browser }) => {
  const names = repeatNames(['Cy', 'Dee']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  await activePage(pages, click);
  // Ambush! on top of the monster deck; a copy dealt elsewhere takes
  // the top card's place (one of a kind).
  await saveEditContinue(pages, names.join(', '), (game) => {
    editSave(game, (save) => {
      const { state } = save as unknown as { state: { decks: number[][]; objects: { card: string }[] } };
      const top = state.decks[MONSTER_DECK]?.at(-1);
      const obj = top === undefined ? undefined : state.objects[top];
      if (!obj) {
        throw new Error('the monster deck is empty');
      }
      for (const o of state.objects) {
        if (o.card === 'ambush') {
          o.card = obj.card;
        }
      }
      obj.card = 'ambush';
    });
  });
  const active = await activePage(pages, click);
  // Attack the monster deck: it reveals Ambush!.
  await click(active, turnAction(active, 'Attack'));
  await untilPick(pages, click, active, Purpose.attackTarget);
  await click(active, active.locator('[data-zone="deck-monster"] .card.hl-pick'));
  // Ambush!: "The active player must attack the monster deck 2 times
  // this turn." Wait for the exact state (5.16.10): 2 deck attacks owed,
  // the revealing attack over, the turn back in its action step.
  await settle(pages, click, async () => {
    const o = await owed(active);
    return o.mustAttacks === 2 && o.mustAttackDeck === 2 && (await free(active));
  });
  const button = turnAction(active, 'Attack');
  await expect(button).toHaveClass(/\bmust\b/);
  await expect(button).toHaveAttribute('title', tip(tr('banner.must_attack_left', { target: tr('banner.monster_deck'), n: 2 })));
  for (const page of pages) {
    await expect(page.locator('[data-zone="deck-monster"]')).toHaveClass(/\bmust-attack\b/);
  }
  await active.mouse.move(1, 1);
  await shot(active, 'must-attack-ambush');
  await button.hover();
  await shot(active, 'must-attack-ambush-hover');
});

// BossSave is the part of a save the Boss Rush! test changes.
interface BossSave {
  state: {
    objects: { card: string }[];
    monsters: { cards: number[] }[];
    turn: { attacks: number; must_attack_one_of_cards?: string[] };
  };
}

test('Boss Rush!: the attack button names both bosses, both framed (19.241)', async ({ browser }) => {
  const names = repeatNames(['Eve', 'Fin']);
  const pages = await startGame(browser, names, { url: `${appURL()}&anim=off` });
  // The turn can't end while the attack is owed, so activePage would
  // play it out: find the active page before the save.
  const active = await activePage(pages, click);
  // The turn as Boss Rush! leaves it with a boss in each of two slots:
  // "The active player must make an additional attack on one of them
  // this turn." One more attack, owed on Monstro or Widow (19.241).
  let bosses: number[] = [];
  await saveEditContinue(pages, names.join(', '), (game) => {
    editSave(game, (save) => {
      const { state } = save as unknown as BossSave;
      bosses = ['monstro', 'widow'].map((card, i) => {
        const top = state.monsters[i]?.cards.at(-1);
        const obj = top === undefined ? undefined : state.objects[top];
        if (top === undefined || !obj) {
          throw new Error(`monster slot ${String(i + 1)} is empty`);
        }
        obj.card = card;
        return top;
      });
      state.turn.attacks += 1;
      // The duty names the bosses' cards (owner, 2026-11-24).
      state.turn.must_attack_one_of_cards = ['monstro', 'widow'];
    });
  });
  const button = turnAction(active, 'Attack');
  await expect(button).toHaveClass(/\bmust\b/);
  // Either boss will do: the tooltip names both.
  const either = tr('banner.either', { a: cardName('monstro'), b: cardName('widow') });
  await expect(button).toHaveAttribute('title', tip(tr('banner.must_attack', { target: either })));
  // The turn can't end before it (R-ATK-03).
  expect(await canEndTurn(active)).toBe(false);
  // Both keep the dashed red frame, for everyone (19.180).
  for (const page of pages) {
    for (const boss of bosses) {
      await expect(page.locator(`.card[data-zone="obj-${String(boss)}"]`)).toHaveClass(/\bmust-attack\b/);
    }
  }
  await active.mouse.move(1, 1);
  await shot(active, 'must-attack-boss-rush');
  await button.hover();
  await shot(active, 'must-attack-boss-rush-hover');
});
