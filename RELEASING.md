# Build the binary for different architectures

Step 1: Build the Binaries

```
make build

```

2. Create a release on GitHub:

Tag release:

```
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### Go to your repository on GitHub.:

- Click on "Releases" on the right sidebar.
- Click "Draft a new release".
- Choose the tag you just created (e.g., v1.0.0).
- Fill in the release title and description.
- Attach the built binaries (nr for amd64 and arm64).
- Click "Publish release".
