// Reconnecting by itself after a lost game connection (N-20, 8.24).
// Pure, so it is tested without a socket.
import { t } from './i18n';
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
        title: t('notice.engine_updated.title'),
        text: t('notice.engine_updated.text'),
      };
    case 'game_lost':
      return {
        title: t('notice.game_lost.title'),
        text: t('notice.game_lost.text'),
      };
    case 'game_idle':
      return {
        title: t('notice.game_idle.title'),
        text: n.reason === 'long' ? t('notice.game_idle.long') : t('notice.game_idle.idle'),
      };
    case 'game_ended':
      return {
        title: t('notice.game_idle.title'),
        text: t('notice.game_ended.text'),
      };
    default:
      return { title: t('notice.title'), text: n.kind };
  }
}

// endsGame: after this notice the game is gone; the app shows the
// server's game list.
export function endsGame(n: Notice): boolean {
  return n.kind === 'game_lost' || n.kind === 'game_idle' || n.kind === 'game_ended';
}
