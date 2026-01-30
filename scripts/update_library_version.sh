#!/bin/bash
HELP_TEXT="update_library_version.sh -h <help> -d <dry-run>"
DRY_RUN=false

while getopts "dh" flag; do
case ${flag} in
d) DRY_RUN=true
   ;;
h) echo $HELP_TEXT; exit 0;
   ;;
esac
done


REPOSITORY_PATH="."
NEXT_VERSION=$(./scripts/calculate_next_version.sh -r $REPOSITORY_PATH)
NEXT_VERSION="v${NEXT_VERSION}"
BRANCH_NAME=$(git branch --show-current)

echo "Updating library version to $NEXT_VERSION..." >&2

if [ "$DRY_RUN" = true ]; then
    echo "[DRY RUN] Would execute: git commit --allow-empty -m \"ci: Release version $NEXT_VERSION\""
    echo "[DRY RUN] Would execute: git tag -a \"$NEXT_VERSION\" -m \"Release $NEXT_VERSION\""
    echo "[DRY RUN] Would execute: git push -u origin \"$BRANCH_NAME\""
    echo "[DRY RUN] Would execute: git push -u origin --tags"
else
    git commit --allow-empty -m "ci: Release version $NEXT_VERSION"
    git tag -a "$NEXT_VERSION" -m "Release $NEXT_VERSION"
    git push -u origin "$BRANCH_NAME"
    git push -u origin --tags
fi

echo "Successfully updated library version to $VERSION and pushed to upstream" >&2