import { describe, expect, it } from 'vitest';
import { botLook, botPaths } from './BotIcon';

// The robot before a bot's name follows its seat's level (ADR 033,
// Decision 10.3).
describe('bot robots', () => {
  it('draws each level its own robot', () => {
    expect(botLook('standard')).toBe('standard');
    expect(botLook('hard')).toBe('hard');
    expect(botPaths.hard).not.toBe(botPaths.standard);
  });

  it("draws a newer server's level as standard", () => {
    expect(botLook('ultra')).toBe('standard');
    expect(botLook('')).toBe('standard');
  });
});
