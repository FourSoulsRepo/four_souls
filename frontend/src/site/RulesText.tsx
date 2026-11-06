// Draws the tree of markdown.ts as React elements: React escapes every
// text, so the rules files can never inject HTML (ADR 026, 5).
import type { ReactNode } from 'react';
import { fileLink } from './rules';
import type { Block, Inline, Item } from './markdown';

interface Links {
  // ruleHref is the site link of a rule ID.
  ruleHref: (id: string) => string;
  // target is the rule the route points at; it is highlighted.
  target: string;
}

// keyed pairs each part of a fixed text with its place as a key: the
// rules text never reorders, so the place is a stable key.
function keyed<T>(list: readonly T[]): [string, T][] {
  return list.map((x, i) => [String(i), x]);
}

function inlines(nodes: Inline[], links: Links): ReactNode[] {
  return keyed(nodes).map(([key, n]) => {
    switch (n.kind) {
      case 'text':
        return n.text;
      case 'bold':
        return <strong key={key}>{inlines(n.children, links)}</strong>;
      case 'em':
        return <em key={key}>{inlines(n.children, links)}</em>;
      case 'code':
        return <code key={key}>{n.text}</code>;
      case 'rule':
        return (
          <a key={key} className="rule-ref" href={links.ruleHref(n.id)}>
            {n.id}
          </a>
        );
      case 'link': {
        const to = fileLink(n.href);
        const label = inlines(n.children, links);
        if (to === undefined) {
          return <span key={key}>{label}</span>;
        }
        return 'url' in to ? (
          <a key={key} href={to.url} target="_blank" rel="noopener noreferrer">
            {label}
          </a>
        ) : (
          <a key={key} href={links.ruleHref(to.rule)}>
            {label}
          </a>
        );
      }
    }
  });
}

function item(it: Item, key: string, links: Links): ReactNode {
  const body = (
    <>
      {inlines(it.text, links)}
      {blocks(it.children, links)}
    </>
  );
  if (!it.rule) {
    return <li key={key}>{body}</li>;
  }
  const target = it.rule === links.target;
  return (
    <li key={key} id={`rule-${it.rule}`} className={target ? 'rule rule-target' : 'rule'} aria-current={target ? 'true' : undefined}>
      <a className="rule-id" href={links.ruleHref(it.rule)}>
        {it.rule}
      </a>{' '}
      {body}
    </li>
  );
}

function blocks(list: Block[], links: Links): ReactNode[] {
  return keyed(list).map(([key, b]) => {
    switch (b.kind) {
      case 'heading':
        return b.level <= 2 ? <h3 key={key}>{inlines(b.text, links)}</h3> : <h4 key={key}>{inlines(b.text, links)}</h4>;
      case 'para':
        return <p key={key}>{inlines(b.text, links)}</p>;
      case 'list': {
        const items = keyed(b.items).map(([k, it]) => item(it, k, links));
        return b.ordered ? (
          <ol key={key} start={b.start}>
            {items}
          </ol>
        ) : (
          <ul key={key}>{items}</ul>
        );
      }
      case 'table':
        return (
          <div key={key} className="table-wrap">
            <table>
              <thead>
                <tr>
                  {keyed(b.head).map(([k, c]) => (
                    <th key={k}>{inlines(c, links)}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {keyed(b.rows).map(([k, row]) => (
                  <tr key={k}>
                    {keyed(row).map(([c, cell]) => (
                      <td key={c}>{inlines(cell, links)}</td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        );
    }
  });
}

// RulesText draws parsed rules text; rule IDs link with ruleHref.
export function RulesText({ blocks: list, ruleHref, target }: { blocks: Block[] } & Links) {
  return <>{blocks(list, { ruleHref, target })}</>;
}
