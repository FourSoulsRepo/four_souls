// The match history (7.9, V-03): one entry per resolved stack item, plus
// turns, loots and the winner. Built from the events this viewer got, so
// it shows only what the viewer may see (hidden cards stay hidden).
import { firstPlayerText, tieText } from './dice';
import type { GameEvent, StackItem } from './protocol';
import type { Names } from './stack';

export interface Entry {
  key: number;
  turn: number;
  icon: string;
  text: string;
  // card is the card the entry is about, for its zoom.
  card?: string;
  // item is the stack item that resolved, for the history strip (15.22).
  item?: StackItem;
}

// entriesFor turns one update's events into history entries; key is the
// first entry's key.
export function entriesFor(events: GameEvent[], turn: number, you: number, n: Names, key: number): Entry[] {
  const out: Entry[] = [];
  let t = turn;
  const who = (p: number) => (p === you ? 'You' : n.player(p));
  // does: "You loot", "Bo loots".
  const does = (p: number, verb: string) => (p === you ? `You ${verb}` : `${n.player(p)} ${verb}s`);
  const add = (icon: string, text: string, card?: string, item?: StackItem) => {
    out.push({ key: key + out.length, turn: t, icon, text, ...(card ? { card } : {}), ...(item ? { item } : {}) });
  };
  // roll is the last resolved roll's entry and its die, for its result.
  let roll: { entry: Entry; die: number } | undefined;
  for (const e of events) {
    const card = e.card ? n.card(e.card) : '';
    switch (e.kind) {
      case 'turn_started':
        t = e.amount ?? t;
        add('▶', `Turn ${String(t)}: ${who(e.player)}`);
        break;
      case 'stack_resolved': {
        const rolled = e.amount ? ` (rolled ${String(e.amount)})` : '';
        const label = e.text && !e.text.startsWith(e.card ?? '\u0000') ? ` — ${e.text}` : '';
        add('✦', `${card || 'Roll'}${label}${rolled}`, e.card, e.item);
        roll = e.amount ? { entry: out[out.length - 1], die: e.amount } : undefined;
        break;
      }
      case 'roll_resolved':
        // Roll changes ("the next attack roll has -2") make the result
        // differ from the die; show both, and the result on the strip.
        if (roll && e.amount && e.amount !== roll.die) {
          const { entry, die } = roll;
          entry.text = entry.text.replace(`(rolled ${String(die)})`, `(rolled ${String(die)} → ${String(e.amount)})`);
          if (entry.item) {
            entry.item = { ...entry.item, roll: e.amount };
          }
        }
        roll = undefined;
        break;
      case 'looted':
        // A private loot comes without its card to everyone but its player.
        add('+', card ? `${does(e.player, 'loot')} ${card}` : `${does(e.player, 'loot')} a card`, e.card);
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
