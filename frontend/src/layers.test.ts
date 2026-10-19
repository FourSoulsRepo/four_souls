import { describe, expect, it } from 'vitest';
import doc from '../../docs/wiki/Layers.md?raw';
import layers from './layers.css?raw';

// The inventory (19.26, ADR 015): layers.css holds every z-index layer,
// top first; docs/wiki/Layers.md lists them in the same order; the rest
// of the code only names them, or tags a number as local to its element.
describe('layer inventory', () => {
  const sources = import.meta.glob<string>(['./**/*.{ts,tsx,css}', '!./**/*.test.ts', '!./layers.css'], {
    query: '?raw',
    import: 'default',
    eager: true,
  });
  const site = /z-index:|zIndex/;
  const named = /var\(--z-([a-z-]+)\)/;

  // defined are the layers in layers.css, top first, with their numbers.
  const defined = [...layers.matchAll(/--z-([a-z-]+):\s*(-?\d+);/g)].map((m) => ({ id: m[1], z: Number(m[2]) }));

  // docIds are the ids in the first column of the inventory table.
  function docIds(): string[] {
    const table = doc.split('## Inventory')[1]?.split('\n## ')[0] ?? '';
    return [...table.matchAll(/^\| `([a-z-]+)` \|/gm)].map((m) => m[1]);
  }

  function scan(): { used: Set<string>; bare: string[] } {
    const used = new Set<string>();
    const bare: string[] = [];
    for (const [file, text] of Object.entries(sources)) {
      text.split('\n').forEach((line, i) => {
        if (!site.test(line)) {
          return;
        }
        const m = named.exec(line);
        if (m?.[1]) {
          used.add(m[1]);
        } else if (!line.includes('layer: local')) {
          bare.push(`${file}:${String(i + 1)}`);
        }
      });
    }
    return { used, bare };
  }

  it('keeps the layers in order, top first', () => {
    expect(defined.length).toBeGreaterThan(0);
    defined.slice(1).forEach((d, i) => {
      expect(d.z, `--z-${d.id}`).toBeLessThan(defined.at(i)?.z ?? 0);
    });
  });

  it('lists every layer in Layers.md, in the same order', () => {
    expect(docIds()).toEqual(defined.map((d) => d.id));
  });

  it('uses only named layers, or numbers tagged local', () => {
    expect(scan().bare).toEqual([]);
  });

  it('names only layers that exist, and uses each one', () => {
    expect([...scan().used].sort()).toEqual(defined.map((d) => d.id).sort());
  });
});
