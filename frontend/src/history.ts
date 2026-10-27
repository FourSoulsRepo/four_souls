// The match history (7.9, V-03, V-51): a line for every play, under turn
// lines: who did what. Built from the events this viewer got, so it shows
// only what the viewer may see (hidden cards stay hidden).
import { firstPlayerText, tieText } from './dice';
import { say, t } from './i18n';
import type { GameEvent, GameView, Rollback, RollbackEnd, RollbackStep, StackItem, Update } from './protocol';
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
  // restarted marks the row where a server restart cut the history.
  restarted?: boolean;
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
  let turnNo = turn;
  // mine: the event is the viewer's own; its line says "You".
  const mine = (p: number) => p === you;
  // current is this update's last headline: the lines after it are
  // what it did, for the popout (15.90).
  let current: Entry | undefined;
  const add = (icon: string, text: string, card?: string, item?: StackItem, headline?: boolean) => {
    const entry: Entry = { key: key + out.length, turn: turnNo, icon, text, ...(card ? { card } : {}), ...(item ? { item } : {}), ...(headline ? { headline } : {}) };
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
  // roll is the last resolved roll's entry, its die, and its line with
  // another number in the die's place, for its result.
  let roll: { entry: Entry; die: number; line: (rolled?: string) => string } | undefined;
  for (const e of events) {
    const card = e.card ? n.card(e.card) : '';
    const player = n.player(e.player);
    switch (e.kind) {
      case 'turn_started':
        turnNo = e.amount ?? turnNo;
        current = undefined;
        add('▶', t('history.turn', { n: turnNo, player: mine(e.player) ? t('history.you') : player }));
        break;
      case 'stack_resolved': {
        // A card's item names the card and its ability's words; any
        // other item (a roll, damage, a death) says what it did in the
        // strip's words (19.52). The words come from the server.
        const words = say(e.label, n);
        const label = words ? ` — ${words}` : '';
        const it = e.item;
        // line says it, with rolled in the die's place, if given.
        const line = (rolled?: string) => {
          if (!card && it) {
            return stackLine(it, view, n, rolled);
          }
          const shown = rolled ?? (e.amount ? String(e.amount) : '');
          return `${card || t('history.roll')}${label}${shown ? ` ${t('history.rolled', { roll: shown })}` : ''}`;
        };
        const entry = add('✦', line(), e.card, e.item, true);
        roll = e.amount ? { entry, die: e.amount, line } : undefined;
        // The popout: a card's item points at its target; damage at the
        // card it hit, with the number (15.90).
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
          // The changes the roll counted, by name (15.94).
          const changes = rollChanges(e.sources, (ref) => (ref ? n.card(ref) : t('stats.hidden_card')));
          entry.text = roll.line(`${String(die)} → ${String(e.amount)}${changes ? ` (${changes})` : ''}`);
          if (entry.item) {
            entry.item = { ...entry.item, roll: e.amount };
          }
        }
        roll = undefined;
        break;
      case 'looted':
        // A private loot comes without its card to everyone but its player.
        if (card) {
          add('+', mine(e.player) ? t('history.loot.you', { card }) : t('history.loot.other', { player, card }), e.card);
        } else {
          add('+', mine(e.player) ? t('history.loot_hidden.you') : t('history.loot_hidden.other', { player }), e.card);
        }
        gave(e.card ? { text: card, card: e.card } : { text: t('history.a_loot_card') });
        break;
      case 'gave_card': {
        // A hand card changed hands (V-14): text is the giver. Only the
        // two hands get the card; the others see that one moved.
        const from = n.player(Number(e.text));
        if (mine(e.player)) {
          add('⇄', card ? t('history.got_card.you', { card, from }) : t('history.got_hidden.you', { from }), e.card);
        } else if (mine(Number(e.text))) {
          add('⇄', card ? t('history.got_card.yours', { player, card }) : t('history.got_hidden.yours', { player }), e.card);
        } else {
          add('⇄', t('history.got_hidden.other', { player, from }));
        }
        break;
      }
      case 'looked_at':
        // Only the looker gets the card (V-07).
        if (card) {
          add('◉', t('history.look.you', { card }), e.card);
        } else {
          add('◉', mine(e.player) ? t('history.look_hidden.you') : t('history.look_hidden.other', { player }), e.card);
        }
        break;
      case 'rewards_gained':
        add('★', mine(e.player) ? t('history.rewards.you', { card }) : t('history.rewards.other', { player, card }), e.card);
        break;
      case 'cursed':
        add('✝', mine(e.player) ? t('history.cursed.you', { card }) : t('history.cursed.other', { player, card }), e.card);
        break;
      case 'game_won':
        add('♛', mine(e.player) ? t('history.won.you') : t('history.won.other', { player }));
        break;
      case 'attack_declared':
        add('⚔', mine(e.player) ? t('history.attack_declared.you') : t('history.attack_declared.other', { player }));
        break;
      case 'attack_target': {
        // Not a headline: the attack roll's tile names the monster (15.89).
        const target = card || t('stack.a_card');
        add('⚔', mine(e.player) ? t('history.attack.you', { card: target }) : t('history.attack.other', { player, card: target }), e.card);
        break;
      }
      case 'purchase_declared':
        add('¢', mine(e.player) ? t('history.purchase_declared.you') : t('history.purchase_declared.other', { player }));
        break;
      case 'purchased':
        first(e.object);
        {
          const deck = e.text === 'treasure deck';
          const text = mine(e.player)
            ? deck
              ? t('history.buy_deck.you', { card })
              : t('history.buy.you', { card })
            : deck
              ? t('history.buy_deck.other', { player, card })
              : t('history.buy.other', { player, card });
          const entry = add('¢', text, e.card, undefined, true);
          entry.badge = 'buy';
          const buyer = character(e.player);
          if (buyer && e.card) {
            entry.from = buyer;
            entry.outcome = [{ text: card, card: e.card }];
          }
        }
        break;
      case 'purchase_failed': {
        const price = e.amount ?? 0;
        add('✗', mine(e.player) ? t('history.cant_buy.you', { n: price }) : t('history.cant_buy.other', { player, n: price }));
        break;
      }
      case 'gained_treasure':
      case 'gained_item':
        if (bought.has(e.object) || !first(e.object)) break;
        if (e.text === 'starting item') {
          add('+', mine(e.player) ? t('history.starts_with.you', { card }) : t('history.starts_with.other', { player, card }), e.card);
        } else if (e.text === 'stolen') {
          add('⇄', mine(e.player) ? t('history.steal_item.you', { card }) : t('history.steal_item.other', { player, card }), e.card);
        } else {
          add('+', mine(e.player) ? t('history.gain.you', { card }) : t('history.gain.other', { player, card }), e.card);
        }
        if (e.card) gave({ text: card, card: e.card });
        break;
      case 'died':
      case 'object_died': {
        first(e.object);
        const by = (e.by ?? 0) - 1;
        const killer = n.player(by);
        let text: string;
        if (e.player >= 0) {
          if (by >= 0 && by !== e.player) {
            text = mine(by) ? t('history.kill.you', { target: player }) : mine(e.player) ? t('history.kill.yours', { player: killer }) : t('history.kill.other', { player: killer, target: player });
          } else {
            text = mine(e.player) ? t('history.die.you') : t('history.die.other', { player });
          }
        } else if (e.kind === 'object_died') {
          text = by >= 0 ? (mine(by) ? t('history.destroy.you', { card }) : t('history.destroy.other', { player: killer, card })) : t('history.destroyed', { card });
        } else {
          text = by >= 0 ? (mine(by) ? t('history.kill.you', { target: card }) : t('history.kill.other', { player: killer, target: card })) : t('history.dies', { card });
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
        const killerCard = by >= 0 ? character(by) : undefined;
        if (killerCard && e.card) {
          death.from = killerCard;
          death.outcome = [{ text: e.player >= 0 ? player : card, card: e.card, mark: 'death' }];
        }
        break;
      }
      case 'destroyed':
        if (e.prev !== undefined && named.has(e.prev)) break;
        if (first(e.object)) add('✗', t('history.destroyed', { card }), e.card);
        break;
      case 'gained_soul': {
        const bonus = e.text === 'bonus soul';
        const text = mine(e.player)
          ? bonus
            ? t('history.bonus_soul.you', { card })
            : t('history.soul.you', { card })
          : bonus
            ? t('history.bonus_soul.other', { player, card })
            : t('history.soul.other', { player, card });
        add('♥', text, e.card);
        if (e.card) gave({ text: card, card: e.card });
        break;
      }
      case 'gained_cents': {
        const cents = e.amount ?? 0;
        add('¢', mine(e.player) ? t('history.cents.you', { n: cents }) : t('history.cents.other', { player, n: cents }));
        gave({ text: `+${String(cents)}¢`, mark: 'cents', amount: cents });
        break;
      }
      case 'lost_cents': {
        const cents = e.amount ?? 0;
        add('¢', mine(e.player) ? t('history.lose_cents.you', { n: cents }) : t('history.lose_cents.other', { player, n: cents }));
        gave({ text: `−${String(cents)}¢`, mark: 'cents', amount: -cents });
        break;
      }
      case 'stole': {
        const cents = e.amount ?? 0;
        const victim = Number(e.text);
        const text = mine(e.player)
          ? t('history.steal.you', { n: cents, target: n.player(victim) })
          : mine(victim)
            ? t('history.steal.yours', { player, n: cents })
            : t('history.steal.other', { player, n: cents, target: n.player(victim) });
        add('⇄', text);
        gave({ text: `+${String(cents)}¢`, mark: 'cents', amount: cents });
        break;
      }
      case 'gave_cents': {
        const cents = e.amount ?? 0;
        const to = (e.by ?? 0) - 1;
        const text = mine(e.player)
          ? t('history.give.you', { n: cents, target: n.player(to) })
          : mine(to)
            ? t('history.give.yours', { player, n: cents })
            : t('history.give.other', { player, n: cents, target: n.player(to) });
        add('⇄', text);
        break;
      }
      case 'discarded':
        // A discard from a hand or a dead player's curses; a stack card's
        // discard is its resolved line (15.74).
        if (e.player >= 0 && !e.text) add('↓', mine(e.player) ? t('history.discard.you', { card }) : t('history.discard.other', { player, card }), e.card);
        break;
      case 'cancelled':
        // Without a card, the text names it; it comes from the server.
        add('✗', t('history.cancelled', { card: card || say(e.label, n) || t('history.an_ability') }), e.card);
        break;
      case 'fizzled':
        add('✗', t('history.fizzled', { card: card || t('history.an_ability') }), e.card);
        break;
      case 'prevented':
        if (!e.text) {
          add('⊘', e.player < 0 ? t('history.prevented') : mine(e.player) ? t('history.prevented_to.you') : t('history.prevented_to.other', { player }));
        } else if (e.text === 'death' && e.player >= 0) {
          add('⊘', mine(e.player) ? t('history.death_prevented.you') : t('history.death_prevented.other', { player }));
        }
        break;
      case 'penalty_paid':
        add('☠', mine(e.player) ? t('history.penalty.you') : t('history.penalty.other', { player }));
        break;
      case 'turn_skipped':
        add('▶', mine(e.player) ? t('history.skipped.you') : t('history.skipped.other', { player }));
        break;
      case 'extra_turn':
        add('▶', mine(e.player) ? t('history.extra_turn.you') : t('history.extra_turn.other', { player }));
        break;
      case 'turn_ended_early':
        add('▶', mine(e.player) ? t('history.ended_early.you') : t('history.ended_early.other', { player }));
        break;
      case 'expanded': {
        // text is the kind of slot; the engine sends "shop" or "monster".
        const slots = e.amount ?? 1;
        add(
          '+',
          e.text === 'shop'
            ? t('history.expanded.shop', { n: slots })
            : e.text === 'monster'
              ? t('history.expanded.monster', { n: slots })
              : t('history.expanded.other', { n: slots, kind: e.text ?? '' }),
        );
        break;
      }
      case 'flipped':
        add('⇅', t('history.flipped', { card: card || t('history.a_card') }), e.card);
        break;
      // The first-player roll keeps every die, ties too (15.43).
      case 'dice_rolled':
        if (e.text === firstRollCode) {
          const rolled = e.amount ?? 0;
          add('⚄', mine(e.player) ? t('history.first_roll.you', { n: rolled }) : t('history.first_roll.other', { player, n: rolled }));
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
  const text = r.seat === you ? t('history.rolled_back.you', { n: r.undid }) : t('history.rolled_back.other', { player: n.player(r.seat), n: r.undid });
  out.push({ key: (history.at(-1)?.key ?? 0) + 1, turn, icon: '↶', text, step, rollback: true });
  return out;
}

// Move is what one update did, for the Undo button (17.7): from and to
// are its first and last step. text names it; barrier, if set, is the
// first barrier it holds, in words for "Nothing to undo since …".
export interface Move {
  from: number;
  to: number;
  text?: string;
  barrier?: string;
}

// barrierWords name the Soft barriers a player can tell apart (ADR 020,
// item 4). The server decides what is one; this only names it.
function barrierWords(kind: string): string | undefined {
  switch (kind) {
    case 'dice_rolled':
    case 'roll_changed':
      return t('history.barrier.die');
    case 'random_pick':
      return t('history.barrier.random');
    case 'damaged':
      return t('history.barrier.damage');
    case 'died':
    case 'object_died':
      return t('history.barrier.died');
    case 'turn_ended':
    case 'turn_ended_early':
      return t('history.barrier.turn_ended');
    case 'turn_skipped':
      return t('history.barrier.turn_skipped');
    case 'looted':
    case 'gained_treasure':
      return t('history.barrier.drawn');
    case 'card_revealed':
      return t('history.barrier.revealed');
    case 'looked_at':
      return t('history.barrier.seen');
    default:
      return undefined;
  }
}

// barrierOf names a barrier event; any other private one shows a hidden
// card, but for a hand card put on the loot deck (item 4.22).
function barrierOf(e: GameEvent): string | undefined {
  return barrierWords(e.kind) ?? (e.private && e.kind !== 'moved_to_deck' ? t('history.barrier.seen') : undefined);
}

// moveOf names what one update did. A played card or a used ability is
// named by its play, which has no history line (15.88): "Bo plays A
// Penny!". An ended turn says so; anything else is its first line.
export function moveOf(events: GameEvent[], entries: Entry[], from: number, to: number, you: number, n: Names): Move {
  const play = events.find((e) => e.kind === 'loot_played' || e.kind === 'activated');
  const ended = events.find((e) => e.kind === 'turn_ended' || e.kind === 'turn_ended_early');
  const text = play
    ? playText(play, you, n)
    : ended
      ? ended.player === you
        ? t('history.move.end_turn.you')
        : t('history.move.end_turn.other', { player: n.player(ended.player) })
      : entries.find((e) => !e.rollback && e.icon !== '▶')?.text;
  const barrier = events.map(barrierOf).find((w) => w !== undefined);
  return { from, to, ...(text ? { text } : {}), ...(barrier ? { barrier } : {}) };
}

// playText names a played card or a used ability: "Bo plays A Penny!".
function playText(e: GameEvent, you: number, n: Names): string {
  const card = e.card ? n.card(e.card) : t('history.an_ability_lower');
  const player = n.player(e.player);
  if (e.kind === 'loot_played') {
    return e.player === you ? t('history.move.play.you', { card }) : t('history.move.play.other', { player, card });
  }
  return e.player === you ? t('history.move.use.you', { card }) : t('history.move.use.other', { player, card });
}

// moveText names the move a target starts: its update's, or an
// answer's up to its last step (a play that asks a target is named by
// the answer's update). A move with no line of its own never borrows
// the next one's: a trigger order or a card's question is "Bo's
// choice" (17.8).
export function moveText(moves: Move[], target: RollbackStep, you: number, name: (p: number) => string): string {
  const last = target.last ?? target.step;
  const text = moves.find((m) => m.to >= target.step && m.from <= last && m.text)?.text;
  return text ?? (target.seat === you ? t('history.choice.you') : t('history.choice.other', { player: name(target.seat) }));
}

// Undo is the History's one Undo button (17.7, ADR 020). step is the
// latest target the server lists; ask: the Soft author asks the others.
// why, if set, is why it is greyed; title names the move, or says why.
export interface Undo {
  label: string;
  step?: number;
  ask: boolean;
  why?: string;
  title: string;
}

// undoButton is the Undo button for this viewer, or null in Disabled
// mode. One press undoes the latest move; Godmode at once, Soft by a
// vote of the others, or the author's ask (items 5 and 6).
export function undoButton(
  u: Pick<Update, 'rollback_mode' | 'rollback_from' | 'rollback_steps' | 'rollback_vote' | 'rollback_cut' | 'pause'> & { view: { over?: boolean } },
  you: number,
  moves: Move[],
  name: (p: number) => string,
): Undo | null {
  if (!u.rollback_mode) {
    return null;
  }
  const target = u.rollback_steps?.at(-1);
  const ask = u.rollback_mode === 'soft' && target?.seat === you;
  const label = ask ? t('history.undo.ask') : t('history.undo.undo');
  const why =
    you < 0
      ? t('history.undo.spectator')
      : u.view.over
        ? t('history.undo.over')
        : u.rollback_vote
          ? t('history.undo.voting')
          : u.pause
            ? t('history.undo.paused')
            : !target
              ? nothingWhy(u, moves)
              : undefined;
  if (why !== undefined || !target) {
    const title = why ?? t('history.undo.nothing');
    return { label, ask, why: title, title };
  }
  const move = moveText(moves, target, you, name);
  return { label, step: target.step, ask, title: ask ? t('history.undo.ask_title', { move }) : t('history.undo.title', { move }) };
}

// nothingWhy says why no move is left to undo: the barrier just before
// rollback_from, if the history names it.
function nothingWhy(u: Pick<Update, 'rollback_mode' | 'rollback_from' | 'rollback_cut'>, moves: Move[]): string {
  const at = (u.rollback_from ?? 1) - 1;
  if (u.rollback_cut && at <= u.rollback_cut) {
    return t('history.undo.restarted');
  }
  const barrier = u.rollback_mode === 'soft' ? moves.find((m) => m.from <= at && at <= m.to)?.barrier : undefined;
  return barrier ? t('history.undo.since', { barrier }) : t('history.undo.yet');
}

// voteEnded adds the row of a Soft vote that ended with no rollback.
export function voteEnded(history: Entry[], end: RollbackEnd, turn: number, you: number, n: Names): Entry[] {
  const text = {
    declined: () => t('history.vote.declined'),
    cancelled: () => (end.seat === you ? t('history.vote.cancelled.you') : t('history.vote.cancelled.other', { player: n.player(end.seat) })),
    timeout: () => t('history.vote.timeout'),
    left: () => t('history.vote.left'),
    failed: () => t('history.vote.failed'),
  }[end.why]();
  return [...history, { key: (history.at(-1)?.key ?? 0) + 1, turn, icon: '↶', text, rollback: true }];
}

// restartCut adds the row where a server restart cut the rollback
// history (ADR 020, item 10.6): nothing before it can be rolled back.
// cut is the step the restart came at; a row for it shows once.
export function restartCut(history: Entry[], cut: number, turn: number): Entry[] {
  if (history.some((e) => e.restarted && e.step === cut)) {
    return history;
  }
  const text = t('history.restarted');
  return [...history, { key: (history.at(-1)?.key ?? 0) + 1, turn, icon: '⟳', text, step: cut, rollback: true, restarted: true }];
}

// firstRollCode marks the first-player roll's dice (dice_rolled text).
export const firstRollCode = 'first';

// Looked is one card this viewer looks at; id is its object. owner is
// the seat whose hand it is in; none for a deck card (19.71).
export interface Looked {
  id: number;
  card: string;
  owner?: number;
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
    return [{ id, card: e.card, ...(e.by ? { owner: e.by - 1 } : {}) }];
  });
}

// lookGroups splits looked-at cards by the hand they are in, in the
// order they came; deck cards make one group with no owner (19.71).
export function lookGroups(cards: Looked[]): { owner?: number; cards: Looked[] }[] {
  const groups: { owner?: number; cards: Looked[] }[] = [];
  for (const c of cards) {
    const g = groups.find((x) => x.owner === c.owner);
    if (g) {
      g.cards.push(c);
    } else {
      groups.push({ ...(c.owner === undefined ? {} : { owner: c.owner }), cards: [c] });
    }
  }
  return groups;
}

// stepName names a turn step, in Go's order (pkg/rules_engine/turn.go).
export function stepName(step: number): string {
  switch (step) {
    case 0:
      return t('turn_badge.step.recharge');
    case 1:
      return t('turn_badge.step.start');
    case 2:
    case 3:
      return t('turn_badge.step.loot');
    case 4:
      return t('turn_badge.step.action');
    case 5:
      return t('turn_badge.step.end');
    case 6:
      return t('turn_badge.step.hand_size');
    case 7:
      return t('turn_badge.step.cleanup');
    default:
      return '';
  }
}
