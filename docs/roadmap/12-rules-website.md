# Step 12. Rules website

Goal: check rules in the browser; no server.

Ideas: A-03.

---

### [ ] 12.1 Engine in WebAssembly

1. Goal: `cmd/website` exposes the engine to JS.
2. Tasks:
   1. Function: payload in, answer with reasons out.
   2. Function: engine version.
   3. Measure and report the `.wasm` size.
3. Done when:
   1. A JS test evaluates a sample payload.

### [ ] 12.2 Site UI

1. Goal: a simple page (A-03).
2. Tasks:
   1. Paste payload or open a link.
   2. Show the answer, reasons, rule links.
   3. Rules text with anchors by rule ID.
3. Done when:
   1. A link copied from the app opens the same situation.

### [ ] 12.3 GitHub Pages

1. Goal: the site is online.
2. Tasks:
   1. Workflow builds and deploys on release.
   2. App config gets the site base URL (11.3).
3. Done when:
   1. The public URL loads and answers payloads.
