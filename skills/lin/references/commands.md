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
- `lin project get <id>` — project summary with lead, labels, milestones, url, and content (markdown body, truncated by default)
- `lin project issues <id> [filters]` — issues in a project
- `lin project update title <id> <value>`
- `lin project update status <id> <value>` — backlog | planned | started | paused | completed | canceled
- `lin project update description <id> <value>`
- `lin project update content <id> <value>` — update the project's markdown body
- `lin project update lead <id> <user-id>`
- `lin project update start-date <id> <YYYY-MM-DD>`
- `lin project update target-date <id> <YYYY-MM-DD>`
- `lin project update priority <id> <value>` — none | urgent | high | medium | low
- `lin project update icon <id> <emoji>`
- `lin project update color <id> <hex>`
- `lin project update labels <id> <labels>` — replace project labels (comma-separated names or UUIDs; resolved via `lin label list --type project`). Replace semantics: previously-set labels not listed are removed.
- `lin project new <name> [--team <keys>] [--description <text>] [--lead <user>] [--start-date <YYYY-MM-DD>] [--target-date <YYYY-MM-DD>] [--status <status>] [--content <markdown>]`
- `lin project delete <id>` — delete project (move to trash)
- `lin project unarchive <id>` — restore trashed/archived project

## Initiatives

- `lin initiative search <text> [--limit] [--cursor]` — search initiatives by name
- `lin initiative list [--status] [--limit] [--cursor]` — list initiatives (status: planned | active | completed)
- `lin initiative get <id>` — initiative summary with status, health, owner, projects
- `lin initiative projects <id> [--limit] [--cursor]` — projects linked to an initiative
- `lin initiative new <name> [--status <status>] [--owner <user>] [--target-date <YYYY-MM-DD>]`
- `lin initiative update status <id> <value>` — planned | active | completed
- `lin initiative update target-date <id> <YYYY-MM-DD>`
- `lin initiative archive <id>` — archive an initiative
- `lin initiative unarchive <id>` — unarchive an initiative
- `lin initiative delete <id>` — delete an initiative

## Documents

- `lin document search <text> [--include-comments] [--include-archived] [--limit] [--cursor]` — full-text search
- `lin document list [--project <name|slug|id>] [--creator <name|email|id>] [--include-archived] [--limit] [--cursor]` — list documents
- `lin document get <id>` — full document details + markdown content (accepts UUID or slug ID)
- `lin document new <title> [--project <name|slug|id>] [--content <markdown>] [--icon <emoji>] [--color <hex>]`
- `lin document update title <id> <value>`
- `lin document update content <id> <value>`
- `lin document update project <id> <project>` — project ID, slug, or name
- `lin document update icon <id> <emoji>`
- `lin document update color <id> <hex>`

## Document History

- `lin document history <id>` — content edit history (actor IDs + timestamps, not paginated)

## Files

- `lin file upload <paths...>` — upload one or more files to Linear's CDN, returns `[{ filename, assetUrl, contentType }]`
- `lin file download <url-or-path> [--output <path>] [--output-dir <dir>] [--stdout] [--force]` — download a file from Linear's CDN

URL formats for download (all equivalent):

- Full URL: `https://uploads.linear.app/<org>/<uuid>/<uuid>`
- Domain + path: `uploads.linear.app/<org>/<uuid>/<uuid>`
- With org: `<org>/<uuid>/<uuid>`
- Without org: `<uuid>/<uuid>` (org inferred from auth)
- Single UUID: `<uuid>` (org inferred from auth)

Flags `--output`, `--output-dir`, and `--stdout` are mutually exclusive. Without `--force`, download refuses to overwrite existing files.

> **Prefer `--file` on comments** when attaching files to issues — `lin issue comment new ENG-123 "text" --file ./image.png` uploads and embeds in one step. Use `lin file upload` when you need a standalone asset URL.

## Issues

- `lin issue search <text> [filters]` — full-text search
- `lin issue list [filters]` — list issues (returns status, assignee, team, branchName)
  - Date filters: `--updated-after`, `--updated-before`, `--created-after`, `--created-before` (YYYY-MM-DD)
- `lin issue get <id>` — full issue details with commentCount, branchName, attachments (PR links)
- `lin issue new <title> --team <team> [--priority <p>] [--status <s>] [--assignee <name|email|id>] [--project <p>] [--labels <names|ids>]`
- `lin issue update title <id> <value>`
- `lin issue update status <id> <value>` — team-specific workflow state names
- `lin issue update assignee <id> <user>` — accepts name, email, or user ID
- `lin issue update priority <id> <priority>` — none | urgent | high | medium | low
- `lin issue update project <id> <project>` — project ID, slug, or name
- `lin issue update estimate <id> <value>` — validated against team estimate scale
- `lin issue update labels <id> <name1,name2,...>` — accepts label names or IDs
- `lin issue update description <id> <value>`
- `lin issue update due-date <id> <YYYY-MM-DD>`
- `lin issue update cycle <id> <cycle-id>` — move issue to a cycle (UUID)
- `lin issue update parent <id> <parent-id>` — set parent issue (make sub-issue)

## Comments

- `lin issue comment list <issue-id> [--limit] [--cursor]` — list comments with authors, parent ref, childCount (paginated)
- `lin issue comment new <issue-id> <body> [--parent <comment-id>] [--file <path>]` — add comment (--parent for threaded reply, 1 level max; --file repeatable for uploads)
- `lin issue comment get <comment-id>` — get comment with author, issue ref, parent ref, and childCount
- `lin issue comment edit <comment-id> <body> [--file <path>]` — edit a comment (--file repeatable)
- `lin issue comment replies <comment-id> [--limit] [--cursor]` — list replies to a comment (paginated)

## Issue Relations

- `lin issue relation list <issue-id>` — list all relations (both directions, includes blocked_by)
- `lin issue relation add <issue-id> --type <type> --related <related-id>` — type: blocks | duplicate | related
- `lin issue relation remove <relation-id>` — delete a relation

## Issue History

- `lin issue history <issue-id> [--limit] [--cursor]` — activity log: status changes, assignee changes, priority, labels, estimate, title, project, due date, archive/trash (paginated)

## Issue Lifecycle

- `lin issue archive <id>` — archive an issue
- `lin issue unarchive <id>` — unarchive an issue
- `lin issue delete <id>` — delete issue (move to trash)

## Issue Attachments

- `lin issue attachment list <issue-id>` — list attachments (any source type)
- `lin issue attachment add <issue-id> <url> [--title <text>]` — link a URL (default: rich link via `attachmentLinkURL` — server detects integrations)
  - `--github-pr` — force GitHub pull request integration (PR status sync)
  - `--github-issue` — force GitHub issue integration
  - `--gitlab-mr` — force GitLab merge request integration (project path + number derived from URL)
  - `--slack [--sync-thread]` — force Slack message integration; `--sync-thread` mirrors the Slack thread to a comment thread
  - `--discord` — force Discord message integration (channel + message IDs derived from URL)
- `lin issue attachment remove <attachment-id>` — remove any attachment, regardless of source type

## Teams

- `lin team list` — list teams (id, name, key)
- `lin team get <id>` — team details + members + estimate config (type, valid values)
- `lin team states <team>` — list workflow states (id, name, type, color, position)

## Users

- `lin user list [--team]` — list users
- `lin user me` — current authenticated user

## Labels

`IssueLabel` and `ProjectLabel` are distinct Linear entities; the same name can exist in both. Pass `--type project` to operate on project labels; the default is `issue`.

- `lin label list [--type issue|project] [--team] [--name <text>] [--is-group[=false]]` — list labels (filterable). Issue-label output includes `team{id,key,name}` and `parent{id,name}` when present; project-label output omits `team` (project labels are workspace-only).
- `lin label search <text> [--type issue|project] [--team]` — substring search by name (case- and accent-insensitive).
- `lin label get <id|name> [--type issue|project] [--team]` — single label by UUID or exact name. Use `--team` (or a UUID) to disambiguate when a name is shared across teams (issue labels only).
- `--team` is rejected with `--type=project` (project labels are workspace-scoped).

## Cycles

- `lin cycle list <team> [--current] [--next] [--previous]` — list cycles
- `lin cycle get <id>` — cycle details + issues

## Usage

- `lin usage` — print concise LLM-optimized top-level docs (~1000 tokens)
- `lin <command> usage` — print detailed docs for a specific command domain:
  - `lin issue usage` — issue search, list, create, update, comment, relation, archive, attachment
  - `lin project usage` — project commands
  - `lin initiative usage` — initiative commands
  - `lin document usage` — document search, list, create, update
  - `lin team usage` — team, user, label, cycle commands
  - `lin auth usage` — authentication + workspace management
  - `lin file usage` — file upload + download commands
  - `lin config usage` — CLI settings (keys, defaults, validation rules)

## API (raw GraphQL)

- `lin api query <graphql> [--variables <json>]` — execute raw GraphQL query (escape hatch for unsupported operations)

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
