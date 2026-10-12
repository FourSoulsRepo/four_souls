import { describe, expect, it } from 'vitest';
import { matChoices, matLook } from './mats';

const basement = { id: 'basement', name: 'Basement', image: '/cardimg/assets/mats/basement.webp' };
const room = { id: 'dark_room', name: 'Dark room', image: '/cardimg/assets/mats/dark_room.webp' };

describe('matChoices', () => {
  it('lists the color mats without images', () => {
    expect(matChoices([]).map((m) => m.id)).toEqual(['felt', 'basement', 'cathedral']);
  });

  it('puts an image in place of its color mat and adds new ones', () => {
    const got = matChoices([room, basement]);
    expect(got.map((m) => m.id)).toEqual(['felt', 'basement', 'cathedral', 'dark_room']);
    expect(got[1]).toBe(basement);
  });
});

describe('matLook', () => {
  it('draws an image over its colors', () => {
    expect(matLook('basement', [basement])).toEqual({
      className: 'mat-basement mat-image',
      style: { '--mat-image': 'url("/cardimg/assets/mats/basement.webp")' },
    });
    expect(matLook('dark_room', [room]).className).toBe('mat-felt mat-image');
  });

  it('falls back to colors when the image is missing', () => {
    expect(matLook('basement', [])).toEqual({ className: 'mat-basement', style: {} });
    expect(matLook('dark_room', [])).toEqual({ className: 'mat-felt', style: {} });
  });
});
