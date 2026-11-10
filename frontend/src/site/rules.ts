// The rules text on the site (ADR 026, 5): every file of
// pkg/rules_engine/docs/rules, read by Vite as text and parsed here.
// open-questions.md is research and stays off the site; README.md's
// file table becomes the topic index the page builds itself.
import { anchorsOf, parseMarkdown, plainText, type Block } from './markdown';

const files = import.meta.glob<string>(
  ['../../../pkg/rules_engine/docs/rules/*.md', '!../../../pkg/rules_engine/docs/rules/README.md', '!../../../pkg/rules_engine/docs/rules/open-questions.md'],
  { query: '?raw', import: 'default', eager: true },
);

// order is README.md's topic order; a new file goes last until added.
const order = ['setup', 'cards', 'zones', 'abilities', 'mechanics', 'turn', 'stack', 'attack', 'shop', 'death', 'win', 'challenges'];

export interface Topic {
  // file is the file's name without .md: "cards".
  file: string;
  title: string;
  blocks: Block[];
  // rules are the IDs the file defines, in order.
  rules: string[];
}

function rank(file: string): number {
  const i = order.indexOf(file);
  return i < 0 ? order.length : i;
}

// topicsOf parses rule files given by path and text.
export function topicsOf(sources: Record<string, string>): Topic[] {
  return Object.entries(sources)
    .map(([path, text]) => {
      const file = (path.split('/').pop() ?? path).replace(/\.md$/, '');
      const blocks = parseMarkdown(text);
      const head = blocks.find((b) => b.kind === 'heading' && b.level === 1);
      return {
        file,
        title: head?.kind === 'heading' ? plainText(head.text) : file,
        // The page shows the title in the index and above the file.
        blocks: blocks.filter((b) => b !== head),
        rules: anchorsOf(blocks),
      };
    })
    .sort((a, b) => rank(a.file) - rank(b.file) || a.file.localeCompare(b.file));
}

export const topics: Topic[] = topicsOf(files);

const ruleIds = new Set(topics.flatMap((t) => t.rules));

// hasRule tells whether the site has the rule ID.
export function hasRule(id: string): boolean {
  return ruleIds.has(id);
}

// fileLink is where a link in the rules text points on the site: the
// first rule of another rules file, by rule ID; undefined for a file
// the site does not carry (sources, open questions). Web links stay.
export function fileLink(href: string): { rule: string } | { url: string } | undefined {
  if (/^https?:\/\//.test(href)) {
    return { url: href };
  }
  const file = href.replace(/^\.\//, '').replace(/#.*$/, '');
  const topic = topics.find((t) => `${t.file}.md` === file);
  const first = topic?.rules[0];
  return first ? { rule: first } : undefined;
}
