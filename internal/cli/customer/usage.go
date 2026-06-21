package customer

const usageText = `lin customer — Customer and customer-request operations

A customer request (Linear "customer need") links a customer to an issue or
project. Requests have no status of their own — their state is the linked
issue's state, so "in triage" / "unassigned" are issue-state filters.

LIST & SEARCH:
  customer list     [--tier <name>] [--status <name>] [--owner <me|name|email|UUID>]
                    [--domain <domain>] [--revenue <min>] [--limit] [--cursor]
  customer search <text>     Name substring match [--limit] [--cursor]

GET:
  customer get <id|slug>...  Detail: tier, status, owner, domains, externalIds,
                             revenue, size, approximateNeedCount
                             NDJSON by default (one line per id; --format json for a single object).
                             Missing ids emit {"@unresolved":{...}} on stdout (exit 0).
                             --format pretty renders a human-readable card (--width <n>).

CUSTOMER REQUESTS:
  customer requests [filters] [--limit] [--cursor]
    --customer <id|slug|name>   Scope to one customer
    --project <id|slug|name>    Scope to one project
    --important                 Only important requests (priority = 1)
    --unassigned                Linked issue has no assignee
    --triage                    Linked issue is in a triage state
    --status <name>             Linked issue status name
    --label <name>              Linked issue label
    --team <key|name>           Linked issue team
    --created-after|before <YYYY-MM-DD>

REFERENCE LISTS:
  customer statuses          Workspace lifecycle statuses (e.g. Active, Churned)
  customer tiers             Workspace tiers/segments (e.g. Enterprise, Pro, Free)

RELATED:
  issue requests <id>        Customer requests linked to a specific issue
  project requests <id>      Customer requests linked to a specific project
  issue get <id>             Includes customerRequestCount, customerImportantCount

IDS: customers accept UUID, slug, or exact name.
PAGINATION: --limit <n> --cursor <token> on all list commands.`
