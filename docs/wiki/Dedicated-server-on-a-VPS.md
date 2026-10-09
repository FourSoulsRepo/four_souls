# Dedicated server on a VPS

Any player can rent a small server on the internet (a VPS) and run the game's dedicated server there. Friends then join its address from anywhere; nobody needs a virtual LAN or port forwarding (see [Playing over the internet](Playing-over-the-internet)).

> Work in progress: until the lobby is ready (roadmap step 6.5) the server runs one game with a fixed number of seats (`-players`). Saving match records comes with step 6.9.

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

It prints the fan-game notice, its versions and the address it listens on, then waits for players. Stop it with Ctrl+C (or `systemctl stop`); it closes every connection cleanly.

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

A JSON config file can set the same values; flags override it.

| Flag | Default | Meaning |
|---|---|---|
| `-port` | 4774 | TCP port for players |
| `-addr` | all addresses | listen only on this address |
| `-records` | `records` | where match records are saved (A-12) |
| `-retention` | 30 | days to keep records; 0 keeps them forever (RP-07) |
| `-players` | 2 | seats of the game, until the lobby exists |
| `-config` | none | a JSON file with the same settings |

```json
{"port": 4774, "records": "/var/lib/four-souls/records", "retention": 30, "players": 4}
```

## Updating

When the game gets a new release, players update their apps and you replace the binary:

```sh
systemctl stop four-souls
scp four-souls-server root@<server address>:/usr/local/bin/
systemctl start four-souls
```

Stopping ends the running games for now; saving them comes with step 6.9.

## Safety

* The game has no accounts and does not encrypt its traffic yet (ADR 006). Anyone who knows the address can connect; share it only with your friends.
* Keep the VPS updated (`apt update && apt upgrade` on Debian or Ubuntu), log in with SSH keys, and do not run the server as root.
* For encryption today, combine both pages: put the VPS into the same Tailscale or ZeroTier network as the players, and do not open port 4774 to the internet at all.
