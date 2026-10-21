// The attack lunge (15.3, V-12, T-23). A copy of the attacker's card
// jumps into the target and returns; a starburst "-N" pops on the hit
// (H-14). The copy lives in a layer of its own, so it is never clipped
// by the table and outlives the effect: playback goes on at the hit,
// and the target's HP drops at once (C-03). Only transform and opacity
// animate (ADR 002). The hit lands with an attack hit sound (15.50).
import { playSound } from '../sounds';

// The normal length of the lunge and of the starburst, in ms.
const lungeMs = 700;
const burstMs = 900;
// The hit lands at this part of the lunge.
const hitAt = 0.5;
// How many points the burst's star has.
const starPoints = 14;

export function find(zones: string[]): HTMLElement | null {
  for (const z of zones) {
    const el = document.querySelector<HTMLElement>(`[data-zone="${CSS.escape(z)}"]`);
    if (el) {
      return el;
    }
  }
  return null;
}

// layer is the page's effect layer, made on first use.
export function layer(): HTMLElement {
  let el = document.querySelector<HTMLElement>('.effects-trail');
  if (!el) {
    el = document.createElement('div');
    el.className = 'effects-trail';
    el.setAttribute('aria-hidden', 'true');
    document.body.append(el);
  }
  return el;
}

const svgNS = 'http://www.w3.org/2000/svg';

// starPolygon is the starburst's points in a 100 × 100 box; the strip's
// popout draws the same star (15.90).
export function starPolygon(): string {
  const pts: string[] = [];
  for (let i = 0; i < starPoints * 2; i++) {
    const r = i % 2 === 0 ? 50 : 34;
    const a = (Math.PI * i) / starPoints - Math.PI / 2;
    pts.push(`${(50 + r * Math.cos(a)).toFixed(1)},${(50 + r * Math.sin(a)).toFixed(1)}`);
  }
  return pts.join(' ');
}

function star(): SVGSVGElement {
  const svg = document.createElementNS(svgNS, 'svg');
  svg.setAttribute('viewBox', '0 0 100 100');
  const poly = document.createElementNS(svgNS, 'polygon');
  poly.setAttribute('points', starPolygon());
  svg.append(poly);
  return svg;
}

// copy is a still copy of the card: no zones, names, or test ids, so
// nothing finds it in place of the card.
function copy(src: HTMLElement, r: DOMRect): HTMLElement {
  const c = src.cloneNode(true) as HTMLElement;
  for (const el of [c, ...c.querySelectorAll<HTMLElement>('*')]) {
    el.removeAttribute('data-zone');
    el.removeAttribute('data-testid');
  }
  c.classList.add('effect-lunge');
  c.dataset.testid = 'lunge';
  Object.assign(c.style, {
    left: `${String(r.left)}px`,
    top: `${String(r.top)}px`,
    width: `${String(r.width)}px`,
    height: `${String(r.height)}px`,
    transform: 'none',
    opacity: '1',
  });
  return c;
}

// burst pops a starburst "−N" on a card (H-14, T-24).
export function burst(at: DOMRect, amount: number, scale: number): void {
  const b = document.createElement('div');
  b.className = 'effect-burst';
  b.dataset.testid = 'burst';
  b.style.left = `${String(at.left + at.width / 2)}px`;
  b.style.top = `${String(at.top + at.height / 2)}px`;
  const n = document.createElement('span');
  n.textContent = `−${String(amount)}`;
  b.append(star(), n);
  layer().append(b);
  // anim: burst
  const a = b.animate(
    [
      { transform: 'scale(0.3)', opacity: 0 },
      { transform: 'scale(1.15)', opacity: 1, offset: 0.2 },
      { transform: 'scale(1)', opacity: 1, offset: 0.75 },
      { transform: 'scale(1)', opacity: 0 },
    ],
    { duration: burstMs * scale },
  );
  a.onfinish = () => {
    b.remove();
  };
}

// lunge plays the attack and calls hit when it lands.
export function lunge(from: string[], to: string[], amount: number, scale: number, hit: () => void): void {
  const dst = find(to);
  if (!dst) {
    hit();
    return;
  }
  const b = dst.getBoundingClientRect();
  const src = find(from);
  if (!src) {
    playSound('attackHit');
    burst(b, amount, scale);
    hit();
    return;
  }
  const a = src.getBoundingClientRect();
  // The card goes into the target, stopping short of its centre.
  const dx = b.left + b.width / 2 - (a.left + a.width / 2);
  const dy = b.top + b.height / 2 - (a.top + a.height / 2);
  const len = Math.hypot(dx, dy);
  const k = len > 0 ? Math.max(0, len - Math.min(b.width, b.height) * 0.35) / len : 0;
  const c = copy(src, a);
  layer().append(c);
  src.style.visibility = 'hidden';
  const ms = lungeMs * scale;
  // anim: lunge
  const anim = c.animate(
    [
      { transform: 'translate(0, 0) scale(1)', easing: 'ease-out' },
      { transform: `translate(${String(-dx * 0.08)}px, ${String(-dy * 0.08)}px) scale(1.1)`, offset: 0.3, easing: 'ease-in' },
      { transform: `translate(${String(dx * k)}px, ${String(dy * k)}px) scale(1.05)`, offset: hitAt, easing: 'ease-in-out' },
      { transform: 'translate(0, 0) scale(1)' },
    ],
    { duration: ms },
  );
  const end = () => {
    c.remove();
    src.style.visibility = '';
  };
  anim.onfinish = end;
  anim.oncancel = end;
  window.setTimeout(() => {
    playSound('attackHit');
    burst(b, amount, scale);
    hit();
  }, ms * hitAt);
}
