# Output format (reference)

## General

All commands print JSON to stdout. Errors print `{ "error": "..." }` to stderr with non-zero exit.

Empty/null fields are pruned automatically — missing keys mean no value, not `null`.

Error messages include valid values when input is invalid (e.g., unknown status names list all valid states).

## Truncation

Fields named `description`, `body`, or `content` are truncated to ~200 characters by default. A companion `*Length` field always shows the full character count.

**Default (truncated):**

```json
{
  "description": "This is the beginning of a long project description that goes on for many paragraphs...",
  "descriptionLength": 1847
}
```

**With `--full` or `--expand description` (expanded):**

```json
{
  "description": "This is the beginning of a long project description that goes on for many paragraphs and includes detailed specifications, requirements, and implementation notes...",
  "descriptionLength": 1847
}
```

The `*Length` field is always present when the source field has content, regardless of truncation. Use it to detect whether content was truncated (`description.length < descriptionLength`).

Truncatable fields: `description`, `body`, `content`. Global flags: `--expand <field,...>` or `--full`.

## List output

List commands return:

```json
{
  "items": [ ... ],
  "pagination": {
    "hasMore": true,
    "nextCursor": "abc123"
  }
}
```

When there are no more pages, the `pagination` key is omitted entirely.

## Single item output

Single-item commands (e.g., `issue get overview`, `team get`) return the object directly:

```json
{
  "id": "...",
  "title": "...",
  "status": "In Progress"
}
```

## Document list items

Documents in list/search output:

```json
{
  "id": "...",
  "slugId": "api-design-doc-a1b2c3",
  "title": "API Design Doc",
  "url": "https://linear.app/.../document/...",
  "project": { "id": "...", "name": "CRM Actions" },
  "creator": "Alice Example",
  "updatedAt": "2025-01-15T10:30:00.000Z"
}
```

## Document detail (`document get`)

Includes full markdown content:

```json
{
  "id": "...",
  "slugId": "api-design-doc-a1b2c3",
  "title": "API Design Doc",
  "content": "# Overview\n\nThis document describes...",
  "url": "https://linear.app/.../document/...",
  "icon": "📄",
  "color": "#5e6ad2",
  "project": { "id": "...", "name": "CRM Actions", "slugId": "crm-actions-d0f9" },
  "creator": { "id": "...", "name": "Alice Example" },
  "updatedBy": { "id": "...", "name": "Bob Example" },
  "createdAt": "2025-01-10T09:00:00.000Z",
  "updatedAt": "2025-01-15T10:30:00.000Z"
}
```

## Issue list items

Issues in list output include inline context to reduce follow-up calls:

```json
{
  "id": "...",
  "identifier": "ENG-123",
  "title": "Fix login redirect",
  "branchName": "alice/eng-123-fix-login-redirect",
  "status": "In Progress",
  "priority": "high",
  "assignee": "Alice Example",
  "team": "Engineering"
}
```

## Issue overview (`issue get overview`)

Includes comment count, branch name, and attachments (e.g., linked GitHub PRs):

```json
{
  "id": "...",
  "identifier": "ENG-123",
  "title": "Fix login redirect",
  "branchName": "alice/eng-123-fix-login-redirect",
  "commentCount": 3,
  "attachments": [{ "title": "PR #456", "url": "https://github.com/...", "sourceType": "github" }]
}
```

## Project list items

```json
{
  "id": "...",
  "slugId": "shipstation-integration-d0f95990a8f1",
  "name": "Shipstation Integration",
  "description": "One-liner project summary",
  "status": "started",
  "progress": 0.45,
  "url": "https://linear.app/.../project/..."
}
```

Use `project get overview <id>` with `--expand content` or `--full` for the full markdown body.

## Roadmap list items

```json
{
  "id": "...",
  "slugId": "q1-2025-roadmap-a1b2c3",
  "url": "https://linear.app/.../roadmap/...",
  "name": "Q1 2025 Roadmap",
  "description": "Key initiatives for Q1",
  "owner": "Alice Example"
}
```

## Roadmap overview (`roadmap get overview`)

```json
{
  "id": "...",
  "slugId": "q1-2025-roadmap-a1b2c3",
  "url": "https://linear.app/.../roadmap/...",
  "name": "Q1 2025 Roadmap",
  "description": "Key initiatives for Q1",
  "owner": { "id": "...", "name": "Alice Example" },
  "creator": { "id": "...", "name": "Bob Example" },
  "createdAt": "2025-01-01T00:00:00.000Z"
}
```

## Roadmap projects (`roadmap get projects`)

```json
{
  "id": "...",
  "slugId": "crm-actions-d0f9",
  "url": "https://linear.app/.../project/...",
  "name": "CRM Actions",
  "status": "started",
  "progress": 0.45,
  "lead": "Alice Example",
  "startDate": "2025-01-15",
  "targetDate": "2025-03-31"
}
```

## Issue relations (`issue relation list`)

```json
[
  { "id": "...", "type": "blocks", "relatedIssue": "ENG-124" },
  { "id": "...", "type": "blocked_by", "relatedIssue": "ENG-122" },
  { "id": "...", "type": "duplicate", "relatedIssue": "ENG-100" }
]
```

`blocked_by` is displayed for inverse "blocks" relations. Relation types: `blocks`, `duplicate`, `related`.

## Workflow states (`team states`)

```json
[
  { "id": "...", "name": "Todo", "type": "unstarted", "color": "#e2e2e2", "position": 0 },
  { "id": "...", "name": "In Progress", "type": "started", "color": "#5e6ad2", "position": 1 },
  { "id": "...", "name": "Done", "type": "completed", "color": "#5e6ad2", "position": 2 }
]
```

State types: `triage` | `backlog` | `unstarted` | `started` | `completed` | `canceled`

## Issue attachments (`issue attachment list`)

```json
[
  {
    "id": "...",
    "title": "PR #456",
    "url": "https://github.com/...",
    "subtitle": "Fixes login bug",
    "sourceType": "github"
  }
]
```

## Priority values

| Value    | Meaning         |
| -------- | --------------- |
| `none`   | No priority set |
| `urgent` | P0 — immediate  |
| `high`   | P1              |
| `medium` | P2              |
| `low`    | P3              |

## Team get output (`team get`)

Includes estimate configuration so callers can discover valid values before updating:

```json
{
  "id": "...",
  "name": "Engineering",
  "key": "ENG",
  "estimates": {
    "type": "fibonacci",
    "allowZero": false,
    "extended": false,
    "default": 0,
    "validValues": [1, 2, 3, 5, 8, 13],
    "display": "1 | 2 | 3 | 5 | 8 | 13"
  },
  "members": [...]
}
```

When `type` is `"notUsed"`, the `estimates` block is pruned.

## Estimate scales

| Type          | Base values                         | Extended adds  |
| ------------- | ----------------------------------- | -------------- |
| `fibonacci`   | 1, 2, 3, 5, 8, 13                   | 21, 34         |
| `linear`      | 1, 2, 3, 4, 5                       | 6, 7, 8, 9, 10 |
| `exponential` | 1, 2, 4, 8, 16                      | 32, 64         |
| `tShirt`      | 1 (XS), 2 (S), 3 (M), 4 (L), 5 (XL) | 6 (XXL)        |

When `allowZero` is true, `0` is also valid. T-shirt sizes are stored as integers.

## Project status values

`backlog` | `planned` | `started` | `paused` | `completed` | `canceled`
