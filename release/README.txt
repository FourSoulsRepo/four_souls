Four Souls
==========

An unofficial, free, fan-made desktop version of the card game
"The Binding of Isaac: Four Souls" for 2-4 players over the network.

The game belongs to Edmund McMillen and Maestro Media.
Please buy the physical game: https://foursouls.com


What is in this archive
-----------------------

  four_souls (.app on macOS, .exe on Windows)   the game
  four_souls_server                              a dedicated server, optional
  LICENSE, NOTICE                                licenses


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


Dedicated server
----------------

  four_souls_server -port 4774

Players join the server's address. Records go to ./records, saved
games to ./saves. Run "four_souls_server -h" for all options.


Licenses
--------

The code is MIT (LICENSE). Third-party items are listed in NOTICE.
