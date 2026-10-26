// Words for the stack and for why the game waits (7.8). Pure, so it is
// tested without a browser.
import { say, t } from './i18n';
import type { GameView, Seat, StackItem } from './protocol';
import { oddsLine, rollChanges, type Name } from './stats';

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

// actionStep and endStep are turn steps, in Go's order
// (pkg/rules_engine/turn.go).
const actionStep = 4;
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

// seatNames names players by their seats' nicknames and cards by card.
export function seatNames(seats: readonly Seat[] | null | undefined, card: (ref: string) => string): Names {
  return {
    player: (p) => seats?.find((s) => s.seat === p)?.name ?? t('table.player_n', { n: p + 1 }),
    card,
  };
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

// partsOf splits a line from t() at its placeholders (ADR 021): each
// {name} given in parts becomes that part (or parts), the words around
// it plain parts. A line keeps its names in their colors in every
// language.
export function partsOf(text: string, parts: Record<string, LinePart | LinePart[]>): LinePart[] {
  const named = new Map(Object.entries(parts));
  return text.split(/(\{\w+\})/).flatMap((piece): LinePart[] => {
    const part = piece.startsWith('{') ? named.get(piece.slice(1, -1)) : undefined;
    if (part) {
      return Array.isArray(part) ? part : [part];
    }
    return piece === '' ? [] : [{ text: piece }];
  });
}

// stackParts says what one stack item is, in parts that keep who is who.
// rolled, if given, stands for a roll's number: the history shows how a
// change turned it into the result (history.ts).
export function stackParts(it: StackItem, view: GameView, n: Names, rolled?: string): LinePart[] {
  const who: LinePart = { text: n.player(it.controller), role: 'player', player: it.controller };
  const card: LinePart = { text: it.card ? n.card(it.card) : '', role: 'card' };
  const targetCard = it.target_card ?? findCard(view, it.target.object ?? 0);
  const target: LinePart = it.target.is_player
    ? { text: n.player(it.target.player), role: 'target', player: it.target.player }
    : { text: targetCard ? n.card(targetCard) : t('stack.a_card'), role: 'target' };
  // An attack's items say so first: "Attack: Ann rolls 4 against Gaper".
  const attack = (line: LinePart[]) => (it.attack ? partsOf(t('stack.attack'), { line }) : line);
  // dealer is who deals damage: the player, or the monster that hits back.
  // Damage a player controls against themselves is "take damage" (Boom
  // Fly, Secret Room): its card deals it (15.47, ADR 011).
  const sourceCard = it.source_card ?? findCard(view, it.source);
  const self = it.target.is_player && it.target.player === it.controller;
  const amount = it.amount ?? 0;
  switch (it.kind) {
    case StackKind.loot:
      return partsOf(t('stack.loot'), { player: who, card });
    case StackKind.ability:
      return partsOf(t('stack.ability'), { card, player: who });
    case StackKind.trigger:
      return partsOf(t('stack.trigger'), { card });
    case StackKind.roll: {
      // A thrown roll names what will change it, and the result (15.94).
      const changes = rollChanges(it.sources, sourceName(n));
      const result = changes ? ` ${t('stack.roll_result', { changes, n: it.final ?? it.roll ?? 0 })}` : '';
      const roll = rolled ?? (it.roll ? String(it.roll) : '');
      const parts = rollParts(it, n, roll, who, target, attack, !!targetCard);
      return result ? [...parts, { text: result }] : parts;
    }
    case StackKind.damage: {
      if (it.controller >= 0 && !self) {
        return attack(partsOf(t('stack.damage', { n: amount }), { player: who, target }));
      }
      if (sourceCard) {
        return attack(partsOf(t('stack.damage_by_card', { n: amount }), { card: { text: n.card(sourceCard), role: 'card' }, target }));
      }
      if (self) {
        // No card to name: the player takes it.
        return attack(partsOf(t('stack.takes_damage', { n: amount }), { target }));
      }
      return attack(partsOf(t('stack.damage_to', { n: amount }), { target }));
    }
    case StackKind.death:
      return partsOf(t('stack.dies'), { target });
    case StackKind.penalty:
      return partsOf(t('stack.penalty'), { player: { text: n.player(it.target.player), role: 'target', player: it.target.player } });
    case StackKind.deathStep:
      return partsOf(t('stack.killed'), { card });
    default:
      return [card.text ? card : { text: t('stack.something') }];
  }
}

// rollParts says who rolls for what: an attack roll names the monster
// it is against, a card's roll the card, other rolls their label. roll
// is the number, '' before the throw.
function rollParts(it: StackItem, n: Names, roll: string, who: LinePart, target: LinePart, attack: (line: LinePart[]) => LinePart[], against: boolean): LinePart[] {
  const rolled: Record<string, LinePart> = { roll: { text: roll } };
  if (it.attack && against) {
    return attack(partsOf(roll ? t('stack.roll_against_n') : t('stack.roll_against'), { player: who, ...rolled, target }));
  }
  if (it.roll_for.card) {
    const card: LinePart = { text: n.card(it.roll_for.card), role: 'card' };
    return partsOf(roll ? t('stack.roll_card_n') : t('stack.roll_card'), { player: who, ...rolled, card });
  }
  if (it.label?.key) {
    // The label comes from the server as a key (ADR 021).
    const label: LinePart = { text: say(it.label, n) };
    return attack(partsOf(roll ? t('stack.roll_label_n') : t('stack.roll_label'), { player: who, ...rolled, label }));
  }
  return attack(partsOf(roll ? t('stack.roll_n') : t('stack.roll'), { player: who, ...rolled }));
}

// sourceName names a stat source's card for rollChanges and oddsLine.
function sourceName(n: Names): Name {
  return (ref) => (ref ? n.card(ref) : t('stats.hidden_card'));
}

// attackTarget names what an attack aims at: an attacked player by name,
// not their character (R-ATK-20), else the monster.
export function attackTarget(view: GameView, target: number, n: Names): LinePart {
  const victim = view.players.find((p) => p.character.id === target);
  if (victim) {
    return { text: n.player(victim.id), role: 'target', player: victim.id };
  }
  const ref = findCard(view, target);
  return { text: ref ? n.card(ref) : t('stack.a_monster'), role: 'target' };
}

// declaration says what the active player declared that is still open
// (15.91, V-51): an attack, by its target once aimed, or a buy. Neither
// is a stack item (R-ATK-02, R-SHOP-02); the stack panel shows it in a
// row of its own, on top. It goes when the attack or the buy ends.
// A declared end of turn too, also one a card declared, until the end
// phase starts (19.74, R-TURN-08, R-TURN-10); a death end is none.
export function declaration(view: GameView, n: Names): LinePart[] | null {
  if (view.over || view.first_roll) {
    return null;
  }
  const { active } = view.turn;
  const who: LinePart = { text: n.player(active), role: 'player', player: active };
  const { attack } = view;
  if (attack) {
    const target = attack.target ?? 0;
    if (target === 0) {
      return partsOf(attack.deck ? t('stack.declared.attack_deck') : t('stack.declared.attack'), { player: who });
    }
    return partsOf(t('stack.declared.attack_target'), { player: who, target: attackTarget(view, target, n) });
  }
  if (view.buying) {
    return partsOf(t('stack.declared.buy'), { player: who });
  }
  const { end_declared: ending, step } = view.turn;
  return ending && step === actionStep ? partsOf(t('stack.declared.end_turn'), { player: who }) : null;
}

// StackRowData is one row of the stack panel: the open declaration, or
// a stack item. key is stable while the row stays, so the row moves
// instead of being drawn again (15.51).
export type StackRowData = { key: string; declared: LinePart[]; odds?: string } | { key: string; item: StackItem };

// stackRows are the stack panel's rows, top first: the open
// declaration (15.91), then the items, newest first (S-01).
export function stackRows(view: GameView, n: Names): StackRowData[] {
  const declared = declaration(view, n);
  const items = [...(view.stack ?? [])].reverse().map((item) => ({ key: `item-${String(item.id)}`, item }));
  if (!declared) {
    return items;
  }
  const odds = attackOdds(view, n);
  return [{ key: 'declared', declared, ...(odds ? { odds } : {}) }, ...items];
}

// attackOdds is the open attack's fight line, once aimed (15.94): what
// the next roll needs and the ATK of a hit, with what changes them.
export function attackOdds(view: GameView, n: Names): string | undefined {
  const odds = view.attack?.odds;
  if (!odds || view.over) {
    return undefined;
  }
  const attacker = view.players.find((p) => p.id === view.turn.active);
  return oddsLine(odds, attacker?.printed_atk ?? odds.atk, sourceName(n));
}

// stackLine says what one stack item is, as plain text; rolled as in
// stackParts.
export function stackLine(it: StackItem, view: GameView, n: Names, rolled?: string): string {
  return stackParts(it, view, n, rolled)
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
  const name = (p: number) => (p === you ? t('stack.you') : n.player(p));
  return { responding: (view.responding ?? []).map(name), passed: (view.passed ?? []).map(name) };
}

// secondsLeft is the response timer, whole seconds, never negative.
export function secondsLeft(deadline: number, now: number): number {
  return Math.max(0, Math.ceil((deadline - now) / 1000));
}
