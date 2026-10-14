import { motion } from 'motion/react';
import { useEffect, useState } from 'react';
import { useApp, type Playing } from '../store';
import { Card } from './Card';
import { Die } from './Die';
import './Effects.css';
import { lunge } from './lunge';
import { hits, skulls } from './marks';

// The size of a table card, for flying cards.
const cardW = 80;
const cardH = 112;

interface Point {
  x: number;
  y: number;
}

// zoneCenter finds the first zone on screen (data-zone) and its centre.
function zoneCenter(zones: string[]): Point | null {
  for (const z of zones) {
    const el = document.querySelector(`[data-zone="${CSS.escape(z)}"]`);
    if (el) {
      const r = el.getBoundingClientRect();
      return { x: r.left + r.width / 2, y: r.top + r.height / 2 };
    }
  }
  return null;
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
  const diceImage = useApp((s) => s.dice.find((d) => d.id === s.settings.dice)?.image ?? '');
  const duration = (seconds: number) => seconds * scale;
  switch (effect.kind) {
    case 'fly':
      return <Fly from={effect.from} to={effect.to} card={effect.card} back={effect.back} playing={playing} />;
    case 'place':
      return <Fly from={effect.from} to={effect.to} card={effect.card} back={effect.back} playing={playing} stay />;
    case 'popup':
      return <Popup at={effect.at} text={effect.text} tone={effect.tone} playing={playing} />;
    case 'dice':
      return (
        // anim: dice
        <motion.div
          className="effect-center"
          initial={{ opacity: 0, scale: 0.4, rotate: -180 }}
          animate={{ opacity: [0, 1, 1, 0], scale: [0.4, 1.15, 1, 1], rotate: [-180, 0, 0, 0] }}
          transition={{ duration: duration(1.1), times: [0, 0.3, 0.8, 1] }}
          onAnimationComplete={done}
        >
          <Die value={effect.value} image={diceImage} />
          {effect.label && <div className="effect-label">{effect.label}</div>}
        </motion.div>
      );
    case 'banner':
      return (
        // anim: banner
        <motion.div
          className="effect-center effect-banner"
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: [0, 1, 1, 0], scale: [0.9, 1, 1, 1] }}
          transition={{ duration: duration(1.1), times: [0, 0.2, 0.8, 1] }}
          onAnimationComplete={done}
        >
          {effect.text}
        </motion.div>
      );
    case 'reveal':
      return (
        // anim: reveal
        <motion.div
          className="effect-center"
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: [0, 1, 1, 0], scale: [0.8, 1, 1, 0.9] }}
          transition={{ duration: duration(1.3), times: [0, 0.2, 0.8, 1] }}
          onAnimationComplete={done}
        >
          <Card cardRef={effect.card} size="zoom" />
        </motion.div>
      );
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
            skulls(effect.at, scale, done);
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
}: {
  from: string[];
  to: string[];
  card: string | undefined;
  back: Parameters<typeof Card>[0]['back'];
  playing: Playing;
  // stay leaves the card where it landed until the next view (19.3).
  stay?: boolean;
}) {
  // Read the zones once, when the effect starts.
  const [path] = useState(() => ({ a: zoneCenter(from), b: zoneCenter(to) }));
  useSkip(!path.a || !path.b, playing.done);
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
      initial={{ x: path.a.x - cardW / 2, y: path.a.y - cardH / 2, scale: 0.9 }}
      animate={{ ...end, scale: 1 }}
      transition={{ duration: 0.45 * playing.scale, ease: 'easeInOut' }}
      onAnimationComplete={land}
    >
      <Card cardRef={card} back={back} size="table" />
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
