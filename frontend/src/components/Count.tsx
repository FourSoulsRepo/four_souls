import { useEffect, useRef, useState, type RefObject } from 'react';
import { countTimes, feel } from '../feel';
import { useApp, viewScale } from '../store';
import { still } from './lunge';

// useCount is a number that counts to value one step at a time (20.11,
// ADR 027 "Counters"): a change of N plays as N steps, 0.2 s apart, or
// faster to fit 1 s (countTimes), at the animation speed. A change
// while it counts goes on from the number on screen. With animations
// off, or while a rollback's table shows, the number jumps.
export function useCount(value: number): number {
  const scale = useApp(viewScale);
  const [shown, setShown] = useState(value);
  const nowRef = useRef(value);
  useEffect(() => {
    const from = nowRef.current;
    if (from === value) {
      return undefined;
    }
    const dir = Math.sign(value - from);
    const times = scale > 0 ? countTimes(value - from) : [0];
    let i = 0;
    const step = () => {
      i++;
      nowRef.current = scale > 0 ? from + dir * i : value;
      setShown(nowRef.current);
      const at = times.at(i);
      if (at !== undefined) {
        timer = window.setTimeout(step, (at - (times.at(i - 1) ?? 0)) * scale * 1000);
      }
    };
    let timer = window.setTimeout(step, (times.at(0) ?? 0) * scale * 1000);
    return () => {
      window.clearTimeout(timer);
    };
  }, [value, scale]);
  return scale > 0 ? shown : value;
}

// useFlip flips el in each time key changes (20.11): it turns in from
// edge-on about its horizontal axis over feel.count.flip, from above
// when key grows, from below when it falls. Speed "off" and reduced
// motion: no flip.
export function useFlip(el: RefObject<HTMLElement | null>, key: number): void {
  const scale = useApp(viewScale);
  const lastRef = useRef(key);
  useEffect(() => {
    if (lastRef.current === key) {
      return;
    }
    const up = key > lastRef.current;
    lastRef.current = key;
    const node = el.current;
    if (!node || scale === 0 || still()) {
      return;
    }
    const { flip, turn, flipCurve } = feel.count;
    // anim: count-flip
    node.animate([{ transform: `rotateX(${String(up ? turn : -turn)}deg)` }, { transform: 'rotateX(0deg)' }], {
      duration: flip * scale * 1000,
      easing: flipCurve,
    });
  }, [el, key, scale]);
}

// Count is a counter's number that counts (20.11): each step flips the
// number in. data-count is the number on screen, for tests.
export function Count({ value, className = '' }: { value: number; className?: string }) {
  const shown = useCount(value);
  const ref = useRef<HTMLSpanElement>(null);
  useFlip(ref, shown);
  return (
    <span ref={ref} className={`count${className ? ` ${className}` : ''}`} data-count={shown}>
      {shown}
    </span>
  );
}
