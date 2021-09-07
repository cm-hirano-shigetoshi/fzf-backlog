#!/usr/bin/env bash
set -eu

TOOL_DIR=$(dirname $(env - python3 -c "import os; print(os.path.realpath('$0'));"))
fzfyml4 run $TOOL_DIR/yml/wiki.yml $TOOL_DIR/backlog "$@"
