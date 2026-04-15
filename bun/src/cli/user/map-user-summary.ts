/** Map a user node to a summary for list/search output */
export function mapUserSummary(u: {
  id: string;
  name: string;
  email: string;
  displayName: string;
}): Record<string, unknown> {
  return {
    id: u.id,
    name: u.name,
    email: u.email,
    displayName: u.displayName,
  };
}
