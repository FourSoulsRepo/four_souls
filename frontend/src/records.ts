// The Records screen's words (M-01, ADR 031).
import type { RecordEntry, RecordWhy } from './bridge';
import { language, t } from './i18n';
import { seatName, seatNameOf } from './lobby';

// whyText says why a record can't be watched; '' when it can.
export function whyText(why: RecordWhy): string {
  switch (why) {
    case 'before_replays':
      return t('error.record_before_replays');
    case 'newer':
      return t('error.record_newer');
    case 'unreadable':
      return t('error.record_unreadable');
    case '':
      return '';
  }
}

// dateText is a record's start in the player's language; '' for none.
export function dateText(ms: number): string {
  if (ms <= 0) {
    return '';
  }
  const options: Intl.DateTimeFormatOptions = { dateStyle: 'medium', timeStyle: 'short' };
  try {
    return new Intl.DateTimeFormat(language(), options).format(ms);
  } catch {
    return new Intl.DateTimeFormat(undefined, options).format(ms);
  }
}

// recordTitle names a listed record: its start and players, a bot as
// in the game (19.216), or its file name when the header could not be
// read.
export function recordTitle(e: RecordEntry): string {
  if (e.started <= 0 || e.players.length === 0) {
    return e.name;
  }
  return t('records.title', { date: dateText(e.started), players: e.players.map(seatName).join(', ') });
}

// resultText is how a game ended: its winners, or that it never did.
export function resultText(finished: boolean, winners: string[]): string {
  if (!finished) {
    return t('records.unfinished');
  }
  return winners.length > 0 ? t('records.won', { players: winners.join(', ') }) : t('records.finished');
}

// entryLine is the line under a listed record: the reason it can't be
// watched, else how it ended and its length.
export function entryLine(e: RecordEntry): string {
  if (e.why) {
    return whyText(e.why);
  }
  const winners = e.winners.map((w) => seatNameOf(e.players, w));
  return `${resultText(e.finished, winners)} · ${t('records.steps', { n: e.steps })}`;
}
