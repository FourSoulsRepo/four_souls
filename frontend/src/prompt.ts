// The prompt banner (15.28, V-20): what the game waits for, in a few
// words. Pure, so every prompt is tested without a browser.
import { choiceTitle, PromptKind, Purpose, type Actions } from './actions';
import type { GameView } from './protocol';
import { findCard, StackKind, stackLine, type LinePart, type Names } from './stack';
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

// activity says what the active player is doing, for every other seat
// (15.30, V-21): their own play or use on top of the stack, the death
// penalty, an attack, a purchase, or thinking while the game waits for
// them. Only public state, so spectators get the same words. Nothing
// while the game waits for others with none of these going on.
export function activity(view: GameView, n: Names): LinePart[] | null {
  if (view.over) {
    return null;
  }
  const { active } = view.turn;
  const who: LinePart = { text: n.player(active), role: 'player', player: active };
  const card = (ref: string | undefined): LinePart => ({ text: ref ? n.card(ref) : 'a card', role: 'card' });
  const w = view.waiting;
  const top = (view.stack ?? []).at(-1);
  if (top?.controller === active && top.kind === StackKind.loot) {
    return [who, { text: ' plays ' }, card(top.card)];
  }
  if (top?.controller === active && top.kind === StackKind.ability) {
    return [who, { text: ' uses ' }, card(top.card ?? top.source_card ?? findCard(view, top.source))];
  }
  const penalty =
    (top?.kind === StackKind.penalty && top.target.player === active) ||
    (w.player === active && w.kind === PromptKind.choose && (w.purpose === Purpose.penaltyItem || w.purpose === Purpose.penaltyLoot));
  if (penalty) {
    return [who, { text: ' pays the death penalty' }];
  }
  const { attack } = view;
  if (attack) {
    const target = attack.target ?? 0;
    if (target === 0) {
      return attack.deck ? [who, { text: ' attacks the monster deck' }] : [who, { text: ' declares an attack' }];
    }
    // An attacked player is named, not their character (R-ATK-20).
    const victim = view.players.find((p) => p.character.id === target);
    const ref = findCard(view, target);
    const aim: LinePart = victim
      ? { text: n.player(victim.id), role: 'target', player: victim.id }
      : { text: ref ? n.card(ref) : 'a monster', role: 'target' };
    return [who, { text: ' attacks ' }, aim];
  }
  if (view.buying) {
    return [who, { text: ' is buying' }];
  }
  if (w.player === active && w.kind !== PromptKind.none) {
    return [who, { text: ' is thinking' }];
  }
  return null;
}

// fighters are the cards in a fight, framed red for everyone (15.31,
// F-19): the active player's character and what it attacks, once the
// target is chosen.
export function fighters(view: GameView): number[] {
  const target = view.attack?.target ?? 0;
  const attacker = view.players.find((p) => p.id === view.turn.active)?.character.id;
  if (view.over || target === 0 || attacker === undefined) {
    return [];
  }
  return [attacker, target];
}
