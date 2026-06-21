# Customer Requests in `lin`

Status: implementing · Owner: `lin` maintainers

## Why

Linear's **Customer Requests** feature is how product feedback from real
customers ("Acme wants SSO") gets attached to the issues and projects that would
deliver it. PMs and support engineers live in questions like:

- *What customer requests came in recently?*
- *Which requests are unassigned / still in triage?*
- *Which requests are flagged important?*
- *Which issues are most demanded by customers — and by important customers?*

Today `lin` has no way to answer these short of hand-writing GraphQL through
`lin api query`. This doc defines a first-class `customer` domain plus customer
visibility on the existing `issue`/`project` commands, so both humans and LLM
agents can route and triage customer feedback with the same ergonomics as every
other entity.

## The Linear model (what we're wrapping)

A **customer request** is `CustomerNeed` in the GraphQL API ("need" internally,
"request" in the product UI — we use **request** in the CLI surface). A request
is a *join entity*:

```
   Customer ◄─────── CustomerNeed ───────► Issue     (issue-based request)
  (optional)         (the request)  ──────► Project   (project-based request)
                                     one OR the other, never both
```

Key facts that shape the design:

- A request links to **exactly one** of an issue or a project (`issueId` *or*
  `projectId`). Issue-based requests have their `project` denormalized from the
  issue's project.
- A request does **not** have its own workflow status — its "state" is the
  linked issue's workflow state. So *"in triage"* is an issue-state question
  (`issue.state.type == "triage"`), not a request field.
- `CustomerNeed.priority` is the **importance flag**: `0 = not important`,
  `1 = important`. (Not the same scale as issue priority.)
- Requests have **no labels of their own.** Categorization comes from three
  places: the linked issue's labels, the customer's **tier** (Enterprise/Pro/…)
  and the customer's **status** (Active/Churned/…).
- One issue can have **many** requests (`Issue.needs` is a connection). That is
  why issues expose `customerCount` and `customerImportantCount` scalars — cheap
  aggregates over that one-to-many.

Supporting entities: `Customer` (the org), `CustomerStatus` (lifecycle, e.g.
Active/Churned), `CustomerTier` (segment, e.g. Enterprise/Pro/Free).

## Question → filter mapping

`CustomerNeedFilter` nests a full `NullableIssueFilter`, so most questions reuse
the issue-filter machinery we already have:

| Question | Filter |
|---|---|
| New requests | `customerNeeds(orderBy: createdAt)` + `createdAt: { gte }` |
| Unassigned | `issue: { assignee: { null: true } }` |
| In triage | `issue: { state: { type: { eq: "triage" } } }` |
| Important | `priority: { eq: 1 }` |
| By issue label / team / status | `issue: { labels / team / state }` |
| For a customer | `customer: { id / name }` |
| Most-requested issues | `issues(sort: { customerCount })` / `customerImportantCount` |

## Command interface

Conventions inherited from the rest of the CLI (do not reinvent):

- **`search <text>`** = full-text (takes a text arg, lighter filters);
  **`list`** = filter-driven (no arg, rich filters). Customers only support a
  name comparator, so `customer list` is primary and `search` maps name-contains.
- **Pagination** is `--limit` / `--cursor` via `output.AddPageFlags` +
  `output.PrintPage`. **ndjson/yaml** come free from the global `--format` flag —
  `PrintPaginated`/`PrintPage` already branch on it. No per-command work.
- **Detail vs collection:** a `get` returns the entity, lightweight inline
  relations, and **counts** of heavy collections — never the collection itself.
  Heavy one-to-many collections get their own paginated subcommand (mirrors
  `issue get` showing `commentCount` while `issue comment list` returns them).
  `--full`/`--expand` are *only* for un-truncating long text — never a toggle for
  including collections.

```
lin customer
├── list                      # filter: --tier --status --owner --domain --revenue; paginated, sortable
├── search <text>             # Customer.name contains; paginated
├── get <id|slug>             # customer detail + approximateNeedCount (NOT the needs list)
├── requests                  # list CustomerNeeds — the workhorse; paginated
│     --customer <id|name>          # scope to one customer
│     --important                    # priority = 1
│     --unassigned                   # linked issue has no assignee
│     --triage                       # linked issue state type = triage
│     --status <name> --label <l> --team <t>   # reuse issue-filter flags
│     --project <id>                 # project-linked requests
│     --created-after / --created-before <YYYY-MM-DD>
├── statuses                  # list workspace CustomerStatuses (lifecycle)
└── tiers                     # list workspace CustomerTiers (segments)

# Entity-scoped collections, mirroring `lin project issues`:
lin issue requests <id>       # CustomerNeeds linked to an issue (Issue.needs)
lin project requests <id>     # CustomerNeeds linked to a project (Project.needs)

# Inline counts added to existing detail/summary output (zero extra query cost):
lin issue get / search        # + customerCount, customerImportantCount
```

## Output shapes

Summary/detail split, pruned of empty fields like every other mapper.

- **`MapCustomerSummary`** — `id, name, slugId, tier, status, owner, revenue,
  approximateNeedCount, url`
- **`MapCustomerDetail`** — summary + `domains, externalIds, size, createdAt,
  updatedAt`
- **`MapCustomerNeedSummary`** — `id, important (bool from priority), customer
  {id,name}, issue {identifier,title} OR project {id,name}, url, createdAt`. The
  mapper shows whichever of issue/project is set (project may be present-but-
  derived on issue-based requests).

## Resolvers

`ResolveCustomer(client, input)` accepts a UUID, slug, or exact name — same
shape as `ResolveProject`: try `customer(id:)` first (it accepts id-or-slug),
fall back to a `customers(filter:)` name lookup. This makes `--customer acme`
and `customer get acme-corp` both work.

## Rollout

1. **Step zero** — `customer.graphql` operations + `make generate` + fixtures.
2. **Issue/project visibility** — counts on issue output, `issue requests` /
   `project requests`. Independently useful, ships first.
3. **Customer domain** — commands, resolvers, filters, mappers, wired into root.

## Out of scope (follow-ups)

- **Mutations.** The read surface answers every question above. The
  highest-value write is `customerNeedCreate` (attach a request to an issue,
  routing feedback into Linear) — a clean follow-up once the read surface lands.
  `customerCreate/Update`, status/tier management follow the same pattern.
- **Priority sort on requests.** The `customerNeeds` query only orders by
  `createdAt`/`updatedAt`; "important first" is a client-side sort if wanted.
