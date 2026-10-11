// Words for the stack and for why the game waits (7.8). Pure, so it is
// tested without a browser.
import { PromptKind } from './actions';
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

// stackLine says what one stack item is.
export function stackLine(it: StackItem, view: GameView, n: Names): string {
  const who = n.player(it.controller);
  const card = it.card ? n.card(it.card) : '';
  const target = it.target.is_player
    ? n.player(it.target.player)
    : n.card(findCard(view, it.target.object ?? 0) ?? 'a card');
  const attack = it.attack ? 'Attack: ' : '';
  switch (it.kind) {
    case StackKind.loot:
      return `${who} plays ${card}`;
    case StackKind.ability:
      return `${card} (${who})`;
    case StackKind.trigger:
      return `${card} triggers`;
    case StackKind.roll: {
      const label = it.label ? ` for ${it.label}` : '';
      return `${attack}${who} rolls${it.roll ? ` ${String(it.roll)}` : ''}${label}`;
    }
    case StackKind.damage:
      return `${attack}${String(it.amount ?? 0)} damage to ${target}`;
    case StackKind.death:
      return `${target} dies`;
    case StackKind.penalty:
      return `${n.player(it.target.player)} pays the death penalty`;
    case StackKind.deathStep:
      return `${card} is killed`;
    default:
      return card || 'Something';
  }
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

// secondsLeft is the response timer, whole seconds, never negative.
export function secondsLeft(deadline: number, now: number): number {
  return Math.max(0, Math.ceil((deadline - now) / 1000));
}
