// How table elements respond to the player (7.7): highlighted and
// clickable when the player may act on them or pick them.
import { aimed } from '../aim';
import { Purpose, type Action, type Actions, type ChoiceOption } from '../actions';
import type { Kind } from '../cards';
import type { Intent } from '../protocol';
import type { Aiming } from '../store';

export interface Handle {
  highlight?: 'act' | 'pick' | 'aim';
  onClick?: () => void;
}

export interface Interact {
  object: (id: number) => Handle;
  player: (id: number) => Handle;
  deck: (kind: Kind) => Handle;
  // play is the play action of a hand card, for drag and drop.
  play: (id: number) => Action | undefined;
}

export const noInteract: Interact = {
  object: () => ({}),
  player: () => ({}),
  deck: () => ({}),
  play: () => undefined,
};

// aimHandle marks a target of the card the player aims with (15.25).
const aimHandle: Handle = { highlight: 'aim' };

// The deck a "deck" option means.
function deckOf(purpose: number): Kind {
  return purpose === Purpose.purchase ? 'treasure' : 'monster';
}

// interact: aiming is the card the player aims with, if any; aim starts
// aiming with a card whose one action asks for a target (15.25).
export function interact(
  actions: Actions,
  act: (i: Intent) => void,
  select: (id: number) => void,
  aiming: Aiming | null = null,
  aim: (a: Aiming) => void = () => undefined,
): Interact {
  const targets = aiming?.targets ?? null;
  const pick = (o: ChoiceOption): Handle => ({
    highlight: 'pick',
    onClick: () => {
      act(o.intent);
    },
  });
  const options = actions.choice?.options ?? [];
  return {
    object: (id) => {
      if (aimed(targets, (t) => !t.is_player && t.object === id)) {
        return aimHandle;
      }
      const o = options.find((x) => x.target?.object === id);
      if (o) {
        return pick(o);
      }
      const list = actions.cards.get(id);
      if (!list) {
        return {};
      }
      const only = list.length === 1 ? list[0] : undefined;
      if (only?.targets && only.targets.length > 0) {
        const { intent, targets: aims, reason } = only;
        return {
          highlight: 'act',
          onClick: () => {
            aim({ source: id, intent, targets: aims, reason });
          },
        };
      }
      return {
        highlight: 'act',
        onClick: () => {
          select(id);
        },
      };
    },
    player: (id) => {
      if (aimed(targets, (t) => !!t.is_player && t.player === id)) {
        return aimHandle;
      }
      const o = options.find((x) => x.target?.is_player && x.target.player === id);
      return o ? pick(o) : {};
    },
    deck: (kind) => {
      const o = options.find((x) => x.target?.deck && deckOf(actions.choice?.purpose ?? 0) === kind);
      return o ? pick(o) : {};
    },
    play: (id) => actions.cards.get(id)?.find((a) => a.label === 'Play'),
  };
}
