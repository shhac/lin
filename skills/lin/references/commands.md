# `lin` command map (reference)

Run `lin --help` (or `lin <command> --help`) for the full option list.
Run `lin usage` for concise LLM-optimized docs.

## Auth

- `lin auth login <api-key>` — validate and store API key
- `lin auth status` — show auth state, user, and workspace

## Projects

- `lin project search <text>` — search projects by name
- `lin project list [--team] [--status] [--limit] [--cursor]` — list projects
- `lin project get overview <id>` — project summary with lead, milestones, url
- `lin project get details <id>` — full project content (markdown body)
- `lin project get issues <id> [filters]` — issues in a project
- `lin project update title <id> <value>`
- `lin project update status <id> <value>` — backlog | planned | started | paused | completed | canceled
- `lin project update description <id> <value>`
- `lin project update lead <id> <user-id>`

## Issues

- `lin issue search <text> [filters]` — full-text search
- `lin issue list [filters]` — list issues (returns status, assignee, team inline)
- `lin issue get overview <id>` — full issue details with url, team, project, labels
- `lin issue get comments <id>` — list comments with authors
- `lin issue new <title> --team <team> [--priority <p>] [--status <s>] [--assignee <a>] [--project <p>] [--label <l>]`
- `lin issue update title <id> <value>`
- `lin issue update status <id> <value>` — team-specific workflow state names
- `lin issue update assignee <id> <user-id>`
- `lin issue update priority <id> <priority>` — none | urgent | high | medium | low
- `lin issue update project <id> <project-id>`
- `lin issue update labels <id> <label1,label2,...>`
- `lin issue update description <id> <value>`

## Comments

- `lin issue comment new <issue-id> <body>` — add comment
- `lin issue comment get <comment-id>` — get a specific comment
- `lin issue comment edit <comment-id> <body>` — edit a comment

## Teams

- `lin team list` — list teams (id, name, key)
- `lin team get <id>` — team details + members

## Users

- `lin user list [--team]` — list users
- `lin user me` — current authenticated user

## Labels

- `lin label list [--team]` — list labels

## Cycles

- `lin cycle list --team <team> [--current] [--next] [--previous]` — list cycles
- `lin cycle get <id>` — cycle details + issues

## Usage

- `lin usage` — print concise LLM-optimized docs (~700 tokens)

## Common filters (list/search commands)

| Flag                           | Description                     |
| ------------------------------ | ------------------------------- |
| `--team <key\|name>`           | Filter by team                  |
| `--status <name>`              | Filter by workflow state        |
| `--assignee <name\|email\|id>` | Filter by assignee              |
| `--priority <level>`           | none, urgent, high, medium, low |
| `--label <name>`               | Filter by label                 |
| `--cycle <id>`                 | Filter by cycle                 |
| `--project <name\|id>`         | Filter by project               |
| `--limit <n>`                  | Max results per page            |
| `--cursor <token>`             | Pagination cursor               |
