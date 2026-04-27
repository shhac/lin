# lin

Linear CLI for humans and LLMs.

- **Structured JSON output** — all output is JSON to stdout, errors to stderr
- **LLM-optimized** — `lin usage` prints concise docs in <1,000 tokens
- **Zero runtime deps** — single static binary via `go build`
- **Smart IDs** — accepts issue keys (`ENG-123`), UUIDs, or URL fragments

**Website:** [lin.paulie.app](https://lin.paulie.app/)

## Installation

```bash
brew install shhac/tap/lin
```

### Claude Code / AI agent skill

```bash
npx skills add shhac/lin
```

This installs the `lin` skill so Claude Code (and other AI agents) can discover and use `lin` automatically. See [skills.sh](https://skills.sh) for details.

## Authentication

Generate a personal API key at **Settings > Account > Security > [Personal API Keys](https://linear.app/settings/account/security)**.

Set an env var (recommended for CI / agent use):

```bash
export LINEAR_API_KEY="lin_api_..."
```

Or store it locally:

```bash
lin auth login <api-key>
lin auth status
```

## Command map

```text
lin
├── auth
│   ├── login <api-key> [--alias <name>]
│   ├── logout [--all]
│   ├── status
│   └── workspace list|switch|remove
├── project
│   ├── search <text>
│   ├── list [--team] [--status]
│   ├── get <id>
│   ├── issues <id> [filters]
│   ├── new <name> --team <teams>
│   ├── update title|status|description|content|lead|start-date|target-date|priority|icon|color <id> <value>
│   ├── delete|unarchive <id>
│   └── usage
├── initiative
│   ├── search <text>
│   ├── list [--status planned|active|completed]
│   ├── get <id>
│   ├── projects <id>
│   ├── new <name> [options]
│   ├── update name|description|owner|status|target-date|content|color|icon <id> <value>
│   ├── archive|unarchive|delete <id>
│   └── usage
├── document
│   ├── search <text>
│   ├── list [--project] [--creator]
│   ├── get <id>
│   ├── new <title> [--project] [--content]
│   ├── update title|content|project|icon|color <id> <value>
│   ├── history <id>
│   └── usage
├── file
│   ├── download <url-or-path> [--output] [--output-dir] [--stdout] [--force]
│   ├── upload <paths...>
│   └── usage
├── issue
│   ├── search <text>
│   ├── list [filters]
│   ├── get <id>
│   ├── new <title> --team <team>
│   ├── update title|status|assignee|priority|project|labels|estimate|description|due-date|cycle|parent <id> <value>
│   ├── comment list|new|get|edit|replies <id> [<body>]
│   ├── relation list|add|remove
│   ├── archive|unarchive|delete <id>
│   ├── attachment list|add|remove   # add: --github-pr|--github-issue|--gitlab-mr|--slack [--sync-thread]|--discord
│   ├── history <id> [--limit] [--cursor]
│   └── usage
├── team
│   ├── list
│   ├── get <id>
│   ├── states <team>
│   └── usage
├── user
│   ├── list
│   ├── me
│   └── usage
├── label
│   ├── list
│   └── usage
├── cycle
│   ├── list <team>
│   ├── get <id>
│   └── usage
├── api
│   ├── query <graphql> [--variables <json>]
│   └── usage
├── config
│   ├── get|set|reset|list-keys
│   └── usage
└── usage                          # LLM-optimized docs (<1k tokens)
```

Each top-level command also has a `usage` subcommand for detailed, LLM-friendly documentation (e.g., `lin issue usage`, `lin project usage`). The top-level `lin usage` gives a broad overview; per-command usage gives full detail on flags, valid values, and return fields.

## Output

- Lists → `{ "items": [...], "pagination"?: { "hasMore", "nextCursor" } }`
- Single items → JSON objects
- Errors → `{ "error": "...", "fixable_by": "agent|human|retry", "hint": "..." }` to stderr + non-zero exit
- Empty/null fields are pruned automatically

## Filters

Most list/search commands accept:

`--team`, `--status`, `--assignee`, `--priority`, `--label`, `--cycle`, `--project`, `--limit`, `--cursor`

## Development

```bash
make dev ARGS="--help"       # run in dev mode
make build                   # build binary
make test                    # run tests
make generate                # regenerate GraphQL code
make lint                    # golangci-lint
```

## License

MIT
