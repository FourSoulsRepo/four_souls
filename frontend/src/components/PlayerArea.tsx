import { Role, type CardView, type PlayerView } from '../protocol';
import type { Position } from '../table';
import { Card, type CardSize } from './Card';
import './Table.css';

interface Props {
  player: PlayerView;
  name: string;
  position: Position;
  you: boolean;
  active: boolean;
}

// spent: a character or item that is used up until it recharges.
function spent(c: CardView): boolean {
  return (c.role === Role.character || c.role === Role.item) && !c.charged;
}

// PlayerArea is one player's place at the table: name and counters,
// character and cards in play, and the hand for the viewer (V-02).
export function PlayerArea({ player, name, position, you, active }: Props) {
  const size: CardSize = position === 'left' || position === 'right' ? 'small' : 'table';
  const label = you ? `${name} (you)` : name;
  return (
    <section
      className={`player player-${position}${active ? ' active' : ''}${player.dead ? ' dead' : ''}`}
      aria-label={label}
      data-zone={`player-${String(player.id)}`}
    >
      <header className="player-header">
        <b>{label}</b>
        <span className="player-stats">
          HP {player.hp} · ATK {player.atk} · {player.cents}¢ · {player.souls} souls
          {!you && ` · ${String(player.hand_size)} in hand`}
          {player.dead && ' · dead'}
        </span>
      </header>
      <div className="player-cards">
        <Card
          cardRef={player.character.card}
          size={size}
          damage={player.character.damage}
          spent={spent(player.character)}
          objectId={player.character.id}
        />
        {(player.in_play ?? []).map((c) => (
          <Card key={c.id} cardRef={c.card} size={size} damage={c.damage} spent={spent(c)} objectId={c.id} />
        ))}
      </div>
      {you && (
        <div className="hand" role="group" aria-label="Your hand" data-zone={`hand-${String(player.id)}`}>
          {(player.hand ?? []).map((c) => (
            <Card key={c.id} cardRef={c.card} size="hand" objectId={c.id} />
          ))}
        </div>
      )}
    </section>
  );
}
