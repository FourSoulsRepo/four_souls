// Game mats (ST-03, V-24). Color mats are CSS; image mats are files in
// the private images, listed by Go. An image mat may share a color mat's
// id: without the image it shows those colors.
import type { CSSProperties } from 'react';
import type { MatInfo } from './bridge';

export const colorMats: readonly MatInfo[] = [
  { id: 'basement', name: 'Basement', image: '' },
  { id: 'cathedral', name: 'Cathedral', image: '' },
];

// matChoices lists the color mats, then the image mats; an image mat
// takes the place of the color mat with its id.
export function matChoices(images: readonly MatInfo[]): MatInfo[] {
  const out = colorMats.map((c) => images.find((m) => m.id === c.id) ?? c);
  return [...out, ...images.filter((m) => !colorMats.some((c) => c.id === m.id))];
}

// droppedMat was the Green felt color mat; it read as a poker table
// (H-05), so it is gone (ST-08) and a saved one is the default.
const droppedMat = 'felt';

// matOf resolves the mat setting: '' is the default, the first mat
// image, or Basement without images.
export function matOf(id: string, images: readonly MatInfo[]): string {
  return (id !== droppedMat && id) || (images[0]?.id ?? 'basement');
}

export interface MatLook {
  className: string;
  style: CSSProperties;
}

// matLook draws a mat: its colors, or Basement for a mat with no colors
// of its own, and its image over them. A missing or unloadable image
// leaves the colors.
export function matLook(setting: string, images: readonly MatInfo[]): MatLook {
  const id = matOf(setting, images);
  const colors = colorMats.some((c) => c.id === id) ? id : 'basement';
  const image = images.find((m) => m.id === id)?.image ?? '';
  if (!image) {
    return { className: `mat-${colors}`, style: {} };
  }
  return { className: `mat-${colors} mat-image`, style: { '--mat-image': `url("${image}")` } as CSSProperties };
}
