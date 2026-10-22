import { describe, expect, it } from 'vitest';
import { onSky, skyLayers, skyStyle } from './sky';

const images = ['sky', 'stars-1', 'stars-2', 'stars-3', 'stars-4', 'stars-5'];

describe('sky', () => {
  it('drifts the nearer stars faster, as on foursouls.ru', () => {
    const layers = skyLayers(images, 'normal');
    expect(layers.map((l) => l.seconds)).toEqual([120, 30, 20, 15, 12, 10]);
    expect(layers[0]).toMatchObject({ image: 'sky', width: 1200, height: 800 });
    expect(layers[5]).toMatchObject({ image: 'stars-5', width: 600, height: 400 });
  });

  it('follows the animation speed, and stands still when it is off', () => {
    expect(skyLayers(images, 'slow')[0]?.seconds).toBe(192);
    expect(skyLayers(images, 'fast')[0]?.seconds).toBe(60);
    const still = skyLayers(images, 'off');
    expect(still.map((l) => l.seconds)).toEqual([0, 0, 0, 0, 0, 0]);
    expect(still[0] && skyStyle(still[0]).animationName).toBe('none');
  });

  it('shows nothing without every layer', () => {
    expect(skyLayers([], 'normal')).toEqual([]);
    expect(skyLayers(images.slice(1), 'normal')).toEqual([]);
  });

  it('is behind the menus, not the splash or the game', () => {
    expect(onSky('menu')).toBe(true);
    expect(onSky('host')).toBe(true);
    expect(onSky('table')).toBe(true);
    expect(onSky('splash')).toBe(false);
    expect(onSky('game')).toBe(false);
  });
});
