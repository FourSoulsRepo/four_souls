// Words for the stack and for why the game waits (7.8). Pure, so it is
// tested without a browser.
import { choiceTitle, PromptKind } from './actions';
import type { GameView, StackItem } from './protocol';

// Stack item kinds, in Go's order (pkg/rules_engine/stack.go).
export const StackKind = {
  loot: 0,
  ability: 1,
  trigger: 2,
  roll: 3,
  damage: 4,
  death: 5,
  penalty: 6,
  deathStep: 7,
} as const;

export interface Names {
  player: (p: number) => string;
  card: (ref: string) => string;
}

// findCard is the card ref of an object on the table, if it is there.
export function findCard(view: GameView, id: number): string | undefined {
  const all = [
    ...view.players.flatMap((p) => [p.character, ...(p.in_play ?? []), ...(p.hand ?? [])]),
    ...[...(view.monsters ?? []), ...(view.shop ?? [])].flatMap((s) => [...(s.covered ?? []), ...(s.top ? [s.top] : [])]),
    ...view.discards.flatMap((d) => d ?? []),
    ...(view.bonus_souls ?? []),
  ];
  return all.find((c) => c.id === id)?.card;
}

// LinePart is a piece of a stack line: plain words, a player's name in
// their color, a target in red (S-03), or a card's name.
export interface LinePart {
  text: string;
  role?: 'player' | 'target' | 'card';
  player?: number;
}

// stackParts says what one stack item is, in parts that keep who is who.
export function stackParts(it: StackItem, view: GameView, n: Names): LinePart[] {
  const who: LinePart = { text: n.player(it.controller), role: 'player', player: it.controller };
  const card: LinePart = { text: it.card ? n.card(it.card) : '', role: 'card' };
  const target: LinePart = it.target.is_player
    ? { text: n.player(it.target.player), role: 'target', player: it.target.player }
    : { text: n.card(findCard(view, it.target.object ?? 0) ?? 'a card'), role: 'target' };
  const attack: LinePart[] = it.attack ? [{ text: 'Attack: ' }] : [];
  switch (it.kind) {
    case StackKind.loot:
      return [who, { text: ' plays ' }, card];
    case StackKind.ability:
      return [card, { text: ' (' }, who, { text: ')' }];
    case StackKind.trigger:
      return [card, { text: ' triggers' }];
    case StackKind.roll: {
      const label = it.label ? ` for ${it.label}` : '';
      return [...attack, who, { text: ` rolls${it.roll ? ` ${String(it.roll)}` : ''}${label}` }];
    }
    case StackKind.damage:
      return [...attack, { text: `${String(it.amount ?? 0)} damage to ` }, target];
    case StackKind.death:
      return [target, { text: ' dies' }];
    case StackKind.penalty:
      return [
        { text: n.player(it.target.player), role: 'target', player: it.target.player },
        { text: ' pays the death penalty' },
      ];
    case StackKind.deathStep:
      return [card, { text: ' is killed' }];
    default:
      return [card.text ? card : { text: 'Something' }];
  }
}

// stackLine says what one stack item is, as plain text.
export function stackLine(it: StackItem, view: GameView, n: Names): string {
  return stackParts(it, view, n)
    .map((p) => p.text)
    .join('');
}

// StackIcon is what a stack row shows first (S-04): the card's art, a
// die for a roll, a mark for damage, or a skull for a death.
export type StackIcon = 'card' | 'roll' | 'damage' | 'death';

export function stackIcon(it: StackItem): StackIcon {
  switch (it.kind) {
    case StackKind.roll:
      return 'roll';
    case StackKind.damage:
      return 'damage';
    case StackKind.death:
    case StackKind.penalty:
      return 'death';
    default:
      return 'card';
  }
}

// stackCard is the card a stack row shows large on hover (S-06): its own
// card, or the card it comes from.
export function stackCard(it: StackItem, view: GameView): string | undefined {
  return it.card ?? findCard(view, it.source);
}

// StackBadge marks a history icon (H-08): a sword for an attack, a bolt
// for a played loot card, our spells.
export type StackBadge = 'attack' | 'spell';

export function stackBadge(it: StackItem): StackBadge | undefined {
  if (it.attack) {
    return 'attack';
  }
  return it.kind === StackKind.loot ? 'spell' : undefined;
}

// waitText says who the game waits for and why, so nobody is lost.
export function waitText(view: GameView, you: number, n: Names, choiceTitle?: string): string {
  const { waiting, turn, over, stack } = view;
  if (over) {
    return 'Game over';
  }
  const mine = waiting.player === you;
  const who = mine ? 'You' : n.player(waiting.player);
  switch (waiting.kind) {
    case PromptKind.choose:
      return mine ? (choiceTitle ?? 'Your choice') : `${who} is choosing`;
    case PromptKind.discard:
      return mine ? 'Discard down to your hand size' : `${who} is discarding`;
    case PromptKind.priority:
      if ((stack ?? []).length > 0) {
        return mine ? 'Respond, or pass' : `${who} may respond`;
      }
      if (turn.active === waiting.player) {
        return mine ? 'Your turn' : `${who}'s turn`;
      }
      return mine ? 'You may act' : `${who} may act`;
    default:
      return '';
  }
}

// StackWait is whom the stack waits for (V-11): the players who may still
// respond, the waiting one first, and those who passed; or the player who
// must answer a prompt, and what.
export interface StackWait {
  responding: string[];
  passed: string[];
  prompt?: string;
}

// stackWait names whom the stack waits for, as the server says.
export function stackWait(view: GameView, you: number, n: Names): StackWait {
  const name = (p: number) => (p === you ? 'You' : n.player(p));
  const w = view.waiting;
  const out: StackWait = { responding: (view.responding ?? []).map(name), passed: (view.passed ?? []).map(name) };
  if (w.kind === PromptKind.choose) {
    out.prompt = `${name(w.player)}: ${choiceTitle(w)}`;
  } else if (w.kind === PromptKind.discard) {
    const count = w.count ?? 0;
    out.prompt = `${name(w.player)}: Discard ${String(count)} loot card${count === 1 ? '' : 's'}`;
  }
  return out;
}

// secondsLeft is the response timer, whole seconds, never negative.
export function secondsLeft(deadline: number, now: number): number {
  return Math.max(0, Math.ceil((deadline - now) / 1000));
}
