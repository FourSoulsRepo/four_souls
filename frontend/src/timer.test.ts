import { describe, expect, it } from 'vitest';
import { feel } from './feel';
import { secondsLeft } from './stack';
import { nextChange, tickSound, timerShows } from './timer';

const deadline = 100_000;
const show = feel.timer.show * 1000;

describe('the response timer (20.15)', () => {
  it('shows only in its last 20 s', () => {
    expect(feel.timer.show).toBe(20);
    expect(timerShows(deadline, deadline - 25_000)).toBe(false);
    expect(timerShows(deadline, deadline - show - 1)).toBe(false);
    expect(timerShows(deadline, deadline - show)).toBe(true);
    expect(timerShows(deadline, deadline - 15_000)).toBe(true);
    expect(timerShows(deadline, deadline + 500)).toBe(true);
    expect(timerShows(0, 0)).toBe(false);
  });

  it('wakes as it shows, then at each new second, and not after 0', () => {
    expect(nextChange(deadline, deadline - 25_000)).toBe(5000);
    expect(nextChange(deadline, deadline - 20_000)).toBe(1000);
    expect(nextChange(deadline, deadline - 15_300)).toBe(300);
    expect(nextChange(deadline, deadline - 400)).toBe(400);
    expect(nextChange(deadline, deadline)).toBeUndefined();
    expect(nextChange(0, 0)).toBeUndefined();
  });

  it('ticks each second from 19 down, and the last at 0, as the waits land', () => {
    const sounds: string[] = [];
    let now = deadline - 20_000;
    let last = secondsLeft(deadline, now);
    for (;;) {
      const wait = nextChange(deadline, now);
      if (wait === undefined) {
        break;
      }
      now += wait;
      const left = secondsLeft(deadline, now);
      sounds.push(tickSound(last, left) ?? 'none');
      last = left;
    }
    expect(sounds).toEqual([...Array<string>(19).fill('timerTick'), 'timerLast']);
  });

  it('ticks once for a jump, and never while the seconds stay', () => {
    expect(tickSound(10, 10)).toBeUndefined();
    expect(tickSound(10, 6)).toBe('timerTick');
    expect(tickSound(3, 0)).toBe('timerLast');
    expect(tickSound(0, 0)).toBeUndefined();
  });
});
