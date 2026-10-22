import { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { zonesAt, type Point } from '../aim';
import { boardBox, empty, fromFrame, onFrame, presenceFor, same, spotElement, type Pointing } from '../presence';
import type { Presence } from '../protocol';
import { presenceRenew, useApp } from '../store';
import { AimArrow } from './Aim';

// sendGap is the least time between two presence messages, in ms: under
// the server's 20 a second, so it never drops one (ADR 016).
const sendGap = 70;

// PresenceSender shares what this player points at (15.36, N-15): the
// aim and its arrow's head while aiming; on the player's own turn, the
// hovered card too. A waiting player shares only aims (owner,
// 2026-10-20). It sends on a change, at most every sendGap, and renews
// a still presence so the others keep it.
export function PresenceSender({ active, hand }: { active: boolean; hand: number[] }) {
  const sendPresence = useApp((s) => s.sendPresence);
  const source = useApp((s) => s.aiming?.source);
  const pointingRef = useRef<Pointing & { hand: number[] }>({ active, hand });
  const pointerRef = useRef<Point | null>(null);
  const lastRef = useRef<Omit<Presence, 'seat'>>({});
  const sentAtRef = useRef(0);
  const laterRef = useRef<number | undefined>(undefined);
  const flushRef = useRef((renew: boolean) => {
    const at = pointerRef.current;
    const box = boardBox();
    const p = presenceFor(
      pointingRef.current,
      at ? zonesAt(at.x, at.y) : [],
      pointingRef.current.hand,
      at && box ? onFrame(at.x, at.y, box) : undefined,
    );
    if (renew ? empty(p) : same(p, lastRef.current)) {
      return;
    }
    const wait = sentAtRef.current + sendGap - performance.now();
    if (wait > 0) {
      laterRef.current ??= window.setTimeout(() => {
        laterRef.current = undefined;
        flushRef.current(renew);
      }, wait);
      return;
    }
    sendPresence(p);
    lastRef.current = p;
    sentAtRef.current = performance.now();
  });
  useEffect(() => {
    pointingRef.current = { source, active, hand };
    flushRef.current(false);
  }, [source, active, hand]);
  useEffect(() => {
    const move = (e: PointerEvent) => {
      if (e.pointerType === 'mouse') {
        pointerRef.current = { x: e.clientX, y: e.clientY };
        flushRef.current(false);
      }
    };
    // A pointer that leaves the window points at nothing.
    const leave = () => {
      pointerRef.current = null;
      flushRef.current(false);
    };
    const renew = window.setInterval(() => {
      flushRef.current(true);
    }, presenceRenew);
    window.addEventListener('pointermove', move);
    window.addEventListener('blur', leave);
    document.documentElement.addEventListener('pointerleave', leave);
    return () => {
      window.removeEventListener('pointermove', move);
      window.removeEventListener('blur', leave);
      document.documentElement.removeEventListener('pointerleave', leave);
      window.clearInterval(renew);
      window.clearTimeout(laterRef.current);
      if (!empty(lastRef.current)) {
        sendPresence({});
      }
    };
  }, [sendPresence]);
  return null;
}

interface OtherArrow {
  seat: number;
  from: Point;
  to: Point;
  on: boolean;
}

// middle is an element's centre on screen.
function middle(el: Element): Point {
  const r = el.getBoundingClientRect();
  return { x: r.left + r.width / 2, y: r.top + r.height / 2 };
}

// PresenceLayer shows what the others point at (15.36): their arrow
// while it moves (T-08), a hand card they aim with as a back lifted
// from its slot (T-09, T-22), and the card their arrow or pointer is
// over marked red (T-10, T-21). Marks are data-presence attributes on
// the table's own elements; a spot that is gone shows nothing.
export function PresenceLayer({ you }: { you: number }) {
  const presence = useApp((s) => s.presence);
  const game = useApp((s) => s.game);
  const [arrows, setArrows] = useState<OtherArrow[]>([]);
  // markedRef holds the elements marked now; the next pass clears them.
  const markedRef = useRef<Element[]>([]);
  // Marks and arrows change together, in one pass on the next frame,
  // so a mark never shows beside an older arrow.
  useLayoutEffect(() => {
    const frame = requestAnimationFrame(() => {
      const marks = new Map<Element, Set<string>>();
      const mark = (el: Element, what: string) => {
        marks.set(el, (marks.get(el) ?? new Set()).add(what));
      };
      const box = boardBox();
      const out: OtherArrow[] = [];
      for (const p of Object.values(presence)) {
        if (p.seat === you) {
          continue;
        }
        const from = spotElement(p.from, p.seat);
        const at = spotElement(p.at, p.seat);
        const to = at ? middle(at) : p.point && box ? fromFrame(p.point, box) : undefined;
        if (from && to) {
          out.push({ seat: p.seat, from: middle(from), to, on: !!at });
        }
        if (from && p.from?.hand) {
          mark(from, 'from');
        }
        if (at) {
          mark(at, 'at');
        }
      }
      unmark(markedRef.current);
      for (const [el, what] of marks) {
        el.setAttribute('data-presence', [...what].join(' '));
      }
      markedRef.current = [...marks.keys()];
      setArrows(out);
    });
    return () => {
      cancelAnimationFrame(frame);
    };
  }, [presence, game, you]);
  useEffect(
    () => () => {
      unmark(markedRef.current);
    },
    [],
  );
  return arrows.map((a) => <AimArrow key={a.seat} from={a.from} to={a.to} on={a.on} others />);
}

function unmark(els: Element[]) {
  for (const el of els) {
    el.removeAttribute('data-presence');
  }
}
