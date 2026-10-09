#!/usr/bin/env bash
# Print the current API-key account (requires DOGABOT_API_KEY).
set -euo pipefail
dogabot account get-me --json
