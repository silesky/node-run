# node-run

An fzf-like npm script runner with monorepo support.

Supports modern package managers:

- npm ✅
- pnpm ✅
- yarn ✅
- Bun ✅

## Installation

```sh
$> brew tap silesky/tap
$> brew install node-run
```

## Usage

```sh
$> cd ~/my-projects/my-node-package
$> nrun
```

## Development

### Test

```sh
$> make test
```

### Run

```sh
$> make run ARGS=--cwd=~/projects/monorepo-example
```

### Release

```sh
$> make version # Bump version via git tag + push (trigger CI job in .github/workflows)
```
