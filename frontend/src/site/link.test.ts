import { describe, expect, it } from 'vitest';
import { compact, decodePayload, encodePayload, fromBase64url, toBase64url } from './link';

const payload = sample();

// sample is a payload like internal/website/testdata, with non-ASCII text.
function sample(): string {
  return JSON.stringify(
    {
      title: 'Может ли игрок сыграть лут? ¢ — "quotes"',
      sets: ['b2'],
      setup: { players: [{ character: 'isaac', hand: ['a_penny'] }, { character: 'maggy' }], active: 0 },
      action: { player: 1, kind: 'play', card: 'a_penny' },
      expect: { allowed: false, rule: 'R-CARD-08' },
    },
    null,
    2,
  );
}

describe('check links (ADR 026, 7)', () => {
  it('decode gives back what encode got', async () => {
    const data = await encodePayload(payload);
    expect(data).toMatch(/^[A-Za-z0-9_-]+$/);
    expect(await decodePayload(data)).toBe(payload);
  });

  it('compresses: a four-player table is shorter than its JSON', async () => {
    const player = (character: string) => ({ character, hand: ['a_penny', 'bomb', 'the_tower'], items: ['the_d6'], cents: 3 });
    const json = JSON.stringify({ sets: ['b2'], setup: { players: ['isaac', 'maggy', 'cain', 'judas'].map(player), active: 0 } });
    expect((await encodePayload(json)).length).toBeLessThan(json.length * 0.75);
  });

  it('compact drops spaces of JSON and keeps other text', () => {
    expect(compact('{ "a": [1, 2] }')).toBe('{"a":[1,2]}');
    expect(compact('not json')).toBe('not json');
  });

  it('base64url has no padding and round-trips every byte', () => {
    const bytes = Uint8Array.from({ length: 256 }, (_, i) => i);
    const text = toBase64url(bytes);
    expect(text).not.toMatch(/[+/=]/);
    expect(fromBase64url(text)).toEqual(bytes);
    expect(toBase64url(Uint8Array.of(0xfb, 0xff))).toBe('-_8');
  });

  it('a cut or foreign link fails to decode', async () => {
    const data = await encodePayload(payload);
    await expect(decodePayload(data.slice(0, data.length / 2))).rejects.toThrow();
    await expect(decodePayload('a+b/')).rejects.toThrow();
  });
});
