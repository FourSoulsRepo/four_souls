// The attack lunge (15.3, V-12, T-23; fast strike 20.2, ADR 027). A
// copy of the attacker's card lifts, then it and a copy of the target
// meet half way and settle back; a starburst "-N" pops on the hit
// (H-14). The copies live in a layer of their own, so they are never
// clipped by the table and outlive the effect: playback goes on at the
// hit, and the target's HP shows on its card once its copy is home
// (C-03). Only transform and opacity animate (ADR 002). The hit lands
// with an attack hit sound (15.50). A bigger hit looks and sounds
// bigger (20.3): its burst grows and glows, the cards near it shake,
// and a lethal one shakes the table.
import { elasticOut, feel, hitSize, type HitSize } from '../feel';
import { playSound, type Sound } from '../sounds';
import { shakeNear, shakeTable } from './shake';

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
// nothing finds it in place of the card. testid names the copy.
function copy(src: HTMLElement, r: DOMRect, testid: string): HTMLElement {
  const c = src.cloneNode(true) as HTMLElement;
  for (const el of [c, ...c.querySelectorAll<HTMLElement>('*')]) {
    el.removeAttribute('data-zone');
    el.removeAttribute('data-testid');
  }
  c.classList.add('effect-lunge');
  c.dataset.testid = testid;
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

// hitSounds are the attack hit sounds by hit size (20.3, 20.17).
const hitSounds: Record<HitSize, Sound> = { 1: 'attackHit', 2: 'attackHit2', 3: 'attackHit3' };

// hitSound is the sound of an attack that deals amount.
export function hitSound(amount: number): Sound {
  return hitSounds[hitSize(amount)];
}

// burstSide is the burst's side in px at size 1 (.effect-burst).
const burstSide = 76;

// stackSpot is a new burst's spot among the live ones on its card
// (20.4): 0 in the middle, then −1 above, 1 below, −2, 2, … The first
// free spot is taken, so a spot frees up when its burst is gone.
export function stackSpot(taken: number[]): number {
  for (let i = 0; ; i++) {
    const spot = i % 2 === 1 ? -(i + 1) / 2 : i / 2;
    if (!taken.includes(spot)) {
      return spot;
    }
  }
}

// live are the bursts on screen: the centre of the card each sits on,
// and its spot in that card's stack.
const live = new Set<{ x: number; y: number; spot: number }>();

// burstFrames are the burst's keyframes (ADR 027 "Damage and death"):
// it grows from 0 with an elastic overshoot, holds, and fades. With
// less motion it shows at full size at once and only fades.
function burstFrames(calm: boolean): Keyframe[] {
  const { grow, hold, fade } = feel.burst;
  const total = hold + fade;
  const frames: Keyframe[] = calm
    ? [{ transform: 'scale(1)', opacity: 1, offset: 0 }]
    : elasticOut.map((s) => ({ transform: `scale(${String(s.value)})`, opacity: 1, offset: (s.offset * grow) / total, easing: s.easing }));
  frames.push({ transform: 'scale(1)', opacity: 1, offset: hold / total }, { transform: 'scale(1)', opacity: 0, offset: 1 });
  return frames;
}

// burst pops a starburst "−N" on a card (H-14, T-24), sized by the
// damage, and shakes the cards near it (20.3). A lethal hit shakes the
// table too. It stays 3 s (20.4), but the caller goes on at once:
// playback never waits for it. Bursts on one card stack, not overlap.
export function burst(at: DOMRect, amount: number, scale: number, lethal = false): void {
  const size = hitSize(amount);
  const grows = feel.hit.burst[size - 1] ?? 1;
  const x = at.left + at.width / 2;
  const y = at.top + at.height / 2;
  const taken = [...live].filter((l) => Math.abs(l.x - x) < at.width / 2 && Math.abs(l.y - y) < at.height / 2).map((l) => l.spot);
  const spot = { x, y, spot: stackSpot(taken) };
  live.add(spot);
  const b = document.createElement('div');
  b.className = `effect-burst size-${String(size)}`;
  b.dataset.testid = 'burst';
  b.dataset.size = String(size);
  b.style.scale = String(grows);
  b.style.left = `${String(x)}px`;
  b.style.top = `${String(y + spot.spot * burstSide * grows * feel.burst.stack)}px`;
  const n = document.createElement('span');
  n.textContent = `−${String(amount)}`;
  b.append(star(), n);
  layer().append(b);
  const { hold, fade } = feel.burst;
  // anim: burst
  const a = b.animate(burstFrames(still()), { duration: (hold + fade) * 1000 * scale, id: 'burst' });
  const end = () => {
    live.delete(spot);
    b.remove();
  };
  a.onfinish = end;
  a.oncancel = end;
  shakeNear(at, size, scale);
  if (lethal) {
    shakeTable(scale);
  }
}

// held maps a card to the lunge whose copy of it still settles home,
// and ends that lunge at once. A next lunge on the card (the monster
// hits back) ends it first: one copy of a card at a time.
const held = new Map<HTMLElement, () => void>();

// still is true when the player asks for less motion: the hit lands
// in place, with no strike.
function still(): boolean {
  return typeof window.matchMedia === 'function' && window.matchMedia('(prefers-reduced-motion: reduce)').matches;
}

// lunge plays the attack (ADR 027 "Attack") and calls hit on contact.
// Windup: the attacker's copy lifts and grows. Strike: it and the
// target's copy meet half way, linear. Contact: the sound, the burst
// and hit, from the strike's own end, not from a timer. Settle: both
// ease back home. A target that is no card (a player's panel) stays;
// the attacker then goes all the way into it. lethal: the hit kills.
export function lunge(from: string[], to: string[], amount: number, scale: number, hit: () => void, lethal = false): void {
  const dst = find(to);
  if (!dst) {
    hit();
    return;
  }
  held.get(dst)?.();
  const src = find(from);
  if (src) {
    held.get(src)?.();
  }
  const b = dst.getBoundingClientRect();
  if (!src || src === dst || still()) {
    playSound(hitSound(amount));
    burst(b, amount, scale, lethal);
    hit();
    return;
  }
  const a = src.getBoundingClientRect();
  const dx = b.left + b.width / 2 - (a.left + a.width / 2);
  const dy = b.top + b.height / 2 - (a.top + a.height / 2);
  const len = Math.hypot(dx, dy);
  // They stop with their centres a little apart, so the cards overlap
  // at the hit, not cover each other. A card target meets the
  // attacker half way.
  const way = len > 0 ? Math.max(0, len - Math.min(b.width, b.height) * 0.35) / len : 0;
  const meets = dst.classList.contains('card');
  const share = meets ? 0.5 : 1;
  const ax = dx * way * share;
  const ay = dy * way * share;
  const tx = -dx * way * (1 - share);
  const ty = -dy * way * (1 - share);
  const { windup, strike, settle, lift, pull, windupCurve, strikeCurve, settleCurve } = feel.lunge;
  const ms = (s: number) => s * 1000 * scale;
  const go = windup + strike;
  const at = (x: number, y: number, k: number) => `translate(${String(x)}px, ${String(y)}px) scale(${String(k)})`;
  const home = at(0, 0, 1);
  const top = at(-dx * pull, -dy * pull, lift);
  const ahit = at(ax, ay, (1 + lift) / 2);
  const thit = at(tx, ty, 1);
  const target = meets ? copy(dst, b, 'lunge-target') : null;
  const c = copy(src, a, 'lunge');
  if (target) {
    layer().append(target);
    dst.style.visibility = 'hidden';
  }
  layer().append(c);
  src.style.visibility = 'hidden';
  // anim: lunge
  const swing = c.animate(
    [
      { transform: home, easing: windupCurve },
      { transform: top, offset: windup / go, easing: strikeCurve },
      { transform: ahit },
    ],
    { duration: ms(go), id: 'strike', fill: 'forwards' },
  );
  const strikes = [swing];
  if (target) {
    strikes.push(
      // anim: lunge
      target.animate(
        [
          { transform: home },
          { transform: home, offset: windup / go, easing: strikeCurve },
          { transform: thit },
        ],
        { duration: ms(go), id: 'strike', fill: 'forwards' },
      ),
    );
  }
  const settles: Animation[] = [];
  let touched = false;
  const end = () => {
    // Ended before the hit (never expected): the hit still lands, so
    // playback never stalls.
    if (!touched) {
      touched = true;
      hit();
    }
    for (const x of [...strikes, ...settles]) {
      x.onfinish = null;
      x.oncancel = null;
    }
    c.remove();
    target?.remove();
    src.style.visibility = '';
    if (target) {
      dst.style.visibility = '';
    }
    for (const el of [src, dst]) {
      if (held.get(el) === end) {
        held.delete(el);
      }
    }
  };
  held.set(src, end);
  if (target) {
    held.set(dst, end);
  }
  const contact = () => {
    if (touched) {
      return;
    }
    touched = true;
    playSound(hitSound(amount));
    burst(b, amount, scale, lethal);
    const back = (el: HTMLElement, from: string) =>
      // anim: lunge
      el.animate([{ transform: from }, { transform: home }], { duration: ms(settle), easing: settleCurve, id: 'settle' });
    const retreat = back(c, ahit);
    settles.push(retreat);
    if (target) {
      settles.push(back(target, thit));
    }
    // Both settle in the same time: the attacker's end ends both.
    retreat.onfinish = end;
    retreat.oncancel = end;
    hit();
  };
  swing.onfinish = contact;
  swing.oncancel = contact;
}
