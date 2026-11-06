// Loads the rules engine (ADR 026, 10): Go's wasm_exec.js, then
// site.wasm.gz, unpacked by the browser's DecompressionStream as it
// downloads and compiled as a stream. Both files sit next to the page;
// make site puts them there. cmd/website sets globalThis.fourSouls.
import type { CardInfo } from '../cards';
import type { Msg } from '../locale';

// Engine is globalThis.fourSouls; every function answers a JSON string
// (docs/wiki/Website-wasm.md).
export interface Engine {
  version(): string;
  sets(): string;
  cards(): string;
  check(payload: string): string;
  notice(): string;
}

// NoticeLine is one line of the fan-game notice; a part with a url is a
// link (internal/legal).
export interface NoticeLine {
  parts: { text: string; url?: string }[];
}

// Answer is the engine's answer to a situation (pkg/rules_engine Answer).
export interface Answer {
  allowed: boolean;
  rule?: string;
  reason?: string;
  say?: Msg;
  events?: { kind: string; player: number; card?: string; amount?: number; label?: Msg }[];
  matches?: boolean;
  mismatch?: string[];
}

// Loaded is what the page reads from the engine once, at start.
export interface Loaded {
  engine: Engine;
  version: { app: string; engine: string };
  cards: CardInfo[];
  notice: NoticeLine[];
}

// LoadError says why the engine did not load: an old browser, or a
// failed download or start (message, for the page to show).
export class LoadError extends Error {
  constructor(
    readonly reason: 'no_decompression' | 'failed',
    message: string,
  ) {
    super(message);
  }
}

declare global {
  // Go is defined by wasm_exec.js.
  var Go: (new () => { importObject: WebAssembly.Imports; run(instance: WebAssembly.Instance): Promise<void> }) | undefined;
  var fourSouls: Engine | undefined;
}

function script(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const el = document.createElement('script');
    el.src = src;
    el.onload = () => {
      resolve();
    };
    el.onerror = () => {
      reject(new Error(`${src} did not load`));
    };
    document.head.append(el);
  });
}

// unpacked is the stream unpacked when it is still gzip (its first two
// bytes are 1f 8b); a stream the browser already unpacked stays as is.
async function unpacked(body: ReadableStream<Uint8Array<ArrayBuffer>>): Promise<ReadableStream<Uint8Array<ArrayBuffer>>> {
  const reader = body.getReader();
  const head = (await reader.read()).value;
  const all = new ReadableStream<Uint8Array<ArrayBuffer>>({
    start(out) {
      if (head) {
        out.enqueue(head);
      }
    },
    async pull(out) {
      const next = await reader.read();
      if (next.done) {
        out.close();
      } else {
        out.enqueue(next.value);
      }
    },
    cancel(reason) {
      return reader.cancel(reason);
    },
  });
  const gzip = head?.[0] === 0x1f && head[1] === 0x8b;
  return gzip ? all.pipeThrough(new DecompressionStream('gzip')) : all;
}

// parse reads a function's JSON answer; {"error"} throws.
export function parse(json: string): unknown {
  const v: unknown = JSON.parse(json);
  if (v !== null && typeof v === 'object' && 'error' in v && typeof v.error === 'string') {
    throw new Error(v.error);
  }
  return v;
}

// loadEngine downloads and starts the engine. progress gets the bytes
// downloaded and the total (0 when the server does not say).
export async function loadEngine(base: string, progress: (loaded: number, total: number) => void): Promise<Loaded> {
  if (typeof DecompressionStream === 'undefined') {
    throw new LoadError('no_decompression', 'DecompressionStream');
  }
  try {
    await script(`${base}wasm_exec.js`);
    const res = await fetch(`${base}site.wasm.gz`);
    if (!res.ok || !res.body) {
      throw new Error(`site.wasm.gz: HTTP ${String(res.status)}`);
    }
    // A host that sends the file with Content-Encoding: gzip (vite
    // preview does) has the browser unpack it; its length is then not
    // the length the page counts.
    const total = res.headers.get('Content-Encoding') ? 0 : Number(res.headers.get('Content-Length') ?? 0);
    let loaded = 0;
    const counted = res.body.pipeThrough(
      new TransformStream<Uint8Array<ArrayBuffer>, Uint8Array<ArrayBuffer>>({
        transform(chunk, out) {
          loaded += chunk.byteLength;
          progress(loaded, total);
          out.enqueue(chunk);
        },
      }),
    );
    const wasm = new Response(await unpacked(counted), { headers: { 'Content-Type': 'application/wasm' } });
    if (!globalThis.Go) {
      throw new Error('wasm_exec.js defined no Go');
    }
    const go = new globalThis.Go();
    const { instance } = await WebAssembly.instantiateStreaming(wasm, go.importObject);
    void go.run(instance);
    const engine = globalThis.fourSouls;
    if (!engine) {
      throw new Error('the engine set no fourSouls');
    }
    return {
      engine,
      version: parse(engine.version()) as Loaded['version'],
      cards: parse(engine.cards()) as CardInfo[],
      notice: parse(engine.notice()) as NoticeLine[],
    };
  } catch (e) {
    throw new LoadError('failed', e instanceof Error ? e.message : String(e));
  }
}
