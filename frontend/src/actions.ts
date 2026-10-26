// What the player may do now, from the update's allowed intents (7.7).
// Pure: the table only offers what the server allows.
import { say, t } from './i18n';
import type { Block, Intent, ModeView, Msg, MsgNames, OptionView, OutcomeView, Prompt, Update } from './protocol';

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
  throw: 9,
} as const;

export const PromptKind = { none: 0, priority: 1, discard: 2, gameOver: 3, choose: 4, order: 5, throw: 6 } as const;

// The engine actions an outcome shows, in Go's order (actions.go).
export const ActionKind = { gainCents: 0, addCounters: 14, pushLingering: 20 } as const;

// purposeTitle is a choice's title by its purpose, in Go's order, for
// a prompt that comes without a question of its own.
function purposeTitle(purpose: number): string {
  switch (purpose) {
    case 0:
      return t('choice.title.replacement');
    case 1:
      return t('choice.title.attack');
    case 2:
      return t('choice.title.monster_slot');
    case 3:
      return t('choice.title.penalty_item');
    case 4:
      return t('choice.title.penalty_loot');
    case 5:
      return t('choice.title.purchase');
    case 6:
      return t('choice.title.target');
    case 7:
      return t('choice.title.trigger_order');
    case 8:
      return t('choice.title.mode');
    case 9:
      return t('choice.title.answer');
    case 10:
      return t('choice.title.starting_item');
    case 11:
      return t('choice.title.curse');
    case 12:
      return t('choice.title.bottom');
    default:
      return t('choice.title.choose');
  }
}

// choiceTitle is the question of a choose prompt.
export function choiceTitle(w: Prompt, names: MsgNames): string {
  return w.text?.key ? say(w.text, names) : purposeTitle(w.purpose ?? 0);
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
    return [t('choice.outcome.again')];
  }
  const lines: string[] = [];
  for (const a of o.actions ?? []) {
    const card = name(a.card ?? '');
    if (a.kind === ActionKind.gainCents) {
      lines.push(t('choice.outcome.cents', { n: a.amount }));
    } else if (a.kind === ActionKind.addCounters) {
      lines.push(t('choice.outcome.level_up', { card, n: a.amount })); // Bum-Bo!
    } else if (a.kind === ActionKind.pushLingering) {
      lines.push(t('choice.outcome.roll_cents', { card, n: a.amount })); // Poker Chip
    }
  }
  if ((o.actions ?? []).length === 0) {
    lines.push(t('choice.outcome.nothing'));
  }
  for (const c of o.unused ?? []) {
    lines.push(t('choice.outcome.unused', { card: name(c) }));
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
export const Purpose = { attackTarget: 1, monsterSlot: 2, penaltyItem: 3, penaltyLoot: 4, purchase: 5, target: 6, triggerOrder: 7, mode: 8, answer: 9 } as const;

export interface Actions {
  // buttons are the actions that need no card: pass.
  buttons: Action[];
  // turn holds the turn actions: icons by own name (F-11) and the big
  // End turn button (H-04).
  turn: TurnActions;
  // cards maps an object to what it can do: play it, use its abilities.
  cards: Map<number, Action[]>;
  // blocked maps a card the player can't play or use now to why (15.29).
  blocked: Map<number, string>;
  choice: Choice | null;
  discard: { count: number; from: number[]; intent: Intent } | null;
  // order puts cards back on a deck (19.16): send intent with the cards
  // top first.
  order: { title: string; intent: Intent } | null;
  // throw throws the die of a roll (15.61); the prompt's text names it.
  throw: { label: string; intent: Intent } | null;
}

// TurnActions are the declare intents and ending the turn, when allowed.
export interface TurnActions {
  attack?: Intent;
  buy?: Intent;
  endTurn?: Intent;
  // idle: nothing is left to do but pass or end the turn (H-20).
  idle?: boolean;
  // why says why each turn action that is not allowed can't be (15.29).
  why?: { attack?: string; buy?: string; endTurn?: string };
  // left names what the player may still do besides passing and ending
  // the turn, for the question before ending it (V-46).
  left?: string[];
}

// leftLabels name the plays End turn's question lists, in its order.
const leftLabels = (): Map<number, string> =>
  new Map([
    [IntentKind.attack, t('turn.left.attack')],
    [IntentKind.purchase, t('turn.left.buy')],
    [IntentKind.playLoot, t('turn.left.loot')],
    [IntentKind.activate, t('turn.left.abilities')],
  ]);

// actionsFor lists the viewer's actions. names names the cards and
// players in the server's lines (ADR 021).
export function actionsFor(u: Update, names: MsgNames): Actions {
  const why = turnWhy(u.blocked ?? [], names);
  const out: Actions = {
    buttons: [],
    turn: Object.keys(why).length > 0 ? { why } : {},
    cards: new Map(),
    blocked: blockedCards(u.blocked ?? [], names),
    choice: null,
    discard: null,
    order: null,
    throw: null,
  };
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
    const aimed = aim ? { ...action, targets: aim.targets, ...(aim.reason?.key ? { reason: say(aim.reason, names) } : {}) } : action;
    out.cards.set(id, [...(out.cards.get(id) ?? []), aimed]);
  };
  for (const [i, a] of allowed.entries()) {
    const id = a.objects?.[0] ?? 0;
    if (a.kind === IntentKind.pass) {
      out.buttons.push({ label: t('stack.pass'), intent: a });
    } else if (a.kind === IntentKind.endTurn) {
      out.turn.endTurn = a;
    } else if (a.kind === IntentKind.attack) {
      out.turn.attack = a;
    } else if (a.kind === IntentKind.purchase) {
      out.turn.buy = a;
    } else if (a.kind === IntentKind.playLoot) {
      add(id, { label: t('actions.play'), intent: a }, i);
    } else if (a.kind === IntentKind.activate) {
      const n = a.choice ?? 0;
      const line = lines.get(id)?.at(n);
      add(id, {
        label: (abilities.get(id) ?? 0) > 1 ? t('actions.use_n', { n: n + 1 }) : t('actions.use'),
        intent: a,
        ...(line ? { line } : {}),
      }, i);
    } else if (a.kind === IntentKind.discard) {
      out.discard = { count: u.view.waiting.count ?? 0, from: a.objects ?? [], intent: a };
    } else if (a.kind === IntentKind.order) {
      out.order = { title: choiceTitle(u.view.waiting, names), intent: a };
    } else if (a.kind === IntentKind.throw) {
      out.throw = { label: say(u.view.waiting.text, names), intent: a };
    }
  }
  if (out.turn.endTurn && allowed.every((a) => a.kind === IntentKind.endTurn || a.kind === IntentKind.pass)) {
    out.turn.idle = true;
  } else if (out.turn.endTurn) {
    out.turn.left = [...leftLabels()].filter(([kind]) => allowed.some((a) => a.kind === kind)).map(([, label]) => label);
  }
  const choose = allowed.filter((a) => a.kind === IntentKind.choose);
  if (choose.length > 0) {
    const w = u.view.waiting;
    out.choice = {
      title: choiceTitle(w, names),
      purpose: w.purpose ?? 0,
      options: choose.map((intent) => {
        const i = intent.choice ?? 0;
        const option = w.options?.at(i);
        const label = option?.key ? say(option, names) : t('choice.option_n', { n: i + 1 });
        const target = u.view.options?.at(i);
        const card = optionCard(option);
        const trigger = w.purpose === Purpose.triggerOrder ? triggerLabel(option) : undefined;
        return {
          label,
          intent,
          ...(card ? { card } : {}),
          ...trigger,
          ...(target && (target.object || target.is_player || target.deck || target.stack_id) ? { target } : {}),
          ...(target?.outcome ? { outcome: target.outcome } : {}),
        };
      }),
    };
  }
  return out;
}

// modeLines is a mode prompt as card lines, every mode in card order
// (15.58): a mode with no valid target now has no intent. cancel is the
// prompt's option that is no mode.
export function modeLines(choice: Choice, modes: ModeView[]): { lines: { text: string; intent?: Intent }[]; cancel?: Intent } {
  const pick = (i: number) => choice.options.find((o) => (o.intent.choice ?? 0) === i)?.intent;
  const lines = modes.map((m) => {
    const intent = m.option >= 0 ? pick(m.option) : undefined;
    return { text: m.text, ...(intent ? { intent } : {}) };
  });
  const cancel = choice.options.find((o) => !modes.some((m) => m.option === (o.intent.choice ?? 0)))?.intent;
  return { lines, ...(cancel ? { cancel } : {}) };
}

// blockedCards maps each blocked card to its reason (15.29). A card
// with two blocked abilities keeps the first one's.
function blockedCards(blocked: Block[], names: MsgNames): Map<number, string> {
  const out = new Map<number, string>();
  for (const b of blocked) {
    if ((b.kind === IntentKind.playLoot || b.kind === IntentKind.activate) && !out.has(b.object)) {
      out.set(b.object, say(b.reason, names));
    }
  }
  return out;
}

// turnWhy is why attack, buy, and end turn can't be used (15.29).
function turnWhy(blocked: Block[], names: MsgNames): NonNullable<TurnActions['why']> {
  const reason = (kind: number) => {
    const b = blocked.find((x) => x.kind === kind);
    return b ? say(b.reason, names) : undefined;
  };
  const attack = reason(IntentKind.attack);
  const buy = reason(IntentKind.purchase);
  const endTurn = reason(IntentKind.endTurn);
  return { ...(attack ? { attack } : {}), ...(buy ? { buy } : {}), ...(endTurn ? { endTurn } : {}) };
}

// optionCard is the card an option is, if it is one ("ask.card").
function optionCard(m?: Msg): string | undefined {
  const v = m?.key === 'ask.card' ? m.args?.card : undefined;
  return typeof v === 'object' && 'card' in v ? v.card : undefined;
}

// triggerLabel splits a trigger order option, an ability ("ask.ability"),
// into its card and the ability's words (19.19).
function triggerLabel(m?: Msg): { card: string; line: string } | undefined {
  const card = m?.key === 'ask.ability' ? m.args?.card : undefined;
  const text = m?.args?.text;
  return typeof card === 'object' && 'card' in card && typeof text === 'string' ? { card: card.card, line: text } : undefined;
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
