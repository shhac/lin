# lin

Linear CLI for humans and LLMs. Go + cobra + genqlient, compiled to static binaries.

## Architecture

```
cmd/lin/main.go              # CLI entry point
internal/
├── cli/                     # Cobra command packages (one per domain)
│   ├── root.go              # Root command, global flags (--expand, --full)
│   ├── issue/               # Largest: search/list/get/new/update/comment/relation/archive/attachment/history
│   ├── project/             # search/list/get/issues/new/update/delete/unarchive
│   ├── initiative/          # search/list/get/projects/new/update/archive/unarchive/delete
│   ├── document/            # search/list/get/new/update/history
│   └── ...                  # api, auth, configcmd, cycle, file, label, team, usage, user
├── linear/                  # GraphQL client layer
│   ├── graphql.go           # GetClient() → authenticated graphql.Client
│   ├── client.go            # Raw HTTP client for `api query`
│   ├── generated.go         # genqlient output (DO NOT EDIT)
│   ├── paginate.go          # FetchAll[T] generic pagination helper
│   └── queries/*.graphql    # GraphQL operations by domain
├── errors/                  # Structured APIError with FixableBy + Hint + GraphQL classification
├── config/                  # ~/.config/lin/ config management
├── credential/              # API key resolution (env → keychain → config)
├── filters/                 # CLI flags → Linear SDK filter objects
├── resolvers/               # Human input → Linear entity ID (per-entity files)
├── mappers/                 # Shared output mappers (issue, project, document summaries)
├── output/                  # JSON printing, pruning, pagination, truncation, error formatting
├── ptr/                     # Generic pointer helpers
├── upload/ + download/      # File operations
└── estimates/ priorities/   # Domain constants and validation
    projectstatuses/
    truncation/
```

## Key patterns

- **Command registration**: `Register(parent *cobra.Command)` per domain, wired in `root.go`
- **Output**: All JSON via `output.PrintJSON()` / `output.PrintPaginated()` / `output.PrintError()`
- **Errors**: Structured `{error, fixable_by, hint}` via `output.WriteError()` / `output.HandleGraphQLError()`
- **Filters**: `filters.BuildIssueFilter(opts)` → typed `*linear.IssueFilter`
- **Resolvers**: `resolvers.ResolveTeam(client, input)` — accepts UUID, key, or name
- **Pagination**: User-facing: `--limit`/`--cursor`. Internal: `linear.FetchAll[T]` for resolvers
- **Auth**: `LINEAR_API_KEY` env → macOS Keychain → config file

## Development

```bash
make dev ARGS="<command>"    # run in dev mode
make build                   # build binary
make test                    # run tests
make generate                # regenerate GraphQL code
make lint                    # golangci-lint
```

## Documentation checklist

When adding or changing commands, flags, or output shapes, update:

- `internal/cli/*/usage.go` — LLM-optimized per-command usage text
- `internal/cli/usage/usage.go` — top-level usage overview
- `skills/lin/SKILL.md` — main skill doc with examples
- `skills/lin/references/commands.md` — full command map
- `skills/lin/references/output.md` — JSON output shapes
- `README.md` — command map tree
