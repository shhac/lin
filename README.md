# lin

Linear CLI for humans and LLMs.

- **Structured output** — JSONL for lists/searches, JSON for single items, errors to stderr
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
npx skills add shhac/agent-skills --skill lin --global
```

Installs the `lin` skill globally so Claude Code (and other AI agents) can discover and use it automatically. It ships from [`shhac/agent-skills`](https://github.com/shhac/agent-skills) — the whole family's skills in one repo, so `npx skills update` checks a single source no matter how many you use. Want several at once? Run `npx skills add shhac/agent-skills --global` and pick from the list.

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
│   ├── list [--team] [--status] [--lead]
│   ├── get <id>...
│   ├── issues <id> [filters]
│   ├── requests <id> [--important]
│   ├── new <name> --team <teams>
│   ├── update title|status|description|content|lead|start-date|target-date|priority|icon|color|labels <id> <value>
│   ├── post new|list|get <project|update-id> [<body>] [--health]
│   ├── delete|unarchive <id>
│   └── usage
├── initiative
│   ├── search <text>
│   ├── list [--status planned|active|completed]
│   ├── get <id>...
│   ├── projects <id>
│   ├── new <name> [options]
│   ├── update name|description|owner|status|target-date|content|color|icon <id> <value>
│   ├── archive|unarchive|delete <id>
│   └── usage
├── document
│   ├── search <text>
│   ├── list [--project] [--creator]
│   ├── get <id>...
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
│   ├── get <id>...
│   ├── new <title> --team <team>
│   ├── update title|status|assignee|priority|project|labels|estimate|description|due-date|cycle|parent <id> <value>
│   ├── comment list|new|get|edit|replies <id> [<body>]
│   ├── relation list|add|remove
│   ├── requests <id> [--important]
│   ├── archive|unarchive|delete <id>
│   ├── attachment list|add|remove   # add: --github-pr|--github-issue|--gitlab-mr|--slack [--sync-thread]|--discord
│   ├── history <id> [--limit] [--cursor]
│   └── usage
├── customer
│   ├── list [--tier] [--status] [--owner] [--domain] [--revenue]
│   ├── search <text>
│   ├── get <id|slug>...
│   ├── requests [--customer] [--project] [--important] [--unassigned] [--triage] [--status] [--label] [--team] [--created-after|before]
│   ├── statuses
│   ├── tiers
│   └── usage
├── team
│   ├── list
│   ├── get <id>...
│   ├── states <team>
│   └── usage
├── user
│   ├── list
│   ├── me
│   └── usage
├── label
│   ├── list [--type issue|project] [--team] [--name] [--is-group]
│   ├── search <text> [--type issue|project] [--team]
│   ├── get <id|name>... [--type issue|project] [--team]
│   └── usage
├── cycle
│   ├── list <team>
│   ├── get <id>...
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

`lin mcp` runs lin as an MCP server. `file download` defaults to the lin cache
(`~/.cache/lin/downloads`); over MCP a built-in read-only `fs` tool lets a client
read those files back without filesystem access — `fs get cache downloads/<name>`
returns the bytes (images as image blocks).

**Multi-user MCP.** One `lin mcp` server can serve several people. Each named
principal's tool calls run pinned to their own stored workspace via the
`--workspace` selector, with `LIN_REQUIRE_IDENTITY` set so a call missing its
identity fails closed instead of falling back to the operator's default. A
principal is bound to a workspace either explicitly
(`mcp pair add <name> --bind workspace=<alias>`) or by self-enrollment: a
principal minted without a binding is prompted, during the browser OAuth
approval, to paste their own Linear API key. The key is validated, stored under
the principal's name (keychain-first), and the binding is written automatically.
A slot converges on one Linear organization — a key resolving to a different org
is refused rather than silently re-pointed.

## Output

- Lists/searches → JSONL by default, one object per line
- Pagination in JSONL → `{"@pagination":{"has_more":true,"next_cursor":"..."}}`
- `get <id>...` → NDJSON by default (one line per id — the record, or `{"@unresolved":{...}}` for a missing id); pass `--format json` for a pretty object
- `--format json|yaml|jsonl` overrides any command; JSON list/get envelopes use `{ "data": [...], "pagination"?: ... }` / `{ "data": [...], "@unresolved": [...] }`
- `--format pretty` (get commands: issue, project, initiative, document, customer) → human-readable terminal card for reading an entity (not for scripting); `--width <n>` sets width, and `--full` adds relations + comments for `issue get`
- Item-level misses (not found) → `@unresolved` line on stdout, exit 0; command-level failures → stderr, exit 1
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

PolyForm Perimeter License 1.0.0 — see [LICENSE](LICENSE).
