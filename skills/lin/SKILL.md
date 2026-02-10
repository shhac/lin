---
name: lin
description: |
  Linear CLI for humans and LLMs. Use when:
  - Looking up Linear issues, projects, documents, cycles, or teams
  - Searching Linear issues, projects, or documents by text
  - Creating or updating Linear issues or documents
  - Adding comments to Linear issues
  - Checking project status, milestones, or team members
  Triggers: "linear issue", "linear project", "linear document", "linear ticket", "linear search", "create issue", "create document", "update issue", "update document", "linear team", "linear cycle"
---

# Linear automation with `lin`

`lin` is a CLI binary installed on `$PATH`. Invoke it directly (e.g. `lin issue list --team ENG`).

All output is JSON to stdout. Errors go to stderr as `{ "error": "..." }` with non-zero exit.

## Quick start (auth)

Set an env var (recommended):

```bash
export LINEAR_API_KEY="lin_api_..."
```

Or store it locally (supports multiple workspaces — each Linear API key is per-org):

```bash
lin auth login <api-key> --alias myorg
lin auth status
lin auth workspace list
lin auth workspace switch <alias>
lin auth logout [--all]
lin auth workspace remove <alias>
```

Generate a personal API key at **Settings > Account > Security > Personal API Keys** in the Linear app.

## Looking up issues

```bash
lin issue search "auth bug"
lin issue list --team ENG --status "In Progress" --assignee "alice@example.com"
lin issue get overview ENG-123    # includes branchName, commentCount, attachments (PR links)
lin issue get comments ENG-123
```

## Creating and updating issues

```bash
lin issue new "Fix login redirect" --team ENG --priority high --status "Todo"
lin issue update status ENG-123 "In Progress"
lin issue update assignee ENG-123 "alice@example.com"
lin issue update priority ENG-123 urgent
lin issue update estimate ENG-123 5      # validated against team's estimate scale
lin issue comment new ENG-123 "Started investigating"
lin issue comment get <comment-id>
lin issue comment edit <comment-id> "Updated analysis"
```

## Issue relations and lifecycle

```bash
# Relations (blocks, duplicate, related)
lin issue relation list ENG-123
lin issue relation add ENG-123 --type blocks --related ENG-124
lin issue relation remove <relation-id>

# Archive and delete
lin issue archive ENG-123
lin issue unarchive ENG-123
lin issue delete ENG-123                  # moves to trash

# Attachments (link PRs, docs, etc.)
lin issue attachment list ENG-123
lin issue attachment add ENG-123 --url "https://github.com/org/repo/pull/456" --title "PR #456" --subtitle "Fixes login bug"
lin issue attachment remove <attachment-id>
```

## Projects

Project commands accept UUID, slug ID, or name.

```bash
lin project search "migration"
lin project list --status started
lin project get overview "CRM Actions"   # accepts UUID, slug, or name
lin project get details <id>             # full markdown content
lin project get issues <id>
lin project new "New Feature" --team ENG --status planned --lead "alice@example.com"
lin project update status <id> completed
```

The `--project` filter on issue commands also accepts slug or name.

## Roadmaps

```bash
lin roadmap list
lin roadmap get overview <id>            # roadmap summary + owner
lin roadmap get projects <id>            # projects linked to a roadmap
```

## Documents

Document commands accept UUID or slug ID.

```bash
lin document search "onboarding"
lin document list --project "CRM Actions" --creator "alice@example.com"
lin document get <id>                    # full markdown content
lin document new "API Design Doc" --project "CRM Actions" --content "# Overview\n..."
lin document update title <id> "New Title"
lin document update content <id> "# Updated content"
lin document update project <id> "Other Project"
```

## Teams, users, labels, cycles

```bash
lin team list
lin team get ENG                         # includes estimate config + valid values
lin team states ENG                      # workflow states (discover valid status values)
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
