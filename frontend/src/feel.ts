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
  // expo-in (2^(10t − 10)): still at first, then a rush at the end. A
  // drawn card leaves its show spot for the hand with it (20.7).
  expoIn: 'cubic-bezier(0.7, 0, 0.84, 0)',
  // quart-in-out: another player's drawn card, deck to hand (20.7).
  quartInOut: 'cubic-bezier(0.76, 0, 0.24, 1)',
  // cubic-in-out: the board dims under a hovered history strip (20.21).
  cubicInOut: 'cubic-bezier(0.65, 0, 0.35, 1)',
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
  // each in seconds. The attacker lifts and grows (windup); it flies
  // into the target, which stays home (strike, 20.19); the hit lands on
  // contact; it settles back home (settle). lift is the attacker's size
  // at the top of the windup; pull is how far it draws back first, as a
  // part of the way to the target; short is how far its centre stops
  // short of a card target's, in its own length along the way (the
  // reference's 0.5 − (−0.25)): cards of one size overlap by a quarter.
  lunge: {
    windup: 0.25,
    strike: 0.12,
    settle: 0.5,
    lift: 1.1,
    pull: 0.06,
    short: 0.75,
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
  // hand: the own hand's fan (20.5, ADR 027 "Hand", fanPose). A hand
  // of fewer than from cards lies straight. From then on the fan spans
  // angle plus perCard per card, in degrees, spread evenly, and the
  // cards drop on the reference's parabola; drop is its unit in px.
  hand: {
    from: 4,
    angle: 40,
    perCard: 2,
    drop: 128,
  },
  // zoom: a hovered card shown large (20.5, ADR 027 principle 2). A
  // hand card's shows at once, with no tween. A table card's waits
  // delay, then grows over grow from the card's size. While one shows,
  // or within grace after one closed, the next shows at once: the
  // pointer moves from card to card and the zoom stays open.
  // While the own arrow aims (20.23), no card zooms on hover: one zooms
  // once the arrow rests still on it for rest, and any move closes it.
  zoom: {
    delay: 0.4,
    grow: 0.15,
    growCurve: curves.expoOut,
    grace: 0.3,
    rest: 0.8,
  },
  // aim: the aiming arrow (20.23, ADR 027 "Target arrow", aim.ts,
  // Aim.tsx). It bows up the screen: its middle stands lift px off the
  // straight line, or liftShare of the line's length on a short one.
  // Only the lift across the line bends it, at full from a line lean of
  // the way off upright (|dx| / length × lean ≥ 1) and less below, so a
  // line straight up stays straight and the bow never flips. Dashes
  // dash px long and wide px across, one every step px, march toward
  // the head at march px a second (the animation speed scales it);
  // stops keyframes sample one dash's way. The dash by the head fades
  // out over its last step as (left / step)^headFade, the tail's fades
  // in as (way / step)^tailFade. The chain ends head px short of the
  // tip, under the head's back, so the last dashes fade into it. The
  // ring on a hovered target grows in over ring. Another player's arrow
  // stays linger seconds after they stop aiming.
  aim: {
    lift: 90,
    liftShare: 0.25,
    lean: 3,
    dash: 40,
    wide: 22,
    step: 48,
    march: 30,
    stops: 8,
    headFade: 6,
    tailFade: 2,
    head: 24,
    ring: 0.25,
    linger: 0.25,
  },
  // drag: a held hand card (20.6, ADR 027 "Drag", drag.ts). Each move
  // tilts it gain degrees per unit moved, up to pitch and roll; unit is
  // the reference's table unit in px, ours, tuned by eye. It springs
  // back to flat over rest, a smooth-damp. perspective is the depth of
  // its 3D turn, in px. Two loops follow it, as mixer gains: the bed
  // fades in to bed over bedIn; the air's level is the log of the speed
  // per 60 fps frame, airScale × units + airBase, eased over airEase at
  // most airRate a second, and plays at 0 to air. Both fade out over
  // out on drop.
  drag: {
    gain: 25,
    unit: 60,
    pitch: 40,
    roll: 45,
    rest: 0.2,
    perspective: 600,
    bed: 0.15,
    bedIn: 0.5,
    air: 0.5,
    airScale: 0.5,
    airBase: 0.92,
    airEase: 0.04,
    airRate: 1,
    out: 0.2,
  },
  // draw: a loot card drawn (20.7, ADR 027 "Draw", draw.ts). An own
  // draw rises from the deck toward a point lift px above it, flies to
  // the show spot over show (spotX and spotY of the window, spot px
  // wide), then into the hand over hand: 1.5 s in all. On the way it
  // turns face up over turn, after turnAt, and grows to the show size
  // over grow, after growAt. Its hand card is playable ready seconds
  // after it lands. Another player's draw flies deck to hand over
  // other, back only. From the fastFrom-th draw of a player in one
  // update on (counted from 0), each is faster by faster, down to
  // fastest of its time (drawPace).
  draw: {
    show: 1,
    hand: 0.5,
    turnAt: 0.1,
    turn: 0.9,
    growAt: 0.1,
    grow: 0.5,
    lift: 80,
    spotX: 0.5,
    spotY: 0.45,
    spot: 200,
    ready: 0.3,
    other: 0.8,
    fastFrom: 3,
    faster: 0.111,
    fastest: 0.556,
    showCurve: curves.sineInExpoOut,
    handCurve: curves.expoIn,
    turnCurve: curves.expoOut,
    growCurve: curves.expoOut,
    otherCurve: curves.quartInOut,
  },
  // turn: a new turn's banner (20.8, ADR 027 "Turn start", turn.ts). It
  // fades in over grow while it grows from from to full size, then
  // punches down to punch over punchTime. Its fade starts hide seconds
  // after it shows and takes out, shrinking back to from. Playback goes
  // on after wait: the turn's draw starts while it still shows. The
  // recharge readies the player's cards as it shows.
  turn: {
    grow: 0.25,
    from: 0.1,
    punch: 0.98,
    punchTime: 0.15,
    hide: 1.5,
    out: 0.25,
    wait: 1,
    curve: curves.expoOut,
  },
  // played: another player's card shown large (20.9, ADR 027 "Others'
  // plays", BigPlayed in Effects.tsx). It rises from their hand to the
  // middle over rise, growing from the hand card's size; playback goes
  // on wait seconds after it shows, while it stays. Its leave into the
  // stack starts stay seconds after it shows and takes leave; its label
  // fades over leave just before. A newer big card, or a later effect
  // in the middle (playedEnds), sends it away at once.
  played: {
    rise: 1,
    wait: 1,
    stay: 3,
    leave: 0.4,
    riseCurve: curves.expoOut,
    leaveCurve: curves.easeInOut,
  },
  // endTurn: the End turn button (20.10, ADR 027 "End Turn button",
  // TurnActions.tsx). Once nothing is left (endTurnLook "done"), the
  // cue plays cue seconds later, once per turn, if nothing became
  // possible meanwhile. Its green halo breathes from low to full
  // opacity and back, each way over breath; on hover it brightens to
  // hover over hoverTime, and back. The cue is a wait, not a motion:
  // the animation speed does not scale it.
  endTurn: {
    cue: 1.5,
    breath: 1.6,
    low: 0.55,
    hover: 1.3,
    hoverTime: 0.15,
    hoverCurve: curves.linear,
  },
  // count: a counter that counts (20.11, ADR 027 "Counters", Count.tsx):
  // cents, souls and health. A change of N plays as N steps, step
  // seconds apart, the first at once; more than fit / step steps speed
  // up to fit fit seconds in all (countTimes). Each step flips the
  // number in over flip with flipCurve, by turn degrees.
  count: {
    step: 0.2,
    fit: 1,
    flip: 0.1,
    turn: 90,
    flipCurve: curves.expoOut,
  },
  // catchUp: a lagging client catches up by weight (20.12, ADR 027
  // principle 7, catchUp in playback.ts). Thresholds count the updates
  // still waiting. From dropMinor on, minor effects are dropped; from
  // shorten on, or with burst or more effects in one update, the rest
  // play at short of their time; from dropNormal on, only key effects
  // play, still at short.
  catchUp: {
    dropMinor: 3,
    shorten: 5,
    burst: 7,
    dropNormal: 9,
    short: 0.4,
  },
  // others: what the other players do, live (20.13, ADR 027 "Other
  // players live", Presence.tsx). A player shares what they point at at
  // most every sample seconds, aimSample while aiming; the server keeps
  // the same gaps. Between samples another player's arrow head follows
  // its new place, closing follow of the way each 60 fps frame
  // (followStep). A hand card they hover lifts hoverLift px toward the
  // table over hover. A card they drag floats holdLift px out of their
  // hand and grows to holdGrow over hold; then it sways sway degrees to
  // each side, swayTime seconds per side, until it is played or put
  // back. These are live answers to a pointer, not effects: the speed
  // setting turns them off, but does not scale them.
  others: {
    sample: 0.35,
    aimSample: 0.25,
    follow: 0.1,
    hover: 0.15,
    hoverLift: 14,
    hoverCurve: curves.expoOut,
    hold: 1,
    holdLift: 56,
    holdGrow: 1.15,
    holdCurve: curves.expoOut,
    sway: 3,
    swayTime: 10,
    swayCurve: curves.easeInOut,
  },
  // mix: the sound mixer (20.14, ADR 027 principle 8, sounds.ts). One
  // sound never starts twice within gap seconds; a start that comes
  // sooner is dropped. While a key cue plays, every other sound goes
  // down to duck of its level over duckIn, linear; once the cue ends,
  // it comes back over duckOut, linear. A cue counts as ended cueMax
  // seconds after it starts at most, so a lost end never keeps the
  // rest quiet. The mixer moves its gain every tick seconds. These are
  // sound, not motion: the animation speed neither scales nor stops
  // them. At most players files play at once (19.135): each holds a
  // media player, which browsers cap per page; a new one past that
  // ends the oldest.
  mix: {
    gap: 0.08,
    duck: 0.2,
    duckIn: 0.7,
    duckOut: 0.7,
    cueMax: 5,
    tick: 0.02,
    players: 16,
  },
  // timer: the response timer (20.15, ADR 027 principle 9, timer.ts,
  // ActionBar.tsx). It shows only in its last show seconds: the seconds
  // left, and a fuse along the prompt's foot that burns down with
  // curve to empty at 0. Each second it loses then ticks for the player
  // it waits for, and 0 ticks timer-last. It is a clock, not an effect:
  // the animation speed neither scales nor stops it. Speed "off" and
  // reduced motion: the fuse steps once a second instead of burning.
  timer: {
    show: 20,
    curve: curves.linear,
  },
  // choice: picking among shown cards (20.18, ADR 027 "Choosing cards",
  // choice.ts, ChoiceDialog.tsx). Cards a player looks at are dealt from
  // their deck: each flies to its place over deal on a curve that rises
  // toward lift px above the deck, and turns face up over turn after
  // turnAt; the first leaves at once, the second first later, each next
  // one gap later; each plays choice-deal sound seconds after it leaves.
  // Other card choices grow at their places from 0 over grow. View board
  // shrinks the cards to 0 and fades the window over hide; back, they
  // grow again over grow. After the pick the cards fly back onto the
  // deck over back, the kept one first, the rest backFirst after it and
  // backGap apart; an order goes bottom card first, top card last. Other
  // players see small backs over the deck while it is looked at; they
  // rise out of it over others and sink back. sizes are the row's card
  // size and gap (in card widths) up to upTo cards; more than row cards
  // wrap (choiceRow).
  choice: {
    deal: 1,
    lift: 80,
    turnAt: 0.094,
    turn: 0.9,
    first: 0.19,
    gap: 0.09,
    sound: 0.04,
    grow: 0.2,
    hide: 0.2,
    back: 0.5,
    backFirst: 0.15,
    backGap: 0.06,
    others: 0.5,
    row: 6,
    sizes: [
      { upTo: 3, scale: 1, gap: 0.25 },
      { upTo: 4, scale: 1, gap: 0.2 },
      { upTo: 5, scale: 0.85, gap: 0.125 },
      { upTo: Infinity, scale: 0.7, gap: 0.1 },
    ],
    dealCurve: curves.sineInExpoOut,
    turnCurve: curves.expoOut,
    growCurve: curves.expoOut,
    backCurve: curves.easeInOut,
    othersCurve: curves.expoOut,
  },
  // start: the start-of-game screen (20.22, ADR 027 "Start of game",
  // ChoiceDialog.tsx). Its banner grows in over banner. A picked card
  // lifts lift px and keeps its bright glow. After Confirm the picked
  // card waits at its place, at most hold seconds, and flies to its
  // player over keep; the rest fly back to the deck (choice.back).
  // Others see "Bo is choosing" low on the board; it grows in over wait.
  start: {
    banner: 0.5,
    lift: 14,
    keep: 0.5,
    hold: 4,
    wait: 0.4,
    bannerCurve: curves.expoOut,
    keepCurve: curves.expoOut,
    waitCurve: curves.expoOut,
  },
  // play: a card played from the own hand (20.20, ADR 027 "Playing from
  // the hand", hold.ts, HandCard.tsx). A press that moves drag px picks
  // the card up; so does a click (a tap on touch lifts it in the fan).
  // Held, it grows to grow over growTime. Let go off the board, or put
  // back, it goes home over back from where it was. Played, it stays
  // where it was let go until the server answers, then shrinks into its
  // stack row over into; it fades out from intoFade of the eased way on,
  // when it is row-sized. The hand closes the gap a played
  // card leaves over gap. These answer the pointer: the speed setting
  // turns them off, but does not scale them.
  play: {
    drag: 8,
    grow: 1.08,
    growTime: 0.2,
    growCurve: curves.expoOut,
    back: 0.5,
    backCurve: curves.expoOut,
    into: 0.4,
    intoFade: 0.9,
    intoCurve: curves.expoOut,
    gap: 0.5,
    gapCurve: curves.expoOut,
  },
  // history: the history strip at the board's left edge (20.21, ADR 027
  // "History", HistoryStrip.tsx). A new tile slides in from left of the
  // strip, from from of its own width away, over slide with slideCurve;
  // the tiles below move down in the same time. The animation speed
  // scales it; "off" places tiles at once. While the pointer is on the
  // strip nothing moves, and new tiles wait till it leaves. Then the
  // rest of the board dims to dim (black's opacity) over dimIn with
  // dimInCurve and clears over dimOut: they answer the pointer, so the
  // speed setting turns them off but does not scale them. The strip's
  // foot fades over fade tile sides when more tiles are below.
  history: {
    slide: 1,
    from: 1.3,
    slideCurve: curves.expoOut,
    dim: 0.3,
    dimIn: 0.5,
    dimInCurve: curves.cubicInOut,
    dimOut: 0.1,
    dimOutCurve: curves.linear,
    fade: 1.2,
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

export interface FanPose {
  // angle is the card's turn in degrees, clockwise.
  angle: number;
  // drop is how far it sits below the fan's middle, in px.
  drop: number;
}

// fanPose is where card i of an own hand of n cards sits in the fan
// (20.5, ADR 027 "Hand"): fewer than feel.hand.from lie straight; a
// bigger hand turns by its place from the middle and drops on the
// parabola (i + 0.5 − n/2)² / (6n).
export function fanPose(i: number, n: number): FanPose {
  const { from, angle, perCard, drop } = feel.hand;
  if (n < from) {
    return { angle: 0, drop: 0 };
  }
  const spread = i + 0.5 - n / 2;
  return {
    angle: ((angle + perCard * n) / (n - 1)) * spread,
    drop: ((spread * spread) / (6 * n)) * drop,
  };
}

// drawPace is the time factor of a player's i-th draw in one update,
// counted from 0 (20.7, ADR 027 principle 6): the first three at full
// time, then each 11 % faster, down to 0.556.
export function drawPace(i: number): number {
  const { fastFrom, faster, fastest } = feel.draw;
  if (i < fastFrom) {
    return 1;
  }
  return Math.max(fastest, 1 - faster * (i + 1 - fastFrom));
}

// countTimes is when each step of a counter's change of n starts, in
// seconds from the change (20.11, ADR 027 "Counters"): the first at
// once, then feel.count.step apart; from more than five steps on they
// come faster, so all of them fit in feel.count.fit.
export function countTimes(n: number): number[] {
  const { step, fit } = feel.count;
  const steps = Math.max(0, Math.round(Math.abs(n)));
  const gap = steps > 0 ? Math.min(step, fit / steps) : 0;
  return Array.from({ length: steps }, (_, i) => i * gap);
}

// followStep is how far another player's arrow head moves toward its
// new place in a frame of dt seconds, as a part of the way (20.13): the
// reference closes feel.others.follow of it each frame at 60 fps; any
// frame rate closes it alike.
export function followStep(dt: number): number {
  if (dt <= 0) {
    return 0;
  }
  return 1 - Math.pow(1 - feel.others.follow, dt * 60);
}
