// What the player may do now, from the update's allowed intents (7.7).
// Pure: the table only offers what the server allows.
import type { Intent, OptionView, OutcomeView, Prompt, Update } from './protocol';

// Engine enums, in Go's order (pkg/rules_engine/prompt.go, choice.go).
export const IntentKind = {
  pass: 0,
  endTurn: 1,
  discard: 2,
  playLoot: 3,
  choose: 4,
  attack: 5,
  purchase: 6,
  activate: 7,
  order: 8,
} as const;

export const PromptKind = { none: 0, priority: 1, discard: 2, gameOver: 3, choose: 4, order: 5 } as const;

// The engine actions an outcome shows, in Go's order (actions.go).
export const ActionKind = { gainCents: 0, addCounters: 14, pushLingering: 20 } as const;

// Choice titles, by purpose, in Go's order.
const purposeTitles = [
  'Choose which card applies first',
  'Choose what to attack',
  'Choose a slot for the monster',
  'Choose an item to lose',
  'Choose a loot card to lose',
  'Choose what to buy',
  'Choose a target',
  'Choose which goes on the stack first; it resolves last',
  'Choose one',
  'Answer',
  'Choose your starting item',
  'Choose who gains the curse',
  'Choose the card for the very bottom',
];

// choiceTitle is the question of a choose prompt.
export function choiceTitle(w: Prompt): string {
  return w.text ?? purposeTitles.at(w.purpose ?? 0) ?? 'Choose';
}

export interface Action {
  label: string;
  intent: Intent;
  // line is the words of the ability it uses, when the card has two or
  // more to pick from (C-05).
  line?: string;
  // targets are what the action asks to target first, if it aims (15.25).
  targets?: OptionView[];
  // reason says what it may target, for a click on anything else (15.33).
  reason?: string;
}

export interface ChoiceOption extends Action {
  // card is the option's card, when the label names one.
  card?: string;
  // target is what the option points at on the table, if anything.
  target?: OptionView;
  // outcome is what picking it leads to (19.17).
  outcome?: OutcomeView;
}

// outcomeLines says what an order of replacement effects gives, for the
// chooser (19.17). name gives a card's title.
export function outcomeLines(o: OutcomeView, name: (card: string) => string): string[] {
  if (o.again) {
    return ['Then you pick the next card'];
  }
  const lines: string[] = [];
  for (const a of o.actions ?? []) {
    const card = name(a.card ?? '');
    if (a.kind === ActionKind.gainCents) {
      lines.push(`You gain ${String(a.amount)}¢`);
    } else if (a.kind === ActionKind.addCounters) {
      lines.push(`${card} levels up by ${String(a.amount)}`); // Bum-Bo!
    } else if (a.kind === ActionKind.pushLingering) {
      lines.push(`${card} rolls for ${String(a.amount)}¢`); // Poker Chip
    }
  }
  if ((o.actions ?? []).length === 0) {
    lines.push('Nothing happens');
  }
  for (const c of o.unused ?? []) {
    lines.push(`${name(c)} does nothing`);
  }
  return lines;
}

export interface Choice {
  title: string;
  purpose: number;
  options: ChoiceOption[];
}

// Purpose.purchase: a "deck" option means the treasure deck, otherwise
// the monster deck.
export const Purpose = { attackTarget: 1, monsterSlot: 2, purchase: 5, target: 6, triggerOrder: 7, answer: 9 } as const;

export interface Actions {
  // buttons are the actions that need no card: pass.
  buttons: Action[];
  // turn holds the turn actions: icons by own name (F-11) and the big
  // End turn button (H-04).
  turn: TurnActions;
  // cards maps an object to what it can do: play it, use its abilities.
  cards: Map<number, Action[]>;
  choice: Choice | null;
  discard: { count: number; from: number[]; intent: Intent } | null;
  // order puts cards back on a deck (19.16): send intent with the cards
  // top first.
  order: { title: string; intent: Intent } | null;
}

// TurnActions are the declare intents and ending the turn, when allowed.
export interface TurnActions {
  attack?: Intent;
  buy?: Intent;
  endTurn?: Intent;
  // idle: nothing is left to do but pass or end the turn (H-20).
  idle?: boolean;
}

const buttonLabels = new Map<number, string>([[IntentKind.pass, 'Pass']]);

// actionsFor lists the viewer's actions. isCard tells whether an option
// label is a card ref.
export function actionsFor(u: Update, isCard: (label: string) => boolean): Actions {
  const out: Actions = { buttons: [], turn: {}, cards: new Map(), choice: null, discard: null, order: null };
  const allowed = u.allowed ?? [];
  const abilities = new Map<number, number>();
  for (const a of allowed) {
    if (a.kind === IntentKind.activate) {
      const id = a.objects?.[0] ?? 0;
      abilities.set(id, (abilities.get(id) ?? 0) + 1);
    }
  }
  const lines = new Map<number, string[]>();
  for (const c of u.view.players.flatMap((p) => [p.character, ...(p.in_play ?? [])])) {
    if (c.abilities) {
      lines.set(c.id, c.abilities);
    }
  }
  const aims = new Map((u.aims ?? []).map((a) => [a.intent, a]));
  const add = (id: number, action: Action, i: number) => {
    const aim = aims.get(i);
    const aimed = aim ? { ...action, targets: aim.targets, ...(aim.reason ? { reason: aim.reason } : {}) } : action;
    out.cards.set(id, [...(out.cards.get(id) ?? []), aimed]);
  };
  for (const [i, a] of allowed.entries()) {
    const id = a.objects?.[0] ?? 0;
    const button = buttonLabels.get(a.kind);
    if (button) {
      out.buttons.push({ label: button, intent: a });
    } else if (a.kind === IntentKind.endTurn) {
      out.turn.endTurn = a;
    } else if (a.kind === IntentKind.attack) {
      out.turn.attack = a;
    } else if (a.kind === IntentKind.purchase) {
      out.turn.buy = a;
    } else if (a.kind === IntentKind.playLoot) {
      add(id, { label: 'Play', intent: a }, i);
    } else if (a.kind === IntentKind.activate) {
      const n = a.choice ?? 0;
      const line = lines.get(id)?.at(n);
      add(id, {
        label: (abilities.get(id) ?? 0) > 1 ? `Use ability ${String(n + 1)}` : 'Use',
        intent: a,
        ...(line ? { line } : {}),
      }, i);
    } else if (a.kind === IntentKind.discard) {
      out.discard = { count: u.view.waiting.count ?? 0, from: a.objects ?? [], intent: a };
    } else if (a.kind === IntentKind.order) {
      out.order = { title: choiceTitle(u.view.waiting), intent: a };
    }
  }
  if (out.turn.endTurn && allowed.every((a) => a.kind === IntentKind.endTurn || a.kind === IntentKind.pass)) {
    out.turn.idle = true;
  }
  const choose = allowed.filter((a) => a.kind === IntentKind.choose);
  if (choose.length > 0) {
    const w = u.view.waiting;
    out.choice = {
      title: choiceTitle(w),
      purpose: w.purpose ?? 0,
      options: choose.map((intent) => {
        const i = intent.choice ?? 0;
        const label = w.options?.at(i) ?? `Option ${String(i + 1)}`;
        const target = u.view.options?.at(i);
        const trigger = w.purpose === Purpose.triggerOrder ? triggerLabel(label, isCard) : undefined;
        return {
          label,
          intent,
          ...(isCard(label) ? { card: label } : {}),
          ...trigger,
          ...(target && (target.object || target.is_player || target.deck || target.stack_id) ? { target } : {}),
          ...(target?.outcome ? { outcome: target.outcome } : {}),
        };
      }),
    };
  }
  return out;
}

// triggerLabel splits a trigger order label, "card_ref: ability text",
// into its card and the ability's words (19.19).
function triggerLabel(label: string, isCard: (label: string) => boolean): { card: string; line: string } | undefined {
  const at = label.indexOf(': ');
  const card = label.slice(0, at);
  return at > 0 && isCard(card) ? { card, line: label.slice(at + 2) } : undefined;
}

// orderIntent puts the cards back in the given order, top first.
export function orderIntent(base: Intent, top: number[]): Intent {
  return { ...base, objects: top };
}

// discardIntent is the discard of exactly the chosen hand cards.
export function discardIntent(base: Intent, chosen: number[]): Intent {
  return { ...base, objects: chosen };
}

// sameTarget: two option views point at the same thing on the table.
export function sameTarget(a: OptionView, b: OptionView): boolean {
  if (a.is_player || b.is_player) {
    return !!a.is_player && !!b.is_player && a.player === b.player;
  }
  if (a.stack_id || b.stack_id) {
    return a.stack_id === b.stack_id;
  }
  return !!a.object && a.object === b.object;
}

// aimedOption is the target prompt's option for a card dropped on target
// (15.25), if the prompt asks for a target and offers it.
export function aimedOption(choice: Choice | null, target: OptionView): Intent | undefined {
  if (choice?.purpose !== Purpose.target) {
    return undefined;
  }
  return choice.options.find((o) => o.target && sameTarget(o.target, target))?.intent;
}
