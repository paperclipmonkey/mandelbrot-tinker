#!/bin/bash

# Create a pre-commit hook that links to the pre-commit.sh script
HOOKS_DIR=".git/hooks"
PRE_COMMIT_HOOK="$HOOKS_DIR/pre-commit"
PRE_COMMIT_SCRIPT="./pre-commit.sh"

# Ensure the hooks directory exists
mkdir -p $HOOKS_DIR

# Create a symbolic link to the pre-commit.sh script
ln -sf $PRE_COMMIT_SCRIPT $PRE_COMMIT_HOOK

echo "Pre-commit hook installed successfully."
