// Damage sounds (15.50), all from Hearthstone. An attack lands with an
// unarmed hit, card to card. Other damage is a fireball: its cast as the
// shot leaves and its impact as it lands. A thrown die clatters on wood
// as it lands (15.69), each die of several on its own (19.39). Each
// looted card is dealt with a card sound (18.7). A player's death sounds
// with its skull (18.8). Each gained treasure sounds as it flies in
// (18.9). Every soul gained sounds with one sound: as it flies to its
// player, a bonus soul too (18.13), and as the Lost Soul turns to one
// (18.14). The reference client's sounds of step 20 (20.17) are here too;
// each later sub-step of step 20 plays its own: a bigger attack hit
// sounds bigger since 20.3; a held card's two loops since 20.6; another
// player's big played card since 20.9; End turn's cue since 20.10; the
// response timer's ticks since 20.15. A sound may have several files:
// each play picks one at random, all alike, as the reference does. The
// files live only in the private images; without them the game is
// silent.
// Every file is levelled to one loudness (ADR 022), so every sound plays
// at the same level. The player's volume (ST-09) scales them all alike.
// One sound per moment (20.14): one sound never starts twice within
// feel.mix.gap, and a key cue ducks every other sound while it plays.
// Ducking is the mixer's gain on the rest, never a sound's own level.
import type { Sounds } from './bridge';
import { feel } from './feel';

export type Sound = keyof Sounds;

// level is the player's volume, 0 to 1 (ST-09).
let level = 1;

let urls = new Map<Sound, string[]>();

// loadSounds takes the sound URLs, a list per sound; an empty list is no file.
export function loadSounds(s: Sounds): void {
  urls = new Map(Object.entries(s) as [Sound, string[]][]);
}

// setSoundVolume sets the player's volume, in percent; 0 is mute.
export function setSoundVolume(percent: number): void {
  level = percent / 100;
  if (master) {
    master.gain.value = level;
  }
  mix();
}

// keyCues duck every other sound while they play (20.14, ADR 027
// principle 8): the voice line, as the reference's voice ducks.
const keyCues: ReadonlySet<Sound> = new Set<Sound>(['endTurnCue']);

// clock is the mixer's time, in seconds.
function clock(): number {
  return performance.now() / 1000;
}

// started is when each sound last started, by clock.
const started = new Map<Sound, number>();

// slack is a µs, so a time exactly at a limit never misses it by rounding.
const slack = 1e-6;

// mayStart tells if a sound that last started at last may start again
// at now (20.14): not within feel.mix.gap.
export function mayStart(last: number | undefined, now: number): boolean {
  return last === undefined || now - last >= feel.mix.gap - slack;
}

// Cue is a key cue that plays from start until end, by clock; end is
// unset while it plays.
export interface Cue {
  start: number;
  end?: number;
}

// duckGain is the mixer's gain on every sound but the key cues at now
// (20.14): each cue takes it down to feel.mix.duck over duckIn, linear,
// and once it ends back up over duckOut, from wherever it got to, as
// the reference does. Several cues: the lowest gain wins.
export function duckGain(now: number, cues: readonly Cue[]): number {
  const { duck, duckIn, duckOut } = feel.mix;
  let gain = 1;
  for (const c of cues) {
    if (now < c.start) {
      continue;
    }
    const down = (t: number) => 1 - (1 - duck) * Math.min(1, (t - c.start) / duckIn);
    let g = down(now);
    if (c.end !== undefined && now > c.end) {
      const from = down(c.end);
      g = from + (1 - from) * Math.min(1, (now - c.end) / duckOut);
    }
    gain = Math.min(gain, g);
  }
  return gain;
}

// cues are the key cues playing or still coming back; ducked are the
// other sounds playing now.
const cues: Cue[] = [];
const ducked = new Set<HTMLAudioElement>();
let ticker: ReturnType<typeof setInterval> | undefined;

// mix sets the gain on the rest from the cues, and forgets the cues
// whose duck is over.
function mix(): void {
  const now = clock();
  const over = cues.filter((c) => c.end !== undefined && now - c.end >= feel.mix.duckOut - slack);
  for (const c of over) {
    cues.splice(cues.indexOf(c), 1);
  }
  const gain = duckGain(now, cues);
  for (const a of ducked) {
    a.volume = level * gain;
  }
  if (mixer && rest) {
    rest.gain.setValueAtTime(gain, mixer.currentTime);
  }
  if (cues.length === 0 && ticker !== undefined) {
    clearInterval(ticker);
    ticker = undefined;
  }
}

// playSound plays a sound once, one of its files at random; sounds may
// overlap, but one sound never starts twice within feel.mix.gap. A
// sound without a file, or one that fails to load or is blocked, stays
// silent.
export function playSound(sound: Sound): void {
  const files = urls.get(sound) ?? [];
  if (files.length === 0 || typeof Audio === 'undefined') {
    return;
  }
  const now = clock();
  if (!mayStart(started.get(sound), now)) {
    return;
  }
  started.set(sound, now);
  const a = new Audio(files[Math.floor(Math.random() * files.length)]);
  if (keyCues.has(sound)) {
    cue(a, now);
    return;
  }
  a.volume = level * duckGain(now, cues);
  ducked.add(a);
  const done = () => {
    ducked.delete(a);
  };
  a.onended = done;
  a.onerror = done;
  a.play().catch(done);
}

// cue plays a key cue at full level and ducks the rest until it ends.
// A cue that cannot play ducks nothing.
function cue(a: HTMLAudioElement, now: number): void {
  a.volume = level;
  const c: Cue = { start: now };
  cues.push(c);
  const end = () => {
    clearTimeout(longest);
    c.end ??= Math.min(clock(), c.start + feel.mix.cueMax);
  };
  const longest = setTimeout(end, feel.mix.cueMax * 1000);
  a.onended = end;
  a.onerror = end;
  a.play().catch(() => {
    clearTimeout(longest);
    const i = cues.indexOf(c);
    if (i >= 0) {
      cues.splice(i, 1);
    }
    mix();
  });
  ticker ??= setInterval(mix, feel.mix.tick * 1000);
  mix();
}

// Loops (20.6): a sound that repeats until it is stopped, as a held
// card's drag-bed and drag-air. They play through Web Audio, so the
// loop has no gap, and a gain node per loop sets its level: a mixer
// gain, never a per-file volume (ADR 022). A master gain is the
// player's volume; the loops reach it through the rest's gain, which a
// key cue ducks (20.14). The context is made at the first loop, inside
// a drag, so the browser counts it as the player's own.
let mixer: AudioContext | undefined;
let master: GainNode | undefined;
let rest: GainNode | undefined;
const decoded = new Map<string, Promise<AudioBuffer | null>>();

// out is the mixer, made at the first call: the rest's gain, which
// leads to the master gain; undefined where the browser has no Web
// Audio.
function out(): { ctx: AudioContext; rest: GainNode } | undefined {
  if (typeof AudioContext === 'undefined') {
    return undefined;
  }
  if (!mixer || !master || !rest) {
    try {
      mixer = new AudioContext();
      master = mixer.createGain();
      master.gain.value = level;
      master.connect(mixer.destination);
      rest = mixer.createGain();
      rest.gain.value = duckGain(clock(), cues);
      rest.connect(master);
    } catch {
      return undefined;
    }
  }
  if (mixer.state === 'suspended') {
    mixer.resume().catch(() => undefined);
  }
  return { ctx: mixer, rest };
}

// decode loads and decodes a file once; a file that fails is null.
function decode(ctx: AudioContext, url: string): Promise<AudioBuffer | null> {
  let p = decoded.get(url);
  if (!p) {
    p = fetch(url)
      .then((r) => (r.ok ? r.arrayBuffer() : Promise.reject(new Error(r.statusText))))
      .then((b) => ctx.decodeAudioData(b))
      .catch(() => null);
    decoded.set(url, p);
  }
  return p;
}

// loopBounds is where a loop's sound starts and ends, in seconds. An
// MP3 may have a few ms of encoder silence at either end, and some
// browsers keep it; the loop skips it, or it would tick at each turn.
// It finds the first and last 2 ms window at a quarter of the loop's
// overall RMS or more.
export function loopBounds(data: Float32Array, rate: number): { start: number; end: number } {
  const n = data.length;
  const win = Math.max(1, Math.round(rate * 0.002));
  let sum = 0;
  for (const v of data) {
    sum += v * v;
  }
  const floor = (sum / Math.max(n, 1)) * 0.25 * 0.25;
  const loud = (from: number): boolean => {
    const part = data.subarray(from, Math.min(from + win, n));
    let s = 0;
    for (const v of part) {
      s += v * v;
    }
    return s / Math.max(part.length, 1) >= floor;
  };
  let start = 0;
  while (start < n && !loud(start)) {
    start += win;
  }
  let end = n;
  while (end > start && !loud(Math.max(end - win, start))) {
    end -= win;
  }
  if (end <= start) {
    return { start: 0, end: n / rate };
  }
  return { start: start / rate, end: end / rate };
}

// Loop is a playing loop.
export interface Loop {
  // set moves its mixer gain, 0 to 1, to gain over `over` seconds.
  set(gain: number, over: number): void;
  // stop fades it out over `over` seconds and ends it.
  stop(over: number): void;
}

const silent: Loop = { set: () => undefined, stop: () => undefined };

// playLoop starts a loop at gain 0, one of its files at random. A sound
// without a file, a file that fails to load, or a browser without Web
// Audio stays silent.
export function playLoop(sound: Sound): Loop {
  const files = urls.get(sound) ?? [];
  const mixed = files.length > 0 ? out() : undefined;
  const url = files[Math.floor(Math.random() * files.length)];
  if (!mixed || !url) {
    return silent;
  }
  const { ctx } = mixed;
  const gain = ctx.createGain();
  gain.gain.value = 0;
  gain.connect(mixed.rest);
  let source: AudioBufferSourceNode | undefined;
  let stopped = false;
  void decode(ctx, url).then((buffer) => {
    if (!buffer || stopped) {
      gain.disconnect();
      return;
    }
    const { start, end } = loopBounds(buffer.getChannelData(0), buffer.sampleRate);
    source = ctx.createBufferSource();
    source.buffer = buffer;
    source.loop = true;
    source.loopStart = start;
    source.loopEnd = end;
    source.connect(gain);
    source.start(0, start);
  });
  const ramp = (to: number, over: number) => {
    const now = ctx.currentTime;
    const g = gain.gain;
    g.cancelScheduledValues(now);
    g.setValueAtTime(g.value, now);
    g.linearRampToValueAtTime(to, now + Math.max(over, 0.001));
  };
  return {
    set: (to, over) => {
      if (!stopped) {
        ramp(to, over);
      }
    },
    stop: (over) => {
      if (stopped) {
        return;
      }
      stopped = true;
      ramp(0, over);
      if (source) {
        const s = source;
        s.onended = () => {
          gain.disconnect();
        };
        s.stop(ctx.currentTime + Math.max(over, 0.001));
      }
    },
  };
}
