import { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { create } from 'zustand';
import type { Kind } from '../cards';
import { backDelays } from '../choice';
import { feel } from '../feel';
import { useApp } from '../store';
import { Card } from './Card';

// A choice's cards fly back onto their deck after the pick (20.18, ADR
// 027 "Choosing cards" 12.10): copies of them leave the window's places
// and land on the deck, turning face down. The first card of the list
// flies first and lies on top of the rest.

interface Flight {
  key: number;
  deck: Kind;
  scale: number;
  cards: { card: string; box: DOMRect }[];
}

const useFlight = create<{ flight: Flight | null }>(() => ({ flight: null }));

let next = 0;

// deckBox is a deck's pile on the board, or null when it is not drawn.
export function deckBox(kind: Kind): DOMRect | null {
  const el = document.querySelector(`[data-zone="deck-${kind}"]`);
  return (el?.querySelector('.card') ?? el)?.getBoundingClientRect() ?? null;
}

// flyBack sends cards back onto a deck: each card from its box on the
// screen, in this order. scale is the animation speed's.
export function flyBack(cards: { card: string; box: DOMRect }[], deck: Kind, scale: number): void {
  next++;
  useFlight.setState({ flight: { key: next, deck, scale, cards } });
}

// Kept is a start-of-game pick after Confirm (20.22): the card waits
// at its place in the window, under zone, until playback flies it to
// its player from there. step is the update it was picked in.
interface Kept {
  key: number;
  card: string;
  box: DOMRect;
  zone: string;
  step: number;
}

const useKept = create<{ kept: Kept | null }>(() => ({ kept: null }));

// keepCard holds a picked card at its place for its flight (20.22).
export function keepCard(k: Omit<Kept, 'key'>): void {
  next++;
  useKept.setState({ kept: { ...k, key: next } });
}

// ChoiceBack draws the cards flying back, while they fly, and a kept
// card while it waits for its flight.
export function ChoiceBack() {
  const flight = useFlight((s) => s.flight);
  const kept = useKept((s) => s.kept);
  return (
    <>
      {flight && <Flying key={flight.key} flight={flight} />}
      {kept && <Held key={kept.key} kept={kept} />}
    </>
  );
}

// Held is the kept card at its place. It hides once its flight starts
// and goes when the flight ends, when the update is drawn with no
// flight (speed "off"), or after feel.start.hold at most.
function Held({ kept }: { kept: Kept }) {
  const flying = useApp((s) => {
    const e = s.playing?.effect;
    return e?.kind === 'fly' && e.from[0] === kept.zone;
  });
  const drawn = useApp((s) => (s.game?.step ?? 0) > kept.step && s.playing === null);
  // seen: its flight started; kept while the render follows it.
  const [seen, setSeen] = useState(false);
  if (flying && !seen) {
    setSeen(true);
  }
  useEffect(() => {
    const end = () => {
      if (useKept.getState().kept?.key === kept.key) {
        useKept.setState({ kept: null });
      }
    };
    if ((seen && !flying) || drawn) {
      end();
      return;
    }
    const id = window.setTimeout(end, feel.start.hold * 1000);
    return () => {
      window.clearTimeout(id);
    };
  }, [seen, flying, drawn, kept.key]);
  const { box } = kept;
  return (
    <div className="choice-back" data-testid="choice-kept">
      <div
        className="choice-back-card choice-kept"
        data-zone={kept.zone}
        style={{ left: box.left, top: box.top, width: box.width, height: box.height, visibility: flying || seen ? 'hidden' : undefined }}
      >
        <Card cardRef={kept.card} size="choice" highlight="pick" />
      </div>
    </div>
  );
}

function Flying({ flight }: { flight: Flight }) {
  const ref = useRef<HTMLDivElement>(null);
  useLayoutEffect(() => {
    const root = ref.current;
    const deck = deckBox(flight.deck);
    const end = () => {
      if (useFlight.getState().flight?.key === flight.key) {
        useFlight.setState({ flight: null });
      }
    };
    if (!root || !deck) {
      end();
      return;
    }
    const c = feel.choice;
    const delays = backDelays(flight.cards.length);
    const parts = flight.cards.flatMap(({ box }, i) => {
      const el = root.querySelector<HTMLElement>(`[data-i="${String(i)}"]`);
      if (!el) {
        return [];
      }
      const dx = deck.left + deck.width / 2 - (box.left + box.width / 2);
      const dy = deck.top + deck.height / 2 - (box.top + box.height / 2);
      const s = deck.width / Math.max(box.width, 1);
      const timing = { duration: c.back * flight.scale * 1000, delay: (delays.at(i) ?? 0) * flight.scale * 1000, easing: c.backCurve, fill: 'both' } as const;
      return [
        // anim: choice-back
        el.animate([{ transform: 'none' }, { transform: `translate(${String(dx)}px, ${String(dy)}px) scale(${String(s)})` }], timing),
        // anim: choice-back
        el.querySelector('.deal-turn')?.animate([{ transform: 'rotateY(0deg)' }, { transform: 'rotateY(180deg)' }], timing),
      ];
    });
    // A cancel (an unmount, or React's StrictMode remount) ends nothing.
    let stopped = false;
    Promise.all(parts.flatMap((p) => (p ? [p.finished] : []))).then(
      () => {
        if (!stopped) {
          end();
        }
      },
      () => undefined,
    );
    return () => {
      stopped = true;
      for (const p of parts) {
        p?.cancel();
      }
    };
  }, [flight]);
  // The first card is drawn last, so it lies on top.
  const order = flight.cards.map((c, i) => ({ ...c, i })).reverse();
  return (
    <div ref={ref} className="choice-back" data-testid="choice-back">
      {order.map(({ card, box, i }) => (
        <div
          key={i}
          className="choice-back-card"
          data-i={i}
          style={{ left: box.left, top: box.top, width: box.width, height: box.height }}
        >
          <div className="deal-turn">
            <div className="deal-face">
              <Card cardRef={card} size="choice" />
            </div>
            <div className="deal-face deal-back">
              <Card back={flight.deck} size="choice" />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
