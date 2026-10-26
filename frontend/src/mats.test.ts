import { describe, expect, it } from 'vitest';
import { matChoices, matLook } from './mats';

const basement = { id: 'basement', name: 'Basement', image: '/cardimg/assets/mats/basement.webp' };
const room = { id: 'dark_room', name: 'Dark room', image: '/cardimg/assets/mats/dark_room.webp' };

describe('matChoices', () => {
  it('lists the color mats without images', () => {
    expect(matChoices([]).map((m) => m.id)).toEqual(['basement', 'cathedral']);
  });

  it('puts an image in place of its color mat and adds new ones', () => {
    const got = matChoices([room, basement]);
    expect(got.map((m) => m.id)).toEqual(['basement', 'cathedral', 'dark_room']);
    expect(got[0]).toEqual(basement);
  });
});

describe('matLook', () => {
  it('draws an image over its colors', () => {
    expect(matLook('basement', [basement])).toEqual({
      className: 'mat-basement mat-image',
      style: { '--mat-image': 'url("/cardimg/assets/mats/basement.webp")' },
    });
    expect(matLook('dark_room', [room]).className).toBe('mat-basement mat-image');
  });

  it('defaults to the first image, or Basement without images', () => {
    expect(matLook('', [room, basement]).className).toBe('mat-basement mat-image');
    expect(matLook('', [basement]).className).toBe('mat-basement mat-image');
    expect(matLook('', [])).toEqual({ className: 'mat-basement', style: {} });
  });

  it('draws a saved Green felt as the default (ST-08)', () => {
    expect(matLook('felt', [room, basement]).className).toBe('mat-basement mat-image');
    expect(matLook('felt', [])).toEqual({ className: 'mat-basement', style: {} });
  });

  it('falls back to colors when the image is missing', () => {
    expect(matLook('basement', [])).toEqual({ className: 'mat-basement', style: {} });
    expect(matLook('dark_room', [])).toEqual({ className: 'mat-basement', style: {} });
  });
});
