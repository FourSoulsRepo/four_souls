import { createContext, use, useLayoutEffect, useMemo, useRef, useState, type CSSProperties, type ReactNode, type RefObject } from 'react';
import { discardIntent, modeLines, orderIntent, outcomeLines, type Actions } from '../actions';
import { claim } from '../animations';
import { deckKinds, type Kind } from '../cards';
import { choiceRow, confirmFirst, dealDelays, dealPose, keptZone, pickToggle } from '../choice';
import { feel } from '../feel';
import { lookGroups, type Looked } from '../history';
import { t } from '../i18n';
import type { CardView, ModeView } from '../protocol';
import { playSound } from '../sounds';
import { seatOf, useApp, viewScale } from '../store';
import { Card } from './Card';
import { CardDialog } from './CardDialog';
import { deckBox, flyBack, keepCard } from './ChoiceBack';
import { still } from './lunge';
import './Actions.css';

// Entry is how a choice's cards come in (20.18, ADR 027 "Choosing
// cards"): dealt from the deck they are looked at in, grown at their
// places, or at once (speed "off", reduced motion).
type Entry = 'deal' | 'grow' | 'none';

interface Frame {
  // deck is the deck the cards are looked at in, if it is on screen.
  deck?: Kind | undefined;
  // fly sends the cards back to the deck after the pick, in this order
  // of their elements: the first flies first (20.18). Those from place
  // under on go under the deck (19.147).
  fly: (order: Element[], under?: number) => void;
  // keep holds the kept card at its place, for its flight to the
  // player, and sends the rest back to the deck (20.22).
  keep: (kept: Element, rest: Element[]) => void;
}

const FrameContext = createContext<Frame>({ fly: () => undefined, keep: () => undefined });

// useLook is the deck the viewer looks at in the prompt waiting now
// (20.18): the server says which deck, and how many cards.
function useLook(): Kind | undefined {
  return useApp((s) => {
    const look = s.game?.view.waiting.look;
    return look ? deckKinds[look.deck ?? 0] : undefined;
  });
}

// dealt keeps, per prompt step, the cards a deal claimed, so a second
// mount of the same choice (React's StrictMode) does not claim again.
const dealt = new Map<number, Set<string>>();

// claimDeal claims each card a deal shows (ADR 012): it returns the
// cards that may fly; a card another animation showed just now comes at
// once instead. A card's object, when known, tells copies apart (19.141).
function claimDeal(step: number, cards: { card: string; object?: number | undefined }[]): Set<string> {
  let ok = dealt.get(step);
  if (!ok) {
    const seen = new Map(cards.map((c) => [c.card, c.object]));
    ok = new Set([...seen].filter(([card, object]) => claim(card, 'show', 'choice-deal', { object })).map(([card]) => card));
    dealt.clear();
    dealt.set(step, ok);
  }
  return ok;
}

const ms = (s: number) => s * 1000;

// sample is a motion as WAAPI keyframes, one per frame of 1/30 s.
function sample<T>(time: number, at: (t: number) => T): T[] {
  const n = Math.max(2, Math.ceil(time * 30));
  return Array.from({ length: n + 1 }, (_, i) => at((time * i) / n));
}

// enter plays the cards' entry: each element marked data-deal, in
// order. It returns a stop for the effect's cleanup.
function enter(root: HTMLElement, entry: Entry, deck: Kind | undefined, scale: number, step: number, back: boolean): () => void {
  const els = [...root.querySelectorAll<HTMLElement>('[data-deal]')];
  const c = feel.choice;
  const parts: (Animation | undefined)[] = [];
  const timers: number[] = [];
  const grow = (el: HTMLElement) =>
    // anim: choice-grow
    el.animate([{ transform: 'scale(0)' }, { transform: 'scale(1)' }], { duration: ms(c.grow * scale), easing: c.growCurve });
  if (back) {
    // anim: choice-grow
    parts.push(root.closest('.choice-backdrop')?.animate([{ opacity: 0 }, { opacity: 1 }], { duration: ms(c.grow * scale), easing: c.growCurve }));
  }
  const box = entry === 'deal' && deck ? deckBox(deck) : null;
  if (entry === 'grow' || (entry === 'deal' && !box)) {
    parts.push(...els.map(grow));
  } else if (entry === 'deal' && box) {
    const may = claimDeal(
      step,
      els.map((el) => ({ card: el.dataset.deal ?? '', ...(el.dataset.object ? { object: Number(el.dataset.object) } : {}) })),
    );
    const delays = dealDelays(els.length);
    root.classList.add('dealing');
    els.forEach((el, i) => {
      if (!may.has(el.dataset.deal ?? '')) {
        return;
      }
      const r = el.getBoundingClientRect();
      const from = { x: box.left + box.width / 2 - (r.left + r.width / 2), y: box.top + box.height / 2 - (r.top + r.height / 2) };
      const poses = sample(c.deal, (tt) => dealPose(tt, from, box.width / Math.max(r.width, 1)));
      const timing = { duration: ms(c.deal * scale), delay: ms((delays.at(i) ?? 0) * scale), fill: 'backwards' } as const;
      // anim: choice-deal
      parts.push(el.animate(poses.map((p) => ({ transform: `translate(${String(p.x)}px, ${String(p.y)}px) scale(${String(p.scale)})` })), timing));
      // anim: choice-deal
      parts.push(el.querySelector('.deal-turn')?.animate(poses.map((p) => ({ transform: `rotateY(${String(p.turn)}deg)` })), timing));
      timers.push(
        window.setTimeout(() => {
          playSound('choiceDeal');
        }, ms(((delays.at(i) ?? 0) + c.sound) * scale)),
      );
    });
    const last = (delays.at(-1) ?? 0) + c.deal;
    timers.push(
      window.setTimeout(() => {
        root.classList.remove('dealing');
      }, ms(last * scale)),
    );
  }
  return () => {
    root.classList.remove('dealing');
    for (const id of timers) {
      window.clearTimeout(id);
    }
    for (const p of parts) {
      p?.cancel();
    }
  };
}

// DealCard is a card that may be dealt face down: its back lies behind
// its face, and the deal turns it about its vertical axis (20.18).
function DealCard({ kind, children }: { kind: Kind | undefined; children: ReactNode }) {
  return (
    <div className="deal-turn">
      <div className="deal-face">{children}</div>
      {kind && (
        <div className="deal-face deal-back" aria-hidden="true">
          <Card back={kind} size="choice" />
        </div>
      )}
    </div>
  );
}

// rowStyle sizes a row of n cards (20.18, choiceRow).
function rowStyle(n: number): { className: string; style: CSSProperties } {
  const r = choiceRow(n);
  return {
    className: `choice-row${r.wrap ? ' wrap' : ''}`,
    style: { '--choice-scale': r.scale, '--choice-gap': r.gap } as CSSProperties,
  };
}

// ChoiceDialog asks the player to pick one option. Card options show as
// cards, Hearthstone-style (V-07); options on the table can be clicked
// there too. subject is the card the choice is about, e.g. the monster
// a deck attack drew (19.9) or the deck top a player looks at (19.11);
// it shows beside the options; option labels name it as a card (ADR 021).
// An option with an outcome says what picking it gives (19.17); a
// trigger option says its ability's words (19.19).
//
// A choice shows big cards over the dimmed board (15.29, H-11). A card
// option is the card itself: a click picks it; it has no name button
// (20.18). Cards looked at in a deck are dealt from it and fly back to
// it after the pick; other cards grow in. "View board" hides it to look
// at the table, where its targets still glow and answer, and "Back to
// the choice" brings it back (F-20).
//
// A start-of-game choice (Eden's) is a screen of its own (20.22, ADR
// 027 "Start of game"): a banner over the cards, a click picks a card
// and Confirm sends it. The kept card then flies to the player.
export function ChoiceDialog({ choice, subject }: { choice: NonNullable<Actions['choice']>; subject?: string | undefined }) {
  const act = useApp((s) => s.act);
  const cardName = (ref: string) => useApp.getState().cards.get(ref)?.name ?? ref;
  const look = useLook();
  const confirm = confirmFirst(choice.purpose);
  const [picked, setPicked] = useState<number | null>(null);
  const character = useApp((s) => s.game?.view.players.find((p) => p.id === seatOf(s))?.character.card);
  const cards = choice.options.filter((o) => o.card).length;
  const row = rowStyle(cards);
  const banner = confirm ? <StartBanner caption={t('start.caption', { card: character ? cardName(character) : '' })} /> : undefined;
  const send = picked === null ? undefined : choice.options.find((o) => (o.intent.choice ?? 0) === picked)?.intent;
  return (
    <ChoiceFrame
      label={choice.title}
      banner={banner}
      className={confirm ? 'start' : ''}
      style={confirm ? ({ '--pick-lift': `${String(feel.start.lift)}px` } as CSSProperties) : undefined}
      look={look}
      cards={cards > 0}
    >
      {subject && (
        <div className="choice-subject">
          <Card cardRef={subject} size="choice" />
        </div>
      )}
      <div className={`choice-options ${row.className}`} style={row.style}>
        {choice.options.map((o) => (
          <div key={o.intent.choice ?? 0} className="choice-option">
            {o.card ? (
              <OptionCard
                card={o.card}
                object={o.target?.object}
                name={cardName(o.card)}
                deck={look}
                pick={() => {
                  act(o.intent);
                }}
                {...(confirm
                  ? {
                      picked: picked === (o.intent.choice ?? 0),
                      toggle: () => {
                        const next = pickToggle(picked, o.intent.choice ?? 0);
                        playSound(next.sound);
                        setPicked(next.picked);
                      },
                    }
                  : {})}
              />
            ) : (
              <button
                type="button"
                onClick={() => {
                  act(o.intent);
                }}
              >
                {o.label}
              </button>
            )}
            {o.line && <p className="choice-line">{o.line}</p>}
            {o.outcome && (
              <ul className="choice-outcome" aria-label={t('choice.result')}>
                {outcomeLines(o.outcome, cardName).map((l) => (
                  <li key={l}>→ {l}</li>
                ))}
              </ul>
            )}
          </div>
        ))}
      </div>
      {confirm && (
        <Confirm
          send={
            send
              ? () => {
                  act(send);
                }
              : undefined
          }
        />
      )}
    </ChoiceFrame>
  );
}

// StartBanner is the start-of-game screen's banner: a headline and a
// caption over the cards (20.22). It grows in as the screen opens.
function StartBanner({ caption }: { caption: string }) {
  const ref = useRef<HTMLElement>(null);
  useLayoutEffect(() => {
    const el = ref.current;
    if (!el || still()) {
      return;
    }
    const s = feel.start;
    // anim: start-banner
    const a = el.animate([{ transform: 'scale(0)' }, { transform: 'scale(1)' }], { duration: ms(s.banner * viewScale(useApp.getState())), easing: s.bannerCurve });
    return () => {
      a.cancel();
    };
  }, []);
  return (
    <header ref={ref} className="start-banner">
      <h2>{t('start.title')}</h2>
      <p>{caption}</p>
    </header>
  );
}

// Confirm sends a start-of-game pick (20.22); it waits disabled until a
// card is picked. The picked card waits at its place for its flight to
// the player; the rest fly back to the deck.
function Confirm({ send }: { send: (() => void) | undefined }) {
  const frame = use(FrameContext);
  return (
    <button
      type="button"
      className="choice-confirm"
      disabled={!send}
      onClick={(e) => {
        const root = e.currentTarget.closest('.choice');
        const kept = root?.querySelector('.choice-card.picked');
        if (!send || !kept) {
          return;
        }
        playSound('choiceConfirm');
        frame.keep(kept, [...(root?.querySelectorAll('.choice-card') ?? [])].filter((el) => el !== kept));
        send();
      }}
    >
      {t('start.confirm')}
    </button>
  );
}

// OptionCard is a card option: the card is the button, named by the
// card (keys, tests, screen readers). Hover lifts it and brightens its
// glow; its zoom waits as usual, or a right click opens it (20.18). A
// pick sounds; cards looked at in a deck then fly back to it, the
// picked one first.
//
// With toggle, a click only picks or unpicks the card, and Confirm sends
// it (20.22): picked lifts it and keeps its bright glow.
function OptionCard({
  card,
  object,
  name,
  deck,
  pick,
  picked = false,
  toggle,
}: {
  card: string;
  object?: number | undefined;
  name: string;
  deck: Kind | undefined;
  pick: () => void;
  picked?: boolean;
  toggle?: () => void;
}) {
  const frame = use(FrameContext);
  return (
    <button
      type="button"
      className={`choice-card${picked ? ' picked' : ''}`}
      aria-label={name}
      aria-pressed={toggle ? picked : undefined}
      data-deal={card}
      data-object={object}
      onClick={(e) => {
        if (toggle) {
          toggle();
          return;
        }
        playSound('choicePick');
        const kept = e.currentTarget;
        const rest = [...(kept.closest('.choice-options')?.querySelectorAll('.choice-card') ?? [])].filter((el) => el !== kept);
        // A look's pick stays on top; the rest go to the bottom (19.147).
        frame.fly([kept, ...rest], 1);
        pick();
      }}
    >
      <DealCard kind={deck}>
        <Card cardRef={card} size="choice" highlight="pick" />
      </DealCard>
    </button>
  );
}

// ChoiceFrame is a choice over the dimmed board, or, while the player
// views the board, a small bar to bring it back (15.29, F-20, H-11).
// title is the heading when it says more than the label. look is the
// deck its cards are looked at in: they are dealt from it (20.18);
// cards says it shows cards, which grow in otherwise. View board
// shrinks the cards away; back, they grow in again, with no deal.
function ChoiceFrame({
  label,
  title,
  banner,
  className = '',
  style,
  look,
  cards = false,
  children,
}: {
  label: string;
  title?: ReactNode;
  // banner replaces the heading: the start-of-game screen's (20.22).
  banner?: ReactNode;
  className?: string;
  style?: CSSProperties | undefined;
  look?: Kind | undefined;
  cards?: boolean;
  children: ReactNode;
}) {
  const [board, setBoard] = useState(false);
  const [shows, setShows] = useState(0);
  const [hiding, setHiding] = useState(false);
  // flown: the pick sent the cards back to the deck; the window waits
  // hidden for the answer. A refusal (nothing sent any more) shows it.
  const [flown, setFlown] = useState(false);
  const waiting = useApp((s) => s.sent !== null);
  const step = useApp((s) => s.game?.step ?? 0);
  const you = useApp(seatOf);
  const [scale] = useState(() => (still() ? 0 : viewScale(useApp.getState())));
  const ref = useRef<HTMLDivElement>(null);
  const entry: Entry = scale === 0 || !cards ? 'none' : shows === 0 && look ? 'deal' : 'grow';
  useEntry(ref, board || (flown && waiting), entry, look, scale, step, shows > 0);
  const frame = useMemo<Frame>(() => ({
    deck: look,
    fly: (order, under) => {
      if (scale === 0 || !look || !deckBox(look)) {
        return;
      }
      flyBack(
        order.map((el) => ({
          card: el.getAttribute('data-deal') ?? el.getAttribute('data-card') ?? '',
          box: (el.querySelector('.deal-face .card') ?? el.querySelector('.card') ?? el).getBoundingClientRect(),
        })),
        look,
        scale,
        under,
      );
      setFlown(true);
    },
    keep: (kept, rest) => {
      if (scale === 0) {
        return;
      }
      const card = kept.getAttribute('data-deal') ?? '';
      const box = (kept.querySelector('.deal-face .card') ?? kept.querySelector('.card') ?? kept).getBoundingClientRect();
      keepCard({ card, box, zone: keptZone(you), step });
      if (look && deckBox(look)) {
        // The rest go to the bottom (19.147).
        flyBack(
          rest.map((el) => ({
            card: el.getAttribute('data-deal') ?? '',
            box: (el.querySelector('.deal-face .card') ?? el.querySelector('.card') ?? el).getBoundingClientRect(),
          })),
          look,
          scale,
          0,
        );
      }
      setFlown(true);
    },
  }), [look, scale, you, step]);
  if (board) {
    return (
      <Peek
        title={title ?? label}
        back={() => {
          setBoard(false);
          setShows(shows + 1);
        }}
      />
    );
  }
  if (flown && waiting) {
    return null;
  }
  const hide = () => {
    const root = ref.current;
    if (scale === 0 || !root || hiding) {
      setBoard(true);
      return;
    }
    setHiding(true);
    const time = ms(feel.choice.hide * scale);
    const easing = feel.choice.growCurve;
    // anim: choice-hide
    const shrink = [...root.querySelectorAll('.card')].map((el) => el.animate([{ transform: 'scale(1)' }, { transform: 'scale(0)' }], { duration: time, easing, fill: 'forwards' }));
    // anim: choice-hide
    const fade = root.closest('.choice-backdrop')?.animate([{ opacity: 1 }, { opacity: 0 }], { duration: time, easing, fill: 'forwards' });
    void (fade ?? shrink.at(0))?.finished
      .catch(() => undefined)
      .then(() => {
        setHiding(false);
        setBoard(true);
      });
  };
  return (
    <FrameContext value={frame}>
      <div className="choice-backdrop">
        <div ref={ref} className={`choice${className ? ` ${className}` : ''}`} style={style} role="dialog" aria-label={label}>
          {banner ?? <h2>{title ?? label}</h2>}
          {children}
          <button type="button" className="choice-view" onClick={hide}>
            {t('choice.view_board')}
          </button>
        </div>
      </div>
    </FrameContext>
  );
}

// useEntry plays the cards' entry once each time the window shows.
function useEntry(ref: RefObject<HTMLDivElement | null>, hidden: boolean, entry: Entry, deck: Kind | undefined, scale: number, step: number, back: boolean) {
  useLayoutEffect(() => {
    const root = ref.current;
    if (hidden || !root || entry === 'none') {
      return;
    }
    return enter(root, entry, deck, scale, step, back);
    // The entry plays as the window shows; later renders keep it.
    // eslint-disable-next-line react-hooks/exhaustive-deps, @eslint-react/exhaustive-deps
  }, [hidden]);
}

// Peek is a hidden choice: its title and a button to bring it back.
function Peek({ title, back }: { title: ReactNode; back: () => void }) {
  return (
    <div className="choice-peek" role="group" aria-label={t('choice.viewing')}>
      <span>{title}</span>
      <button type="button" onClick={back}>
        {t('choice.back')}
      </button>
    </div>
  );
}

// ModeDialog asks a "choose one-" card's mode the way The Bone asks its
// ability: the card large, one line per mode (15.58, V-45). A mode with
// no valid target now stays on the list, greyed. Cancel and a click
// outside send the prompt's cancel; "View board" hides it, like any
// choice (F-20).
export function ModeDialog({ choice, card, modes }: { choice: NonNullable<Actions['choice']>; card: string | undefined; modes: ModeView[] }) {
  const act = useApp((s) => s.act);
  const [board, setBoard] = useState(false);
  const name = useApp((s) => (card ? (s.cards.get(card)?.name ?? card) : choice.title));
  if (board) {
    return (
      <Peek
        title={choice.title}
        back={() => {
          setBoard(false);
        }}
      />
    );
  }
  const { lines, cancel } = modeLines(choice, modes);
  return (
    <CardDialog
      cardRef={card}
      name={name}
      label={choice.title}
      lines={lines.map(({ text, intent }) => ({
        text,
        why: t('choice.no_target'),
        ...(intent
          ? {
              pick: () => {
                act(intent);
              },
            }
          : {}),
      }))}
      onCancel={() => {
        if (cancel) {
          act(cancel);
        }
      }}
      backdrop="choice-backdrop"
      buttons={
        <button
          type="button"
          className="ability-view"
          onClick={() => {
            setBoard(true);
          }}
        >
          {t('choice.view_board')}
        </button>
      }
    />
  );
}

// LookDialog shows the cards an effect let this player look at,
// Hearthstone-style (V-07), until they close it. Hand cards come in
// one labelled row per hand, so the player sees whose they are (19.71).
export function LookDialog({ cards, you, name }: { cards: Looked[]; you: number; name: (seat: number) => string }) {
  const endLook = useApp((s) => s.endLook);
  const row = (cs: Looked[], key?: string) => (
    <div key={key} className="choice-options">
      {cs.map((c) => (
        <Card key={c.id} cardRef={c.card} size="hand" />
      ))}
    </div>
  );
  return (
    <div className="choice look" role="dialog" aria-label={t('look.title')}>
      <h2>{t('look.title')}</h2>
      {lookGroups(cards).map((g) => {
        if (g.owner === undefined) {
          return row(g.cards, 'deck');
        }
        const label = g.owner === you ? t('look.your_hand') : t('look.hand_of', { player: name(g.owner) });
        return (
          <section key={g.owner} className="look-hand" role="group" aria-label={label}>
            <h3>{label}</h3>
            {row(g.cards)}
          </section>
        );
      })}
      <button type="button" onClick={endLook}>
        {t('look.ok')}
      </button>
    </div>
  );
}

// OrderDialog puts the cards of an order prompt back on the deck (19.16).
// A click puts a card in the next free place, top first; a click on a
// placed card takes it back. Nothing is sent until every card has a
// place, so the player sees the whole order before it counts. The cards
// are dealt from the deck, at a choice's size when they fit one row;
// "Put back" sends them back onto it in their order, the top card last
// (20.18).
export function OrderDialog({ order, cards }: { order: NonNullable<Actions['order']>; cards: CardView[] }) {
  const act = useApp((s) => s.act);
  const look = useLook();
  const [placed, setPlaced] = useState<number[]>([]);
  const left = cards.filter((c) => !placed.includes(c.id));
  // One place per card, numbered from the top.
  const places = cards.map((_, i) => ({ n: i + 1, card: cards.find((x) => x.id === placed.at(i)) }));
  const row = rowStyle(cards.length);
  const big = !choiceRow(cards.length).wrap;
  const size = big ? 'choice' : 'table';
  const name = (ref: string) => useApp.getState().cards.get(ref)?.name ?? ref;
  return (
    <ChoiceFrame label={order.title} className={`order${big ? ' order-big' : ''}`} style={big ? row.style : undefined} look={look} cards>
      <p className="order-hint">{t('order.hint')}</p>
      <div className="order-deck">
        {places.map(({ n, card: c }) => (
          <div key={n} className="order-place" role="group" aria-label={t('order.place', { n })}>
            {c ? (
              <div className="order-card" data-card={c.card}>
                <Card
                  cardRef={c.card}
                  size={size}
                  onClick={() => {
                    setPlaced(placed.filter((id) => id !== c.id));
                  }}
                />
              </div>
            ) : (
              <div className="order-empty" />
            )}
            <span>{n === 1 ? t('order.top') : String(n)}</span>
          </div>
        ))}
      </div>
      {left.length > 0 && (
        <div className={`choice-options${big ? ` ${row.className}` : ''}`} role="group" aria-label={t('order.not_placed')}>
          {left.map((c) => (
            <button
              key={c.id}
              type="button"
              className="choice-card"
              aria-label={name(c.card)}
              data-deal={c.card}
              data-object={c.id}
              onClick={() => {
                playSound('choicePick');
                setPlaced([...placed, c.id]);
              }}
            >
              <DealCard kind={look}>
                <Card cardRef={c.card} size={size} highlight="pick" />
              </DealCard>
            </button>
          ))}
        </div>
      )}
      <PutBack
        disabled={left.length > 0}
        send={() => {
          act(orderIntent(order.intent, placed));
        }}
      />
    </ChoiceFrame>
  );
}

// PutBack sends an order; the placed cards fly back onto the deck, the
// bottom one first and the top one last (20.18).
function PutBack({ disabled, send }: { disabled: boolean; send: () => void }) {
  const frame = use(FrameContext);
  return (
    <button
      type="button"
      disabled={disabled}
      onClick={(e) => {
        const root = e.currentTarget.closest('.choice');
        frame.fly([...(root?.querySelectorAll('.order-card') ?? [])].reverse());
        send();
      }}
    >
      {t('order.put_back')}
    </button>
  );
}

// DiscardDialog asks for exactly count cards from the hand. Its cards
// grow in, in one row up to six (20.18).
export function DiscardDialog({
  discard,
  hand,
}: {
  discard: NonNullable<Actions['discard']>;
  hand: CardView[];
}) {
  const act = useApp((s) => s.act);
  const [chosen, setChosen] = useState<number[]>([]);
  const cards = hand.filter((c) => discard.from.includes(c.id));
  const toggle = (id: number) => {
    setChosen(chosen.includes(id) ? chosen.filter((x) => x !== id) : [...chosen, id]);
  };
  const row = rowStyle(cards.length);
  return (
    <ChoiceFrame label={t('discard.label')} title={t('discard.title', { n: discard.count })} cards={cards.length > 0}>
      <div className={`choice-options ${row.className}`} style={row.style}>
        {cards.map((c) => (
          <div key={c.id} className="choice-grow" data-deal={c.card} data-object={c.id}>
            <Card
              cardRef={c.card}
              size="choice"
              highlight={chosen.includes(c.id) ? 'pick' : 'act'}
              onClick={() => {
                toggle(c.id);
              }}
            />
          </div>
        ))}
      </div>
      <button
        type="button"
        disabled={chosen.length !== discard.count}
        onClick={() => {
          act(discardIntent(discard.intent, chosen));
        }}
      >
        {t('discard.button', { chosen: chosen.length, n: discard.count })}
      </button>
    </ChoiceFrame>
  );
}
