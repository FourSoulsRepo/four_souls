import { afterEach, describe, expect, it } from 'vitest';
import { px, rem, scaleOf, setUiScale, uiScale } from './uiScale';

// The interface scale (ST-10, 18.18).
describe('uiScale', () => {
  afterEach(() => {
    setUiScale(100);
  });

  it('reads a percent out of bounds as 100', () => {
    expect(scaleOf(150)).toBe(150);
    expect(scaleOf(75)).toBe(75);
    expect(scaleOf(200)).toBe(200);
    expect(scaleOf(0)).toBe(100);
    expect(scaleOf(74)).toBe(100);
    expect(scaleOf(201)).toBe(100);
    expect(scaleOf(Number.NaN)).toBe(100);
  });

  it('scales px numbers by the setting', () => {
    expect(px(12)).toBe(12);
    setUiScale(150);
    expect(uiScale()).toBe(1.5);
    expect(px(12)).toBe(18);
  });

  it('writes CSS lengths in rem, 16 px a rem', () => {
    expect(rem(24)).toBe('1.5rem');
    expect(rem(-70)).toBe('-4.375rem');
  });
});
