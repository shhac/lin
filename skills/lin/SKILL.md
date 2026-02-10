---
name: lin
description: |
  Linear CLI for humans and LLMs. Use when:
  - Looking up Linear issues, projects, cycles, or teams
  - Searching Linear issues or projects by text
  - Creating or updating Linear issues
  - Adding comments to Linear issues
  - Checking project status, milestones, or team members
  Triggers: "linear issue", "linear project", "linear ticket", "linear search", "create issue", "update issue", "linear team", "linear cycle"
---

# Linear automation with `lin`

`lin` is a CLI binary installed on `$PATH`. Invoke it directly (e.g. `lin issue list --team ENG`).

All output is JSON to stdout. Errors go to stderr as `{ "error": "..." }` with non-zero exit.

## Quick start (auth)

Set an env var (recommended):

```bash
export LINEAR_API_KEY="lin_api_..."
```

Or store it locally:

```bash
lin auth login <api-key>
lin auth status
```

Generate a personal API key at **Settings > Account > Security > Personal API Keys** in the Linear app.

## Looking up issues

```bash
lin issue search "auth bug"
lin issue list --team ENG --status "In Progress" --assignee "alice@example.com"
lin issue get overview ENG-123
lin issue get comments ENG-123
```

## Creating and updating issues

```bash
lin issue new "Fix login redirect" --team ENG --priority high --status "Todo"
lin issue update status ENG-123 "In Progress"
lin issue update assignee ENG-123 "alice@example.com"
lin issue update priority ENG-123 urgent
lin issue comment new ENG-123 "Started investigating"
```

## Projects

```bash
lin project search "migration"
lin project list --status started
lin project get overview <id>
lin project get details <id>       # full markdown content
lin project get issues <id>
lin project update status <id> completed
```

## Teams, users, labels, cycles

```bash
lin team list
lin team get ENG
lin user me
lin user list --team ENG
lin label list --team ENG
lin cycle list --team ENG --current
lin cycle get <id>
```

## IDs

All commands accept multiple ID formats:

- Issue keys: `ENG-123`
- UUIDs: `aaaaaaaa-1111-2222-3333-444444444444`
- URL slugs: `fix-login-redirect-abc123`
- `--team` accepts team key (`ENG`) or name (`Engineering`)

## Pagination

List commands return `{ "items": [...], "pagination"?: { "hasMore": true, "nextCursor": "..." } }`.

Use `--limit <n>` and `--cursor <token>` to paginate.

## References

- [references/commands.md](references/commands.md): full command map + all flags
- [references/output.md](references/output.md): JSON output shapes + field details
