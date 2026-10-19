// Words for the stack and for why the game waits (7.8). Pure, so it is
// tested without a browser.
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

// endStep is the end of turn step, in Go's order (pkg/rules_engine/turn.go).
const endStep = 5;

// waitTurn is "Wait for my turn" (N-16): on, offered on another player's
// turn, or not offered once the turn before this player's is ending.
export function waitTurn(view: GameView, you: number, on: boolean): 'offer' | 'on' | 'none' {
  const { turn, over, players } = view;
  if (on) {
    return 'on';
  }
  if (over || you < 0 || turn.active === you) {
    return 'none';
  }
  const before = (turn.active + 1) % players.length === you;
  const ending = (turn.end_declared ?? false) || (turn.death_end ?? false) || turn.step >= endStep;
  return before && ending ? 'none' : 'offer';
}

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
  const targetCard = it.target_card ?? findCard(view, it.target.object ?? 0);
  const target: LinePart = it.target.is_player
    ? { text: n.player(it.target.player), role: 'target', player: it.target.player }
    : { text: targetCard ? n.card(targetCard) : 'a card', role: 'target' };
  const attack: LinePart[] = it.attack ? [{ text: 'Attack: ' }] : [];
  // dealer is who deals damage: the player, or the monster that hits back.
  // Damage a player controls against themselves is "take damage" (Boom
  // Fly, Secret Room): its card deals it (15.47, ADR 011).
  const sourceCard = it.source_card ?? findCard(view, it.source);
  const self = it.target.is_player && it.target.player === it.controller;
  const byCard: LinePart[] = sourceCard ? [{ text: n.card(sourceCard), role: 'card' }, { text: ' deals ' }] : [];
  const dealer: LinePart[] = it.controller >= 0 && !self ? [who, { text: ' deals ' }] : byCard;
  switch (it.kind) {
    case StackKind.loot:
      return [who, { text: ' plays ' }, card];
    case StackKind.ability:
      return [card, { text: ' (' }, who, { text: ')' }];
    case StackKind.trigger:
      return [card, { text: ' triggers' }];
    case StackKind.roll: {
      const rolls: LinePart = { text: ` rolls${it.roll ? ` ${String(it.roll)}` : ''}` };
      if (it.attack && targetCard) {
        // An attack roll names the monster it is against.
        return [...attack, who, rolls, { text: ' against ' }, target];
      }
      // A card's roll names the card; other rolls their label.
      if (it.roll_for.card) {
        return [who, { text: `${rolls.text} for ` }, { text: n.card(it.roll_for.card), role: 'card' }];
      }
      return [...attack, who, { text: `${rolls.text}${it.label ? ` for ${it.label}` : ''}` }];
    }
    case StackKind.damage: {
      if (dealer.length === 0 && self) {
        // No card to name: the player takes it.
        return [...attack, target, { text: ` takes ${String(it.amount ?? 0)} damage` }];
      }
      return [...attack, ...dealer, { text: `${String(it.amount ?? 0)} damage to ` }, target];
    }
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
// card, the card it comes from, or else the card it targets: the monster
// an attack roll is against (15.22).
export function stackCard(it: StackItem, view: GameView): string | undefined {
  return (
    it.card ??
    it.source_card ??
    findCard(view, it.source) ??
    it.target_card ??
    (it.target.is_player ? undefined : findCard(view, it.target.object ?? 0))
  );
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

// playedCard is the newest loot card another player has on the stack:
// it shows large by the stack until it resolves (15.31, H-10, T-16).
export function playedCard(view: GameView, you: number): StackItem | undefined {
  return (view.stack ?? []).findLast((it) => it.kind === StackKind.loot && it.controller !== you && !!it.card);
}

// StackWait is whom the stack waits for (V-11): the players who may still
// respond, the waiting one first, and those who passed. Other prompts
// get the banner's words (15.28, stackPrompt).
export interface StackWait {
  responding: string[];
  passed: string[];
}

// stackWait names whom the stack waits for, as the server says.
export function stackWait(view: GameView, you: number, n: Names): StackWait {
  const name = (p: number) => (p === you ? 'You' : n.player(p));
  return { responding: (view.responding ?? []).map(name), passed: (view.passed ?? []).map(name) };
}

// secondsLeft is the response timer, whole seconds, never negative.
export function secondsLeft(deadline: number, now: number): number {
  return Math.max(0, Math.ceil((deadline - now) / 1000));
}
