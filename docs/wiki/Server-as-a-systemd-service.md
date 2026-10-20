# Server as a systemd service

The smallest setup: the server binary, run by systemd, so it starts on boot and restarts after a crash. Read [Dedicated server on a VPS](Dedicated-server-on-a-VPS) first: the port, options and safety are there.

## 1. Install

Log in and create a user that runs the server, so it never runs as root:

```sh
ssh root@<server address>
adduser --system --group --home /var/lib/four-souls four-souls
```

Download the latest server (use `linux-arm64` on an ARM VPS) and install it:

```sh
cd /tmp
curl -fsSL https://github.com/FourSoulsRepo/four_souls/releases/latest/download/four-souls-server-linux-amd64.tar.gz | tar -xz
install -m 755 four-souls-server-linux-amd64/four_souls_server /usr/local/bin/
```

Or build it yourself on any machine with Go (version in `go.mod`), for the server's CPU, and copy it with `scp`:

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags embed -o four_souls_server ./cmd/server   # or GOARCH=arm64
```

## 2. Try it

```sh
cd /var/lib/four-souls
sudo -u four-souls four_souls_server
```

It prints the fan-game notice, its versions and the address it listens on, then waits for players. Stop it with Ctrl+C; it closes every connection cleanly and closes the records of running games as unfinished.

## 3. Run it as a service

Download the unit file from the repository (`deploy/systemd/four-souls.service`):

```sh
curl -fsSL -o /etc/systemd/system/four-souls.service \
  https://raw.githubusercontent.com/FourSoulsRepo/four_souls/main/deploy/systemd/four-souls.service
systemctl daemon-reload
systemctl enable --now four-souls
systemctl status four-souls        # is it running?
journalctl -u four-souls -f        # its output
```

Records go to `/var/lib/four-souls/records`, saved games to `/var/lib/four-souls/saves`.

## Settings

Put `FOUR_SOULS_*` lines (see the Options table) into `/etc/four-souls.env`; the unit reads it if it exists:

```sh
cat > /etc/four-souls.env <<'EOF'
FOUR_SOULS_NAME=Basement
FOUR_SOULS_RETENTION=60
EOF
systemctl restart four-souls
```

## Updating

Download the newest server over the old one and restart:

```sh
cd /tmp
curl -fsSL https://github.com/FourSoulsRepo/four_souls/releases/latest/download/four-souls-server-linux-amd64.tar.gz | tar -xz
install -m 755 four-souls-server-linux-amd64/four_souls_server /usr/local/bin/
systemctl restart four-souls
journalctl -u four-souls -n 5 | grep app   # the new version shows at start
```

Restarting ends the running games; their records stay, marked unfinished.
