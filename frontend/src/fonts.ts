// The game's own fonts (15.20), from the private images: font-souls for
// Latin text, and the Russian site's card font for Cyrillic. CSS asks for
// them through --font-game; without the images neither family exists,
// and the stack falls back to the system font.
import type { Fonts } from './bridge';

export const families = { title: 'Four Souls Title', cyrillic: 'Four Souls Cyrillic' } as const;

// loadFonts registers the fonts that have a URL. A face loads when text
// first uses it; one that fails to load is skipped by the stack.
export function loadFonts(fonts: Fonts): void {
  if (typeof document === 'undefined' || typeof FontFace === 'undefined') {
    return;
  }
  for (const [family, url] of [
    [families.title, fonts.title],
    [families.cyrillic, fonts.cyrillic],
  ] as const) {
    if (url !== '') {
      document.fonts.add(new FontFace(family, `url("${url}")`, { display: 'swap' }));
    }
  }
}
