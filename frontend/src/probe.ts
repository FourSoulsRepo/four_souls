// The page side of `--diagnose` (ADR 035, item 3): a probe window shows
// the main menu over the drifting sky, counts its frames, and sends what
// it saw to Go, which prints it for the report and quits.
import { sendProbe, type ProbeResult } from './bridge';
import { useApp } from './store';

// settleMs lets the window, the images and the sky start; measureMs is
// how long frames are counted (ADR 035, item 3.6).
const settleMs = 1000;
const measureMs = 5000;

// slowMs matches diag.SlowMs: a frame that missed two refreshes at 60 Hz.
const slowMs = 34;

// skyLayers is how many images the sky needs (sky.ts).
const skyLayers = 6;

export interface FrameStats {
  frames: number;
  seconds: number;
  slow: number;
  worstMs: number;
}

// frameStats sums up requestAnimationFrame times: the frames between the
// first and the last, how many came late, and the longest wait.
export function frameStats(times: number[]): FrameStats {
  let slow = 0;
  let worstMs = 0;
  for (let i = 1; i < times.length; i++) {
    const gap = (times.at(i) ?? 0) - (times.at(i - 1) ?? 0);
    if (gap > slowMs) {
      slow++;
    }
    worstMs = Math.max(worstMs, gap);
  }
  const frames = Math.max(0, times.length - 1);
  const seconds = frames > 0 ? ((times.at(-1) ?? 0) - (times.at(0) ?? 0)) / 1000 : 0;
  return { frames, seconds, slow, worstMs };
}

// soundTypes are the files the game plays (bindings/assets: MP3) and
// two common others; WebKitGTK needs GStreamer plugins for each.
const soundTypes: [string, string][] = [
  ['mp3', 'audio/mpeg'],
  ['ogg', 'audio/ogg; codecs="vorbis"'],
  ['wav', 'audio/wav'],
];

// playable lists the sound types the webview says it can play.
export function playable(canPlay: (type: string) => string): string[] {
  return soundTypes.filter(([, type]) => canPlay(type) !== '').map(([name]) => name);
}

// gpuNames asks WebGL which renderer draws; "llvmpipe" means the CPU.
function gpuNames(): { renderer: string; vendor: string } {
  try {
    const gl = document.createElement('canvas').getContext('webgl');
    if (!gl) {
      return { renderer: 'no WebGL', vendor: '' };
    }
    const info = gl.getExtension('WEBGL_debug_renderer_info');
    const renderer: unknown = gl.getParameter(info ? info.UNMASKED_RENDERER_WEBGL : gl.RENDERER);
    const vendor: unknown = gl.getParameter(info ? info.UNMASKED_VENDOR_WEBGL : gl.VENDOR);
    return { renderer: String(renderer), vendor: String(vendor) };
  } catch {
    return { renderer: 'no WebGL', vendor: '' };
  }
}

function wait(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// frameTimes collects requestAnimationFrame times for ms.
function frameTimes(ms: number): Promise<number[]> {
  return new Promise((resolve) => {
    const times: number[] = [];
    const tick = (now: number) => {
      times.push(now);
      if (times.length > 1 && now - (times.at(0) ?? now) >= ms) {
        resolve(times);
        return;
      }
      requestAnimationFrame(tick);
    };
    requestAnimationFrame(tick);
  });
}

// runProbe waits for the settings and card data, shows the menu at the
// normal animation speed (never saved), measures, and reports.
export async function runProbe(loading: Promise<void>[]): Promise<void> {
  await Promise.all(loading);
  useApp.setState((s) => ({ settings: { ...s.settings, animation: 'normal' } }));
  useApp.getState().show('menu');
  await wait(settleMs);
  const stats = frameStats(await frameTimes(measureMs));
  const reducedMotion = matchMedia('(prefers-reduced-motion: reduce)').matches;
  const result: ProbeResult = {
    ...stats,
    ...gpuNames(),
    pixelRatio: devicePixelRatio,
    width: innerWidth,
    height: innerHeight,
    sky: useApp.getState().icons.sky.length === skyLayers && !reducedMotion,
    reducedMotion,
    sounds: playable((type) => new Audio().canPlayType(type)),
  };
  await sendProbe(result);
}
