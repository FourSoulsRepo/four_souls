import { describe, expect, it } from 'vitest';
import { Purpose, type Choice } from './actions';
import { aimClick, arrow, blocks, cancelOption, promptAim, targetAt } from './aim';
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

describe('aimClick', () => {
  const targets = [monster, ann];

  it('picks a target', () => {
    expect(aimClick(['obj-40', 'board'], 7, targets)).toEqual({ kind: 'pick', target: monster });
  });

  it('cancels on empty space', () => {
    expect(aimClick([], 7, targets)).toEqual({ kind: 'cancel' });
  });

  it('does nothing on the source, even over a player it may target', () => {
    expect(aimClick(['obj-7', 'hand-0', 'player-0'], 7, targets)).toEqual({ kind: 'none' });
  });

  it('calls any other card, player, pile, or stack item wrong', () => {
    expect(aimClick(['obj-8', 'obj-7', 'hand-0'], 7, targets)).toEqual({ kind: 'wrong' });
    expect(aimClick(['player-1'], 7, targets)).toEqual({ kind: 'wrong' });
    expect(aimClick(['deck-monster'], 7, targets)).toEqual({ kind: 'wrong' });
    expect(aimClick(['stack-3'], 7, targets)).toEqual({ kind: 'wrong' });
  });

  it('lets a source that is its own target pick itself', () => {
    expect(aimClick(['obj-40'], 40, targets)).toEqual({ kind: 'pick', target: monster });
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

describe('promptAim', () => {
  const choose = (choice: number) => ({ player: 0, kind: 4, choice });
  const target: Choice = {
    title: 'Choose a target',
    purpose: Purpose.target,
    options: [
      { label: 'Ann', intent: choose(0), target: ann },
      { label: 'cancel', intent: choose(1) },
    ],
  };

  it('aims a placed card at its target prompt', () => {
    expect(promptAim(7, target)).toEqual({ source: 7, choice: target, targets: [ann] });
  });

  it('aims from object 0 too', () => {
    expect(promptAim(0, target)?.source).toBe(0);
  });

  it('keeps the dialog when the card is unknown or nothing aims', () => {
    expect(promptAim(undefined, target)).toBeNull();
    expect(promptAim(7, null)).toBeNull();
    expect(promptAim(7, { ...target, purpose: Purpose.answer })).toBeNull();
    expect(promptAim(7, { ...target, options: target.options.slice(1) })).toBeNull();
  });

  it('cancels by the option that points at nothing', () => {
    expect(cancelOption(target)).toEqual(choose(1));
  });
});
