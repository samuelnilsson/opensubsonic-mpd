# opensubsonic-mpd

> ⚠️ **This project is a work in progress and not yet ready for production use.**

[![Status: Work in Progress](https://img.shields.io/badge/Status-Work_in_Progress-orange)](#)

A headless music player that bridges any [OpenSubsonic](https://github.com/Mromani/opensubsonic) server with MPD clients via [libmpv](https://mpv.io/).

## What it does

`opensubsonic-mpd` acts as a drop-in MPD-compatible server. Any MPD client — can connect to it and control playback of music stored on your OpenSubsonic server. Playback is handled by MPV through libmpv, so you get high-quality audio without needing a graphical interface.

In short: your OpenSubsonic library, played through any MPD client you like.

## Why

OpenSubsonic provides great remote music hosting, but its built-in web player is limited and many users prefer dedicated MPD clients. Rather than writing yet another UI, this project lets you use the clients you already love. It's lightweight, terminal-friendly, and fully scriptable via the standard MPD protocol.

## OpenSubsonic specification

This project follows the [OpenSubsonic API specification](https://github.com/opensubsonic/opensubsonic.github.io) for all server communication.

## Getting Started

### Prerequisites

- Go 1.22+
- libmpv (development headers)

### Building

```bash
go build ./...
```

### Running

```bash
opensubsonic-mpd --server <opensubsonic-url> --username <user> --password <pass>
```

The server starts listening on port `6600` by default. Connect your favourite MPD client and start playing.

## Nix

### Building

```bash
nix build .#default
nix run .#default -- --help
```

### Formatting

```bash
# Format all files (Go + Nix)
nix fmt
```

## Contributing

PR titles must follow [Conventional Commits](https://www.conventionalcommits.org/) format. This is enforced by CI since the repo uses squash merging — the PR title becomes the commit on `main`.

## License

[MIT](LICENSE)
