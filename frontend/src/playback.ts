// Event playback (7.6, V-08). Updates queue up and play one at a time:
// first their events as short effects, then the update's view is shown.
// Playback only changes what is drawn; the game state is the server's.
import { kindOf, type CardInfo, type Kind } from './cards';
import type { GameEvent, GameView, Update } from './protocol';

export type Speed = 'normal' | 'fast' | 'slow' | 'off';

// An effect names the places it plays at by zone (data-zone on the
// table), most precise first; the first one on screen is used.
export type Effect =
  // fly: middle starts the card large in the middle, where the reveal
  // before it held it (15.73).
  | { kind: 'fly'; from: string[]; to: string[]; card?: string; back: Kind; middle?: boolean }
  | { kind: 'popup'; at: string[]; text: string; tone: 'good' | 'bad' }
  // dice: thrown for a fresh roll or a reroll (15.62); a changed roll
  // spins in.
  | { kind: 'dice'; value: number; label: string; player: number; thrown?: boolean }
  | { kind: 'banner'; text: string }
  // reveal: the card grows out of its deck to the middle; it holds
  // there when the next effect moves it, else shrinks back (15.73).
  | { kind: 'reveal'; card: string; deck: string[]; hold?: boolean }
  // played: another player's loot card rises from their hand, shows
  // large for a moment and goes on the stack (15.31, H-10, 15.74).
  | { kind: 'played'; card: string; who: string; from: string[]; to: string[] }
  // place: a card flies into a slot and stays there until the view is
  // shown, so later effects play over it (19.3).
  | { kind: 'place'; from: string[]; to: string[]; card: string; back: Kind; middle?: boolean }
  // lunge: the attacker's card hits the target and returns (15.3, V-12).
  | { kind: 'lunge'; from: string[]; to: string[]; amount: number }
  // hits: other damage flies from its source and bursts on every
  // damaged card at once (15.37, T-15, T-24).
  | { kind: 'hits'; hits: Hit[] }
  // deaths: a big skull on each dying card before it leaves (T-26).
  | { kind: 'deaths'; at: string[][] };

export interface Hit {
  from: string[];
  to: string[];
  amount: number;
}

const noPlayer = -1;

interface Context {
  card: (ref: string) => CardInfo | undefined;
  name: (player: number) => string;
  // character is a player's character object, if known.
  character: (player: number) => number | undefined;
  // slot is the zone of the slot an object is the top card of, if any.
  slot: (object: number) => string | undefined;
  // seen tells if the update before revealed the card (19.4).
  seen: (card: string) => boolean;
  // you is the viewer's seat; -1 for none.
  you: number;
}

function playerZones(p: number): string[] {
  return [`player-${String(p)}`];
}

function handZones(p: number): string[] {
  return [`hand-${String(p)}`, `player-${String(p)}`];
}

// targetZones: an object, else a player.
function targetZones(e: GameEvent): string[] {
  return e.object ? [`obj-${String(e.object)}`, ...playerZones(e.player)] : playerZones(e.player);
}

// hitZones: the damaged object, else the player's character card.
function hitZones(e: GameEvent, ctx: Context): string[] {
  if (e.object) {
    return [`obj-${String(e.object)}`];
  }
  const c = ctx.character(e.player);
  return c === undefined ? playerZones(e.player) : [`obj-${String(c)}`, ...playerZones(e.player)];
}

// slotZone is the zone of the monster or shop slot whose top card is
// object, as the board names it.
export function slotZone(view: GameView, object: number): string | undefined {
  for (const [name, slots] of [
    ['monster', view.monsters],
    ['shop', view.shop],
  ] as const) {
    const i = (slots ?? []).findIndex((s) => s.top?.id === object);
    if (i >= 0) {
      return `${name}-slot-${String(i)}`;
    }
  }
  return undefined;
}

// effectsFor turns one update's events into effects. Events without an
// effect still change the table when the view is shown.
export function effectsFor(events: GameEvent[], ctx: Context): Effect[] {
  const out: Effect[] = [];
  for (const e of events) {
    const amount = e.amount ?? 0;
    switch (e.kind) {
      case 'looted':
        out.push({
          kind: 'fly',
          from: [e.text === 'from the discard' ? 'discard-loot' : 'deck-loot'],
          to: handZones(e.player),
          card: e.card,
          back: 'loot',
        });
        break;
      case 'discarded': {
        // A loot card leaves the stack from its large card under the
        // stack, else its stack row (15.74).
        if (e.text === 'stack' && e.prev) {
          const from = [`played-${String(e.prev)}`, `stack-src-${String(e.prev)}`, `obj-${String(e.prev)}`, 'deck-loot'];
          out.push({ kind: 'fly', from, to: ['discard-loot'], card: e.card, back: 'loot' });
          break;
        }
        const kind = e.card ? kindOf(ctx.card(e.card)?.type ?? '') : 'loot';
        const from = [`obj-${String(e.prev ?? e.object ?? 0)}`];
        if (e.player !== noPlayer) {
          from.push(...handZones(e.player));
        }
        from.push(`deck-${kind}`);
        out.push({ kind: 'fly', from, to: [`discard-${kind}`], card: e.card, back: kind });
        break;
      }
      case 'dice_rolled':
        out.push({ kind: 'dice', value: amount, label: e.text ?? '', player: e.player, thrown: true });
        break;
      case 'roll_changed':
        out.push({ kind: 'dice', value: amount, label: e.text ?? '', player: e.player });
        break;
      case 'gained_cents':
        out.push({ kind: 'popup', at: playerZones(e.player), text: `+${String(amount)}¢`, tone: 'good' });
        break;
      case 'healed':
        out.push({ kind: 'popup', at: targetZones(e), text: `+${String(amount)} HP`, tone: 'good' });
        break;
      case 'damaged': {
        // Damage shows when it lands. Attack damage: its source lunges
        // at the target. Other damage flies from its source: a card on
        // the table, else its stack row.
        if (e.text === 'combat' && e.source) {
          out.push({ kind: 'lunge', from: [`obj-${String(e.source)}`], to: hitZones(e, ctx), amount });
          break;
        }
        const from = e.source ? [`obj-${String(e.source)}`, `stack-src-${String(e.source)}`] : [];
        const hit = { from, to: hitZones(e, ctx), amount };
        const last = out.at(-1);
        if (last?.kind === 'hits') {
          last.hits.push(hit);
        } else {
          out.push({ kind: 'hits', hits: [hit] });
        }
        break;
      }
      case 'died': {
        // A monster's card left its place under a new ID (prev).
        const at = e.prev ? [`obj-${String(e.prev)}`] : hitZones(e, ctx);
        const last = out.at(-1);
        if (last?.kind === 'deaths') {
          last.at.push(at);
        } else {
          out.push({ kind: 'deaths', at: [at] });
        }
        break;
      }
      case 'prevented':
        out.push({ kind: 'popup', at: targetZones(e), text: 'Prevented', tone: 'good' });
        break;
      case 'stole':
        out.push({ kind: 'popup', at: playerZones(e.player), text: `+${String(amount)}¢ stolen`, tone: 'good' });
        break;
      case 'card_revealed': {
        if (!e.card) {
          break;
        }
        // A deck attack reveals its card, then puts it in a slot in the
        // next update: it shows large once (19.4).
        const kind = kindOf(ctx.card(e.card)?.type ?? '');
        if (!ctx.seen(e.card)) {
          out.push({ kind: 'reveal', card: e.card, deck: [`deck-${kind}`] });
        }
        // A card put into a slot, e.g. a deck attack's monster, is there
        // before the attack roll that follows it.
        const slot = e.object ? ctx.slot(e.object) : undefined;
        if (slot) {
          out.push({ kind: 'place', from: [`deck-${kind}`], to: [slot], card: e.card, back: kind });
        }
        break;
      }
      case 'loot_played':
        // The player's own play needs no reminder (15.31).
        if (e.card && e.player !== ctx.you) {
          out.push({ kind: 'played', card: e.card, who: ctx.name(e.player), from: handZones(e.player), to: ['stack'] });
        }
        break;
      case 'turn_started':
        out.push({ kind: 'banner', text: `${ctx.name(e.player)}'s turn` });
        break;
      case 'game_won':
        out.push({ kind: 'banner', text: `${ctx.name(e.player)} wins!` });
        break;
    }
  }
  holdReveals(out);
  return out;
}

// holdReveals joins a reveal to a move of its card right after it: the
// reveal holds the card large in the middle, and the move starts there
// (15.73). The card is seen on one path, from its deck to its place.
function holdReveals(effects: Effect[]): void {
  effects.forEach((e, i) => {
    const next = effects.at(i + 1);
    if (e.kind === 'reveal' && (next?.kind === 'fly' || next?.kind === 'place') && next.card === e.card) {
      e.hold = true;
      next.middle = true;
    }
  });
}

function speedScale(speed: Speed): number {
  switch (speed) {
    case 'slow':
      return 1.6;
    case 'fast':
      return 0.5;
    case 'off':
      return 0;
    default:
      return 1;
  }
}

// scaleFor is how long effects take, as a factor of their normal time;
// 0 skips them. A long queue or a big burst plays faster, so the table
// catches up with the game.
export function scaleFor(speed: Speed, pending: number, effects: number): number {
  const base = speedScale(speed);
  if (base === 0 || pending > 8) {
    return 0;
  }
  if (pending > 2 || effects > 6) {
    return base * 0.4;
  }
  return base;
}

interface Deps {
  // begin starts an update, before its effects (animations.ts).
  begin: () => void;
  // show draws an update's view.
  show: (u: Update) => void;
  // run plays one effect and resolves when it ends.
  run: (e: Effect, scale: number) => Promise<void>;
  effects: (u: Update) => Effect[];
  speed: () => Speed;
  // shown is the update on screen now, if any.
  shown: () => Update | null;
}

// Playback is the queue of updates waiting to be shown.
export interface Playback {
  push: (u: Update) => void;
  clear: () => void;
  // idle resolves when everything queued is shown.
  idle: () => Promise<void>;
  // pending is how many updates wait to be shown.
  pending: () => number;
}

export function createPlayback(deps: Deps): Playback {
  let queue: Update[] = [];
  // busy is set while drain runs. An update pushed while drain is
  // finishing would wait forever, so drain starts again if one did.
  let busy = false;
  // waiters are idle() calls waiting for the queue to empty.
  let waiters: (() => void)[] = [];

  function settle(): void {
    if (!busy && queue.length === 0) {
      for (const w of waiters) {
        w();
      }
      waiters = [];
    }
  }

  async function drain(): Promise<void> {
    try {
      await play();
    } finally {
      busy = false;
      if (queue.length > 0) {
        start();
      }
      settle();
    }
  }

  function start(): void {
    busy = true;
    void drain();
  }

  async function play(): Promise<void> {
    for (let u = queue.shift(); u; u = queue.shift()) {
      deps.begin();
      // The first update has no table to play on yet.
      if (deps.shown() !== null) {
        const effects = deps.effects(u);
        const scale = scaleFor(deps.speed(), queue.length, effects.length);
        if (scale > 0) {
          for (const e of effects) {
            await deps.run(e, scale);
          }
        }
      }
      deps.show(u);
    }
  }

  return {
    push: (u) => {
      queue.push(u);
      if (!busy) {
        start();
      }
    },
    clear: () => {
      queue = [];
    },
    pending: () => queue.length,
    idle: () =>
      new Promise((resolve) => {
        waiters.push(resolve);
        settle();
      }),
  };
}
