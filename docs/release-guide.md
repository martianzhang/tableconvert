# Release Guide

This guide explains how to create releases for tableconvert.

## Quick Start

### Local Release Build

```bash
# Build all platform binaries and create checksums
make release

# Output will be in ./release/ directory:
# - tableconvert-linux-amd64
# - tableconvert-linux-arm64
# - tableconvert-darwin-amd64
# - tableconvert-darwin-arm64
# - tableconvert-windows-amd64.exe
# - checksums.txt
# - RELEASE_INFO.txt
```

### Optional: Create Zip Archives

```bash
# Requires zip command to be installed
make release-zip

# Creates individual zip files for each platform
```

### Optional: Generate Release Notes

```bash
# Generate release notes from git history
make release-notes

# Creates RELEASE_NOTES.md with commit history
```

## Automated GitHub Releases

### Workflow

When you push a tag starting with `v`, the CI/CD pipeline automatically:

1. **Runs tests** - Ensures all tests pass
2. **Runs linting** - Checks code quality
3. **Builds binaries** - Creates binaries for all platforms
4. **Generates checksums** - SHA256 checksums for verification
5. **Creates release** - GitHub Release with changelog and files

### Creating a Release

```bash
# 1. Update version (if needed)
# Edit any version references in code/docs

# 2. Create and push tag
git tag v1.0.0
git push origin v1.0.0

# 3. Wait for CI/CD
# The release workflow will automatically:
# - Build all binaries
# - Create checksums
# - Generate changelog
# - Create GitHub Release
```

### What Gets Released

The GitHub Release includes:

- **Binaries** for all platforms:
  - `tableconvert-linux-amd64`
  - `tableconvert-linux-arm64`
  - `tableconvert-darwin-amd64`
  - `tableconvert-darwin-arm64`
  - `tableconvert-windows-amd64.exe`

- **Security files**:
  - `checksums.txt` - SHA256 checksums for verification

- **Documentation**:
  - `RELEASE_INFO.txt` - Build metadata
  - `CHANGELOG.md` - Changes in this release

## Manual Release Process

If you need to release manually:

```bash
# 1. Build locally
make release

# 2. Verify binaries
./release/tableconvert-linux-amd64 --help

# 3. Verify checksums
cd release && sha256sum -c checksums.txt

# 4. Create GitHub Release manually
# Go to: https://github.com/martianzhang/tableconvert/releases/new
# Upload all files from ./release/
```

## Release Checklist

Before creating a release:

- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated (if manual release)
- [ ] Version numbers are updated
- [ ] Breaking changes are documented

After release:

- [ ] Release notes are readable
- [ ] Binaries work on target platforms
- [ ] Checksums verify correctly
- [ ] Installation instructions are clear

## Platform Support

| Platform | Architecture | Binary Name |
|----------|--------------|-------------|
| Linux | x64 | `tableconvert-linux-amd64` |
| Linux | ARM64 | `tableconvert-linux-arm64` |
| macOS | x64 | `tableconvert-darwin-amd64` |
| macOS | ARM64 | `tableconvert-darwin-arm64` |
| Windows | x64 | `tableconvert-windows-amd64.exe` |

## Installation from Release

### Linux/macOS

```bash
# Download binary
wget https://github.com/martianzhang/tableconvert/releases/download/v1.0.0/tableconvert-linux-amd64

# Make executable
chmod +x tableconvert-linux-amd64

# Move to PATH
sudo mv tableconvert-linux-amd64 /usr/local/bin/tableconvert

# Verify
tableconvert --help
```

### Windows

```powershell
# Download from GitHub Releases
# Move to desired location
# Add to PATH or run directly
.\tableconvert-windows-amd64.exe --help
```

## Verification

Always verify downloaded binaries using checksums:

```bash
# Download both binary and checksums.txt
# Then verify:
sha256sum -c checksums.txt

# Should output:
# tableconvert-linux-amd64: OK
```

## Troubleshooting

### Build fails on Windows
- Ensure you have Go installed
- Use PowerShell or WSL

### Zip command not found
- Install zip: `sudo apt-get install zip` (Ubuntu/Debian)
- Or skip zip step: `make release` only

### Checksums don't match
- Rebuild: `make release-clean && make release`
- Ensure no file corruption during download

### Release workflow doesn't trigger
- Check tag format: must start with `v` (e.g., `v1.0.0`)
- Verify tag was pushed: `git push origin --tags`

## CI/CD Configuration

The release workflow is defined in `.github/workflows/ci.yml`:

```yaml
release:
  name: Create Release
  runs-on: ubuntu-latest
  if: startsWith(github.ref, 'refs/tags/v')
  needs: [test, lint]
  # ... build and release steps
```

Key features:
- **Triggers on tags**: `v*` pattern
- **Requires tests**: Won't release if tests fail
- **Auto-generates**: Changelog, checksums, release info
- **Multi-platform**: Linux, macOS, Windows (x64 + ARM64)

## Make Targets Reference

| Target | Description |
|--------|-------------|
| `make release` | Build all binaries + checksums + info |
| `make release-clean` | Clean release directory |
| `make release-checksums` | Generate checksums only |
| `make release-zip` | Create zip archives (requires zip) |
| `make release-notes` | Generate release notes from git |

## Next Steps

After setting up releases:

1. **Test the workflow**: Create a pre-release tag (e.g., `v1.0.0-beta`)
2. **Monitor CI**: Check GitHub Actions tab for workflow progress
3. **Verify release**: Download and test binaries from GitHub
4. **Share**: Announce release with installation instructions

## Related Documentation

- [README.md](../README.md) - User documentation
- [CLAUDE.md](../CLAUDE.md) - Developer guide
- [arguments.md](arguments.md) - Format parameters
