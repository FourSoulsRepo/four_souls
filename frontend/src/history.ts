// The match history (7.9, V-03): one entry per resolved stack item, plus
// turns, loots and the winner. Built from the events this viewer got, so
// it shows only what the viewer may see (hidden cards stay hidden).
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
  for (const e of events) {
    const card = e.card ? n.card(e.card) : '';
    switch (e.kind) {
      case 'turn_started':
        t = e.amount ?? t;
        add('▶', `Turn ${String(t)}: ${who(e.player)}`);
        break;
      case 'stack_resolved': {
        const roll = e.amount ? ` (rolled ${String(e.amount)})` : '';
        const label = e.text && !e.text.startsWith(e.card ?? '\u0000') ? ` — ${e.text}` : '';
        add('✦', `${card || 'Roll'}${label}${roll}`, e.card, e.item);
        break;
      }
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
    }
  }
  return out;
}

// Looked is one card this viewer looks at; id is its object.
export interface Looked {
  id: number;
  card: string;
}

// lookedAt is the cards this viewer looks at in one update's events,
// for the look dialog (V-07).
export function lookedAt(events: GameEvent[], you: number): Looked[] {
  return events.flatMap((e) => (e.kind === 'looked_at' && e.player === you && e.card ? [{ id: e.object ?? 0, card: e.card }] : []));
}

// Turn step names, in Go's order (pkg/rules_engine/turn.go).
const steps = ['Recharge', 'Start of turn', 'Loot', 'Loot', 'Action', 'End of turn', 'Hand size', 'Cleanup'];

export function stepName(step: number): string {
  return steps.at(step) ?? '';
}
