// Card display data and the pure helpers the Card component uses.
// The data comes from Go (App.Cards), from card_db.

export interface CardInfo {
  ref: string;
  set: string;
  type: string; // e.g. "Eternal Treasure Card"
  name: string;
  effects: string[];
  footnotes: { title: string; text: string }[];
  stats: { name: string; value: string }[];
  rewards: { count: string; kind: string }[];
  image: string; // URL; '' when missing: the card shows as text (A-11)
}

// Kind is the deck family a card belongs to; it sets the blank card's
// colour and the card back.
export type Kind = 'character' | 'treasure' | 'loot' | 'monster' | 'soul';

export function kindOf(type: string): Kind {
  if (type === 'Character Card') {
    return 'character';
  }
  if (type === 'Bonus Soul Card') {
    return 'soul';
  }
  if (type.includes('Treasure')) {
    return 'treasure';
  }
  if (/Monster|Boss|Event|Curse Card/.test(type)) {
    return 'monster';
  }
  return 'loot';
}

// TextPart is plain text or an icon written as {Icon Name} in card text.
export type TextPart = { text: string } | { icon: string };

export function textParts(s: string): TextPart[] {
  const parts: TextPart[] = [];
  let rest = s;
  for (;;) {
    const open = rest.indexOf('{');
    const close = rest.indexOf('}', open);
    if (open < 0 || close < 0) {
      break;
    }
    if (open > 0) {
      parts.push({ text: rest.slice(0, open) });
    }
    parts.push({ icon: rest.slice(open + 1, close) });
    rest = rest.slice(close + 1);
  }
  if (rest !== '') {
    parts.push({ text: rest });
  }
  return parts;
}

// iconText is the short label of a text icon on a blank card.
export function iconText(icon: string): string {
  switch (icon) {
    case 'Tap Effect':
      return '↷';
    case 'Paid Effect':
      return '$';
    case 'Curse Effect':
      return 'Curse';
    default:
      return icon.replace(/ Effect$/, '');
  }
}

// unknownCard stands in for a ref the app has no data for.
export function unknownCard(ref: string): CardInfo {
  return { ref, set: '', type: '', name: ref, effects: [], footnotes: [], stats: [], rewards: [], image: '' };
}
