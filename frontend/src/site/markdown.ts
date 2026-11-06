// A small Markdown reader for the rules text (ADR 026, 5): headings,
// numbered and bulleted lists (nested by indent), paragraphs, tables,
// bold, italics, code, links and rule IDs. It builds a tree; RulesText.tsx
// draws it as React elements, so no HTML string ever reaches the page.

export type Inline =
  | { kind: 'text'; text: string }
  | { kind: 'bold'; children: Inline[] }
  | { kind: 'em'; children: Inline[] }
  | { kind: 'code'; text: string }
  | { kind: 'link'; href: string; children: Inline[] }
  | { kind: 'rule'; id: string };

export interface Item {
  // rule is the ID the item defines: a **R-XXX-NN** at its start.
  rule?: string;
  text: Inline[];
  children: Block[];
}

export type Block =
  | { kind: 'heading'; level: number; text: Inline[] }
  | { kind: 'para'; text: Inline[] }
  | { kind: 'list'; ordered: boolean; start: number; items: Item[] }
  | { kind: 'table'; head: Inline[][]; rows: Inline[][][] };

// ruleId matches a rule ID anywhere in text: R-STACK-03.
export const ruleId = /R-[A-Z]+-\d+/;

const itemLine = /^( *)(?:(\d+)\.|[-*]) +(.*)$/;
const headingLine = /^(#{1,6}) +(.*)$/;
const tableLine = /^\s*\|.*\|\s*$/;

// parseMarkdown reads a whole file.
export function parseMarkdown(src: string): Block[] {
  const lines = src.replace(/\r\n/g, '\n').split('\n');
  const blocks: Block[] = [];
  let i = 0;
  while (i < lines.length) {
    const line = lines.at(i) ?? '';
    if (line.trim() === '') {
      i++;
      continue;
    }
    const h = headingLine.exec(line);
    if (h) {
      blocks.push({ kind: 'heading', level: h[1].length, text: parseInline(h[2]) });
      i++;
      continue;
    }
    if (itemLine.test(line)) {
      const got = parseList(lines, i);
      blocks.push(got.list);
      i = got.next;
      continue;
    }
    if (tableLine.test(line)) {
      const rows: string[] = [];
      while (i < lines.length && tableLine.test(lines.at(i) ?? '')) {
        rows.push(lines.at(i) ?? '');
        i++;
      }
      blocks.push(table(rows));
      continue;
    }
    const para: string[] = [];
    while (i < lines.length) {
      const l = lines.at(i) ?? '';
      if (l.trim() === '' || headingLine.test(l) || itemLine.test(l) || tableLine.test(l)) {
        break;
      }
      para.push(l.trim());
      i++;
    }
    blocks.push({ kind: 'para', text: parseInline(para.join(' ')) });
  }
  return blocks;
}

function indentOf(line: string): number {
  return line.length - line.trimStart().length;
}

// parseList reads the list that starts at line i, with its nested lists:
// an item indented deeper than its list belongs to the item above.
function parseList(lines: string[], start: number): { list: Block & { kind: 'list' }; next: number } {
  const first = itemLine.exec(lines.at(start) ?? '');
  const indent = (first?.[1] ?? '').length;
  const ordered = first?.[2] !== undefined;
  const list: Block & { kind: 'list' } = { kind: 'list', ordered, start: Number(first?.[2] ?? 1), items: [] };
  let i = start;
  let current: { item: Item; words: string[] } | undefined;
  const close = () => {
    if (current) {
      finish(current.item, current.words.join(' '));
      list.items.push(current.item);
      current = undefined;
    }
  };
  while (i < lines.length) {
    const line = lines.at(i) ?? '';
    if (line.trim() === '') {
      // A blank line ends the list unless the list goes on after it.
      let j = i + 1;
      while (j < lines.length && (lines.at(j) ?? '').trim() === '') {
        j++;
      }
      const after = lines.at(j) ?? '';
      if (j < lines.length && itemLine.test(after) && indentOf(after) >= indent) {
        i = j;
        continue;
      }
      break;
    }
    const m = itemLine.exec(line);
    const at = indentOf(line);
    if (m && at === indent) {
      close();
      current = { item: { text: [], children: [] }, words: [m[3]] };
      i++;
    } else if (m && at > indent && current) {
      const sub = parseList(lines, i);
      current.item.children.push(sub.list);
      i = sub.next;
    } else if (!m && at > indent && current) {
      current.words.push(line.trim());
      i++;
    } else {
      break;
    }
  }
  close();
  return { list, next: i };
}

// finish parses an item's text; a bold rule ID at its start makes it
// the anchor of that rule.
function finish(item: Item, text: string): void {
  const m = /^\*\*(R-[A-Z]+-\d+)\*\* */.exec(text);
  if (m?.[1]) {
    item.rule = m[1];
    item.text = parseInline(text.slice(m[0].length));
  } else {
    item.text = parseInline(text);
  }
}

function cells(row: string): string[] {
  return row.trim().replace(/^\|/, '').replace(/\|$/, '').split('|').map((c) => c.trim());
}

function table(rows: string[]): Block {
  const [head = '', ...rest] = rows;
  const body = rest.filter((r) => !/^[\s|:-]+$/.test(r));
  return { kind: 'table', head: cells(head).map(parseInline), rows: body.map((r) => cells(r).map(parseInline)) };
}

// Inline marks, tried at each position in this order. A word mark
// only starts a word: "XR-A-1" holds no rule ID.
const marks: { re: RegExp; word?: true; make: (m: RegExpExecArray) => Inline }[] = [
  { re: /^\*\*(.+?)\*\*/, make: (m) => ({ kind: 'bold', children: parseInline(m[1]) }) },
  { re: /^\*(?!\s)(.+?)\*/, make: (m) => ({ kind: 'em', children: parseInline(m[1]) }) },
  { re: /^`([^`]+)`/, make: (m) => ({ kind: 'code', text: m[1] }) },
  { re: /^\[([^\]]+)\]\(([^)\s]+)\)/, make: (m) => ({ kind: 'link', href: m[2], children: parseInline(m[1]) }) },
  { re: /^R-[A-Z]+-\d+\b/, word: true, make: (m) => ({ kind: 'rule', id: m[0] }) },
];

// parseInline reads one line of text.
export function parseInline(src: string): Inline[] {
  const out: Inline[] = [];
  let text = '';
  let i = 0;
  while (i < src.length) {
    const rest = src.slice(i);
    const wordStart = i === 0 || !/[A-Za-z0-9-]/.test(src.charAt(i - 1));
    let hit: { len: number; node: Inline } | undefined;
    for (const { re, word, make } of marks) {
      if (word && !wordStart) {
        continue;
      }
      const m = re.exec(rest);
      if (m) {
        hit = { len: m[0].length, node: make(m) };
        break;
      }
    }
    if (hit) {
      if (text !== '') {
        out.push({ kind: 'text', text });
        text = '';
      }
      out.push(hit.node);
      i += hit.len;
    } else {
      text += src.charAt(i);
      i++;
    }
  }
  if (text !== '') {
    out.push({ kind: 'text', text });
  }
  return out;
}

// anchorsOf lists the rule IDs the blocks define, in order.
export function anchorsOf(blocks: Block[]): string[] {
  return blocks.flatMap((b) => (b.kind === 'list' ? b.items.flatMap((it) => [...(it.rule ? [it.rule] : []), ...anchorsOf(it.children)]) : []));
}

// plainText is the text of inline nodes, without marks.
export function plainText(nodes: Inline[]): string {
  return nodes
    .map((n) => {
      switch (n.kind) {
        case 'text':
        case 'code':
          return n.text;
        case 'rule':
          return n.id;
        default:
          return plainText(n.children);
      }
    })
    .join('');
}
