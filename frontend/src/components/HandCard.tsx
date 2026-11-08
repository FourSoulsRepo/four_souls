// Playing a card from the own hand (20.20, ADR 027 "Playing from the
// hand"). A card is picked up by a drag or a click; on touch a tap lifts
// it in the fan. Held, it grows a little and its glow turns gold over
// the board. A release or a click there plays it; anywhere else, a
// right click or Esc puts it back, and it goes home from where it was.
// A card that needs a target is dragged too: on the board's edge it
// stops and the arrow takes over; back down, the drag goes on. Played,
// it stays where it was let go until the server answers; then it
// shrinks into its stack row (DropLayer), or goes home when refused.
// The hand closes the gap it leaves (useHandGap). Speed "off" and
// reduced motion: no growth, no flight home, no shrink, no gap tween.
import { motion, useMotionValue } from 'motion/react';
import { createPortal } from 'react-dom';
import { useEffect, useLayoutEffect, useRef, useState, type CSSProperties, type MouseEvent as ReactMouseEvent, type RefObject } from 'react';
import { targetAt, zonesAt } from '../aim';
import { arrivedBy } from '../animations';
import { feel, type FanPose } from '../feel';
import { answer, dragged, gapMoves, into, release } from '../hold';
import { useKey } from '../keys';
import { scaleFor } from '../playback';
import type { CardView } from '../protocol';
import { playSound } from '../sounds';
import { useApp, type Dropped } from '../store';
import { swallowAimClick } from './Aim';
import { Card } from './Card';
import { holdCard, type Held } from './dragFeel';
import type { Interact } from './Interact';
import { still } from './lunge';

// Mode is how a hand card is held: dragged with the button down;
// following the pointer after a click; lifted in the fan after a tap;
// stopped at the board's edge while it aims; or played and waiting for
// the server's answer. null: in the fan.
type Mode = 'drag' | 'follow' | 'lift' | 'stopped' | 'wait';

interface Point {
  x: number;
  y: number;
}

// onPlayArea: the point is over the board, where hand cards are played.
function onPlayArea(x: number, y: number): boolean {
  return document.elementsFromPoint(x, y).some((el) => el.closest('[data-drop="play"]') !== null);
}

// Pose is the parts of a hand card's place that its motions tween.
interface Pose {
  transform: string;
  translate: string;
  scale: string;
}

function poseOf(el: HTMLElement): Pose {
  const cs = getComputedStyle(el);
  return { transform: cs.transform, translate: cs.translate, scale: cs.scale };
}

// calm: animation speed "off" or reduced motion, so no new motion. It
// is read as each motion starts.
function calm(): boolean {
  return scaleFor(useApp.getState().settings.animation, 0, 0) === 0 || still();
}

// HandCard is a card in the viewer's hand, turned and dropped by its
// place in the fan (pose). Drag is plain pointer events with pointer
// capture, so mouse and touch share it and the release is never missed
// (ADR 008). A held card tilts by its motion and two loops follow it
// (20.6, holdCard); its inner box turns, so its own move stays apart.
export function HandCard({ card, ui, pose }: { card: CardView; ui: Interact; pose: FanPose }) {
  const act = useApp((s) => s.act);
  const playAt = useApp((s) => s.playAt);
  const setZoom = useApp((s) => s.setZoom);
  const setHeld = useApp((s) => s.setHeld);
  const setAiming = useApp((s) => s.setAiming);
  const setDropped = useApp((s) => s.setDropped);
  const aiming = useApp((s) => s.aiming?.source === card.id);
  const [mode, setMode] = useState<Mode | null>(null);
  const [offset, setOffset] = useState<Point | null>(null);
  const [over, setOver] = useState(false);
  // returning counts the flights home; each runs once (hand-return).
  const [returning, setReturning] = useState(0);
  const modeRef = useRef<Mode | null>(null);
  const elRef = useRef<HTMLDivElement>(null);
  const startRef = useRef<{ x: number; y: number; pointer: number } | null>(null);
  // pressRef is the kind of pointer of the last press: a click by mouse
  // picks the card up, a tap lifts it.
  const pressRef = useRef('mouse');
  // originRef is where a click picked the card up; it follows from there.
  const originRef = useRef<Point>({ x: 0, y: 0 });
  // lastRef is the drag's last offset off the board: a card that aims
  // stops there.
  const lastRef = useRef<Point>({ x: 0, y: 0 });
  // draggedRef swallows the click that ends a drag.
  const draggedRef = useRef(false);
  // heldRef is the held card's tilt and loops.
  const heldRef = useRef<Held | null>(null);
  // fromRef is the pose a flight home starts from.
  const fromRef = useRef<Pose | null>(null);
  // waitRef is the step the card was played at.
  const waitRef = useRef(0);
  const pitch = useMotionValue(0);
  const roll = useMotionValue(0);
  // A card just drawn is playable only a moment after it lands (20.7).
  const [landing, setLanding] = useState(() => arrivedBy(card.card) === 'loot-draw');
  useEffect(() => {
    if (!landing) {
      return;
    }
    const id = window.setTimeout(() => {
      setLanding(false);
    }, feel.draw.ready * 1000);
    return () => {
      window.clearTimeout(id);
    };
  }, [landing]);
  const action = landing ? undefined : ui.play(card.id);
  const play = action && !action.targets?.length ? action.intent : undefined;
  const aims = action?.targets?.length ? action : undefined;
  const handle = ui.object(card.id);
  const usable = !landing && handle.highlight === 'act';
  // liveRef is what the handlers outside a render read: the latest
  // card and its play.
  const liveRef = useRef({ card, play });
  useLayoutEffect(() => {
    liveRef.current = { card, play };
  });

  const to = (m: Mode | null) => {
    modeRef.current = m;
    setMode(m);
  };
  const feelOn = (x: number, y: number) => {
    heldRef.current?.end();
    heldRef.current = holdCard(x, y, calm(), (p, r) => {
      pitch.set(p);
      roll.set(r);
    });
  };
  const feelOff = () => {
    heldRef.current?.end();
    heldRef.current = null;
  };
  // pickUp holds the card: it grows to its held size and sounds.
  const pickUp = (m: 'drag' | 'follow' | 'lift', x: number, y: number) => {
    playSound('cardGrab');
    setZoom(null);
    setHeld(card.id);
    to(m);
    if (m === 'lift') {
      return;
    }
    originRef.current = { x, y };
    feelOn(x, y);
    const el = elRef.current;
    if (el && !calm()) {
      // anim: hand-hold
      el.animate([{ scale: '1' }, { scale: String(feel.play.grow) }], { duration: feel.play.growTime * 1000, easing: feel.play.growCurve });
    }
  };
  // goHome puts the card back in the fan, flying from where it is.
  const goHome = () => {
    const el = elRef.current;
    startRef.current = null;
    feelOff();
    setHeld(null);
    if (el && !calm() && modeRef.current !== null) {
      fromRef.current = poseOf(el);
      for (const a of el.getAnimations()) {
        a.cancel();
      }
      setReturning((n) => n + 1);
    }
    to(null);
    setOffset(null);
    setOver(false);
  };
  // toWait holds a played card where it is until the server answers.
  const toWait = () => {
    feelOff();
    setHeld(null);
    setOver(false);
    startRef.current = null;
    waitRef.current = useApp.getState().game?.step ?? 0;
    to('wait');
    const el = elRef.current;
    const c = liveRef.current.card;
    if (el && !calm()) {
      const r = el.getBoundingClientRect();
      setDropped({ id: c.id, card: c.card, step: waitRef.current, x: r.left, y: r.top, w: r.width, h: r.height });
    }
  };

  // The flight home, once the fan's pose is drawn.
  useLayoutEffect(() => {
    const el = elRef.current;
    const from = fromRef.current;
    if (returning === 0 || !el || !from) {
      return;
    }
    fromRef.current = null;
    // anim: hand-return
    const a = el.animate([{ ...from }, { ...poseOf(el) }], { duration: feel.play.back * 1000, easing: feel.play.backCurve });
    const n = returning;
    a.onfinish = () => {
      setReturning((now) => (now === n ? 0 : now));
    };
  }, [returning]);

  // The store says how a play goes: an aim that ends with a pick plays
  // the card, which then waits; a cancelled aim sends a card stopped at
  // the board's edge home. A waiting card goes home when the server
  // refuses it, or answers and leaves it in the hand; one that leaves
  // the hand is DropLayer's.
  useEffect(
    () =>
      useApp.subscribe((s, prev) => {
        const id = liveRef.current.card.id;
        if (!(s.game?.view.players ?? []).some((p) => (p.hand ?? []).some((c) => c.id === id))) {
          return;
        }
        const now = s.aiming?.source === id;
        const m = modeRef.current;
        if (prev.aiming?.source === id && !now) {
          if (s.sent !== null && s.sent === s.game?.step) {
            toWait();
          } else if (m === 'stopped') {
            goHome();
          }
          return;
        }
        if (m !== 'wait') {
          return;
        }
        // Its target is asked on the table: the aim lifts it.
        if (now) {
          setDropped(null);
          goHome();
          return;
        }
        if (s.game?.step === prev.game?.step && s.sent === prev.sent) {
          return;
        }
        const step = s.game?.step ?? 0;
        const more = s.game?.view.activating === id && (s.aimed !== null || s.sent === step);
        const a = answer({ step: waitRef.current, now: step, sent: s.sent, more });
        if (a === 'wait') {
          waitRef.current = step;
          return;
        }
        if (a === 'refused') {
          playSound('playRefused');
        }
        setDropped(null);
        goHome();
      }),
    // The handlers read only refs and setters that never change.
    // eslint-disable-next-line react-hooks/exhaustive-deps, @eslint-react/exhaustive-deps
    [],
  );

  // A card held by a click follows the pointer (a tapped one stays
  // lifted); the next click plays it on the board or puts it back, and
  // so does a right click.
  useEffect(() => {
    if (mode !== 'follow' && mode !== 'lift') {
      return;
    }
    const move = (e: PointerEvent) => {
      if (modeRef.current !== 'follow') {
        return;
      }
      const o = originRef.current;
      heldRef.current?.move(e.clientX, e.clientY);
      setOffset({ x: e.clientX - o.x, y: e.clientY - o.y });
      setOver(onPlayArea(e.clientX, e.clientY));
    };
    const click = (e: MouseEvent) => {
      e.stopPropagation();
      e.preventDefault();
      const button = e.target instanceof Element && e.target.closest('button') !== null;
      const intent = liveRef.current.play;
      if (release(onPlayArea(e.clientX, e.clientY), button) === 'play' && intent) {
        act(intent);
        toWait();
        return;
      }
      goHome();
    };
    const menu = (e: MouseEvent) => {
      e.preventDefault();
      e.stopPropagation();
      goHome();
    };
    window.addEventListener('pointermove', move);
    window.addEventListener('click', click, true);
    window.addEventListener('contextmenu', menu, true);
    return () => {
      window.removeEventListener('pointermove', move);
      window.removeEventListener('click', click, true);
      window.removeEventListener('contextmenu', menu, true);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps, @eslint-react/exhaustive-deps
  }, [mode]);

  useKey('escape', () => {
    const m = modeRef.current;
    if (m === 'drag' || m === 'follow' || m === 'lift') {
      goHome();
    }
  });

  // A card that leaves the hand while held drops its loops; the others
  // no longer see it held.
  useEffect(
    () => () => {
      if (heldRef.current) {
        setHeld(null);
      }
      heldRef.current?.end();
    },
    [setHeld],
  );

  // A click: a playable card is picked up (a tap lifts it); one that
  // needs a target starts aiming (15.25). A card that can't be played
  // says why (handle). A card in a pick answers the pick.
  const onClick = (e?: ReactMouseEvent) => {
    if (modeRef.current !== null) {
      return;
    }
    if (play) {
      pickUp(pressRef.current === 'mouse' ? 'follow' : 'lift', e?.clientX ?? 0, e?.clientY ?? 0);
      return;
    }
    if (aims) {
      playSound('cardGrab');
      setAiming({ source: card.id, intent: aims.intent, targets: aims.targets ?? [], ...(aims.reason ? { reason: aims.reason } : {}) });
    }
  };
  const click = handle.highlight === 'act' ? onClick : handle.onClick;

  const off = offset !== null && (mode === 'drag' || mode === 'follow' || mode === 'stopped' || mode === 'wait');
  const cls = [
    'hand-card',
    play || aims ? 'playable' : '',
    usable ? 'usable' : '',
    landing ? 'landing' : '',
    mode === 'drag' || mode === 'follow' ? 'dragging' : '',
    mode === 'follow' || mode === 'lift' ? 'held' : '',
    mode === 'lift' || (mode === 'wait' && offset === null) ? 'lifted' : '',
    mode === 'stopped' ? 'stopped' : '',
    mode === 'wait' ? 'waiting' : '',
    over ? 'over' : '',
    returning > 0 ? 'returning' : '',
    aiming ? 'aiming' : '',
  ]
    .filter(Boolean)
    .join(' ');
  const style: CSSProperties = off
    ? {
        transform: 'none',
        translate: `${String(offset.x)}px ${String(offset.y)}px`,
        scale: String(feel.play.grow),
        perspective: `${String(feel.drag.perspective)}px`,
      }
    : ({ '--fan-angle': `${String(pose.angle)}deg`, '--fan-drop': `${String(pose.drop)}px` } as CSSProperties);
  return (
    <div
      ref={elRef}
      className={cls}
      data-hand={card.id}
      style={style}
      onPointerDown={(e) => {
        pressRef.current = e.pointerType;
        const m = modeRef.current;
        if ((m !== null && m !== 'lift') || e.button !== 0 || (!play && !aims)) {
          return;
        }
        startRef.current = { x: e.clientX, y: e.clientY, pointer: e.pointerId };
        draggedRef.current = false;
      }}
      onPointerMove={(e) => {
        const start = startRef.current;
        if (start?.pointer !== e.pointerId) {
          return;
        }
        const x = e.clientX - start.x;
        const y = e.clientY - start.y;
        const m = modeRef.current;
        if (m === null || m === 'lift') {
          if (!dragged(x, y)) {
            return;
          }
          e.currentTarget.setPointerCapture(e.pointerId);
          draggedRef.current = true;
          if (m === 'lift') {
            to('drag');
            feelOn(e.clientX, e.clientY);
          } else {
            pickUp('drag', start.x, start.y);
          }
        }
        const board = onPlayArea(e.clientX, e.clientY);
        // A card that needs a target stops at the board's edge and aims;
        // back below it, the drag goes on.
        if (aims && board && modeRef.current === 'drag') {
          feelOff();
          setHeld(null);
          setOver(false);
          setOffset(lastRef.current);
          to('stopped');
          setAiming({ source: card.id, intent: aims.intent, targets: aims.targets ?? [], ...(aims.reason ? { reason: aims.reason } : {}) });
          return;
        }
        if (modeRef.current === 'stopped') {
          if (board) {
            return;
          }
          // The mode first: the aim's end then sends nothing home.
          to('drag');
          setAiming(null);
          setHeld(card.id);
          feelOn(e.clientX, e.clientY);
        }
        lastRef.current = { x, y };
        heldRef.current?.move(e.clientX, e.clientY);
        setOffset({ x, y });
        setOver(board);
      }}
      onPointerUp={(e) => {
        const m = modeRef.current;
        if (startRef.current?.pointer !== e.pointerId) {
          return;
        }
        startRef.current = null;
        if (m === 'drag') {
          if (play && onPlayArea(e.clientX, e.clientY)) {
            act(play);
            toWait();
          } else {
            goHome();
          }
        } else if (m === 'stopped') {
          // Let go on a target, it plays there; elsewhere the aim stays
          // and the next click picks.
          swallowAimClick();
          const a = useApp.getState().aiming;
          const target = a ? targetAt(zonesAt(e.clientX, e.clientY), a.targets) : undefined;
          if (a?.intent && target) {
            playAt(a.intent, target);
          }
        }
      }}
      onPointerCancel={() => {
        if (modeRef.current === 'drag') {
          goHome();
        }
        startRef.current = null;
      }}
      onClickCapture={(e) => {
        if (draggedRef.current) {
          draggedRef.current = false;
          e.stopPropagation();
        }
      }}
    >
      {/* anim: drag-tilt */}
      <motion.div className="drag-tilt" style={{ rotateX: pitch, rotateY: roll }}>
        <Card cardRef={card.card} size="hand" objectId={card.id} rise {...handle} onClick={click} />
      </motion.div>
    </div>
  );
}

// DropLayer shows an own played card shrinking into its stack row once
// the server took it (20.20, ADR 027 "Playing from the hand" 12): from
// where it was held, as the view without it is drawn.
export function DropLayer({ you }: { you: number }) {
  const dropped = useApp((s) => s.dropped);
  const gone = useApp((s) => {
    const d = s.dropped;
    if (!d || !s.game || s.game.step === d.step) {
      return false;
    }
    return !(s.game.view.players.at(you)?.hand ?? []).some((c) => c.id === d.id);
  });
  if (!dropped || !gone) {
    return null;
  }
  return <DropGhost key={`${String(dropped.id)}-${String(dropped.step)}`} dropped={dropped} />;
}

// rowOf is where a played card goes: its stack row (the newest item of
// its card), else its place on the table (a trinket), else the loot
// discard (it resolved at once).
function rowOf(d: Dropped): DOMRect | null {
  const items = (useApp.getState().game?.view.stack ?? []).filter((it) => it.card === d.card);
  const top = items.reduce<number | undefined>((m, it) => (m === undefined || it.id > m ? it.id : m), undefined);
  const sels = [top === undefined ? '' : `[data-stack="${String(top)}"]`, `[data-zone="obj-${String(d.id)}"]`, '[data-zone="discard-loot"]'];
  for (const sel of sels.filter(Boolean)) {
    const el = document.querySelector(sel);
    if (el) {
      return el.getBoundingClientRect();
    }
  }
  return null;
}

function DropGhost({ dropped }: { dropped: Dropped }) {
  const setDropped = useApp((s) => s.setDropped);
  const ref = useRef<HTMLDivElement>(null);
  useLayoutEffect(() => {
    const el = ref.current;
    if (!el) {
      return;
    }
    const row = rowOf(dropped);
    const p = into(
      { x: dropped.x, y: dropped.y, w: dropped.w, h: dropped.h },
      { w: el.offsetWidth, h: el.offsetHeight },
      row ? { x: row.left, y: row.top, w: row.width, h: row.height } : null,
    );
    // anim: hand-play
    const a = el.animate(
      [
        { transform: `translate(0px, 0px) scale(${String(p.from)})`, opacity: 1 },
        // It fades only at the very end, once it is row-sized.
        { offset: feel.play.intoFade, opacity: 1 },
        { transform: `translate(${String(p.dx)}px, ${String(p.dy)}px) scale(${String(p.to)})`, opacity: 0 },
      ],
      { duration: feel.play.into * 1000, easing: feel.play.intoCurve, fill: 'forwards' },
    );
    a.onfinish = () => {
      if (useApp.getState().dropped === dropped) {
        setDropped(null);
      }
    };
    return () => {
      a.cancel();
    };
  }, [dropped, setDropped]);
  // On the body: inside the table a transformed box would move it.
  return createPortal(
    <div ref={ref} className="hand-drop" style={{ left: dropped.x, top: dropped.y, transform: 'scale(0)' }} aria-hidden="true">
      <Card cardRef={dropped.card} size="hand" noZoom />
    </div>,
    document.body,
  );
}

// useHandGap closes the gap a card leaves in the own hand (20.20, ADR
// 027 "Playing from the hand" 13): the cards that stay move from their
// old places to their new ones. A held card keeps its own place.
export function useHandGap(hand: RefObject<HTMLDivElement | null>, ids: string): void {
  const lastRef = useRef(new Map<number, number>());
  useLayoutEffect(() => {
    const el = hand.current;
    if (!el) {
      return;
    }
    const now = new Map<number, number>();
    const cards = [...el.querySelectorAll<HTMLElement>(':scope > .hand-card')];
    for (const c of cards) {
      now.set(Number(c.dataset.hand), c.offsetLeft);
    }
    const moves = calm() ? [] : gapMoves(lastRef.current, now);
    lastRef.current = now;
    for (const m of moves) {
      const c = cards.find((x) => Number(x.dataset.hand) === m.id);
      if (!c || c.classList.contains('waiting') || c.classList.contains('dragging')) {
        continue;
      }
      // anim: hand-gap
      c.animate([{ translate: `${String(m.dx)}px 0px` }, { translate: '0px 0px' }], { duration: feel.play.gap * 1000, easing: feel.play.gapCurve });
    }
  }, [hand, ids]);
}
