// The one WebSocket to the game server (ADR 006). Thin on purpose: the
// lobby logic lives in lobby.ts and the store.
import { t } from './i18n';
import type { Envelope } from './protocol';

let socket: WebSocket | null = null;
let nextID = 1;

// connectTimeout: a host that never answers (or a port the browser
// blocks) must not leave the player waiting forever.
const connectTimeout = 5000;

// connect opens the socket; it resolves once open. onClose runs when the
// server or the network closes it, not on disconnect().
export function connect(
  url: string,
  onMessage: (env: Envelope) => void,
  onClose: () => void,
): Promise<void> {
  disconnect();
  return new Promise((resolve, reject) => {
    const ws = new WebSocket(url);
    socket = ws;
    let open = false;
    const timer = window.setTimeout(() => {
      if (!open && socket === ws) {
        socket = null;
        ws.close();
        reject(new Error(t('error.no_connection')));
      }
    }, connectTimeout);
    ws.onopen = () => {
      open = true;
      window.clearTimeout(timer);
      resolve();
    };
    ws.onmessage = (e: MessageEvent<string>) => {
      onMessage(JSON.parse(e.data) as Envelope);
    };
    ws.onclose = () => {
      if (socket !== ws) {
        return; // closed by disconnect()
      }
      socket = null;
      window.clearTimeout(timer);
      if (open) {
        onClose();
      } else {
        reject(new Error(t('error.no_connection')));
      }
    };
  });
}

// connectFirst tries each URL in turn, e.g. wss:// then ws:// (8.18),
// and resolves with the one that opened. It rejects with the last
// error when none does.
export async function connectFirst(
  urls: string[],
  onMessage: (env: Envelope) => void,
  onClose: () => void,
): Promise<string> {
  let last: unknown = new Error(t('error.no_connection'));
  for (const url of urls) {
    try {
      await connect(url, onMessage, onClose);
      return url;
    } catch (e) {
      last = e;
    }
  }
  throw last;
}

// send sends one message and returns its ID; replies carry the same ID.
export function send(type: string, data: unknown = {}): number {
  const id = nextID++;
  socket?.send(JSON.stringify({ type, id, data }));
  return id;
}

export function disconnect(): void {
  const ws = socket;
  socket = null;
  ws?.close();
}

// drop closes the socket as if the network failed, so onClose runs.
// Tests use it to check reconnecting.
export function drop(): void {
  socket?.close();
}
