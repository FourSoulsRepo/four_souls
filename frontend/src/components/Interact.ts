// How table elements respond to the player (7.7): highlighted and
// clickable when the player may act on them or pick them.
import { aimed } from "../aim";
import {
  IntentKind,
  Purpose,
  type Action,
  type Actions,
  type ChoiceOption,
} from "../actions";
import type { Kind } from "../cards";
import type { Intent } from "../protocol";
import type { Aiming } from "../store";

export interface Handle {
  highlight?: "act" | "pick" | "aim";
  onClick?: () => void;
  // why says why the player can't use the card now (15.29).
  why?: string;
}

export interface Interact {
  object: (id: number) => Handle;
  player: (id: number) => Handle;
  // character is a player's character card: aimed at, it is the player.
  character: (player: number, id: number) => Handle;
  deck: (kind: Kind) => Handle;
  // play is the play action of a hand card, for drag and drop.
  play: (id: number) => Action | undefined;
}

export const noInteract: Interact = {
  object: () => ({}),
  player: () => ({}),
  character: () => ({}),
  deck: () => ({}),
  play: () => undefined,
};

// aimHandle marks a target of the card the player aims with (15.25).
const aimHandle: Handle = { highlight: "aim" };

// The deck a "deck" option means.
function deckOf(purpose: number): Kind {
  return purpose === Purpose.purchase ? "treasure" : "monster";
}

// interact: aiming is the card the player aims with, if any; aim starts
// aiming with a card whose one action asks for a target (15.25).
export function interact(
  actions: Actions,
  act: (i: Intent) => void,
  select: (id: number) => void,
  aiming: Aiming | null = null,
  aim: (a: Aiming) => void = () => undefined,
  explain: (id: number) => void = () => undefined,
  mine: (id: number) => boolean = () => false,
): Interact {
  const targets = aiming?.targets ?? null;
  const pick = (o: ChoiceOption): Handle => ({
    highlight: "pick",
    onClick: () => {
      act(o.intent);
    },
  });
  const options = actions.choice?.options ?? [];
  const object = (id: number): Handle => {
    if (aimed(targets, (t) => !t.is_player && t.object === id)) {
      return aimHandle;
    }
    // While aiming, only targets answer; the aiming card keeps its
    // glow (15.28).
    if (targets) {
      return id === aiming?.source ? { highlight: "act" } : {};
    }
    const o = options.find((x) => x.target?.object === id);
    if (o) {
      return pick(o);
    }
    const list = actions.cards.get(id);
    if (!list) {
      // No action: a click on one of the player's own cards asks why
      // (the reasons come only then). Other cards keep no click.
      const why = actions.blocked.get(id);
      if (why) {
        return { why };
      }
      return mine(id)
        ? {
            onClick: () => {
              explain(id);
            },
          }
        : {};
    }
    const only = list.length === 1 ? list[0] : undefined;
    if (only?.targets && only.targets.length > 0) {
      const { intent, targets: aims, reason } = only;
      return {
        highlight: "act",
        onClick: () => {
          aim({ source: id, intent, targets: aims, reason });
        },
      };
    }
    return {
      highlight: "act",
      onClick: () => {
        select(id);
      },
    };
  };
  return {
    object,
    // While aiming, a player is picked by the character card (15.34).
    player: (id) => {
      if (targets) {
        return {};
      }
      const o = options.find(
        (x) => x.target?.is_player && x.target.player === id,
      );
      return o ? pick(o) : {};
    },
    // A player is aimed at by the character card (15.33).
    character: (player, id) =>
      aimed(targets, (t) => !!t.is_player && t.player === player)
        ? aimHandle
        : object(id),
    deck: (kind) => {
      const o = options.find(
        (x) => x.target?.deck && deckOf(actions.choice?.purpose ?? 0) === kind,
      );
      return o ? pick(o) : {};
    },
    play: (id) => actions.cards.get(id)?.find((a) => a.intent.kind === IntentKind.playLoot),
  };
}
