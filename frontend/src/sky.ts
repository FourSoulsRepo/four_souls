// The menu background (15.42, V-35), after the foursouls.ru page: a
// night sky with five see-through star layers in front, all drifting
// down and to the left. The nearer a layer, the faster it moves, so the
// stars seem deep. On the site each layer is a background tile moved a
// whole number of tiles per 120 s loop, so the loop has no seam; here a
// layer moves one tile per its own loop, which looks the same.
import type { CSSProperties } from 'react';
import { scaleFor, type Speed } from './playback';
import type { Screen } from './store';

// SkyLayer is one drifting tile: its picture, its drawn size, and the
// seconds it takes to move by one tile at the normal speed.
export interface SkyLayer {
  image: string;
  width: number;
  height: number;
  seconds: number;
}

// skyLoop is the site's loop: the sky moves one tile in it.
const skyLoop = 120;

// tilesPerLoop is how far each layer moves in one loop, in its own
// tiles: the sky one, the star layers four to twelve, far to near.
const tilesPerLoop = [1, 4, 6, 8, 10, 12];

// skyLayers places the images (the sky, then the star layers, far to
// near) as drifting tiles, far first. The sky is drawn twice its
// 600×400 size, as on the site. Animations off: still. A partial set
// gives none: half a sky would not match the site.
export function skyLayers(images: string[], speed: Speed): SkyLayer[] {
  if (images.length !== tilesPerLoop.length) {
    return [];
  }
  const scale = scaleFor(speed, 0, 0);
  return images.map((image, i) => {
    const size = i === 0 ? 2 : 1;
    const seconds = scale > 0 ? (skyLoop * scale) / (tilesPerLoop.at(i) ?? 1) : 0;
    return { image, width: 600 * size, height: 400 * size, seconds };
  });
}

// skyStyle is a layer's inline style; Sky.css draws and moves it.
export function skyStyle(l: SkyLayer): CSSProperties {
  return {
    '--tile-w': `${String(l.width)}px`,
    '--tile-h': `${String(l.height)}px`,
    backgroundImage: `url("${l.image}")`,
    animationDuration: l.seconds > 0 ? `${String(l.seconds)}s` : undefined,
    animationName: l.seconds > 0 ? undefined : 'none',
  } as CSSProperties;
}

// onSky tells the screens drawn on the sky: the home screen and the
// screens that set up a game. The splash keeps its plain notice, and the
// game its mat.
export function onSky(screen: Screen): boolean {
  return screen !== 'splash' && screen !== 'game';
}
