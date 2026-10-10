import { Board } from '../components/Board';
import { PlayerArea } from '../components/PlayerArea';
import type { Update } from '../protocol';
import { useApp } from '../store';
import { seatPositions } from '../table';
import '../components/Table.css';
import './screens.css';

// GameScreen is the game table (7.5), drawn from the latest update.
export function GameScreen() {
  const game = useApp((s) => s.game);
  const you = useApp((s) => s.lobby.table?.you ?? 0);
  const quit = useApp((s) => s.quit);
  if (!game) {
    return (
      <div className="screen">
        <h1>Setting up the match…</h1>
        <p className="muted">Bans and picks get their own screen in step 7.14.</p>
        <button
          type="button"
          onClick={() => {
            void quit();
          }}
        >
          Leave
        </button>
      </div>
    );
  }
  return <Table game={game} you={you} />;
}

function Table({ game, you }: { game: Update; you: number }) {
  const { view } = game;
  const name = (player: number) => game.seats?.find((s) => s.seat === player)?.name ?? `Player ${String(player + 1)}`;
  return (
    <main className="game" aria-label={`Turn ${String(view.turn.number)}`}>
      {seatPositions(view.players.length, you).map(({ position, player }) => {
        const p = view.players.at(player);
        return (
          p && (
            <PlayerArea
              key={p.id}
              player={p}
              name={name(p.id)}
              position={position}
              you={p.id === you}
              active={view.turn.active === p.id}
            />
          )
        );
      })}
      <Board view={view} />
    </main>
  );
}
