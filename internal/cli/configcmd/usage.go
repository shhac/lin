package configcmd

const usageText = `lin config — View and update persistent CLI settings

SUBCOMMANDS:
  config get [key...]         Show current settings as NDJSON (omit key for all)
  config set <key> <value>    Update a setting
  config reset [key]          Reset to defaults (omit key to reset all)
  config list-keys            List all available setting keys with descriptions

SETTING KEYS:
  truncation.maxLength         Max chars before truncating description/body/content fields
                               Default: 200. Must be a non-negative integer (0 = no truncation).
  pagination.defaultPageSize   Default number of results for list/search commands
                               Default: 50. Must be an integer between 1 and 250.
  output.defaultFormat          Default output format when --format is omitted
                               Values: json, yaml, jsonl.
  request.timeoutMS             Default request timeout in milliseconds
                               Default: 0 (no client timeout). Must be non-negative.

EXAMPLES:
  config set truncation.maxLength 500       Show more content before truncating
  config set pagination.defaultPageSize 20  Fetch fewer results per page
  config set output.defaultFormat json       Return JSON envelopes by default
  config set request.timeoutMS 10000         Set a 10 second API timeout
  config get truncation.maxLength           Check current truncation setting
  config reset truncation.maxLength         Reset truncation to default (200)
  config reset                              Reset all settings to defaults

STORAGE: Settings persisted in ~/.config/lin/config.json alongside auth credentials.

OUTPUT: NDJSON to stdout by default — config get with no args emits one line per
setting; config get <key>... emits one {"key","value"} line per key (or an
{"@unresolved"} line for an unknown key). Pass --format json for a {"data":[...]}
envelope. config set/reset still report status as before.`
