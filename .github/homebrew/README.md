# Homebrew Formula for tableconvert

This directory contains the Homebrew formula for tableconvert.

## Installation

### For users (once published to Homebrew core):

```bash
# Install tableconvert
brew install tableconvert

# Upgrade to latest version
brew upgrade tableconvert
```

### For testing the formula (before publishing):

```bash
# Install directly from this formula
brew install --HEAD .github/homebrew/tableconvert.rb

# Or tap the repository (once available)
brew tap martianzhang/tableconvert
brew install tableconvert
```

## Development

### Testing the formula locally:

```bash
# Create a test version
brew install --build-from-source --HEAD .github/homebrew/tableconvert.rb

# Run tests
brew test tableconvert

# Check for issues
brew audit --strict tableconvert.rb
brew style tableconvert.rb
```

### Publishing to Homebrew core:

1. **Fork homebrew-core**:
   ```bash
   git clone https://github.com/Homebrew/homebrew-core.git
   cd homebrew-core
   ```

2. **Create new formula**:
   ```bash
   # Create formula file
   brew create --go --set-version=1.0.0-pre \
     https://github.com/martianzhang/tableconvert/archive/refs/tags/v1.0.0-pre.tar.gz
   ```

3. **Copy our formula**:
   ```bash
   cp /path/to/tableconvert/.github/homebrew/tableconvert.rb \
      Formula/t/tableconvert.rb
   ```

4. **Calculate SHA256**:
   ```bash
   curl -L https://github.com/martianzhang/tableconvert/archive/refs/tags/v1.0.0-pre.tar.gz | shasum -a 256
   ```

5. **Update formula**:
   - Update `url` and `sha256`
   - Update `version` if needed
   - Ensure `desc` and `homepage` are correct

6. **Test locally**:
   ```bash
   brew audit --strict Formula/t/tableconvert.rb
   brew style Formula/t/tableconvert.rb
   brew install --build-from-source tableconvert
   brew test tableconvert
   ```

7. **Submit PR**:
   ```bash
   git add Formula/t/tableconvert.rb
   git commit -m "tableconvert 1.0.0-pre (new formula)"
   git push origin tableconvert-1.0.0-pre
   ```

## Formula Details

### Build Process
- Uses Go's standard build process
- Compiles from `./cmd/tableconvert`
- Strips symbols for smaller binary (`-s -w`)

### Dependencies
- Only requires Go 1.23+ to build
- No runtime dependencies

### Test Coverage
- Basic conversion test (CSV → JSON)
- Help command verification
- Version check

## References

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [tableconvert GitHub](https://github.com/martianzhang/tableconvert)
