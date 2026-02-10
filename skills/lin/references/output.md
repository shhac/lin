# Output format (reference)

## General

All commands print JSON to stdout. Errors print `{ "error": "..." }` to stderr with non-zero exit.

Empty/null fields are pruned automatically — missing keys mean no value, not `null`.

Error messages include valid values when input is invalid (e.g., unknown status names list all valid states).

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

Use `project get details <id>` for the full markdown body.

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
  { "id": "...", "title": "PR #456", "url": "https://github.com/...", "subtitle": "Fixes login bug", "sourceType": "github" }
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
