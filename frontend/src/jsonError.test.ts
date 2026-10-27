import { describe, expect, it } from 'vitest';
import { explainJson, jsonErrorAt, jsonErrorText, lineColumn, parseJson } from './jsonError';

describe('jsonErrorAt', () => {
  it.each([
    ['{}', -1],
    ['{"a": "x", "b": [1, -2.5e3, true, false, null, {"c": "\\"\\u00e9\\n"}]}', -1],
    [' [ ] ', -1],
    ['', 0],
    ['\uFEFF{}', 0],
    ['{"a": 1,}', 8],
    ['{"a": 1 "b": 2}', 8],
    ['{"a": "x\ny"}', 8],
    ['{a: 1}', 1],
    ["{'a': 1}", 1],
    ['{"a": nope}', 7],
    ['{"a": 01}', 7],
    ['[-]', 2],
    ['[1.]', 3],
    ['[1e+]', 4],
    ['[-0.5E-3, 10]', -1],
    ['{"a": "\\q"}', 8],
    ['{"a": "\\u12g4"}', 11],
    ['{"a": [1, 2', 11],
    ['{} x', 3],
  ])('%j stops at %i', (text, at) => {
    expect(jsonErrorAt(text)).toBe(at);
  });
});

describe('explainJson', () => {
  const hint = (text: string) => explainJson(text, jsonErrorAt(text));
  it.each([
    ['an empty file', '', 0, 'the file is empty'],
    ['spaces only', ' \n', 0, 'the file is empty'],
    ['a BOM', '\uFEFF{"a": "x"}', 0, 'byte order mark'],
    ['a trailing comma before }', '{\n  "a": "x",\n}\n', 12, 'a comma before }: delete this comma'],
    ['a trailing comma before ]', '{"a": [1, 2, ]}', 11, 'a comma before ]'],
    ['a missing comma between lines', '{\n  "a": "x"\n  "b": "y"\n}', 12, 'a comma is missing here'],
    ['a missing comma after a number', '{"a": 1 "b": 2}', 7, 'a comma is missing here'],
    ['a missing comma after an object', '{"a": {} "b": 2}', 8, 'a comma is missing here'],
    ['a missing closing quote', '{\n  "a": "x,\n  "b": "y"\n}', 12, 'closing quote " is missing'],
    ['a tab in a text', '{"a": "x\ty"}', 8, 'a tab inside a text'],
    ['a quote inside a text', '{"a": "say "hi" now"}', 12, 'write \\" for a "'],
    ['a key without the closing quote', '{"a: "x"}', 6, 'a quote is missing'],
    ['a key without quotes', '{\n  a: "x"\n}', 4, 'a key needs double quotes'],
    ['a key without quotes after a comma', '{"a": "x", b: "y"}', 11, 'a key needs double quotes'],
    ['a text without quotes', '{"a": hello}', 6, 'a text needs double quotes'],
    ['a misspelt null', '{"a": nope}', 6, 'a text needs double quotes'],
    ['single quotes', "{'a': 'x'}", 1, 'double quotes " only'],
    ['single quotes for a text', "{\"a\": 'x'}", 6, 'double quotes " only'],
    ['the end comes too early', '{"a": "x"', 9, 'ends too early'],
    ['text after the end', '{"a": "x"}\n}', 11, 'more text after the last }'],
    ['no object', 'hello', 0, 'starts with {'],
    ['anything else', '{"a": "x",, "b": "y"}', 10, '"," does not fit here'],
  ])('%s', (_, text, offset, part) => {
    const p = hint(text);
    expect(p.hint).toContain(part);
    expect(p.offset).toBe(offset);
  });
});

describe('lineColumn', () => {
  it.each([
    ['abc', 0, 1, 1],
    ['abc', 2, 1, 3],
    ['a\nbc', 2, 2, 1],
    ['a\nbc', 4, 2, 3],
  ])('%j at %i is %i:%i', (text, offset, line, column) => {
    expect(lineColumn(text, offset)).toEqual({ line, column });
  });
});

describe('jsonErrorText', () => {
  it('names the file, line and column and puts a caret under the place', () => {
    const text = '{\n  "a": "x",\n  "b": "y",\n}\n';
    expect(jsonErrorText('de.json', text).split('\n')).toEqual([
      'de.json:3:11: JSON error: a comma before }: delete this comma; the last entry has none',
      '  3 |   "b": "y",',
      '    |           ^',
    ]);
  });

  it('keeps tabs in the caret line, so it lines up', () => {
    expect(jsonErrorText('de.json', '{\n\t"a": 1\n\t"b": 2\n}').split('\n').slice(1)).toEqual(['  2 | \t"a": 1', '    | \t      ^']);
  });

  it('falls back to the parser position when the scan finds none', () => {
    expect(jsonErrorText('de.json', '{}', 'Unexpected token in JSON at position 1')).toContain('de.json:1:2:');
  });
});

describe('parseJson', () => {
  it('gives the value of good JSON', () => {
    expect(parseJson('de.json', '{"a": "x"}')).toEqual({ value: { a: 'x' } });
  });

  it('gives the message for bad JSON', () => {
    expect(parseJson('de.json', '{"a": "x",}')).toEqual({ error: expect.stringContaining('de.json:1:10: JSON error: a comma before }') as unknown });
  });
});
