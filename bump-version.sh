#!/bin/bash
# Usage: "bump-version.sh major|minor|patch"
# Exit on error
set -e

# Get the latest tag (default to v0.0.0 if none exists)
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Extract major, minor, and patch versions
MAJOR=$(echo "$LATEST_TAG" | cut -d. -f1 | tr -d 'v')
MINOR=$(echo "$LATEST_TAG" | cut -d. -f2)
PATCH=$(echo "$LATEST_TAG" | cut -d. -f3)

# Choose the version bump type (patch, minor, major)
if [ "$1" == "major" ]; then
  MAJOR=$((MAJOR + 1))
  MINOR=0
  PATCH=0
elif [ "$1" == "minor" ]; then
  MINOR=$((MINOR + 1))
  PATCH=0
else
  PATCH=$((PATCH + 1))
fi

# Generate new version tag
NEW_TAG="v$MAJOR.$MINOR.$PATCH"

# Ask for confirmation
echo "⚠️  Are you sure you want to create and push tag: $NEW_TAG? (y/N)"
read -r CONFIRM

# Convert input to lowercase and check if it's "y"
if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
  echo "🚫 Tag creation canceled."
  exit 1
fi

# Create an annotated tag with a message
ANNOTATION=${2:-"Release $NEW_TAG"} # Default message if none provided
git tag -a "$NEW_TAG" -m "$ANNOTATION"

# Push the new tag
git push origin "$NEW_TAG"

echo "✅ New annotated tag created and pushed: $NEW_TAG"
echo "📜 Message: $ANNOTATION"
