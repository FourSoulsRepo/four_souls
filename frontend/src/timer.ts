// The response timer (20.15, ADR 027 principle 9): time pressure shows
// only near the end. The timer shows in its last feel.timer.show
// seconds; each whole second it loses then ticks, and 0 ticks the last.
import { feel } from './feel';
import type { Sound } from './sounds';
import { secondsLeft } from './stack';

// timerShows tells if the response timer that answers at deadline, Unix
// ms, shows at now: only in its last feel.timer.show seconds. 0 is no
// timer.
export function timerShows(deadline: number, now: number): boolean {
  return deadline > 0 && deadline - now <= feel.timer.show * 1000;
}

// nextChange is how many ms from now the timer next changes on screen:
// it shows, or its seconds go down by one. Undefined once it is at 0,
// or with no timer.
export function nextChange(deadline: number, now: number): number | undefined {
  const left = deadline - now;
  if (deadline <= 0 || left <= 0) {
    return undefined;
  }
  if (!timerShows(deadline, now)) {
    return left - feel.timer.show * 1000;
  }
  return left - (secondsLeft(deadline, now) - 1) * 1000;
}

// tickSound is the sound for the timer going from last seconds left to
// left: a tick for each new second, the last tick at 0; none while the
// seconds stay. A jump of several seconds ticks once.
export function tickSound(last: number, left: number): Sound | undefined {
  if (left >= last) {
    return undefined;
  }
  return left === 0 ? 'timerLast' : 'timerTick';
}
