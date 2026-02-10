# lin

Linear CLI for humans and LLMs.

- **Structured JSON output** — all output is JSON to stdout, errors to stderr
- **LLM-optimized** — `lin usage` prints concise docs in <1,000 tokens
- **Zero runtime deps** — single compiled binary via `bun build --compile`
- **Smart IDs** — accepts issue keys (`ENG-123`), UUIDs, or URL fragments

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
│   ├── get overview|details|issues <id>
│   ├── new <name> --team <teams>
│   └── update title|status|description|lead <id> <value>
├── roadmap
│   ├── list
│   └── get overview|projects <id>
├── document
│   ├── search <text>
│   ├── list [--project] [--creator]
│   ├── get <id>
│   ├── new <title> [--project] [--content]
│   └── update title|content|project <id> <value>
├── issue
│   ├── search <text>
│   ├── list [filters]
│   ├── get overview|comments <id>
│   ├── new <title> --team <team>
│   ├── update title|status|assignee|priority|project|labels|estimate|description <id> <value>
│   ├── comment new|get|edit <id> [<body>]
│   ├── relation list|add|remove
│   ├── archive|unarchive|delete <id>
│   └── attachment list|add|remove
├── team
│   ├── list
│   ├── get <id>
│   └── states <team>
├── user
│   ├── list
│   └── me
├── label
│   └── list
├── cycle
│   ├── list --team <team>
│   └── get <id>
└── usage                          # LLM-optimized docs (<1k tokens)
```

## Output

- Lists → JSON arrays
- Single items → JSON objects
- Errors → `{ "error": "message" }` to stderr + non-zero exit
- Empty/null fields are pruned automatically

## Filters

Most list/search commands accept:

`--team`, `--status`, `--assignee`, `--priority`, `--label`, `--cycle`, `--project`, `--limit`, `--cursor`

## Development

```bash
bun install
bun run dev -- --help        # run in dev mode
bun run typecheck             # type check
bun test                      # run tests
bun run lint                  # lint
```

## License

MIT
