# 17. Secure WebSocket through a reverse proxy

* **Status:** Accepted
* **Date:** 2026-10-21
* **Authors:** @HardDie

---

## Context

1. ADR 006 left encryption for later.
2. Operators want a server on a domain with HTTPS (N-19).
   1. Caddy or a similar proxy gets the certificate on its own.
   2. It serves `wss://` on 443 and passes plain WebSocket on.
3. The client only spoke `ws://host:port/ws`.
4. Players type an address by hand.
   1. They should not need to know the scheme.

## Considered options

1. **TLS in the Go server**
   1. Lost. Certificates, renewal and ports become our code.
   2. A proxy already does it well.
2. **A proxy in front; the client speaks `wss://`**
   1. Won. No server change; any proxy works.
3. **A bare domain means `wss://` only**
   1. Lost. A plain server on a domain would need `:4774` typed.

## Decision

Use option 2.

1. The server stays plain `ws://`; the proxy adds TLS.
2. A typed scheme is kept.
   1. `wss://` and `https://` mean `wss://`, port 443 by default.
   2. `ws://` and `http://` mean `ws://`, port 4774 by default.
3. A bare domain tries `wss://` on 443, then `ws://` on 4774.
   1. Each try waits at most the connect timeout (5 s).
4. Plain `ws://` only, with no `wss://` try:
   1. A typed port.
   2. An IP address.
   3. A name without a dot, e.g. `localhost`.
   4. A LAN name: `.local`, `.lan`, `.home`, `.internal`, `.home.arpa`.
5. The record link follows the socket that opened.
   1. `https://` for `wss://`, `http://` for `ws://`.
6. Code: `wsURLs` in `protocol.ts`, `connectFirst` in `net.ts`.

## Consequences

### Positive

1. A server behind Caddy is joined by its domain alone.
2. Traffic to such a server is encrypted.
3. Old addresses keep working as before.

### Negative and risks

1. A domain whose 443 drops packets waits 5 s before the fallback.
2. Someone in the path can block 443 and force plain `ws://`.
   1. Typing `wss://` turns the fallback off.
3. Still no accounts; the token stays the only proof of a seat.

### Neutral

1. Plain servers stay unencrypted, as in ADR 006.
