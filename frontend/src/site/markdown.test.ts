import { createElement } from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { anchorsOf, parseInline, parseMarkdown } from './markdown';
import { RulesText } from './RulesText';
import { fileLink, topics, topicsOf } from './rules';

// The raw rules files, as rules.ts reads them, to count IDs by hand.
const raw = import.meta.glob<string>(
  ['../../../pkg/rules_engine/docs/rules/*.md', '!../../../pkg/rules_engine/docs/rules/README.md', '!../../../pkg/rules_engine/docs/rules/open-questions.md'],
  { query: '?raw', import: 'default', eager: true },
);

const html = (src: string, target = '') =>
  renderToStaticMarkup(createElement(RulesText, { blocks: parseMarkdown(src), ruleHref: (id: string) => `#/en/rules/${id}`, target }));

describe('parseInline', () => {
  it('reads bold, italics, code, links and rule IDs', () => {
    expect(parseInline('a **b** *c* `d` [e](f.md) R-CARD-08.')).toEqual([
      { kind: 'text', text: 'a ' },
      { kind: 'bold', children: [{ kind: 'text', text: 'b' }] },
      { kind: 'text', text: ' ' },
      { kind: 'em', children: [{ kind: 'text', text: 'c' }] },
      { kind: 'text', text: ' ' },
      { kind: 'code', text: 'd' },
      { kind: 'text', text: ' ' },
      { kind: 'link', href: 'f.md', children: [{ kind: 'text', text: 'e' }] },
      { kind: 'text', text: ' ' },
      { kind: 'rule', id: 'R-CARD-08' },
      { kind: 'text', text: '.' },
    ]);
  });

  it('finds a rule ID only at the start of a word', () => {
    expect(parseInline('XR-CARD-08 (R-ABIL-20)')).toEqual([
      { kind: 'text', text: 'XR-CARD-08 (' },
      { kind: 'rule', id: 'R-ABIL-20' },
      { kind: 'text', text: ')' },
    ]);
  });
});

describe('parseMarkdown', () => {
  const src = [
    '# Title',
    '',
    'Sources: S-OFF',
    'next line.',
    '',
    '## Part',
    '',
    '1. **R-X-01** First, see R-X-02.',
    '2. **R-X-02** Second',
    '    1. Nested one.',
    '    2. Nested two',
    '       goes on.',
    '',
    '3. **R-X-03** *(later sets)* Third.',
    '',
    '| A | B |',
    '|---|---|',
    '| 1 | **2** |',
  ].join('\n');
  const blocks = parseMarkdown(src);

  it('reads headings, paragraphs, nested lists and tables', () => {
    expect(blocks.map((b) => b.kind)).toEqual(['heading', 'para', 'heading', 'list', 'table']);
    expect(blocks[1]).toEqual({ kind: 'para', text: [{ kind: 'text', text: 'Sources: S-OFF next line.' }] });
    const list = blocks.at(3);
    expect(list?.kind === 'list' && list.items.length).toBe(3);
    expect(list?.kind === 'list' && list.items.at(1)?.children.at(0)).toMatchObject({ kind: 'list', start: 1, items: [{}, { text: [{ kind: 'text', text: 'Nested two goes on.' }] }] });
    expect(blocks[4]).toMatchObject({ kind: 'table', rows: [[[{ kind: 'text', text: '1' }], [{ kind: 'bold' }]]] });
  });

  it('a bold rule ID at the start of an item is its anchor', () => {
    expect(anchorsOf(blocks)).toEqual(['R-X-01', 'R-X-02', 'R-X-03']);
  });

  it('draws anchors, rule links and the target, with no raw HTML', () => {
    const out = html(src, 'R-X-02');
    expect(out).toContain('<li id="rule-R-X-01" class="rule"><a class="rule-id" href="#/en/rules/R-X-01">R-X-01</a>');
    expect(out).toContain('<a class="rule-ref" href="#/en/rules/R-X-02">R-X-02</a>');
    expect(out).toContain('id="rule-R-X-02" class="rule rule-target" aria-current="true"');
    expect(out).toContain('<em>(later sets)</em>');
    expect(out).toContain('<table>');
    expect(html('1. <b>x</b> & <script>')).toContain('&lt;b&gt;x&lt;/b&gt; &amp; &lt;script&gt;');
  });

  it('links between rule files stay on the site; unknown files are text', () => {
    expect(fileLink('cards.md')).toEqual({ rule: 'R-CARD-01' });
    expect(fileLink('https://example.com/x')).toEqual({ url: 'https://example.com/x' });
    expect(fileLink('../sources.md')).toBeUndefined();
    expect(html('See [cards](cards.md) and [sources](../sources.md).')).toBe('<p>See <a href="#/en/rules/R-CARD-01">cards</a> and <span>sources</span>.</p>');
  });
});

describe('the real rules files', () => {
  it('leave out open-questions.md and README.md', () => {
    expect(topics.map((t) => t.file)).not.toContain('open-questions');
    expect(topics.map((t) => t.file)).not.toContain('README');
    expect(topics.map((t) => t.file).slice(0, 2)).toEqual(['setup', 'cards']);
    expect(topicsOf(raw).length).toBe(topics.length);
  });

  it('give every rule ID exactly one anchor', () => {
    const defined = Object.values(raw).flatMap((text) => [...text.matchAll(/^\s*\d+\.\s+\*\*(R-[A-Z]+-\d+)\*\*/gm)].map((m) => m[1]));
    const anchors = topics.flatMap((t) => t.rules);
    expect(anchors.length).toBeGreaterThan(250);
    expect([...anchors].sort()).toEqual([...defined].sort());
    expect(new Set(anchors).size).toBe(anchors.length);
    const page = topics.map((t) => renderToStaticMarkup(createElement(RulesText, { blocks: t.blocks, ruleHref: (id: string) => id, target: '' }))).join('');
    for (const id of anchors) {
      expect(page.split(`id="rule-${id}"`).length - 1, id).toBe(1);
    }
  });

  it('link only to rules they define', () => {
    const anchors = new Set(topics.flatMap((t) => t.rules));
    const mentioned = Object.values(raw).flatMap((text) => [...text.matchAll(/\bR-[A-Z]+-\d+\b/g)].map((m) => m[0]));
    expect(mentioned.filter((id) => !anchors.has(id))).toEqual([]);
  });
});
