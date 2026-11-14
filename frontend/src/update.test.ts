import { describe, expect, it } from 'vitest';
import { manualReason, megabytes, percentOf, progressText, reasonText } from './update';

// The Update screen (18.24): the download line and the reasons.
describe('update', () => {
  it('counts MB with one decimal', () => {
    expect(megabytes(0)).toBe('0.0');
    expect(megabytes(1 << 20)).toBe('1.0');
    expect(megabytes(125 * (1 << 20) + (1 << 19))).toBe('125.5');
    expect(megabytes(-5)).toBe('0.0');
  });

  it('gives a percent only with a known total', () => {
    expect(percentOf({ done: 0, total: 100 })).toBe(0);
    expect(percentOf({ done: 999, total: 1000 })).toBe(99);
    expect(percentOf({ done: 1000, total: 1000 })).toBe(100);
    expect(percentOf({ done: 2000, total: 1000 })).toBe(100);
    expect(percentOf({ done: 500, total: 0 })).toBeNull();
  });

  it('shows percent and MB, or MB only', () => {
    const mb = 1 << 20;
    expect(progressText({ done: 5 * mb, total: 20 * mb })).toBe('25% · 5.0 of 20.0 MB');
    expect(progressText({ done: 5 * mb, total: 0 })).toBe('5.0 MB');
  });

  it('says every reason in words', () => {
    for (const r of ['not_version', 'not_newer', 'no_archive', 'not_writable', 'other'] as const) {
      expect(reasonText(r)).not.toBe('');
    }
    expect(reasonText('')).toBe('');
  });

  it('offers the releases page unless the app is up to date', () => {
    expect(manualReason('not_newer')).toBe(false);
    expect(manualReason('')).toBe(false);
    expect(manualReason('not_writable')).toBe(true);
    expect(manualReason('other')).toBe(true);
  });
});
