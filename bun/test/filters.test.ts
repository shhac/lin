import { describe, expect, test } from "bun:test";
import {
  buildIssueFilter,
  buildProjectFilter,
  buildTeamFilter,
  nonEmptyFilter,
} from "../src/lib/filters.ts";

describe("nonEmptyFilter", () => {
  test("returns undefined for empty object", () => {
    expect(nonEmptyFilter({})).toBeUndefined();
  });

  test("returns the object when non-empty", () => {
    const result = nonEmptyFilter({ status: "active" });
    expect(result).toBeDefined();
    expect(result).toHaveProperty("status", "active");
  });
});

describe("buildTeamFilter", () => {
  test("matches by key or name", () => {
    expect(buildTeamFilter("ENG")).toEqual({
      or: [{ key: { eqIgnoreCase: "ENG" } }, { name: { eqIgnoreCase: "ENG" } }],
    });
  });
});

describe("buildProjectFilter", () => {
  test("matches by slugId or name when input is not a UUID", () => {
    const result = buildProjectFilter("my-project");
    expect(result).toEqual({
      or: [{ slugId: { eq: "my-project" } }, { name: { eqIgnoreCase: "my-project" } }],
    });
  });

  test("includes id branch when input is a UUID", () => {
    const result = buildProjectFilter("a1b2c3d4-e5f6-7890-abcd-ef1234567890");
    expect(result).toEqual({
      or: [
        { id: { eq: "a1b2c3d4-e5f6-7890-abcd-ef1234567890" } },
        { slugId: { eq: "a1b2c3d4-e5f6-7890-abcd-ef1234567890" } },
        { name: { eqIgnoreCase: "a1b2c3d4-e5f6-7890-abcd-ef1234567890" } },
      ],
    });
  });
});

describe("buildIssueFilter", () => {
  test("returns empty object for no options", () => {
    expect(buildIssueFilter({})).toEqual({});
  });

  test("builds project filter", () => {
    const result = buildIssueFilter({ project: "abc-123" });
    expect(result.project).toEqual(buildProjectFilter("abc-123"));
  });

  test("builds team filter", () => {
    const result = buildIssueFilter({ team: "ENG" });
    expect(result.team).toEqual(buildTeamFilter("ENG"));
  });

  test("handles assignee 'me' special case", () => {
    const result = buildIssueFilter({ assignee: "me" });
    expect(result.assignee).toEqual({ isMe: { eq: true } });
  });

  test("handles assignee 'me' case-insensitively", () => {
    const result = buildIssueFilter({ assignee: "ME" });
    expect(result.assignee).toEqual({ isMe: { eq: true } });
  });

  test("handles assignee by name/email (no id branch for non-UUID)", () => {
    const result = buildIssueFilter({ assignee: "alice@example.com" });
    expect(result.assignee).toEqual({
      or: [
        { name: { eqIgnoreCase: "alice@example.com" } },
        { displayName: { eqIgnoreCase: "alice@example.com" } },
        { email: { eqIgnoreCase: "alice@example.com" } },
      ],
    });
  });

  test("includes id branch for assignee when input is a UUID", () => {
    const uuid = "a1b2c3d4-e5f6-7890-abcd-ef1234567890";
    const result = buildIssueFilter({ assignee: uuid });
    expect(result.assignee).toEqual({
      or: [
        { id: { eq: uuid } },
        { name: { eqIgnoreCase: uuid } },
        { displayName: { eqIgnoreCase: uuid } },
        { email: { eqIgnoreCase: uuid } },
      ],
    });
  });

  test("builds status filter", () => {
    const result = buildIssueFilter({ status: "In Progress" });
    expect(result.state).toEqual({ name: { eqIgnoreCase: "In Progress" } });
  });

  test("builds priority filter for valid priority", () => {
    const result = buildIssueFilter({ priority: "high" });
    expect(result.priority).toEqual({ eq: 2 });
  });

  test("ignores invalid priority silently", () => {
    const result = buildIssueFilter({ priority: "super-duper" });
    expect(result.priority).toBeUndefined();
  });

  test("builds label filter", () => {
    const result = buildIssueFilter({ label: "Bug" });
    expect(result.labels).toEqual({ name: { eqIgnoreCase: "Bug" } });
  });

  test("builds cycle filter", () => {
    const result = buildIssueFilter({ cycle: "cycle-uuid" });
    expect(result.cycle).toEqual({ id: { eq: "cycle-uuid" } });
  });

  test("builds updatedAt range with both bounds", () => {
    const result = buildIssueFilter({
      "updated-after": "2024-01-01",
      "updated-before": "2024-01-31",
    });
    expect(result.updatedAt).toEqual({ gte: "2024-01-01", lte: "2024-01-31" });
  });

  test("builds updatedAt with only after", () => {
    const result = buildIssueFilter({ "updated-after": "2024-01-01" });
    expect(result.updatedAt).toEqual({ gte: "2024-01-01" });
  });

  test("builds createdAt range", () => {
    const result = buildIssueFilter({
      "created-after": "2024-06-01",
      "created-before": "2024-06-30",
    });
    expect(result.createdAt).toEqual({ gte: "2024-06-01", lte: "2024-06-30" });
  });

  test("omits date filters when not provided", () => {
    const result = buildIssueFilter({ team: "ENG" });
    expect(result.updatedAt).toBeUndefined();
    expect(result.createdAt).toBeUndefined();
  });

  test("combines multiple filters", () => {
    const result = buildIssueFilter({
      team: "ENG",
      status: "Done",
      priority: "urgent",
      "updated-after": "2024-01-01",
    });
    expect(result.team).toBeDefined();
    expect(result.state).toBeDefined();
    expect(result.priority).toEqual({ eq: 1 });
    expect(result.updatedAt).toEqual({ gte: "2024-01-01" });
  });
});
