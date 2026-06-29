#!/bin/bash

set -e

SKILL_NAME="generate-flow"
SKILL_FILE="$(dirname "$0")/SKILL.md"
TARGET_DIR=".claude/commands"
TARGET_FILE="$TARGET_DIR/$SKILL_NAME.md"

if [ ! -f "$SKILL_FILE" ]; then
  echo "error: SKILL.md not found at $SKILL_FILE"
  exit 1
fi

mkdir -p "$TARGET_DIR"
cp "$SKILL_FILE" "$TARGET_FILE"

echo "✓ installed /$SKILL_NAME → $TARGET_FILE"
