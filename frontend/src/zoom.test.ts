import { describe, expect, it } from 'vitest';
import { keyword, placeZoom, zoomGap, zoomMargin } from './zoom';

const view = { width: 1280, height: 720 };
const size = { width: 320, height: 450 };
const card = (left: number, top: number) => ({ left, top, right: left + 80, bottom: top + 112 });

describe('placeZoom', () => {
  it('opens right beside a card in the left half, tips to the right', () => {
    expect(placeZoom(card(300, 200), size, view, false)).toEqual({ x: 380 + zoomGap, y: 256 - 225, tips: 'right' });
  });

  it('opens left beside a card in the right half, tips to the left', () => {
    const p = placeZoom(card(900, 200), size, view, false);
    expect(p).toEqual({ x: 900 - zoomGap - 320, y: 31, tips: 'left' });
  });

  it('stays inside the window', () => {
    expect(placeZoom(card(10, 0), size, view, false).y).toBe(zoomMargin);
    expect(placeZoom(card(10, 650), size, view, false).y).toBe(720 - 450 - zoomMargin);
    expect(placeZoom(card(1000, 200), { width: 1100, height: 450 }, view, false).x).toBe(zoomMargin);
  });

  it('rises from a hand card: centered over it, on the bottom edge', () => {
    expect(placeZoom(card(600, 640), size, view, true)).toEqual({ x: 640 - 160, y: 720 - 450 - zoomMargin, tips: 'left' });
    expect(placeZoom(card(20, 640), size, view, true).x).toBe(zoomMargin);
    expect(placeZoom(card(1190, 640), size, view, true).x).toBe(1280 - 320 - zoomMargin);
  });
});

describe('keyword', () => {
  it('drops the dashes around a footnote title', () => {
    expect(keyword('-Eternal-')).toBe('Eternal');
    expect(keyword('Guppy')).toBe('Guppy');
  });
});
