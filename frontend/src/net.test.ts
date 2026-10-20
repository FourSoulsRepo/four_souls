import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { connectFirst, disconnect } from './net';

// FakeSocket opens when its URL is in `listening`, and fails otherwise.
const listening = new Set<string>();
const tried: string[] = [];

class FakeSocket {
  onopen: (() => void) | null = null;
  onclose: (() => void) | null = null;
  onmessage: ((e: MessageEvent<string>) => void) | null = null;
  constructor(readonly url: string) {
    tried.push(url);
    queueMicrotask(() => {
      if (listening.has(url)) {
        this.onopen?.();
      } else {
        this.onclose?.();
      }
    });
  }
  close() {
    // closed by the client; nothing to report
  }
}

beforeEach(() => {
  vi.stubGlobal('window', globalThis);
  vi.stubGlobal('WebSocket', FakeSocket);
  listening.clear();
  tried.length = 0;
});

afterEach(() => {
  disconnect();
  vi.unstubAllGlobals();
});

describe('connectFirst', () => {
  const wss = 'wss://games.example.com/ws';
  const ws = 'ws://games.example.com:4774/ws';

  it('takes wss:// when it opens', async () => {
    listening.add(wss).add(ws);
    await expect(connectFirst([wss, ws], vi.fn(), vi.fn())).resolves.toBe(wss);
    expect(tried).toEqual([wss]);
  });

  it('falls back to ws:// when wss:// fails (8.18)', async () => {
    listening.add(ws);
    await expect(connectFirst([wss, ws], vi.fn(), vi.fn())).resolves.toBe(ws);
    expect(tried).toEqual([wss, ws]);
  });

  it('fails when none opens', async () => {
    await expect(connectFirst([wss, ws], vi.fn(), vi.fn())).rejects.toThrow('Could not connect to the server.');
    expect(tried).toEqual([wss, ws]);
  });
});
