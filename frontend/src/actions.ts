// What the player may do now, from the update's allowed intents (7.7).
// Pure: the table only offers what the server allows.
import type { Intent, OptionView, Prompt, Update } from './protocol';

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
} as const;

export const PromptKind = { none: 0, priority: 1, discard: 2, gameOver: 3, choose: 4 } as const;

// Choice titles, by purpose, in Go's order.
const purposeTitles = [
  'Choose which happens first',
  'Choose what to attack',
  'Choose a slot for the monster',
  'Choose an item to lose',
  'Choose a loot card to lose',
  'Choose what to buy',
  'Choose a target',
  'Choose which goes on the stack first',
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
}

export interface ChoiceOption extends Action {
  // card is the option's card, when the label names one.
  card?: string;
  // target is what the option points at on the table, if anything.
  target?: OptionView;
}

export interface Choice {
  title: string;
  purpose: number;
  options: ChoiceOption[];
}

// Purpose.purchase: a "deck" option means the treasure deck, otherwise
// the monster deck.
export const Purpose = { attackTarget: 1, purchase: 5, answer: 9 } as const;

export interface Actions {
  // buttons are the actions that need no card: pass, end turn, attack, buy.
  buttons: Action[];
  // cards maps an object to what it can do: play it, use its abilities.
  cards: Map<number, Action[]>;
  choice: Choice | null;
  discard: { count: number; from: number[]; intent: Intent } | null;
}

const buttonLabels = new Map<number, string>([
  [IntentKind.pass, 'Pass'],
  [IntentKind.endTurn, 'End turn'],
  [IntentKind.attack, 'Attack'],
  [IntentKind.purchase, 'Buy'],
]);

// actionsFor lists the viewer's actions. isCard tells whether an option
// label is a card ref.
export function actionsFor(u: Update, isCard: (label: string) => boolean): Actions {
  const out: Actions = { buttons: [], cards: new Map(), choice: null, discard: null };
  const allowed = u.allowed ?? [];
  const abilities = new Map<number, number>();
  for (const a of allowed) {
    if (a.kind === IntentKind.activate) {
      const id = a.objects?.[0] ?? 0;
      abilities.set(id, (abilities.get(id) ?? 0) + 1);
    }
  }
  const add = (id: number, action: Action) => {
    out.cards.set(id, [...(out.cards.get(id) ?? []), action]);
  };
  for (const a of allowed) {
    const id = a.objects?.[0] ?? 0;
    const button = buttonLabels.get(a.kind);
    if (button) {
      out.buttons.push({ label: button, intent: a });
    } else if (a.kind === IntentKind.playLoot) {
      add(id, { label: 'Play', intent: a });
    } else if (a.kind === IntentKind.activate) {
      const n = a.choice ?? 0;
      add(id, { label: (abilities.get(id) ?? 0) > 1 ? `Use ability ${String(n + 1)}` : 'Use', intent: a });
    } else if (a.kind === IntentKind.discard) {
      out.discard = { count: u.view.waiting.count ?? 0, from: a.objects ?? [], intent: a };
    }
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
        return {
          label,
          intent,
          ...(isCard(label) ? { card: label } : {}),
          ...(target && (target.object || target.is_player || target.deck || target.stack_id) ? { target } : {}),
        };
      }),
    };
  }
  return out;
}

// discardIntent is the discard of exactly the chosen hand cards.
export function discardIntent(base: Intent, chosen: number[]): Intent {
  return { ...base, objects: chosen };
}
