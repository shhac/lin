/**
 * Sanitized example payloads documenting the shape of real Linear API responses
 * as returned by the lin CLI.
 *
 * All personally-identifiable information has been replaced:
 *   - Names   → Alice, Bob, Charlie, etc.
 *   - Emails  → alice@example.com, bob@example.com
 *   - IDs     → deterministic fake UUIDs (aaaaaaaa-…, bbbbbbbb-…)
 *   - URLs    → example.com equivalents
 *   - Text    → lorem-ipsum placeholders
 */

// ── auth status ──────────────────────────────────────────────────────

export const authStatusAuthenticated = {
  authenticated: true,
  source: "config",
  user: {
    id: "aaaaaaaa-1111-2222-3333-444444444444",
    name: "Alice Example",
    email: "alice@example.com",
  },
  organization: {
    id: "bbbbbbbb-1111-2222-3333-444444444444",
    name: "acme-corp",
    urlKey: "acme-corp",
  },
};

export const authStatusUnauthenticated = {
  authenticated: false,
};

// ── user me ──────────────────────────────────────────────────────────

export const userMe = {
  id: "aaaaaaaa-1111-2222-3333-444444444444",
  name: "Alice Example",
  email: "alice@example.com",
  displayName: "alice",
  organization: {
    id: "bbbbbbbb-1111-2222-3333-444444444444",
    name: "acme-corp",
  },
};

// ── team list ────────────────────────────────────────────────────────

export const teamList = [
  { id: "cccccccc-1111-2222-3333-444444444444", name: "Engineering", key: "ENG" },
  { id: "cccccccc-2222-3333-4444-555555555555", name: "Design", key: "DES" },
  { id: "cccccccc-3333-4444-5555-666666666666", name: "Operations", key: "OPS" },
];

// ── team get ─────────────────────────────────────────────────────────

export const teamGet = {
  id: "cccccccc-1111-2222-3333-444444444444",
  name: "Engineering",
  key: "ENG",
  members: [
    {
      id: "aaaaaaaa-1111-2222-3333-444444444444",
      name: "Alice Example",
      email: "alice@example.com",
    },
    {
      id: "dddddddd-1111-2222-3333-444444444444",
      name: "Bob Builder",
      email: "bob@example.com",
    },
    {
      id: "dddddddd-2222-3333-4444-555555555555",
      name: "Charlie Delta",
      email: "charlie@example.com",
    },
  ],
};

// ── user list ────────────────────────────────────────────────────────

export const userList = [
  {
    id: "aaaaaaaa-1111-2222-3333-444444444444",
    name: "Alice Example",
    email: "alice@example.com",
    displayName: "alice",
  },
  {
    id: "dddddddd-1111-2222-3333-444444444444",
    name: "Bob Builder",
    email: "bob@example.com",
    displayName: "bob",
  },
  {
    id: "dddddddd-2222-3333-4444-555555555555",
    name: "Charlie Delta",
    email: "charlie@example.com",
    displayName: "charlie",
  },
];

// ── label list ───────────────────────────────────────────────────────

export const labelList = [
  { id: "eeeeeeee-1111-2222-3333-444444444444", name: "bug", color: "#eb5757" },
  { id: "eeeeeeee-2222-3333-4444-555555555555", name: "feature", color: "#5795DF" },
  { id: "eeeeeeee-3333-4444-5555-666666666666", name: "chore", color: "#4cb782" },
];

// ── project list ─────────────────────────────────────────────────────

export const projectList = [
  {
    id: "ffffffff-1111-2222-3333-444444444444",
    slugId: "aaa111bbb222",
    url: "https://linear.app/example/project/project-alpha-aaa111bbb222",
    name: "Project Alpha",
    description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
    status: "started",
    progress: 0.42,
    startDate: "2026-01-15",
    targetDate: "2026-06-30",
  },
  {
    id: "ffffffff-2222-3333-4444-555555555555",
    slugId: "ccc333ddd444",
    url: "https://linear.app/example/project/project-beta-ccc333ddd444",
    name: "Project Beta",
    description: "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
    status: "planned",
    progress: 0,
    targetDate: "2026-09-01",
  },
];

// ── project search ───────────────────────────────────────────────────

export const projectSearch = [
  {
    id: "ffffffff-1111-2222-3333-444444444444",
    slugId: "aaa111bbb222",
    url: "https://linear.app/example/project/project-alpha-aaa111bbb222",
    name: "Project Alpha",
    description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
    status: "started",
    progress: 0.42,
  },
];

// ── project get overview ─────────────────────────────────────────────

export const projectOverview = {
  id: "ffffffff-1111-2222-3333-444444444444",
  slugId: "aaa111bbb222",
  url: "https://linear.app/example/project/project-alpha-aaa111bbb222",
  name: "Project Alpha",
  description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
  status: "started",
  progress: 0.42,
  lead: {
    id: "aaaaaaaa-1111-2222-3333-444444444444",
    name: "Alice Example",
  },
  startDate: "2026-01-15",
  targetDate: "2026-06-30",
  milestones: [
    {
      id: "11111111-aaaa-bbbb-cccc-dddddddddddd",
      name: "Milestone 1 - MVP",
      targetDate: "2026-03-15",
    },
    {
      id: "11111111-bbbb-cccc-dddd-eeeeeeeeeeee",
      name: "Milestone 2 - Beta",
      targetDate: "2026-05-01",
    },
  ],
};

// ── project get issues ───────────────────────────────────────────────

export const projectIssues = [
  {
    id: "22222222-aaaa-bbbb-cccc-dddddddddddd",
    identifier: "ENG-101",
    title: "Implement authentication flow",
    priority: 2,
    priorityLabel: "High",
  },
  {
    id: "22222222-bbbb-cccc-dddd-eeeeeeeeeeee",
    identifier: "ENG-102",
    title: "Set up CI/CD pipeline",
    priority: 3,
    priorityLabel: "Medium",
  },
];

// ── issue list ───────────────────────────────────────────────────────

export const issueList = [
  {
    id: "22222222-aaaa-bbbb-cccc-dddddddddddd",
    identifier: "ENG-101",
    title: "Implement authentication flow",
    priority: 2,
    priorityLabel: "High",
  },
  {
    id: "22222222-bbbb-cccc-dddd-eeeeeeeeeeee",
    identifier: "ENG-102",
    title: "Set up CI/CD pipeline",
    priority: 3,
    priorityLabel: "Medium",
  },
  {
    id: "22222222-cccc-dddd-eeee-ffffffffffff",
    identifier: "ENG-103",
    title: "Fix login redirect bug",
    priority: 1,
    priorityLabel: "Urgent",
  },
];

// ── issue search ─────────────────────────────────────────────────────

export const issueSearch = [
  {
    id: "22222222-cccc-dddd-eeee-ffffffffffff",
    identifier: "ENG-103",
    title: "Fix login redirect bug",
    priority: 1,
    priorityLabel: "Urgent",
  },
];

// ── issue get overview ───────────────────────────────────────────────

export const issueOverviewMinimal = {
  id: "22222222-aaaa-bbbb-cccc-dddddddddddd",
  identifier: "ENG-101",
  title: "Implement authentication flow",
  description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
  status: {
    id: "33333333-aaaa-bbbb-cccc-dddddddddddd",
    name: "Backlog",
    type: "backlog",
  },
  priority: 0,
  priorityLabel: "No priority",
};

export const issueOverviewFull = {
  id: "22222222-cccc-dddd-eeee-ffffffffffff",
  identifier: "ENG-103",
  title: "Fix login redirect bug",
  description:
    "Users are being redirected to the wrong page after logging in.\n\nSteps to reproduce:\n1. Log in with valid credentials\n2. Observe redirect goes to /dashboard instead of /home",
  status: {
    id: "33333333-bbbb-cccc-dddd-eeeeeeeeeeee",
    name: "Triage",
    type: "triage",
  },
  assignee: {
    id: "dddddddd-1111-2222-3333-444444444444",
    name: "Bob Builder",
  },
  priority: 1,
  priorityLabel: "Urgent",
  labels: [
    { id: "eeeeeeee-1111-2222-3333-444444444444", name: "bug" },
    { id: "eeeeeeee-4444-5555-6666-777777777777", name: "auth" },
  ],
  parent: {
    id: "22222222-aaaa-bbbb-cccc-dddddddddddd",
    identifier: "ENG-101",
  },
  estimate: 3,
  dueDate: "2026-03-01",
  createdAt: "2026-02-01T10:00:00.000Z",
  updatedAt: "2026-02-10T14:30:00.000Z",
};

// ── issue get comments ───────────────────────────────────────────────

export const issueComments = [
  {
    id: "44444444-aaaa-bbbb-cccc-dddddddddddd",
    body: "I can reproduce this on staging. The redirect URL is hardcoded in the auth callback handler.",
    createdAt: "2026-02-02T09:15:00.000Z",
    updatedAt: "2026-02-02T09:15:00.000Z",
  },
  {
    id: "44444444-bbbb-cccc-dddd-eeeeeeeeeeee",
    body: "Fixed in PR #42. The callback now reads the redirect from the session.",
    createdAt: "2026-02-03T16:45:00.000Z",
    updatedAt: "2026-02-03T16:45:00.000Z",
  },
];

// ── issue new (create response) ──────────────────────────────────────

export const issueCreated = {
  id: "22222222-dddd-eeee-ffff-000000000000",
  identifier: "ENG-104",
  title: "Add rate limiting to API endpoints",
  created: true,
};

// ── issue update (generic) ───────────────────────────────────────────

export const issueUpdated = {
  updated: true,
};

// ── issue comment new ────────────────────────────────────────────────

export const commentCreated = {
  id: "44444444-cccc-dddd-eeee-ffffffffffff",
  body: "This looks good to merge after CI passes.",
  created: true,
};

// ── cycle list ───────────────────────────────────────────────────────

export const cycleList = [
  {
    id: "55555555-aaaa-bbbb-cccc-dddddddddddd",
    number: 6,
    name: "Sprint 6",
    startsAt: "2026-02-03T00:00:00.000Z",
    endsAt: "2026-02-17T00:00:00.000Z",
  },
  {
    id: "55555555-bbbb-cccc-dddd-eeeeeeeeeeee",
    number: 5,
    name: "Sprint 5",
    startsAt: "2026-01-20T00:00:00.000Z",
    endsAt: "2026-02-03T00:00:00.000Z",
  },
];

// ── cycle get ────────────────────────────────────────────────────────

export const cycleGet = {
  id: "55555555-aaaa-bbbb-cccc-dddddddddddd",
  number: 6,
  name: "Sprint 6",
  startsAt: "2026-02-03T00:00:00.000Z",
  endsAt: "2026-02-17T00:00:00.000Z",
  issues: [
    {
      id: "22222222-aaaa-bbbb-cccc-dddddddddddd",
      identifier: "ENG-101",
      title: "Implement authentication flow",
      priority: 2,
      priorityLabel: "High",
    },
    {
      id: "22222222-cccc-dddd-eeee-ffffffffffff",
      identifier: "ENG-103",
      title: "Fix login redirect bug",
      priority: 1,
      priorityLabel: "Urgent",
    },
  ],
};

// ── project update ───────────────────────────────────────────────────

export const projectTitleUpdated = {
  id: "ffffffff-1111-2222-3333-444444444444",
  name: "Project Alpha (Renamed)",
  updated: true,
};

export const projectUpdated = {
  updated: true,
};

// ── error response ───────────────────────────────────────────────────

export const errorResponse = {
  error: "Not authenticated. Run: lin auth login <api-key>",
};

export const errorNotFound = {
  error: "Entity not found",
};
