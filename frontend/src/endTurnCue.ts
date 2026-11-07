// End turn's cue (20.10, ADR 027 "End Turn button"): once nothing is
// left to do, a short voice line plays feel.endTurn.cue seconds later.
// Each turn tries it once: the first time the button turns "done". If
// by then it is no longer done (a play became possible, the turn was
// ended or another player's turn began), the cue is skipped, and that
// turn has no other.
import { feel } from './feel';

export interface Clock {
  set: (fn: () => void, ms: number) => number;
  clear: (id: number) => void;
}

const windowClock: Clock = {
  set: (fn, ms) => window.setTimeout(fn, ms),
  clear: (id) => {
    window.clearTimeout(id);
  },
};

export interface EndTurnCue {
  // update tells the button's state: the turn's number and if it is done.
  update: (turn: number, done: boolean) => void;
  // stop drops a waiting cue, as the button goes away.
  stop: () => void;
}

// createEndTurnCue calls play for the cue; clock is for tests.
export function createEndTurnCue(play: () => void, clock: Clock = windowClock): EndTurnCue {
  let tried: number | undefined;
  let now = { turn: -1, done: false };
  let timer: number | undefined;
  return {
    update: (turn, done) => {
      now = { turn, done };
      if (!done || tried === turn) {
        return;
      }
      tried = turn;
      if (timer !== undefined) {
        clock.clear(timer);
      }
      timer = clock.set(() => {
        timer = undefined;
        if (now.done && now.turn === turn) {
          play();
        }
      }, feel.endTurn.cue * 1000);
    },
    stop: () => {
      if (timer !== undefined) {
        clock.clear(timer);
        timer = undefined;
      }
    },
  };
}
