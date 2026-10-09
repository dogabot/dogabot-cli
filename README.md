# dogabot CLI

Official command-line client for the [dogabot public API](https://docs.dogabot.com/).

Covers the **full REST allowlist** (same `operationId`s as the official SDKs). MCP remains a separate agent protocol — see [MCP docs](https://docs.dogabot.com/mcp/).

## Install

```bash
go install github.com/dogabot/dogabot-cli/cmd/dogabot@latest
```

Requires Go 1.22+. `@latest` tracks the newest release tag; pin `@vX.Y.Z` for reproducible installs. Prebuilt GitHub Release binaries and Homebrew may follow later.

## Auth

```bash
export DOGABOT_API_KEY=dbk_live_...
# optional:
# export DOGABOT_BASE_URL=https://api.dogabot.com
# export DOGABOT_YES=1   # confirm writes without --yes
```

Or `~/.config/dogabot/config.json`:

```json
{ "api_key": "dbk_live_...", "base_url": "https://api.dogabot.com" }
```

Prefer env vars over putting keys in shell history via `--api-key`.

Color defaults to **auto** (TTY only). Use `--color=always|never`, or set non-empty `NO_COLOR` to disable.
`--json` never adds decorative color.

## Quick start

```bash
dogabot version
dogabot whoami --json
dogabot account get-me --json
dogabot markets get-markets --json
dogabot ops --tag terminal
```

Writes require confirmation and send an `Idempotency-Key` (auto-generated if omitted):

```bash
dogabot terminal post-terminal-placeorder \
  --body @order.json \
  --yes
```

Invoke by exact SDK `operationId`:

```bash
dogabot call getMe --json
dogabot call getFollowersById --id YOUR_ID --json
```

## CLI vs SDK vs MCP

| Surface | Best for |
|---------|----------|
| **CLI** (`dogabot`) | Shell scripts, terminals, one-off ops |
| **SDKs** | App / service code |
| **MCP** | AI agents (Cursor, Claude, …) — curated tool set |

Docs: https://docs.dogabot.com/cli/

## License

MIT
