import { motion } from 'motion/react';
import { useEffect, useRef, useState } from 'react';
import { diceDelays, diceImage } from '../dice';
import { t } from '../i18n';
import { playSound, type Sound } from '../sounds';
import { useApp, type Playing } from '../store';
import { Card } from './Card';
import { Die, DieCube } from './Die';
import { playerColor } from '../table';
import { Coin } from './PlayerArea';
import './Effects.css';
import { lunge } from './lunge';
import { hits, skulls } from './marks';

// The size of a table card, for flying cards.
const cardW = 80;
const cardH = 112;
// The width of a large card (size-zoom); its height follows.
const zoomW = 320;
const zoomH = (zoomW * cardH) / cardW;
// The width of a hand card (size-hand).
const handW = 100;

interface Point {
  x: number;
  y: number;
}

// zoneRect finds the first zone on screen (data-zone) and its box. The
// large played card's zone is its caption too; its card is used (15.74).
function zoneRect(zones: string[]): DOMRect | null {
  for (const z of zones) {
    const el = document.querySelector(`[data-zone="${CSS.escape(z)}"]`);
    if (el) {
      const card = z.startsWith('played-') ? el.querySelector('.card') : null;
      return (card ?? el).getBoundingClientRect();
    }
  }
  return null;
}

// zoneCenter finds the first zone on screen (data-zone) and its centre.
function zoneCenter(zones: string[]): Point | null {
  const r = zoneRect(zones);
  return r && { x: r.left + r.width / 2, y: r.top + r.height / 2 };
}

// playedWidth is the width of the large played card when it is the
// first of the zones on screen, else 0.
function playedWidth(zones: string[]): number {
  const z = zones.find((z) => document.querySelector(`[data-zone="${CSS.escape(z)}"]`));
  return z?.startsWith('played-') ? (zoneRect([z])?.width ?? 0) : 0;
}

// middle is the centre of the window, where large cards show.
function middle(): Point {
  return { x: window.innerWidth / 2, y: window.innerHeight / 2 };
}

// Effects draws the event effect playing now (7.6). Each effect calls
// done when its animation ends; playback then goes on. Cards that place
// effects left on their slots stay below it (19.3).
export function Effects() {
  const playing = useApp((s) => s.playing);
  const placed = useApp((s) => s.placed);
  return (
    <>
      {placed.length > 0 && (
        <div className="effects">
          {placed.map((p) => (
            <div
              key={p.key}
              className="effect-fly"
              style={{ transform: `translate(${String(p.x)}px, ${String(p.y)}px)` }}
              data-testid="placed"
            >
              <Card cardRef={p.card} size="table" />
            </div>
          ))}
        </div>
      )}
      {playing && (
        <div className="effects" data-testid={`effect-${playing.effect.kind}`}>
          <Effect key={playing.key} playing={playing} />
        </div>
      )}
    </>
  );
}

function Effect({ playing }: { playing: Playing }) {
  const { effect, scale, done } = playing;
  // A die is drawn in its roller's look (15.6).
  const image = useApp((s) => (effect.kind === 'dice' ? diceImage(s.dice, s.game?.seats, effect.player) : ''));
  const duration = (seconds: number) => seconds * scale;
  switch (effect.kind) {
    case 'fly':
    case 'place':
      if (effect.middle && effect.card) {
        return <FlyFromMiddle to={effect.to} card={effect.card} playing={playing} stay={effect.kind === 'place'} />;
      }
      return (
        <Fly
          from={effect.from}
          to={effect.to}
          card={effect.card}
          back={effect.back}
          playing={playing}
          stay={effect.kind === 'place'}
          sound={effect.kind === 'fly' ? effect.sound : undefined}
        />
      );
    case 'steal':
      return <Steal from={effect.from} to={effect.to} amount={effect.amount} playing={playing} />;
    case 'popup':
      return <Popup at={effect.at} text={effect.text} tone={effect.tone} playing={playing} />;
    case 'dice':
      if (effect.thrown) {
        return <ThrownDie values={effect.dice ?? [effect.value]} label={effect.label} image={image} playing={playing} />;
      }
      return (
        <ChangedDie
          value={effect.value}
          old={effect.old}
          source={effect.source}
          label={effect.label}
          player={effect.player}
          image={image}
          playing={playing}
        />
      );
    case 'banner':
      return (
        // anim: banner
        <motion.div
          className="effect-center effect-banner"
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: [0, 1, 1, 0], scale: [0.9, 1, 1, 1] }}
          transition={{ duration: duration(effect.long ? 2.2 : 1.1), times: effect.long ? [0, 0.1, 0.9, 1] : [0, 0.2, 0.8, 1] }}
          onAnimationComplete={done}
        >
          {effect.text}
        </motion.div>
      );
    case 'reveal':
      return <Reveal card={effect.card} deck={effect.deck} hold={effect.hold ?? false} playing={playing} />;
    case 'played':
      return <Played card={effect.card} who={effect.who} from={effect.from} to={effect.to} playing={playing} />;
    case 'lunge':
      return (
        <Layered
          playing={playing}
          start={(done) => {
            lunge(effect.from, effect.to, effect.amount, scale, done);
          }}
        />
      );
    case 'hits':
      return (
        <Layered
          playing={playing}
          start={(done) => {
            hits(effect.hits, scale, done);
          }}
        />
      );
    case 'deaths':
      return (
        <Layered
          playing={playing}
          start={(done) => {
            skulls(effect.at, scale, done, effect.sound);
          }}
        />
      );
    case 'sound':
      // Nothing moves: the sound plays and a short beat holds the rest (18.14).
      return (
        <Layered
          playing={playing}
          start={(done) => {
            playSound(effect.sound);
            window.setTimeout(done, duration(0.6));
          }}
        />
      );
  }
}

// started is the last layered effect started: Strict Mode runs effects
// twice.
let started = 0;

// Layered starts an effect that plays in its own layer (lunge.ts,
// marks.ts) and calls done when playback may go on.
function Layered({ playing, start }: { playing: Playing; start: (done: () => void) => void }) {
  useEffect(() => {
    if (playing.key > started) {
      started = playing.key;
      start(playing.done);
    }
  }, [playing, start]);
  return null;
}

// useSkip ends an effect at once when it has nowhere to play.
function useSkip(skip: boolean, done: () => void) {
  useEffect(() => {
    if (skip) {
      done();
    }
  }, [skip, done]);
}

function Fly({
  from,
  to,
  card,
  back,
  playing,
  stay = false,
  sound,
}: {
  from: string[];
  to: string[];
  card: string | undefined;
  back: Parameters<typeof Card>[0]['back'];
  playing: Playing;
  // stay leaves the card where it landed until the next view (19.3).
  stay?: boolean;
  // sound plays once as the card leaves (18.7).
  sound?: Sound;
}) {
  // Read the zones once, when the effect starts. A card leaving the
  // large played card starts at its size (15.74).
  const [path] = useState(() => ({ a: zoneCenter(from), b: zoneCenter(to), grow: playedWidth(from) / cardW }));
  useSkip(!path.a || !path.b, playing.done);
  const flies = Boolean(path.a && path.b);
  // A timer, so StrictMode's second mount does not play it twice.
  useEffect(() => {
    if (!flies || !sound) {
      return;
    }
    const id = window.setTimeout(() => {
      playSound(sound);
    }, 0);
    return () => {
      window.clearTimeout(id);
    };
  }, [flies, sound]);
  if (!path.a || !path.b) {
    return null;
  }
  const end = { x: path.b.x - cardW / 2, y: path.b.y - cardH / 2 };
  const land = () => {
    if (stay && card) {
      useApp.setState((s) => ({ placed: [...s.placed, { key: playing.key, card, ...end }] }));
    }
    playing.done();
  };
  return (
    // anim: fly
    <motion.div
      className="effect-fly"
      initial={{ x: path.a.x - cardW / 2, y: path.a.y - cardH / 2, scale: path.grow || 0.9 }}
      animate={{ ...end, scale: 1 }}
      transition={{ duration: 0.45 * playing.scale, ease: 'easeInOut' }}
      onAnimationComplete={land}
    >
      <Card cardRef={card} back={back} size="table" />
    </motion.div>
  );
}

// Reveal shows a card large (15.73): it grows out of its deck to the
// middle and stays a moment. Held, it ends there for the move after it;
// else it shrinks back onto the deck, where the card stays.
function Reveal({ card, deck, hold, playing }: { card: string; deck: string[]; hold: boolean; playing: Playing }) {
  const [from] = useState(() => {
    const r = zoneRect(deck);
    const m = middle();
    return r && { x: r.left + r.width / 2 - m.x, y: r.top + r.height / 2 - m.y, scale: Math.min(1, r.width / zoomW) };
  });
  // Without the deck on screen, the card fades in and out in the middle.
  const off = from ?? { x: 0, y: 0, scale: 0.8 };
  const fade = from ? 1 : 0;
  const end = hold ? { x: 0, y: 0, scale: 1, opacity: 1 } : { ...off, opacity: 0 };
  return (
    // anim: reveal
    <motion.div
      className="effect-center"
      initial={{ x: off.x, y: off.y, scale: off.scale, opacity: fade }}
      animate={{
        x: [off.x, 0, 0, end.x, end.x],
        y: [off.y, 0, 0, end.y, end.y],
        scale: [off.scale, 1, 1, end.scale, end.scale],
        opacity: [fade, 1, 1, hold ? 1 : fade, end.opacity],
      }}
      transition={{ duration: 1.5 * playing.scale, times: [0, 0.25, 0.75, 0.95, 1], ease: 'easeInOut' }}
      onAnimationComplete={playing.done}
    >
      <Card cardRef={card} size="zoom" />
    </motion.div>
  );
}

// cardOffset is where a zone's card is, from the middle, and the scale
// that makes a large card its size; null when the zone is not drawn.
function cardOffset(zones: string[], width: (r: DOMRect) => number): { x: number; y: number; scale: number } | null {
  const r = zoneRect(zones);
  if (!r) {
    return null;
  }
  const m = middle();
  return { x: r.left + r.width / 2 - m.x, y: r.top + r.height / 2 - m.y, scale: Math.min(1, width(r) / zoomW) };
}

// Played shows another player's loot card (15.31, 15.74): it rises from
// their hand to the middle, stays a moment with "Ann plays", then
// shrinks into the stack panel, to the size of a hand card.
function Played({ card, who, from, to, playing }: { card: string; who: string; from: string[]; to: string[]; playing: Playing }) {
  const [path] = useState(() => ({
    a: cardOffset(from, () => cardW),
    b: cardOffset(to, () => handW),
  }));
  const a = path.a ?? { x: 0, y: 0, scale: 0.6 };
  const b = path.b ?? { x: 0, y: 0, scale: 0.5 };
  return (
    // anim: played
    <motion.div
      className="effect-center"
      initial={{ x: a.x, y: a.y, scale: a.scale, opacity: 0 }}
      animate={{
        x: [a.x, 0, 0, b.x],
        y: [a.y, 0, 0, b.y],
        scale: [a.scale, 1, 1, b.scale],
        opacity: [0, 1, 1, path.b ? 1 : 0],
      }}
      transition={{ duration: 1.5 * playing.scale, times: [0, 0.25, 0.75, 1], ease: 'easeInOut' }}
      onAnimationComplete={playing.done}
      data-testid="played-pop"
    >
      {/* anim: played */}
      <motion.div
        className="effect-label"
        animate={{ opacity: [1, 1, 0] }}
        transition={{ duration: 1.5 * playing.scale, times: [0, 0.75, 0.85] }}
      >
        {t('effect.plays', { player: who })}
      </motion.div>
      <Card cardRef={card} size="zoom" />
    </motion.div>
  );
}

// FlyFromMiddle moves a card a held reveal left large in the middle to
// its place, shrinking to a table card on the way (15.73).
function FlyFromMiddle({ to, card, playing, stay }: { to: string[]; card: string; playing: Playing; stay: boolean }) {
  const [path] = useState(() => ({ a: middle(), b: zoneCenter(to) }));
  useSkip(!path.b, playing.done);
  if (!path.b) {
    return null;
  }
  const b = path.b;
  const land = () => {
    if (stay) {
      useApp.setState((s) => ({ placed: [...s.placed, { key: playing.key, card, x: b.x - cardW / 2, y: b.y - cardH / 2 }] }));
    }
    playing.done();
  };
  return (
    // anim: fly
    <motion.div
      className="effect-fly"
      initial={{ x: path.a.x - zoomW / 2, y: path.a.y - zoomH / 2, scale: 1 }}
      animate={{ x: b.x - zoomW / 2, y: b.y - zoomH / 2, scale: cardW / zoomW }}
      transition={{ duration: 0.55 * playing.scale, ease: 'easeInOut' }}
      onAnimationComplete={land}
    >
      <Card cardRef={card} size="zoom" />
    </motion.div>
  );
}

function Popup({ at, text, tone, playing }: { at: string[]; text: string; tone: string; playing: Playing }) {
  const [point] = useState(() => zoneCenter(at));
  useSkip(!point, playing.done);
  if (!point) {
    return null;
  }
  return (
    // anim: popup
    <motion.div
      className={`effect-popup effect-${tone}`}
      style={{ left: point.x, top: point.y }}
      initial={{ opacity: 0, y: 0 }}
      animate={{ opacity: [0, 1, 1, 0], y: [0, -12, -24, -32] }}
      transition={{ duration: 0.8 * playing.scale }}
      onAnimationComplete={playing.done}
    >
      {text}
    </motion.div>
  );
}

// When a changed roll's die is in, its shot lands, and its face turns,
// as parts of the effect (15.76).
const dieIn = 0.15;
const shotLands = 0.4;
const faceTurns = 0.47;

// ChangedDie shows a roll change (15.76): the die spins in with its old
// value, a shot flies to it from the card that changes it, and it turns
// to the new value. Without an old value it shows the new one only.
function ChangedDie({
  value,
  old,
  source,
  label,
  player,
  image,
  playing,
}: {
  value: number;
  old: number | undefined;
  source: string[] | undefined;
  label: string;
  player: number;
  image: string;
  playing: Playing;
}) {
  const total = 1.8 * playing.scale;
  const [face, setFace] = useState(old ?? value);
  const dieRef = useRef<HTMLDivElement>(null);
  const [shot, setShot] = useState<{ a: Point; b: Point } | null>(null);
  useEffect(() => {
    const el = dieRef.current;
    const a = source ? zoneCenter(source) : null;
    if (el && a) {
      const r = el.getBoundingClientRect();
      setShot({ a, b: { x: r.left + r.width / 2, y: r.top + r.height / 2 } });
    }
    const turn = window.setTimeout(() => {
      setFace(value);
    }, faceTurns * total * 1000);
    return () => {
      window.clearTimeout(turn);
    };
  }, [source, value, total]);
  return (
    <>
      <div className="effect-center">
        {/* anim: dice */}
        <motion.div
          ref={dieRef}
          initial={{ opacity: 0, scale: 0.4, rotate: -180 }}
          animate={{ opacity: [0, 1, 1, 0], scale: [0.4, 1, 1, 1.15, 1, 1], rotate: [-180, 0, 0, 360, 360] }}
          transition={{
            opacity: { duration: total, times: [0, dieIn, 0.88, 1] },
            scale: { duration: total, times: [0, dieIn, shotLands, faceTurns, 0.55, 1] },
            rotate: { duration: total, times: [0, dieIn, shotLands, 0.55, 1] },
          }}
          onAnimationComplete={playing.done}
          data-testid="changed-die"
          data-face={face}
        >
          <Die value={face} image={image} />
        </motion.div>
        {label && (
          <div className="effect-label" style={{ color: playerColor(player) }}>
            {label}
          </div>
        )}
      </div>
      {shot && (
        // anim: roll-shot
        <motion.div
          className="effect-shot effect-shot-roll"
          style={{ left: shot.a.x, top: shot.a.y }}
          initial={{ x: 0, y: 0, scale: 0.4, opacity: 0 }}
          animate={{ x: [0, 0, shot.b.x - shot.a.x], y: [0, 0, shot.b.y - shot.a.y], scale: [0.4, 1, 1.2], opacity: [0, 1, 1, 0] }}
          transition={{
            duration: total * shotLands,
            times: [0, dieIn / shotLands, 1],
            ease: 'easeIn',
            opacity: { duration: total * shotLands, times: [0, dieIn / shotLands, 0.95, 1] },
          }}
          data-testid="roll-shot"
        />
      )}
    </>
  );
}

// coinAt is the coin counter of a player's zone, else the zone's
// centre (15.75).
function coinAt(zones: string[]): Point | null {
  for (const z of zones) {
    const el = document.querySelector(`[data-zone="${CSS.escape(z)}"]`);
    if (el) {
      const r = (el.querySelector('.coin') ?? el).getBoundingClientRect();
      return { x: r.left + r.width / 2, y: r.top + r.height / 2 };
    }
  }
  return null;
}

// Steal shows cents changing hands (15.75): "−N¢" rises on the victim's
// coins as a coin leaves them, and "+N¢ stolen" rises on the thief's
// when it lands.
function Steal({ from, to, amount, playing }: { from: string[]; to: string[]; amount: number; playing: Playing }) {
  const [path] = useState(() => ({ a: coinAt(from), b: coinAt(to) }));
  useSkip(!path.a || !path.b, playing.done);
  if (!path.a || !path.b) {
    return null;
  }
  const { a, b } = path;
  const total = 1.5 * playing.scale;
  // The coin flies in the first part; the thief's text follows.
  const fly = 0.45;
  return (
    <>
      {/* anim: popup */}
      <motion.div
        className="effect-popup effect-bad"
        style={{ left: a.x, top: a.y }}
        initial={{ opacity: 0, y: 0 }}
        animate={{ opacity: [0, 1, 1, 0], y: [0, -12, -24, -32] }}
        transition={{ duration: 0.8 * playing.scale }}
      >
        −{amount}¢
      </motion.div>
      {/* anim: coin-shot */}
      <motion.div
        className="effect-coin"
        style={{ left: a.x, top: a.y }}
        initial={{ x: 0, y: 0, scale: 1, opacity: 0 }}
        animate={{ x: [0, (b.x - a.x) / 2, b.x - a.x], y: [0, (b.y - a.y) / 2 - 40, b.y - a.y], scale: [1, 1.6, 1], opacity: [1, 1, 0] }}
        transition={{ duration: total * fly, ease: 'easeInOut', opacity: { duration: total * fly, times: [0, 0.9, 1] } }}
        data-testid="coin-shot"
      >
        <Coin />
      </motion.div>
      {/* anim: popup */}
      <motion.div
        className="effect-popup effect-good"
        style={{ left: b.x, top: b.y }}
        initial={{ opacity: 0, y: 0 }}
        animate={{ opacity: [0, 0, 1, 1, 0], y: [0, 0, -12, -24, -32] }}
        transition={{ duration: total, times: [0, fly, fly + 0.1, 0.85, 1] }}
        onAnimationComplete={playing.done}
      >
        {t('effect.stolen', { n: amount })}
      </motion.div>
    </>
  );
}

// The normal length of a throw, in ms.
const throwMs = 1600;
// When the die lands, bounces and lands again, as parts of the throw.
const landings = [0.28, 0.5] as const;

// fall is the thrown die's path (15.62, V-17): it falls from above the
// screen, bounces twice, lies still, then fades.
const fall: Keyframe[] = [
  { offset: 0, transform: 'translateY(-60vh)', opacity: 0, easing: 'ease-in' },
  { offset: 0.05, opacity: 1, easing: 'ease-in' },
  { offset: landings[0], transform: 'translateY(0)', easing: 'ease-out' },
  { offset: 0.39, transform: 'translateY(-70px)', easing: 'ease-in' },
  { offset: landings[1], transform: 'translateY(0)', easing: 'ease-out' },
  { offset: 0.56, transform: 'translateY(-16px)', easing: 'ease-in' },
  { offset: 0.62, transform: 'translateY(0)' },
  { offset: 0.88, transform: 'translateY(0)', opacity: 1 },
  { offset: 1, transform: 'translateY(0)', opacity: 0 },
];

// tumble turns the cube as it falls (15.69): fast in the air, slower at
// each bounce, then at rest with the value on top, seen from above with
// two sides showing (15.82).
const rest = 'rotateX(-58deg) rotateY(-28deg) rotateZ(0deg)';
const tumble: Keyframe[] = [
  { offset: 0, transform: 'rotateX(740deg) rotateY(-560deg) rotateZ(150deg)', easing: 'linear' },
  { offset: landings[0], transform: 'rotateX(200deg) rotateY(-170deg) rotateZ(40deg)', easing: 'ease-out' },
  { offset: landings[1], transform: 'rotateX(-30deg) rotateY(-55deg) rotateZ(8deg)', easing: 'ease-out' },
  { offset: 0.62, transform: rest },
  { offset: 1, transform: rest },
];

// shadow is the die's shadow on the table: small and faint while the
// die is high, full where it lands.
const shadow: Keyframe[] = [
  { offset: 0, transform: 'scale(0.3)', opacity: 0, easing: 'ease-in' },
  { offset: landings[0], transform: 'scale(1)', opacity: 0.6, easing: 'ease-out' },
  { offset: 0.39, transform: 'scale(0.75)', opacity: 0.35, easing: 'ease-in' },
  { offset: landings[1], transform: 'scale(1)', opacity: 0.6 },
  { offset: 0.88, transform: 'scale(1)', opacity: 0.6 },
  { offset: 1, transform: 'scale(1)', opacity: 0 },
];

// The label shows once the die has landed.
const labelIn: Keyframe[] = [
  { offset: 0, opacity: 0 },
  { offset: landings[0], opacity: 0 },
  { offset: 0.4, opacity: 1 },
  { offset: 0.88, opacity: 1 },
  { offset: 1, opacity: 0 },
];

// ThrownDie throws a fresh roll's die: a cube falls, bounces and tumbles
// to the engine's value (15.62, 15.69), and clatters as it first lands.
// A roll of several dice throws them all side by side, at small uneven
// gaps, and each die clatters as it lands (19.39). Its label says the value all
// along, so the result never differs from the engine. Web Animations on
// transform and opacity only (ADR 002); scale follows the animation
// speed (ST-05).
function ThrownDie({ values, label, image, playing }: { values: number[]; label: string; image: string; playing: Playing }) {
  const rowRef = useRef<HTMLDivElement>(null);
  const wordsRef = useRef<HTMLDivElement>(null);
  const { scale, done } = playing;
  const count = values.length;
  useEffect(() => {
    const spots = [...(rowRef.current?.children ?? [])];
    const delays = diceDelays(spots.length);
    const ms = throwMs * scale;
    const start = performance.now();
    const parts: (Animation | undefined)[] = [];
    spots.forEach((spot, i) => {
      const delay = (delays.at(i) ?? 0) * ms;
      const timing = { duration: ms, delay, fill: 'both' } as const;
      // anim: dice-throw
      const a = spot.querySelector('.effect-thrown')?.animate(fall, timing);
      if (a && i === spots.length - 1) {
        a.onfinish = done;
      }
      // anim: dice-throw
      parts.push(a, spot.querySelector('.die-cube')?.animate(tumble, timing), spot.querySelector('.effect-thrown-shadow')?.animate(shadow, timing));
    });
    // anim: dice-throw
    parts.push(wordsRef.current?.animate(labelIn, { duration: ms, fill: 'forwards' }));
    // Each die clatters as it first lands, one after another.
    const lands = delays.map((at) => (at + landings[0]) * ms);
    let heard = 0;
    const clatter = window.setInterval(() => {
      const now = performance.now() - start;
      while (heard < lands.length && (lands.at(heard) ?? 0) <= now) {
        playSound('diceRoll');
        heard++;
      }
      if (heard >= lands.length) {
        window.clearInterval(clatter);
      }
    }, 15);
    return () => {
      for (const p of parts) {
        p?.cancel();
      }
      window.clearInterval(clatter);
    };
  }, [count, scale, done]);
  return (
    <div className="effect-center">
      <div ref={rowRef} className="effect-throw-row">
        {values.map((value, i) => (
          // The dice keep their order while they fall.
          // eslint-disable-next-line @eslint-react/no-array-index-key
          <div key={i} className="effect-throw-spot">
            <div className="effect-thrown-shadow" />
            <div className="effect-thrown" data-testid="thrown-die">
              <DieCube value={value} image={image} />
            </div>
          </div>
        ))}
      </div>
      {label && (
        <div ref={wordsRef} className="effect-label">
          {label}
        </div>
      )}
    </div>
  );
}
