import { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { zonesAt, type Point } from '../aim';
import { feel, followStep } from '../feel';
import { playerColor } from '../table';
import { scaleFor } from '../playback';
import { boardBox, empty, fromFrame, marksOf, onFrame, presenceFor, same, sendGap, spotElement, type Pointing } from '../presence';
import type { Presence } from '../protocol';
import { presenceRenew, useApp } from '../store';
import { AimArrow } from './Aim';
import { still } from './lunge';

// PresenceSender shares what this player points at (15.36, N-15): the
// aim and its arrow's head while aiming; on the player's own turn, the
// hovered card too; a dragged hand card, as its slot (20.13). A waiting
// player shares only aims and drags (owner, 2026-10-20). It sends on a
// change, at most every sendGap (0.35 s, 0.25 s while aiming; ADR 027
// "Other players live"), and renews a still presence so the others
// keep it.
// An attack's pick aims from the attacker too (attack, 20.23).
export function PresenceSender({ active, hand, attack }: { active: boolean; hand: number[]; attack?: number }) {
  const sendPresence = useApp((s) => s.sendPresence);
  const source = useApp((s) => s.aiming?.source) ?? attack;
  const held = useApp((s) => s.held ?? undefined);
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
    const wait = sentAtRef.current + sendGap(p) - performance.now();
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
    pointingRef.current = { source, held, active, hand };
    flushRef.current(false);
  }, [source, held, active, hand]);
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

// near: two points within half a pixel.
function near(a: Point, b: Point): boolean {
  return Math.abs(a.x - b.x) < 0.5 && Math.abs(a.y - b.y) < 0.5;
}

// PresenceLayer shows what the others point at (15.36): their arrow
// while it moves (T-08), a hand card they aim with as a back lifted
// from its slot (T-09, T-22), and the card their arrow is over marked
// red (T-10, T-21). Since 20.13 (ADR 027 "Other players live"): a card
// they hover glows in their color, and a hand slot lifts; a card they
// drag floats out of their hand and sways; their arrow's head follows
// each new sample smoothly (followStep), never jumping. Since 20.23
// (ADR 027 "Target arrow"): a new arrow's head glides out from its
// source, and an arrow stays feel.aim.linger after they stop. Marks are
// data-presence attributes on the table's own elements, with the
// player's color as --presence-color; a spot that is gone shows
// nothing. Speed "off" and reduced motion: the head jumps, nothing
// floats or sways; a new arrow is drawn whole at once.
export function PresenceLayer({ you }: { you: number }) {
  const presence = useApp((s) => s.presence);
  const game = useApp((s) => s.game);
  const calm = useApp((s) => scaleFor(s.settings.animation, 0, 0) === 0) || still();
  const [arrows, setArrows] = useState<OtherArrow[]>([]);
  // gone are arrows that stopped, kept until a while after (feel.aim.linger).
  const [gone, setGone] = useState<(OtherArrow & { until: number })[]>([]);
  useEffect(() => {
    if (gone.length === 0) {
      return undefined;
    }
    const first = Math.min(...gone.map((a) => a.until));
    const timer = window.setTimeout(
      () => {
        const now = performance.now();
        setGone((g) => g.filter((a) => a.until > now));
      },
      Math.max(0, first - performance.now()),
    );
    return () => {
      window.clearTimeout(timer);
    };
  }, [gone]);
  // markedRef holds the elements marked now; the next pass clears them.
  const markedRef = useRef<Element[]>([]);
  // targetsRef are the arrows as the last samples place them; shownRef
  // each seat's head as drawn now, on its way to the target.
  const targetsRef = useRef<OtherArrow[]>([]);
  const shownRef = useRef(new Map<number, Point>());
  const frameRef = useRef(0);
  const calmRef = useRef(calm);
  useLayoutEffect(() => {
    calmRef.current = calm;
  }, [calm]);
  // follow draws the arrows, each head a step closer to its target, and
  // asks for the next frame until every head is there.
  const followRef = useRef((dt: number): boolean => {
    const k = calmRef.current ? 1 : followStep(dt);
    const shown = new Map<number, Point>();
    let moving = false;
    const out = targetsRef.current.map((a) => {
      // A new arrow glides out from its source.
      const was = shownRef.current.get(a.seat) ?? (k < 1 ? a.from : undefined);
      let to = a.to;
      if (was && k < 1) {
        to = { x: was.x + (a.to.x - was.x) * k, y: was.y + (a.to.y - was.y) * k };
        if (near(to, a.to)) {
          to = a.to;
        } else {
          moving = true;
        }
      }
      shown.set(a.seat, to);
      return { ...a, to };
    });
    shownRef.current = shown;
    setArrows(out);
    return moving;
  });
  useEffect(
    () => () => {
      cancelAnimationFrame(frameRef.current);
    },
    [],
  );
  // Marks and arrows change together, in one pass on the next frame,
  // so a mark never shows beside an older arrow.
  useLayoutEffect(() => {
    const frame = requestAnimationFrame(() => {
      const marks = new Map<Element, { what: Set<string>; seat: number }>();
      const mark = (el: Element, what: string, seat: number) => {
        const m = marks.get(el) ?? { what: new Set<string>(), seat };
        m.what.add(what);
        marks.set(el, m);
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
        for (const [spot, what] of marksOf(p)) {
          const el = spotElement(spot, p.seat);
          if (el) {
            mark(el, what, p.seat);
          }
        }
      }
      unmark(markedRef.current);
      for (const [el, m] of marks) {
        el.setAttribute('data-presence', [...m.what].join(' '));
        if (el instanceof HTMLElement) {
          el.style.setProperty('--presence-color', playerColor(m.seat) ?? '');
        }
      }
      markedRef.current = [...marks.keys()];
      // An arrow that stopped stays a while where it was drawn; one
      // that comes back takes its place.
      const seats = new Set(out.map((a) => a.seat));
      const until = performance.now() + feel.aim.linger * 1000;
      const stopped = targetsRef.current.filter((a) => !seats.has(a.seat)).map((a) => ({ ...a, to: shownRef.current.get(a.seat) ?? a.to, until }));
      setGone((g) => {
        const kept = g.filter((a) => !seats.has(a.seat) && !stopped.some((s) => s.seat === a.seat));
        return kept.length === g.length && stopped.length === 0 ? g : [...kept, ...stopped];
      });
      targetsRef.current = out;
      cancelAnimationFrame(frameRef.current);
      let last = performance.now();
      const tick = (now: number) => {
        const dt = Math.min((now - last) / 1000, 0.1);
        last = now;
        if (followRef.current(dt)) {
          frameRef.current = requestAnimationFrame(tick);
        }
      };
      // A known arrow moves on from where it is; a new one from its source.
      if (followRef.current(0)) {
        frameRef.current = requestAnimationFrame(tick);
      }
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
  return [
    ...arrows.map((a) => <AimArrow key={a.seat} from={a.from} to={a.to} on={a.on} others />),
    ...gone.map((a) => <AimArrow key={`gone-${String(a.seat)}`} from={a.from} to={a.to} on={a.on} others />),
  ];
}

function unmark(els: Element[]) {
  for (const el of els) {
    el.removeAttribute('data-presence');
    if (el instanceof HTMLElement) {
      el.style.removeProperty('--presence-color');
    }
  }
}
