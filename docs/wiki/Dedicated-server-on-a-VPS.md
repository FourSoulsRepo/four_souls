# Dedicated server on a VPS

Any player can rent a small server on the internet (a VPS) and run the game's dedicated server there. Friends then join its address from anywhere; nobody needs a virtual LAN or port forwarding (see [Playing over the internet](Playing-over-the-internet)).

> The dedicated server is built in roadmap step 6.3. Commands and flags on this page are the planned ones; they are checked against the real binary when it ships.

## What you need

* A Linux VPS. The server is small: 1 CPU and 512 MB–1 GB of memory are plenty for a few games. Any provider works (Hetzner, DigitalOcean, Vultr, OVH, Linode, a local one); the cheapest plan is usually enough. Pick a location close to your friends.
* SSH access to it (the provider shows the address and the login).
* The `four-souls-server` binary for Linux from the game's releases, matching your CPU (`x64` or `arm64`). Everyone's app must speak the same protocol version as the server; using the same release is safest.

## 1. Install

Log in and create a user that runs the server, so it never runs as root:

```sh
ssh root@<server address>
adduser --system --group --home /var/lib/four-souls four-souls
```

Copy the binary and make it runnable (from your computer):

```sh
scp four-souls-server root@<server address>:/usr/local/bin/
ssh root@<server address> chmod 755 /usr/local/bin/four-souls-server
```

Or build it yourself on any machine with Go (version in `go.mod`), for the server's CPU:

```sh
GOOS=linux GOARCH=amd64 go build -o four-souls-server ./cmd/server   # or GOARCH=arm64
```

## 2. Try it

```sh
sudo -u four-souls four-souls-server -port 4774 -records /var/lib/four-souls/records
```

It prints the fan-game notice and its versions, then waits for players. Stop it with Ctrl+C; running games are saved first.

## 3. Open the port

The game uses **TCP port 4774** (or the one you chose). With `ufw`:

```sh
ufw allow OpenSSH
ufw allow 4774/tcp
ufw enable
```

Many providers also have a firewall in their web panel; allow the port there too.

## 4. Run it as a service

So it starts on boot and restarts after a crash, add `/etc/systemd/system/four-souls.service`:

```ini
[Unit]
Description=Four Souls game server
After=network-online.target
Wants=network-online.target

[Service]
User=four-souls
Group=four-souls
ExecStart=/usr/local/bin/four-souls-server -port 4774 -records /var/lib/four-souls/records
Restart=on-failure
NoNewPrivileges=true
ProtectSystem=strict
ReadWritePaths=/var/lib/four-souls

[Install]
WantedBy=multi-user.target
```

Then:

```sh
systemctl daemon-reload
systemctl enable --now four-souls
systemctl status four-souls        # is it running?
journalctl -u four-souls -f        # its output
```

## 5. Play

Friends open the game, choose **Join**, and enter the server's address (and the port, if it is not 4774). One of them creates the game in the lobby; the others join it.

## Options

Planned flags; a JSON config file can set the same values.

| Flag | Default | Meaning |
|---|---|---|
| `-port` | 4774 | TCP port for players |
| `-addr` | all addresses | listen only on this address |
| `-records` | `records/` | where match records are saved (A-12) |
| `-retention` | 30 days | how long records are kept; 0 keeps them forever (RP-07) |
| `-config` | none | a JSON file with the same settings |

## Updating

When the game gets a new release, players update their apps and you replace the binary:

```sh
systemctl stop four-souls
scp four-souls-server root@<server address>:/usr/local/bin/
systemctl start four-souls
```

Stopping saves running games; finished ones stay in the records folder.

## Safety

* The game has no accounts and does not encrypt its traffic yet (ADR 006). Anyone who knows the address can connect; share it only with your friends.
* Keep the VPS updated (`apt update && apt upgrade` on Debian or Ubuntu), log in with SSH keys, and do not run the server as root.
* For encryption today, combine both pages: put the VPS into the same Tailscale or ZeroTier network as the players, and do not open port 4774 to the internet at all.
