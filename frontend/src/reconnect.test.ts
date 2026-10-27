import { describe, expect, it } from 'vitest';
import { endsGame, noticeText, retryDelay } from './reconnect';

describe('retryDelay (N-20)', () => {
  it('tries at once, then 1, 2, 4, 8 seconds, then every 10', () => {
    expect([0, 1, 2, 3, 4, 5, 6, 20].map(retryDelay)).toEqual([0, 1000, 2000, 4000, 8000, 10000, 10000, 10000]);
  });
});

describe('notices (N-20)', () => {
  it('the engine notice keeps the game; a lost or idle one ends it', () => {
    expect(endsGame({ kind: 'engine_updated' })).toBe(false);
    expect(endsGame({ kind: 'game_lost' })).toBe(true);
    expect(endsGame({ kind: 'game_idle', reason: 'empty' })).toBe(true);
    expect(endsGame({ kind: 'game_ended' })).toBe(true);
  });

  it('says what happened', () => {
    expect(noticeText({ kind: 'engine_updated' }).text).toContain('behave differently');
    expect(noticeText({ kind: 'game_lost' }).text).toContain('start a new one');
    expect(noticeText({ kind: 'game_idle', reason: 'long' }).text).toContain('24 hours');
    expect(noticeText({ kind: 'game_idle', reason: 'empty' }).text).toContain('Nobody played');
    expect(noticeText({ kind: 'game_ended' }).text).toContain('The host ended');
  });
});
