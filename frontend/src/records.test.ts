import { describe, expect, it } from 'vitest';
import type { RecordEntry } from './bridge';
import { dateText, entryLine, recordTitle, resultText, whyText } from './records';

const entry: RecordEntry = {
  path: '/r/g1.fsrec',
  name: 'g1.fsrec',
  started: Date.UTC(2026, 10, 16, 12, 0),
  players: ['Ann', 'Bo'],
  finished: true,
  winners: ['Bo'],
  steps: 840,
  challenge: '',
  why: '',
};

// The Records screen (M-01): its lines for each kind of record.
describe('records', () => {
  it('says why a record cannot be watched', () => {
    expect(whyText('')).toBe('');
    for (const why of ['before_replays', 'newer', 'unreadable'] as const) {
      expect(whyText(why)).not.toBe('');
    }
    expect(entryLine({ ...entry, why: 'newer' })).toBe(whyText('newer'));
  });

  it('names a record by its start and players, else its file', () => {
    expect(dateText(0)).toBe('');
    expect(recordTitle(entry)).toContain('Ann, Bo');
    expect(recordTitle(entry)).toContain(dateText(entry.started));
    expect(recordTitle({ ...entry, started: 0, players: [] })).toBe('g1.fsrec');
  });

  it('says how a game ended', () => {
    expect(resultText(true, ['Bo'])).toBe('Won by Bo');
    expect(resultText(true, [])).toBe('Finished');
    expect(resultText(false, [])).toBe('Unfinished');
    expect(entryLine(entry)).toBe('Won by Bo · 840 steps');
    expect(entryLine({ ...entry, finished: false, steps: 1 })).toBe('Unfinished · 1 step');
  });
});
