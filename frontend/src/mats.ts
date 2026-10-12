// Game mats (ST-03, V-24). Color mats are CSS; image mats are files in
// the private images, listed by Go. An image mat may share a color mat's
// id: without the image it shows those colors.
import type { CSSProperties } from 'react';
import type { MatInfo } from './bridge';

export const colorMats: readonly MatInfo[] = [
  { id: 'felt', name: 'Green felt', image: '' },
  { id: 'basement', name: 'Basement', image: '' },
  { id: 'cathedral', name: 'Cathedral', image: '' },
];

// matChoices lists the color mats, then the image mats; an image mat
// takes the place of the color mat with its id.
export function matChoices(images: readonly MatInfo[]): MatInfo[] {
  const out = colorMats.map((c) => images.find((m) => m.id === c.id) ?? c);
  return [...out, ...images.filter((m) => !colorMats.some((c) => c.id === m.id))];
}

export interface MatLook {
  className: string;
  style: CSSProperties;
}

// matLook draws a mat: its colors, or felt for a mat with no colors of
// its own, and its image over them. A missing or unloadable image
// leaves the colors.
export function matLook(id: string, images: readonly MatInfo[]): MatLook {
  const colors = colorMats.some((c) => c.id === id) ? id : 'felt';
  const image = images.find((m) => m.id === id)?.image ?? '';
  if (!image) {
    return { className: `mat-${colors}`, style: {} };
  }
  return { className: `mat-${colors} mat-image`, style: { '--mat-image': `url("${image}")` } as CSSProperties };
}
