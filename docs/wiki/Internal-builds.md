# Internal builds

An internal build is the app for a play test with friends: card images included, the version naming the commit (`internal-1a2b3c4`, shown in the main menu corner). It needs the private card images submodule (`pkg/card_db/images`).

```sh
make internal                          # for this machine
make internal PLATFORM=windows/amd64   # a Windows .exe, from macOS or Linux
make internal TAGS=webkit2_41          # on Ubuntu 24.04
```

The result is in `build/bin/`: `FourSouls.app` on macOS, `FourSouls.exe` for Windows, `FourSouls` on Linux. A Linux build has to be made on Linux (WebKitGTK is linked natively); macOS builds need a Mac.

## Sharing

The build contains the licensed card images. Share it only with the people in the play test, never publicly (see `NOTICE` and the card images README).

* macOS: zip `FourSouls.app`. The build is not signed with a developer ID, so on first start use right click → Open.
* Windows: send `FourSouls.exe`; SmartScreen may ask to confirm. The PC needs the WebView2 runtime (built into Windows 11).

## A LAN play test

1. One player hosts: **Play → LAN → Host a game**, players and setup, **Host**. The table shows the addresses to join.
2. The others: **Play → LAN**, then **Join** next to the game on the list, or type one of those addresses and **Connect**; then **Join** the game.
3. Everyone presses **Ready**; the game starts.

Over the internet, use a virtual LAN: see [Playing over the internet](Playing-over-the-internet).

## What the browser tests cover

`make e2e` runs every UI test in Chromium and in WebKit. WebKit is the engine of macOS's WKWebView and the nearest check for Linux's WebKitGTK; it already caught a connection that never failed on a blocked port. Speed on Linux still needs a real Linux machine.
