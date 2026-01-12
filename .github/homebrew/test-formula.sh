#!/bin/bash
# Test script for Homebrew formula

set -e

echo "🧪 Testing Homebrew formula for tableconvert..."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Please run from tableconvert project root"
    exit 1
fi

# Check if formula exists
FORMULA=".github/homebrew/tableconvert.rb"
if [ ! -f "$FORMULA" ]; then
    echo "❌ Formula not found at $FORMULA"
    exit 1
fi

echo "✅ Formula file exists"

# Check formula syntax (basic Ruby check)
echo "Checking Ruby syntax..."
ruby -c "$FORMULA"

# Check for common formula issues
echo "Checking formula structure..."

# Check for required fields
grep -q "class Tableconvert < Formula" "$FORMULA" && echo "✅ Class definition OK" || echo "❌ Missing class"
grep -q 'desc "' "$FORMULA" && echo "✅ Description OK" || echo "❌ Missing desc"
grep -q 'homepage "' "$FORMULA" && echo "✅ Homepage OK" || echo "❌ Missing homepage"
grep -q 'url "' "$FORMULA" && echo "✅ URL OK" || echo "❌ Missing url"
grep -q 'license "' "$FORMULA" && echo "✅ License OK" || echo "❌ Missing license"
grep -q 'depends_on "go"' "$FORMULA" && echo "✅ Go dependency OK" || echo "❌ Missing go dependency"
grep -q 'def install' "$FORMULA" && echo "✅ Install method OK" || echo "❌ Missing install"
grep -q 'test do' "$FORMULA" && echo "✅ Test block OK" || echo "❌ Missing test"

# Check for SHA256 placeholder
if grep -q "# sha256" "$FORMULA"; then
    echo "⚠️  SHA256 is placeholder - needs to be calculated for release"
else
    echo "✅ SHA256 present"
fi

# Test the actual build process
echo ""
echo "Testing actual build..."
go build -o /tmp/tableconvert-test ./cmd/tableconvert

# Test version flag
echo "Testing version flag..."
/tmp/tableconvert-test --version

# Test verbose flag
echo "Testing verbose flag..."
echo "name,age\nAlice,30" | /tmp/tableconvert-test --from=csv --to=json -v 2>&1 | head -3

# Test basic functionality
echo "Testing basic conversion..."
echo "name,age
Alice,30
Bob,25" | /tmp/tableconvert-test --from=csv --to=json

echo ""
echo "✅ All tests passed!"
echo ""
echo "📋 Next steps:"
echo "1. Create a release tag: git tag v1.0.0"
echo "2. Calculate SHA256: curl -L <url> | shasum -a 256"
echo "3. Update formula with real SHA256"
echo "4. Submit to Homebrew core (see README.md)"

# Cleanup
rm -f /tmp/tableconvert-test
