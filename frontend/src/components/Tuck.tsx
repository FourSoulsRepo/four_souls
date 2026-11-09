import { useLayoutEffect, useRef, useState } from 'react';
import type { Kind } from '../cards';
import { dealDelays } from '../choice';
import { feel } from '../feel';
import type { Playing } from '../store';
import { raised, tuckClip, tuckPose, tuckSide, tuckTime, type Spot } from '../tuck';
import { Card } from './Card';
import { deckBox } from './ChoiceBack';
import { still } from './lunge';

// The size of a table card, which the moving cards are drawn at.
const cardW = 80;
const cardH = 112;
// Most cards drawn; the tag counts them all.
const most = 5;
const ms = (s: number) => s * 1000;

// lookBacks are the backs another player's look holds over the deck
// (Pile.tsx, 20.18): the last n of them, which this effect takes over.
function lookBacks(deck: Kind, n: number): HTMLElement[] {
  const all = [...document.querySelectorAll<HTMLElement>(`[data-zone="deck-${deck}"] .deck-look-back`)];
  return n > 0 ? all.slice(-n) : [];
}

// spotOf is an element's spot on the screen: its centre, its width
// without its tilt, and its tilt from its inline transform.
function spotOf(el: HTMLElement): Spot {
  const r = el.getBoundingClientRect();
  const rotate = Number(/rotate\((-?[\d.]+)deg\)/.exec(el.style.transform)?.[1] ?? 0);
  return { x: r.left + r.width / 2, y: r.top + r.height / 2, w: el.offsetWidth || r.width, rotate };
}

// frame is a pose as a keyframe of a card drawn at table size.
function frame(t: number, from: Spot, pile: Spot, bottom: boolean, side: 1 | -1): Keyframe {
  const p = tuckPose(t, from, pile, bottom, side);
  const at = `translate(${(p.x - cardW / 2).toFixed(1)}px, ${(p.y - cardH / 2).toFixed(1)}px) rotate(${p.rotate.toFixed(1)}deg) scale(${(p.w / cardW).toFixed(4)})`;
  return { transform: at, clipPath: tuckClip(p.hide, side), opacity: p.opacity };
}

// Tuck shows where looked-at deck cards went (19.147): each drops onto
// the pile or slides out and under it; a tag on the pile says where.
// Others' backs over the deck are taken over where they are; the
// chooser's own cards rise from a small fan over the pile.
export function Tuck({ deck, cards, tag, playing }: { deck: Kind; cards: { card?: string; bottom: boolean }[]; tag: string; playing: Playing }) {
  const ref = useRef<HTMLDivElement>(null);
  const [pile] = useState(() => {
    const r = deckBox(deck);
    return r && { x: r.left + r.width / 2, y: r.top + r.height / 2, w: r.width, rotate: 0 };
  });
  // Read once, when the effect starts, as the pile.
  const [shown] = useState(() => (still() ? [] : cards.slice(0, most)));
  useLayoutEffect(() => {
    const root = ref.current;
    if (!root || !pile) {
      playing.done();
      return;
    }
    const k = feel.tuck;
    const scale = playing.scale;
    const side = tuckSide(pile, window.innerWidth);
    const backs = lookBacks(deck, shown.length);
    const delays = dealDelays(shown.length);
    const parts: Animation[] = [];
    shown.forEach((c, i) => {
      const el = root.querySelector<HTMLElement>(`[data-i="${String(i)}"]`);
      const back = backs.at(i - (shown.length - backs.length));
      const from = back ? spotOf(back) : raised(pile, i, shown.length);
      if (back) {
        back.style.visibility = 'hidden';
      }
      if (!el) {
        return;
      }
      const time = tuckTime(c.bottom);
      const n = Math.max(2, Math.ceil(time * 60));
      const frames = Array.from({ length: n + 1 }, (_, j) => frame((time * j) / n, from, pile, c.bottom, side));
      // anim: tuck
      parts.push(el.animate(frames, { duration: ms(time * scale), delay: ms((delays.at(i) ?? 0) * scale), fill: 'both' }));
    });
    const fade = k.fade / k.tag;
    // anim: tuck
    const label = root.querySelector('.effect-tuck-tag')?.animate(
      [
        { opacity: 0, offset: 0 },
        { opacity: 1, offset: fade },
        { opacity: 1, offset: 1 - fade },
        { opacity: 0, offset: 1 },
      ],
      { duration: ms(k.tag * scale), fill: 'both' },
    );
    // A cancel (an unmount, or React's StrictMode remount) ends nothing.
    let stopped = false;
    const all = [...parts, ...(label ? [label] : [])];
    Promise.all(all.map((a) => a.finished)).then(
      () => {
        if (!stopped) {
          playing.done();
        }
      },
      () => undefined,
    );
    return () => {
      stopped = true;
      for (const a of all) {
        a.cancel();
      }
    };
  }, [deck, pile, playing, shown]);
  if (!pile) {
    return null;
  }
  return (
    <div ref={ref} className="effect-tuck" data-deck={deck}>
      {shown.map((c, i) => (
        // eslint-disable-next-line @eslint-react/no-array-index-key
        <div key={i} className="effect-fly" data-i={i} data-where={c.bottom ? 'bottom' : 'top'} data-face={c.card}>
          <Card cardRef={c.card} back={deck} size="table" />
        </div>
      ))}
      <div className="effect-tuck-tag" style={{ left: pile.x, top: pile.y }}>
        {tag}
      </div>
    </div>
  );
}
