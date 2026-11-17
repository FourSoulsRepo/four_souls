// Dev only, for the fake bridge (?fake): the steps of a fake record, as
// Go's Records.Step gives them (10.2). A turn every ten steps; the
// active player gains a cent at step 3 of each turn; a finished record
// ends with its winner. Views hold the Base Game's real cards, so the
// table draws them as text or with the private images.

const characters = ['isaac', 'cain', 'judas', 'eve'];
const hands = [['a_penny', 'two_cents'], ['bomb', 'three_cents', 'a_nickel'], ['four_cents'], ['a_dime', 'soul_heart']];

// FakeRecording is what a fake record's steps are made from.
export interface FakeRecording {
  players: number;
  steps: number;
  winners: number[]; // seats; none while unfinished
}

// stepsPerTurn is how long each fake turn is.
const stepsPerTurn = 10;

// fakeStep is live step index of a fake record, or undefined past its end.
export function fakeStep(r: FakeRecording, index: number): object | undefined {
  if (index < 0 || index >= r.steps) {
    return undefined;
  }
  const turn = Math.floor(index / stepsPerTurn);
  const active = turn % r.players;
  const last = index === r.steps - 1 && r.winners.length > 0;
  const events: object[] = [];
  if (index % stepsPerTurn === 0 && index > 0) {
    events.push({ kind: 'turn_started', player: active });
  }
  if (index % stepsPerTurn === 3) {
    events.push({ kind: 'gained_cents', player: active, amount: 1 });
  }
  if (last) {
    events.push(...r.winners.map((w) => ({ kind: 'game_won', player: w })));
  }
  return { index, n: index + 1, events, view: fakeView(r, index, last) };
}

// fakeView is the judge view after step index: every hand shows.
function fakeView(r: FakeRecording, index: number, over: boolean): object {
  const turn = Math.floor(index / stepsPerTurn);
  let id = 100;
  const monster = { hp: 2, dc: 3, atk: 1, printed_hp: 2, printed_dc: 3, printed_atk: 1 };
  const card = (ref: string, more: object = {}) => ({ id: ++id, card: ref, charged: true, ...more });
  const players = Array.from({ length: r.players }, (_, p) => {
    // Each turn of this player gives a cent at its step 3.
    let cents = 3;
    for (let t = 0; t <= turn; t++) {
      if (t % r.players === p && (t < turn || index % stepsPerTurn >= 3)) {
        cents++;
      }
    }
    const hand = (hands.at(p) ?? []).map((ref) => card(ref));
    return {
      id: p,
      character: card(characters.at(p) ?? 'isaac', { role: 1 }),
      hp: 2,
      atk: 1,
      cents,
      souls: over && r.winners.includes(p) ? 4 : 0,
      in_play: [],
      hand_size: hand.length,
      hand,
      loot_plays: 1,
      attacks: 1,
      purchases: 1,
    };
  });
  return {
    players,
    shop: [{ top: card('breakfast') }, { top: card('the_d4') }],
    monsters: [{ top: card('big_spider', { stats: monster }) }, { top: card('clotty', { stats: monster }) }],
    deck_sizes: [80, 60, 90, 0],
    discards: [null, null, null, null],
    discard_sizes: [0, 0, 0, 0],
    stack: null,
    turn: { active: turn % r.players, number: turn + 1, step: 4, loot_plays: 1, attacks: 1, purchases: 1 },
    waiting: { kind: over ? 3 : 1, player: turn % r.players },
    win_souls: 4,
    ...(over ? { over: true, winners: r.winners } : {}),
  };
}
