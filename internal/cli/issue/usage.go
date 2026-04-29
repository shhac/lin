package issue

const usageText = `lin issue — Issue operations (search, list, create, update, comment, relate, archive, attach, history)

SEARCH & LIST:
  issue search <text> [--project] [--team] [--assignee] [--status] [--priority]
  issue list [--project] [--team] [--assignee] [--status] [--priority] [--label] [--cycle]
    [--updated-after|before] [--created-after|before] (YYYY-MM-DD)

GET:
  issue get <id>             Full details + labels, attachments, branchName, commentCount

CREATE:
  issue new <title> --team <key|name|UUID>    --team required
    [--project <name|slug|UUID>] [--assignee] [--priority] [--status]
    [--labels <names|ids>] [--description <md>] [--cycle] [--parent <id>] [--estimate <n>]

UPDATE (each is a subcommand):
  issue update title|status|assignee|priority|project|labels|description|estimate <id> <value>
  issue update due-date <id> <YYYY-MM-DD>
  issue update cycle <id> <cycle-id>  |  parent <id> <parent-id>
  status: team-scoped (use "team states <team>"). labels: comma-separated names or IDs.
  estimate: validated against team scale. assignee: name/email/ID.

COMMENTS (--parent 1 level):
  issue comment list <id> [--limit] [--cursor]
  issue comment new <id> <body> [--parent <cid>]
  issue comment edit <cid> <body>
  issue comment get <cid>  |  replies <cid>

RELATIONS:
  issue relation list <id>      Both directions (blocks, blocked_by, duplicate, related)
  issue relation add <id> --type <t> --related <id>  Types: blocks|duplicate|related
  issue relation remove <id>

HISTORY:
  issue history <id> [--limit] [--cursor]     Activity log (status, assignee, priority, labels, etc.)

LIFECYCLE:  issue archive|unarchive|delete <id>

ATTACHMENTS:
  issue attachment list <id>            All attachments (any source type)
  issue attachment remove <attachment-id>   Works for any attachment (URL, GitHub PR, Slack, …)
  issue attachment add <id> <url> [--title <t>]    Default: rich link via attachmentLinkURL
    --github-pr        Force GitHub pull request integration
    --github-issue     Force GitHub issue integration
    --gitlab-mr        Force GitLab merge request integration (project + number derived from URL)
    --slack            Force Slack message integration
      --sync-thread    (with --slack) sync the Slack thread to a comment thread
    --discord          Force Discord message integration (channel + message derived from URL)

IDS: Issue keys (ENG-123) or UUIDs. PRIORITY: none|urgent|high|medium|low
PAGINATION: --limit <n> --cursor <token> on search and list.`
