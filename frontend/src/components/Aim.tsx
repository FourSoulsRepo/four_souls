import { useEffect, useLayoutEffect, useRef, useState, type CSSProperties } from 'react';
import { aimedOption } from '../actions';
import {
  aimClick,
  arc,
  cancelOption,
  chainOf,
  dashCount,
  dashFrames,
  dashPose,
  lastPointer,
  marchTime,
  setArrowShows,
  targetAt,
  targetElement,
  zonesAt,
  type Point,
} from '../aim';
import { feel } from '../feel';
import { t } from '../i18n';
import { useKey } from '../keys';
import { scaleFor } from '../playback';
import type { Intent, OptionView } from '../protocol';
import { playSound } from '../sounds';
import { useApp, type Aiming } from '../store';
import { still } from './lunge';
import { rem } from '../uiScale';

// AimLayer is the player's aim (15.25): a click on a card that needs a
// target starts it, the next click picks the target. A click on empty
// space, a right click, or Esc cancels it. A click on the card again
// does nothing, unless it may target itself. A click on a card, pile,
// or stack row that is no target keeps aiming and says why (15.33), in
// the shared WhyNote (15.29). While a mouse aims, an arrow runs from
// the card to the pointer (useArrow). A card already placed answers its
// target prompt; cancelling it sends the prompt's cancel, and the card
// goes back (15.34, T-13). A hand card dragged onto the board aims too
// (20.20): letting it go there decides by itself, so the click that
// follows the release is no pick (swallowAimClick).
export function AimLayer() {
  const aiming = useApp((s) => s.aiming);
  return aiming ? <Aim key={aiming.source} aiming={aiming} /> : null;
}

// swallow is set when a drag that aims ends: the click its release
// makes is no click of the aim. A new press clears it.
let swallow = false;

// swallowAimClick drops the click that ends a drag which aims (20.20).
export function swallowAimClick(): void {
  swallow = true;
}

function Aim({ aiming }: { aiming: Aiming }) {
  const setAiming = useApp((s) => s.setAiming);
  const playAt = useApp((s) => s.playAt);
  const act = useApp((s) => s.act);
  const showWhy = useApp((s) => s.showWhy);
  const line = useArrow(aiming.source, aiming.targets);
  useKey('escape', () => {
    cancel(aiming, act, setAiming);
  });
  useEffect(() => {
    // Capture: the click picks or cancels and reaches nothing else.
    const click = (e: MouseEvent) => {
      e.stopPropagation();
      e.preventDefault();
      if (swallow) {
        swallow = false;
        return;
      }
      const what = aimClick(zonesAt(e.clientX, e.clientY), aiming.source, aiming.targets);
      if (what.kind === 'pick') {
        pick(aiming, what.target);
      } else if (what.kind === 'wrong') {
        showWhy(aiming.reason ?? t('actions.bad_target'), e.clientX, e.clientY);
      } else if (what.kind === 'cancel') {
        cancel(aiming, act, setAiming);
      }
    };
    const pick = (a: Aiming, target: OptionView) => {
      if (a.intent) {
        playAt(a.intent, target);
        return;
      }
      // Sent first: the card that aimed sees its play sent as the aim
      // ends, and waits for the answer (20.20).
      const answer = a.choice && aimedOption(a.choice, target);
      if (answer) {
        act(answer);
      }
      setAiming(null);
    };
    const menu = (e: MouseEvent) => {
      e.preventDefault();
      cancel(aiming, act, setAiming);
    };
    const press = () => {
      swallow = false;
    };
    window.addEventListener('pointerdown', press, true);
    window.addEventListener('click', click, true);
    window.addEventListener('contextmenu', menu, true);
    return () => {
      window.removeEventListener('pointerdown', press, true);
      window.removeEventListener('click', click, true);
      window.removeEventListener('contextmenu', menu, true);
    };
  }, [aiming, act, playAt, setAiming, showWhy]);
  return line && <AimArrow {...line} />;
}

// cancel ends an aim; a placed card's prompt is cancelled with it. It
// sounds (20.20, ADR 027 "Playing from the hand" 10).
function cancel(aiming: Aiming, act: (i: Intent) => void, setAiming: (a: Aiming | null) => void) {
  playSound('aimCancel');
  const c = aiming.choice && cancelOption(aiming.choice);
  setAiming(null);
  if (c) {
    act(c);
  }
}

// AttackAim is the arrow of an attack's pick (20.23, ADR 027 "Target
// arrow" 16.14): from the attacking character to the pointer, while the
// player picks what to attack on the table. The pick stays as it was:
// the targets glow blue and a click on one answers it.
export function AttackAim({ source, targets }: { source: number; targets: OptionView[] }) {
  const line = useArrow(source, targets);
  return line && <AimArrow {...line} />;
}

export interface Line {
  from: Point;
  to: Point;
  on: boolean;
}

// useArrow draws the own arrow from the source card to the mouse, at
// once from where the mouse is when the aim starts (20.23, ADR 027
// "Target arrow"). It follows the card as it moves, e.g. while it
// lifts in the fan. A target under the pointer turns it "on" and gets
// data-aim-over, which brightens its glow. While it shows, the system
// pointer hides (html.aim-cursor), a mouse zoom closes on each move,
// and cards zoom only once it rests (arrowShows, useZoom). No arrow on
// touch: its pointer is no mouse.
function useArrow(source: number, targets: OptionView[]): Line | null {
  const [line, setLine] = useState<Line | null>(null);
  const targetsRef = useRef(targets);
  useLayoutEffect(() => {
    targetsRef.current = targets;
  }, [targets]);
  useEffect(() => {
    let at = lastPointer();
    let on = false;
    let over: Element | null = null;
    let drawn: Line | null = null;
    let frame = 0;
    const hover = (p: Point | null) => {
      const target = p ? targetAt(zonesAt(p.x, p.y), targetsRef.current) : undefined;
      const el = target ? targetElement(target) : null;
      if (el !== over) {
        over?.removeAttribute('data-aim-over');
        el?.setAttribute('data-aim-over', '');
        over = el;
      }
      on = target !== undefined;
    };
    // draw runs each frame: the card may move while the pointer rests.
    const draw = () => {
      frame = requestAnimationFrame(draw);
      const r = document.querySelector(`[data-zone="obj-${String(source)}"]`)?.getBoundingClientRect();
      const next = r && at ? { from: { x: r.left + r.width / 2, y: r.top + r.height / 2 }, to: at, on } : null;
      if (!same(next, drawn)) {
        drawn = next;
        setLine(next);
      }
    };
    const move = (e: PointerEvent) => {
      at = e.pointerType === 'mouse' ? { x: e.clientX, y: e.clientY } : null;
      hover(at);
      const { zoom, setZoom } = useApp.getState();
      if (at && zoom && !zoom.touch) {
        setZoom(null);
      }
    };
    hover(at);
    frame = requestAnimationFrame(draw);
    window.addEventListener('pointermove', move);
    return () => {
      cancelAnimationFrame(frame);
      window.removeEventListener('pointermove', move);
      over?.removeAttribute('data-aim-over');
    };
  }, [source]);
  const shows = line !== null;
  useEffect(() => {
    if (!shows) {
      return undefined;
    }
    setArrowShows(true);
    document.documentElement.classList.add('aim-cursor');
    return () => {
      setArrowShows(false);
      document.documentElement.classList.remove('aim-cursor');
    };
  }, [shows]);
  return line;
}

// same: two lines are drawn alike.
function same(a: Line | null, b: Line | null): boolean {
  if (!a || !b) {
    return a === b;
  }
  return a.on === b.on && a.from.x === b.from.x && a.from.y === b.from.y && a.to.x === b.to.x && a.to.y === b.to.y;
}

// AimArrow (15.32, 20.23): the aiming arrow, after the reference's
// (ADR 027 "Target arrow"). A curve bowing up the screen (arc) carries
// a chain of dashes from the card to the head, whose tip is at the
// pointer; each dash and the head turn along the curve. The dashes
// march toward the head (anim: aim-march): each runs one step along
// the curve and fades in at the tail and out by the head, so the loop
// has no seam. On a valid target of the own aim a ring marks the tip.
// Speed "off" and reduced motion: the dashes stand still. Another
// player's arrow (others, 15.36) sits under the windows, with no ring.
export function AimArrow({ from, to, on, others }: Line & { others?: boolean }) {
  const scale = useApp((s) => scaleFor(s.settings.animation, 0, 0));
  const moving = scale > 0 && !still();
  const way = arc(from, to);
  const n = dashCount(chainOf(way));
  const dashesRef = useRef(new Map<number, HTMLDivElement>());
  const animsRef = useRef(new Map<HTMLDivElement, Animation>());
  // startRef is when the march began: every dash keeps its phase.
  const startRef = useRef<number | null>(null);
  useLayoutEffect(() => {
    const anims = animsRef.current;
    const els = Array.from({ length: n }, (_, i) => dashesRef.current.get(i));
    for (const [el, a] of anims) {
      if (!moving || !els.includes(el)) {
        a.cancel();
        anims.delete(el);
      }
    }
    els.forEach((el, i) => {
      if (!el) {
        return;
      }
      if (!moving) {
        const p = dashPose(way, from, i, 0.5);
        el.style.transform = p.transform;
        el.style.opacity = String(p.opacity);
        return;
      }
      const frames = dashFrames(way, from, i).map((f) => ({ ...f }));
      const a = anims.get(el);
      if (a?.effect instanceof KeyframeEffect) {
        a.effect.setKeyframes(frames);
        a.effect.updateTiming({ duration: marchTime(scale) });
        return;
      }
      // anim: aim-march
      const made = el.animate(frames, { duration: marchTime(scale), iterations: Infinity });
      const now = document.timeline.currentTime;
      startRef.current ??= typeof now === 'number' ? now : performance.now();
      made.startTime = startRef.current;
      anims.set(el, made);
    });
  });
  useEffect(() => {
    const anims = animsRef.current;
    return () => {
      for (const a of anims.values()) {
        a.cancel();
      }
      anims.clear();
    };
  }, []);
  const { x, y, angle } = way.end;
  const cls = ['aim-arrow', on ? 'on' : '', moving ? '' : 'still', others ? 'others' : ''].filter(Boolean).join(' ');
  const look = {
    '--aim-dash': rem(feel.aim.dash),
    '--aim-wide': rem(feel.aim.wide),
    '--aim-ring': `${String(feel.aim.ring * scale)}s`,
  } as CSSProperties;
  return (
    <div className={cls} style={look} aria-hidden="true">
      <div className="aim-chain" style={{ transform: `translate(${String(from.x)}px, ${String(from.y)}px)` }}>
        {Array.from({ length: n }, (_, i) => (
          <div
            key={i}
            ref={(el) => {
              if (el) {
                dashesRef.current.set(i, el);
              } else {
                dashesRef.current.delete(i);
              }
            }}
            className="aim-dash"
          >
            <div />
          </div>
        ))}
      </div>
      <div className="aim-head" style={{ transform: `translate(${String(x)}px, ${String(y)}px) rotate(${String(angle)}deg)` }} />
      <div className="aim-ring" style={{ transform: `translate(${String(to.x)}px, ${String(to.y)}px)` }}>
        {!others && <div />}
      </div>
    </div>
  );
}
