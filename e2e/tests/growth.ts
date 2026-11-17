// Long games stay fast (TS-12, 22.9): games.spec.ts measures each page
// while a whole game plays, and fails a page that grows or slows with
// the game's length, as 19.135 and 19.138 did.
//
// Three measures per page:
//  1. elements: the page's element count at the game's end.
//  2. script time per step, early and late in the game: Chromium's own
//     ScriptDuration (CDP Performance.getMetrics), read every sampleEvery
//     steps. Its change over a window, divided by the steps in it, is the
//     page's script time per step there.
//  3. sounds: the most Audio elements that held a media player at once
//     (19.135), read at each sample.
import type { CDPSession, Page } from '@playwright/test';

// Sample is one reading of a page: the page's step, its script time so
// far in seconds, its element count and its live sounds.
export interface Sample {
  step: number;
  script: number;
  elements: number;
  sounds: number;
}

// sampleEvery: a page is read when its step has moved this far.
export const sampleEvery = 100;

// Windows of the script time per step. early starts at skipSteps,
// after the deal and the first roll; late is the game's last
// windowSteps. Each spans at least windowSteps; a game too short for
// two apart has no growth check.
export const skipSteps = 200;
export const windowSteps = 300;

// limits: about twice today's values (2026-11-17; 15 games, every set,
// 2-4 players and the challenges, up to 10547 steps; one Mac, Chromium,
// one game at a time, speed "fast"). Measured:
//  - elements at the end: 580-869;
//  - script ms a step, early or late: 5.6-12.9;
//  - growth, late / early: 0.64-1.19; 19.138 was about 5;
//  - live sounds: 0-1; the test build has no sound files, so each
//    play fails and is released at once.
// growth is the main check: machine load slows both windows alike.
// scriptMs is a backstop for a page slow from the start; it depends on
// the machine.
export const limits = {
  elements: 1700,
  scriptMs: 25,
  growth: 2,
  sounds: 4,
};

// openMetrics turns Chromium's performance counters on for a page.
export async function openMetrics(cdp: CDPSession): Promise<void> {
  await cdp.send('Performance.enable');
}

// sample reads one page. ScriptDuration counts the page's own tasks
// only, not a test's evaluate: neither this read nor the bot's polling
// adds to it.
export async function sample(page: Page, cdp: CDPSession, step: number): Promise<Sample> {
  const { metrics } = await cdp.send('Performance.getMetrics');
  const script = metrics.find((m) => m.name === 'ScriptDuration')?.value ?? 0;
  const { elements, sounds } = await page.evaluate(() => {
    const w = window as unknown as { fourSouls?: { sounds?: { live: () => number } } };
    return { elements: document.getElementsByTagName('*').length, sounds: w.fourSouls?.sounds?.live() ?? 0 };
  });
  return { step, script, elements, sounds };
}

// perStep is the script time per step, in ms, between two samples.
function perStep(a: Sample, b: Sample): number {
  return ((b.script - a.script) * 1000) / (b.step - a.step);
}

// Measure is what a whole game left on one page.
export interface Measure {
  page: string;
  steps: number;
  elements: number;
  sounds: number;
  early?: number;
  late?: number;
}

// measured sums up one page's samples, oldest first.
export function measured(page: string, samples: Sample[]): Measure {
  const end = samples[samples.length - 1];
  const m: Measure = {
    page,
    steps: end?.step ?? 0,
    elements: end?.elements ?? 0,
    sounds: Math.max(0, ...samples.map((s) => s.sounds)),
  };
  const a = samples.find((s) => s.step >= skipSteps);
  const b = a && samples.find((s) => s.step >= a.step + windowSteps);
  if (a && b) {
    m.early = perStep(a, b);
  }
  const d = end;
  const c = d && samples.findLast((s) => s.step <= d.step - windowSteps);
  if (b && c && d && c.step >= b.step) {
    m.late = perStep(c, d);
  }
  return m;
}

// describe is one line for the log.
export function describe(m: Measure): string {
  const ms = (v?: number) => (v === undefined ? '-' : v.toFixed(1));
  const growth = m.early && m.late ? ` (x${(m.late / m.early).toFixed(2)})` : '';
  return `${m.page}: ${String(m.steps)} steps, ${String(m.elements)} elements, ${String(m.sounds)} sounds, script ${ms(m.early)} ms a step early, ${ms(m.late)} late${growth}`;
}

// overLimits names each limit a page broke in the game.
export function overLimits(m: Measure, game: string): string[] {
  const out: string[] = [];
  const over = (what: string, value: number, limit: number, text: string) => {
    if (value > limit) {
      out.push(`${m.page}: ${what} ${text}, limit ${String(limit)}; ${game}`);
    }
  };
  over('elements at the end', m.elements, limits.elements, String(m.elements));
  over('live sounds', m.sounds, limits.sounds, String(m.sounds));
  if (m.early !== undefined) {
    over('script ms a step early', m.early, limits.scriptMs, m.early.toFixed(1));
  }
  if (m.late !== undefined) {
    over('script ms a step late', m.late, limits.scriptMs, m.late.toFixed(1));
  }
  if (m.early !== undefined && m.late !== undefined) {
    over('script growth late / early', m.late / m.early, limits.growth, (m.late / m.early).toFixed(2));
  }
  return out;
}
