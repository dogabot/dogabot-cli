#!/usr/bin/env bash
# List markets (requires DOGABOT_API_KEY).
set -euo pipefail
dogabot markets get-markets --json
