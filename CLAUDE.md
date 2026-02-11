# lin

Linear CLI for humans and LLMs. TypeScript + Bun, compiled to standalone binaries.

## Architecture

```
src/
├── index.ts                 # CLI entry — registers all commands via commander
├── cli/
│   ├── auth-command.ts      # auth login / auth status
│   ├── project-command.ts   # project search/list/get/update
│   ├── issue-command.ts     # issue search/list/get/new/update
│   ├── comment-command.ts   # issue comment new/get/edit
│   ├── team-command.ts      # team list/get
│   ├── user-command.ts      # user list/me
│   ├── label-command.ts     # label list
│   ├── cycle-command.ts     # cycle list/get
│   └── usage-command.ts     # LLM-optimized usage text (<1k tokens)
└── lib/
    ├── client.ts            # LinearClient factory (exits if no API key)
    ├── config.ts            # ~/.config/lin/ config + API key storage
    ├── filters.ts           # Shared buildIssueFilter() for CLI opts → Linear SDK filters
    ├── output.ts            # pruneEmpty, printJson, printPaginated, printError
    └── version.ts           # Version from build-time define / env / package.json
```

## Key patterns

- **Command registration**: Each `*-command.ts` exports `registerXyzCommand({ program })` called from `index.ts`
- **Output**: All commands use `printJson()` or `printPaginated()` from `lib/output.ts`. Errors use `printError()`. All output is JSON, empty/null fields auto-pruned.
- **Pagination**: List commands return `{ items: [...], pagination?: { hasMore, nextCursor } }` via `printPaginated()`
- **Filters**: `lib/filters.ts` builds Linear SDK `IssueFilter` objects from CLI flags (`--team`, `--status`, `--assignee`, etc.)
- **Auth**: `LINEAR_API_KEY` env var takes precedence, otherwise stored in `~/.config/lin/config.json`
- **Error messages**: Include valid values so LLMs can self-correct (e.g., `"Unknown status: X. Valid values: ..."`)
- **Usage subcommands**: Each command has a `usage` subcommand (`src/cli/*/usage.ts`) providing LLM-friendly docs. When modifying a command's behavior, options, or flags, update its usage text too. Sub-usage texts are tested to stay under 500 tokens each.

## Commands

Run `bun run dev -- usage` for the full command reference. Each command also supports `<command> usage` for detailed per-command docs.

## Development

```bash
bun install
bun run dev -- <command>     # run in dev mode
bun run typecheck            # tsc --noEmit
bun test                     # bun:test
bun run lint                 # oxlint
bun run format               # oxfmt
```

## Release

```bash
bun run release patch        # bumps version, commits, tags, pushes
bun run build:release        # cross-platform binaries in release/
```

Then create GitHub release, update homebrew-tap formula at `shhac/homebrew-tap` with new sha256s.

## Conventions

- TypeScript strict mode, ES2022 target, Bun bundler resolution
- `type` over `interface` (enforced by oxlint)
- kebab-case filenames (enforced by oxlint)
- Max 350 lines per file, max 2 params per function (oxlint warnings)
- Pre-commit hook: oxlint fix + oxfmt
- Tests: bun:test, no mocking libraries, inline fixtures, pure functions preferred
