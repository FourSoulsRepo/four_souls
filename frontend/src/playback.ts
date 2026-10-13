// Event playback (7.6, V-08). Updates queue up and play one at a time:
// first their events as short effects, then the update's view is shown.
// Playback only changes what is drawn; the game state is the server's.
import { kindOf, type CardInfo, type Kind } from './cards';
import type { GameEvent, Update } from './protocol';

export type Speed = 'normal' | 'fast' | 'slow' | 'off';

// An effect names the places it plays at by zone (data-zone on the
// table), most precise first; the first one on screen is used.
export type Effect =
  | { kind: 'fly'; from: string[]; to: string[]; card?: string; back: Kind }
  | { kind: 'popup'; at: string[]; text: string; tone: 'good' | 'bad' }
  | { kind: 'dice'; value: number; label: string }
  | { kind: 'banner'; text: string }
  | { kind: 'reveal'; card: string }
  // lunge: the attacker's card hits the target and returns (15.3, V-12).
  | { kind: 'lunge'; from: string[]; to: string[]; amount: number };

const noPlayer = -1;

interface Context {
  card: (ref: string) => CardInfo | undefined;
  name: (player: number) => string;
  // character is a player's character object, if known.
  character: (player: number) => number | undefined;
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
      case 'roll_changed':
        out.push({ kind: 'dice', value: amount, label: e.text ?? '' });
        break;
      case 'gained_cents':
        out.push({ kind: 'popup', at: playerZones(e.player), text: `+${String(amount)}¢`, tone: 'good' });
        break;
      case 'healed':
        out.push({ kind: 'popup', at: targetZones(e), text: `+${String(amount)} HP`, tone: 'good' });
        break;
      case 'damage_pending':
        // Attack damage shows when it lands, as its lunge.
        if (e.text === 'combat') {
          break;
        }
        out.push({ kind: 'popup', at: targetZones(e), text: `−${String(amount)}`, tone: 'bad' });
        break;
      case 'damaged':
        // Attack damage: its source lunges at the target. Other damage
        // shows with its stack item.
        if (e.text === 'combat' && e.source) {
          out.push({ kind: 'lunge', from: [`obj-${String(e.source)}`], to: hitZones(e, ctx), amount });
        }
        break;
      case 'prevented':
        out.push({ kind: 'popup', at: targetZones(e), text: 'Prevented', tone: 'good' });
        break;
      case 'stole':
        out.push({ kind: 'popup', at: playerZones(e.player), text: `+${String(amount)}¢ stolen`, tone: 'good' });
        break;
      case 'card_revealed':
        if (e.card) {
          out.push({ kind: 'reveal', card: e.card });
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
  return out;
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
