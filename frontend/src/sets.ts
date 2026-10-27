// Card sets: display names and the order the host screen lists them
// (oldest first, ADR 010).

import { t } from './i18n';

// baseGame is in every game; expansions need it (CD-02).
export const baseGame = 'b2';

// names are the sets' names in the chosen language.
const names = new Map([
  ['b2', () => t('set.base_game')],
  ['g2', () => t('set.gold_box')],
  ['fsp2', () => t('set.four_souls_plus')],
  ['r', () => t('set.requiem')],
]);

const order = ['b2', 'g2', 'fsp2', 'r'];

// setName names a set code; unknown codes stay as they are.
export function setName(code: string): string {
  return names.get(code)?.() ?? code;
}

// sortSets puts known sets in release order, then the rest by code.
export function sortSets(codes: string[]): string[] {
  const rank = (c: string) => {
    const i = order.indexOf(c);
    return i < 0 ? order.length : i;
  };
  return [...codes].sort((a, b) => rank(a) - rank(b) || a.localeCompare(b));
}

// setList names sets for a game list: "Base Game, Gold Box".
export function setList(codes: string[] | null): string {
  return sortSets(codes ?? []).map(setName).join(', ');
}
