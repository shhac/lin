/** Map a cycle node to a summary for list/get output */
export function mapCycleSummary(c: {
  id: string;
  number: number;
  name?: string;
  startsAt: Date;
  endsAt: Date;
}): Record<string, unknown> {
  return {
    id: c.id,
    number: c.number,
    name: c.name,
    startsAt: c.startsAt,
    endsAt: c.endsAt,
  };
}
