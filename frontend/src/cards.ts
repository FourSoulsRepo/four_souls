// Card display data and the pure helpers the Card component uses.
// The data comes from Go (App.Cards), from card_db.

export interface CardInfo {
  ref: string;
  set: string;
  type: string; // e.g. "Eternal Treasure Card"
  name: string; // in the chosen language; the store sets it (inLanguage)
  names?: Record<string, string>; // by language, English too (LO-02)
  translators?: Record<string, string[]>; // credited by language (CD-06)
  effects: string[];
  footnotes: { title: string; text: string }[];
  stats: { name: string; value: string }[];
  rewards: { count: string; kind: string }[];
  image: string; // URL in the chosen language; the store sets it (inLanguage)
  // images holds the URL per language with an image folder; '' or none
  // when missing: the card shows as text (A-11, 16.6).
  images?: Record<string, string>;
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

// inLanguage gives the card named and pictured in lang, English if it
// has no name or image there (LO-06, 16.6).
export function inLanguage(card: CardInfo, lang: string): CardInfo {
  const names = new Map(Object.entries(card.names ?? {}));
  const name = names.get(lang) ?? names.get('en') ?? card.name;
  const images = new Map(Object.entries(card.images ?? {}));
  const image = images.get(lang) ?? images.get('en') ?? card.image;
  return name === card.name && image === card.image ? card : { ...card, name, image };
}

// translatorsOf credits who translated the card into lang (CD-06): none
// for English or a card with no name in lang.
export function translatorsOf(card: CardInfo, lang: string): string[] {
  if (lang === 'en' || !Object.keys(card.names ?? {}).includes(lang)) {
    return [];
  }
  return new Map(Object.entries(card.translators ?? {})).get(lang) ?? [];
}

// unknownCard stands in for a ref the app has no data for.
export function unknownCard(ref: string): CardInfo {
  return { ref, set: '', type: '', name: ref, effects: [], footnotes: [], stats: [], rewards: [], image: '' };
}

// AbilityLine is one ability's words split like a line of card text:
// its cost icon, the cost, and what it does (15.24).
export interface AbilityLine {
  // icon: 'tap' for a ↷ ability, 'paid' for one with another cost.
  icon?: 'tap' | 'paid';
  cost?: string;
  effect: string;
}

// abilityLine splits "↷: Gain 3¢." or "Pay 5¢: Deal 1 damage." at the
// first colon. Words with no cost before a colon stay whole.
export function abilityLine(text: string): AbilityLine {
  const colon = text.indexOf(': ');
  if (colon < 0) {
    return { effect: text };
  }
  const cost = text.slice(0, colon).trim();
  const effect = text.slice(colon + 2).trim();
  if (cost === '↷' || cost === '{Tap Effect}') {
    return { icon: 'tap', effect };
  }
  return { icon: 'paid', cost, effect };
}
