// Reconnecting by itself after a lost game connection (N-20, 8.24).
// Pure, so it is tested without a socket.
import type { Notice } from './protocol';

// retryDelay is the wait in ms before reconnect attempt n (from 0): at
// once, then 1, 2, 4, 8 seconds, then every 10.
export function retryDelay(n: number): number {
  return n <= 0 ? 0 : Math.min(1000 * 2 ** (n - 1), 10000);
}

// lostDialogAfter is how long the small badge shows before the
// "Connection lost" dialog; tests shorten it.
let dialogAfter = 30000;

export function lostDialogAfter(): number {
  return dialogAfter;
}

export function setLostDialogAfter(ms: number): void {
  dialogAfter = ms;
}

// noticeText is the dialog for a server notice.
export function noticeText(n: Notice): { title: string; text: string } {
  switch (n.kind) {
    case 'engine_updated':
      return {
        title: 'The server was updated',
        text: 'The game goes on with a newer rules engine. Some cards may behave differently now.',
      };
    case 'game_lost':
      return {
        title: 'Sorry, the game ended',
        text: 'The server was updated and this game cannot go on. Please start a new one.',
      };
    case 'game_idle':
      return {
        title: 'The game ended',
        text: n.reason === 'long' ? 'It ran for 24 hours, the longest a game may run.' : 'Nobody played it for too long.',
      };
    default:
      return { title: 'Notice', text: n.kind };
  }
}

// endsGame: after this notice the game is gone; the app shows the
// server's game list.
export function endsGame(n: Notice): boolean {
  return n.kind === 'game_lost' || n.kind === 'game_idle';
}
