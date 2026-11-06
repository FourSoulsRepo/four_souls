// Smoke test of the rules website's .wasm (12.1): loads site.wasm.gz the
// way the page will, then calls every exported function. make site runs it.
//
// Usage: node cmd/website/smoke.mjs <wasm_exec.js> <site.wasm.gz>
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';
import { gunzipSync } from 'node:zlib';

const [execJS, wasmGz] = process.argv.slice(2);
if (!execJS || !wasmGz) {
  console.error('usage: node smoke.mjs <wasm_exec.js> <site.wasm.gz>');
  process.exit(2);
}

// The same payloads as internal/website's Go test.
const testdata = join(dirname(fileURLToPath(import.meta.url)), '../../internal/website/testdata');
const payload = (name) => readFileSync(join(testdata, name), 'utf8');

await import(pathToFileURL(execJS).href);
const go = new globalThis.Go();
const { instance } = await WebAssembly.instantiate(gunzipSync(readFileSync(wasmGz)), go.importObject);
void go.run(instance);
const api = globalThis.fourSouls;
assert.ok(api, 'globalThis.fourSouls is not set');

const version = JSON.parse(api.version());
assert.match(version.engine, /^\d+\.\d+\.\d+$/);

const sets = JSON.parse(api.sets());
assert.ok(sets.includes('b2'), `sets: ${sets}`);

const cards = JSON.parse(api.cards());
assert.ok(Array.isArray(cards) && cards.length > 0, `cards: ${JSON.stringify(cards).slice(0, 200)}`);
assert.ok(cards.every((c) => sets.includes(c.set) && c.image === '' && Object.keys(c.images).length === 0));

const notice = JSON.parse(api.notice());
assert.ok(notice.length > 0 && notice[0].parts.length > 0);

const refused = JSON.parse(api.check(payload('loot-needs-a-loot-play.json')));
assert.equal(refused.allowed, false);
assert.equal(refused.rule, 'R-CARD-08');
assert.ok(refused.say?.key, `refusal without a message key: ${JSON.stringify(refused)}`);
assert.equal(refused.matches, true);

const allowed = JSON.parse(api.check(payload('penny-gives-a-cent.json')));
assert.equal(allowed.allowed, true);
assert.equal(allowed.matches, true, JSON.stringify(allowed.mismatch));
assert.ok(allowed.events.length > 0);

assert.ok(JSON.parse(api.check('not json')).error);
assert.ok(JSON.parse(api.check()).error);

console.log(`site.wasm: engine ${version.engine}, ${cards.length} cards, sets ${sets.join(' ')}: ok`);
process.exit(0);
