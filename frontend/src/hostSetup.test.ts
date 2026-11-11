import { describe, expect, it } from 'vitest';
import { createOptions, pickChallenge, type ChallengeChoice } from './hostSetup';
import type { Options } from './protocol';

const krampus: ChallengeChoice = { id: 'krampus', difficulties: ['normal', 'hard', 'ultra'] };
const detailed: Options = { picking: 'random', draft_size: 3, response_timer: 90, no_bonus_souls: true, room_deck: true };

describe('pickChallenge', () => {
  it('keeps the remembered challenge and difficulty', () => {
    expect(pickChallenge([krampus], 'krampus', 'hard')).toEqual({ challenge: 'krampus', difficulty: 'hard' });
  });
  it('falls back to the first challenge at its easiest difficulty', () => {
    expect(pickChallenge([krampus], '', '')).toEqual({ challenge: 'krampus', difficulty: 'normal' });
    expect(pickChallenge([krampus], 'gone', 'hard')).toEqual({ challenge: 'krampus', difficulty: 'hard' });
    expect(pickChallenge([krampus], 'krampus', 'insane')).toEqual({ challenge: 'krampus', difficulty: 'normal' });
  });
  it('is null when the app knows no challenge', () => {
    expect(pickChallenge([], 'krampus', 'hard')).toBeNull();
    expect(pickChallenge([{ id: 'x', difficulties: [] }], 'x', '')).toBeNull();
  });
});

describe('createOptions', () => {
  it('sends no options in a simple regular game (GS-10)', () => {
    expect(createOptions(false, detailed, true, null)).toBeUndefined();
  });
  it('sends the detailed options in a regular game; rooms only with rooms', () => {
    expect(createOptions(true, detailed, true, null)).toEqual(detailed);
    expect(createOptions(true, detailed, false, null)).toEqual({ ...detailed, room_deck: false });
  });
  it('sends a challenge from the simple setup too', () => {
    expect(createOptions(false, detailed, true, { challenge: 'krampus', difficulty: 'hard' })).toEqual({ challenge: 'krampus', difficulty: 'hard' });
  });
  it('sends a challenge with the detailed options, without bonus souls', () => {
    const got = createOptions(true, detailed, true, { challenge: 'krampus', difficulty: 'ultra' });
    expect(got).toEqual({ picking: 'random', draft_size: 3, response_timer: 90, room_deck: true, challenge: 'krampus', difficulty: 'ultra' });
    expect(got).not.toHaveProperty('no_bonus_souls');
  });
});
