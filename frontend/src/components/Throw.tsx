import type { Actions } from '../actions';
import { PromptKind } from '../actions';
import { useKey } from '../keys';
import type { Prompt } from '../protocol';
import type { Names } from '../stack';
import { useApp } from '../store';
import { playerColor } from '../table';
import { Die } from './Die';
import './Throw.css';

interface Props {
  waiting: Prompt;
  you: number;
  names: Names;
  // action throws, when the engine asks this player now.
  action: Actions['throw'];
  disabled: boolean;
}

// The face a die waiting in the hand shows with a texture.
const handFace = 6;

// Throw shows a roll waiting for its throw (15.61, V-17) where the die
// will land. The roller gets a big glowing Throw button with their die;
// Space presses it. Everyone else sees the roller's name and the die
// shaking in their hand. It hides while an effect plays, so the die
// that lands is the only one.
export function Throw({ waiting, you, names, action, disabled }: Props) {
  const act = useApp((s) => s.act);
  const effect = useApp((s) => s.playing !== null);
  const diceImage = useApp((s) => s.dice.find((d) => d.id === s.settings.dice)?.image ?? '');
  const lit = action !== null && !disabled;
  const press = () => {
    if (lit) {
      act(action.intent);
    }
  };
  useKey(' ', press);
  if (waiting.kind !== PromptKind.throw || effect) {
    return null;
  }
  const mine = waiting.player === you;
  const label = waiting.text ?? '';
  if (lit) {
    return (
      <div className="throw" role="region" aria-label="Throw">
        <button type="button" className="throw-button" title="Throw (Space)" onClick={press}>
          <Die value={handFace} image={diceImage} label="Your die" hidden />
          <span className="throw-word">Throw</span>
        </button>
        {label && <div className="throw-label">{label}</div>}
      </div>
    );
  }
  const who = mine ? 'You are' : `${names.player(waiting.player)} is`;
  return (
    <div className="throw" role="region" aria-label="Throw">
      <div className="throw-shake">
        <Die value={handFace} image={diceImage} label={mine ? 'Your die' : `${names.player(waiting.player)}'s die`} hidden />
      </div>
      <div className="throw-label" style={{ color: playerColor(waiting.player) }}>
        {who} throwing…
      </div>
    </div>
  );
}
