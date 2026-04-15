import { describe, expect, test } from "bun:test";
import * as fixtures from "./fixtures.ts";

/**
 * Validates that all fixture payloads conform to the expected structural shapes.
 * These tests ensure the fixtures themselves are well-formed and serve as
 * living documentation of the CLI's JSON output contracts.
 */

// ── Helpers ──────────────────────────────────────────────────────────

function expectUuid(value: unknown): void {
  expect(typeof value).toBe("string");
  expect(value as string).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/);
}

function expectNonEmptyString(value: unknown): void {
  expect(typeof value).toBe("string");
  expect((value as string).length).toBeGreaterThan(0);
}

// ── Auth ─────────────────────────────────────────────────────────────

describe("fixture: auth", () => {
  test("authStatusAuthenticated has correct shape", () => {
    const d = fixtures.authStatusAuthenticated;
    expect(d.authenticated).toBe(true);
    expect(["config", "environment"]).toContain(d.source);
    expectUuid(d.user.id);
    expectNonEmptyString(d.user.name);
    expectNonEmptyString(d.user.email);
    expectUuid(d.organization.id);
    expectNonEmptyString(d.organization.name);
    expectNonEmptyString(d.organization.urlKey);
  });

  test("authStatusUnauthenticated has correct shape", () => {
    expect(fixtures.authStatusUnauthenticated).toEqual({ authenticated: false });
  });
});

// ── User ─────────────────────────────────────────────────────────────

describe("fixture: user", () => {
  test("userMe has correct shape", () => {
    const d = fixtures.userMe;
    expectUuid(d.id);
    expectNonEmptyString(d.name);
    expectNonEmptyString(d.email);
    expectNonEmptyString(d.displayName);
    expectUuid(d.organization.id);
    expectNonEmptyString(d.organization.name);
  });

  test("userList items have correct shape", () => {
    expect(fixtures.userList.length).toBeGreaterThan(0);
    for (const u of fixtures.userList) {
      expectUuid(u.id);
      expectNonEmptyString(u.name);
      expectNonEmptyString(u.email);
      expectNonEmptyString(u.displayName);
    }
  });
});

// ── Team ─────────────────────────────────────────────────────────────

describe("fixture: team", () => {
  test("teamList items have correct shape", () => {
    expect(fixtures.teamList.length).toBeGreaterThan(0);
    for (const t of fixtures.teamList) {
      expectUuid(t.id);
      expectNonEmptyString(t.name);
      expectNonEmptyString(t.key);
    }
  });

  test("teamGet has members array", () => {
    const d = fixtures.teamGet;
    expectUuid(d.id);
    expectNonEmptyString(d.name);
    expectNonEmptyString(d.key);
    expect(d.members.length).toBeGreaterThan(0);
    for (const m of d.members) {
      expectUuid(m.id);
      expectNonEmptyString(m.name);
      expectNonEmptyString(m.email);
    }
  });
});

// ── Label ────────────────────────────────────────────────────────────

describe("fixture: label", () => {
  test("labelList items have id, name, color", () => {
    expect(fixtures.labelList.length).toBeGreaterThan(0);
    for (const l of fixtures.labelList) {
      expectUuid(l.id);
      expectNonEmptyString(l.name);
      expect(l.color).toMatch(/^#[0-9a-fA-F]{6}$/);
    }
  });
});

// ── Project ──────────────────────────────────────────────────────────

describe("fixture: project", () => {
  test("projectList items have correct shape", () => {
    for (const p of fixtures.projectList) {
      expectUuid(p.id);
      expectNonEmptyString(p.slugId);
      expect(p.url).toContain("http");
      expectNonEmptyString(p.name);
      expectNonEmptyString(p.status);
      expect(typeof p.progress).toBe("number");
    }
  });

  test("projectSearch returns same shape as projectList", () => {
    for (const p of fixtures.projectSearch) {
      expectUuid(p.id);
      expectNonEmptyString(p.name);
      expectNonEmptyString(p.status);
      expect(typeof p.progress).toBe("number");
    }
  });

  test("projectOverview includes lead and milestones", () => {
    const d = fixtures.projectOverview;
    expectUuid(d.id);
    expectNonEmptyString(d.name);
    expectNonEmptyString(d.status);
    expect(typeof d.progress).toBe("number");
    expect(d.lead).toBeDefined();
    expectUuid(d.lead!.id);
    expect(d.milestones.length).toBeGreaterThan(0);
    for (const m of d.milestones) {
      expectUuid(m.id);
      expectNonEmptyString(m.name);
    }
  });

  test("projectIssues items have issue shape", () => {
    for (const i of fixtures.projectIssues) {
      expectUuid(i.id);
      expect(i.identifier).toMatch(/^[A-Z]+-\d+$/);
      expectNonEmptyString(i.title);
      expect(typeof i.priority).toBe("number");
      expectNonEmptyString(i.priorityLabel);
    }
  });

  test("projectTitleUpdated includes updated flag", () => {
    const d = fixtures.projectTitleUpdated;
    expectUuid(d.id);
    expectNonEmptyString(d.name);
    expect(d.updated).toBe(true);
  });

  test("projectUpdated is simple success flag", () => {
    expect(fixtures.projectUpdated).toEqual({ updated: true });
  });
});

// ── Issue ────────────────────────────────────────────────────────────

describe("fixture: issue", () => {
  test("issueList items have correct shape", () => {
    expect(fixtures.issueList.length).toBeGreaterThan(0);
    for (const i of fixtures.issueList) {
      expectUuid(i.id);
      expect(i.identifier).toMatch(/^[A-Z]+-\d+$/);
      expectNonEmptyString(i.title);
      expect(typeof i.priority).toBe("number");
      expect(i.priority).toBeGreaterThanOrEqual(0);
      expect(i.priority).toBeLessThanOrEqual(4);
      expectNonEmptyString(i.priorityLabel);
    }
  });

  test("issueSearch has same shape as issueList items", () => {
    for (const i of fixtures.issueSearch) {
      expectUuid(i.id);
      expect(i.identifier).toMatch(/^[A-Z]+-\d+$/);
      expectNonEmptyString(i.title);
      expect(typeof i.priority).toBe("number");
    }
  });

  test("issueOverviewMinimal has status but no assignee/labels/parent", () => {
    const d = fixtures.issueOverviewMinimal;
    expectUuid(d.id);
    expect(d.identifier).toMatch(/^[A-Z]+-\d+$/);
    expectNonEmptyString(d.title);
    expect(d.status).toBeDefined();
    expectUuid(d.status.id);
    expectNonEmptyString(d.status.name);
    expectNonEmptyString(d.status.type);
    // Minimal: no assignee, labels, parent, estimate, dueDate
    expect(d).not.toHaveProperty("assignee");
    expect(d).not.toHaveProperty("labels");
    expect(d).not.toHaveProperty("parent");
  });

  test("issueOverviewFull has assignee, labels, parent, dates", () => {
    const d = fixtures.issueOverviewFull;
    expectUuid(d.id);
    expect(d.identifier).toMatch(/^[A-Z]+-\d+$/);
    expectNonEmptyString(d.title);
    expectNonEmptyString(d.description);

    // Status
    expectUuid(d.status.id);
    expectNonEmptyString(d.status.name);
    expectNonEmptyString(d.status.type);

    // Assignee
    expect(d.assignee).toBeDefined();
    expectUuid(d.assignee!.id);
    expectNonEmptyString(d.assignee!.name);

    // Labels
    expect(d.labels!.length).toBeGreaterThan(0);
    for (const l of d.labels!) {
      expectUuid(l.id);
      expectNonEmptyString(l.name);
    }

    // Parent
    expect(d.parent).toBeDefined();
    expectUuid(d.parent!.id);
    expect(d.parent!.identifier).toMatch(/^[A-Z]+-\d+$/);

    // Dates
    expect(typeof d.estimate).toBe("number");
    expectNonEmptyString(d.dueDate);
    expectNonEmptyString(d.createdAt);
    expectNonEmptyString(d.updatedAt);
  });

  test("issueComments have body and timestamps", () => {
    expect(fixtures.issueComments.length).toBeGreaterThan(0);
    for (const c of fixtures.issueComments) {
      expectUuid(c.id);
      expectNonEmptyString(c.body);
      expectNonEmptyString(c.createdAt);
      expectNonEmptyString(c.updatedAt);
    }
  });

  test("issueCreated has success shape", () => {
    const d = fixtures.issueCreated;
    expectUuid(d.id);
    expect(d.identifier).toMatch(/^[A-Z]+-\d+$/);
    expectNonEmptyString(d.title);
    expect(d.created).toBe(true);
  });

  test("issueUpdated is simple success flag", () => {
    expect(fixtures.issueUpdated).toEqual({ updated: true });
  });

  test("commentCreated has success shape", () => {
    const d = fixtures.commentCreated;
    expectUuid(d.id);
    expectNonEmptyString(d.body);
    expect(d.created).toBe(true);
  });
});

// ── Cycle ────────────────────────────────────────────────────────────

describe("fixture: cycle", () => {
  test("cycleList items have correct shape", () => {
    expect(fixtures.cycleList.length).toBeGreaterThan(0);
    for (const c of fixtures.cycleList) {
      expectUuid(c.id);
      expect(typeof c.number).toBe("number");
    }
  });

  test("cycleGet includes issues array", () => {
    const d = fixtures.cycleGet;
    expectUuid(d.id);
    expect(typeof d.number).toBe("number");
    expect(d.issues.length).toBeGreaterThan(0);
    for (const i of d.issues) {
      expectUuid(i.id);
      expect(i.identifier).toMatch(/^[A-Z]+-\d+$/);
    }
  });
});

// ── Error ────────────────────────────────────────────────────────────

describe("fixture: errors", () => {
  test("errorResponse has error string", () => {
    expectNonEmptyString(fixtures.errorResponse.error);
  });

  test("errorNotFound has error string", () => {
    expectNonEmptyString(fixtures.errorNotFound.error);
  });
});
