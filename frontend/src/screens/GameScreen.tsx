import { useEffect, useLayoutEffect, useMemo, useState } from 'react';
import { actionsFor, aimedOption, IntentKind, Purpose, type Actions } from '../actions';
import { onScreen, promptAim } from '../aim';
import { ActionBar } from '../components/ActionBar';
import { ActionMenu } from '../components/ActionMenu';
import { AimLayer } from '../components/Aim';
import { Board } from '../components/Board';
import { ChoiceDialog, DiscardDialog, LookDialog, OrderDialog } from '../components/ChoiceDialog';
import { interact, noInteract } from '../components/Interact';
import { PlayerArea } from '../components/PlayerArea';
import { PresenceLayer, PresenceSender } from '../components/Presence';
import { HistoryPanel, TurnBadge } from '../components/InfoPanels';
import { GameMenu } from '../components/GameMenu';
import { HistoryStrip } from '../components/HistoryStrip';
import { EndOverlay, LostOverlay, PauseOverlay } from '../components/Overlays';
import { PlayedCard, StackPanel } from '../components/StackPanel';
import { Throw } from '../components/Throw';
import { WhyNote } from '../components/Why';
import { EndTurn, TurnActions } from '../components/TurnActions';
import { matLook } from '../mats';
import type { Update } from '../protocol';
import { seatOf, useApp } from '../store';
import { activity, picking, promptBanner, stackPrompt } from '../prompt';
import { waitTurn, type Names } from '../stack';
import { seatPositions } from '../table';
import '../components/Table.css';
import './screens.css';

const noActions: Actions = { buttons: [], turn: {}, cards: new Map(), blocked: new Map(), choice: null, discard: null, order: null, throw: null };

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
  const looking = useApp((s) => s.looking);
  const aiming = useApp((s) => s.aiming);
  const aimed = useApp((s) => s.aimed);
  const takeAimed = useApp((s) => s.takeAimed);
  // The full history opens from the strip; the stack has its room (15.22).
  const [historyOpen, setHistoryOpen] = useState(false);
  // The player acts only on the latest state, and once per step. While
  // waiting, buttons stay but are disabled, so they do not blink.
  const canAct = !behind && sent !== game.step;
  const actions = useMemo(() => actionsFor(game, (label) => cards.has(label)), [game, cards]);
  const setAiming = useApp((s) => s.setAiming);
  const ui = canAct ? interact(actions, act, select, aiming, setAiming) : noInteract;
  // A new state ends an aim: its targets may be gone (15.25). Before
  // paint, so a placed card's aim below replaces it without a blink.
  useLayoutEffect(() => {
    setAiming(null);
  }, [game.step, setAiming]);
  // A target prompt of the player's own play or use aims from its card,
  // with no dialog (15.34, T-13). A card aimed already gets its pick
  // below; a target off the table keeps the dialog.
  const placed = useMemo(() => promptAim(view.activating, actions.choice), [view.activating, actions.choice]);
  useLayoutEffect(() => {
    if (placed && canAct && !aimed && onScreen(placed)) {
      setAiming(placed);
    }
  }, [placed, canAct, aimed, setAiming]);
  // A card aimed at a target: its target prompt picks that target at
  // once (15.25). Only the first state after the play gets the chance.
  useEffect(() => {
    if (!aimed || !canAct || game.step <= aimed.step) {
      return;
    }
    takeAimed();
    const pick = aimedOption(actions.choice, aimed.target);
    if (pick) {
      act(pick);
    }
  }, [aimed, canAct, game.step, actions, act, takeAimed]);
  const offered = canAct ? actions : noActions;
  const name = (player: number) => game.seats?.find((s) => s.seat === player)?.name ?? `Player ${String(player + 1)}`;
  const names: Names = { player: name, card: (ref) => cards.get(ref)?.name ?? ref };
  const me = view.players.at(you);
  // The hand's cards in order: a hand card is shared as its slot (A-16).
  const hand = useMemo(() => (me?.hand ?? []).map((c) => c.id), [me]);
  // The banner and the stack say what the game waits for (15.28). While
  // the player picks a target, the cards that answer glow; the rest dim.
  const banner = promptBanner(view, you, names, actions, aiming);
  const dim = canAct && picking(actions, aiming);
  // What the active player does, on their seat for everyone else (15.30).
  const doing = activity(view, names);
  return (
    <main className={`game ${look.className}${dim ? ' picking' : ''}`} style={look.style} aria-label={`Turn ${String(view.turn.number)}`}>
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
              doing={doing}
              winSouls={view.win_souls ?? 0}
              ui={ui}
              tools={
                p.id === you && !view.over ? (
                  <TurnActions
                    allowed={actions.turn}
                    disabled={!canAct}
                    step={game.step}
                    waiting={view.waiting}
                    you={you}
                    left={view.turn.active === you ? { buy: p.purchases, attack: p.attacks } : undefined}
                  />
                ) : undefined
              }
            />
          )
        );
      })}
      <Board
        view={view}
        ui={ui}
        strip={
          <HistoryStrip
            view={view}
            names={names}
            open={historyOpen}
            toggle={() => {
              setHistoryOpen(!historyOpen);
            }}
          />
        }
        end={
          me && !view.over ? (
            <EndTurn
              intent={actions.turn.endTurn}
              idle={actions.turn.idle ?? false}
              left={actions.turn.left}
              own={view.turn.active === you}
              disabled={!canAct}
              step={game.step}
              why={actions.turn.why?.endTurn}
            />
          ) : undefined
        }
      >
        <TurnBadge view={view} names={names} />
        <StackPanel
          view={view}
          you={you}
          names={names}
          pass={actions.buttons.find((b) => b.intent.kind === IntentKind.pass)?.intent}
          disabled={!canAct}
          prompt={stackPrompt(banner, names)}
          answer={banner?.mine === true && (banner.kind === 'respond' || banner.kind === 'act')}
        />
        <PlayedCard view={view} you={you} names={names} />
        {historyOpen && (
          <HistoryPanel
            hide={() => {
              setHistoryOpen(false);
            }}
          />
        )}
      </Board>
      {!view.over && !game.pause && (
        <Throw waiting={view.waiting} you={you} names={names} action={actions.throw} disabled={!canAct} />
      )}
      <ActionBar
        banner={banner}
        paused={!!game.pause}
        name={name}
        deadline={game.deadline ?? 0}
        wait={waitTurn(view, you, game.wait_turn ?? false)}
      />
      <ActionMenu actions={offered.cards} />
      {canAct && <AimLayer />}
      {me && !view.over && <PresenceSender active={view.turn.active === you} hand={hand} />}
      <PresenceLayer you={you} />
      <WhyNote />
      {offered.choice && !aiming?.choice && (
        <ChoiceDialog
          key={game.step}
          choice={offered.choice}
          subject={offered.choice.purpose === Purpose.monsterSlot ? view.revealed?.card : view.subject?.card}
        />
      )}
      {offered.discard && me && <DiscardDialog key={game.step} discard={offered.discard} hand={me.hand ?? []} />}
      {offered.order && <OrderDialog key={game.step} order={offered.order} cards={view.ordering ?? []} />}
      {looking.length > 0 && <LookDialog cards={looking} />}
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
