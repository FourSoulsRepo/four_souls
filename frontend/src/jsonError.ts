// Plain JSON errors for translators (16.11): where a file stops being
// JSON, and a hint for the usual mistakes. make i18n uses it on the
// language and glossary files; pkg/card_db has the same hints in Go.
// JSON.parse words its errors differently in each Node version, and
// some name no position, so jsonErrorAt finds the place itself.

// JsonProblem is one syntax error: the 0-based offset to show and a hint.
export interface JsonProblem {
  offset: number;
  hint: string;
}

class Stop extends Error {
  constructor(readonly at: number) {
    super(`JSON stops at ${String(at)}`);
  }
}

const space = ' \t\n\r';
const escapes = '"\\/bfnrtu';
const isDigit = (c: string) => c >= '0' && c <= '9';

// jsonErrorAt is the offset of the first character that breaks the JSON
// text, text.length when it ends too early, or -1 when it is JSON.
export function jsonErrorAt(text: string): number {
  let i = 0;
  const stop = (): never => {
    throw new Stop(i);
  };
  const skip = () => {
    while (i < text.length && space.includes(text.charAt(i))) {
      i++;
    }
  };
  const digits = () => {
    const from = i;
    while (isDigit(text.charAt(i))) {
      i++;
    }
    if (i === from) {
      stop();
    }
  };
  // num reads -0.5e3: a sign, the whole part, a fraction, an exponent.
  const num = () => {
    if (text.charAt(i) === '-') {
      i++;
    }
    if (text.charAt(i) === '0') {
      i++;
    } else {
      digits();
    }
    if (text.charAt(i) === '.') {
      i++;
      digits();
    }
    if (text.charAt(i) === 'e' || text.charAt(i) === 'E') {
      i++;
      if (text.charAt(i) === '+' || text.charAt(i) === '-') {
        i++;
      }
      digits();
    }
  };
  const str = () => {
    i++;
    for (;;) {
      const c = text.charAt(i);
      if (c === '' || c < ' ') {
        stop();
      }
      i++;
      if (c === '"') {
        return;
      }
      if (c === '\\') {
        const e = text.charAt(i);
        if (e === '' || !escapes.includes(e)) {
          stop();
        }
        i++;
        if (e === 'u') {
          for (let k = 0; k < 4; k++, i++) {
            if (!/[0-9a-fA-F]/.test(text.charAt(i))) {
              stop();
            }
          }
        }
      }
    }
  };
  // list reads an object or an array after its opening bracket.
  const list = (close: string, item: () => void) => {
    i++;
    skip();
    if (text.charAt(i) === close) {
      i++;
      return;
    }
    for (;;) {
      item();
      skip();
      const c = text.charAt(i);
      if (c !== ',' && c !== close) {
        stop();
      }
      i++;
      if (c === close) {
        return;
      }
    }
  };
  const value = (): void => {
    skip();
    const c = text.charAt(i);
    if (c === '{') {
      list('}', () => {
        skip();
        if (text.charAt(i) !== '"') {
          stop();
        }
        str();
        skip();
        if (text.charAt(i) !== ':') {
          stop();
        }
        i++;
        value();
      });
    } else if (c === '[') {
      list(']', value);
    } else if (c === '"') {
      str();
    } else if (c === '-' || isDigit(c)) {
      num();
    } else {
      const word = ['true', 'false', 'null'].find((w) => w.startsWith(c) && c !== '') ?? stop();
      for (const w of word) {
        if (text.charAt(i) !== w) {
          stop();
        }
        i++;
      }
    }
  };
  try {
    value();
    skip();
    return i < text.length ? i : -1;
  } catch (e) {
    if (e instanceof Stop) {
      return e.at;
    }
    throw e;
  }
}

const isSpace = (c: string) => c !== '' && space.includes(c);
const isLetter = (c: string) => /[\p{L}\p{N}_$]/u.test(c);

// explainJson says what is wrong at offset, the place jsonErrorAt or a
// parser found, in plain words, and which place to show: the comma to
// delete, or where a comma is missing.
export function explainJson(text: string, offset: number): JsonProblem {
  if (text.trim() === '') {
    return { offset: 0, hint: 'the file is empty; an empty language file is {}' };
  }
  if (text.startsWith('\uFEFF')) {
    return { offset: 0, hint: 'the file starts with an invisible byte order mark (BOM); save it as UTF-8 without BOM' };
  }
  if (offset >= text.length) {
    return { offset: text.length, hint: 'the file ends too early: a closing }, ] or quote is missing' };
  }
  let at = offset;
  while (at > 0 && isLetter(text.charAt(at)) && isLetter(text.charAt(at - 1))) {
    at--;
  }
  let before = at - 1;
  while (before >= 0 && isSpace(text.charAt(before))) {
    before--;
  }
  const c = text.charAt(at);
  const prev = before < 0 ? '' : text.charAt(before);
  if (prev === '') {
    return { offset: at, hint: 'a language file starts with { and ends with }' };
  }
  if (jsonErrorAt(text.slice(0, before + 1)) === -1) {
    return { offset: at, hint: 'more text after the last }: delete it, or move it inside' };
  }
  if ((c === '}' || c === ']') && prev === ',') {
    return { offset: before, hint: `a comma before ${c}: delete this comma; the last entry has none` };
  }
  if (c === "'") {
    return { offset: at, hint: 'JSON takes double quotes " only, not single quotes \'' };
  }
  if (c === '\t') {
    return { offset: at, hint: 'a tab inside a text: write \\t, or a space' };
  }
  if (c < ' ') {
    return { offset: at, hint: 'a text runs to the end of the line: its closing quote " is missing (a line break inside a text is \\n)' };
  }
  if ((c === '"' || c === '{' || c === '[') && (prev === '"' || prev === '}' || prev === ']' || isLetter(prev))) {
    return { offset: before + 1, hint: 'a comma is missing here, between two entries' };
  }
  if (isLetter(c) && prev === '"') {
    return { offset: at, hint: 'a quote is missing, or a quote inside a text lacks its backslash: write \\" for a " in a text' };
  }
  if (isLetter(c) && (prev === '{' || prev === ',')) {
    return { offset: at, hint: 'a key needs double quotes: "key": "text"' };
  }
  if (isLetter(c) && (prev === ':' || prev === '[')) {
    return { offset: at, hint: 'a text needs double quotes: "key": "text"' };
  }
  return { offset: at, hint: `${JSON.stringify(c)} does not fit here` };
}

// lineColumn turns an offset into a 1-based line and column.
export function lineColumn(text: string, offset: number): { line: number; column: number } {
  const head = text.slice(0, offset).split('\n');
  return { line: head.length, column: (head.at(-1) ?? '').length + 1 };
}

// jsonErrorText is the message for a file that is no JSON: its path,
// line and column, the hint, and the line with a caret under the place.
// The parser's own position is the fallback, should the scan find none.
export function jsonErrorText(file: string, text: string, parserMessage = ''): string {
  const found = jsonErrorAt(text);
  const p = explainJson(text, found >= 0 ? found : Number(/position (\d+)/.exec(parserMessage)?.[1] ?? 0));
  const { line, column } = lineColumn(text, p.offset);
  const shown = (text.split('\n')[line - 1] ?? '').replace(/\r$/, '');
  const pad = Array.from(shown.slice(0, column - 1), (ch) => (ch === '\t' ? '\t' : ' ')).join('');
  const num = String(line);
  return [`${file}:${num}:${String(column)}: JSON error: ${p.hint}`, `  ${num} | ${shown}`, `  ${' '.repeat(num.length)} | ${pad}^`].join('\n');
}

// parseJson reads a JSON text, or gives the message for its error.
export function parseJson(file: string, text: string): { value: unknown } | { error: string } {
  try {
    return { value: JSON.parse(text) as unknown };
  } catch (e) {
    return { error: jsonErrorText(file, text, e instanceof Error ? e.message : '') };
  }
}
