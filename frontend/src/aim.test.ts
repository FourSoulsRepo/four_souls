import { describe, expect, it } from 'vitest';
import { arrow, blocks, targetAt } from './aim';
import type { OptionView } from './protocol';

const monster: OptionView = { object: 40, player: -1 };
const ann: OptionView = { player: 0, is_player: true };
const roll: OptionView = { stack_id: 3, player: 1 };

describe('targetAt', () => {
  const targets = [monster, ann, roll];

  it('finds a card, a player, or a stack item', () => {
    expect(targetAt(['obj-40', 'board'], targets)).toBe(monster);
    expect(targetAt(['obj-12', 'player-0'], targets)).toBe(ann);
    expect(targetAt(['stack-3'], targets)).toBe(roll);
  });

  it('prefers the topmost zone', () => {
    expect(targetAt(['obj-40', 'player-0'], targets)).toBe(monster);
  });

  it('finds nothing off the targets', () => {
    expect(targetAt(['obj-41', 'player-1', 'deck-loot'], targets)).toBeUndefined();
    expect(targetAt([], targets)).toBeUndefined();
  });

  it('never mistakes a player for an object with the same number', () => {
    expect(targetAt(['obj-0'], [ann])).toBeUndefined();
    expect(targetAt(['player-40'], [monster])).toBeUndefined();
  });
});

describe('arrow', () => {
  it('measures length and angle', () => {
    expect(arrow({ x: 0, y: 0 }, { x: 3, y: 4 }).length).toBe(5);
    expect(arrow({ x: 0, y: 0 }, { x: 0, y: -10 }).angle).toBe(-90);
  });
});

describe('blocks', () => {
  it('fills the chain short of the head', () => {
    expect(blocks(100, 30, 40)).toBe(2);
    expect(blocks(130, 30, 40)).toBe(3);
    expect(blocks(20, 30, 40)).toBe(0);
  });
});
