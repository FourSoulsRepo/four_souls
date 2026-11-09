// The interface scale (ST-10, 18.18). The CSS sizes everything in rem,
// so the root font size scales the whole window. Pixel numbers in code
// (feel.ts, card sizes) are px at 100 %; px() turns one into screen px
// at the current scale.

// The setting's bounds and default, in percent, as internal/settings.
export const minScale = 75;
export const maxScale = 200;
export const defaultScale = 100;

let factor = 1;

// scaleOf reads a saved percent; anything out of bounds is the default.
export function scaleOf(percent: number): number {
  return Number.isFinite(percent) && percent >= minScale && percent <= maxScale ? percent : defaultScale;
}

// setUiScale sets the root font size to percent of the browser's 16 px.
export function setUiScale(percent: number): void {
  const p = scaleOf(percent);
  factor = p / 100;
  if (typeof document !== 'undefined') {
    document.documentElement.style.fontSize = p === defaultScale ? '' : `${String(p)}%`;
  }
}

// uiScale is the current scale as a factor: 1 at 100 %.
export function uiScale(): number {
  return factor;
}

// px is n px at 100 %, in screen px at the current scale.
export function px(n: number): number {
  return n * factor;
}

// rem is n px at 100 % as a CSS length, which the browser scales.
export function rem(n: number): string {
  return `${String(n / 16)}rem`;
}
