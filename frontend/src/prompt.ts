// The prompt banner (15.28, V-20): what the game waits for, in a few
// words. Pure, so every prompt is tested without a browser.
import { choiceTitle, PromptKind, Purpose, type Actions } from './actions';
import { say, t } from './i18n';
import type { GameView } from './protocol';
import { attackTarget, findCard, partsOf, StackKind, stackLine, type LinePart, type Names } from './stack';
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
      title: t('banner.aim'),
      hint: card ? t('banner.aim_for', { card: n.card(card) }) : t('banner.aim_cancel'),
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
  return { kind: ask.kind, mine: false, player: w.player, title: t('banner.waiting_for', { player: n.player(w.player) }), hint: ask.other };
}

// asked is the prompt in words: title and hint for the asked player,
// other for everyone else.
function asked(view: GameView, n: Names): { kind: BannerKind; title: string; hint?: string; other: string } | null {
  const { waiting: w, turn } = view;
  switch (w.kind) {
    case PromptKind.priority: {
      const top = (view.stack ?? []).at(-1);
      if (top) {
        const title = t('banner.respond');
        return { kind: 'respond', title, hint: t('banner.on_stack', { line: stackLine(top, view, n) }), other: title };
      }
      if (turn.active === w.player) {
        return { kind: 'turn', title: t('banner.your_turn'), other: t('banner.their_turn') };
      }
      return { kind: 'act', title: t('banner.act'), hint: t('banner.act_hint'), other: t('banner.act_other') };
    }
    case PromptKind.choose: {
      const title = choiceTitle(w, n);
      return { kind: 'choose', title, other: title };
    }
    case PromptKind.discard: {
      const count = w.count ?? 0;
      const title = t('discard.title', { n: count });
      return { kind: 'discard', title, other: title };
    }
    case PromptKind.order: {
      const title = choiceTitle(w, n);
      return { kind: 'order', title, other: title };
    }
    case PromptKind.throw: {
      const title = t('banner.throw');
      const label = say(w.text, n);
      return label
        ? { kind: 'throw', title, hint: t('banner.throw_for', { label }), other: t('banner.throw_for_other', { label }) }
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
    return t('banner.idle');
  }
  const can = [
    cards.size > 0 ? t('banner.can.play') : '',
    turn.attack ? t('banner.can.attack') : '',
    turn.buy ? t('banner.can.buy') : '',
    turn.endTurn ? t('banner.can.end_turn') : '',
  ].filter(Boolean);
  if (can.length === 0) {
    return t('banner.pass');
  }
  const last = can.pop() ?? '';
  const words = can.length > 0 ? t('banner.can.or', { list: can.join(', '), last }) : last;
  return words.charAt(0).toLocaleUpperCase() + words.slice(1);
}

// picking: the player picks a target on the table now, so the cards
// that answer glow and the rest dim (15.28). A choice counts once one
// of its options is on the table.
export function picking(actions: Actions, aiming: Aiming | null): boolean {
  return aiming !== null || (actions.choice?.options.some((o) => o.target) ?? false);
}

// myMove: the board glows green, as the game waits for the viewer
// (15.95, V-52). Not while they aim or pick a target: then only the
// answers glow.
export function myMove(b: Banner | null, picks: boolean): boolean {
  return b?.mine === true && b.kind !== 'aim' && !picks;
}

// askedSeat is the player the game waits for, marked on their seat for
// everyone (15.95); -1 while it waits for nobody.
export function askedSeat(view: GameView): number {
  return view.over || view.waiting.kind === PromptKind.none ? -1 : view.waiting.player;
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
  // Before turn 1 nobody is active: the first-player roll (15.43).
  if (view.over || view.first_roll) {
    return null;
  }
  const { active } = view.turn;
  const who: LinePart = { text: n.player(active), role: 'player', player: active };
  const card = (ref: string | undefined): LinePart => ({ text: ref ? n.card(ref) : t('stack.a_card'), role: 'card' });
  const w = view.waiting;
  const top = (view.stack ?? []).at(-1);
  if (top?.controller === active && top.kind === StackKind.loot) {
    return partsOf(t('activity.plays'), { player: who, card: card(top.card) });
  }
  if (top?.controller === active && top.kind === StackKind.ability) {
    return partsOf(t('activity.uses'), { player: who, card: card(top.card ?? top.source_card ?? findCard(view, top.source)) });
  }
  const penalty =
    (top?.kind === StackKind.penalty && top.target.player === active) ||
    (w.player === active && w.kind === PromptKind.choose && (w.purpose === Purpose.penaltyItem || w.purpose === Purpose.penaltyLoot));
  if (penalty) {
    return partsOf(t('activity.penalty'), { player: who });
  }
  const { attack } = view;
  if (attack) {
    const target = attack.target ?? 0;
    if (target === 0) {
      return partsOf(attack.deck ? t('activity.attack_deck') : t('activity.attack'), { player: who });
    }
    return partsOf(t('activity.attack_target'), { player: who, target: attackTarget(view, target, n) });
  }
  if (view.buying) {
    return partsOf(t('activity.buying'), { player: who });
  }
  if (w.player === active && w.kind !== PromptKind.none) {
    return partsOf(t('activity.thinking'), { player: who });
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
