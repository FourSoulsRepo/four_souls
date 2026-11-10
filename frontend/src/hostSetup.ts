// The host screen's remembered setup (UC-03), kept in the settings.
import type { Difficulty, Options } from './protocol';

// HostSetup is a remembered host screen; seats 0 and sets null mean
// "not set yet". The password is never kept (N-14).
export interface HostSetup {
  seats: number;
  sets: string[] | null;
  detailed: boolean;
  response_timer: number;
  no_bonus_souls: boolean;
  no_deck_ratios: boolean;
  // room_deck plays with the room deck (R-SETUP-16); off by default.
  room_deck: boolean;
  // random_seats seats the players in a random order at the start.
  random_seats: boolean;
  // rollback is the rollback mode (ADR 020); '': Disabled.
  rollback: string;
  // picking is how characters are dealt (GS-01): '' one random each,
  // 'draft' a draft.
  picking: string;
  // draft_size is the characters a draft offers; 0: the default.
  draft_size: number;
  // mode is the game mode (ADR 028): '' a regular game, 'challenge' a
  // co-op challenge. An older settings file loads a regular game.
  mode: string;
  // challenge and difficulty are the last challenge and its difficulty,
  // e.g. 'krampus' and 'hard'; '': not set yet.
  challenge: string;
  difficulty: string;
}

// noHostSetup is the setup before the first host.
export const noHostSetup: HostSetup = {
  seats: 0,
  sets: null,
  detailed: false,
  response_timer: 0,
  no_bonus_souls: false,
  no_deck_ratios: false,
  room_deck: false,
  random_seats: false,
  rollback: '',
  picking: '',
  draft_size: 0,
  mode: '',
  challenge: '',
  difficulty: '',
};

// ChallengeChoice is a challenge the host screen offers, with its
// difficulties, easiest first (the bridge's ChallengeInfo).
export interface ChallengeChoice {
  id: string;
  difficulties: Difficulty[];
}

// ChallengePick is the chosen challenge and difficulty.
export interface ChallengePick {
  challenge: string;
  difficulty: Difficulty;
}

// pickChallenge picks the remembered challenge and difficulty when the
// app still knows them, else the first challenge and its easiest
// difficulty; null when the app knows none.
export function pickChallenge(list: ChallengeChoice[], challenge: string, difficulty: string): ChallengePick | null {
  const c = list.find((x) => x.id === challenge) ?? list.at(0);
  if (!c) {
    return null;
  }
  const d = c.difficulties.find((x) => x === difficulty) ?? c.difficulties.at(0);
  return d ? { challenge: c.id, difficulty: d } : null;
}

// challengeSeats are a challenge's seats: exactly 2 (ADR 028, R-COOP).
export const challengeSeats = 2;

// createOptions are the create message's options (GS-10, ADR 028).
// The simple setup sends none, unless it is a challenge: the challenge
// and difficulty always go. A challenge sends no bonus souls choice:
// the server plays it without them. The room deck goes only with rooms
// in the sets (R-SETUP-16).
export function createOptions(detailed: boolean, options: Options, rooms: boolean, challenge: ChallengePick | null): Options | undefined {
  const sent: Options | undefined = detailed ? { ...options, room_deck: rooms && !!options.room_deck } : undefined;
  if (!challenge) {
    return sent;
  }
  const out: Options = { ...sent, challenge: challenge.challenge, difficulty: challenge.difficulty };
  delete out.no_bonus_souls;
  return out;
}
