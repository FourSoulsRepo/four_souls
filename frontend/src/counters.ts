import { t } from './i18n';

// Counter colours (18.15). A general counter is green, a gold counter
// is gold, fear is ghost white (18.26). The owner picks the colour of
// each new kind; until then a kind shows grey.
const colors = new Map<string, string>([
  ['', '#5cb83a'],
  ['gold', '#e0b43a'],
  ['fear', '#e6ecf5'],
]);

// counterColor is the fill of a counter badge of this kind; '' is the
// general kind.
export function counterColor(name = ''): string {
  return colors.get(name) ?? '#8a8f98';
}

// ghostPath is fear's shape (18.26), in a 100 × 100 box: a ghost with a
// wavy hem, so fear differs from other counters without its colour.
export const ghostPath = 'M50 6C75 6 91 25 91 49L91 91C85 84 78 84 71 92C64 84 57 84 50 92C43 84 36 84 29 92C22 84 15 84 9 91L9 49C9 25 25 6 50 6Z';

// fearOf is a player's fear counters (R-TRICK-08); 0 when none.
export function fearOf(list: { name?: string; count: number }[] | undefined): number {
  return list?.find((c) => c.name === 'fear')?.count ?? 0;
}

// PlayerCounter is one kind of counter on a player (14.31): word names
// the kind beside the number; title is its tooltip and screen-reader
// name.
export interface PlayerCounter {
  name: string;
  count: number;
  word: string;
  title: string;
}

// playerCounters lists a player's counters to show by their name: each
// kind with any, in the view's order. Fear is not among them: it shows
// on the player's character card only (owner, 18.26).
export function playerCounters(list: { name?: string; count: number }[] | undefined): PlayerCounter[] {
  return (list ?? [])
    .filter((c) => c.count > 0 && c.name !== 'fear')
    .map(({ name = '', count }) => ({ name, count, word: name, title: t('table.player_counters_title', { name, n: count }) }));
}
