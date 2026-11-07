import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { Sounds } from './bridge';
import { duckGain, loadSounds, loopBounds, mayStart, playLoop, playSound, setSoundVolume } from './sounds';

const hit = '/cardimg/assets/sounds/attack-hit.mp3';
const cast = '/cardimg/assets/sounds/effect-cast.mp3';
const impact = '/cardimg/assets/sounds/effect-impact.mp3';
const dice = '/cardimg/assets/sounds/dice-roll.mp3';
const deal = '/cardimg/assets/sounds/loot-draw.mp3';
const death = '/cardimg/assets/sounds/player-died.mp3';
const gain = '/cardimg/assets/sounds/treasure-gain.mp3';
const soul = '/cardimg/assets/sounds/soul-gain.mp3';

const lands = ['1', '2', '3'].map((i) => `/cardimg/assets/sounds/loot-land-${i}.mp3`);

const none: Sounds = {
  attackHit: [],
  effectCast: [],
  effectImpact: [],
  diceRoll: [],
  lootDraw: [],
  playerDied: [],
  treasureGain: [],
  soulGain: [],
  attackHit2: [],
  attackHit3: [],
  lootLand: [],
  lootDrawOther: [],
  dragBed: [],
  dragAir: [],
  turnStart: [],
  cardPlayOther: [],
  endTurnCue: [],
  timerTick: [],
  timerLast: [],
};

const all: Sounds = {
  ...none,
  attackHit: [hit],
  effectCast: [cast],
  effectImpact: [impact],
  diceRoll: [dice],
  lootDraw: [deal],
  playerDied: [death],
  treasureGain: [gain],
  soulGain: [soul],
  lootLand: lands,
};

// FakeAudio is a played file; onended ends it.
interface FakeAudio {
  src: string;
  volume: number;
  onended: (() => void) | null;
}

// fakeAudio records each played URL; play fails with the given error.
// elements gets each played file.
function fakeAudio(fail?: DOMException, volumes: number[] = [], elements: FakeAudio[] = []): string[] {
  const played: string[] = [];
  vi.stubGlobal(
    'Audio',
    class {
      volume = 0.5;
      onended: (() => void) | null = null;
      constructor(public src: string) {}
      play(): Promise<void> {
        played.push(this.src);
        volumes.push(this.volume);
        elements.push(this);
        return fail ? Promise.reject(fail) : Promise.resolve();
      }
    },
  );
  return played;
}

// clock is the mixer's time in ms; each test starts well after the last
// one, so no start is within the playback limit (20.14) of an earlier.
let clock = 0;

beforeEach(() => {
  clock += 60_000;
  vi.spyOn(performance, 'now').mockImplementation(() => clock);
});

afterEach(() => {
  vi.useRealTimers();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe('playSound', () => {
  it('plays each sound from its file', () => {
    loadSounds(all);
    const played = fakeAudio();
    playSound('attackHit');
    playSound('effectCast');
    playSound('effectImpact');
    playSound('diceRoll');
    playSound('lootDraw');
    playSound('playerDied');
    playSound('treasureGain');
    playSound('soulGain');
    expect(played).toEqual([hit, cast, impact, dice, deal, death, gain, soul]);
  });

  // The files are levelled to one loudness (ADR 022); none gets its own volume.
  it('plays every sound at full volume', () => {
    loadSounds(all);
    const volumes: number[] = [];
    fakeAudio(undefined, volumes);
    playSound('attackHit');
    playSound('effectCast');
    playSound('effectImpact');
    playSound('diceRoll');
    playSound('lootDraw');
    playSound('playerDied');
    playSound('treasureGain');
    playSound('soulGain');
    expect(volumes).toEqual([1, 1, 1, 1, 1, 1, 1, 1]);
  });

  // The player's volume (ST-09) scales every sound alike; mute plays nothing audible.
  it('plays every sound at the player volume', () => {
    loadSounds(all);
    const volumes: number[] = [];
    fakeAudio(undefined, volumes);
    setSoundVolume(40);
    playSound('attackHit');
    playSound('diceRoll');
    setSoundVolume(0);
    playSound('treasureGain');
    playSound('soulGain');
    setSoundVolume(100);
    playSound('lootDraw');
    expect(volumes).toEqual([0.4, 0.4, 0, 0, 1]);
  });

  it('is silent without the private files', () => {
    loadSounds(none);
    const played = fakeAudio();
    for (const sound of Object.keys(none) as (keyof Sounds)[]) {
      playSound(sound);
    }
    expect(played).toEqual([]);
  });

  // A sound with several files plays one of them, each alike (20.17).
  it('picks among the variants at random', () => {
    loadSounds(all);
    const played = fakeAudio();
    const random = vi.spyOn(Math, 'random');
    for (const r of [0, 0.32, 0.34, 0.65, 0.67, 0.999]) {
      random.mockReturnValueOnce(r);
      playSound('lootLand');
      clock += 100;
    }
    expect(played).toEqual([lands[0], lands[0], lands[1], lands[1], lands[2], lands[2]]);
    random.mockRestore();
  });

  // A single file is a list of one: it always plays.
  it('plays a single file every time', () => {
    loadSounds(all);
    const played = fakeAudio();
    const random = vi.spyOn(Math, 'random').mockReturnValue(0.999);
    playSound('attackHit');
    clock += 100;
    playSound('attackHit');
    expect(played).toEqual([hit, hit]);
    random.mockRestore();
  });

  it('swallows a file that fails to load', async () => {
    loadSounds(all);
    const played = fakeAudio(new DOMException('no file', 'NotSupportedError'));
    playSound('attackHit');
    await Promise.resolve();
    expect(played).toEqual([hit]);
  });
});

// One sound per moment (20.14, ADR 027 principle 8).
describe('the playback limit', () => {
  it('lets a sound start again only 80 ms after its last start', () => {
    expect(mayStart(undefined, 5)).toBe(true);
    expect(mayStart(5, 5)).toBe(false);
    expect(mayStart(5, 5.079)).toBe(false);
    expect(mayStart(5, 5.08)).toBe(true);
  });

  it('drops a second start within 80 ms; other sounds still play', () => {
    loadSounds(all);
    const played = fakeAudio();
    playSound('diceRoll');
    clock += 50;
    playSound('diceRoll');
    playSound('attackHit');
    clock += 30;
    playSound('diceRoll');
    expect(played).toEqual([dice, hit, dice]);
  });
});

describe('ducking', () => {
  it('takes the rest to 20 % over 0.7 s and back over 0.7 s, linear', () => {
    const cue = { start: 10, end: 12 };
    const at = (t: number) => duckGain(t, [cue]);
    expect(duckGain(10.35, [])).toBe(1);
    expect(at(9)).toBe(1);
    expect(at(10)).toBe(1);
    expect(at(10.35)).toBeCloseTo(0.6);
    expect(at(10.7)).toBeCloseTo(0.2);
    expect(at(11.5)).toBeCloseTo(0.2);
    expect(at(12)).toBeCloseTo(0.2);
    expect(at(12.35)).toBeCloseTo(0.6);
    expect(at(12.7)).toBeCloseTo(1);
    expect(at(20)).toBe(1);
  });

  it('comes back from where it got to when a cue ends early', () => {
    const at = (t: number) => duckGain(t, [{ start: 0, end: 0.35 }]);
    expect(at(0.35)).toBeCloseTo(0.6);
    expect(at(0.7)).toBeCloseTo(0.8);
    expect(at(1.05)).toBeCloseTo(1);
  });

  it('follows the lowest of several cues', () => {
    const cues = [
      { start: 0, end: 0.7 },
      { start: 1.05 },
    ];
    expect(duckGain(1.05, cues)).toBeCloseTo(0.6);
    expect(duckGain(1.4, cues)).toBeCloseTo(0.6);
    expect(duckGain(1.75, cues)).toBeCloseTo(0.2);
    expect(duckGain(30, cues)).toBeCloseTo(0.2);
  });

  // Ducking is the mixer's gain on the rest (ADR 022): the cue keeps the
  // player's volume; every other sound plays at it times the duck.
  it('ducks the other sounds while End turn speaks, not the cue', () => {
    vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval', 'setTimeout', 'clearTimeout'] });
    loadSounds({ ...all, endTurnCue: ['/cardimg/assets/sounds/end-turn-cue.mp3'] });
    const volumes: number[] = [];
    const elements: FakeAudio[] = [];
    fakeAudio(undefined, volumes, elements);
    setSoundVolume(50);
    playSound('attackHit');
    playSound('endTurnCue');
    const hitting = elements.at(0);
    const speaking = elements.at(1);
    if (!hitting || !speaking) {
      throw new Error('no sounds');
    }
    expect(volumes).toEqual([0.5, 0.5]);
    const step = (ms: number) => {
      clock += ms;
      vi.advanceTimersByTime(20);
    };
    step(350);
    expect(hitting.volume).toBeCloseTo(0.3);
    step(350);
    expect(hitting.volume).toBeCloseTo(0.1);
    // A sound that starts during the duck starts ducked.
    playSound('diceRoll');
    expect(volumes.at(-1)).toBeCloseTo(0.1);
    expect(speaking.volume).toBe(0.5);
    speaking.onended?.();
    step(350);
    expect(elements.at(-1)?.volume).toBeCloseTo(0.3);
    step(350);
    expect(elements.at(-1)?.volume).toBeCloseTo(0.5);
    // The duck is over: the mixer stops its ticks.
    expect(vi.getTimerCount()).toBe(0);
    // Back at full level, new sounds play at the player's volume again.
    playSound('soulGain');
    expect(volumes.at(-1)).toBe(0.5);
    setSoundVolume(100);
    vi.useRealTimers();
  });

  it('ducks nothing when the cue cannot play', async () => {
    loadSounds({ ...all, endTurnCue: ['/cardimg/assets/sounds/end-turn-cue.mp3'] });
    const volumes: number[] = [];
    fakeAudio(new DOMException('blocked', 'NotAllowedError'), volumes);
    setSoundVolume(100);
    playSound('endTurnCue');
    await Promise.resolve();
    await Promise.resolve();
    clock += 1000;
    playSound('attackHit');
    expect(volumes).toEqual([1, 1]);
  });
});

// signal is a test loop: `quiet` silent samples at each end around
// `loud` samples of a square wave.
function signal(quiet: number, loud: number): Float32Array {
  const data = new Float32Array(quiet * 2 + loud);
  for (let i = 0; i < loud; i++) {
    data[quiet + i] = i % 2 === 0 ? 0.5 : -0.5;
  }
  return data;
}

describe('loopBounds', () => {
  it('skips encoder silence at either end', () => {
    expect(loopBounds(signal(100, 800), 1000)).toEqual({ start: 0.1, end: 0.9 });
  });

  it('keeps a loop that has none', () => {
    expect(loopBounds(signal(0, 1000), 1000)).toEqual({ start: 0, end: 1 });
  });

  it('keeps a silent file whole', () => {
    expect(loopBounds(new Float32Array(500), 1000)).toEqual({ start: 0, end: 0.5 });
  });
});

class FakeParam {
  value = 0;
  ramps: [number, number][] = [];
  cancelScheduledValues = (): void => undefined;
  setValueAtTime(v: number): void {
    this.value = v;
  }
  linearRampToValueAtTime(v: number, at: number): void {
    this.ramps.push([v, at]);
  }
}

class FakeGain {
  gain = new FakeParam();
  connected = true;
  connect = (): void => undefined;
  disconnect(): void {
    this.connected = false;
  }
}

class FakeSource {
  buffer: unknown;
  loop = false;
  loopStart = 0;
  loopEnd = 0;
  started: number[] = [];
  stopped = -1;
  onended: (() => void) | null = null;
  connect = (): void => undefined;
  start(when: number, offset: number): void {
    this.started = [when, offset];
  }
  stop(at: number): void {
    this.stopped = at;
  }
}

class FakeContext {
  static last: FakeContext | undefined;
  state = 'running';
  currentTime = 2;
  destination = {};
  gains: FakeGain[] = [];
  sources: FakeSource[] = [];
  constructor() {
    FakeContext.last = this;
  }
  createGain(): FakeGain {
    const g = new FakeGain();
    this.gains.push(g);
    return g;
  }
  createBufferSource(): FakeSource {
    const s = new FakeSource();
    this.sources.push(s);
    return s;
  }
  decodeAudioData(): Promise<unknown> {
    const data = signal(100, 800);
    return Promise.resolve({ sampleRate: 1000, getChannelData: () => data });
  }
  resume(): Promise<void> {
    return Promise.resolve();
  }
}

// Loops run last: the first one makes the mixer, which then stays.
describe('playLoop', () => {
  const bed = '/cardimg/assets/sounds/drag-bed.mp3';

  it('is silent without the private files or Web Audio', () => {
    loadSounds({ ...none, dragBed: [bed] });
    const quiet = playLoop('dragAir');
    quiet.set(0.5, 0.1);
    quiet.stop(0.2);
    // No AudioContext in Node: a file alone plays nothing either.
    playLoop('dragBed').stop(0.2);
  });

  // The levels are mixer gains (ADR 022): a gain node per loop, the
  // player's volume on the master.
  it('loops past the silence, at its mixer gain, and fades out', async () => {
    vi.stubGlobal('AudioContext', FakeContext);
    vi.stubGlobal('fetch', () => Promise.resolve({ ok: true, arrayBuffer: () => Promise.resolve(new ArrayBuffer(8)) }));
    loadSounds({ ...none, dragBed: [bed] });
    setSoundVolume(40);
    const loop = playLoop('dragBed');
    const ctx = FakeContext.last;
    if (!ctx) {
      throw new Error('no context');
    }
    await vi.waitFor(() => {
      expect(ctx.sources).toHaveLength(1);
    });
    // The master gain, the rest's gain a key cue ducks (20.14), the loop's.
    const master = ctx.gains.at(0);
    const rest = ctx.gains.at(1);
    const gain = ctx.gains.at(2);
    const source = ctx.sources.at(0);
    if (!master || !rest || !gain || !source) {
      throw new Error('no nodes');
    }
    expect(master.gain.value).toBe(0.4);
    expect(rest.gain.value).toBe(1);
    expect(source.loop).toBe(true);
    expect([source.loopStart, source.loopEnd]).toEqual([0.1, 0.9]);
    expect(source.started).toEqual([0, 0.1]);
    loop.set(0.15, 0.5);
    expect(gain.gain.ramps).toEqual([[0.15, 2.5]]);
    loop.stop(0.2);
    expect(gain.gain.ramps.at(-1)).toEqual([0, 2.2]);
    expect(source.stopped).toBe(2.2);
    source.onended?.();
    expect(gain.connected).toBe(false);
    // A stopped loop takes no more gain.
    loop.set(0.5, 0.1);
    expect(gain.gain.ramps).toHaveLength(2);
    setSoundVolume(100);
    expect(master.gain.value).toBe(1);
  });
});
