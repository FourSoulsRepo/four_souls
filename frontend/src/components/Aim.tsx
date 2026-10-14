import { useEffect, useState, type CSSProperties } from 'react';
import { aimedOption } from '../actions';
import { aimClick, arrow, blocks, cancelOption, targetAt, zonesAt, type Point } from '../aim';
import { useKey } from '../keys';
import { scaleFor } from '../playback';
import type { Intent, OptionView } from '../protocol';
import { useApp, type Aiming } from '../store';

// AimLayer is the player's aim (15.25): a click on a card that needs a
// target starts it, the next click picks the target. A click on empty
// space, a right click, or Esc cancels it. A click on the card again
// does nothing, unless it may target itself. A click on a card, pile,
// or stack row that is no target keeps aiming and says why (15.33). While the mouse moves, an arrow
// runs from the card to the pointer. A card already placed answers its
// target prompt; cancelling it sends the prompt's cancel, and the card
// goes back (15.34, T-13).
export function AimLayer() {
  const aiming = useApp((s) => s.aiming);
  return aiming ? <Aim key={aiming.source} aiming={aiming} /> : null;
}

function Aim({ aiming }: { aiming: Aiming }) {
  const setAiming = useApp((s) => s.setAiming);
  const playAt = useApp((s) => s.playAt);
  const act = useApp((s) => s.act);
  const [line, setLine] = useState<{ from: Point; to: Point; on: boolean } | null>(null);
  // why is the reason shown beside a wrong click; n restarts its fade.
  const [why, setWhy] = useState<{ at: Point; n: number } | null>(null);
  useKey('escape', () => {
    cancel(aiming, act, setAiming);
  });
  useEffect(() => {
    const move = (e: PointerEvent) => {
      // The arrow starts at the card's middle, wherever it is drawn now.
      const r = document.querySelector(`[data-zone="obj-${String(aiming.source)}"]`)?.getBoundingClientRect();
      if (r && e.pointerType === 'mouse') {
        const on = !!targetAt(zonesAt(e.clientX, e.clientY), aiming.targets);
        setLine({ from: { x: r.left + r.width / 2, y: r.top + r.height / 2 }, to: { x: e.clientX, y: e.clientY }, on });
      }
    };
    // Capture: the click picks or cancels and reaches nothing else.
    const click = (e: MouseEvent) => {
      e.stopPropagation();
      e.preventDefault();
      const what = aimClick(zonesAt(e.clientX, e.clientY), aiming.source, aiming.targets);
      if (what.kind === 'pick') {
        pick(aiming, what.target);
      } else if (what.kind === 'wrong') {
        setWhy((w) => ({ at: { x: e.clientX, y: e.clientY }, n: (w?.n ?? 0) + 1 }));
      } else if (what.kind === 'cancel') {
        cancel(aiming, act, setAiming);
      }
    };
    const pick = (a: Aiming, target: OptionView) => {
      if (a.intent) {
        playAt(a.intent, target);
        return;
      }
      const answer = a.choice && aimedOption(a.choice, target);
      setAiming(null);
      if (answer) {
        act(answer);
      }
    };
    const menu = (e: MouseEvent) => {
      e.preventDefault();
      cancel(aiming, act, setAiming);
    };
    window.addEventListener('pointermove', move);
    window.addEventListener('click', click, true);
    window.addEventListener('contextmenu', menu, true);
    return () => {
      window.removeEventListener('pointermove', move);
      window.removeEventListener('click', click, true);
      window.removeEventListener('contextmenu', menu, true);
    };
  }, [aiming, act, playAt, setAiming]);
  // The reason goes after a while; a new wrong click shows it again.
  useEffect(() => {
    if (!why) {
      return;
    }
    const t = setTimeout(() => {
      setWhy(null);
    }, whyMs);
    return () => {
      clearTimeout(t);
    };
  }, [why]);
  return (
    <>
      {line && <AimArrow {...line} />}
      {why && (
        <div
          key={why.n}
          className="aim-why"
          role="status"
          style={{ transform: `translate(${String(why.at.x)}px, ${String(why.at.y)}px)` }}
        >
          <div>{aiming.reason ?? 'Not a valid target.'}</div>
        </div>
      )}
    </>
  );
}

// cancel ends an aim; a placed card's prompt is cancelled with it.
function cancel(aiming: Aiming, act: (i: Intent) => void, setAiming: (a: Aiming | null) => void) {
  const c = aiming.choice && cancelOption(aiming.choice);
  setAiming(null);
  if (c) {
    act(c);
  }
}

// How long a wrong click's reason stays, in milliseconds.
const whyMs = 2500;

// The head's length and the chain's step, as in Table.css.
const headLength = 48;
const step = 36;

// AimArrow (15.32): a chain of thick blocks from the card to the head,
// and a head whose tip is at the pointer; on a valid target a ring
// marks the tip. The blocks march toward the head at the animation
// speed; the first fades in and the last fades out, so the loop has no
// seam. Placed and moved with transform and opacity only.
function AimArrow({ from, to, on }: { from: Point; to: Point; on: boolean }) {
  const scale = useApp((s) => scaleFor(s.settings.animation, 0, 0));
  const { length, angle } = arrow(from, to);
  const at = (p: Point) => `translate(${String(p.x)}px, ${String(p.y)}px) rotate(${String(angle)}deg)`;
  const n = blocks(length, step, headLength);
  const cls = ['aim-arrow', on ? 'on' : '', scale > 0 ? '' : 'still'].filter(Boolean).join(' ');
  const march = { '--march': `${String(0.5 * scale)}s` } as CSSProperties;
  return (
    <div className={cls} style={march} aria-hidden="true">
      <div className="aim-chain" style={{ transform: at(from) }}>
        <div className="aim-march">
          {Array.from({ length: n }, (_, i) => (
            <div key={i} className="aim-block" style={{ transform: `translateX(${String(i * step)}px)` }}>
              <div className={i === 0 ? 'first' : i === n - 1 ? 'last' : undefined} />
            </div>
          ))}
        </div>
      </div>
      <div className="aim-head" style={{ transform: at(to) }} />
      <div className="aim-ring" style={{ transform: `translate(${String(to.x)}px, ${String(to.y)}px)` }}>
        <div />
      </div>
    </div>
  );
}
