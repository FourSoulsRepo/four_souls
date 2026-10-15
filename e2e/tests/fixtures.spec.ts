// Guards the clash check (19.6): a spec that imports test from
// @playwright/test would skip it.
import { readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';
import { expect, test } from './fixtures';

test('every spec takes test from fixtures.ts', () => {
  const dir = import.meta.dirname;
  const wrong = readdirSync(dir)
    .filter((f) => f.endsWith('.spec.ts'))
    .filter((f) => /import \{[^}]*\btest\b[^}]*\} from '@playwright\/test'/.test(readFileSync(join(dir, f), 'utf8')));
  expect(wrong).toEqual([]);
});
