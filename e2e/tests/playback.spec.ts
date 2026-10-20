// Event playback and animations (7.6).
import type { Page } from '@playwright/test';
import { expect, test } from './fixtures';
import { startGame } from './game';
import { shot } from './shot';

// listenSounds records the file of every sound played from now on, in
// window.played (15.50). Each sound still plays.
async function listenSounds(page: Page): Promise<void> {
  await page.evaluate(() => {
    const w = window as unknown as { played: string[] };
    w.played = [];
    const play = HTMLMediaElement.prototype.play;
    HTMLMediaElement.prototype.play = function (this: HTMLMediaElement) {
      w.played.push(new URL(this.src).pathname);
      return play.call(this);
    };
  });
}

function played(page: Page): Promise<string[]> {
  return page.evaluate(() => (window as unknown as { played: string[] }).played);
}

test('the animation speed setting survives a restart', async ({ page }) => {
  await page.goto('/?fake&nosplash');
  await page.getByRole('button', { name: 'Settings' }).click();
  await page.getByLabel('Animations').selectOption('fast');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('status')).toHaveText('Saved.');
  await page.reload();
  await page.getByRole('button', { name: 'Settings' }).click();
  await expect(page.getByLabel('Animations')).toHaveValue('fast');
});

test('a fast burst of updates plays in order and ends on the last one', async ({ browser }) => {
  const [kai] = await startGame(browser, ['Kai', 'Lu']);
  if (!kai) {
    throw new Error('no pages');
  }
  // Mark one card: if updates redrew the table from scratch, the mark
  // would be gone and every card would pop in again.
  const character = kai.locator('.player-bottom .player-cards > .card').first();
  await character.evaluate((el) => el.setAttribute('data-mark', 'kept'));
  // Replay the current view 20 times as new turns, each with a dice roll
  // and a looted card, as fast as the page takes them.
  await kai.evaluate(() => {
    interface Store {
      getState: () => {
        game: { view: { turn: { number: number; active: number } } } | null;
        lobby: { table: { you: number } | null };
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const you = s.lobby.table?.you ?? 0;
    for (let n = 2; n <= 21; n++) {
      const game = structuredClone(s.game);
      if (!game) {
        throw new Error('no game');
      }
      game.view.turn.number = n;
      s.receive({
        type: 'update',
        data: {
          ...game,
          step: n,
          events: [
            { kind: 'dice_rolled', player: you, amount: (n % 6) + 1, text: 'roll' },
            { kind: 'looted', player: you, object: 1000 + n },
          ],
        },
      });
    }
  });
  await expect(kai.getByTestId('effect-dice')).toBeVisible();
  await kai.waitForTimeout(300); // let the dice settle for the picture
  await shot(kai, 'effect-dice');
  // A fly plays between throws only briefly in a fast burst: look for it
  // on every frame (15.62).
  await kai.waitForFunction(() => document.querySelector('[data-testid="effect-fly"]') !== null);
  await expect(kai.getByRole('main', { name: 'Turn 21' })).toBeVisible({ timeout: 20_000 });
  await expect(kai.locator('.effects')).toHaveCount(0);
  await expect(character).toHaveAttribute('data-mark', 'kept');
  await expect(character).toHaveCSS('opacity', '1');
});

test('a thrown die falls, bounces and lands on the engine\'s value (15.62)', async ({ browser }) => {
  const [mira] = await startGame(browser, ['Mira', 'Noor']);
  if (!mira) {
    throw new Error('no pages');
  }
  await listenSounds(mira);
  // Replay the view with one fresh attack roll of a 4.
  await mira.evaluate(() => {
    interface Store {
      getState: () => {
        game: { step: number } | null;
        lobby: { table: { you: number } | null };
        receive: (env: unknown) => void;
      };
    }
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    if (!game) {
      throw new Error('no game');
    }
    const you = s.lobby.table?.you ?? 0;
    s.receive({ type: 'update', data: { ...game, step: game.step + 1, events: [{ kind: 'dice_rolled', player: you, amount: 4, text: 'attack roll' }] } });
  });
  const thrown = mira.getByTestId('thrown-die');
  await expect(thrown).toBeAttached();
  // The die is a cube of six faces (15.69); only the cube is named.
  await expect(thrown.locator('.die-face')).toHaveCount(6);
  await expect(thrown.getByRole('img')).toHaveAttribute('aria-label', 'Rolled 4');
  // seek holds the throw at a part of its length: the fall, the cube's
  // turn, its shadow and the label.
  const effect = mira.getByTestId('effect-dice');
  const seek = (at: number) =>
    effect.evaluate((el, part) => {
      for (const a of el.getAnimations({ subtree: true })) {
        a.pause();
        a.currentTime = part * Number(a.effect?.getTiming().duration ?? 0);
      }
    }, at);
  // In the air: still falling, well above where it lands. The die says
  // the engine's value from the start.
  await seek(0.15);
  const air = await thrown.boundingBox();
  const middle = (mira.viewportSize()?.height ?? 0) / 2;
  expect(air && air.y + air.height / 2).toBeLessThan(middle - 100);
  await expect(thrown.locator('.die')).toHaveAttribute('aria-label', 'Rolled 4');
  await shot(mira, 'dice-in-air');
  // Landed: it lies still in the middle, with the roll's name below.
  await seek(0.75);
  const landed = await thrown.boundingBox();
  expect(landed && Math.abs(landed.y + landed.height / 2 - middle)).toBeLessThan(80);
  await expect(mira.getByTestId('effect-dice')).toContainText('attack roll');
  await expect(thrown.locator('.die')).toHaveAttribute('aria-label', 'Rolled 4');
  await shot(mira, 'dice-landed');
  // It clattered on the table as it landed (15.69). The clatter keeps
  // real time, not the paused throw: wait for it before the throw ends,
  // or the ended effect takes its timer along (19.34).
  await expect.poll(() => played(mira)).toEqual(['/cardimg/assets/sounds/dice-roll.mp3']);
  await effect.evaluate((el) => {
    for (const a of el.getAnimations({ subtree: true })) {
      a.play();
    }
  });
  await expect(mira.locator('.effects')).toHaveCount(0);
  expect(await played(mira)).toEqual(['/cardimg/assets/sounds/dice-roll.mp3']);
});

test('an attack lunges at the monster and drops its HP at the hit', async ({ browser }) => {
  const [vee] = await startGame(browser, ['Vee', 'Wyn']);
  if (!vee) {
    throw new Error('no pages');
  }
  await listenSounds(vee);
  // Replay the view with one hit of the own character on the first
  // monster, whose HP is one lower after it.
  const ids = await vee.evaluate(() => {
    interface Card {
      id: number;
      stats?: { hp: number };
    }
    interface Store {
      getState: () => {
        game: { view: { players: { id: number; character: Card }[]; monsters: { top?: Card }[] } } | null;
        lobby: { table: { you: number } | null };
        receive: (env: unknown) => void;
      };
    }
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table?.you ?? 0;
    const character = game?.view.players.find((p) => p.id === you)?.character;
    const monster = game?.view.monsters.find((m) => m.top?.stats)?.top;
    if (!game || !character || !monster?.stats) {
      throw new Error('no game');
    }
    const hp = monster.stats.hp;
    monster.stats.hp = hp - 1;
    s.receive({
      type: 'update',
      data: {
        ...game,
        step: 2,
        events: [{ kind: 'damaged', player: -1, object: monster.id, amount: 1, source: character.id, text: 'combat' }],
      },
    });
    return { character: character.id, monster: monster.id, hp };
  });
  const character = vee.locator(`[data-zone="obj-${String(ids.character)}"]`);
  const monster = vee.locator(`[data-zone="obj-${String(ids.monster)}"]`);
  // The card's copy flies; the card itself waits hidden at home.
  await expect(vee.getByTestId('lunge')).toBeVisible();
  await expect(character).toBeHidden();
  if (process.env.SCREENSHOTS) {
    // Hold the copy on its way into the monster for the picture.
    const lunge = vee.getByTestId('lunge');
    await lunge.evaluate((el) => {
      for (const a of el.getAnimations()) {
        a.pause();
        a.currentTime = 0.42 * Number(a.effect?.getTiming().duration ?? 0);
      }
    });
    await shot(vee, 'attack-lunge');
    await lunge.evaluate((el) => {
      for (const a of el.getAnimations()) {
        a.play();
      }
    });
  }
  await expect(vee.getByTestId('burst')).toHaveText('−1');
  // The hit lands with the attack hit sound.
  expect(await played(vee)).toEqual(['/cardimg/assets/sounds/attack-hit.mp3']);
  await vee.waitForTimeout(200); // the burst pops
  await shot(vee, 'attack-hit');
  await expect(monster.getByLabel(`${String(ids.hp - 1)} health`)).toBeVisible();
  await expect(vee.locator('.effects-trail > *')).toHaveCount(0);
  await expect(character).toBeVisible();
});

// The view parts the replays below change.
interface Card {
  id: number;
  card: string;
  damage?: number;
  stats?: { hp: number };
}
interface View {
  players: { id: number; hp: number; character: Card }[];
  monsters: { top?: Card }[];
}
interface Store {
  getState: () => {
    game: { view: View } | null;
    lobby: { table: { you: number } | null };
    receive: (env: unknown) => void;
  };
}

test('damage from a card hits every target at once and reddens their HP', async ({ browser }) => {
  const [mo] = await startGame(browser, ['Mo', 'Nia']);
  if (!mo) {
    throw new Error('no pages');
  }
  await listenSounds(mo);
  // The own character's ability deals 1 damage to two monsters and to
  // its own player.
  const ids = await mo.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table?.you ?? 0;
    const me = game?.view.players.find((p) => p.id === you);
    const tops = (game?.view.monsters ?? []).flatMap((m) => (m.top?.stats ? [m.top] : [])).slice(0, 2);
    if (!game || !me || tops.length < 2) {
      throw new Error('no game');
    }
    const hps = tops.map((t) => t.stats?.hp ?? 0);
    const events: unknown[] = tops.map((t) => {
      if (t.stats) {
        t.stats.hp -= 1;
      }
      t.damage = 1;
      return { kind: 'damaged', player: -1, object: t.id, amount: 1, source: me.character.id };
    });
    me.hp -= 1;
    me.character.damage = 1;
    events.push({ kind: 'damaged', player: you, amount: 1, source: me.character.id });
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
    return { character: me.character.id, monsters: tops.map((t) => t.id), hps };
  });
  // A shot flies from the source to each target.
  await expect(mo.getByTestId('shot')).toHaveCount(3);
  if (process.env.SCREENSHOTS) {
    // Hold the shots on their way for the picture.
    await mo.evaluate(() => {
      for (const a of document.getAnimations()) {
        a.pause();
        a.currentTime = 0.6 * Number(a.effect?.getTiming().duration ?? 0);
      }
    });
    await shot(mo, 'damage-shots');
    await mo.evaluate(() => {
      for (const a of document.getAnimations()) {
        a.play();
      }
    });
  }
  // Every target bursts at once.
  await expect(mo.getByTestId('burst')).toHaveCount(3);
  // One fireball for all shots: the cast as they leave, the impact as
  // they land.
  expect(await played(mo)).toEqual(['/cardimg/assets/sounds/effect-cast.mp3', '/cardimg/assets/sounds/effect-impact.mp3']);
  await mo.waitForTimeout(200); // the bursts pop
  await shot(mo, 'damage-bursts');
  for (const [i, id] of ids.monsters.entries()) {
    const health = mo.locator(`[data-zone="obj-${String(id)}"]`).getByLabel(`${String((ids.hps[i] ?? 0) - 1)} health`);
    await expect(health).toHaveClass(/hurt/);
  }
  await expect(mo.locator(`[data-zone="obj-${String(ids.character)}"] .card-hp`)).toHaveClass(/hurt/);
  // The red health is the damage mark; no corner mark.
  await expect(mo.locator(`[data-zone="obj-${String(ids.character)}"] .card-damage`)).toHaveCount(0);
  await expect(mo.locator('.effects-trail > *')).toHaveCount(0);
  await shot(mo, 'damage-red-hp');
});

test('a dying monster gets a skull before it leaves', async ({ browser }) => {
  const [ola] = await startGame(browser, ['Ola', 'Pim']);
  if (!ola) {
    throw new Error('no pages');
  }
  const id = await ola.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const slot = game?.view.monsters.find((m) => m.top?.stats);
    const top = slot?.top;
    if (!game || !slot || !top) {
      throw new Error('no game');
    }
    slot.top = undefined;
    const events = [{ kind: 'died', player: -1, object: 999, card: top.card, prev: top.id }];
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
    return top.id;
  });
  const monster = ola.locator(`[data-zone="obj-${String(id)}"]`);
  await expect(ola.getByTestId('skull')).toBeVisible();
  await expect(monster).toBeVisible();
  await ola.waitForTimeout(250); // the skull pops
  await shot(ola, 'death-skull');
  await expect(monster).toHaveCount(0);
  await expect(ola.locator('.effects-trail > *')).toHaveCount(0);
});

test('a deck attack places the monster in its slot before the attack roll', async ({ browser }) => {
  const [dax] = await startGame(browser, ['Dax', 'Ivo']);
  if (!dax) {
    throw new Error('no pages');
  }
  // Replay the slot pick of a deck attack: the revealed card covers the
  // first slot, the attack starts, and the first roll is made.
  const ids = await dax.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table?.you ?? 0;
    const [first, second] = game?.view.monsters ?? [];
    const old = first?.top;
    if (!game || !first || !old || !second?.top) {
      throw new Error('no game');
    }
    const top = { ...second.top, id: 990 };
    first.top = top;
    const events = [
      { kind: 'card_revealed', player: -1, object: top.id, card: top.card },
      { kind: 'attack_target', player: you, object: top.id, card: top.card },
      { kind: 'dice_rolled', player: you, amount: 4, text: 'attack roll' },
    ];
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
    return { old: old.id, top: top.id };
  });
  // The card flies to the slot and stays there while the dice rolls.
  await expect(dax.getByTestId('effect-dice')).toBeVisible({ timeout: 5000 });
  await dax.waitForTimeout(300); // let the dice settle for the picture
  await shot(dax, 'deck-attack-roll');
  const placed = dax.getByTestId('placed');
  await expect(placed).toBeVisible();
  const slot = await dax.locator(`[data-zone="obj-${String(ids.old)}"]`).boundingBox();
  const box = await placed.boundingBox();
  expect(Math.abs((box?.x ?? 0) - (slot?.x ?? 0))).toBeLessThan(2);
  expect(Math.abs((box?.y ?? 0) - (slot?.y ?? 0))).toBeLessThan(2);
  // Then the board shows the card in its slot, and the copy is gone.
  await expect(dax.locator(`[data-zone="obj-${String(ids.top)}"]`)).toBeVisible();
  await expect(placed).toHaveCount(0);
});

test('a deck attack shows the revealed monster large once (19.4)', async ({ browser }) => {
  const [vix] = await startGame(browser, ['Vix', 'Wen']);
  if (!vix) {
    throw new Error('no pages');
  }
  // Replay a deck attack: the top card is revealed, then it covers the
  // first slot and the first roll is made. Every effect is recorded.
  await vix.evaluate(() => {
    interface Effects {
      subscribe: (f: (s: { playing: { key: number; effect: { kind: string } } | null }) => void) => void;
    }
    const w = window as unknown as { fourSouls: { store: Store & Effects }; effects: string[] };
    const store = w.fourSouls.store;
    w.effects = [];
    let last = 0;
    store.subscribe((s) => {
      if (s.playing && s.playing.key !== last) {
        last = s.playing.key;
        w.effects.push(s.playing.effect.kind);
      }
    });
    const s = store.getState();
    const you = s.lobby.table?.you ?? 0;
    const drawn = structuredClone(s.game);
    const placed = structuredClone(s.game);
    const [first, second] = placed?.view.monsters ?? [];
    const card = second?.top?.card;
    if (!drawn || !placed || !first || !second?.top || !card) {
      throw new Error('no game');
    }
    first.top = { ...second.top, id: 990 };
    s.receive({
      type: 'update',
      data: { ...drawn, step: 2, events: [{ kind: 'card_revealed', player: you, object: 989, card }] },
    });
    const events = [
      { kind: 'card_revealed', player: -1, object: 990, card },
      { kind: 'attack_target', player: you, object: 990, card },
      { kind: 'dice_rolled', player: you, amount: 4, text: 'attack roll' },
    ];
    s.receive({ type: 'update', data: { ...placed, step: 3, events } });
  });
  // The drawn card shows large; in its slot it only flies there.
  await expect(vix.getByTestId('effect-dice')).toBeVisible({ timeout: 8000 });
  await vix.waitForTimeout(300); // let the dice settle for the picture
  await shot(vix, 'deck-attack-once');
  await expect(vix.locator('[data-zone="obj-990"]')).toBeVisible();
  const effects = await vix.evaluate(() => (window as unknown as { effects: string[] }).effects);
  expect(effects).toEqual(['reveal', 'place', 'dice']);
});

test('a deck attack monster does not pop in after it lands (19.6)', async ({ browser }) => {
  const [yara] = await startGame(browser, ['Yara', 'Bram']);
  if (!yara) {
    throw new Error('no pages');
  }
  // Replay the slot pick of a deck attack, as above. When the board
  // draws the card in its slot, record how see-through it is then.
  await yara.evaluate(() => {
    const w = window as unknown as { fourSouls: { store: Store }; landed: Promise<number> };
    w.landed = new Promise((resolve) => {
      const seen = new MutationObserver(() => {
        const el = document.querySelector('[data-zone="obj-990"]');
        if (el) {
          seen.disconnect();
          requestAnimationFrame(() => {
            resolve(Number(getComputedStyle(el).opacity));
          });
        }
      });
      seen.observe(document.body, { childList: true, subtree: true });
    });
    const s = w.fourSouls.store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table?.you ?? 0;
    const [first, second] = game?.view.monsters ?? [];
    if (!game || !first || !second?.top) {
      throw new Error('no game');
    }
    first.top = { ...second.top, id: 990 };
    const events = [
      { kind: 'card_revealed', player: -1, object: 990, card: second.top.card },
      { kind: 'attack_target', player: you, object: 990, card: second.top.card },
      { kind: 'dice_rolled', player: you, amount: 4, text: 'attack roll' },
    ];
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
  });
  const card = yara.locator('[data-zone="obj-990"]');
  await expect(card).toBeVisible({ timeout: 8000 });
  await yara.waitForTimeout(60); // a pop would be half way now
  await shot(yara, 'deck-attack-landed');
  // The card that flew there is already whole: no second arrival.
  const opacity = await yara.evaluate(() => (window as unknown as { landed: Promise<number> }).landed);
  expect(opacity).toBe(1);
});

test('a deck attack keeps the drawn monster in view while its slot is picked (19.9)', async ({ browser }) => {
  const [odo] = await startGame(browser, ['Odo', 'Kye']);
  if (!odo) {
    throw new Error('no pages');
  }
  // Replay the draw of a deck attack: the top card is revealed and the
  // attacker is asked for a slot to put it in.
  const name = await odo.evaluate(() => {
    const store = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store;
    const s = store.getState();
    const game = structuredClone(s.game);
    const you = s.lobby.table?.you ?? 0;
    const slots = game?.view.monsters ?? [];
    const card = slots[1]?.top?.card;
    if (!game || !card || slots.some((m) => !m.top)) {
      throw new Error('no game');
    }
    game.view.revealed = { id: 989, card, charged: false };
    game.view.waiting = { kind: 4, player: you, purpose: 2, options: slots.map((_, i) => `slot ${String(i + 1)}`) };
    game.view.options = slots.map((m) => ({ object: m.top?.id ?? 0, player: -1 }));
    game.allowed = slots.map((_, i) => ({ player: you, kind: 4, choice: i }));
    const events = [{ kind: 'card_revealed', player: you, object: 989, card }];
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
    return s.cards.get(card)?.name ?? card;
  });
  // After the large reveal, the card stays face up on the monster deck
  // and in the dialog until the slot is picked.
  const dialog = odo.getByRole('dialog', { name: 'Choose a slot for the monster' });
  await expect(dialog).toBeVisible({ timeout: 8000 });
  await expect(odo.getByTestId('effect-reveal')).toHaveCount(0);
  await odo.waitForTimeout(400); // the deck card pops in
  await shot(odo, 'deck-attack-pick-slot');
  await expect(odo.locator('[data-zone="deck-monster"] [data-zone="obj-989"]')).toBeVisible();
  await expect(dialog.getByLabel(name, { exact: true })).toBeVisible();
  // Once it flies into its slot, the deck no longer shows it.
  await odo.evaluate(() => {
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const [first] = game?.view.monsters ?? [];
    const drawn = game?.view.revealed;
    if (!game || !first || !drawn) {
      throw new Error('no game');
    }
    first.top = { ...drawn, id: 990 };
    delete game.view.revealed;
    game.view.waiting = { kind: 1, player: s.lobby.table?.you ?? 0 };
    game.allowed = [];
    const events = [{ kind: 'card_revealed', player: -1, object: 990, card: drawn.card }];
    s.receive({ type: 'update', data: { ...game, step: 3, events } });
  });
  await expect(odo.getByTestId('effect-place')).toBeVisible({ timeout: 5000 });
  await expect(odo.locator('[data-zone="obj-989"]')).toHaveCount(0);
  await expect(odo.locator('[data-zone="obj-990"]')).toBeVisible();
  await expect(dialog).toHaveCount(0);
});

// trackCards records, every frame, the box of the large revealed card
// and of a card flying to its place, in window.boxes (15.73). It stops
// when the boxes are read, or after 10 s: the test's pages stay open
// until the worker ends, and a loop left running slows later tests.
async function trackCards(page: Page): Promise<void> {
  await page.evaluate(() => {
    interface Box {
      x: number;
      y: number;
      w: number;
    }
    const w = window as unknown as { boxes: { reveal: Box[]; place: Box[] }; tracking: boolean };
    w.boxes = { reveal: [], place: [] };
    w.tracking = true;
    const end = performance.now() + 10_000;
    const box = (sel: string): Box | null => {
      const el = document.querySelector(sel);
      if (!el) {
        return null;
      }
      const r = el.getBoundingClientRect();
      return { x: r.left + r.width / 2, y: r.top + r.height / 2, w: r.width };
    };
    const frame = () => {
      const reveal = box('[data-testid="effect-reveal"] .card');
      const place = box('[data-testid="effect-place"] .effect-fly');
      if (reveal) {
        w.boxes.reveal.push(reveal);
      }
      if (place) {
        w.boxes.place.push(place);
      }
      if (w.tracking && performance.now() < end) {
        requestAnimationFrame(frame);
      }
    };
    requestAnimationFrame(frame);
  });
}

interface Boxes {
  reveal: { x: number; y: number; w: number }[];
  place: { x: number; y: number; w: number }[];
}

function boxes(page: Page): Promise<Boxes> {
  return page.evaluate(() => {
    const w = window as unknown as { boxes: Boxes; tracking: boolean };
    w.tracking = false;
    return w.boxes;
  });
}

// centre is where an element's middle is, once its box holds still: the
// coin, for one, is drawn anew when its image fails to load.
async function centre(page: Page, selector: string): Promise<{ x: number; y: number }> {
  let last = '';
  let box: { x: number; y: number; width: number; height: number } | null = null;
  await expect
    .poll(
      async () => {
        box = await page.locator(selector).boundingBox();
        const now = JSON.stringify(box);
        const still = box !== null && now === last;
        last = now;
        return still;
      },
      { intervals: [100] },
    )
    .toBe(true);
  const b = box as { x: number; y: number; width: number; height: number } | null;
  if (!b) {
    throw new Error(`no ${selector}`);
  }
  return { x: b.x + b.width / 2, y: b.y + b.height / 2 };
}

test('a refilled slot\'s card grows out of its deck, then flies from the middle (15.73)', async ({ browser }) => {
  const [lia] = await startGame(browser, ['Lia', 'Otto']);
  if (!lia) {
    throw new Error('no pages');
  }
  const deck = await centre(lia, '[data-zone="deck-monster"]');
  await trackCards(lia);
  // Replay a slot refill: a new monster is revealed into the first slot.
  await lia.evaluate(() => {
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const [first, second] = game?.view.monsters ?? [];
    if (!game || !first || !second?.top) {
      throw new Error('no game');
    }
    first.top = { ...second.top, id: 990 };
    const events = [{ kind: 'card_revealed', player: -1, object: 990, card: second.top.card }];
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
  });
  await expect(lia.getByTestId('effect-reveal')).toBeVisible();
  await shot(lia, 'reveal-from-deck');
  await lia.waitForTimeout(700);
  await shot(lia, 'reveal-held');
  await expect(lia.getByTestId('effect-place')).toBeVisible({ timeout: 5000 });
  await shot(lia, 'reveal-to-slot');
  await expect(lia.locator('[data-zone="obj-990"]')).toBeVisible();
  // The reveal starts small on the deck; the flight starts large in
  // the middle of the window.
  const { reveal, place } = await boxes(lia);
  const start = reveal[0];
  const flight = place[0];
  if (!start || !flight) {
    throw new Error('no boxes');
  }
  expect(Math.abs(start.x - deck.x)).toBeLessThan(30);
  expect(Math.abs(start.y - deck.y)).toBeLessThan(30);
  expect(start.w).toBeLessThan(160);
  expect(Math.abs(flight.x - 640)).toBeLessThan(30);
  expect(Math.abs(flight.y - 360)).toBeLessThan(30);
  expect(flight.w).toBeGreaterThan(200);
});

test('a revealed card that stays in its deck shrinks back onto it (15.73)', async ({ browser }) => {
  const [noa] = await startGame(browser, ['Noa', 'Pim']);
  if (!noa) {
    throw new Error('no pages');
  }
  const deck = await centre(noa, '[data-zone="deck-monster"]');
  await trackCards(noa);
  // Replay a reveal of the deck's top card that leaves it there.
  await noa.evaluate(() => {
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const card = game?.view.monsters?.[1]?.top?.card;
    if (!game || !card) {
      throw new Error('no game');
    }
    const events = [{ kind: 'card_revealed', player: -1, object: 991, card, text: 'monster deck' }];
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
  });
  await expect(noa.getByTestId('effect-reveal')).toBeVisible();
  await noa.waitForTimeout(1250);
  await shot(noa, 'reveal-back-to-deck');
  await expect(noa.getByTestId('effect-reveal')).toHaveCount(0, { timeout: 5000 });
  const { reveal } = await boxes(noa);
  const end = reveal.at(-1);
  if (!end) {
    throw new Error('no boxes');
  }
  expect(Math.abs(end.x - deck.x)).toBeLessThan(30);
  expect(Math.abs(end.y - deck.y)).toBeLessThan(30);
  expect(end.w).toBeLessThan(160);
});

// trackBoxes records, every frame, the box of the first element each
// selector finds, in window.tracked (15.74). It stops as trackCards does.
async function trackBoxes(page: Page, selectors: Record<string, string>): Promise<void> {
  await page.evaluate((sel) => {
    const w = window as unknown as { tracked: Record<string, { x: number; y: number; w: number; hidden: boolean }[]>; tracking: boolean };
    w.tracked = Object.fromEntries(Object.keys(sel).map((k) => [k, []]));
    w.tracking = true;
    const end = performance.now() + 10_000;
    const frame = () => {
      for (const [k, s] of Object.entries(sel)) {
        const el = document.querySelector(s);
        if (el) {
          const r = el.getBoundingClientRect();
          const hidden = getComputedStyle(el).visibility === 'hidden';
          w.tracked[k]?.push({ x: r.left + r.width / 2, y: r.top + r.height / 2, w: r.width, hidden });
        }
      }
      if (w.tracking && performance.now() < end) {
        requestAnimationFrame(frame);
      }
    };
    requestAnimationFrame(frame);
  }, selectors);
}

function tracked(page: Page): Promise<Record<string, { x: number; y: number; w: number; hidden: boolean }[]>> {
  return page.evaluate(() => {
    const w = window as unknown as { tracked: Record<string, { x: number; y: number; w: number; hidden: boolean }[]>; tracking: boolean };
    w.tracking = false;
    return w.tracked;
  });
}

// putLoot shows a loot card of the viewer's hand on the stack as
// object 700, played by the viewer (mine) or the other player. With
// played, the update says it was just played.
async function putLoot(page: Page, mine: boolean, played: boolean): Promise<void> {
  await page.evaluate(
    ({ mine, played }) => {
      const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
      const game = structuredClone(s.game) as unknown as {
        view: { players: { id: number; hand?: { card: string }[] }[]; stack?: unknown[] };
      } | null;
      const you = s.lobby.table?.you ?? 0;
      const other = game?.view.players.find((p) => p.id !== you)?.id ?? 1;
      const card = game?.view.players.find((p) => p.id === you)?.hand?.[0]?.card;
      if (!game || !card) {
        throw new Error('no game');
      }
      const controller = mine ? you : other;
      const none = { card: '', index: 0 };
      game.view.stack = [{ id: 50, kind: 0, controller, source: 700, card, target: { player: -1 }, ability: none, roll_for: none }];
      const events = played ? [{ kind: 'loot_played', player: controller, object: 700, card }] : [];
      s.receive({ type: 'update', data: { ...game, step: 2, events } });
    },
    { mine, played },
  );
}

test("another player's loot card rises from their hand and shrinks into the stack (15.74)", async ({ browser }) => {
  const [ula] = await startGame(browser, ['Ula', 'Vic']);
  if (!ula) {
    throw new Error('no pages');
  }
  const other = await ula.evaluate(() => {
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    return s.lobby.table?.you === 0 ? 1 : 0;
  });
  const hand = await centre(ula, `[data-zone="hand-${String(other)}"]`);
  const stack = await centre(ula, '[data-zone="stack"]');
  await trackBoxes(ula, { pop: '[data-testid="played-pop"] .card' });
  await putLoot(ula, false, true);
  await expect(ula.getByTestId('played-pop')).toBeVisible();
  await shot(ula, 'played-from-hand');
  await ula.waitForTimeout(1250);
  await shot(ula, 'played-to-stack');
  await expect(ula.locator('[data-zone="played-700"]')).toBeVisible({ timeout: 5000 });
  await shot(ula, 'played-under-stack');
  const { pop } = await tracked(ula);
  const first = pop?.[0];
  const last = pop?.at(-1);
  if (!first || !last) {
    throw new Error('no boxes');
  }
  expect(Math.abs(first.x - hand.x)).toBeLessThan(40);
  expect(Math.abs(first.y - hand.y)).toBeLessThan(40);
  expect(first.w).toBeLessThan(120);
  expect(Math.abs(last.x - stack.x)).toBeLessThan(40);
  expect(Math.abs(last.y - stack.y)).toBeLessThan(40);
  expect(last.w).toBeLessThan(140);
});

for (const mine of [false, true]) {
  const from = mine ? 'its stack row' : 'the large card under the stack';
  test(`a resolved loot card flies to the discard from ${from} (15.74)`, async ({ browser }) => {
    const [wes] = await startGame(browser, ['Wes', 'Xia']);
    if (!wes) {
      throw new Error('no pages');
    }
    await putLoot(wes, mine, false);
    const at = mine ? '[data-zone="stack-src-700"]' : '[data-zone="played-700"] .card';
    await expect(wes.locator(at)).toBeVisible();
    const start = await centre(wes, at);
    await trackBoxes(wes, { fly: '[data-testid="effect-fly"] .effect-fly .card', played: '[data-zone="played-700"]' });
    // The loot resolves: it leaves the stack for the loot discard.
    await wes.evaluate(() => {
      const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
      const game = structuredClone(s.game) as unknown as { view: { stack?: { controller: number; card: string }[] } } | null;
      const item = game?.view.stack?.[0];
      if (!game || !item) {
        throw new Error('no game');
      }
      game.view.stack = [];
      const events = [{ kind: 'discarded', player: item.controller, object: 900, prev: 700, card: item.card, text: 'stack' }];
      s.receive({ type: 'update', data: { ...game, step: 3, events } });
    });
    await expect(wes.getByTestId('effect-fly')).toBeVisible();
    await shot(wes, mine ? 'discard-from-stack-row' : 'discard-from-played');
    await expect(wes.getByTestId('effect-fly')).toHaveCount(0, { timeout: 5000 });
    const boxes = await tracked(wes);
    const first = boxes.fly?.[0];
    if (!first) {
      throw new Error('no boxes');
    }
    expect(Math.abs(first.x - start.x)).toBeLessThan(30);
    expect(Math.abs(first.y - start.y)).toBeLessThan(30);
    if (!mine) {
      // The large card hides while its copy flies, so it is seen once.
      expect(boxes.played?.some((b) => b.hidden)).toBe(true);
    }
  });
}

test('stolen cents fly as a coin from the victim to the thief (15.75)', async ({ browser }) => {
  const [yan] = await startGame(browser, ['Yan', 'Zed']);
  if (!yan) {
    throw new Error('no pages');
  }
  const ids = await yan.evaluate(() => {
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    const you = s.lobby.table?.you ?? 0;
    return { you, other: you === 0 ? 1 : 0 };
  });
  const thief = await centre(yan, `[data-zone="player-${String(ids.you)}"] .coin`);
  const victim = await centre(yan, `[data-zone="player-${String(ids.other)}"] .coin`);
  await trackBoxes(yan, { coin: '[data-testid="coin-shot"]' });
  // The viewer steals 2¢ from the other player.
  await yan.evaluate((ids) => {
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const events = [{ kind: 'stole', player: ids.you, amount: 2, text: String(ids.other) }];
    s.receive({ type: 'update', data: { ...game, step: 2, events } });
  }, ids);
  await expect(yan.locator('.effects .effect-popup').first()).toBeVisible();
  await yan.waitForTimeout(150);
  await shot(yan, 'steal-coin');
  await expect(yan.getByText('+2¢ stolen')).toBeVisible();
  await shot(yan, 'steal-landed');
  await expect(yan.locator('.effects .effect-popup')).toHaveCount(0, { timeout: 5000 });
  const { coin } = await tracked(yan);
  const first = coin?.[0];
  const last = coin?.at(-1);
  if (!first || !last) {
    throw new Error('no boxes');
  }
  expect(Math.abs(first.x - victim.x)).toBeLessThan(30);
  expect(Math.abs(first.y - victim.y)).toBeLessThan(30);
  expect(Math.abs(last.x - thief.x)).toBeLessThan(30);
  expect(Math.abs(last.y - thief.y)).toBeLessThan(30);
});

test('a changed roll turns from its old value, hit by the card that changed it (15.76)', async ({ browser }) => {
  const [abe] = await startGame(browser, ['Abe', 'Bea']);
  if (!abe) {
    throw new Error('no pages');
  }
  // A roll of 3 by the viewer is on the stack.
  const ids = await abe.evaluate(() => {
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    const game = structuredClone(s.game) as unknown as {
      view: { players: { id: number; character: { id: number } }[]; stack?: unknown[] };
    } | null;
    const you = s.lobby.table?.you ?? 0;
    const character = game?.view.players.find((p) => p.id === you)?.character.id;
    if (!game || character === undefined) {
      throw new Error('no game');
    }
    const none = { card: '', index: 0 };
    game.view.stack = [{ id: 60, kind: 3, controller: you, source: 0, roll: 3, target: { player: you }, ability: none, roll_for: none }];
    s.receive({ type: 'update', data: { ...game, step: 2, events: [] } });
    return { you, character };
  });
  const from = await centre(abe, `[data-zone="obj-${String(ids.character)}"]`);
  await trackBoxes(abe, { shot: '[data-testid="roll-shot"]', die: '[data-testid="changed-die"]' });
  // The viewer's character changes it to 5.
  await abe.evaluate((ids) => {
    const s = (window as unknown as { fourSouls: { store: Store } }).fourSouls.store.getState();
    const game = structuredClone(s.game);
    const events = [{ kind: 'roll_changed', player: ids.you, amount: 5, stack_id: 60, source: ids.character }];
    s.receive({ type: 'update', data: { ...game, step: 3, events } });
  }, ids);
  await expect(abe.locator('.effects')).toBeVisible();
  await shot(abe, 'roll-change-start');
  const die = abe.getByTestId('changed-die');
  await expect(die).toHaveAttribute('data-face', /[35]/);
  await expect(abe.getByText("Abe's roll: 3 → 5")).toBeVisible();
  await expect(abe.getByTestId('roll-shot')).toBeVisible();
  await shot(abe, 'roll-change-shot');
  await expect(die).toHaveAttribute('data-face', '5');
  await shot(abe, 'roll-change-turned');
  await expect(die).toHaveCount(0, { timeout: 5000 });
  // The shot starts on the character and ends on the die.
  const boxes = await tracked(abe);
  const first = boxes.shot?.[0];
  const last = boxes.shot?.at(-1);
  const at = boxes.die?.at(-1);
  if (!first || !last || !at) {
    throw new Error('no boxes');
  }
  expect(Math.abs(first.x - from.x)).toBeLessThan(30);
  expect(Math.abs(first.y - from.y)).toBeLessThan(30);
  expect(Math.abs(last.x - at.x)).toBeLessThan(40);
  expect(Math.abs(last.y - at.y)).toBeLessThan(40);
});
