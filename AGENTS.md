# lin

Linear CLI for humans and LLMs. Go + cobra + genqlient, compiled to static binaries.

## Architecture

```
cmd/lin/main.go              # CLI entry point
internal/
├── cli/                     # Cobra command packages (one per domain)
│   ├── root.go              # Root command, global flags (--expand, --full)
│   ├── issue/               # Largest: search/list/get/new/update/comment/relation/archive/attachment/history
│   ├── project/             # search/list/get/issues/new/update
│   ├── document/            # search/list/get/new/update/history
│   └── ...                  # api, auth, configcmd, cycle, file, label, roadmap, team, usage, user
├── linear/                  # GraphQL client layer
│   ├── graphql.go           # GetClient() → authenticated graphql.Client
│   ├── client.go            # Raw HTTP client for `api query`
│   ├── generated.go         # genqlient output (DO NOT EDIT)
│   └── queries/*.graphql    # GraphQL operations by domain
├── config/                  # ~/.config/lin/ config management
├── credential/              # API key resolution (env → keychain → config)
├── filters/                 # CLI flags → Linear SDK filter objects
├── resolvers/               # Human input → Linear entity ID (per-entity files)
├── mappers/                 # Shared output mappers (issue, project, document summaries)
├── output/                  # JSON printing, pruning, pagination, truncation
├── ptr/                     # Generic pointer helpers
├── upload/ + download/      # File operations
└── estimates/ priorities/   # Domain constants and validation
    projectstatuses/
    truncation/
```

## Key patterns

- **Command registration**: `Register(parent *cobra.Command)` per domain, wired in `root.go`
- **Output**: All JSON via `output.PrintJSON()` / `output.PrintPaginated()` / `output.PrintError()`
- **Filters**: `filters.BuildIssueFilter(opts)` → typed `*linear.IssueFilter`
- **Resolvers**: `resolvers.ResolveTeam(client, input)` — accepts UUID, key, or name
- **Auth**: `LINEAR_API_KEY` env → macOS Keychain → config file
- **Error messages**: Include valid values for LLM self-correction

## Development

```bash
make dev ARGS="<command>"    # run in dev mode
make build                   # build binary
make test                    # run tests
make generate                # regenerate GraphQL code
```
