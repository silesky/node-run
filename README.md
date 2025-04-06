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

### Commands

```
$> nrun --help
Usage: nrun [options] [initial search input]
Example:
- nrun
- nrun my-search-string
Options:
  --cwd
        Specify current working directory
  --debug
        Turn debug logging on
  --help
        Show help text
  --version
        Show version
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
