// The match history (7.9, V-03, V-51): a line for every play, under turn
// lines: who did what. Built from the events this viewer got, so it shows
// only what the viewer may see (hidden cards stay hidden).
import { firstPlayerText, tieText } from './dice';
import type { GameEvent, GameView, Rollback, StackItem } from './protocol';
import { findCard, StackKind, stackLine, type Names } from './stack';
import { rollChanges } from './stats';

export interface Entry {
  key: number;
  turn: number;
  icon: string;
  text: string;
  // card is the card the entry is about, for its zoom.
  card?: string;
  // item is the stack item that resolved, for the history strip (15.22).
  item?: StackItem;
  // headline: a play, an ability, a buy or a kill (15.88); the strip
  // shows these (15.89).
  headline?: boolean;
  // badge marks a strip tile that is no plain stack item: a buy, or a
  // death that shows the dead card (15.89).
  badge?: 'buy' | 'kill';
  // from is the card the strip's popout shows large: the played card,
  // the dealer, the killer or the buyer (15.90).
  from?: string;
  // outcome is what the headline did, for the popout (15.90).
  outcome?: Outcome[];
  // step is the first step of the update that brought the entry: a
  // rollback to before it undoes it (ADR 020). A rollback row has its
  // own number.
  step?: number;
  // undone: a rollback undid it; the history greys it.
  undone?: boolean;
  // rollback marks the row that tells who rolled back.
  rollback?: boolean;
}

// Outcome is one result in the strip's popout (15.90): a card, with a
// red damage number or a skull, or cents as a plain mark.
export interface Outcome {
  // text names it, for screen readers and cards without a ref.
  text: string;
  card?: string;
  mark?: 'damage' | 'death' | 'cents';
  // amount: the damage, or the cents won (negative if lost).
  amount?: number;
}

// told marks an event kind that gives a history line in historyKinds.
const told = 'line';

// historyKinds says what the history does with each engine event kind
// (15.88): told gives a line; any other text is why it gives none. A kind
// missing here fails history.test.ts, so no play goes untold.
export const historyKinds: Record<string, string> = {
  activated: 'Its stack item resolves as a line, or is cancelled.',
  attack_declared: told,
  attack_ended: 'The rolls, the damage and the kill tell how it ended.',
  attack_roll_repeated: 'A trigger cause; the trigger resolves as a line.',
  attack_target: told,
  boosted: 'The boosting card resolves as a line; the card shows its stats.',
  cancelled: told,
  card_revealed: 'The table shows it; an attack on it names it.',
  character_dealt: 'Setup: the pick screen shows the characters.',
  counters: 'The card shows its counters (15.23).',
  cursed: told,
  damage_pending: 'A stack item; it resolves as a line.',
  damaged: 'The damage stack item resolves as a line.',
  deactivated: 'The causing card resolves as a line; the table shows the cards.',
  death_pending: 'A stack item; it resolves as a line.',
  destroyed: told,
  dice_rolled: told,
  died: told,
  discarded: told,
  end_of_turn: 'A trigger window, not a play.',
  entered_play: 'The gain, the buy or the play names the card.',
  expanded: told,
  extra_turn: told,
  first_player: told,
  first_roll_tied: told,
  fizzled: told,
  flipped: told,
  gained_cents: told,
  gained_item: told,
  gained_soul: told,
  gained_treasure: told,
  game_started: 'The first turn line follows.',
  game_won: told,
  gave_card: told,
  gave_cents: told,
  healed: 'Healing at the end of a turn is for everyone; HP shows on the card.',
  looked_at: told,
  loot_done: 'Each looted card has its line.',
  loot_played: 'The loot card resolves as a line, or is cancelled.',
  looted: told,
  lost_cents: told,
  moved_to_deck: 'Cards put back into a deck are mostly hidden.',
  object_died: told,
  paid: 'The buy or the ability names what was paid for.',
  passed_priority: 'Passing is no play.',
  penalty_paid: told,
  phase: 'The turn badge shows the step.',
  prevented: told,
  purchase_declared: told,
  purchase_failed: told,
  purchased: told,
  random_pick: 'Its card resolves as a line, with the results after it.',
  recharged: 'The table shows charged cards.',
  replaced: 'The results that replace the event have their lines.',
  rewards_gained: told,
  roll_changed: "The roll's line shows the change.",
  roll_repeated: 'A trigger cause; the trigger resolves as a line.',
  roll_resolved: "It changes the roll's line; no new line.",
  roll_would_resolve: 'The roll resolves as a line.',
  shielded: 'The shielding card resolves as a line; a prevent has its line.',
  stack_added: 'The stack panel shows it; it resolves as a line.',
  stack_resolved: told,
  start_of_turn: 'A trigger window, not a play.',
  stole: told,
  turn_ended: 'The next turn line follows.',
  turn_ended_early: told,
  turn_skipped: told,
  turn_started: told,
};

// entriesFor turns one update's events into history entries; key is the
// first entry's key.
// view is the table after them: a stack item's words may name its cards
// from it, as on the strip (19.52).
export function entriesFor(events: GameEvent[], turn: number, you: number, n: Names, key: number, view: GameView): Entry[] {
  const out: Entry[] = [];
  let t = turn;
  const who = (p: number) => (p === you ? 'You' : n.player(p));
  // does: "You loot", "Bo loots".
  const does = (p: number, verb: string) => (p === you ? `You ${verb}` : `${n.player(p)} ${verb}s`);
  // whom: "you", "Bo"; whose: "Your", "Bo's".
  const whom = (p: number) => (p === you ? 'you' : n.player(p));
  const whose = (p: number) => (p === you ? 'Your' : `${n.player(p)}'s`);
  // current is this update's last headline: the lines after it are
  // what it did, for the popout (15.90).
  let current: Entry | undefined;
  const add = (icon: string, text: string, card?: string, item?: StackItem, headline?: boolean) => {
    const entry: Entry = { key: key + out.length, turn: t, icon, text, ...(card ? { card } : {}), ...(item ? { item } : {}), ...(headline ? { headline } : {}) };
    out.push(entry);
    if (headline) current = entry;
    return entry;
  };
  // gave adds a result to the current headline.
  const gave = (o: Outcome) => {
    if (current) current.outcome = [...(current.outcome ?? []), o];
  };
  // character is a player's character card.
  const character = (p: number) => view.players.find((pl) => pl.id === p)?.character.card;
  // targetOf is the card a stack item aims at, a player's character
  // for a player.
  const targetOf = (it: StackItem): Outcome | undefined => {
    const ref = it.target.is_player ? character(it.target.player) : (it.target_card ?? findCard(view, it.target.object ?? 0));
    if (it.target.is_player) return { text: n.player(it.target.player), ...(ref ? { card: ref } : {}) };
    return ref ? { text: n.card(ref), card: ref } : undefined;
  };
  // A bought item also comes as a gain; the buy line names it alone.
  const bought = new Set(events.filter((e) => e.kind === 'purchased').map((e) => e.object));
  // named: objects a line names already; a second event about one (a
  // gained treasure is also a gained item; a dead item is destroyed)
  // gives no second line.
  const named = new Set<number>();
  const first = (o?: number) => {
    if (o === undefined) return true;
    if (named.has(o)) return false;
    named.add(o);
    return true;
  };
  // roll is the last resolved roll's entry and its die, for its result.
  let roll: { entry: Entry; die: number } | undefined;
  for (const e of events) {
    const card = e.card ? n.card(e.card) : '';
    switch (e.kind) {
      case 'turn_started':
        t = e.amount ?? t;
        current = undefined;
        add('▶', `Turn ${String(t)}: ${who(e.player)}`);
        break;
      case 'stack_resolved': {
        // A card's item names the card and its ability's words; any
        // other item (a roll, damage, a death) says what it did in the
        // strip's words (19.52).
        const rolled = e.amount ? ` (rolled ${String(e.amount)})` : '';
        const label = e.text && !e.text.startsWith(e.card ?? '\u0000') ? ` — ${e.text}` : '';
        const text = !card && e.item ? stackLine(e.item, view, n) : `${card || 'Roll'}${label}${rolled}`;
        const entry = add('✦', text, e.card, e.item, true);
        roll = e.amount ? { entry, die: e.amount } : undefined;
        // The popout: a card's item points at its target; damage at the
        // card it hit, with the number (15.90).
        const it = e.item;
        const aim = it ? targetOf(it) : undefined;
        if (it && card && aim) {
          entry.from = e.card;
          entry.outcome = [aim];
        } else if (it?.kind === StackKind.damage && aim) {
          // The dealer: the damage's card, else its player (stackCard
          // would fall back to the target).
          const from = it.source_card ?? findCard(view, it.source) ?? character(it.controller);
          if (from) entry.from = from;
          entry.outcome = [{ ...aim, mark: 'damage', amount: it.amount ?? 0 }];
        }
        break;
      }
      case 'roll_resolved':
        // Roll changes ("the next attack roll has -2") make the result
        // differ from the die; show both, and the result on the strip.
        if (roll && e.amount && e.amount !== roll.die) {
          const { entry, die } = roll;
          const was = entry.item && !entry.card ? ` rolls ${String(die)}` : `(rolled ${String(die)}`;
          // The changes the roll counted, by name (15.94).
          const changes = rollChanges(e.sources, (ref) => (ref ? n.card(ref) : 'a hidden card'));
          entry.text = entry.text.replace(was, `${was} → ${String(e.amount)}${changes ? ` (${changes})` : ''}`);
          if (entry.item) {
            entry.item = { ...entry.item, roll: e.amount };
          }
        }
        roll = undefined;
        break;
      case 'looted':
        // A private loot comes without its card to everyone but its player.
        add('+', card ? `${does(e.player, 'loot')} ${card}` : `${does(e.player, 'loot')} a card`, e.card);
        gave(e.card ? { text: card, card: e.card } : { text: 'a loot card' });
        break;
      case 'gave_card': {
        // A hand card changed hands (V-14): text is the giver. Only the
        // two hands get the card; the others see that one moved.
        const from = Number(e.text);
        if (e.player === you) {
          add('⇄', card ? `You get ${card} from ${who(from)}` : `You get a card from ${who(from)}`, e.card);
        } else if (from === you) {
          add('⇄', card ? `${n.player(e.player)} gets your ${card}` : `${n.player(e.player)} gets a card from you`, e.card);
        } else {
          add('⇄', `${n.player(e.player)} gets a card from ${n.player(from)}`);
        }
        break;
      }
      case 'looked_at':
        // Only the looker gets the card (V-07).
        add('◉', card ? `You look at ${card}` : `${does(e.player, 'look')} at a card`, e.card);
        break;
      case 'rewards_gained':
        add('★', `${does(e.player, 'gain')} the rewards of ${card}`, e.card);
        break;
      case 'cursed':
        add('✝', `${does(e.player, 'gain')} the curse ${card}`, e.card);
        break;
      case 'game_won':
        add('♛', `${does(e.player, 'win')} the game`);
        break;
      case 'attack_declared':
        add('⚔', `${does(e.player, 'declare')} an attack`);
        break;
      case 'attack_target':
        // Not a headline: the attack roll's tile names the monster (15.89).
        add('⚔', `${does(e.player, 'attack')} ${card || 'a card'}`, e.card);
        break;
      case 'purchase_declared':
        add('¢', `${does(e.player, 'declare')} a purchase`);
        break;
      case 'purchased':
        first(e.object);
        {
          const entry = add('¢', `${does(e.player, 'buy')} ${card}${e.text === 'treasure deck' ? ' from the treasure deck' : ''}`, e.card, undefined, true);
          entry.badge = 'buy';
          const buyer = character(e.player);
          if (buyer && e.card) {
            entry.from = buyer;
            entry.outcome = [{ text: card, card: e.card }];
          }
        }
        break;
      case 'purchase_failed':
        add('✗', `${who(e.player)} can't buy: the price is ${String(e.amount ?? 0)}¢`);
        break;
      case 'gained_treasure':
      case 'gained_item':
        if (bought.has(e.object) || !first(e.object)) break;
        if (e.text === 'starting item') {
          add('+', `${who(e.player)} ${e.player === you ? 'start' : 'starts'} with ${card}`, e.card);
        } else if (e.text === 'stolen') {
          add('⇄', `${does(e.player, 'steal')} ${card}`, e.card);
        } else {
          add('+', `${does(e.player, 'gain')} ${card}`, e.card);
        }
        if (e.card) gave({ text: card, card: e.card });
        break;
      case 'died':
      case 'object_died': {
        first(e.object);
        const by = (e.by ?? 0) - 1;
        let text: string;
        if (e.player >= 0) {
          text = by >= 0 && by !== e.player ? `${does(by, 'kill')} ${whom(e.player)}` : `${who(e.player)} ${e.player === you ? 'die' : 'dies'}`;
        } else if (e.kind === 'object_died') {
          text = by >= 0 ? `${does(by, 'destroy')} ${card}` : `${card} is destroyed`;
        } else {
          text = by >= 0 ? `${does(by, 'kill')} ${card}` : `${card} dies`;
        }
        // Every death resolves a death stack item first: its entry
        // becomes the kill, with the dead card, so the strip and the
        // History name it once (15.89).
        let death = out.findLast((d) => d.item?.kind === StackKind.death && dies(d.item, e));
        if (death) {
          death.icon = '☠';
          death.text = text;
          if (e.card) death.card = e.card;
          current = death;
        } else {
          death = add('☠', text, e.card, undefined, true);
        }
        death.badge = 'kill';
        // The popout: the killer, then the dead card with its skull; the
        // rewards follow (15.90). Without a killer the dead card is large.
        const killer = by >= 0 ? character(by) : undefined;
        if (killer && e.card) {
          death.from = killer;
          death.outcome = [{ text: e.player >= 0 ? n.player(e.player) : card, card: e.card, mark: 'death' }];
        }
        break;
      }
      case 'destroyed':
        if (e.prev !== undefined && named.has(e.prev)) break;
        if (first(e.object)) add('✗', `${card} is destroyed`, e.card);
        break;
      case 'gained_soul':
        add('♥', `${does(e.player, 'gain')} ${e.text === 'bonus soul' ? 'the bonus soul' : 'the soul of'} ${card}`, e.card);
        if (e.card) gave({ text: card, card: e.card });
        break;
      case 'gained_cents':
        add('¢', `${does(e.player, 'gain')} ${String(e.amount ?? 0)}¢`);
        gave({ text: `+${String(e.amount ?? 0)}¢`, mark: 'cents', amount: e.amount ?? 0 });
        break;
      case 'lost_cents':
        add('¢', `${who(e.player)} ${e.player === you ? 'lose' : 'loses'} ${String(e.amount ?? 0)}¢`);
        gave({ text: `−${String(e.amount ?? 0)}¢`, mark: 'cents', amount: -(e.amount ?? 0) });
        break;
      case 'stole':
        add('⇄', `${does(e.player, 'steal')} ${String(e.amount ?? 0)}¢ from ${whom(Number(e.text))}`);
        gave({ text: `+${String(e.amount ?? 0)}¢`, mark: 'cents', amount: e.amount ?? 0 });
        break;
      case 'gave_cents':
        add('⇄', `${does(e.player, 'give')} ${String(e.amount ?? 0)}¢ to ${whom((e.by ?? 0) - 1)}`);
        break;
      case 'discarded':
        // A discard from a hand or a dead player's curses; a stack card's
        // discard is its resolved line (15.74).
        if (e.player >= 0 && !e.text) add('↓', `${does(e.player, 'discard')} ${card}`, e.card);
        break;
      case 'cancelled':
        add('✗', `${card || (e.text ?? 'An ability')} is cancelled`, e.card);
        break;
      case 'fizzled':
        add('✗', `${card || 'An ability'} fizzles: its target is gone`, e.card);
        break;
      case 'prevented':
        if (!e.text) {
          add('⊘', e.player >= 0 ? `Damage to ${whom(e.player)} is prevented` : 'Damage is prevented');
        } else if (e.text === 'death' && e.player >= 0) {
          add('⊘', `${whose(e.player)} death is prevented`);
        }
        break;
      case 'penalty_paid':
        add('☠', `${who(e.player)} ${e.player === you ? 'pay' : 'pays'} the death penalty`);
        break;
      case 'turn_skipped':
        add('▶', `${whose(e.player)} turn is skipped`);
        break;
      case 'extra_turn':
        add('▶', `${who(e.player)} ${e.player === you ? 'take' : 'takes'} an extra turn`);
        break;
      case 'turn_ended_early':
        add('▶', `${whose(e.player)} turn ends early`);
        break;
      case 'expanded':
        add('+', `The table gets ${String(e.amount ?? 1)} more ${e.text ?? ''} ${e.amount === 1 ? 'slot' : 'slots'}`);
        break;
      case 'flipped':
        add('⇅', `${card || 'A card'} flips`, e.card);
        break;
      // The first-player roll keeps every die, ties too (15.43).
      case 'dice_rolled':
        if (e.text === firstRollLabel) {
          add('⚄', `${does(e.player, 'roll')} ${String(e.amount ?? 0)} for the first player`);
        }
        break;
      case 'first_roll_tied':
        add('⚄', tieText(e, you, n.player));
        break;
      case 'first_player':
        add('①', firstPlayerText(e, you, n.player, n.card), e.card);
        break;
    }
  }
  return out;
}

// dies says whether a death stack item is the one a death event ends:
// the same player, or the object that moved away as it died.
function dies(it: StackItem, e: GameEvent): boolean {
  return it.target.is_player ? e.player >= 0 && it.target.player === e.player : e.player < 0 && it.target.object === e.prev;
}

// stripEntries are the history strip's tiles, newest on top (15.89):
// the headlines. Each keeps its entry's key, so an older tile slides
// down when a new one comes, instead of being drawn again (15.51).
export function stripEntries(history: Entry[]): Entry[] {
  return history.filter((e) => e.headline).reverse();
}

// rolledBack is the history after a rollback (ADR 020, item 8): every
// entry from the target step on is undone, and a row tells who rolled
// back. step is the rollback's own number.
export function rolledBack(history: Entry[], r: Rollback, step: number, turn: number, you: number, n: Names): Entry[] {
  const out = history.map((e) => (e.step !== undefined && e.step >= r.to && !e.undone ? { ...e, undone: true } : e));
  const who = r.seat === you ? 'You roll' : `${n.player(r.seat)} rolls`;
  const text = `${who} back ${String(r.undid)} ${r.undid === 1 ? 'step' : 'steps'}`;
  out.push({ key: (history.at(-1)?.key ?? 0) + 1, turn, icon: '↶', text, step, rollback: true });
  return out;
}

// rollbackTargets are the keys of the entries that get a rollback
// button: the first live entry of each step since from (ADR 020, item
// 6). An undone entry or a rollback row is no target.
export function rollbackTargets(history: Entry[], from: number): Set<number> {
  const out = new Set<number>();
  let last: number | undefined;
  for (const e of history) {
    if (e.step !== undefined && e.step !== last && e.step >= from && !e.undone && !e.rollback) {
      out.add(e.key);
    }
    last = e.step;
  }
  return out;
}

// firstRollLabel is the label of the first-player roll's dice.
export const firstRollLabel = 'first player roll';

// Looked is one card this viewer looks at; id is its object.
export interface Looked {
  id: number;
  card: string;
}

// lookedAt is the cards this viewer looks at in one update's events,
// for the look dialog (V-07). A card looked at twice in one update
// shows once (19.13).
export function lookedAt(events: GameEvent[], you: number): Looked[] {
  const seen = new Set<number>();
  return events.flatMap((e) => {
    if (e.kind !== 'looked_at' || e.player !== you || !e.card) return [];
    const id = e.object ?? 0;
    if (seen.has(id)) return [];
    seen.add(id);
    return [{ id, card: e.card }];
  });
}

// Turn step names, in Go's order (pkg/rules_engine/turn.go).
const steps = ['Recharge', 'Start of turn', 'Loot', 'Loot', 'Action', 'End of turn', 'Hand size', 'Cleanup'];

export function stepName(step: number): string {
  return steps.at(step) ?? '';
}
