#!/bin/bash
# Usage: "bump-version.sh major|minor|patch"
# Exit on error
set -e

# Check if the current branch is main
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
  echo "Error: You are not on the main branch. Please switch to the main branch and try again."
  exit 1
fi

# Pull, to ensure caught up
git pull --ff-only

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
echo "--------"
echo "Create and push tag: $NEW_TAG"
read -p "Are you sure you want to proceed? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo "Aborted."
  exit 1
fi

# Commit changes with a message
COMMIT_MESSAGE=${2:-"Release $NEW_TAG"} # Default message if none provided
git add . || echo "nothing to add"
git commit -m "$COMMIT_MESSAGE" --allow-empty

# Create an annotated tag with a messageg
git tag -a "$NEW_TAG" -m "$COMMIT_MESSAGE"

# Push the commit and the new tag
git push --follow-tags

echo "✅ New commit and annotated tag created and pushed: $NEW_TAG"
echo "📜 Message: $COMMIT_MESSAGE"
