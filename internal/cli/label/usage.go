package label

const usageText = `lin label — Search, list, and inspect Linear labels (issue or project)

SUBCOMMANDS:
  label list [--type <t>] [--team <team>] [--name <text>] [--is-group[=false]]   List labels (filterable)
  label search <text> [--type <t>] [--team <team>]                               Substring search (case- and accent-insensitive)
  label get <id|name> [--type <t>] [--team <team>]                               Single label by UUID or exact name

OPTIONS:
  --type <t>        Label type: "issue" (default) or "project"
  --team <team>     Filter/disambiguate by team key, name, or UUID (issue labels only)
  --name <text>     Exact match (case-insensitive)
  --is-group        Only group labels (--is-group=false for non-groups)
  --limit <n>       Limit results (list, search)
  --cursor <token>  Pagination cursor (list, search)

OUTPUT FIELDS:
  list/search (issue)   → id, name, color, [description, isGroup, team{id,key,name}, parent{id,name}]
  list/search (project) → id, name, color, [description, isGroup, parent{id,name}]
  get                   → same fields, single object

NOTES:
  Issue labels can be workspace-wide or team-scoped; --team filters/disambiguates.
  Project labels are workspace-only — --team is rejected with --type=project.
  IssueLabel and ProjectLabel are distinct Linear entities; the same name can exist in both.
  Use the resulting issue label name (with --team) or UUID with "issue new --labels" / "issue update labels".
  Use a project label name or UUID with "project update labels".`
