// Where the zoomed card shows (15.27): right beside the hovered card
// (F-16, T-27), or, for a hand card, rising from it (H-13). Pure, so
// tests can check it without a window.

// Box is a rectangle in window pixels.
export interface Box {
  left: number;
  top: number;
  right: number;
  bottom: number;
}

export interface Size {
  width: number;
  height: number;
}

export interface Place {
  x: number;
  y: number;
  // tips is the side of the card its keyword tips stand on: away from
  // the hovered card, toward the window's middle.
  tips: 'left' | 'right';
}

// Space kept between the zoom and the card, and to the window's edges.
export const zoomGap = 12;
export const zoomMargin = 8;

// boxOf is an element's rectangle as a plain box.
export function boxOf(el: Element): Box {
  const r = el.getBoundingClientRect();
  return { left: r.left, top: r.top, right: r.right, bottom: r.bottom };
}

function clamp(v: number, lo: number, hi: number): number {
  return Math.max(lo, Math.min(v, Math.max(lo, hi)));
}

// placeZoom puts a zoom of the given size by the anchor, the hovered
// card, inside the view. Beside: on the side toward the window's middle,
// centered on the card's height. Rise: centered over the card, standing
// on the window's bottom edge.
export function placeZoom(anchor: Box, size: Size, view: Size, rise: boolean): Place {
  const middle = (anchor.left + anchor.right) / 2;
  const tips = middle < view.width / 2 ? 'right' : 'left';
  const maxX = view.width - size.width - zoomMargin;
  const maxY = view.height - size.height - zoomMargin;
  if (rise) {
    return { x: clamp(middle - size.width / 2, zoomMargin, maxX), y: clamp(maxY, zoomMargin, maxY), tips };
  }
  const x = tips === 'right' ? anchor.right + zoomGap : anchor.left - zoomGap - size.width;
  const y = (anchor.top + anchor.bottom) / 2 - size.height / 2;
  return { x: clamp(x, zoomMargin, maxX), y: clamp(y, zoomMargin, maxY), tips };
}

// keyword is a footnote title without its dashes: "-Eternal-" → "Eternal".
export function keyword(title: string): string {
  return title.replace(/^-+|-+$/g, '').trim();
}

// RestWatch opens a card's zoom while the own arrow aims (20.23, ADR 027
// "Target arrow"): only once the pointer rests still on the card for a
// while. move tells it where the pointer is; a place that differs from
// the last restarts the wait. stop ends the watch.
export interface RestWatch {
  move: (x: number, y: number) => void;
  stop: () => void;
}

// restWatch calls open once the pointer stays at one place for wait ms
// after its last move. Each move starts the wait again; a move to the
// same place, as a browser may send, does not.
export function restWatch(wait: number, open: () => void): RestWatch {
  let timer: ReturnType<typeof setTimeout> | undefined;
  let at: { x: number; y: number } | null = null;
  return {
    move: (x, y) => {
      if (at?.x === x && at.y === y) {
        return;
      }
      at = { x, y };
      clearTimeout(timer);
      timer = setTimeout(() => {
        timer = undefined;
        open();
      }, wait);
    },
    stop: () => {
      clearTimeout(timer);
      timer = undefined;
      at = null;
    },
  };
}
