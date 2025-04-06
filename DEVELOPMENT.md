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
