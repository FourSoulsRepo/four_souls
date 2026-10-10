# Playing over the internet

The game is built for a LAN: one player hosts, the others join by address (ADR 006). Over the internet, the host's home router blocks incoming connections. Until the game has its own relay or online server, there are three ways around that:

1. **A virtual LAN** — every player installs the same free tool; the game then sees a LAN. Easiest; nothing to configure on routers.
2. **Port forwarding** — the host opens the game's port on their router.
3. **A dedicated server** — one player rents a small server and runs the game server there for everyone. See [Dedicated server on a VPS](Dedicated-server-on-a-VPS).

The game uses one TCP port, **4774** by default.

## Option 1: a virtual LAN

A virtual LAN gives every player's computer an extra private address (for example `100.101.102.103`). The host gives that address to the friends; they join it like a LAN address. Traffic between the players is encrypted by the tool.

All of these have free plans for a small group of friends. Their limits change; check the current ones on their sites.

| Tool | Systems | Good to know |
|---|---|---|
| [Tailscale](https://tailscale.com) | Windows, macOS, Linux, iOS, Android | Simplest setup; each player logs in with an existing account (Google, Microsoft, GitHub, …). |
| [ZeroTier](https://www.zerotier.com) | Windows, macOS, Linux, iOS, Android | One network ID for the group; the host approves who joins. |
| [Radmin VPN](https://www.radmin-vpn.com) | Windows only | No account; popular for games. |
| [Hamachi](https://vpn.net) | Windows, macOS | The classic one; the free plan allows only a few members per network. |

### Tailscale

1. Every player installs Tailscale and signs in.
2. The host invites the friends into their network: in the [admin console](https://login.tailscale.com/admin/users), **Invite users**, or share just the host's computer with them (**Share** on the machine's page).
3. Each friend accepts the invite; their Tailscale app now lists the host's computer.
4. The host looks up their own Tailscale address: in the Tailscale menu, or `tailscale ip -4` in a terminal. It starts with `100.`.
5. The host starts hosting in the game. Friends join `100.x.y.z` (port 4774 is the default).

### ZeroTier

1. The host creates a free account at [my.zerotier.com](https://my.zerotier.com) and clicks **Create A Network**. The network ID is 16 characters, like `8056c2e21c000001`.
2. Every player (the host too) installs ZeroTier and joins that network ID.
3. In the network's page on my.zerotier.com, the host ticks **Auth** for each member.
4. Each member gets a managed IP, shown on that page. The host gives their managed IP to the friends.
5. The host starts hosting; friends join that address.

### Radmin VPN

1. Every player installs Radmin VPN (Windows).
2. The host: **Network → Create network**, with a name and a password.
3. Friends: **Network → Join network**, with the same name and password.
4. The host's Radmin address (it starts with `26.`) is shown next to their name. Friends join it in the game.

### Hamachi

1. Every player installs Hamachi and creates a LogMeIn account.
2. The host: **Network → Create a new network**, with an ID and a password.
3. Friends: **Network → Join an existing network**.
4. The host's Hamachi address (it starts with `25.`) is shown at the top of the window. Friends join it in the game.

### If friends cannot connect

* **Firewall.** The first time the game hosts, Windows asks whether to allow it on networks: allow it. On macOS, allow incoming connections when asked. A virtual LAN adapter may count as a "public" network on Windows; allow both private and public, or set that adapter to private.
* **Same tool, same network.** Everyone must be in the same Tailscale network, ZeroTier network or Radmin/Hamachi network, and online in it.
* **The right address.** Use the host's virtual address, not their home address (`192.168.…`).
* **Check the path.** A friend can run `ping <host address>`; if the ping works but the game does not, it is the firewall.

## Option 2: port forwarding

The host tells their router to send the game's port to their computer.

1. Give the host's computer a fixed LAN address (in the router's DHCP settings, "address reservation").
2. In the router's settings, find **Port forwarding** (also **Virtual server** or **NAT**). Forward **TCP port 4774** to that LAN address.
3. Allow the game in the host's firewall.
4. Friends join the host's public address; the host can see it at sites such as [ifconfig.me](https://ifconfig.me).

It does not work when the internet provider puts the host behind its own NAT ("CGNAT"; the router's WAN address then differs from the public one, often `100.64.…`). Use a virtual LAN or a server instead.

## Privacy

The game has no accounts and does not encrypt its own traffic yet (ADR 006). With a virtual LAN the tool encrypts everything between players. With port forwarding or a server, anyone who knows the address can join a game that is open. Share the address only with your friends.
