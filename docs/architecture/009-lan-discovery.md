# 9. LAN game discovery

* **Status:** Accepted
* **Date:** 2026-10-12
* **Authors:** @HardDie

---

## Context

1. Joining a LAN game needs the host's IP today (N-13).
   1. Players read it off the host's screen and type it.
2. Goal: the join screen lists games on the network.
   1. Typing an address keeps working.
3. Only an app that hosts announces.
   1. A dedicated server never does.
   2. Announcements stop when the app has no game left.
4. Firewalls and OS prompts.
   1. Windows asks once per app for listening sockets.
   2. The hosting app already listens on TCP 4774.
   3. macOS 15 asks for "Local Network" access before LAN traffic.
5. Two app copies may run on one computer (UC-16).
6. Few dependencies; each one goes into `NOTICE`.

## Considered options

1. **UDP broadcast, query and answer**
   1. Won. Standard library only; one small message each way.
   2. The joining app sends from a random port and never listens.
   3. Windows lets unicast answers to a broadcast through.
   4. Works over virtual LANs that pass broadcasts (ZeroTier, Hamachi).
2. **UDP broadcast beacons from the host**
   1. Lost. Every joining app must listen on a fixed UDP port.
   2. That is one more firewall prompt, on every player's computer.
   3. Two copies on one computer fight over the port.
3. **mDNS / DNS-SD (Bonjour)**
   1. Lost. Needs a library and multicast group handling.
   2. Shares port 5353 with the OS responder.
   3. macOS also wants `NSBonjourServices` declared.
   4. Nothing it gives matters for a list of a few games.

## Decision

Use option 1.

1. Package `internal/lan`.
2. Port: UDP 4774, the same number as the TCP game port.
3. Probe, joining app → network.
   1. JSON `{"magic":"four-souls","kind":"probe"}`.
   2. Sent to 255.255.255.255, each interface's broadcast, and 127.0.0.1.
4. Answer, hosting app → the probe's sender, unicast.
   1. JSON with `magic`, `kind: "games"`, a server `id`, `protocol`, `app`, `port`.
   2. `games`: id, host nickname, seats, taken, locked, started.
   3. Every game but a finished one: waiting, full, or started (15.97).
   4. A started game lets a player who held a seat rejoin it.
   5. Only sent when the app has at least one such game.
5. The joining app collects answers for one second.
   1. Answers from one server `id` merge; a LAN address beats loopback.
   2. Each found game shows its address; a click connects there.
6. `server.Config.Discovery` turns answering on.
   1. Only the app's `Host` sets it.
   2. It has no JSON name, so a dedicated server cannot turn it on.
7. Messages are capped at 8 KiB; anything else is ignored.
8. macOS: `Info.plist` has `NSLocalNetworkUsageDescription`.

## Consequences

### Positive

1. No new dependency.
2. Joining players get no extra firewall prompt on Windows.
3. Several copies on one computer find each other's games.

### Negative and risks

1. Routers and guest Wi-Fi may block broadcasts; the address field stays.
2. Any computer on the LAN can see that a game waits, and who hosts.
   1. Game passwords (N-14) keep strangers out.
   2. It also sees games in progress; only a seat's token gets in.
3. macOS asks for Local Network access once; a "no" hides LAN games.

### Neutral

1. IPv4 only, like `LocalAddresses`.
2. Two machines on one Wi-Fi must be checked by hand.
