# node-run

An fzf-like npm script runner with monorepo support. Written in go.

Supports all modern package managers:

- npm ✅
- pnpm ✅
- yarn ✅

## Installation

```
brew tap silesky/tap
brew install node-run
```

## Usage

```
$> cd ~/projects/my-node-package
$> nr
```

## Development

### Build & Execute

```sh
make run ARGS=--cwd=~/projects/monorepo-example
```

### Build Only

```sh
make build # See artifacts in `./bin`
```

### Releasing

```sh
make version # Bumps the tag, which
```
