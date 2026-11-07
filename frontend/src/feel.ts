// Game feel (step 20, ADR 027): every timing and easing curve of the
// table's motion lives here, so tuning one never means a hunt. Times
// are in seconds at normal speed; playback's scale stretches them (the
// animation speed setting, scaleFor). Parts of a motion are fractions
// of its time, as keyframe offsets are.

// curves are the easing curves as CSS strings: WAAPI keyframes and CSS
// take them as they are, motion takes bezier(curve).
export const curves = {
  linear: 'linear',
  easeIn: 'ease-in',
  easeOut: 'ease-out',
  easeInOut: 'ease-in-out',
  // expo-out (1 − 2^(−10t)): a fast start and a long, calm settle.
  // Cards settle with it after a hit or a move (ADR 027).
  expoOut: 'cubic-bezier(0.16, 1, 0.3, 1)',
  // sine-in-expo-out: the reference's draw path (ADR 027 "Draw"). Its
  // switch to expo-out compares the time with the start value, which
  // is 0 on a path, so it runs as expo-out from the first frame.
  sineInExpoOut: 'cubic-bezier(0.16, 1, 0.3, 1)',
} as const;

export type Curve = keyof typeof curves;

// keywords are CSS's named curves as cubic Béziers.
const keywords = new Map<string, [number, number, number, number]>([
  ['linear', [0, 0, 1, 1]],
  ['ease', [0.25, 0.1, 0.25, 1]],
  ['ease-in', [0.42, 0, 1, 1]],
  ['ease-out', [0, 0, 0.58, 1]],
  ['ease-in-out', [0.42, 0, 0.58, 1]],
]);

// bezier is a CSS curve as motion's [x1, y1, x2, y2]. It throws on a
// string CSS would not take.
export function bezier(css: string): [number, number, number, number] {
  const named = keywords.get(css);
  if (named) {
    return [...named];
  }
  const m = /^cubic-bezier\(([^)]*)\)$/.exec(css);
  const n = m?.[1]?.split(',').map((s) => (s.trim() === '' ? NaN : Number(s))) ?? [];
  const [x1 = NaN, y1 = NaN, x2 = NaN, y2 = NaN] = n;
  if (n.length !== 4 || !n.every(Number.isFinite) || x1 < 0 || x1 > 1 || x2 < 0 || x2 > 1) {
    throw new Error(`bad easing curve: ${css}`);
  }
  return [x1, y1, x2, y2];
}

// elastic is easeOutElastic: it overshoots, swings back, and settles
// at 1 (period 0.3, as in ADR 027 "Damage and death").
export function elastic(t: number): number {
  if (t <= 0) {
    return 0;
  }
  if (t >= 1) {
    return 1;
  }
  const p = 0.3;
  return Math.pow(2, -10 * t) * Math.sin(((t - p / 4) * 2 * Math.PI) / p) + 1;
}

export interface Stop {
  offset: number;
  value: number;
  easing: string;
}

// elasticStops is elastic-out as a keyframe list, since CSS has no
// elastic curve: n straight steps, value 0 → 1 with the overshoot.
export function elasticStops(n: number): Stop[] {
  const stops: Stop[] = [];
  for (let i = 0; i <= n; i++) {
    const t = i / n;
    stops.push({ offset: t, value: elastic(t), easing: curves.linear });
  }
  return stops;
}

export const elasticOut = elasticStops(30);

// feel holds the times of each motion. Today's are kept as they were
// (20.1); later step 20 sub-steps tune them to ADR 027.
export const feel = {
  // lunge: the attack (lunge.ts, ADR 027 "Attack"), in three phases,
  // each in seconds. The attacker lifts and grows (windup); attacker
  // and target meet half way (strike); the hit lands on contact; both
  // settle back home (settle). lift is the attacker's size at the top
  // of the windup; pull is how far it draws back first, as a part of
  // the way to the target.
  lunge: {
    windup: 0.25,
    strike: 0.12,
    settle: 0.5,
    lift: 1.1,
    pull: 0.06,
    windupCurve: curves.easeOut,
    strikeCurve: curves.linear,
    settleCurve: curves.expoOut,
  },
  // burst: the starburst "−N", the damage number (lunge.ts, 20.4, ADR
  // 027 "Damage and death"). It grows from 0 over grow, elastic-out
  // (elasticOut); its fade starts hold seconds after it shows and takes
  // fade, as the reference's: 3 s on screen in all. Playback never
  // waits for it. Bursts on one card stack, stack burst heights apart.
  burst: {
    grow: 1,
    hold: 2,
    fade: 1,
    stack: 1,
  },
  // fly: a card from one zone to another (Effects.tsx); fromMiddle
  // from a held reveal, shrinking to a table card.
  fly: {
    time: 0.45,
    fromMiddle: 0.55,
    curve: curves.easeInOut,
  },
  // hit: a hit's size by its damage (20.3, ADR 027 principle 4).
  // Each list is by size, 1 to 3 (hitSize). burst is the starburst's
  // scale; push is how far the cards near the hit move away from it,
  // as a part of reach (px), as the reference's shakes 0.1, 0.5 and 1;
  // near is how close a card must be, its centre to the hit's, in the
  // hit card's larger side. The cards push out, swing back past home
  // a little and stop, over shake seconds.
  hit: {
    burst: [1, 1.25, 1.5],
    push: [0.1, 0.5, 1],
    reach: 12,
    near: 2.5,
    shake: 0.35,
    shakeCurve: curves.easeOut,
  },
  // table: the whole table jitters on a lethal hit or a player's death
  // (20.3): random offsets up to reach px that fade to 0 over time, in
  // steps jolts.
  table: {
    reach: 6,
    time: 0.4,
    steps: 8,
  },
  // banner: the text in the middle (Effects.tsx). times are its
  // opacity keyframes: in, held, out. A long one (start of game) stays
  // twice as long, with fades as long.
  banner: {
    time: 1.1,
    long: 2.2,
    times: [0, 0.2, 0.8, 1],
    longTimes: [0, 0.1, 0.9, 1],
  },
} as const;

export type HitSize = 1 | 2 | 3;

// hitSize is a hit's size by its damage (20.3): 1 is small, 2–3 mid,
// 4 and more large.
export function hitSize(amount: number): HitSize {
  if (amount >= 4) {
    return 3;
  }
  if (amount >= 2) {
    return 2;
  }
  return 1;
}
