// Hit and death marks (15.37). Damage that is not an attack flies from
// its source as a glowing shot (T-15); every damaged card then bursts
// at once (T-24). A dying card gets a big skull before it leaves (T-26).
// Both play in the lunge's layer and end the effect early, so the table
// changes while the marks fade. Only transform and opacity animate
// (ADR 002). A shot sounds as a fireball: cast, then impact (15.50).
import type { Hit } from '../playback';
import { playSound } from '../sounds';
import { burst, find, layer } from './lunge';

// The normal length of a shot and of a skull, in ms.
const shotMs = 450;
const skullMs = 1100;
// The dying card leaves at this part of the skull.
const leaveAt = 0.6;

// skullPath is the stack's skull icon (24 × 24); the big skull reuses it.
export const skullPath =
  'M12 2.5a8 8 0 0 0-8 8c0 2.8 1.4 4.6 3 5.6V19a1.5 1.5 0 0 0 1.5 1.5h7A1.5 1.5 0 0 0 17 19v-2.9c1.6-1 3-2.8 3-5.6a8 8 0 0 0-8-8Zm-3.2 7.2a1.9 1.9 0 1 1 0 3.8 1.9 1.9 0 0 1 0-3.8Zm6.4 0a1.9 1.9 0 1 1 0 3.8 1.9 1.9 0 0 1 0-3.8ZM12 14l1.3 2.2h-2.6Z';

const svgNS = 'http://www.w3.org/2000/svg';

function centre(r: DOMRect): { x: number; y: number } {
  return { x: r.left + r.width / 2, y: r.top + r.height / 2 };
}

// place puts a mark's centre on a point of the page.
function place(el: HTMLElement, x: number, y: number): void {
  el.style.left = `${String(x)}px`;
  el.style.top = `${String(y)}px`;
}

// hits flies a shot from each source to its target, then bursts on
// every target at once and calls done.
export function hits(list: Hit[], scale: number, done: () => void): void {
  const landed: { at: DOMRect; amount: number }[] = [];
  let flying = false;
  for (const h of list) {
    const dst = find(h.to);
    if (!dst) {
      continue;
    }
    const at = dst.getBoundingClientRect();
    landed.push({ at, amount: h.amount });
    const src = find(h.from);
    if (!src) {
      continue;
    }
    flying = true;
    const a = centre(src.getBoundingClientRect());
    const b = centre(at);
    const shot = document.createElement('div');
    shot.className = 'effect-shot';
    shot.dataset.testid = 'shot';
    place(shot, a.x, a.y);
    layer().append(shot);
    const dx = b.x - a.x;
    const dy = b.y - a.y;
    const anim = shot.animate(
      [
        { transform: 'translate(0, 0) scale(0.4)', opacity: 0 },
        { transform: `translate(${String(dx * 0.1)}px, ${String(dy * 0.1)}px) scale(1)`, opacity: 1, offset: 0.15 },
        { transform: `translate(${String(dx)}px, ${String(dy)}px) scale(1.2)`, opacity: 1 },
      ],
      { duration: shotMs * scale, easing: 'ease-in' },
    );
    const end = () => {
      shot.remove();
    };
    anim.onfinish = end;
    anim.oncancel = end;
  }
  if (flying) {
    playSound('effectCast');
  }
  const land = () => {
    if (landed.length > 0) {
      playSound('effectImpact');
    }
    for (const l of landed) {
      burst(l.at, l.amount, scale);
    }
    done();
  };
  if (flying) {
    window.setTimeout(land, shotMs * scale);
  } else {
    land();
  }
}

function skull(): SVGSVGElement {
  const svg = document.createElementNS(svgNS, 'svg');
  svg.setAttribute('viewBox', '0 0 24 24');
  const path = document.createElementNS(svgNS, 'path');
  path.setAttribute('d', skullPath);
  svg.append(path);
  return svg;
}

// skulls pops a skull on each dying card; done is called when the cards
// should leave, and the skulls fade after.
export function skulls(at: string[][], scale: number, done: () => void): void {
  let shown = false;
  for (const zones of at) {
    const el = find(zones);
    if (!el) {
      continue;
    }
    shown = true;
    const r = el.getBoundingClientRect();
    const p = centre(r);
    const s = document.createElement('div');
    s.className = 'effect-skull';
    s.dataset.testid = 'skull';
    s.style.width = `${String(Math.min(r.width, r.height) * 0.9)}px`;
    place(s, p.x, p.y);
    s.append(skull());
    layer().append(s);
    const anim = s.animate(
      [
        { transform: 'scale(0.3) rotate(-12deg)', opacity: 0 },
        { transform: 'scale(1.15) rotate(4deg)', opacity: 1, offset: 0.2 },
        { transform: 'scale(1) rotate(0deg)', opacity: 1, offset: leaveAt },
        { transform: 'scale(1.1) rotate(0deg)', opacity: 0 },
      ],
      { duration: skullMs * scale },
    );
    const end = () => {
      s.remove();
    };
    anim.onfinish = end;
    anim.oncancel = end;
  }
  if (shown) {
    window.setTimeout(done, skullMs * scale * leaveAt);
  } else {
    done();
  }
}
