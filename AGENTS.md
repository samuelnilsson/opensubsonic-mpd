# AGENTS.md

Instructions and conventions for working with this repository.

## Building

### Nix (recommended)

```bash
# Build the binary
nix build .#default

# Run directly
nix run .#default -- --help

# Development shell (go + libmpv headers)
nix develop
```

The binary is placed at `result/bin/opensubsonic-mpd` (a symlink into `/nix/store`).
To copy it out: `cp -L result/bin/opensubsonic-mpd ~/bin/`.

### Traditional Go

```bash
go build ./cmd/opensubsonic-mpd/
```

Requires Go 1.22+ and libmpv development headers.

## Nix dependency management

When adding Go dependencies (i.e. when `go.mod` acquires entries and `go.sum` is generated):

1. Set `vendorHash = null;` in `flake.nix` (if not already)
2. Run `nix build .#default` locally
3. Copy the `expected vendorHash = sha256-...` line from the error output
4. Update `vendorHash` in `flake.nix` with that value
5. Commit both `flake.nix` and the updated `flake.lock`

## Formatting

The project uses [`treefmt-nix`](https://github.com/nix-community/treefmt) for unified formatting across both Go and Nix files.

```bash
# Format all files (Go + Nix)
nix fmt

# Check formatting (runs automatically via `nix flake check` in CI)
nix flake check
```

This runs:
- **gofmt** — Go code formatting
- **nixfmt** — Nix code formatting

## GitHub Actions

A CI workflow builds the project with Nix on push and PR. No secrets or cache setup required for basic operation.

### Automatic releases

Releases are created automatically on every push to `main` using [`semantic-release`](https://semantic-release.gitbook.io/semantic-release). It:
1. Analyzes conventional commits to determine the next semver version
2. Generates `CHANGELOG.md` entries
3. Commits the version bump + changelog back to the branch
4. Creates a GitHub Release with the changelog as release notes
5. Attaches the Nix-built binary to the release

Version bump rules:
| Commit type | Bump       |
|-------------|------------|
| `feat`      | minor      |
| `fix`       | patch      |
| `perf`      | patch      |
| `chore`, `docs`, `style`, `refactor`, `test`, `build`, `ci` | none |

The config lives in `.releaserc.json`.

## Branching convention

**Never commit directly to `main` — always use a PR.**

1. Create a feature branch: `git checkout -b feat/<description>`
2. Make changes and commit with conventional commit messages
3. Push the branch and open a PR against `main`
4. CI will run `nix flake check` (build + formatting) on the PR

## Conventional Commits

Use [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages. Format:

```
<type>(<scope>): <description>
```

Types:
| Type         | When to use                                                                 |
|--------------|-----------------------------------------------------------------------------|
| `feat`       | A new feature                                                               |
| `fix`        | A bug fix                                                                   |
| `docs`       | Documentation changes only                                                  |
| `style`      | Code style changes (formatting, semicolons, etc. — no logic change)         |
| `refactor`   | Code changes that neither fix a bug nor add a feature                       |
| `perf`       | Performance improvements                                                    |
| `test`       | Adding or correcting tests                                                  |
| `build`      | Build system or dependency changes (e.g. Nix, go.mod)                       |
| `ci`         | CI pipeline changes                                                         |
| `chore`      | Other changes that don't modify src or test files                           |

Examples:
- `feat(nix): add buildGoModule package to flake.nix`
- `fix: handle missing server URL flag`
- `ci: add GitHub Actions workflow for Nix build`

### PR titles

Since this repo uses squash merging, **PR titles must also follow the Conventional Commits format** (enforced by a CI check). The squash-merged commit — which becomes the official commit on `main` — takes the PR title, so it must be valid.

## Project overview

`opensubsonic-mpd` is a headless music player that bridges an [OpenSubsonic](https://github.com/Mromani/opensubsonic) server with MPD clients via libmpv. The main entry point is `cmd/opensubsonic-mpd/main.go`.
