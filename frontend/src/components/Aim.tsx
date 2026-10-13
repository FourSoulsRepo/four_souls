import { useEffect, useState } from 'react';
import { arrow, targetAt, zonesAt, type Point } from '../aim';
import { useKey } from '../keys';
import { useApp, type Aiming } from '../store';

// AimLayer is the player's aim (15.25): a click on a card that needs a
// target starts it, the next click picks the target. A click on
// anything else, a right click, or Esc cancels it. While the mouse
// moves, an arrow runs from the card to the pointer.
export function AimLayer() {
  const aiming = useApp((s) => s.aiming);
  return aiming ? <Aim key={aiming.source} aiming={aiming} /> : null;
}

function Aim({ aiming }: { aiming: Aiming }) {
  const setAiming = useApp((s) => s.setAiming);
  const playAt = useApp((s) => s.playAt);
  const [line, setLine] = useState<{ from: Point; to: Point } | null>(null);
  useKey('escape', () => {
    setAiming(null);
  });
  useEffect(() => {
    const move = (e: PointerEvent) => {
      // The arrow starts at the card's middle, wherever it is drawn now.
      const r = document.querySelector(`[data-zone="obj-${String(aiming.source)}"]`)?.getBoundingClientRect();
      if (r && e.pointerType === 'mouse') {
        setLine({ from: { x: r.left + r.width / 2, y: r.top + r.height / 2 }, to: { x: e.clientX, y: e.clientY } });
      }
    };
    // Capture: the click picks or cancels and reaches nothing else.
    const click = (e: MouseEvent) => {
      e.stopPropagation();
      e.preventDefault();
      const target = targetAt(zonesAt(e.clientX, e.clientY), aiming.targets);
      if (target) {
        playAt(aiming.intent, target);
      } else {
        setAiming(null);
      }
    };
    const menu = (e: MouseEvent) => {
      e.preventDefault();
      setAiming(null);
    };
    window.addEventListener('pointermove', move);
    window.addEventListener('click', click, true);
    window.addEventListener('contextmenu', menu, true);
    return () => {
      window.removeEventListener('pointermove', move);
      window.removeEventListener('click', click, true);
      window.removeEventListener('contextmenu', menu, true);
    };
  }, [aiming, playAt, setAiming]);
  return line ? <AimArrow from={line.from} to={line.to} /> : null;
}

// The head's length and the shaft's drawn width, as in Table.css.
const headLength = 30;
const shaftWidth = 100;

// AimArrow: a shaft from the card to the head's base, and a head whose
// tip is at the pointer; placed with transform only. 15.32 gives it its
// final look.
function AimArrow({ from, to }: { from: Point; to: Point }) {
  const { length, angle } = arrow(from, to);
  const at = (p: Point) => `translate(${String(p.x)}px, ${String(p.y)}px) rotate(${String(angle)}deg)`;
  const shaft = Math.max(0, length - headLength + 2) / shaftWidth;
  return (
    <div className="aim-arrow" aria-hidden="true">
      <div className="aim-shaft" style={{ transform: `${at(from)} scaleX(${String(shaft)})` }} />
      <div className="aim-head" style={{ transform: at(to) }} />
    </div>
  );
}
