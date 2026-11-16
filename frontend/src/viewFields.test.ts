// The client reads every field of the game update (TS-07, 22.4). Go's
// TestViewFieldsInProtocolTS checks that protocol.ts declares each field
// the server sends in an Update; this test checks that the app reads
// each one, or lists why not. Lobby messages are out of scope.
//
// A read is the field's name in a property access (x.copy, x?.copy),
// an element access with a literal (x['copy']) or a destructuring
// ({ copy }) in app code. It goes by name, not by type: a read of
// `id` on any object counts for every `id`. Reads in protocol.ts, in
// tests and in the dev-only fakes do not count.
import ts from 'typescript';
import { describe, expect, it } from 'vitest';
import protocol from './protocol.ts?raw';

const sources = import.meta.glob<string>(['./**/*.{ts,tsx}', '!./**/*.test.ts', '!./**/*.d.ts', '!./protocol.ts', '!./fakeBridge.ts', '!./fakeReplay.ts'], {
  query: '?raw',
  import: 'default',
  eager: true,
});

// unread lists the fields the app does not read, by interface and
// field, with why. A gap names its 19.N item until it is fixed; the
// rest stay unread on purpose.
const unread: Record<string, string> = {
  // The server drops a private event for anyone it does not name.
  'GameEvent.seen_by': 'the server filters private events by it',
  // A trigger keeps its event for the engine; the history line before
  // it shows that event.
  'StackItem.event_player': 'the engine resolves the trigger with it',
  'StackItem.event_amount': 'the engine resolves the trigger with it',
  'StackItem.event_stack': 'the engine resolves the trigger with it',
  'StackItem.event_object': 'the engine resolves the trigger with it',
  'StackItem.event_source': 'the engine resolves the trigger with it',
  'StackItem.event_by': 'the engine resolves the trigger with it',
  'StackItem.reward': 'its label says "reward roll"; the monster card shows the reward',
  'AbilityRef.lingering': 'the engine finds the ability by it; text gives the words',
  'StackItem.roll_result': '19.186: the number a roll result acts on is not shown',
  // The engine's turn bookkeeping; what it leads to reaches the client
  // another way.
  'Turn.attack_rolls': 'the odds and stat sources count it ("first attack roll")',
  'Turn.bonus_loot_used': "what is left is the player's loot_plays",
  'Turn.bonus_attacks_used': "what is left is the player's attacks",
  'Turn.bonus_purchases_used': "what is left is the player's purchases",
  'Turn.random_attack_made': 'the Q-31 duty itself is not in the view',
  'Turn.shop_bought': 'the purchase prompt shows each price (20.11)',
  'Turn.free_shop_item': 'the purchase prompt shows the 0¢ price (20.11)',
  'Turn.monster_died': 'the room replace comes as a prompt',
  'Turn.room_replaced': 'the room replace comes as a prompt',
  'Turn.late_rooms': 'blocked says why such a room waits (R-TURN-18)',
  'Turn.entered': "the engine's own step state",
  'Seat.kicked': '19.187: a kicked player is not shown as gone',
};

// viewFields lists every field reachable from Update in protocol.ts, as
// "Interface.field"; an inline object type is named by its field
// ("Prompt.look.deck").
function viewFields(): string[] {
  const file = ts.createSourceFile('protocol.ts', protocol, ts.ScriptTarget.Latest, true);
  const interfaces = new Map<string, ts.NodeArray<ts.TypeElement>>();
  for (const s of file.statements) {
    if (ts.isInterfaceDeclaration(s)) {
      interfaces.set(s.name.text, s.members);
    }
  }
  const seen = new Set<string>();
  const out: string[] = [];
  const walkMembers = (owner: string, members: ts.NodeArray<ts.TypeElement>) => {
    if (seen.has(owner)) {
      return;
    }
    seen.add(owner);
    for (const m of members) {
      if (ts.isPropertySignature(m) && m.type) {
        const name = m.name.getText(file);
        out.push(`${owner}.${name}`);
        walkType(`${owner}.${name}`, m.type);
      }
    }
  };
  const walkType = (at: string, node: ts.Node) => {
    if (ts.isTypeReferenceNode(node)) {
      const members = interfaces.get(node.typeName.getText(file));
      if (members) {
        walkMembers(node.typeName.getText(file), members);
        return;
      }
    }
    if (ts.isTypeLiteralNode(node)) {
      walkMembers(at, node.members);
      return;
    }
    node.forEachChild((c) => {
      walkType(at, c);
    });
  };
  const root = interfaces.get('Update');
  if (root) {
    walkMembers('Update', root);
  }
  return out;
}

// readNames lists every name the app reads as a property.
function readNames(): Set<string> {
  const names = new Set<string>();
  for (const [path, text] of Object.entries(sources)) {
    const file = ts.createSourceFile(path, text, ts.ScriptTarget.Latest, true, path.endsWith('.tsx') ? ts.ScriptKind.TSX : ts.ScriptKind.TS);
    const visit = (n: ts.Node) => {
      if (ts.isPropertyAccessExpression(n)) {
        names.add(n.name.text);
      } else if (ts.isElementAccessExpression(n) && ts.isStringLiteralLike(n.argumentExpression)) {
        names.add(n.argumentExpression.text);
      } else if (ts.isBindingElement(n) && ts.isObjectBindingPattern(n.parent)) {
        names.add((n.propertyName ?? n.name).getText(file));
      }
      n.forEachChild(visit);
    };
    visit(file);
  }
  return names;
}

describe('view fields', () => {
  const fields = viewFields();
  const read = readNames();
  const name = (field: string) => field.slice(field.lastIndexOf('.') + 1);

  it('finds the update and the app sources', () => {
    expect(fields.length).toBeGreaterThan(100);
    expect(fields).toContain('CardView.copy');
    expect(fields).toContain('Prompt.look.deck');
    expect(Object.keys(sources).length).toBeGreaterThan(50);
  });

  it('the app reads every field it is sent, or lists why not', () => {
    const missed = fields.filter((f) => !read.has(name(f)) && !(f in unread));
    expect(missed, 'read these fields in the app, or list them in unread with a reason').toEqual([]);
  });

  it('lists only fields that are sent and still unread', () => {
    const stale = Object.keys(unread).filter((f) => !fields.includes(f) || read.has(name(f)));
    expect(stale, 'drop these from unread').toEqual([]);
  });
});
