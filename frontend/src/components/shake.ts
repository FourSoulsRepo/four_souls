// Shakes (20.3, ADR 027 principle 4): a hit pushes the table cards
// near it away, harder for a bigger hit; a lethal hit or a player's
// death jitters the whole table. They move the cards' and the table's
// own translate, apart from their transform (ADR 002). Animation speed
// "off" (scale 0) and reduced motion shake nothing.
import { feel, type HitSize } from '../feel';
import { px } from '../uiScale';

// calm is true when nothing may shake.
function calm(scale: number): boolean {
  return (
    scale <= 0 ||
    (typeof window.matchMedia === 'function' && window.matchMedia('(prefers-reduced-motion: reduce)').matches)
  );
}

// shakeNear pushes the table cards near a hit at `at` away from it.
export function shakeNear(at: DOMRect, size: HitSize, scale: number): void {
  if (calm(scale)) {
    return;
  }
  const { push, reach, near, shake, shakeCurve } = feel.hit;
  const cx = at.left + at.width / 2;
  const cy = at.top + at.height / 2;
  const radius = near * Math.max(at.width, at.height);
  const amp = (push[size - 1] ?? 0) * px(reach);
  for (const el of document.querySelectorAll<HTMLElement>('.card[data-zone^="obj-"]')) {
    const r = el.getBoundingClientRect();
    const dx = r.left + r.width / 2 - cx;
    const dy = r.top + r.height / 2 - cy;
    const d = Math.hypot(dx, dy);
    if (d < 1 || d > radius || r.width === 0) {
      continue;
    }
    const x = (dx / d) * amp;
    const y = (dy / d) * amp;
    const to = (k: number) => `${String(x * k)}px ${String(y * k)}px`;
    // anim: hit-shake
    el.animate(
      [
        { translate: to(0) },
        { translate: to(1), offset: 0.2 },
        { translate: to(-0.35), offset: 0.55 },
        { translate: to(0.1), offset: 0.8 },
        { translate: to(0) },
      ],
      { duration: shake * 1000 * scale, easing: shakeCurve, id: 'hit-shake' },
    );
  }
}

// shakeTable jitters the whole table, as a camera shake: random
// offsets that fade to nothing.
export function shakeTable(scale: number): void {
  const table = document.querySelector<HTMLElement>('main.game');
  if (!table || calm(scale)) {
    return;
  }
  for (const a of table.getAnimations()) {
    if (a.id === 'table-shake') {
      a.cancel();
    }
  }
  const { reach, time, steps } = feel.table;
  const jolt = () => (Math.random() * 2 - 1) * px(reach);
  const frames: Keyframe[] = [{ translate: '0px 0px' }];
  for (let i = 1; i < steps; i++) {
    const k = 1 - i / steps;
    frames.push({ translate: `${String(jolt() * k)}px ${String(jolt() * k)}px` });
  }
  frames.push({ translate: '0px 0px' });
  // anim: table-shake
  table.animate(frames, { duration: time * 1000 * scale, id: 'table-shake' });
}
