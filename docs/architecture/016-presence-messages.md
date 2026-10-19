# 16. Presence messages

* **Status:** Accepted
* **Date:** 2026-10-20
* **Authors:** @HardDie

---

## Context

1. Others should see the active player aim and hover live (N-15).
   1. The arrow while it moves (T-08).
   2. The hovered or aimed-at card (T-10, T-21).
   3. A hidden card as a back from its hand slot (T-09, T-22).
2. Intents change the game; events report changes (ADR 006).
3. Records hold only applied intents and their events (ADR 007).
4. Hidden cards never leave the server (A-07, A-16).
5. A pointer moves many times a second.

## Considered options

1. **A new message kind, relayed by the server room**
   1. Won. The engine and the record stay untouched.
2. **Engine intents for aim and hover**
   1. Lost. Every move would be a step, a checksum and a record line.
3. **Players' own pointer positions in screen pixels**
   1. Lost. Each seat lays out the table differently.
4. **Spots by game object, the board frame for the rest**
   1. Won. Every seat can find a card, player, deck or stack item.
5. **Send the hand card's object and let clients hide it**
   1. Lost. A modified client would learn the card (A-16).

## Decision

Use options 1 and 4.

1. Message
   1. Type `presence`, both ways, same shape.
   2. `from`: what aims; `at`: what the arrow or pointer is over.
   3. `point`: the pointer on the board frame, 0 to 1 each way.
   4. Each one replaces the sender's last; an empty one means nothing.
2. Spot
   1. Exactly one of: `object`, `hand`, `player`, `deck`, `stack_id`.
   2. `object` only for a card every player can see.
   3. `hand` is the sender's own hand slot, from 1.
3. Server
   1. Only a seated player, only while the game runs.
   2. The room fills in `seat`; the client's value is ignored.
   3. It sends it to every other client of the game, not the sender.
   4. Spots are checked against the spectator view of this step.
   5. A spot that fails the check becomes no spot.
   6. A point is moved inside the board frame.
   7. At most `MaxPresence` (20) a second per seat; the rest are dropped.
   8. It is never answered, not even with an error.
   9. A dropped connection sends the others an empty presence of its seat.
4. Engine, records and saves never see it.
5. Spectators get it too, once step 9 seats them (owner, 2026-10-20).
6. Adding the type keeps protocol version 3.
   1. An older client ignores a type it does not know.
7. The operator may turn it off.
   1. `-presence=false`, or `"presence": false` in the config file.
   2. Then the server drops every presence and relays none.
   3. `welcome.presence` says it is on; a client sends none otherwise.
   4. On by default; the app's own host always has it on.
8. Client (15.36)
   1. An aim is always sent: `from`, and `at` or `point`.
   2. A hover is sent only on the sender's own turn (owner, 2026-10-20).
   3. A waiting player's hover is not sent: the table stays quiet.
   4. At most one every 70 ms, under the server's limit.
   5. A still presence is sent again every 2 s.
   6. A receiver drops a presence not renewed in 5 s.
   7. Others' arrows sit under the windows (layer `presence`).

```json
{"type": "presence", "id": 17, "data": {"seat": 0, "from": {"hand": 2}, "point": {"x": 0.25, "y": 0.5}}}
{"type": "presence", "data": {"seat": 1, "from": {"object": 314}, "at": {"player": 0}}}
```

## Consequences

### Positive

1. Live aiming costs the engine and the record nothing.
2. A hidden card cannot leak through presence, whatever the client sends.
3. Lost or late messages heal with the next one.

### Negative and risks

1. The room builds a spectator view once per step for the check.
   1. A server short of CPU or bandwidth turns presence off.
2. Replays never show where a player aimed.
3. A spot between steps may name a card that just left; clients ignore it.

### Neutral

1. Clients draw it since 15.36; see decision 8.
