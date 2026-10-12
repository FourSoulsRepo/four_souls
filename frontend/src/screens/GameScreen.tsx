import { useMemo } from 'react';
import { actionsFor, type Actions } from '../actions';
import { ActionBar } from '../components/ActionBar';
import { ActionMenu } from '../components/ActionMenu';
import { Board } from '../components/Board';
import { ChoiceDialog, DiscardDialog } from '../components/ChoiceDialog';
import { interact, noInteract } from '../components/Interact';
import { PlayerArea } from '../components/PlayerArea';
import { HistoryPanel, TurnBadge } from '../components/InfoPanels';
import { GameMenu } from '../components/GameMenu';
import { EndOverlay, LostOverlay, PauseOverlay } from '../components/Overlays';
import { StackPanel } from '../components/StackPanel';
import { matLook } from '../mats';
import type { Update } from '../protocol';
import { seatOf, useApp } from '../store';
import { waitText, type Names } from '../stack';
import { seatPositions } from '../table';
import '../components/Table.css';
import './screens.css';

const noActions: Actions = { buttons: [], cards: new Map(), choice: null, discard: null };

// GameScreen is the game table (7.5), drawn from the update on screen.
export function GameScreen() {
  const game = useApp((s) => s.game);
  const you = useApp(seatOf);
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
  const cards = useApp((s) => s.cards);
  const mat = useApp((s) => s.settings.mat);
  const mats = useApp((s) => s.mats);
  const look = matLook(mat, mats);
  const connected = useApp((s) => s.connected);
  const behind = useApp((s) => s.behind);
  const sent = useApp((s) => s.sent);
  const act = useApp((s) => s.act);
  const select = useApp((s) => s.select);
  // The player acts only on the latest state, and once per step. While
  // waiting, buttons stay but are disabled, so they do not blink.
  const canAct = !behind && sent !== game.step;
  const actions = useMemo(() => actionsFor(game, (label) => cards.has(label)), [game, cards]);
  const ui = canAct ? interact(actions, act, select) : noInteract;
  const offered = canAct ? actions : noActions;
  const name = (player: number) => game.seats?.find((s) => s.seat === player)?.name ?? `Player ${String(player + 1)}`;
  const names: Names = { player: name, card: (ref) => cards.get(ref)?.name ?? ref };
  const me = view.players.at(you);
  const stacked = (view.stack ?? []).length > 0;
  const skip = game.skip_all ? 'on' : stacked && !view.over ? 'offer' : 'none';
  return (
    <main className={`game ${look.className}`} style={look.style} aria-label={`Turn ${String(view.turn.number)}`}>
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
              ui={ui}
            />
          )
        );
      })}
      <Board view={view} ui={ui}>
        <TurnBadge view={view} names={names} />
        <StackPanel view={view} you={you} names={names} />
        <HistoryPanel />
      </Board>
      <ActionBar
        status={game.pause ? 'Paused' : waitText(view, you, names, actions.choice?.title)}
        buttons={actions.buttons}
        disabled={!canAct}
        deadline={game.deadline ?? 0}
        skip={skip}
      />
      <ActionMenu actions={offered.cards} />
      {offered.choice && <ChoiceDialog key={game.step} choice={offered.choice} />}
      {offered.discard && me && <DiscardDialog key={game.step} discard={offered.discard} hand={me.hand ?? []} />}
      {connected && !view.over && <GameMenu host={you === 0} />}
      {!connected ? (
        <LostOverlay />
      ) : view.over ? (
        <EndOverlay game={game} you={you} name={name} />
      ) : (
        <PauseOverlay game={game} you={you} name={name} />
      )}
    </main>
  );
}
