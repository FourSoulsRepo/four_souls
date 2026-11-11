import { t } from './i18n';

// Counter colours (18.15). A general counter is green, a gold counter
// is gold. The owner picks the colour of each new kind; until then a
// kind shows grey.
const colors = new Map<string, string>([
  ['', '#5cb83a'],
  ['gold', '#e0b43a'],
]);

// counterColor is the fill of a counter badge of this kind; '' is the
// general kind.
export function counterColor(name = ''): string {
  return colors.get(name) ?? '#8a8f98';
}

// PlayerCounter is one kind of counter on a player (14.31, 14.33):
// fear in Trick or Treat. word names the kind beside the number; title
// says what the counters do, as tooltip and screen-reader name.
export interface PlayerCounter {
  name: string;
  count: number;
  word: string;
  title: string;
}

// playerCounters lists a player's counters to show by their name: each
// kind with any, in the view's order. Fear says what it does (R-TRICK-09).
export function playerCounters(list: { name?: string; count: number }[] | undefined): PlayerCounter[] {
  return (list ?? [])
    .filter((c) => c.count > 0)
    .map(({ name = '', count }) =>
      name === 'fear'
        ? { name, count, word: t('table.fear'), title: t('table.fear_title', { n: count }) }
        : { name, count, word: name, title: t('table.player_counters_title', { name, n: count }) },
    );
}
