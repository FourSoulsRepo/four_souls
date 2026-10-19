// The prompt banner (15.28, V-20): what the game waits for, in a few
// words. Pure, so every prompt is tested without a browser.
import { choiceTitle, PromptKind, type Actions } from './actions';
import type { GameView } from './protocol';
import { findCard, stackLine, type Names } from './stack';
import type { Aiming } from './store';

// BannerKind is the prompt a banner is about. "aim" is the player's own
// aim, before anything is sent (15.25).
export type BannerKind = 'turn' | 'respond' | 'act' | 'choose' | 'discard' | 'order' | 'throw' | 'aim';

// Banner: title says what to do, or whom the game waits for; hint adds
// what for. mine: the game waits for this player.
export interface Banner {
  kind: BannerKind;
  mine: boolean;
  player: number;
  title: string;
  hint?: string;
}

// promptBanner says what the game waits for now; nothing once it is
// over. aiming is the viewer's aim, if any.
export function promptBanner(
  view: GameView,
  you: number,
  n: Names,
  actions: Actions,
  aiming: Aiming | null = null,
): Banner | null {
  if (view.over) {
    return null;
  }
  if (aiming) {
    const card = findCard(view, aiming.source);
    return {
      kind: 'aim',
      mine: true,
      player: you,
      title: 'Choose a target',
      hint: `${card ? `For ${n.card(card)}. ` : ''}Esc cancels`,
    };
  }
  const w = view.waiting;
  const ask = asked(view, n);
  if (!ask) {
    return null;
  }
  if (w.player === you) {
    const hint = ask.kind === 'turn' ? turnHint(actions) : ask.hint;
    return { kind: ask.kind, mine: true, player: w.player, title: ask.title, ...(hint ? { hint } : {}) };
  }
  return { kind: ask.kind, mine: false, player: w.player, title: `Waiting for ${n.player(w.player)}`, hint: ask.other };
}

// asked is the prompt in words: title and hint for the asked player,
// other for everyone else.
function asked(view: GameView, n: Names): { kind: BannerKind; title: string; hint?: string; other: string } | null {
  const { waiting: w, turn } = view;
  switch (w.kind) {
    case PromptKind.priority: {
      const top = (view.stack ?? []).at(-1);
      if (top) {
        const title = 'Respond, or pass';
        return { kind: 'respond', title, hint: `On the stack: ${stackLine(top, view, n)}`, other: title };
      }
      if (turn.active === w.player) {
        return { kind: 'turn', title: 'Your turn', other: 'Playing their turn' };
      }
      return { kind: 'act', title: 'You may act', hint: 'Use an item or a loot card, or pass', other: 'May act, or pass' };
    }
    case PromptKind.choose: {
      const title = choiceTitle(w);
      return { kind: 'choose', title, other: title };
    }
    case PromptKind.discard: {
      const count = w.count ?? 0;
      const title = `Discard ${String(count)} loot card${count === 1 ? '' : 's'}`;
      return { kind: 'discard', title, other: title };
    }
    case PromptKind.order: {
      const title = choiceTitle(w);
      return { kind: 'order', title, other: title };
    }
    case PromptKind.throw: {
      const title = 'Throw the die';
      return w.text
        ? { kind: 'throw', title, hint: `For the ${w.text}`, other: `${title} for the ${w.text}` }
        : { kind: 'throw', title, other: title };
    }
    default:
      return null;
  }
}

// turnHint lists what the player may do on their own turn, or says that
// only ending it is left (H-20).
function turnHint(actions: Actions): string {
  const { turn, cards } = actions;
  if (turn.idle) {
    return 'Nothing left to do: end the turn';
  }
  const can = [
    cards.size > 0 ? 'play a card' : '',
    turn.attack ? 'attack' : '',
    turn.buy ? 'buy' : '',
    turn.endTurn ? 'end the turn' : '',
  ].filter(Boolean);
  if (can.length === 0) {
    return 'Pass to go on';
  }
  const last = can.pop() ?? '';
  const words = can.length > 0 ? `${can.join(', ')}, or ${last}` : last;
  return words.charAt(0).toUpperCase() + words.slice(1);
}

// picking: the player picks a target on the table now, so the cards
// that answer glow and the rest dim (15.28). A choice counts once one
// of its options is on the table.
export function picking(actions: Actions, aiming: Aiming | null): boolean {
  return aiming !== null || (actions.choice?.options.some((o) => o.target) ?? false);
}

// stackPrompt is the banner's prompt as the stack panel's line (V-11):
// only for prompts other than priority, whose players the panel lists.
export function stackPrompt(b: Banner | null, n: Names): string | undefined {
  if (!b || b.kind === 'turn' || b.kind === 'respond' || b.kind === 'act') {
    return undefined;
  }
  return b.mine ? b.title : `${n.player(b.player)}: ${b.hint ?? ''}`;
}
