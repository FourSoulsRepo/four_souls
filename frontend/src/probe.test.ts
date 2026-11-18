import { describe, expect, it } from 'vitest';
import { frameStats, playable } from './probe';

describe('frameStats', () => {
  // ADR 035 item 3.7: frames, seconds, and frames over 34 ms (diag.SlowMs).
  it('counts the frames between the first and last time', () => {
    expect(frameStats([0, 16, 32, 82, 100])).toEqual({ frames: 4, seconds: 0.1, slow: 1, worstMs: 50 });
  });

  it('gives zeros without two frames', () => {
    expect(frameStats([])).toEqual({ frames: 0, seconds: 0, slow: 0, worstMs: 0 });
    expect(frameStats([5])).toEqual({ frames: 0, seconds: 0, slow: 0, worstMs: 0 });
  });
});

describe('playable', () => {
  // canPlayType answers '', 'maybe' or 'probably' (HTML standard).
  it('keeps the types the webview may play', () => {
    expect(playable((type) => (type === 'audio/mpeg' ? 'maybe' : ''))).toEqual(['mp3']);
  });
});
