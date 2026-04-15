export const PRIORITY_MAP: Record<string, number> = {
  none: 0,
  urgent: 1,
  high: 2,
  medium: 3,
  low: 4,
};

export const PRIORITY_VALUES = "none | urgent | high | medium | low";

/**
 * Resolve a priority string to its numeric value.
 * Returns undefined if the input is not a valid priority name.
 */
export function resolvePriority(input: string): number | undefined {
  return PRIORITY_MAP[input.toLowerCase()];
}
