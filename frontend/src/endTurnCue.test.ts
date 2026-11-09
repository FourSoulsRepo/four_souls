import { describe, expect, it } from 'vitest';
import { createEndTurnCue, type Clock } from './endTurnCue';
import { feel } from './feel';

// fakeClock runs timers by hand: advance moves its time on.
function fakeClock() {
  let now = 0;
  let next = 1;
  const timers = new Map<number, { at: number; fn: () => void }>();
  const clock: Clock = {
    set: (fn, ms) => {
      timers.set(next, { at: now + ms, fn });
      return next++;
    },
    clear: (id) => {
      timers.delete(id);
    },
  };
  const advance = (ms: number) => {
    now += ms;
    for (const [id, t] of [...timers]) {
      if (t.at <= now) {
        timers.delete(id);
        t.fn();
      }
    }
  };
  return { clock, advance, pending: () => timers.size };
}

const wait = feel.endTurn.cue * 1000;

describe('createEndTurnCue (20.10)', () => {
  it('plays the cue feel.endTurn.cue after the button turns done', () => {
    const { clock, advance } = fakeClock();
    let played = 0;
    const cue = createEndTurnCue(() => played++, clock);
    cue.update(3, true);
    advance(wait - 1);
    expect(played).toBe(0);
    advance(1);
    expect(played).toBe(1);
  });

  it('plays once per turn at most', () => {
    const { clock, advance } = fakeClock();
    let played = 0;
    const cue = createEndTurnCue(() => played++, clock);
    cue.update(3, true);
    advance(wait);
    cue.update(3, false);
    cue.update(3, true);
    advance(wait * 2);
    expect(played).toBe(1);
    // The next own turn has a cue of its own.
    cue.update(5, true);
    advance(wait);
    expect(played).toBe(2);
  });

  it('skips the cue if a play became possible meanwhile, and tries no more that turn', () => {
    const { clock, advance } = fakeClock();
    let played = 0;
    const cue = createEndTurnCue(() => played++, clock);
    cue.update(3, true);
    advance(wait / 2);
    cue.update(3, false);
    advance(wait);
    expect(played).toBe(0);
    cue.update(3, true);
    advance(wait);
    expect(played).toBe(0);
  });

  it('skips the cue once the turn changed', () => {
    const { clock, advance } = fakeClock();
    let played = 0;
    const cue = createEndTurnCue(() => played++, clock);
    cue.update(3, true);
    cue.update(4, false);
    advance(wait);
    expect(played).toBe(0);
  });

  it('plays nothing while the button is not done, and stop drops a waiting cue', () => {
    const { clock, advance, pending } = fakeClock();
    let played = 0;
    const cue = createEndTurnCue(() => played++, clock);
    cue.update(3, false);
    expect(pending()).toBe(0);
    cue.update(3, true);
    cue.stop();
    expect(pending()).toBe(0);
    advance(wait);
    expect(played).toBe(0);
  });
});
