# Four Souls

Unofficial, free, fan-made desktop version of the card game
*The Binding of Isaac: Four Souls* for 2–4 players over the network.
Built with [Wails](https://wails.io) (Go + React/TypeScript).

The game belongs to Edmund McMillen and Maestro Media.
Please [buy the physical game](https://foursouls.com).

## Play

Download the archive for your system from
[Releases](https://github.com/FourSoulsRepo/four_souls/releases).
The app is not code-signed: `README.txt` inside, or the wiki page
[First start](docs/wiki/First-start.md), shows how to open it the first time.

Several copies can run on one computer; each is its own player.
On macOS start another copy with `open -n four_souls.app`.

## Run in dev mode

```sh
wails dev
```

## Build

```sh
wails build -tags "embed cardimages"
```

Card images come from a private submodule (`pkg/card_db/images`). Without access, build with `-tags embed`; cards then show as text.

The binary lands in `build/bin/`.

## Develop

```sh
make build vet test lint wasm
cd frontend && npm run typecheck && npm run lint
```

Developer notes: [docs/wiki/Home.md](docs/wiki/Home.md).
Plan: [docs/roadmap.md](docs/roadmap.md).

## License

Code: MIT ([LICENSE](LICENSE)). Third-party items: [NOTICE](NOTICE).
