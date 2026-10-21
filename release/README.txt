Four Souls
==========

An unofficial, free, fan-made desktop version of the card game
"The Binding of Isaac: Four Souls" for 2-4 players over the network.

The game belongs to Edmund McMillen and Maestro Media.
Please buy the physical game: https://foursouls.com


What is in this archive
-----------------------

  FourSouls (.app on macOS, .exe on Windows)    the game
  LICENSE, NOTICE                                licenses

A "for-poor" archive is the same game without card images: cards
show as text. Players with either archive can play together.


First start
-----------

The game is not signed by Apple or Microsoft (it is free and fan-made),
so the system warns you the first time. You only do this once.

macOS
  1. Unzip the archive and move "FourSouls" (the app) to Applications.
  2. Double-click it. macOS says it cannot check the app: press Done
     (or Cancel).
  3. Open System Settings > Privacy & Security. Near the bottom it says
     "FourSouls" was blocked: press "Open Anyway", then "Open Anyway"
     again and enter your password.
  On macOS 14 and older, right-click the app, choose Open, then Open.

Windows
  1. Unzip the archive (right-click > Extract All). Do not run it from
     inside the zip.
  2. Double-click FourSouls.exe. If "Windows protected your PC" shows,
     press "More info", then "Run anyway".
  3. The game needs Microsoft Edge WebView2. Windows 10 and 11 usually
     have it; if not, the game offers to install it.

Linux
  1. Install WebKitGTK 4.1, e.g. on Ubuntu or Debian:
       sudo apt install libwebkit2gtk-4.1-0
     Fedora: sudo dnf install webkit2gtk4.1
     Arch: sudo pacman -S webkit2gtk-4.1
  2. Unpack: tar -xzf four-souls-*.tar.gz
  3. Run: ./FourSouls
     If it does not start: chmod +x FourSouls

The firewall may ask whether the game may use the network the first
time you host: allow it, or friends cannot connect.


Play
----

1. Start the game on every computer.
2. One player chooses "Host a game". The table shows the addresses
   the others can use.
3. The others choose "Join a game", type one of those addresses,
   press Connect, then Join.
4. Everyone presses Ready.

Over the internet, use a virtual LAN (Tailscale, ZeroTier, Radmin VPN)
or a dedicated server: see the wiki pages "Playing over the internet"
and "Dedicated server on a VPS".

Several copies of the game can run on one computer; each one is its
own player with its own nickname. On macOS start another copy from
Terminal: open -n /Applications/FourSouls.app


Dedicated server
----------------

The dedicated server has its own archives in the same release:
four-souls-server-<system>-<cpu>. Players join the server's address.
See the wiki page "Dedicated server on a VPS".


Licenses
--------

The code is MIT (LICENSE). Third-party items are listed in NOTICE.
