// How table elements respond to the player (7.7): highlighted and
// clickable when the player may act on them or pick them.
import { Purpose, type Actions, type ChoiceOption } from '../actions';
import type { Kind } from '../cards';
import type { Intent } from '../protocol';

export interface Handle {
  highlight?: 'act' | 'pick';
  onClick?: () => void;
}

export interface Interact {
  object: (id: number) => Handle;
  player: (id: number) => Handle;
  deck: (kind: Kind) => Handle;
  // play is the play intent of a hand card, for drag and drop.
  play: (id: number) => Intent | undefined;
}

export const noInteract: Interact = {
  object: () => ({}),
  player: () => ({}),
  deck: () => ({}),
  play: () => undefined,
};

// The deck a "deck" option means.
function deckOf(purpose: number): Kind {
  return purpose === Purpose.purchase ? 'treasure' : 'monster';
}

export function interact(actions: Actions, act: (i: Intent) => void, select: (id: number) => void): Interact {
  const pick = (o: ChoiceOption): Handle => ({
    highlight: 'pick',
    onClick: () => {
      act(o.intent);
    },
  });
  const options = actions.choice?.options ?? [];
  return {
    object: (id) => {
      const o = options.find((x) => x.target?.object === id);
      if (o) {
        return pick(o);
      }
      return actions.cards.has(id)
        ? {
            highlight: 'act',
            onClick: () => {
              select(id);
            },
          }
        : {};
    },
    player: (id) => {
      const o = options.find((x) => x.target?.is_player && x.target.player === id);
      return o ? pick(o) : {};
    },
    deck: (kind) => {
      const o = options.find((x) => x.target?.deck && deckOf(actions.choice?.purpose ?? 0) === kind);
      return o ? pick(o) : {};
    },
    play: (id) => actions.cards.get(id)?.find((a) => a.label === 'Play')?.intent,
  };
}
