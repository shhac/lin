package initiative

const usageText = `lin initiative — Manage Linear initiatives (replaces roadmap)

SEARCH & LIST:
  initiative search <text>                Search initiatives by name
    [--limit] [--cursor]
  initiative list                         List all initiatives
    [--status planned|active|completed]
    [--limit] [--cursor]
    Returns per item: id, slugId, url, name, status, health, targetDate, owner

GET:
  initiative get <id>            Initiative summary: id, slugId, url, name,
                                 description, content, status, health, targetDate,
                                 startedAt, completedAt, createdAt, updatedAt,
                                 owner{id,name}, creator{id,name}
  initiative projects <id>       Projects linked to an initiative
    [--limit] [--cursor]
    Returns per item: id, slugId, url, name, status, progress,
    lead, startDate, targetDate

CREATE:
  initiative new <name>          Create initiative
    [--description] [--owner] [--status] [--target-date]
    [--content] [--color] [--icon]

UPDATE:
  initiative update name <id> <value>
  initiative update description <id> <value>
  initiative update content <id> <value>
  initiative update status <id> <value>        planned|active|completed
  initiative update owner <id> <user>          name, email, or user ID
  initiative update target-date <id> <date>    YYYY-MM-DD
  initiative update color <id> <color>
  initiative update icon <id> <icon>

LIFECYCLE:
  initiative archive <id>        Archive initiative
  initiative unarchive <id>      Restore archived initiative
  initiative delete <id>         Permanently delete initiative

IDS: <id> accepts UUID, slug ID, or initiative name.
PAGINATION: --limit <n> --cursor <token> on list and projects.
STATUS: planned | active | completed`
