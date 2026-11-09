import { describe, expect, it } from 'vitest';
import { PromptKind, Purpose } from './actions';
import { backDelays, choiceRow, choosers, confirmFirst, dealDelays, dealPose, keptZone, pickToggle, startChoice, startWaits } from './choice';
import { feel } from './feel';
import { t } from './i18n';
import type { GameView } from './protocol';

describe('choiceRow', () => {
  it('keeps one row up to six, smaller from five on (20.18)', () => {
    expect(choiceRow(1)).toEqual({ scale: 1, gap: 0.25, wrap: false });
    expect(choiceRow(3)).toEqual({ scale: 1, gap: 0.25, wrap: false });
    expect(choiceRow(4)).toEqual({ scale: 1, gap: 0.2, wrap: false });
    expect(choiceRow(5)).toEqual({ scale: 0.85, gap: 0.125, wrap: false });
    expect(choiceRow(6)).toEqual({ scale: 0.7, gap: 0.1, wrap: false });
    expect(choiceRow(9)).toEqual({ scale: 0.7, gap: 0.1, wrap: true });
  });
});

describe('dealDelays', () => {
  it('deals the first at once, the second after 0.19 s, then 0.09 s apart', () => {
    const d = dealDelays(5);
    expect(d[0]).toBe(0);
    expect(d[1]).toBeCloseTo(0.19);
    expect(d[2]).toBeCloseTo(0.28);
    expect(d[4]).toBeCloseTo(0.46);
    expect(dealDelays(0)).toEqual([]);
  });

  it('flies back the first at once, the rest after it', () => {
    const d = backDelays(3);
    expect(d[0]).toBe(0);
    expect(d[1]).toBeCloseTo(feel.choice.backFirst);
    expect(d[2]).toBeCloseTo(feel.choice.backFirst + feel.choice.backGap);
  });
});

describe('dealPose', () => {
  const from = { x: 300, y: -100 };
  it('starts on the deck, face down, at its size', () => {
    expect(dealPose(0, from, 0.4)).toEqual({ x: 300, y: -100, scale: 0.4, turn: 180 });
  });

  it('rises above the deck early on', () => {
    expect(dealPose(0.02, from, 0.4).y).toBeLessThan(-100);
  });

  it('ends at its place, face up, full size', () => {
    const p = dealPose(feel.choice.deal, from, 0.4);
    expect(p.x).toBeCloseTo(0);
    expect(p.y).toBeCloseTo(0);
    expect(p.scale).toBeCloseTo(1);
    expect(p.turn).toBeCloseTo(0);
  });
});

describe('start-of-game choices (20.22)', () => {
  const start = { kind: PromptKind.choose, player: 0, purpose: Purpose.startingItem, players: [0, 2] };
  const view = { waiting: start } as unknown as GameView;
  const name = (p: number) => ['Ann', 'Bo', 'Cy'].at(p) ?? '?';

  it('asks Confirm only for a start-of-game pick', () => {
    expect(confirmFirst(Purpose.startingItem)).toBe(true);
    expect(confirmFirst(Purpose.target)).toBe(false);
    expect(startChoice(start)).toBe(true);
    expect(startChoice({ ...start, kind: PromptKind.order })).toBe(false);
  });

  it('picks, switches and unpicks a card, each with its sound', () => {
    expect(pickToggle(null, 1)).toEqual({ picked: 1, sound: 'choicePick' });
    expect(pickToggle(1, 2)).toEqual({ picked: 2, sound: 'choicePick' });
    expect(pickToggle(2, 2)).toEqual({ picked: null, sound: 'choiceUnpick' });
  });

  it('names every chooser for the others, none for a chooser', () => {
    expect(choosers(start)).toEqual([0, 2]);
    expect(choosers({ kind: PromptKind.choose, player: 1 })).toEqual([1]);
    expect(startWaits(view, 1, name)).toEqual([t('start.choosing', { player: 'Ann' }), t('start.choosing', { player: 'Cy' })]);
    expect(startWaits(view, -1, name)).toHaveLength(2);
    expect(startWaits(view, 2, name)).toEqual([]);
    expect(startWaits({ ...view, over: true }, 1, name)).toEqual([]);
    expect(startWaits({ waiting: { kind: PromptKind.choose, player: 0, purpose: Purpose.target } } as unknown as GameView, 1, name)).toEqual([]);
  });

  it('keeps the picked card under its player zone', () => {
    expect(keptZone(2)).toBe('choice-kept-2');
  });
});
