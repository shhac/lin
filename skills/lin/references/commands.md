# `lin` command map (reference)

Run `lin --help` (or `lin <command> --help`) for the full option list.
Run `lin usage` for concise LLM-optimized docs.
Run `lin <command> usage` for detailed per-command docs (e.g., `lin issue usage`).

## Auth

- `lin auth login <api-key> [--alias <name>]` — validate and store API key (auto-detects workspace)
- `lin auth logout [--all]` — remove active workspace credentials (--all: clear all workspaces)
- `lin auth status` — show auth state, active workspace, and other stored workspaces
- `lin auth workspace list` — list all stored workspaces
- `lin auth workspace switch <alias>` — set default workspace
- `lin auth workspace remove <alias>` — remove a stored workspace

## Projects

- `lin project search <text>` — search projects by name
- `lin project list [--team] [--status] [--limit] [--cursor]` — list projects
- `lin project get overview <id>` — project summary with lead, milestones, url, and content (markdown body, truncated by default)
- `lin project get issues <id> [filters]` — issues in a project
- `lin project update title <id> <value>`
- `lin project update status <id> <value>` — backlog | planned | started | paused | completed | canceled
- `lin project update description <id> <value>`
- `lin project update lead <id> <user-id>`
- `lin project new <name> [--team <keys>] [--description <text>] [--lead <user>] [--start-date <YYYY-MM-DD>] [--target-date <YYYY-MM-DD>] [--status <status>] [--content <markdown>]`

## Roadmaps

- `lin roadmap list [--limit] [--cursor]` — list roadmaps
- `lin roadmap get overview <id>` — roadmap summary with owner, creator, url
- `lin roadmap get projects <id> [--limit] [--cursor]` — projects linked to a roadmap

## Documents

- `lin document search <text> [--include-comments] [--include-archived] [--limit] [--cursor]` — full-text search
- `lin document list [--project <name|slug|id>] [--creator <name|email|id>] [--include-archived] [--limit] [--cursor]` — list documents
- `lin document get <id>` — full document details + markdown content (accepts UUID or slug ID)
- `lin document new <title> [--project <name|slug|id>] [--content <markdown>] [--icon <emoji>] [--color <hex>]`
- `lin document update title <id> <value>`
- `lin document update content <id> <value>`
- `lin document update project <id> <project>` — project ID, slug, or name

## Issues

- `lin issue search <text> [filters]` — full-text search
- `lin issue list [filters]` — list issues (returns status, assignee, team, branchName)
- `lin issue get overview <id>` — full issue details with commentCount, branchName, attachments (PR links)
- `lin issue get comments <id>` — list comments with authors
- `lin issue new <title> --team <team> [--priority <p>] [--status <s>] [--assignee <name|email|id>] [--project <p>] [--label <l>]`
- `lin issue update title <id> <value>`
- `lin issue update status <id> <value>` — team-specific workflow state names
- `lin issue update assignee <id> <user>` — accepts name, email, or user ID
- `lin issue update priority <id> <priority>` — none | urgent | high | medium | low
- `lin issue update project <id> <project>` — project ID, slug, or name
- `lin issue update estimate <id> <value>` — validated against team estimate scale
- `lin issue update labels <id> <label1,label2,...>`
- `lin issue update description <id> <value>`

## Comments

- `lin issue comment new <issue-id> <body> [--parent <comment-id>] [--file <path>]` — add comment (--parent for threaded reply, 1 level max; --file repeatable for uploads)
- `lin issue comment get <comment-id>` — get comment with author, issue ref, parent ref, and childCount
- `lin issue comment edit <comment-id> <body> [--file <path>]` — edit a comment (--file repeatable)
- `lin issue comment replies <comment-id> [--limit] [--cursor]` — list replies to a comment (paginated)

## Issue Relations

- `lin issue relation list <issue-id>` — list all relations (both directions, includes blocked_by)
- `lin issue relation add <issue-id> --type <type> --related <related-id>` — type: blocks | duplicate | related
- `lin issue relation remove <relation-id>` — delete a relation

## Issue Lifecycle

- `lin issue archive <id>` — archive an issue
- `lin issue unarchive <id>` — unarchive an issue
- `lin issue delete <id>` — delete issue (move to trash)

## Issue Attachments

- `lin issue attachment list <issue-id>` — list attachments
- `lin issue attachment add <issue-id> --url <url> --title <title> [--subtitle <text>]` — add URL attachment
- `lin issue attachment remove <attachment-id>` — remove attachment

## Teams

- `lin team list` — list teams (id, name, key)
- `lin team get <id>` — team details + members + estimate config (type, valid values)
- `lin team states <team>` — list workflow states (id, name, type, color, position)

## Users

- `lin user list [--team]` — list users
- `lin user me` — current authenticated user

## Labels

- `lin label list [--team]` — list labels

## Cycles

- `lin cycle list --team <team> [--current] [--next] [--previous]` — list cycles
- `lin cycle get <id>` — cycle details + issues

## Usage

- `lin usage` — print concise LLM-optimized top-level docs (~1000 tokens)
- `lin <command> usage` — print detailed docs for a specific command domain:
  - `lin issue usage` — issue search, list, create, update, comment, relation, archive, attachment
  - `lin project usage` — project + roadmap commands
  - `lin document usage` — document search, list, create, update
  - `lin team usage` — team, user, label, cycle commands
  - `lin auth usage` — authentication + workspace management
  - `lin config usage` — CLI settings (keys, defaults, validation rules)

## Common filters (list/search commands)

| Flag                           | Description                           |
| ------------------------------ | ------------------------------------- |
| `--team <key\|name>`           | Filter by team                        |
| `--status <name>`              | Filter by workflow state              |
| `--assignee <name\|email\|id>` | Filter by assignee                    |
| `--priority <level>`           | none, urgent, high, medium, low       |
| `--label <name>`               | Filter by label                       |
| `--cycle <id>`                 | Filter by cycle                       |
| `--project <name\|slug\|id>`   | Filter by project (ID, slug, or name) |
| `--limit <n>`                  | Max results per page                  |
| `--cursor <token>`             | Pagination cursor                     |

## Global flags

| Flag                   | Description                                                         |
| ---------------------- | ------------------------------------------------------------------- |
| `--expand <field,...>` | Expand specific truncated fields (e.g. `--expand description,body`) |
| `--full`               | Expand all truncated fields                                         |
